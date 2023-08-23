package float16

import (
	"math"
	"runtime"
	"testing"
)

var negZero = math.Float64frombits(1 << 63)

func TestFromFloat32(t *testing.T) {
	tests := []struct {
		f float32
		r Float16
	}{
		// from https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		{0, 0x0000},
		{0x1p-24, 0x0001},     // smallest positive subnormal number
		{0x1.ff8p-15, 0x03ff}, // largest positive subnormal number
		{0x1p-14, 0x0400},     // smallest positive normal number
		{0x1.554p-02, 0x3555}, // nearest value to 1/3
		{0x1.ffcp-01, 0x3bff}, // largest number less than one
		{0x1p+00, 0x3c00},     // one
		{0x1.004p+00, 0x3c01}, // smallest number larger than one
		{0x1.ffcp+15, 0x7bff}, // largest normal number
		{float32(negZero), 0x8000},
		{-2, 0xc000},

		// rounds to nearest even
		{0x1.002p+00, 0x3c00},
		{math.Nextafter32(0x1.002p+00, 2), 0x3c01},
		{math.Nextafter32(0x1.006p+00, 0), 0x3c01},
		{0x1.006p+00, 0x3c02},

		// underflow
		{math.Nextafter32(0x1p-24, 0), 0x0000},
		{0x1p-126, 0x0000},
		{0x1.fffffcp-127, 0x0000},

		// overflow
		{0x1p+16, 0x7c00},
		{0x1p+17, 0x7c00},
		{-0x1p+16, 0xfc00},
		{-0x1p+17, 0xfc00},

		// infinities
		{float32(math.Inf(1)), 0x7c00},
		{float32(math.Inf(-1)), 0xfc00},
	}
	for _, tt := range tests {
		r := FromFloat32(tt.f)
		if r != tt.r {
			t.Errorf("expected %x, got %x", tt.r, r)
		}
	}
}

func TestFromFloat64(t *testing.T) {
	tests := []struct {
		f float64
		r Float16
	}{
		// from https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		{0, 0x0000},
		{0x1p-24, 0x0001},     // smallest positive subnormal number
		{0x1.ff8p-15, 0x03ff}, // largest positive subnormal number
		{0x1p-14, 0x0400},     // smallest positive normal number
		{0x1.554p-02, 0x3555}, // nearest value to 1/3
		{0x1.ffcp-01, 0x3bff}, // largest number less than one
		{0x1p+00, 0x3c00},     // one
		{0x1.004p+00, 0x3c01}, // smallest number larger than one
		{0x1.ffcp+15, 0x7bff}, // largest normal number
		{negZero, 0x8000},
		{-2, 0xc000},

		// infinities
		{math.Inf(1), 0x7c00},
		{math.Inf(-1), 0xfc00},
	}
	for _, tt := range tests {
		r := FromFloat64(tt.f)
		if r != tt.r {
			t.Errorf("expected %x, got %x", tt.r, r)
		}
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		f Float16
		r float32
	}{
		// from https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		{0x0000, 0},
		{0x0001, 0x1p-24},     // smallest positive subnormal number
		{0x03ff, 0x1.ff8p-15}, // largest positive subnormal number
		{0x0400, 0x1p-14},     // smallest positive normal number
		{0x3555, 0x1.554p-02}, // nearest value to 1/3
		{0x3bff, 0x1.ffcp-01}, // largest number less than one
		{0x3c00, 0x1p+00},     // one
		{0x3c01, 0x1.004p+00}, // smallest number larger than one
		{0x7bff, 0x1.ffcp+15}, // largest normal number
		{0x8000, -0},
		{0xc000, -2},
	}

	for _, tt := range tests {
		r := tt.f.Float32()
		if r != tt.r {
			t.Errorf("expected %x, got %x", tt.r, r)
		}
	}
}

func TestFloat32_Specials(t *testing.T) {
	// infinity
	if r := Inf(1).Float32(); !math.IsInf(float64(r), 1) {
		t.Errorf("expected +Inf, got %x", r)
	}

	// negative infinity
	if r := Inf(-1).Float32(); !math.IsInf(float64(r), -1) {
		t.Errorf("expected -Inf, got %x", r)
	}

	// NaN
	if r := NaN().Float32(); !math.IsNaN(float64(r)) {
		t.Errorf("expected NaN, got %x", r)
	}

	// zero
	if r := Float16(0x0000).Float32(); r != 0 || !math.IsInf(float64(1/r), 1) {
		t.Errorf("expected +0, got %x", r)
	}

	// negative zero
	if r := Float16(0x8000).Float32(); r != 0 || !math.IsInf(float64(1/r), -1) {
		t.Errorf("expected -0, got %x", r)
	}
}

func BenchmarkFloat32(b *testing.B) {
	f := Float16(0x3555)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(f.Float32())
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		f Float16
		r float64
	}{
		// from https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		{0x0000, 0},
		{0x0001, 0x1p-24},     // smallest positive subnormal number
		{0x03ff, 0x1.ff8p-15}, // largest positive subnormal number
		{0x0400, 0x1p-14},     // smallest positive normal number
		{0x3555, 0x1.554p-02}, // nearest value to 1/3
		{0x3bff, 0x1.ffcp-01}, // largest number less than one
		{0x3c00, 0x1p+00},     // one
		{0x3c01, 0x1.004p+00}, // smallest number larger than one
		{0x7bff, 0x1.ffcp+15}, // largest normal number
		{0x8000, -0},
		{0xc000, -2},
	}

	for _, tt := range tests {
		r := tt.f.Float64()
		if r != tt.r {
			t.Errorf("expected %x, got %x", tt.r, r)
		}
	}
}

func TestFloat64_Specials(t *testing.T) {
	// infinity
	if r := Inf(1).Float64(); !math.IsInf(r, 1) {
		t.Errorf("expected +Inf, got %x", r)
	}

	// negative infinity
	if r := Inf(-1).Float64(); !math.IsInf(r, -1) {
		t.Errorf("expected -Inf, got %x", r)
	}

	// NaN
	if r := NaN().Float64(); !math.IsNaN(r) {
		t.Errorf("expected NaN, got %x", r)
	}

	// zero
	if r := Float16(0x0000).Float64(); r != 0 || !math.IsInf(1/r, 1) {
		t.Errorf("expected +0, got %x", r)
	}

	// negative zero
	if r := Float16(0x8000).Float64(); r != 0 || !math.IsInf(1/r, -1) {
		t.Errorf("expected -0, got %x", r)
	}
}
