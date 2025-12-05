package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type User struct {
	Name string `json:"name" xml:"name"`
	Age  int    `json:"age" xml:"age"`
}

func main() {
	u := User{Name: "Go", Age: 10}

	j, _ := json.Marshal(u)
	fmt.Println("json:", string(j))

	x, _ := xml.Marshal(u)
	fmt.Println("xml:", string(x))
}
