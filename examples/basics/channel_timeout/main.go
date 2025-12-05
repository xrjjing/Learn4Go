package main

import (
	"fmt"
	"time"
)

// 演示 select + time.After 进行超时控制。
func main() {
	result := make(chan string)

	go func() {
		time.Sleep(150 * time.Millisecond)
		result <- "done"
	}()

	select {
	case v := <-result:
		fmt.Println("got:", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout")
	}
}
