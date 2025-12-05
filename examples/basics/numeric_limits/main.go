package main

import (
	"fmt"
	"math"
)

func main() {
	var u8 uint8 = math.MaxUint8
	var i8 int8 = math.MaxInt8
	fmt.Println("uint8 max:", u8, "int8 max:", i8)

	// 溢出示例（仅打印，不实际溢出避免编译器报错）
	fmt.Println("int32 max:", math.MaxInt32, "int64 max:", math.MaxInt64)
	fmt.Printf("float32 max: %e, float64 max: %e\n", math.MaxFloat32, math.MaxFloat64)
}
