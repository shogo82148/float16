package float16

import (
	"runtime"
	"testing"
)

func TestMul(t *testing.T) {
	tests := []struct {
		a, b, r Float16
	}{
		{0x3c00, 0x3c00, 0x3c00}, // 1*1 = 1
		{0x3c00, 0x4000, 0x4000}, // 1*2 = 2
	}
	for _, tt := range tests {
		r := tt.a.Mul(tt.b)
		if r != tt.r {
			t.Errorf("%x*%x: expected %x, got %x", tt.a, tt.b, tt.r, r)
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
