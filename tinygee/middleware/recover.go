package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/xrjjing/Learn4Go/tinygee"
)

// RecoverConfig 用于配置恢复响应格式。
type RecoverConfig struct {
	JSON bool // true 返回 JSON，false 返回纯文本
}

// Recover 捕获 panic，避免服务器崩溃。
func Recover(cfgs ...RecoverConfig) tinygee.HandlerFunc {
	cfg := RecoverConfig{JSON: true}
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	return func(c *tinygee.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v\n%s", err, debug.Stack())
				if cfg.JSON {
					c.JSON(http.StatusInternalServerError, map[string]any{
						"error": "internal server error",
					})
				} else {
					c.String(http.StatusInternalServerError, "internal server error")
				}
			}
		}()
		c.Next()
	}
}
