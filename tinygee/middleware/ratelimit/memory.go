package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
)

// 简易滑动窗口限流（内存版）
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func New(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (r *RateLimiter) Middleware() tinygee.HandlerFunc {
	return func(c *tinygee.Context) {
		key := clientIP(c.Req)
		if !r.allow(key) {
			c.JSON(http.StatusTooManyRequests, map[string]string{"error": "too many requests"})
			return
		}
		c.Next()
	}
}

func (r *RateLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	requests := r.clients[key]
	cutoff := now.Add(-r.window)
	keep := requests[:0]
	for _, ts := range requests {
		if ts.After(cutoff) {
			keep = append(keep, ts)
		}
	}
	if len(keep) >= r.limit {
		r.clients[key] = keep
		return false
	}
	keep = append(keep, now)
	r.clients[key] = keep
	return true
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}
