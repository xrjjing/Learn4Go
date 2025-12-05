package main

import (
	"expvar"
	"flag"
	"log"
	"net/http"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
)

// Demo 入口：可选开启 /metrics。
func main() {
	port := flag.String("port", ":9999", "listen address")
	enableProm := flag.Bool("prom", false, "enable /metrics")
	flag.Parse()

	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	r.GET("/", func(c *tinygee.Context) {
		c.String(http.StatusOK, "welcome to tinygee")
	})
	r.GET("/ping", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	if *enableProm {
		// 使用标准库 expvar 提供基础指标
		r.GET("/metrics", func(c *tinygee.Context) {
			expvar.Handler().ServeHTTP(c.Writer, c.Req)
		})
	}

	log.Printf("TinyGee demo on %s (prometheus: %v)", *port, *enableProm)
	if err := r.Run(*port); err != nil {
		log.Fatal(err)
	}
}
