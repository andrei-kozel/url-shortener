package random

import "testing"

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "empty string", length: 0},
		{name: "short string", length: 1},
		{name: "medium string", length: 8},
		{name: "long string", length: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.length)
			if len(got) != tt.length {
				t.Errorf("NewRandomString() = %v, want %v", len(got), tt.length)
			}
		})
	}
}
