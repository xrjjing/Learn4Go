package batchrename

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDryAndApply(t *testing.T) {
	dir := t.TempDir()
	// 准备文件
	files := []string{"a.txt", "b.log", "c.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("x"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	logger := log.New(os.Stdout, "", 0)
	cfg := Config{Dir: dir, Prefix: "new_", Suffix: ".txt", Apply: false, Logger: logger}
	if err := Run(cfg); err != nil {
		t.Fatalf("dry-run err: %v", err)
	}

	// dry-run 不应改名
	for _, name := range files {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected file still there: %s", name)
		}
	}

	// 实际改名
	cfg.Apply = true
	if err := Run(cfg); err != nil {
		t.Fatalf("apply err: %v", err)
	}

	// 仅匹配后缀的文件被改名
	checkExists(t, filepath.Join(dir, "new_a.txt"))
	checkExists(t, filepath.Join(dir, "new_c.txt"))
	if _, err := os.Stat(filepath.Join(dir, "b.log")); err != nil {
		t.Fatalf("non-matching file should remain: %v", err)
	}
}

func checkExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist, err=%v", path, err)
	}
	if !strings.Contains(filepath.Base(path), "new_") {
		t.Fatalf("name not updated: %s", path)
	}
}
