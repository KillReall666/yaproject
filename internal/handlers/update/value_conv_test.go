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
		{name: "int with dot", val: "123.0", intWant: 123},
		{name: "negative int with dot", val: "-123.0", intWant: -123},
		{name: "zero with dot", val: "0.0", intWant: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intValueConv(tt.val); got != tt.intWant {
				t.Errorf("ValueConv() = %v, want %v", got, tt.intWant)
			}
		})
	}
}

func TestValueConvFloat(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want float64
	}{

		{name: "simple float", val: "1.23", want: 1.23},
		{name: "negative float", val: "-1.23", want: -1.23},
		{name: "zero float", val: "0.0", want: 0.000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := floatValueConv(tt.val); got != tt.want {
				t.Errorf("ValueConv() = %v want %v", got, tt.want)
			}

		})
	}
}
