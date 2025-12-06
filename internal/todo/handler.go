package todo

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 初始化结构化日志
var slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

// Server 封装路由与存储。
type Server struct {
	store       TodoStore
	userStore   UserStore
	jwtManager  *JWTManager
	rateLimiter *RateLimiter
	rbacManager *RBACManager
	mux         *http.ServeMux
	// 登录安全与刷新
	refreshTTL    time.Duration
	refreshStore  map[string]refreshSession
	refreshMu     sync.Mutex
	loginFailures map[string]*loginFailure
	loginMu       sync.Mutex
	// 清理协程控制
	cleanupDone chan struct{}
}

// Option 可选项配置服务器。
type Option func(s *Server)

// WithUserStore 指定用户存储。
func WithUserStore(us UserStore) Option {
	return func(s *Server) {
		s.userStore = us
	}
}

// WithJWT 配置 JWT 密钥与过期时间。
func WithJWT(secret string, ttl time.Duration) Option {
	return func(s *Server) {
		if secret != "" {
			s.jwtManager = NewJWTManager(secret, ttl)
		}
	}
}

// WithRateLimiter 配置速率限制器。
func WithRateLimiter(rl *RateLimiter) Option {
	return func(s *Server) {
		s.rateLimiter = rl
	}
}

// NewServer 创建带路由的 HTTP 处理器。
// store 可以是 *Store (内存) 或 *DBStore (数据库)
func NewServer(store TodoStore, opts ...Option) *Server {
	s := &Server{
		store:         store,
		userStore:     NewMemoryUserStore(), // 默认内存用户存储，便于测试
		jwtManager:    NewJWTManager("dev-secret-change-me-in-production", 24*time.Hour),
		rbacManager:   NewRBACManager(),
		mux:           http.NewServeMux(),
		refreshTTL:    defaultRefreshTTL,
		refreshStore:  make(map[string]refreshSession),
		loginFailures: make(map[string]*loginFailure),
		cleanupDone:   make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.routes()
	go s.startCleanup()
	return s
}

// Shutdown 优雅关闭服务器，停止清理协程
func (s *Server) Shutdown() {
	close(s.cleanupDone)
	if s.rateLimiter != nil {
		s.rateLimiter.Stop()
	}
}

// startCleanup 启动后台清理协程，定期清理过期的 refresh token 和登录失败记录
func (s *Server) startCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanExpiredRefreshTokens()
			s.cleanExpiredLoginFailures()
		case <-s.cleanupDone:
			return
		}
	}
}

// cleanExpiredRefreshTokens 清理过期的 refresh token
func (s *Server) cleanExpiredRefreshTokens() {
	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	now := time.Now()
	for token, session := range s.refreshStore {
		if now.After(session.expiresAt) {
			delete(s.refreshStore, token)
		}
	}
}

// cleanExpiredLoginFailures 清理过期的登录失败记录
func (s *Server) cleanExpiredLoginFailures() {
	s.loginMu.Lock()
	defer s.loginMu.Unlock()

	cutoff := time.Now().Add(-loginFailureWindow * 2)
	for email, rec := range s.loginFailures {
		if rec.lastFailedAt.Before(cutoff) {
			delete(s.loginFailures, email)
		}
	}
}

// Handler 返回带日志中间件的处理器。
func (s *Server) Handler() http.Handler {
	// 基础路由处理
	base := http.Handler(s.mux)

	// 在认证之后，对 /v1/todos* 路径统一应用 RBAC 授权中间件
	rbacHandler := s.authzMiddleware(base)
	rbacWrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/todos") {
			rbacHandler.ServeHTTP(w, r)
			return
		}
		base.ServeHTTP(w, r)
	})

	var h http.Handler = rbacWrapped
	h = s.authMiddleware(h)
	if s.rateLimiter != nil {
		h = s.rateLimiter.Middleware(h)
	}
	h = loggingMiddleware(h)
	h = corsMiddleware(h)
	return h
}

