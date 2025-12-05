package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// 演示一键启用 pprof，访问 http://localhost:6060/debug/pprof/
// 对照 Java Flight Recorder/VisualVM。
func main() {
	// 模拟一个工作负载
	go func() {
		sum := 0
		for i := 0; i < 1e7; i++ {
			sum += i
		}
		fmt.Println("work done", sum)
	}()

	addr := "localhost:6060"
	log.Println("pprof listening on", addr)
	// 使用默认 mux，已注册 pprof 路由
	log.Fatal(http.ListenAndServe(addr, nil))
	// 启动后可用：go tool pprof http://localhost:6060/debug/pprof/profile
	// 或浏览器打开 /debug/pprof/
	time.Sleep(time.Second)
}
