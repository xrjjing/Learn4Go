package todo

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

// ctxUserIDKey 用于在上下文中存放用户 ID
type ctxUserIDKey struct{}

// JWTManager 管理令牌生成与校验。
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

// Generate 生成用户访问令牌。
func (m *JWTManager) Generate(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.TTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// Parse 验证令牌并返回用户 ID。
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

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserID 从上下文中获取用户ID
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(ctxUserIDKey{}).(uint)
	return userID, ok
}

// authMiddleware 保护 /todos* 路径。
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 公开路径
		if r.URL.Path == "/" || r.URL.Path == "/healthz" ||
			r.URL.Path == "/register" || r.URL.Path == "/login" || r.URL.Path == "/refresh" {
			next.ServeHTTP(w, r)
			return
		}

		// 需要认证的路径
		if !strings.HasPrefix(r.URL.Path, "/todos") {
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
