package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/xrjjing/Learn4Go/internal/todo"
)

// 入口：启动内存版 TODO API
func main() {
	rand.Seed(time.Now().UnixNano())
	s := todo.NewServer(todo.NewStore())
	addr := ":8080"
	log.Printf("TODO API listening on %s", addr)
	if err := http.ListenAndServe(addr, s.Handler()); err != nil {
		log.Fatal(err)
	}
}
