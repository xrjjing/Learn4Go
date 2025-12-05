package main

import (
	"log"
	"net/http"
	"time"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
	"github.com/xrjjing/Learn4Go/tinygee/middleware/ratelimit"
)

func main() {
	app := tinygee.New()
	app.Use(middleware.Logger(), middleware.Recover())

	limiter := ratelimit.New(2*time.Second, 3) // 2s 内最多 3 次
	app.Use(limiter.Middleware())

	app.GET("/ping", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	log.Println("tinygee ratelimit demo on :8090")
	log.Fatal(app.Run(":8090"))
}
