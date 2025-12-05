package render

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"

	"github.com/xrjjing/Learn4Go/tinygee"
)

func TestTemplateRender(t *testing.T) {
	r := tinygee.New()
	tr, err := New("testdata/*.html", template.FuncMap{"upper": strings.ToUpper})
	if err != nil {
		t.Fatalf("load template: %v", err)
	}
	r.GET("/hello", func(c *tinygee.Context) {
		tr.HTML(c, http.StatusOK, "hello.html", map[string]string{"Name": "go"})
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, "GO") {
		t.Fatalf("want upper name, got %s", body)
	}
}

func TestStatic(t *testing.T) {
	r := tinygee.New()
	r.GET("/static/*filepath", Static("/static/", http.Dir("testdata/static")))

	req := httptest.NewRequest(http.MethodGet, "/static/readme.txt", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "static") {
		t.Fatalf("unexpected body %s", w.Body.String())
	}
}
