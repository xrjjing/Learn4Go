package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func find(id int) (string, error) {
	if id == 0 {
		return "", fmt.Errorf("query id=%d: %w", id, ErrNotFound)
	}
	return "ok", nil
}

func main() {
	_, err := find(0)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("wrapped not found:", err)
	}
}
