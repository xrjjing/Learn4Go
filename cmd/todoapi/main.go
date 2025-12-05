package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
	rand.Seed(time.Now().UnixNano())

	// 读取配置
	storage := getEnv("TODO_STORAGE", "memory")
	addr := getEnv("TODO_ADDR", ":8080")

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

	// 启动服务
	s := todo.NewServer(store)
	log.Printf("TODO API 启动: http://localhost%s", addr)
	log.Println("API 端点:")
	log.Println("  GET    /todos      - 列表")
	log.Println("  POST   /todos      - 创建")
	log.Println("  PUT    /todos/{id} - 更新状态")
	log.Println("  DELETE /todos/{id} - 删除")
	log.Println("  GET    /healthz    - 健康检查")

	if err := http.ListenAndServe(addr, s.Handler()); err != nil {
		log.Fatal(err)
	}
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
