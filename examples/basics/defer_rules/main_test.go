package main

import (
	"io"
	"os"
	"testing"
)

func TestDeferOrderAndCapture(t *testing.T) {
	out := captureOutput(func() {
		deferValue()
	})
	expected := "value capture: 1\nvalue capture: 0\n"
	if out != expected {
		t.Fatalf("defer value capture order mismatch:\nwant %q\ngot  %q", expected, out)
	}

	out = captureOutput(func() {
		deferClosure()
	})
	// 同一 i 的闭包，LIFO 顺序输出 1 再 0（迭代时的 i 值）
	expected = "closure capture: 1\nclosure capture: 0\n"
	if out != expected {
		t.Fatalf("defer closure capture mismatch:\nwant %q\ngot  %q", expected, out)
	}
}

// 简易输出捕获，避免引入额外依赖。
func captureOutput(fn func()) string {
	pr, pw, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = pw
	fn()
	pw.Close()
	os.Stdout = stdout
	buf, _ := io.ReadAll(pr)
	return string(buf)
}
