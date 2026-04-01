package todo

// 本文件聚焦“认证基础能力”，不直接做业务路由分发。
//
// 它提供三块最核心的能力：
// 1. JWT 的生成与解析
// 2. 密码的 bcrypt 加密与校验
// 3. 鉴权中间件，把 Bearer Token 里的 userID 写回请求上下文
//
// 页面排查建议：
// - 登录成功但后续接口 401：优先看 authMiddleware / Parse
// - token 看起来正确却解析失败：看 Generate/Parse 是否使用了同一套密钥与 TTL
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ctxUserIDKey 是私有上下文 key，避免和其他 context key 冲突。
// Handler 层通过 GetUserID 读取它，从而把“认证结果”传给后续业务处理。
type ctxUserIDKey struct{}

// JWTManager 是访问令牌的统一入口。
// 上游由 handleLogin / handleRefresh 调用 Generate，下游由 authMiddleware 调用 Parse。
type JWTManager struct {
	cfg JWTConfig
}

// JWTConfig 配置密钥与过期时间。
type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

// NewJWTManager 创建 JWT 管理器。
func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		cfg: JWTConfig{
			Secret: secret,
			TTL:    ttl,
		},
	}
}

// Generate 只负责生成 access token。refresh token 的生成与轮换逻辑在 security.go。
// Generate：登录成功后由 security.go 中的 issueTokens() 调用，用来签发 access token。
func (m *JWTManager) Generate(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.TTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// Parse 校验签名和 claims，并把 subject 还原成 uint userID。
// Parse：认证中间件每次放行受保护接口前都会走到这里。
func (m *JWTManager) Parse(tokenStr string) (uint, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok || !t.Valid {
		return 0, errors.New("invalid token claims")
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// HashPassword 在注册和后台创建用户时使用，避免明文密码进入存储层。
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 用于登录时比对用户输入和存储的 hash。
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserID 是 handler/rbac 层读取当前登录用户的统一入口。
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(ctxUserIDKey{}).(uint)
	return userID, ok
}

// authMiddleware 负责把“登录态”转换为“业务上下文”。
//
// 链路位置：
// 前端 Authorization 头 → authMiddleware → JWTManager.Parse → context 写入 userID → 业务 handler / RBAC 继续处理
//
// 如果 /v1/me、/v1/todos 等接口返回 401，优先看这里。
// authMiddleware：统一保护受限接口，并把认证结果写进上下文供后续 handler/RBAC 读取。
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 公开路径（无需认证）
		switch r.URL.Path {
		case "/", "/healthz", "/v1/register", "/v1/login", "/v1/refresh":
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, http.StatusUnauthorized, "authorization required")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := s.jwtManager.Parse(token)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
