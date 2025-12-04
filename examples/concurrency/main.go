// 并发: goroutine 与 channel 示例
// 对应章节: 08_并发_goroutine_channel.md
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ========== Goroutine ==========
	fmt.Println("=== Goroutine ===")

	// 启动 goroutine
	go func() {
		fmt.Println("  goroutine 执行中...")
	}()

	// 等待 goroutine 执行
	time.Sleep(100 * time.Millisecond)

	// ========== Channel ==========
	fmt.Println("\n=== Channel 基础 ===")

	// 无缓冲 channel（同步）
	ch := make(chan int)
	go func() {
		ch <- 42 // 发送
	}()
	val := <-ch // 接收（阻塞直到有数据）
	fmt.Println("  收到:", val)

	// 有缓冲 channel（异步）
	buffered := make(chan int, 2)
	buffered <- 1
	buffered <- 2
	fmt.Println("  缓冲:", <-buffered, <-buffered)

	// ========== Select ==========
	fmt.Println("\n=== Select ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "来自 ch1"
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "来自 ch2"
	}()

	// select 等待多个 channel
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("  ", msg)
		case msg := <-ch2:
			fmt.Println("  ", msg)
		}
	}

	// ========== 超时控制 ==========
	fmt.Println("\n=== 超时控制 ===")

	timeout := make(chan int)
	go func() {
		time.Sleep(200 * time.Millisecond)
		timeout <- 1
	}()

	select {
	case <-timeout:
		fmt.Println("  收到数据")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  超时了！")
	}

	// ========== WaitGroup ==========
	fmt.Println("\n=== WaitGroup ===")

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  Worker %d 完成\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("  所有 worker 完成")

	// ========== Mutex ==========
	fmt.Println("\n=== Mutex ===")

	var (
		counter int
		mu      sync.Mutex
	)

	for i := 0; i < 100; i++ {
		go func() {
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("  计数器:", counter)

	// ========== 关闭 Channel ==========
	fmt.Println("\n=== 关闭 Channel ===")

	jobs := make(chan int, 3)
	jobs <- 1
	jobs <- 2
	jobs <- 3
	close(jobs) // 关闭 channel

	// range 遍历 channel 直到关闭
	for job := range jobs {
		fmt.Println("  处理任务:", job)
	}

	// 检查 channel 是否关闭
	_, ok := <-jobs
	fmt.Println("  channel 已关闭:", !ok)
}
