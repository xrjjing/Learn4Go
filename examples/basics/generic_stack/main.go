package main

import "fmt"

// 一个最简泛型栈示例
type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	last := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return last, true
}

func main() {
	var s Stack[string]
	s.Push("a")
	s.Push("b")
	if v, ok := s.Pop(); ok {
		fmt.Println("pop:", v)
	}
}
