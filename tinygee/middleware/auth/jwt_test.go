package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
)

func TestJWTMiddleware(t *testing.T) {
	cfg := JWTConfig{Secret: "secret", TTL: time.Second}
	token, err := GenerateToken(cfg, 1, "admin")
	if err != nil {
		t.Fatalf("gen token: %v", err)
	}

	engine := tinygee.New()
	engine.Use(NewJWTMiddleware(cfg))
	engine.GET("/ping", func(c *tinygee.Context) {
		c.String(http.StatusOK, "ok")
	})

	// valid
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", w.Code)
	}

	// expired
	time.Sleep(time.Second)
	req2 := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	engine.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("want 401 expired, got %d", w2.Code)
	}
}
