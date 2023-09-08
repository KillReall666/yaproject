package update

import (
	"testing"
)

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
			if got := FloatValueConv(tt.val); got != tt.want {
				t.Errorf("ValueConv() = %v want %v", got, tt.want)
			}

		})
	}
}
