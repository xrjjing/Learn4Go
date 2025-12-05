package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
	"github.com/xrjjing/Learn4Go/tinygee/render"
)

func main() {
	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	funcMap := template.FuncMap{
		"upper": func(s string) string { return template.HTMLEscapeString(s) },
	}
	renderer, err := render.New("examples/tinygee/day6/templates/*.html", funcMap)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/hello", func(c *tinygee.Context) {
		renderer.HTML(c, http.StatusOK, "hello.html", map[string]any{"Name": "TinyGee"})
	})

	// 静态文件
	r.GET("/static/*filepath", render.Static("/static/", http.Dir("examples/tinygee/day6/static")))

	log.Println("tinygee day6 on :8088")
	log.Fatal(r.Run(":8088"))
}
