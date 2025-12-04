// 测试与基准示例 - 主程序入口
// 对应章节: 10_测试与基准.md
package main

import "fmt"

func main() {
	fmt.Println("=== Go 测试示例 ===")

	// 演示被测函数
	fmt.Println("Add(2, 3) =", Add(2, 3))
	fmt.Println("Subtract(5, 3) =", Subtract(5, 3))
	fmt.Println("Multiply(4, 5) =", Multiply(4, 5))
	fmt.Println("Divide(10, 2) =", Divide(10, 2))
	fmt.Println("Divide(10, 0) =", Divide(10, 0))

	fmt.Println("\nFibonacci(10) =", Fibonacci(10))
	fmt.Println("FibonacciIterative(10) =", FibonacciIterative(10))

	fmt.Println("\n--- 测试命令 ---")
	fmt.Println("运行测试:     go test -v ./examples/testing")
	fmt.Println("运行基准测试: go test -bench=. ./examples/testing")
	fmt.Println("查看覆盖率:   go test -cover ./examples/testing")
	fmt.Println("生成覆盖报告: go test -coverprofile=coverage.out ./examples/testing")
	fmt.Println("              go tool cover -html=coverage.out")
}
