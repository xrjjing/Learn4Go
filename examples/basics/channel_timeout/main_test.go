package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestChannelTimeout(t *testing.T) {
	out := capture(func() {
		main()
	})
	if !(strings.Contains(out, "got:") || strings.Contains(out, "timeout")) {
		t.Fatalf("expected output got/timeout, got %q", out)
	}
}

// capture standard output for tests
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
