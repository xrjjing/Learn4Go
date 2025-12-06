// 标准库 API 网关示例
//
// 本示例使用 Go 标准库的 httputil.ReverseProxy 实现一个简单的 API 网关，
// 展示了以下功能：
//   - 反向代理：将请求转发到后端服务
//   - 中间件模式：认证、日志等横切关注点
//   - 路由分发：根据路径前缀转发到不同后端
//
// 运行方式：
//  1. 先启动后端服务: go run ./cmd/todoapi
//  2. 启动网关: go run ./examples/gateway/stdlib
//  3. 通过网关访问: curl http://localhost:8888/api/v1/todos
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// 后端服务地址配置
// 通过路径前缀 /api/v1/* 转发到实际的 /v1/* 接口
var backends = map[string]string{
	"/api/v1/todos": "http://localhost:8080", // TODO API 服务
	"/api/v1/users": "http://localhost:8081", // 用户服务（示例）
}

// middleware 定义中间件函数类型
// 中间件是一个接收 http.Handler 并返回 http.Handler 的函数
// 这种模式允许我们像洋葱一样层层包装处理逻辑
type middleware func(http.Handler) http.Handler

// chain 将多个中间件串联起来
// 执行顺序：第一个中间件最外层，最后一个最内层
// 例如: chain(a, b, c)(handler) => a(b(c(handler)))
func chain(middlewares ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// loggingMiddleware 请求日志中间件
// 记录每个请求的方法、路径、耗时和状态码
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 使用自定义 ResponseWriter 来捕获状态码
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 调用下一个处理器
		next.ServeHTTP(rw, r)

		// 记录请求日志
		log.Printf("[%s] %s %s - %d (%v)",
			r.Method,
			r.Host,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// authMiddleware 简单的认证中间件
// 检查请求头中是否包含有效的 Authorization token
// 生产环境应该使用 JWT 验证等更安全的方式
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 跳过不需要认证的路径
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// 检查 Authorization 头
		auth := r.Header.Get("Authorization")
		if auth == "" {
			// 开发模式：允许无认证访问，但记录警告
			log.Printf("[WARN] 无认证请求: %s %s", r.Method, r.URL.Path)
			// 如需强制认证，取消下面注释：
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			// return
		}

		next.ServeHTTP(w, r)
	})
}

// corsMiddleware 跨域资源共享中间件
// 允许前端从不同域名访问 API
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// createProxy 创建反向代理处理器
// target: 后端服务地址，如 "http://localhost:8080"
func createProxy(target string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("解析后端地址失败: %v", err)
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 自定义 Director 函数，可以修改转发的请求
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 添加自定义请求头，标识请求来自网关
		req.Header.Set("X-Forwarded-By", "Learn4Go-Gateway")
		req.Header.Set("X-Real-IP", req.RemoteAddr)
	}

	// 自定义错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ERROR] 代理请求失败: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}

// router 简单的路由分发器
func router() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 健康检查端点
		if path == "/health" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}

		// 根据路径前缀找到对应的后端
		for prefix, backend := range backends {
			if strings.HasPrefix(path, prefix) {
				// 将 /api 前缀剥离，保持 /v1/... 路径与后端一致
				if strings.HasPrefix(path, "/api/") {
					r.URL.Path = strings.TrimPrefix(path, "/api")
				}

				// 创建代理并转发请求
				proxy := createProxy(backend)
				proxy.ServeHTTP(w, r)
				return
			}
		}

		// 未匹配到任何后端
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}

func main() {
	// 组装中间件链
	// 执行顺序: cors -> logging -> auth -> router
	handler := chain(
		corsMiddleware,
		loggingMiddleware,
		authMiddleware,
	)(router())

	// 配置服务器
	server := &http.Server{
		Addr:         ":8888",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("API 网关启动，监听 :8888")
	log.Println("后端服务映射:")
	for prefix, backend := range backends {
		log.Printf("  %s -> %s", prefix, backend)
	}
	log.Println("\n测试命令:")
	log.Println("  curl http://localhost:8888/health")
	log.Println("  curl http://localhost:8888/api/v1/todos")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("服务器退出: %v", err)
	}
}
