package float16

import (
	"math"
	"runtime"
	"testing"
)

func TestSqrt(t *testing.T) {
	tests := []struct {
		x float64
		y float64
	}{
		// special cases
		{0, 0},
		{negZero, negZero},
		{math.Inf(1), math.Inf(1)},
		{-1, math.NaN()},

		// normal numbers
		{1, 1},
		{2, 0x1.6ap+00},
		{3, 0x1.bb8p+00},
		{4, 0x2p+00},
	}

	for _, tt := range tests {
		x := FromFloat64(tt.x)
		y := Sqrt(x)
		if y.IsNaN() && math.IsNaN(tt.y) {
			continue
		}
		if y.Float64() != tt.y {
			t.Errorf("expected %x, got %x", tt.y, y.Float64())
		}
	}
}

func BenchmarkSqrt(b *testing.B) {
	r := newXorshift32()
	for i := 0; i < b.N; i++ {
		f, _ := r.Float16Pair()
		runtime.KeepAlive(Sqrt(f))
	}
}
