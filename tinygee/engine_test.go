package tinygee

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngineGET(t *testing.T) {
	engine := New()
	engine.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello %s", "tinygee")
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200 got %d", w.Code)
	}
	if body := w.Body.String(); body != "hello tinygee" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestEngineNotFound(t *testing.T) {
	engine := New()
	req := httptest.NewRequest(http.MethodGet, "/absent", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}
