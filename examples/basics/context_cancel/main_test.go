package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestContextCancel(t *testing.T) {
	out := capture(func() { main() })
	if !(strings.Contains(out, "canceled:") || strings.Contains(out, "work finished")) {
		t.Fatalf("unexpected output: %q", out)
	}
}

// capture is reused to keep dependencies minimal.
func capture(fn func()) string {
	pr, pw, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = pw
	fn()
	pw.Close()
	os.Stdout = old
	buf, _ := io.ReadAll(pr)
	return string(buf)
}
