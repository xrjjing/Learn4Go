package middleware

import (
	"log"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
)

// Logger 记录请求方法、路径、耗时与状态码。
func Logger() tinygee.HandlerFunc {
	return func(c *tinygee.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s -> %d (%v)", c.Method, c.Path, c.StatusCode, time.Since(start))
	}
}
