package logzip

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"time"
)

// BuildZip 将日志内容写入 zip，使用固定时间保证可测。
func BuildZip(logs []string, modTime time.Time) ([]byte, [32]byte, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for i, line := range logs {
		h := &zip.FileHeader{
			Name:     fmt.Sprintf("log_%d.txt", i+1),
			Method:   zip.Deflate,
			Modified: modTime,
		}
		f, err := zw.CreateHeader(h)
		if err != nil {
			return nil, [32]byte{}, err
		}
		if _, err := io.WriteString(f, line+"\n"); err != nil {
			return nil, [32]byte{}, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, [32]byte{}, err
	}
	sum := sha256.Sum256(buf.Bytes())
	return buf.Bytes(), sum, nil
}
