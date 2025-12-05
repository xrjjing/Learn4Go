package main

import "fmt"

//go:generate sh -c "echo '// generated value' > generated.txt && echo 'const Generated = 42' >> generated.txt"

func main() {
	fmt.Println("run `go generate ./...` to create generated.txt")
}
