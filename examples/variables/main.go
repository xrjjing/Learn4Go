// 变量与类型示例
// 对应章节: 02_变量_常量_类型.md
package main

import "fmt"

func main() {
	// 变量声明方式
	var a int = 10        // 显式类型
	var b = 20            // 类型推断
	c := 30               // 短声明（仅函数内）
	var d, e int = 40, 50 // 多变量声明

	fmt.Println("变量声明:", a, b, c, d, e)

	// 常量
	const Pi = 3.14159
	const (
		StatusOK    = 200
		StatusError = 500
	)
	fmt.Println("常量:", Pi, StatusOK, StatusError)

	// 基本类型
	var (
		intVal    int     = 42
		floatVal  float64 = 3.14
		boolVal   bool    = true
		stringVal string  = "Hello Go"
	)
	fmt.Printf("类型: int=%d, float=%.2f, bool=%t, string=%s\n",
		intVal, floatVal, boolVal, stringVal)

	// 类型转换（Go 需要显式转换）
	var x int = 10
	var y float64 = float64(x)
	fmt.Println("类型转换: int->float64:", y)

	// 零值（未初始化变量的默认值）
	var (
		zeroInt    int
		zeroFloat  float64
		zeroBool   bool
		zeroString string
	)
	fmt.Printf("零值: int=%d, float=%f, bool=%t, string=%q\n",
		zeroInt, zeroFloat, zeroBool, zeroString)
}
