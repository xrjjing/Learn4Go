package main

import "testing"

func TestPointerVsValue(t *testing.T) {
	c := Counter{n: 1}
	c.AddValue(5)
	if c.n != 1 {
		t.Fatalf("value receiver should not change original, got %d", c.n)
	}
	c.AddPtr(5)
	if c.n != 6 {
		t.Fatalf("pointer receiver should change original, got %d", c.n)
	}
}
