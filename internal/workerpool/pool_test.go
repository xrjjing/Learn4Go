package workerpool

import (
	"sync/atomic"
	"testing"
)

func TestPool_RunAndSubmit(t *testing.T) {
	var counter int32
	pool := New(3, 5)
	pool.Run()

	for i := 0; i < 10; i++ {
		pool.Submit(func() error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	pool.Close()
	pool.Wait()

	if counter != 10 {
		t.Fatalf("expected 10 tasks executed, got %d", counter)
	}
}

func TestPool_ZeroTasks(t *testing.T) {
	pool := New(2, 5)
	pool.Run()
	pool.Close()
	pool.Wait()
}
