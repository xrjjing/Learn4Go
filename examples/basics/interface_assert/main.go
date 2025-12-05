package main

import "fmt"

type Walker interface{ Walk() }

type Person struct{}

func (Person) Walk() {}

func main() {
	var w Walker = Person{}

	// 类型断言安全写法
	if p, ok := w.(Person); ok {
		fmt.Println("assert ok", p)
	} else {
		fmt.Println("assert fail")
	}

	switch v := w.(type) {
	case Person:
		fmt.Println("type switch Person", v)
	default:
		fmt.Println("unknown type")
	}
}
