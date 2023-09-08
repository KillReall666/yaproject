package update

import "testing"

func TestValueConvInt(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		intWant int64
	}{

		{name: "simple int", val: "123", intWant: 123},
		{name: "negative int", val: "-123", intWant: -123},
		{name: "zero", val: "0", intWant: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntValueConv(tt.val); got != tt.intWant {
				t.Errorf("ValueConv() = %v, want %v", got, tt.intWant)
			}
		})
	}
}
