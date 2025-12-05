package main

import "testing"

func TestMapInit(t *testing.T) {
	var m map[string]int
	if m != nil {
		t.Fatalf("expected nil map initially")
	}
	m = make(map[string]int)
	m["x"] = 1
	if m["x"] != 1 {
		t.Fatalf("expected value 1, got %d", m["x"])
	}
	if _, ok := m["absent"]; ok {
		t.Fatalf("expected absent key not ok")
	}
}
