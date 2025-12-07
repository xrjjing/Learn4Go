package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/xrjjing/Learn4Go/rabbitmq-demo/internal/logstore"
	"github.com/xrjjing/Learn4Go/rabbitmq-demo/internal/rabbit"
)

// AppConfig 汇总运行所需配置。
type AppConfig struct {
	Rabbit      rabbit.Config
	Port        string
	StaticDir   string
	DocsDir     string
	MockPath    string
	Mode        string
	LogCapacity int
}

func loadConfig() AppConfig {
	cfg := AppConfig{Rabbit: rabbit.DefaultConfig(), LogCapacity: 400}
	if v := os.Getenv("RABBITMQ_URL"); v != "" {
		cfg.Rabbit.URL = v
	}
	if v := os.Getenv("RABBITMQ_EXCHANGE"); v != "" {
		cfg.Rabbit.Exchange = v
	}
	if v := os.Getenv("RABBITMQ_DELAY_EXCHANGE"); v != "" {
		cfg.Rabbit.DelayExchange = v
	}
	if v := os.Getenv("RABBITMQ_DLX_EXCHANGE"); v != "" {
		cfg.Rabbit.DLXExchange = v
	}
	if v := os.Getenv("RABBITMQ_WORK_QUEUE"); v != "" {
		cfg.Rabbit.WorkQueue = v
	}
	if v := os.Getenv("RABBITMQ_DELAY_QUEUE"); v != "" {
		cfg.Rabbit.DelayQueue = v
	}
	if v := os.Getenv("RABBITMQ_DLX_QUEUE"); v != "" {
		cfg.Rabbit.DLXQueue = v
	}
	if v := os.Getenv("RABBITMQ_ROUTING_WORK"); v != "" {
		cfg.Rabbit.RoutingWork = v
	}
	if v := os.Getenv("RABBITMQ_ROUTING_DLX"); v != "" {
		cfg.Rabbit.RoutingDLX = v
	}
	if v := os.Getenv("RABBITMQ_PREFETCH"); v != "" {
		if n, err := parseInt(v); err == nil {
			cfg.Rabbit.Prefetch = n
		}
	}
	if v := os.Getenv("LOG_CAPACITY"); v != "" {
		if n, err := parseInt(v); err == nil {
			if n > 0 && n <= 10000 {
				cfg.LogCapacity = n
			} else {
				log.Printf("⚠️  LOG_CAPACITY 值无效 (%d)，使用默认值 400", n)
				cfg.LogCapacity = 400
			}
		} else {
			log.Printf("⚠️  LOG_CAPACITY 解析失败: %v，使用默认值 400", err)
			cfg.LogCapacity = 400
		}
	}
	cfg.Port = getenvDefault("PORT", "8088")

	root := os.Getenv("RABBITMQ_DEMO_ROOT")
	if root == "" {
		// 默认使用仓库相对路径
		cwd, _ := os.Getwd()
		root = filepath.Join(cwd, "rabbitmq-demo")
	}
	cfg.StaticDir = filepath.Join(root, "web", "rabbitmq-demo")
	cfg.DocsDir = filepath.Join(root, "docs")
	cfg.MockPath = filepath.Join(root, "mock", "messages.json")
	cfg.Mode = "real"
	return cfg
}

func parseInt(v string) (int, error) {
	var n int
	_, err := fmt.Sscanf(v, "%d", &n)
	return n, err
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func loadMockMessages(path string) ([]rabbit.DemoMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var msgs []rabbit.DemoMessage
	if err := json.Unmarshal(data, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := loadConfig()
	mockMessages, err := loadMockMessages(cfg.MockPath)
	if err != nil {
		log.Fatalf("读取 mock 数据失败: %v", err)
	}

	logStore := logstore.New(cfg.LogCapacity)

	var mq rabbit.MQ
	if os.Getenv("RABBITMQ_FAKE") == "1" {
		mq = rabbit.NewMock(cfg.Rabbit)
		log.Println("使用内存版 RabbitMQ Mock，未连接真实 MQ")
		cfg.Mode = "fake"
	} else {
		mq, err = rabbit.New(cfg.Rabbit)
		if err != nil {
			log.Fatalf("初始化 RabbitMQ 客户端失败: %v", err)
		}
	}
	defer mq.Close()

	if err := mq.DeclareTopology(ctx); err != nil {
		log.Fatalf("声明拓扑失败: %v", err)
	}

	mux, err := setupHTTP(ctx, cfg, mq, mockMessages, logStore)
	if err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: mux}

	log.Printf("RabbitMQ Demo 服务启动，端口 %s，静态目录 %s", cfg.Port, cfg.StaticDir)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP 服务异常: %v", err)
	}
}

