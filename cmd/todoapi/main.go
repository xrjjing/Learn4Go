package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/xrjjing/Learn4Go/internal/todo"
)

// TODO API 服务入口
//
// 环境变量:
//   TODO_STORAGE:  存储类型 (memory | sqlite | mysql)，默认 memory
//   TODO_ADDR:     监听地址，默认 :8080
//
// SQLite 配置:
//   TODO_DB_PATH:  SQLite 数据库路径，默认 todos.db
//
// MySQL 配置 (Docker):
//   TODO_DB_HOST:  MySQL 主机，默认 localhost
//   TODO_DB_PORT:  MySQL 端口，默认 3306
//   TODO_DB_USER:  MySQL 用户，默认 root
//   TODO_DB_PASS:  MySQL 密码，默认 root
//   TODO_DB_NAME:  数据库名，默认 learn4go
//
// 运行示例:
//   go run ./cmd/todoapi                                    # 内存存储
//   TODO_STORAGE=sqlite go run ./cmd/todoapi                # SQLite
//   TODO_STORAGE=mysql TODO_DB_PASS=secret go run ./cmd/todoapi  # MySQL
//
// Docker MySQL:
//   docker run -d --name mysql -p 3306:3306 \
//     -e MYSQL_ROOT_PASSWORD=root \
//     -e MYSQL_DATABASE=learn4go \
//     mysql:8
func main() {
	// 读取配置
	storage := getEnv("TODO_STORAGE", "memory")
	addr := getEnv("TODO_ADDR", ":8080")
	jwtSecret := getEnv("JWT_SECRET", "")

	// 创建存储
	var store todo.TodoStore
	var err error

	switch storage {
	case "mysql":
		host := getEnv("TODO_DB_HOST", "localhost")
		port := getEnvInt("TODO_DB_PORT", 3306)
		user := getEnv("TODO_DB_USER", "root")
		pass := getEnv("TODO_DB_PASS", "root")
		dbName := getEnv("TODO_DB_NAME", "learn4go")

		store, err = todo.NewMySQLStore(host, port, user, pass, dbName)
		if err != nil {
			log.Fatalf("MySQL 连接失败: %v", err)
		}
		log.Printf("使用 MySQL 存储: %s@%s:%d/%s", user, host, port, dbName)

	case "sqlite":
		dbPath := getEnv("TODO_DB_PATH", "todos.db")
		store, err = todo.NewSQLiteStore(dbPath)
		if err != nil {
			log.Fatalf("SQLite 初始化失败: %v", err)
		}
		log.Printf("使用 SQLite 存储: %s", dbPath)

	default:
		store = todo.NewStore()
		log.Println("使用内存存储")
	}

// JWT 密钥配置检查
	if jwtSecret == "" {
		if storage == "memory" {
			jwtSecret = "dev-secret-for-memory-mode-only"
			log.Println("警告: 使用内置开发密钥，仅限内存模式测试")
		} else {
			log.Fatal("错误: 生产模式必须设置 JWT_SECRET 环境变量")
		}
	}

	// 启动服务
	s := todo.NewServer(store, todo.WithJWT(jwtSecret, 24*time.Hour))

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("TODO API 启动: http://localhost%s", addr)
		log.Println("API 端点 (v1):")
		log.Println("  POST   /v1/register    - 注册")
		log.Println("  POST   /v1/login       - 登录")
		log.Println("  POST   /v1/refresh     - 刷新令牌")
		log.Println("  GET    /v1/todos       - 列表")
		log.Println("  POST   /v1/todos       - 创建")
		log.Println("  PUT    /v1/todos/{id}  - 更新状态")
		log.Println("  DELETE /v1/todos/{id}  - 删除")
		log.Println("  GET    /healthz        - 健康检查")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅关闭：监听 SIGINT 和 SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到关闭信号，正在优雅关闭...")

	// 给予 30 秒超时时间完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP 服务器关闭错误: %v", err)
	}

	// 关闭业务层资源（清理协程等）
	s.Shutdown()

	// 关闭数据库连接
	if closer, ok := store.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			log.Printf("数据库关闭错误: %v", err)
		}
	}

	log.Println("服务已安全关闭")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
