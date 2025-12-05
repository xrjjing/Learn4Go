package main

import (
	"errors"
	"testing"
)

func TestErrorWrap(t *testing.T) {
	_, err := find(0)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected errors.Is to match ErrNotFound, got %v", err)
	}
	if err == ErrNotFound {
		t.Fatalf("wrapped error should not be identical")
	}
}
