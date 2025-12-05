package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xrjjing/Learn4Go/tinygee"
)

func TestRBACReject(t *testing.T) {
	app := tinygee.New()
	// 模拟 JWT 中写入 role
	app.Use(func(c *tinygee.Context) { c.Params = map[string]string{"role": "user"}; c.Next() })
	app.Use(RBAC(RBACConfig{
		RolePermissions: map[string][]string{
			"admin": {"/api"},
			"user":  {"/api/public"},
		},
	}))
	app.GET("/api/secret", func(c *tinygee.Context) {
		c.String(http.StatusOK, "secret")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/secret", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}