// setupHTTP 负责启动消费者并返回 HTTP mux，便于测试复用。
func setupHTTP(ctx context.Context, cfg AppConfig, mq rabbit.MQ, mockMessages []rabbit.DemoMessage, logStore *logstore.Store) (*http.ServeMux, error) {
	// 启动工作队列消费者
	if err := mq.Consume(ctx, cfg.Rabbit.WorkQueue, func(msg rabbit.DemoMessage) error {
		logStore.Add(logstore.Entry{Time: time.Now(), Kind: "consume", ID: msg.ID, Type: msg.Type, Message: "工作队列消息已消费"})
		if msg.Type == "order.fail" {
			return errors.New("模拟业务失败")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// 启动死信队列消费者
	if err := mq.Consume(ctx, cfg.Rabbit.DLXQueue, func(msg rabbit.DemoMessage) error {
		logStore.Add(logstore.Entry{Time: time.Now(), Kind: "dlx", ID: msg.ID, Type: msg.Type, Message: "死信队列收到消息"})
		return nil
	}); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	})

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		workMsg, workCons, err1 := mq.Stats(r.Context(), cfg.Rabbit.WorkQueue)
		delayMsg, delayCons, err2 := mq.Stats(r.Context(), cfg.Rabbit.DelayQueue)
		dlxMsg, dlxCons, err3 := mq.Stats(r.Context(), cfg.Rabbit.DLXQueue)
		if err1 != nil || err2 != nil || err3 != nil {
			http.Error(w, "队列状态查询失败", http.StatusInternalServerError)
			return
		}
		respondJSON(w, map[string]interface{}{
			"mode":       cfg.Mode,
			"rabbit_url": cfg.Rabbit.URL,
			"queues": []map[string]interface{}{
				{"name": cfg.Rabbit.WorkQueue, "messages": workMsg, "consumers": workCons},
				{"name": cfg.Rabbit.DelayQueue, "messages": delayMsg, "consumers": delayCons},
				{"name": cfg.Rabbit.DLXQueue, "messages": dlxMsg, "consumers": dlxCons},
			},
		})
	})

	// 文档静态托管
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(cfg.DocsDir))))

	mux.HandleFunc("/api/mock", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, mockMessages)
	})

	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, logStore.List())
	})

	mux.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}
		var msg rabbit.DemoMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "请求格式错误", http.StatusBadRequest)
			return
		}
		if msg.ID == "" {
			msg.ID = uuid.NewString()
		}
		if msg.Type == "" {
			msg.Type = "order.created"
		}
		if err := mq.Publish(r.Context(), msg); err != nil {
			logStore.Add(logstore.Entry{Time: time.Now(), Kind: "error", ID: msg.ID, Type: msg.Type, Message: err.Error()})
			http.Error(w, "发布失败", http.StatusInternalServerError)
			return
		}
		logStore.Add(logstore.Entry{Time: time.Now(), Kind: "send", ID: msg.ID, Type: msg.Type, Message: "消息已发布"})
		respondJSON(w, map[string]string{"id": msg.ID})
	})

	mux.HandleFunc("/api/messages/batch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}
		count := 0
		for _, msg := range mockMessages {
			localMsg := msg
			if localMsg.ID == "" {
				localMsg.ID = uuid.NewString()
			}
			if err := mq.Publish(r.Context(), localMsg); err != nil {
				logStore.Add(logstore.Entry{Time: time.Now(), Kind: "error", ID: localMsg.ID, Type: localMsg.Type, Message: err.Error()})
				http.Error(w, "部分发布失败", http.StatusInternalServerError)
				return
			}
			logStore.Add(logstore.Entry{Time: time.Now(), Kind: "send", ID: localMsg.ID, Type: localMsg.Type, Message: "批量消息已发布"})
			count++
		}
		respondJSON(w, map[string]int{"published": count})
	})

	// 静态文件
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/", fs)
	return mux, nil
}

func respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
