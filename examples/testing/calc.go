// 测试与基准示例 - 被测试的代码
// 对应章节: 10_测试与基准.md
package main

// Add 两数相加
func Add(a, b int) int {
	return a + b
}

// Subtract 两数相减
func Subtract(a, b int) int {
	return a - b
}

// Multiply 两数相乘
func Multiply(a, b int) int {
	return a * b
}

// Divide 两数相除，除数为0返回0
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

// Fibonacci 返回第n个斐波那契数
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// FibonacciIterative 迭代版斐波那契（更高效）
func FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
