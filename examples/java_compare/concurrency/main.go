package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// 演示 WaitGroup + channel 实现并发抓取，对照 Java ExecutorService+Future。
func main() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/get",
		"https://httpbin.org/uuid",
	}

	type result struct {
		url     string
		status  int
		err     error
		latency time.Duration
	}

	var wg sync.WaitGroup
	out := make(chan result, len(urls))

	client := &http.Client{Timeout: 3 * time.Second}

	for _, u := range urls {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(target)
			if err != nil {
				out <- result{url: target, err: err}
				return
			}
			resp.Body.Close()
			out <- result{url: target, status: resp.StatusCode, latency: time.Since(start)}
		}(u)
	}

	wg.Wait()
	close(out)

	for r := range out {
		if r.err != nil {
			fmt.Printf("%s -> error: %v\n", r.url, r.err)
			continue
		}
		fmt.Printf("%s -> %d in %v\n", r.url, r.status, r.latency)
	}
}
