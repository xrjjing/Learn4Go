package main

import (
	"log"
	"net/http"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
)

// Day3：动态路由与通配符
func main() {
	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	r.GET("/hello/:name", func(c *tinygee.Context) {
		c.String(http.StatusOK, "hello %s", c.Param("name"))
	})

	r.GET("/assets/*filepath", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"path": c.Param("filepath")})
	})

	log.Println("tinygee day3 on :8086")
	log.Fatal(r.Run(":8086"))
}
