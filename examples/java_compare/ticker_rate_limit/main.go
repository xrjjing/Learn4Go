package main

import (
	"fmt"
	"time"
)

// 演示 time.Ticker 实现简单限速，对照 Java ScheduledExecutorService。
func main() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	tasks := []int{1, 2, 3, 4, 5}
	for _, id := range tasks {
		<-ticker.C // 每 200ms 才放行一次
		fmt.Printf("process task %d at %v\n", id, time.Now().Format("15:04:05.000"))
	}
	fmt.Println("done")
}
