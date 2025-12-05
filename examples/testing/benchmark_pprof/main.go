package main

import (
	"bytes"
	"strings"
)

// 提供一个可跑 benchmark 的函数
func JoinRepeat(n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteString("hello")
	}
	return buf.String()
}

// Dummy 使用，避免 "no non-test Go files" 警告
var _ = strings.Builder{}
