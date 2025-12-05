package todo

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter 基于内存的滑动窗口限流（简化版，便于演示）
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*ClientInfo
	window  time.Duration
	limit   int
	done    chan struct{}
}

// ClientInfo 客户端信息
type ClientInfo struct {
	requests []time.Time
	lastSeen time.Time
}

// NewRateLimiter 创建限流器。
func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	limiter := &RateLimiter{
		clients: make(map[string]*ClientInfo),
		window:  window,
		limit:   limit,
		done:    make(chan struct{}),
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// Stop 停止限流器的后台清理协程
func (r *RateLimiter) Stop() {
	close(r.done)
}

// cleanup 定期清理过期的客户端信息
func (r *RateLimiter) cleanup() {
	ticker := time.NewTicker(r.window * 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			now := time.Now()
			for key, client := range r.clients {
				if now.Sub(client.lastSeen) > r.window*2 {
					delete(r.clients, key)
				}
			}
			r.mu.Unlock()
		case <-r.done:
			return
		}
	}
}

// Middleware 限流中间件。
func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := clientIP(req)
		allowed := r.Allow(req.Context(), key)
		if !allowed {
			respondError(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		next.ServeHTTP(w, req)
	})
}

// Allow 判断当前请求是否通过限流。
func (r *RateLimiter) Allow(ctx context.Context, key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	client, exists := r.clients[key]

	if !exists {
		client = &ClientInfo{
			requests: []time.Time{now},
			lastSeen: now,
		}
		r.clients[key] = client
		return true
	}

	// 清理过期的请求记录
	cutoff := now.Add(-r.window)
	validRequests := make([]time.Time, 0, len(client.requests))
	for _, reqTime := range client.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// 检查是否超过限制
	if len(validRequests) >= r.limit {
		return false
	}

	// 添加当前请求
	client.requests = append(validRequests, now)
	client.lastSeen = now
	return true
}

// clientIP 提取客户端 IP（优先 X-Forwarded-For）
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}
