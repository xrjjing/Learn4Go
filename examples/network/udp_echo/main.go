package main

import (
	"fmt"
	"net"
)

func main() {
	go startUDPServer()
	startUDPClient()
}

func startUDPServer() {
	addr, _ := net.ResolveUDPAddr("udp", ":9002")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		conn.WriteToUDP(append([]byte("echo: "), buf[:n]...), remote)
	}
}

func startUDPClient() {
	conn, err := net.Dial("udp", "127.0.0.1:9002")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("hello-udp"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("udp client recv:", string(buf[:n]))
}
