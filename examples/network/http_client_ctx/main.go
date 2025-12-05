package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// 演示 http.Client + context 超时
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/1", nil)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("status:", resp.StatusCode)
}
