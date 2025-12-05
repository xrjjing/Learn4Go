package main

import "fmt"

// 演示 iota 生成位掩码
const (
	FlagRead = 1 << iota
	FlagWrite
	FlagExec
)

func main() {
	fmt.Println("read", FlagRead, "write", FlagWrite, "exec", FlagExec)
	mask := FlagRead | FlagExec
	fmt.Println("mask read?", mask&FlagRead != 0, "write?", mask&FlagWrite != 0)
}
