package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xrjjing/Learn4Go/tinygee"
)

// JWTConfig 配置
type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

// Claims 自定义声明
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTMiddleware 返回验证中间件。
func NewJWTMiddleware(cfg JWTConfig) tinygee.HandlerFunc {
	return func(c *tinygee.Context) {
		authHeader := c.Req.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization required"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return []byte(cfg.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			return
		}
		// 将用户信息放入 Context
		if c.Params == nil {
			c.Params = map[string]string{}
		}
		c.Params["uid"] = claims.Subject
		c.Params["role"] = claims.Role
		c.Next()
	}
}

// GenerateToken 签发 token（示例用途）
func GenerateToken(cfg JWTConfig, uid uint, role string) (string, error) {
	claims := Claims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(uid), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.TTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}
