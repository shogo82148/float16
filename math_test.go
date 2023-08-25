package float16

import (
	"runtime"
	"testing"
)

func TestMul(t *testing.T) {
	tests := []struct {
		a, b float64
	}{
		// normal * normal = normal
		{1, 1}, // 1 * 1 = 1
		{1, 2}, // 1 * 2 = 2
		{0x1.f44p-01, 0x1.fa8p-01},
		{0x1.efp-01, 0x1.08cp+00},

		// subnormal * normal = normal
		{0x1p-15, 2}, // 0x1p-15 * 2  = 0x1p-14

		// normal * subnormal = normal
		{2, 0x1p-15}, // 0x1p-15 * 2 = 0x1p-14

		// subnormal * normal = subnormal
		{0x1p-24, 2}, // 0x1p-24 * 2 = 0x1p-23

		// subnormal * subnormal = subnormal
		{0, 0}, // 0 * 0 = 0
	}
	for _, tt := range tests {
		fa := FromFloat64(tt.a)
		fb := FromFloat64(tt.b)
		fr := FromFloat64(tt.a * tt.b)
		fc := fa.Mul(fb)
		if fc != fr {
			t.Errorf("%x * %x: expected %x (0x%04x), got %x (0x%04x)", tt.a, tt.b, fr.Float64(), fr, fc.Float64(), fc)
		}
	}
}

func BenchmarkMul(b *testing.B) {
	a := Float16(0x3c00)
	bb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(a.Mul(bb))
	}
}

func FuzzMul(f *testing.F) {
	f.Add(uint16(0x3c00), uint16(0x3c00))

	f.Fuzz(func(t *testing.T, a, b uint16) {
		fa := Float16(a)
		fb := Float16(b)
		fc := fa.Mul(fb)

		want := FromFloat64(fa.Float64() * fb.Float64())
		if fc != want {
			t.Errorf("%x * %x: expected %x, got %x", fa.Float64(), fb.Float64(), want.Float64(), fc.Float64())
		}
	})
}