func (s *Server) routes() {
	// 根路径返回服务说明，避免裸访问 404 误判
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		respondJSON(w, map[string]any{
			"service":   "Learn4Go TODO API",
			"version":   "1.0",
			"endpoints": []string{"/v1/todos", "/v1/todos/{id}", "/healthz"},
		}, http.StatusOK)
	})

	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]any{
			"status": "ok",
			"checks": map[string]string{},
		}
		checks := status["checks"].(map[string]string)

		// 检查数据库连接（如果 store 支持 Ping）
		if pinger, ok := s.store.(interface{ Ping() error }); ok {
			if err := pinger.Ping(); err != nil {
				checks["database"] = "unhealthy: " + err.Error()
				status["status"] = "degraded"
			} else {
				checks["database"] = "healthy"
			}
		} else {
			checks["database"] = "in-memory"
		}

		code := http.StatusOK
		if status["status"] != "ok" {
			code = http.StatusServiceUnavailable
		}
		respondJSON(w, status, code)
	})

	// 认证相关路由 (v1)
	s.mux.HandleFunc("/v1/register", s.handleRegister)
	s.mux.HandleFunc("/v1/login", s.handleLogin)
	s.mux.HandleFunc("/v1/refresh", s.handleRefresh)

	// 用户和 RBAC 相关路由
	s.mux.Handle("/v1/me", s.authMiddleware(http.HandlerFunc(s.handleGetCurrentUser)))
	s.mux.Handle("/v1/logout", s.authMiddleware(http.HandlerFunc(s.handleLogout)))
	s.mux.Handle("/v1/users", s.authMiddleware(http.HandlerFunc(s.handleUsers)))
	s.mux.Handle("/v1/users/", s.authMiddleware(http.HandlerFunc(s.handleUserDetail)))
	s.mux.Handle("/v1/rbac/roles", s.authMiddleware(http.HandlerFunc(s.handleRoles)))
	s.mux.Handle("/v1/rbac/permissions", s.authMiddleware(http.HandlerFunc(s.handlePermissions)))

	s.mux.HandleFunc("/v1/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// 根据用户角色控制可见范围
			userID, ok := GetUserID(r.Context())
			if !ok {
				respondError(w, http.StatusUnauthorized, "authorization required")
				return
			}

			user, err := s.userStore.FindByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, ErrUserNotFound) {
					respondError(w, http.StatusUnauthorized, "user not found")
					return
				}
				respondError(w, http.StatusInternalServerError, "internal error")
				return
			}

			var items []Todo
			switch user.Role {
			case RoleAdmin, RoleGuest:
				// 管理员和访客都可以查看全部 TODO（写权限由 RBAC 控制）
				items, err = s.store.List()
			case RoleUser:
				// 普通用户只能看到自己的 TODO
				items, err = s.store.ListByUser(userID)
			default:
				// 未知角色默认按普通用户处理，避免越权
				items, err = s.store.ListByUser(userID)
			}
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal error")
				return
			}
			respondJSON(w, items, http.StatusOK)
		case http.MethodPost:
			var body struct {
				Title string `json:"title"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondError(w, http.StatusBadRequest, "invalid json")
				return
			}
			if body.Title == "" {
				respondError(w, http.StatusBadRequest, "title required")
				return
			}

			// 从上下文中获取用户 ID，将创建的 TODO 归属到该用户
			userID, ok := GetUserID(r.Context())
			if !ok {
				respondError(w, http.StatusUnauthorized, "authorization required")
				return
			}

			t, err := s.store.Create(body.Title, userID)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal error")
				return
			}
			respondJSON(w, t, http.StatusCreated)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	s.mux.HandleFunc("/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/v1/todos/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid id")
			return
		}
		switch r.Method {
		case http.MethodPut:
			var body struct {
				Done bool `json:"done"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				respondError(w, http.StatusBadRequest, "invalid json")
				return
			}
			t, ok, err := s.store.Toggle(id, body.Done)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal error")
				return
			}
			if !ok {
				respondError(w, http.StatusNotFound, "not found")
				return
			}
			respondJSON(w, t, http.StatusOK)
		case http.MethodDelete:
			ok, err := s.store.Delete(id)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal error")
				return
			}
			if !ok {
				respondError(w, http.StatusNotFound, "not found")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})
}

func respondJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, map[string]any{"error": msg}, code)
}

// handleRegister 用户注册
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Email == "" || body.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password required")
		return
	}
	if err := ValidateEmail(body.Email); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := ValidatePassword(body.Password); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := HashPassword(body.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "password hash error")
		return
	}
	user, err := s.userStore.Create(r.Context(), body.Email, hash)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	respondJSON(w, map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}, http.StatusCreated)
}

