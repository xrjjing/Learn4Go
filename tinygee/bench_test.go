package tinygee

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// 基准：路由匹配 + 中间件链开销
// 基线：map 路由（无 trie）
type mapRouter struct {
	m map[string]HandlerFunc
}

func newMapRouter() *mapRouter {
	return &mapRouter{m: make(map[string]HandlerFunc)}
}
func (r *mapRouter) add(path string, h HandlerFunc) { r.m[path] = h }
func (r *mapRouter) serve(c *Context) {
	if h, ok := r.m[c.Path]; ok {
		h(c)
		return
	}
	c.Status(http.StatusNotFound)
}

func BenchmarkRouterMatch(b *testing.B) {
	engine := New()
	for i := 0; i < 100; i++ {
		engine.GET("/api/v1/item/"+strconv.Itoa(i), func(c *Context) { c.String(http.StatusOK, "ok") })
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item/42", nil)
	w := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkRouterMatch_MapBaseline(b *testing.B) {
	r := newMapRouter()
	for i := 0; i < 100; i++ {
		path := "/api/v1/item/" + strconv.Itoa(i)
		r.add(path, func(c *Context) { c.String(http.StatusOK, "ok") })
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/item/42", nil)
	w := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := NewContext(w, req)
		r.serve(ctx)
	}
}

func BenchmarkMiddlewareChain(b *testing.B) {
	engine := New()
	engine.Use(func(c *Context) { c.Next() }, func(c *Context) { c.Next() }, func(c *Context) { c.Next() })
	engine.GET("/ping", func(c *Context) { c.String(http.StatusOK, "pong") })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ServeHTTP(w, req)
	}
}
