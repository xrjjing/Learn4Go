package main

import "fmt"

// 演示值接收者 vs 指针接收者对状态的影响。

type Counter struct {
	n int
}

// 值接收者：修改不会影响原对象
func (c Counter) AddValue(delta int) {
	c.n += delta
}

// 指针接收者：可以修改原对象
func (c *Counter) AddPtr(delta int) {
	c.n += delta
}

func main() {
	c := Counter{n: 1}
	c.AddValue(5)
	fmt.Println("after AddValue:", c.n) // 仍然是 1
	c.AddPtr(5)
	fmt.Println("after AddPtr:", c.n) // 变为 6
}