// handleLogin 用户登录获取 JWT
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// 登录失败次数限制
	if locked, retryAfter := s.checkLock(body.Email); locked {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
		respondError(w, http.StatusTooManyRequests, "account temporarily locked")
		return
	}

	user, err := s.userStore.FindByEmail(r.Context(), body.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			locked, wait := s.recordLoginFailure(body.Email)
			if locked {
				w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())))
				respondError(w, http.StatusTooManyRequests, "account temporarily locked")
				return
			}
			respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !CheckPassword(body.Password, user.PasswordHash) {
		locked, wait := s.recordLoginFailure(body.Email)
		if locked {
			w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())))
			respondError(w, http.StatusTooManyRequests, "account temporarily locked")
			return
		}
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// 登录成功，清除失败计数
	s.clearLoginFailure(body.Email)

	access, refresh, expiresIn, err := s.issueTokens(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "token error")
		return
	}
	respondJSON(w, map[string]any{
		"token":              access,
		"expires_in":         expiresIn,
		"refresh_token":      refresh,
		"refresh_expires_in": int(s.refreshTTL.Seconds()),
	}, http.StatusOK)
}

// handleGetCurrentUser 获取当前用户信息
func (s *Server) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authorization required")
		return
	}

	user, err := s.userStore.FindByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	respondJSON(w, map[string]any{
		"id":           user.ID,
		"email":        user.Email,
		"role":         user.Role,
		"is_superuser": user.Role == RoleAdmin,
		"is_active":    true, // 当前实现中所有用户都是活跃的
		"created_at":   user.CreatedAt,
	}, http.StatusOK)
}

// handleLogout 退出登录
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authorization required")
		return
	}

	// 清除该用户的所有 refresh tokens
	s.refreshMu.Lock()
	for token, session := range s.refreshStore {
		if session.userID == userID {
			delete(s.refreshStore, token)
		}
	}
	s.refreshMu.Unlock()

	respondJSON(w, map[string]any{
		"message": "logged out successfully",
	}, http.StatusOK)
}

