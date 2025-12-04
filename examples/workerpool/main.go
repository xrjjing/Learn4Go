package main

import (
	"fmt"
	"time"

	"github.com/xrjjing/Learn4Go/internal/workerpool"
)

func main() {
	pool := workerpool.New(2, 5)
	pool.Run()

	for i := 1; i <= 5; i++ {
		jobID := i
		pool.Submit(func() error {
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("Task %d completed\n", jobID)
			return nil
		})
	}

	pool.Close()
	pool.Wait()
	fmt.Println("All tasks completed")
}
