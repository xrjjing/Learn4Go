package main

import (
	"bufio"
	"fmt"
	"net"
)

// 最小 TCP Echo：先启动 server，再运行 client。
func main() {
	go startServer()
	startClient()
}

func startServer() {
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			reader := bufio.NewReader(c)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				c.Write([]byte("echo: " + line))
			}
		}(conn)
	}
}

func startClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Fprintf(conn, "hello\n")
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("client recv:", reply)
}
