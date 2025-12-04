package main

import (
	"flag"
	"log"
	"os"

	"github.com/xrjjing/Learn4Go/internal/cli/batchrename"
)

// 简易批量重命名工具：默认 dry-run，需 --apply 才会真正重命名
func main() {
	dir := flag.String("dir", ".", "目标目录")
	prefix := flag.String("prefix", "new_", "新文件名前缀")
	suffix := flag.String("suffix", "", "仅处理指定后缀的文件，如 .txt")
	apply := flag.Bool("apply", false, "执行实际重命名，默认仅打印")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	cfg := batchrename.Config{
		Dir:    *dir,
		Prefix: *prefix,
		Suffix: *suffix,
		Apply:  *apply,
		Logger: logger,
	}
	if err := batchrename.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
