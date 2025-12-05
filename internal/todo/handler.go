package todo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

	// 在认证之后，对 /todos* 路径统一应用 RBAC 授权中间件
	rbacHandler := s.authzMiddleware(base)
	rbacWrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/todos") {
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
			"endpoints": []string{"/todos", "/todos/{id}", "/healthz"},
		}, http.StatusOK)
	})

	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// 认证相关路由
	s.mux.HandleFunc("/register", s.handleRegister)
	s.mux.HandleFunc("/login", s.handleLogin)
	s.mux.HandleFunc("/refresh", s.handleRefresh)

	s.mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
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

	s.mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/todos/"):]
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
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	}, http.StatusOK)
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

// 简易日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
