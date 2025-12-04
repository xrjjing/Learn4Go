package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xrjjing/Learn4Go/internal/logzip"
)

// 将若干日志字符串写入 ZIP，并输出校验和示例
func main() {
	logs := []string{
		"INFO boot service at " + time.Now().Format(time.RFC3339),
		"WARN slow response in handler X",
		"ERROR sample error detail",
	}

	data, sum, err := logzip.BuildZip(logs, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	const out = "logs.zip"
	if err := os.WriteFile(out, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("已生成 %s, SHA256=%x, 大小=%d bytes\n", out, sum, len(data))
}
