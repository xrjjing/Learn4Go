package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 演示 net/http 中间件链，对照 Java Servlet Filter。

// 简单日志中间件
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

// 认证占位中间件
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello middleware")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	// 组合中间件：logging -> auth -> handler
	var h http.Handler = mux
	h = auth(h)
	h = logging(h)

	addr := ":8090"
	log.Println("listen on", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}
