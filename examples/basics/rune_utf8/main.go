package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "Go语言"
	fmt.Println("bytes:", []byte(s))
	fmt.Println("runes:", []rune(s))
	fmt.Println("len bytes:", len(s))
	fmt.Println("rune count:", utf8.RuneCountInString(s))
}
