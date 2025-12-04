// 测试与基准示例 - 测试代码
// 对应章节: 10_测试与基准.md
//
// 运行测试: go test -v ./examples/testing
// 运行基准: go test -bench=. ./examples/testing
// 覆盖率:   go test -cover ./examples/testing
package main

import (
	"fmt"
	"testing"
)

// ========== 基本测试 ==========

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(5, 3)
	if result != 2 {
		t.Errorf("Subtract(5, 3) = %d; want 2", result)
	}
}

// ========== 表驱动测试 ==========

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正数乘正数", 2, 3, 6},
		{"正数乘零", 5, 0, 0},
		{"负数乘正数", -2, 3, -6},
		{"负数乘负数", -2, -3, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; want %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正常除法", 10, 2, 5},
		{"除以零", 10, 0, 0},
		{"负数除法", -10, 2, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Divide(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Divide(%d, %d) = %d; want %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// ========== 基准测试 ==========

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkFibonacciIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FibonacciIterative(20)
	}
}

// ========== 示例测试 ==========

func ExampleAdd() {
	result := Add(1, 2)
	fmt.Println(result)
	// Output: 3
}
