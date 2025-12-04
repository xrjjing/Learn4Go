package batchrename

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config 配置项。
type Config struct {
	Dir    string
	Prefix string
	Suffix string
	Apply  bool
	Logger *log.Logger
}

// Run 执行批量重命名。
func Run(cfg Config) error {
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	entries, err := os.ReadDir(cfg.Dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if cfg.Suffix != "" && !strings.HasSuffix(name, cfg.Suffix) {
			continue
		}
		newName := cfg.Prefix + name
		oldPath := filepath.Join(cfg.Dir, name)
		newPath := filepath.Join(cfg.Dir, newName)

		if _, err := os.Stat(newPath); err == nil {
			cfg.Logger.Printf("跳过: 目标已存在 %s", newPath)
			continue
		}

		if cfg.Apply {
			if err := os.Rename(oldPath, newPath); err != nil {
				cfg.Logger.Printf("重命名失败 %s -> %s: %v", oldPath, newPath, err)
				continue
			}
			cfg.Logger.Printf("已重命名: %s -> %s", name, newName)
		} else {
			cfg.Logger.Printf("(dry-run) %s -> %s", name, newName)
		}
	}
	return nil
}

// ListFiles 列出目录文件（测试/演示用）。
func ListFiles(dir string) ([]fs.DirEntry, error) {
	return os.ReadDir(dir)
}