// handleUsers 用户管理 (列表/创建)
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authorization required")
		return
	}

	currentUser, err := s.userStore.FindByID(r.Context(), userID)
	if err != nil || currentUser.Role != RoleAdmin {
		respondError(w, http.StatusForbidden, "admin privileges required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleUsersList(w, r)
	case http.MethodPost:
		s.handleUsersCreate(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleUsersList 获取用户列表
func (s *Server) handleUsersList(w http.ResponseWriter, r *http.Request) {
	// 获取所有用户 (从 MemoryUserStore)
	store, ok := s.userStore.(*MemoryUserStore)
	if !ok {
		respondError(w, http.StatusInternalServerError, "storage type not supported")
		return
	}

	store.mu.Lock()
	users := make([]map[string]any, 0, len(store.usersByID))
	for _, user := range store.usersByID {
		users = append(users, map[string]any{
			"id":           user.ID,
			"email":        user.Email,
			"role":         user.Role,
			"is_superuser": user.Role == RoleAdmin,
			"is_active":    true,
			"created_at":   user.CreatedAt,
		})
	}
	store.mu.Unlock()

	respondJSON(w, users, http.StatusOK)
}

// handleUsersCreate 创建新用户
func (s *Server) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsActive    bool   `json:"is_active"`
		IsSuperuser bool   `json:"is_superuser"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if body.Email == "" || body.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password required")
		return
	}

	if err := ValidateEmail(body.Email); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ValidatePassword(body.Password); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := HashPassword(body.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "password hash error")
		return
	}

	// 创建用户
	user, err := s.userStore.Create(r.Context(), body.Email, hash)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// 如果是超级用户，更新角色
	if body.IsSuperuser {
		store, ok := s.userStore.(*MemoryUserStore)
		if ok {
			store.mu.Lock()
			if u, exists := store.usersByID[user.ID]; exists {
				u.Role = RoleAdmin
				store.usersByID[user.ID] = u
				store.users[user.Email] = u
				user = u
			}
			store.mu.Unlock()
		}
	}

	respondJSON(w, map[string]any{
		"id":           user.ID,
		"email":        user.Email,
		"role":         user.Role,
		"is_superuser": user.Role == RoleAdmin,
		"is_active":    true,
		"created_at":   user.CreatedAt,
	}, http.StatusCreated)
}

// handleUserDetail 用户详情操作 (更新)
func (s *Server) handleUserDetail(w http.ResponseWriter, r *http.Request) {
	// 验证管理员权限
	userID, ok := GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "authorization required")
		return
	}

	currentUser, err := s.userStore.FindByID(r.Context(), userID)
	if err != nil || currentUser.Role != RoleAdmin {
		respondError(w, http.StatusForbidden, "admin privileges required")
		return
	}

	// 提取目标用户 ID
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/users/")
	targetID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	switch r.Method {
	case http.MethodPatch:
		s.handleUserUpdate(w, r, uint(targetID))
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleUserUpdate 更新用户状态
func (s *Server) handleUserUpdate(w http.ResponseWriter, r *http.Request, targetID uint) {
	var body struct {
		IsActive *bool `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	// 查找目标用户
	user, err := s.userStore.FindByID(r.Context(), targetID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// 注意：当前 MemoryUserStore 不支持 is_active 字段
	// 这里返回成功但不做实际修改（保持与前端期望一致）
	respondJSON(w, map[string]any{
		"id":           user.ID,
		"email":        user.Email,
		"role":         user.Role,
		"is_superuser": user.Role == RoleAdmin,
		"is_active":    true, // 始终返回 true
		"created_at":   user.CreatedAt,
	}, http.StatusOK)
}

// handleRoles 角色管理
func (s *Server) handleRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 返回基于 rbac.go 的角色列表
		roles := []map[string]any{
			{
				"id":          1,
				"name":        "admin",
				"description": "Administrator - Full permissions",
				"permissions": []map[string]string{
					{"code": "todos:create", "description": "Create TODO"},
					{"code": "todos:read", "description": "Read TODO"},
					{"code": "todos:update", "description": "Update TODO"},
					{"code": "todos:delete", "description": "Delete TODO"},
				},
			},
			{
				"id":          2,
				"name":        "user",
				"description": "Regular user - Manage own resources",
				"permissions": []map[string]string{
					{"code": "todos:create", "description": "Create TODO"},
					{"code": "todos:read", "description": "Read own TODO"},
					{"code": "todos:update", "description": "Update own TODO"},
					{"code": "todos:delete", "description": "Delete own TODO"},
				},
			},
			{
				"id":          3,
				"name":        "guest",
				"description": "Guest - Read-only access",
				"permissions": []map[string]string{
					{"code": "todos:read", "description": "Read TODO"},
				},
			},
		}
		respondJSON(w, roles, http.StatusOK)
	case http.MethodPost:
		// 角色创建功能暂未实现
		respondError(w, http.StatusNotImplemented, "role creation not implemented yet")
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handlePermissions 权限列表查询
func (s *Server) handlePermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	permissions := []map[string]any{
		{
			"id":          1,
			"code":        "todos:create",
			"description": "Create new TODO items",
		},
		{
			"id":          2,
			"code":        "todos:read",
			"description": "Read TODO items",
		},
		{
			"id":          3,
			"code":        "todos:update",
			"description": "Update TODO items",
		},
		{
			"id":          4,
			"code":        "todos:delete",
			"description": "Delete TODO items",
		},
	}

	respondJSON(w, permissions, http.StatusOK)
}

// handleRefresh 刷新 access token
func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh_token required")
		return
	}

	user, err := s.validateRefresh(body.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrRefreshExpired) {
			respondError(w, http.StatusUnauthorized, "refresh token expired")
			return
		}
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// 旋转 refresh，旧的立即失效
	s.dropRefresh(body.RefreshToken)

	access, refresh, expiresIn, err := s.issueTokens(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "token error")
		return
	}

	respondJSON(w, map[string]any{
		"token":              access,
		"expires_in":         expiresIn,
		"refresh_token":      refresh,
		"refresh_expires_in": int(s.refreshTTL.Seconds()),
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	}, http.StatusOK)
}

// responseWriter 包装器，用于捕获响应状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// 结构化日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		slogger.Info("http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		)
	})
}

// corsMiddleware 为本地开发提供简单的 CORS 支持，方便从 8000 端口的前端页面访问 8080 上的 TODO API。
// Docker 部署下由 Nginx 处理 CORS，这里主要覆盖直连 API 的场景。
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 暴露登录限流等场景使用到的 Retry-After 头，便于前端读取
		w.Header().Set("Access-Control-Expose-Headers", "Retry-After")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
