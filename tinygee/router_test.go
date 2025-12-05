package tinygee

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteParams(t *testing.T) {
	r := newRouter()
	r.addRoute(http.MethodGet, "/hello/:name", func(c *Context) {})

	n, params := r.getRoute(http.MethodGet, "/hello/go")
	if n == nil || n.pattern != "/hello/:name" {
		t.Fatalf("route not matched")
	}
	if params["name"] != "go" {
		t.Fatalf("param mismatch %v", params)
	}
}

func TestWildcardRoute(t *testing.T) {
	r := newRouter()
	r.addRoute(http.MethodGet, "/assets/*filepath", func(c *Context) {})

	n, params := r.getRoute(http.MethodGet, "/assets/img/logo.png")
	if n == nil || n.pattern != "/assets/*filepath" {
		t.Fatalf("wildcard not matched")
	}
	if params["filepath"] != "img/logo.png" {
		t.Fatalf("wildcard value mismatch %v", params)
	}
}

func TestGroupMiddlewareOrder(t *testing.T) {
	engine := New()
	api := engine.Group("/api")
	api.Use(func(c *Context) { c.Next() })
	api.GET("/ping", func(c *Context) { c.String(http.StatusOK, "pong") })

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}
}
