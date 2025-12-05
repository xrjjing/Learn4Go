package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
)

// Day4/5：分组与中间件链
func main() {
	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	api := r.Group("/api")
	api.Use(func(c *tinygee.Context) {
		c.SetHeader("X-API", "v1")
		c.Next()
	})

	api.GET("/slow", func(c *tinygee.Context) {
		time.Sleep(50 * time.Millisecond)
		c.JSON(http.StatusOK, map[string]any{"ok": true})
	})

	log.Println("tinygee day5 on :8087")
	log.Fatal(r.Run(":8087"))
}
