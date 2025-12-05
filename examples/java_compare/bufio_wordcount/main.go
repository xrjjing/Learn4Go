package main

import (
	"bufio"
	"fmt"
	"strings"
)

// 演示 bufio.Scanner 逐行读取并统计词频，对照 Java BufferedReader。
func main() {
	text := `Go makes it easy to build simple, reliable, and efficient software.
Go concurrency patterns are powerful.
Go error handling is explicit.`

	wordCount := map[string]int{}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		w := strings.Trim(strings.ToLower(scanner.Text()), ".,")
		if w == "" {
			continue
		}
		wordCount[w]++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for k, v := range wordCount {
		fmt.Printf("%s -> %d\n", k, v)
	}
}
