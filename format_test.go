package float16

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		format string
		x      Float16
		want   string
	}{
		{"%b", FromFloat64(0), "0p-24"},

		{"%f", FromFloat64(0.5), "0.5"},
		{"%f", FromFloat64(-0.5), "-0.5"},
		{"%+f", FromFloat64(0.5), "+0.5"},
		{"%+f", FromFloat64(-0.5), "-0.5"},
		{"% f", FromFloat64(0.5), " 0.5"},
		{"% f", FromFloat64(-0.5), "-0.5"},
		{"%8f", FromFloat64(0.5), "     0.5"},
		{"%-8f", FromFloat64(0.5), "0.5     "},

		{"%.6e", FromFloat64(0.5), "5.000000e-01"},

		{"%g", FromFloat64(0.5), "0.5"},

		{"%x", FromFloat64(0.5), "0x1p-01"},
		{"%#x", FromFloat64(0.5), "0x1p-01"},

		{"%X", FromFloat64(0.5), "0X1P-01"},
		{"%#X", FromFloat64(0.5), "0X1P-01"},

		{"%v", FromFloat64(0.5), "0.5"},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.x)
		if got != tt.want {
			t.Errorf("expected %s, got %s", tt.want, got)
		}
	}
}
