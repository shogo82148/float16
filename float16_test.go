package float16

import (
	"math"
	"runtime"
	"testing"
)

func TestFloat32(t *testing.T) {
	tests := []struct {
		f Float16
		r float32
	}{
		// from https://en.wikipedia.org/wiki/Half-precision_floating-point_format
		{0x0000, 0},
		{0x0001, 0x1p-24},     // smallest positive subnormal number
		{0x03ff, 0x1.ff8p-24}, // largest positive subnormal number
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
	if r := Float16(0x0000).Float32(); !math.IsInf(float64(1/r), 1) {
		t.Errorf("expected +0, got %x", r)
	}

	// negative zero
	if r := Float16(0x8000).Float32(); !math.IsInf(float64(1/r), -1) {
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
		{0x03ff, 0x1.ff8p-24}, // largest positive subnormal number
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
	if r := Float16(0x0000).Float64(); !math.IsInf(1/r, 1) {
		t.Errorf("expected +0, got %x", r)
	}

	// negative zero
	if r := Float16(0x8000).Float64(); !math.IsInf(1/r, -1) {
		t.Errorf("expected -0, got %x", r)
	}
}
