package main

import "fmt"

// 演示 defer 参数捕获与 LIFO 顺序。
func deferValue() {
	for i := 0; i < 2; i++ {
		defer fmt.Println("value capture:", i) // i 按值复制
	}
}

func deferClosure() {
	for i := 0; i < 2; i++ {
		defer func() {
			fmt.Println("closure capture:", i) // 捕获同一个 i，最终执行顺序 LIFO：先 1 后 0
		}()
	}
}

func main() {
	deferValue()
	deferClosure()
}
