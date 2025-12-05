package tinygee

import (
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 确保 expvar /metrics 路由可正常工作
func TestMetricsRoute(t *testing.T) {
	r := New()
	r.GET("/metrics", func(c *Context) {
		expvar.Handler().ServeHTTP(c.Writer, c.Req)
	})
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("metrics status %d", w.Code)
	}
}
