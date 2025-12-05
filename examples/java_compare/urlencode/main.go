package main

import (
	"fmt"
	"net/url"
)

// 演示 Go 标准库 urlencode/urldecode，与 Java URLEncoder 类似。
func main() {
	raw := "https://example.com/search?q=Go 语言&lang=zh"
	encoded := url.QueryEscape(raw)
	fmt.Println("encoded:", encoded)

	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println("decoded:", decoded)
}
