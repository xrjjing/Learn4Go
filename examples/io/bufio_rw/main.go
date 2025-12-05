package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 演示 bufio 逐行读取与写入。
func main() {
	input := "Go makes it easy.\nGo is expressive.\n"
	r := strings.NewReader(input)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println("line:", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scan error:", err)
	}

	// 写入文件
	f, err := os.CreateTemp("", "bufio-demo")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	w := bufio.NewWriter(f)
	_, _ = w.WriteString("hello bufio\nsecond line\n")
	w.Flush()
}
