// 结构体、方法、组合示例
// 对应章节: 06_结构体_方法_组合.md
package main

import "fmt"

// 定义结构体
type Person struct {
	Name string
	Age  int
}

// 方法（值接收者）
func (p Person) Greet() string {
	return fmt.Sprintf("你好，我是%s，今年%d岁", p.Name, p.Age)
}

// 方法（指针接收者）- 可以修改结构体字段
func (p *Person) SetAge(age int) {
	p.Age = age
}

// 嵌入结构体（组合）
type Employee struct {
	Person   // 匿名嵌入，继承字段和方法
	Company  string
	Position string
}

// Employee 的方法
func (e Employee) Info() string {
	return fmt.Sprintf("%s 在 %s 担任 %s", e.Name, e.Company, e.Position)
}

// 带标签的结构体（常用于 JSON）
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

func main() {
	// 创建结构体
	p1 := Person{Name: "张三", Age: 25}
	p2 := Person{"李四", 30} // 按顺序
	p3 := new(Person)      // 返回指针
	p3.Name = "王五"
	p3.Age = 35

	fmt.Println("结构体创建:")
	fmt.Println("  p1:", p1)
	fmt.Println("  p2:", p2)
	fmt.Println("  p3:", p3)

	// 调用方法
	fmt.Println("\n方法调用:")
	fmt.Println("  ", p1.Greet())

	// 指针接收者方法
	p1.SetAge(26)
	fmt.Println("  修改年龄后:", p1)

	// 组合（继承）
	fmt.Println("\n结构体组合:")
	emp := Employee{
		Person:   Person{Name: "赵六", Age: 28},
		Company:  "Tech Inc",
		Position: "工程师",
	}
	fmt.Println("  ", emp.Info())
	fmt.Println("  直接访问:", emp.Name) // 可以直接访问嵌入字段

	// 结构体比较
	fmt.Println("\n结构体比较:")
	a := Person{Name: "test", Age: 20}
	b := Person{Name: "test", Age: 20}
	fmt.Println("  a == b:", a == b) // true（所有字段相等）

	// 匿名结构体
	fmt.Println("\n匿名结构体:")
	point := struct {
		X, Y int
	}{10, 20}
	fmt.Println("  point:", point)
}
