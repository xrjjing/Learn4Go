package main

import (
	"context"
	"fmt"
	"time"
)

// 演示 context.WithTimeout 取消 goroutine，避免泄露。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		// 模拟外部调用
		time.Sleep(200 * time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("work finished")
	case <-ctx.Done():
		fmt.Println("canceled:", ctx.Err())
	}
}
