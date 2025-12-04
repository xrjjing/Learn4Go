// 函数与错误处理示例
// 对应章节: 04_函数与错误处理.md
package main

import (
	"errors"
	"fmt"
)

// 基本函数
func add(a, b int) int {
	return a + b
}

// 多返回值
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// 命名返回值
func swap(a, b int) (x, y int) {
	x = b
	y = a
	return // 裸返回
}

// 可变参数
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 函数作为参数
func apply(fn func(int) int, val int) int {
	return fn(val)
}

// 闭包
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// defer 示例
func deferDemo() {
	fmt.Println("开始")
	defer fmt.Println("defer 1") // 后进先出
	defer fmt.Println("defer 2")
	fmt.Println("结束")
}

// panic 和 recover
func safeCall() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获 panic:", r)
		}
	}()
	panic("出错了！")
}

func main() {
	// 基本调用
	fmt.Println("add(3, 5) =", add(3, 5))

	// 多返回值与错误处理
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("10/2 =", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("错误:", err)
	}

	// 命名返回值
	x, y := swap(1, 2)
	fmt.Println("swap(1,2) =", x, y)

	// 可变参数
	fmt.Println("sum(1,2,3,4,5) =", sum(1, 2, 3, 4, 5))

	// 函数作为参数
	double := func(n int) int { return n * 2 }
	fmt.Println("apply(double, 5) =", apply(double, 5))

	// 闭包
	c := counter()
	fmt.Println("counter:", c(), c(), c())

	// defer
	deferDemo()

	// panic/recover
	safeCall()
	fmt.Println("程序继续执行")
}
