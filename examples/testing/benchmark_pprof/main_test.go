package main

import "testing"

func BenchmarkJoinRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinRepeat(1000)
	}
}
