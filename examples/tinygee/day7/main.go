package main

import (
	"log"

	"github.com/xrjjing/Learn4Go/tinygee"
	"github.com/xrjjing/Learn4Go/tinygee/middleware"
)

func main() {
	r := tinygee.New()
	r.Use(middleware.Logger(), middleware.Recover())

	r.GET("/panic", func(c *tinygee.Context) {
		panic("boom")
	})

	log.Println("tinygee day7 on :8091")
	log.Fatal(r.Run(":8091"))
}
