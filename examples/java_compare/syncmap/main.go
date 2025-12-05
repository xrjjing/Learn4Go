package main

import (
	"fmt"
	"sync"
)

// 演示 sync.Map 适合读多写少场景，对照 Java ConcurrentHashMap。
func main() {
	var m sync.Map
	keys := []string{"a", "b", "c", "a"}

	var wg sync.WaitGroup
	for _, k := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			m.LoadOrStore(key, 0)
			val, _ := m.Load(key)
			m.Store(key, val.(int)+1)
		}(k)
	}
	wg.Wait()

	m.Range(func(k, v any) bool {
		fmt.Printf("%s -> %d\n", k, v.(int))
		return true
	})
}
