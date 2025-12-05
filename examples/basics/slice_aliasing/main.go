package main

import "fmt"

// 演示切片共享底层数组导致的数据联动。
func main() {
	base := []int{1, 2, 3, 4}
	a := base[:2] // [1 2]
	b := base[1:] // [2 3 4]

	fmt.Println("before:", a, b, base)

	// 修改共享元素
	a[1] = 99
	fmt.Println("after a change:", a, b, base)

	// append 触发扩容（取决于容量）
	a = append(a, 5) // 可能触发新底层
	a[0] = 42
	fmt.Println("after append:", a, b, base)
}
