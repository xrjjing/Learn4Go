package logzip

import (
	"archive/zip"
	"bytes"
	"testing"
	"time"
)

func TestBuildZipDeterministic(t *testing.T) {
	logs := []string{"a", "b"}
	mod := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	data, sum, err := BuildZip(logs, mod)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("no data")
	}

	// 再次生成应得到相同哈希
	data2, sum2, err := BuildZip(logs, mod)
	if err != nil {
		t.Fatalf("build2: %v", err)
	}
	if sum != sum2 {
		t.Fatalf("hash mismatch")
	}
	if !bytes.Equal(data, data2) {
		t.Fatalf("data mismatch")
	}

	// 校验 zip 内容
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read: %v", err)
	}
	if len(r.File) != len(logs) {
		t.Fatalf("file count mismatch")
	}
}
