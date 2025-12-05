package main

import "fmt"

// 演示接口隐式实现，对照 Java interface+implements。

type Storage interface {
	Save(data string) string
}

type FileStorage struct{}

func (FileStorage) Save(data string) string { return "file://" + data }

type MemoryStorage struct{}

func (MemoryStorage) Save(data string) string { return "mem://" + data }

// 不需要关键字 implements，只要方法匹配即可视为实现。
func persist(s Storage, data string) {
	fmt.Println(s.Save(data))
}

func main() {
	persist(FileStorage{}, "report.txt")
	persist(MemoryStorage{}, "session-1")
}
