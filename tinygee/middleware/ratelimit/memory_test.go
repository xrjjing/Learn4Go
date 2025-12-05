package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
)

func TestRateLimit(t *testing.T) {
	app := tinygee.New()
	limiter := New(200*time.Millisecond, 2)
	app.Use(limiter.Middleware())
	app.GET("/ping", func(c *tinygee.Context) { c.String(http.StatusOK, "pong") })

	makeReq := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		return w
	}

	if w := makeReq(); w.Code != http.StatusOK {
		t.Fatalf("req1 status %d", w.Code)
	}
	if w := makeReq(); w.Code != http.StatusOK {
		t.Fatalf("req2 status %d", w.Code)
	}
	if w := makeReq(); w.Code != http.StatusTooManyRequests {
		t.Fatalf("req3 expect 429 got %d", w.Code)
	}
	time.Sleep(250 * time.Millisecond)
	if w := makeReq(); w.Code != http.StatusOK {
		t.Fatalf("after window status %d", w.Code)
	}
}
