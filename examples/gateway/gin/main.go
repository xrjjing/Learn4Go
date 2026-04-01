// Gin 框架 API 网关示例
//
// 本示例使用 Gin 框架实现 API 网关，展示了：
//   - Gin 中间件的使用方式
//   - 路由分组（Router Group）
//   - 反向代理集成
//   - 优雅的错误处理
//
// 与标准库版本的对比：
//   - Gin 提供更简洁的 API
//   - 内置的中间件管理
//   - 更强大的路由功能
//   - 更好的性能（基于 httprouter）
//
// 运行方式：
//  1. 安装 Gin: go get -u github.com/gin-gonic/gin
//  2. 启动后端: go run ./cmd/todoapi
//  3. 启动网关: go run ./examples/gateway/gin
//  4. 访问: curl http://localhost:8888/api/v1/todos
package main

// Package main 演示基于 Gin 的 API 网关实现。
//
// 核心职责：
// 1. 对外暴露统一入口 `/api/*`
// 2. 通过中间件做 CORS、日志、认证占位、限流
// 3. 将 `/api/v1/*` 代理到真实后端 `/v1/*`
//
// 当前最常见的实际链路是：前端 -> Gin Gateway(:8888) -> TODO API(:8080)。
// 如果页面通过网关访问失败，优先看 setupRouter 和 createReverseProxy。
import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// BackendConfig 描述单个后端服务的名称、目标地址和超时。
type BackendConfig struct {
	Name    string // 服务名称
	Target  string // 服务地址
	Timeout time.Duration
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// backends 定义网关可代理的服务清单。
var backends = map[string]BackendConfig{
	"todos": {Name: "TODO API", Target: getEnv("TODO_API_URL", "http://localhost:8080"), Timeout: 10 * time.Second},
	"users": {Name: "User API", Target: getEnv("USER_API_URL", "http://localhost:8081"), Timeout: 10 * time.Second},
}

// LoggerMiddleware Gin 日志中间件
// 记录请求详情、耗时和状态码
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 获取状态码
		status := c.Writer.Status()

		// 根据状态码选择日志级别
		if status >= 500 {
			log.Printf("[ERROR] %s %s - %d (%v)", c.Request.Method, path, status, latency)
		} else if status >= 400 {
			log.Printf("[WARN] %s %s - %d (%v)", c.Request.Method, path, status, latency)
		} else {
			log.Printf("[INFO] %s %s - %d (%v)", c.Request.Method, path, status, latency)
		}
	}
}

// AuthMiddleware 认证中间件
// 检查 Authorization 头部
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要认证的路径
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			// 开发模式：记录警告但不阻止请求
			log.Printf("[WARN] 无认证请求: %s %s from %s",
				c.Request.Method, c.Request.URL.Path, c.ClientIP())
			// 生产环境取消注释以下代码：
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			//     "error": "missing authorization token",
			// })
			// return
		}

		// 将认证信息存入上下文，供后续处理器使用
		c.Set("auth_token", token)
		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 简单的限流中间件
// 使用令牌桶算法限制请求频率
func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	// 简化实现：使用 channel 作为令牌桶
	tokens := make(chan struct{}, requestsPerSecond)

	// 定期补充令牌
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer ticker.Stop()
		for range ticker.C {
			select {
			case tokens <- struct{}{}:
			default: // 桶满了，丢弃令牌
			}
		}
	}()

	// 初始填满令牌桶
	for i := 0; i < requestsPerSecond; i++ {
		tokens <- struct{}{}
	}

	return func(c *gin.Context) {
		select {
		case <-tokens:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
		}
	}
}

// createReverseProxy 是最关键的桥接点：它把前端看到的 `/api/...` 改写为后端真正识别的路径。
// createReverseProxy：真正的转发核心，会把 /api 前缀剥离后再送往下游。
func createReverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("解析后端地址失败: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Director 会在请求真正发到后端前改写 URL、Host 和透传头。
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// 剥离 /api 前缀，保留资源路径
		// 例如：/api/v1/todos -> /v1/todos, /api/v1/todos/123 -> /v1/todos/123
		if strings.HasPrefix(req.URL.Path, "/api/") {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")
		}
		req.Header.Set("X-Forwarded-By", "Learn4Go-Gin-Gateway")
		req.Host = targetURL.Host
	}

	// 错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ERROR] 代理失败: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"bad gateway"}`))
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// setupRouter 负责组装网关对外暴露的全部路由和中间件链。
// setupRouter：拼装根路由、健康检查和 /api 分组下的反向代理规则。
func setupRouter() *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New() // 不使用默认中间件

	// 全局中间件
	r.Use(gin.Recovery()) // 恢复 panic
	r.Use(CORSMiddleware())
	r.Use(LoggerMiddleware())

	// 根路径说明，避免裸访问 404
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Learn4Go API Gateway",
			"version": "1.0",
			"routes":  []string{"/health", "/api/v1/todos", "/api/v1/users"},
		})
	})

	// 健康检查（无需认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"gateway": "gin",
		})
	})

	// `/api` 组是网关的业务入口；前端通常访问的就是这一组路径。
	api := r.Group("/api")
	api.Use(AuthMiddleware())
	api.Use(RateLimitMiddleware(100)) // 每秒 100 请求
	{
		// TODO API 代理：/api/v1/todos -> /v1/todos
		todosProxy := createReverseProxy(backends["todos"].Target)
		api.Any("/v1/todos", todosProxy)
		api.Any("/v1/todos/*path", todosProxy)
		// TODO 健康检查代理：/api/healthz -> /healthz
		api.Any("/healthz", todosProxy)

		// Users API 代理（示例）：/api/v1/users -> /v1/users
		usersProxy := createReverseProxy(backends["users"].Target)
		api.Any("/v1/users", usersProxy)
		api.Any("/v1/users/*path", usersProxy)
	}

	return r
}

// main 只负责启动 Gin 网关进程，具体代理规则已经在 setupRouter 中定义。
// main：启动 Gin 网关进程。排查页面通过 8888 访问异常时，先看这里打印的后端映射。
func main() {
	r := setupRouter()

	addr := getEnv("GATEWAY_ADDR", ":8888")
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Gin API 网关启动，监听 %s", addr)
	log.Println("后端服务映射:")
	log.Printf("  /api/v1/todos  -> %s (%s)", backends["todos"].Target, backends["todos"].Name)
	log.Printf("  /api/v1/users  -> %s (%s)", backends["users"].Target, backends["users"].Name)
	log.Println("\n测试命令:")
	log.Println("  curl http://localhost:8888/health")
	log.Println("  curl http://localhost:8888/api/v1/todos")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("服务器退出: %v", err)
	}
}
