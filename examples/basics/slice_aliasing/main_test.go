package main

import "testing"

func TestSliceAliasing(t *testing.T) {
	base := []int{1, 2, 3, 4}
	a := base[:2]
	b := base[1:]
	a[1] = 99
	if b[0] != 99 {
		t.Fatalf("expected shared element change, got %d", b[0])
	}
	// append may reallocate; ensure base unaffected after append when cap exceeded
	a = append(a, 5, 6, 7) // force new backing
	a[0] = 42
	if base[0] != 1 {
		t.Fatalf("base should stay 1 after append reallocation, got %d", base[0])
	}
}
