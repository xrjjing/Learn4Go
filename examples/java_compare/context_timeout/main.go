package main

import (
	"context"
	"fmt"
	"time"
)

// 演示 context.WithTimeout 取消协程，对照 Java Future.get(timeout)。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		done <- "work done"
	}()

	select {
	case msg := <-done:
		fmt.Println(msg)
	case <-ctx.Done():
		// 超时则走这里
		fmt.Println("timeout:", ctx.Err())
	}
}
