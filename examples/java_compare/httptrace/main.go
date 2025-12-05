package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

// 演示 httptrace 捕获 DNS/连接/首字节时间，对照 Java HttpClient + Listener。
func main() {
	req, _ := http.NewRequest("GET", "https://httpbin.org/delay/1", nil)

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) { fmt.Println("DNS start:", info.Host) },
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fmt.Println("DNS done:", info.Addrs, "err:", info.Err)
		},
		ConnectStart: func(network, addr string) { fmt.Println("Dial start:", network, addr) },
		ConnectDone: func(network, addr string, err error) {
			fmt.Println("Dial done:", addr, "err:", err)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Println("GotConn reused:", info.Reused)
		},
		GotFirstResponseByte: func() { fmt.Println("First byte at", time.Now().Format(time.RFC3339Nano)) },
	}

	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	client := &http.Client{Timeout: 3 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("status:", resp.Status, "total:", time.Since(start))
}
