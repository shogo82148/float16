package float16

import (
	"strconv"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		x Float16
		s string
	}{
		// special cases
		{0, "0"},
		{0x8000, "-0"},
		{Inf(1), "+Inf"},
		{Inf(-1), "-Inf"},
		{NaN(), "NaN"},

		{FromFloat64(1), "1"},
		{FromFloat64(1.5), "1.5"},
		{FromFloat64(1.25), "1.25"},
		{FromFloat64(1.125), "1.125"},
		{FromFloat64(1.0625), "1.0625"},
		{FromFloat64(1.03125), "1.03125"},
		{FromFloat64(1.015625), "1.015625"},
		{FromFloat64(1.0078125), "1.0078125"},
		{FromFloat64(1.00390625), "1.00390625"},
		{FromFloat64(1.001953125), "1.001953125"},
		{FromFloat64(1.0009765625), "1.0009765625"},

		{FromFloat64(2), "2"},
		{FromFloat64(4), "4"},
		{FromFloat64(8), "8"},
		{FromFloat64(16), "16"},
		{FromFloat64(32), "32"},
		{FromFloat64(64), "64"},
		{FromFloat64(128), "128"},
		{FromFloat64(256), "256"},
		{FromFloat64(512), "512"},
		{FromFloat64(1024), "1024"},
		{FromFloat64(2048), "2048"},
		{FromFloat64(4096), "4096"},
		{FromFloat64(8192), "8192"},
		{FromFloat64(16384), "16384"},
		{FromFloat64(32768), "32768"},
		{FromFloat64(65504), "65504"}, // max normal
	}

	for _, tt := range tests {
		got := tt.x.String()
		if got != tt.s {
			t.Errorf("expected %s, got %s", tt.s, got)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		x    Float16
		fmt  byte
		prec int
		s    string
	}{
		// binary exponent formats

		// {0, 'b', -1, "0"},
		// {0x8000, 'b', -1, "-0"},
		{Inf(1), 'b', -1, "+Inf"},
		{Inf(-1), 'b', -1, "-Inf"},
		{NaN(), 'b', -1, "NaN"},

		// decimal exponent formats
		{0, 'e', -1, "0e+00"},
		{0x8000, 'e', -1, "-0e+00"},
		{Inf(1), 'e', -1, "+Inf"},
		{Inf(-1), 'e', -1, "-Inf"},
		{NaN(), 'e', -1, "NaN"},

		{0, 'E', -1, "0E+00"},
		{0x8000, 'E', -1, "-0E+00"},
		{Inf(1), 'E', -1, "+Inf"},
		{Inf(-1), 'E', -1, "-Inf"},
		{NaN(), 'E', -1, "NaN"},

		// decimal formats
		{0, 'f', -1, "0"},
		{0x8000, 'f', -1, "-0"},
		{Inf(1), 'f', -1, "+Inf"},
		{Inf(-1), 'f', -1, "-Inf"},
		{NaN(), 'f', -1, "NaN"},

		// alternate formats
		{0, 'g', -1, "0"},
		{0x8000, 'g', -1, "-0"},
		{Inf(1), 'g', -1, "+Inf"},
		{Inf(-1), 'g', -1, "-Inf"},
		{NaN(), 'g', -1, "NaN"},

		{0, 'G', -1, "0"},
		{0x8000, 'G', -1, "-0"},
		{Inf(1), 'G', -1, "+Inf"},
		{Inf(-1), 'G', -1, "-Inf"},
		{NaN(), 'G', -1, "NaN"},

		// hexadecimal formats
		{0, 'x', -1, "0x0p+00"},
		{0x8000, 'x', -1, "-0x0p+00"},
		{Inf(1), 'x', -1, "+Inf"},
		{Inf(-1), 'x', -1, "-Inf"},
		{NaN(), 'x', -1, "NaN"},

		{FromFloat64(0), 'x', 0, "0x0p+00"},
		{FromFloat64(0x1.0p0), 'x', 0, "0x1p+00"},
		{FromFloat64(0x1.7p0), 'x', 0, "0x1p+00"},
		{FromFloat64(0x1.8p0), 'x', 0, "0x1p+01"},

		{FromFloat64(0x0.0p0), 'x', 1, "0x0.0p+00"},
		{FromFloat64(0x1.0p0), 'x', 1, "0x1.0p+00"},
		{FromFloat64(0x1.fp0), 'x', 1, "0x1.fp+00"},
		{FromFloat64(0x1.ffc0p0), 'x', 1, "0x1.0p+01"},
		{FromFloat64(0x1.08p0), 'x', 1, "0x1.0p+00"},
		{FromFloat64(0x1.18p0), 'x', 1, "0x1.2p+00"},

		{FromFloat64(0x0.00p0), 'x', 2, "0x0.00p+00"},
		{FromFloat64(0x1.00p0), 'x', 2, "0x1.00p+00"},
		{FromFloat64(0x1.0fp0), 'x', 2, "0x1.0fp+00"},
		{FromFloat64(0x1.ffc0p0), 'x', 2, "0x1.00p+01"},
		{FromFloat64(0x1.008p0), 'x', 2, "0x1.00p+00"},
		{FromFloat64(0x1.018p0), 'x', 2, "0x1.02p+00"},

		{FromFloat64(0x0.000p0), 'x', 3, "0x0.000p+00"},
		{FromFloat64(0x1.000p0), 'x', 3, "0x1.000p+00"},
		{FromFloat64(0x1.00c0p0), 'x', 3, "0x1.00cp+00"},
		{FromFloat64(0x1.ffc0p0), 'x', 3, "0x1.ffcp+00"},

		{FromFloat64(0x1.ffc0p0), 'x', 4, "0x1.ffc0p+00"},
		{FromFloat64(0x1p-24), 'x', 4, "0x1.0000p-24"},

		{FromFloat64(0x1p0), 'x', -1, "0x1p+00"},
		{FromFloat64(0x1.8p0), 'x', -1, "0x1.8p+00"},
		{FromFloat64(0x1.08p0), 'x', -1, "0x1.08p+00"},
		{FromFloat64(0x1.00cp0), 'x', -1, "0x1.00cp+00"},

		{0, 'X', -1, "0X0P+00"},
		{0x8000, 'X', -1, "-0X0P+00"},
		{Inf(1), 'X', -1, "+Inf"},
		{Inf(-1), 'X', -1, "-Inf"},
		{NaN(), 'X', -1, "NaN"},
	}

	for _, tt := range tests {
		got := tt.x.Format(tt.fmt, tt.prec)
		if got != tt.s {
			t.Errorf("%#v: expected %s, got %s", tt, tt.s, got)
		}
	}
}

func FuzzFormat(f *testing.F) {
	f.Add(uint16(0), 'b', 1)
	f.Add(uint16(0), 'e', 1)
	f.Add(uint16(0), 'E', 1)
	f.Add(uint16(0), 'f', 1)
	f.Add(uint16(0), 'g', 1)
	f.Add(uint16(0), 'G', 1)
	f.Add(uint16(0), 'x', 1)
	f.Add(uint16(0), 'X', 1)
	f.Fuzz(func(t *testing.T, x uint16, fmt rune, prec int) {
		if prec < 0 || prec > 100 {
			return
		}
		switch fmt {
		case 'b', 'e', 'E', 'f', 'g', 'G', 'x', 'X':
		default:
			return
		}
		if x&signMask16 == 0 && fmt == 'b' {
			return
		}

		f := FromBits(x)
		got := f.Format(byte(fmt), prec)
		want := strconv.FormatFloat(f.Float64(), byte(fmt), prec, 64)
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
}
