package todo

// 本文件实现一个教学用途的内存限流器。
//
// 在 TODO API 中它属于可选中间件：
// - 若 NewServer 注入 WithRateLimiter，则会进入 Handler() 的中间件链
// - 若未注入，则业务接口不会经过这一层
//
// 如果接口突然大量返回 429，可从这里确认窗口大小、限额和 clientIP 的取值。
import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter 按客户端维度记录时间窗内的请求时间列表。
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
// NewRateLimiter：创建限流器并启动后台清理协程。
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
// cleanup：定期清理长时间未访问的客户端状态，避免内存无限增长。
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

// Middleware 把 Allow 的判断结果翻译成 HTTP 429。
// Middleware：把限流判断封装成标准 http.Handler 中间件。
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

// Allow 是核心判定逻辑：清理窗口外记录，判断当前请求是否超限。
// Allow：滑动窗口核心判断逻辑。
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

// clientIP 优先取代理透传头，便于网关或反向代理场景下识别真实来源。
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return r.RemoteAddr
}
