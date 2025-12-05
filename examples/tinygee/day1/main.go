package main

import (
	"log"
	"net/http"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
)

// Day1/2：基础路由与 JSON/String
func main() {
	app := tinygee.New()
	app.Use(middleware.Logger(), middleware.Recover())

	app.GET("/", func(c *tinygee.Context) {
		c.String(http.StatusOK, "hello tinygee")
	})

	app.POST("/echo", func(c *tinygee.Context) {
		c.JSON(http.StatusOK, map[string]string{"path": c.Path, "method": c.Method})
	})

	log.Println("tinygee day1 on :8085")
	log.Fatal(app.Run(":8085"))
}
