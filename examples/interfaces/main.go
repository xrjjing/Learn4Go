// 接口与多态示例
// 对应章节: 07_接口_多态.md
package main

import (
	"fmt"
	"math"
)

// 定义接口
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 实现接口的结构体：圆形
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 实现接口的结构体：矩形
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 使用接口的函数
func PrintShapeInfo(s Shape) {
	fmt.Printf("  面积: %.2f, 周长: %.2f\n", s.Area(), s.Perimeter())
}

// 空接口 interface{} 可以接收任何类型
func PrintAny(v interface{}) {
	fmt.Printf("  类型: %T, 值: %v\n", v, v)
}

// Stringer 接口（类似 Java 的 toString）
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("Person{Name: %s, Age: %d}", p.Name, p.Age)
}

func main() {
	// 多态：不同类型实现同一接口
	fmt.Println("多态示例:")
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 3},
	}

	for _, s := range shapes {
		PrintShapeInfo(s)
	}

	// 类型断言
	fmt.Println("\n类型断言:")
	var s Shape = Circle{Radius: 3}

	// 断言成功
	if c, ok := s.(Circle); ok {
		fmt.Println("  是圆形，半径:", c.Radius)
	}

	// 断言失败
	if _, ok := s.(Rectangle); !ok {
		fmt.Println("  不是矩形")
	}

	// type switch
	fmt.Println("\ntype switch:")
	checkType := func(i interface{}) {
		switch v := i.(type) {
		case int:
			fmt.Println("  int:", v)
		case string:
			fmt.Println("  string:", v)
		case Circle:
			fmt.Println("  Circle with radius:", v.Radius)
		default:
			fmt.Printf("  未知类型: %T\n", v)
		}
	}

	checkType(42)
	checkType("hello")
	checkType(Circle{Radius: 2})
	checkType(3.14)

	// 空接口
	fmt.Println("\n空接口 interface{}:")
	PrintAny(123)
	PrintAny("Go语言")
	PrintAny([]int{1, 2, 3})

	// Stringer 接口
	fmt.Println("\nStringer 接口:")
	p := Person{Name: "张三", Age: 25}
	fmt.Println(" ", p) // 自动调用 String() 方法

	// 接口组合
	fmt.Println("\n接口可以组合:")
	type ReadWriter interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
	}
	fmt.Println("  ReadWriter 接口组合了 Read 和 Write 方法")
}
