package main

import "fmt"

// 演示 nil map 写入会 panic，以及正确初始化。
func main() {
	var m map[string]int // nil map
	fmt.Println("nil map len:", len(m))
	// 写入会 panic，如果解除注释会崩溃
	// m["a"] = 1

	// 正确做法
	m = make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	if v, ok := m["c"]; !ok {
		fmt.Println("c not found, ok:", ok, "value:", v)
	}
	fmt.Println("map:", m)
}
