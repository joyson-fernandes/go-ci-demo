package main

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("Joyson")
	want := "Hello, Joyson!"
	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 3},
		{0, 0, 0},
		{-1, 1, 0},
		{100, 200, 300},
	}
	for _, tt := range tests {
		got := Add(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
