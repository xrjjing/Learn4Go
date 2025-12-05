package main

import (
	"fmt"
	"io"
	"os"
)

// 简单文件拷贝示例
func main() {
	src := "main.go"
	dst := "main_copy.tmp"

	in, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer func() {
		out.Close()
		os.Remove(dst)
	}()

	written, err := io.Copy(out, in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("copied %d bytes\n", written)
}
