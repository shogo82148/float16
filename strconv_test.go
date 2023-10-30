package float16

import (
	"cmp"
	"fmt"
	"strconv"
	"testing"
)

// exact returns the exact float16 representation of f.
// it panics if f doesn't have an exact float16 representation.
func exact(f float64) Float16 {
	ret := FromFloat64(f)
	if cmp.Compare(ret.Float64(), f) != 0 {
		panic(fmt.Sprintf("%f doesn't have exact float16 representation", f))
	}
	return ret
}

func TestParse(t *testing.T) {
	tests := []struct {
		s string
		x Float16
	}{
		{"0", 0},
		{"-0", 0x8000},
		{"+Inf", Inf(1)},
		{"-Inf", Inf(-1)},
		{"+infinity", Inf(1)},
		{"-infinity", Inf(-1)},
		{"NaN", NaN()},

		// greater than one
		{"1", exact(1)},
		{"2", exact(2)},
		{"4", exact(4)},
		{"8", exact(8)},
		{"16", exact(16)},
		{"32", exact(32)},
		{"64", exact(64)},
		{"128", exact(128)},
		{"65504", exact(65504)},

		// less than one
		{"0.5", exact(0.5)},
		{"0.25", exact(0.25)},
		{"0.125", exact(0.125)},
		{"0.00024414062", exact(0x1p-12)},
		{"0.00012207031", exact(0x1p-13)},
		{"0.00006103515625", exact(0x1p-14)},

		// subnormal numbers
		{"0.000030517578125", exact(0x1p-15)},
		{"0.000000059604644775390625", exact(0x1p-24)},

		// test rounding
		{"1.0009765625", exact(1.0009765625)}, // minimum value greater than one
		{"1.00048828125", exact(1)},
		{"1.0004882812500001", exact(1.0009765625)},
		{"1.0014648437499999", exact(1.0009765625)},
		{"1.00146484375", exact(1.001953125)},

		// hexadecimal
		// {"0x1p0", exact(1)},
		{"0x1p-14", exact(0x1p-14)}, // minimum nominal
		// {"0x1p-24", exact(0x1p-14)}, // minimum subnormal greater than zero
	}

	for _, tt := range tests {
		got, err := Parse(tt.s)
		if err != nil {
			t.Errorf("%q: expected no error, got %v", tt.s, err)
		}
		if got != tt.x {
			t.Errorf("%q: expected %v, got %v", tt.s, tt.x, got)
		}
	}
}

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

		{exact(1), "1"},
		{exact(1.5), "1.5"},
		{exact(1.25), "1.25"},
		{exact(1.125), "1.125"},

		{exact(2), "2"},
		{exact(4), "4"},
		{exact(8), "8"},
		{exact(16), "16"},
		{exact(32), "32"},
		{exact(64), "64"},
		{exact(128), "128"},
		{exact(256), "256"},
		{exact(512), "512"},
		{exact(1024), "1024"},
		{exact(2048), "2048"},
		{exact(4096), "4096"},
		{exact(8192), "8192"},
		{exact(16384), "16384"},
		{exact(32768), "32768"},
		{exact(65504), "65504"}, // max normal
	}

	for _, tt := range tests {
		got := tt.x.String()
		if got != tt.s {
			t.Errorf("expected %s, got %s", tt.s, got)
		}
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		x    Float16
		fmt  byte
		prec int
		s    string
	}{
		/****** binary exponent formats ******/
		{0, 'b', -1, "0p-24"},
		{0x8000, 'b', -1, "-0p-24"},
		{Inf(1), 'b', -1, "+Inf"},
		{Inf(-1), 'b', -1, "-Inf"},
		{NaN(), 'b', -1, "NaN"},

		{0x0001, 'b', -1, "1p-24"},
		{exact(1), 'b', -1, "1024p-10"},
		{exact(1.5), 'b', -1, "1536p-10"},
		{exact(65504), 'b', -1, "2047p+5"},

		/****** decimal exponent formats ******/
		{0, 'e', -1, "0e+00"},
		{0x8000, 'e', -1, "-0e+00"},
		{Inf(1), 'e', -1, "+Inf"},
		{Inf(-1), 'e', -1, "-Inf"},
		{NaN(), 'e', -1, "NaN"},

		{exact(0), 'e', 16, "0.0000000000000000e+00"},
		{exact(0x1p-24), 'e', 16, "5.9604644775390625e-08"},
		{exact(0x2p-24), 'e', 16, "1.1920928955078125e-07"},
		{exact(0x3p-24), 'e', 16, "1.7881393432617188e-07"},
		{exact(0x4p-24), 'e', 16, "2.3841857910156250e-07"},
		{exact(0x5p-24), 'e', 16, "2.9802322387695312e-07"},
		{exact(0x6p-24), 'e', 16, "3.5762786865234375e-07"},
		{exact(0x7p-24), 'e', 16, "4.1723251342773438e-07"},
		{exact(0x8p-24), 'e', 16, "4.7683715820312500e-07"},
		{exact(0x9p-24), 'e', 16, "5.3644180297851562e-07"},
		{exact(0xap-24), 'e', 16, "5.9604644775390625e-07"},

		{exact(0x1p-24), 'e', -1, "6e-08"},
		{exact(0x2p-24), 'e', -1, "1e-07"},
		{exact(0x3p-24), 'e', -1, "2e-07"},
		{exact(0x4p-24), 'e', -1, "2.4e-07"},
		{exact(0x5p-24), 'e', -1, "3e-07"},
		{exact(0x6p-24), 'e', -1, "3.6e-07"},
		{exact(0x7p-24), 'e', -1, "4e-07"},
		{exact(0x8p-24), 'e', -1, "5e-07"},
		{exact(0x9p-24), 'e', -1, "5.4e-07"},
		{exact(0xap-24), 'e', -1, "6e-07"},
		{exact(10), 'e', -1, "1e+01"},
		{exact(100), 'e', -1, "1e+02"},
		{exact(1000), 'e', -1, "1e+03"},
		{exact(10000), 'e', -1, "1e+04"},
		{exact(65504), 'e', -1, "6.55e+04"},

		// random numbers
		{0x816f, 'e', -1, "-2.19e-05"},
		{0xadaa, 'e', -1, "-8.85e-02"},
		{0x92dc, 'e', -1, "-8.373e-04"},
		{0xf4ca, 'e', -1, "-1.962e+04"},
		{0x7aa7, 'e', -1, "5.45e+04"},
		{0x3734, 'e', -1, "4.502e-01"},
		{0x2ac3, 'e', -1, "5.283e-02"},
		{0xe6d6, 'e', -1, "-1.75e+03"},
		{0x700f, 'e', -1, "8.31e+03"},

		{0, 'E', -1, "0E+00"},
		{0x8000, 'E', -1, "-0E+00"},
		{Inf(1), 'E', -1, "+Inf"},
		{Inf(-1), 'E', -1, "-Inf"},
		{NaN(), 'E', -1, "NaN"},

		/****** decimal formats ******/
		{0, 'f', -1, "0"},
		{0x8000, 'f', -1, "-0"},
		{Inf(1), 'f', -1, "+Inf"},
		{Inf(-1), 'f', -1, "-Inf"},
		{NaN(), 'f', -1, "NaN"},

		{exact(0x1p-24), 'f', 24, "0.000000059604644775390625"},
		{exact(0x2p-24), 'f', 24, "0.000000119209289550781250"},
		{exact(0x3p-24), 'f', 24, "0.000000178813934326171875"},
		{exact(0x4p-24), 'f', 24, "0.000000238418579101562500"},
		{exact(0x5p-24), 'f', 24, "0.000000298023223876953125"},
		{exact(0x6p-24), 'f', 24, "0.000000357627868652343750"},
		{exact(0x7p-24), 'f', 24, "0.000000417232513427734375"},
		{exact(0x8p-24), 'f', 24, "0.000000476837158203125000"},
		{exact(0x9p-24), 'f', 24, "0.000000536441802978515625"},
		{exact(0xap-24), 'f', 24, "0.000000596046447753906250"},

		{exact(0x1p-24), 'f', 23, "0.00000005960464477539062"},
		{exact(0x2p-24), 'f', 23, "0.00000011920928955078125"},
		{exact(0x3p-24), 'f', 23, "0.00000017881393432617188"},
		{exact(0x4p-24), 'f', 23, "0.00000023841857910156250"},
		{exact(0x5p-24), 'f', 23, "0.00000029802322387695312"},
		{exact(0x6p-24), 'f', 23, "0.00000035762786865234375"},
		{exact(0x7p-24), 'f', 23, "0.00000041723251342773438"},
		{exact(0x8p-24), 'f', 23, "0.00000047683715820312500"},
		{exact(0x9p-24), 'f', 23, "0.00000053644180297851562"},
		{exact(0xap-24), 'f', 23, "0.00000059604644775390625"},

		{exact(0x1p-24), 'f', -1, "0.00000006"},
		{exact(0x2p-24), 'f', -1, "0.0000001"},
		{exact(0x3p-24), 'f', -1, "0.0000002"},
		{exact(0x4p-24), 'f', -1, "0.00000024"},
		{exact(0x5p-24), 'f', -1, "0.0000003"},
		{exact(0x6p-24), 'f', -1, "0.00000036"},
		{exact(0x7p-24), 'f', -1, "0.0000004"},
		{exact(0x8p-24), 'f', -1, "0.0000005"},
		{exact(0x9p-24), 'f', -1, "0.00000054"},
		{exact(0xap-24), 'f', -1, "0.0000006"},

		{exact(10), 'f', -1, "10"},
		{exact(100), 'f', -1, "100"},
		{exact(1000), 'f', -1, "1000"},
		{exact(10000), 'f', -1, "10000"},
		{exact(65504), 'f', -1, "65504"},
		{exact(-1), 'f', -1, "-1"},
		{exact(-10), 'f', -1, "-10"},
		{exact(-100), 'f', -1, "-100"},
		{exact(-1000), 'f', -1, "-1000"},
		{exact(-10000), 'f', -1, "-10000"},
		{exact(-65504), 'f', -1, "-65504"},

		// random numbers
		{0x816f, 'f', -1, "-0.0000219"},
		{0xadaa, 'f', -1, "-0.0885"},
		{0x92dc, 'f', -1, "-0.0008373"},
		{0xf4ca, 'f', -1, "-19616"},
		{0x7aa7, 'f', -1, "54496"},
		{0x3734, 'f', -1, "0.4502"},
		{0x2ac3, 'f', -1, "0.05283"},
		{0xe6d6, 'f', -1, "-1750"},
		{0x700f, 'f', -1, "8312"},

		/******* alternate formats *******/
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

		{exact(0), 'x', 0, "0x0p+00"},
		{exact(0x1.0p0), 'x', 0, "0x1p+00"},
		{exact(0x1.7p0), 'x', 0, "0x1p+00"},
		{exact(0x1.8p0), 'x', 0, "0x1p+01"},

		{exact(0x0.0p0), 'x', 1, "0x0.0p+00"},
		{exact(0x1.0p0), 'x', 1, "0x1.0p+00"},
		{exact(0x1.fp0), 'x', 1, "0x1.fp+00"},
		{exact(0x1.ffc0p0), 'x', 1, "0x1.0p+01"},
		{exact(0x1.08p0), 'x', 1, "0x1.0p+00"},
		{exact(0x1.18p0), 'x', 1, "0x1.2p+00"},

		{exact(0x0.00p0), 'x', 2, "0x0.00p+00"},
		{exact(0x1.00p0), 'x', 2, "0x1.00p+00"},
		{exact(0x1.0fp0), 'x', 2, "0x1.0fp+00"},
		{exact(0x1.ffc0p0), 'x', 2, "0x1.00p+01"},
		{exact(0x1.008p0), 'x', 2, "0x1.00p+00"},
		{exact(0x1.018p0), 'x', 2, "0x1.02p+00"},

		{exact(0x0.000p0), 'x', 3, "0x0.000p+00"},
		{exact(0x1.000p0), 'x', 3, "0x1.000p+00"},
		{exact(0x1.00c0p0), 'x', 3, "0x1.00cp+00"},
		{exact(0x1.ffc0p0), 'x', 3, "0x1.ffcp+00"},

		{exact(0x1.ffc0p0), 'x', 4, "0x1.ffc0p+00"},
		{exact(0x1p-24), 'x', 4, "0x1.0000p-24"},

		{exact(0x1p0), 'x', -1, "0x1p+00"},
		{exact(0x1.8p0), 'x', -1, "0x1.8p+00"},
		{exact(0x1.08p0), 'x', -1, "0x1.08p+00"},
		{exact(0x1.00cp0), 'x', -1, "0x1.00cp+00"},

		{0, 'X', -1, "0X0P+00"},
		{0x8000, 'X', -1, "-0X0P+00"},
		{Inf(1), 'X', -1, "+Inf"},
		{Inf(-1), 'X', -1, "-Inf"},
		{NaN(), 'X', -1, "NaN"},
		{exact(0x1.ffc0p0), 'X', 3, "0X1.FFCP+00"},
	}

	for _, tt := range tests {
		got := tt.x.Text(tt.fmt, tt.prec)
		if got != tt.s {
			t.Errorf("%#v: expected %s, got %s", tt, tt.s, got)
		}
	}
}

func TestString_RoundTrip(t *testing.T) {
	for i := 0; i < 0x10000; i++ {
		f := FromBits(uint16(i))
		str := f.String()
		got, err := Parse(str)
		if err != nil {
			t.Errorf("%04x: expected no error, got %v", i, err)
		}
		if got.IsNaN() && f.IsNaN() {
			continue
		}
		if got != f {
			t.Errorf("%04x: expected %v, got %v", i, f, got)
		}
	}
}

func FuzzParse(f *testing.F) {
	f.Add("0")
	f.Add("-0")
	f.Add("+Inf")
	f.Add("-Inf")
	f.Add("NaN")
	f.Add("1")
	f.Add("1.0009765625")
	f.Add("1.00048828125")

	f.Fuzz(func(t *testing.T, s string) {
		x0, err := Parse(s)
		if err != nil {
			return
		}
		s1 := x0.String()

		x1, err := Parse(s1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if x0.IsNaN() && x1.IsNaN() {
			return
		}
		if x0 != x1 {
			t.Errorf("expected %v, got %v", x0, x1)
		}
	})
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
		got := f.Text(byte(fmt), prec)
		want := strconv.FormatFloat(f.Float64(), byte(fmt), prec, 64)
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
}
