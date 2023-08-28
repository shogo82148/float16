package float16

import (
	"cmp"
	"errors"
	"math"
	"runtime"
	"testing"
	"testing/quick"
)

func TestMul(t *testing.T) {
	tests := []struct {
		a, b float64
	}{
		// normal * normal => normal
		{1, 1}, // 1 * 1 = 1
		{1, 2}, // 1 * 2 = 2
		{0x1.f44p-01, 0x1.fa8p-01},
		{0x1.efp-01, 0x1.08cp+00},
		{0x1p-14, 0x1p-10},
		{0x1p-10, 0x1p-14},

		// normal * normal => subnormal
		{-0x1.c14p-12, -0x1.32cp-12},

		// subnormal * normal => normal
		{0x1p-15, 2}, // 0x1p-15 * 2  = 0x1p-14

		// normal * subnormal => normal
		{2, 0x1p-15}, // 0x1p-15 * 2 = 0x1p-14

		// subnormal * normal => subnormal
		{0x1p-24, 2}, // 0x1p-24 * 2 = 0x1p-23
		{-0x1.1d8p-09, 0x1.07p-13},

		// subnormal * subnormal => subnormal
		{0, 0}, // 0 * 0 = 0
		{negZero, 0},

		// underflow
		{0x1p-14, 0x1p-11},
		{0x1p-11, 0x1p-14},

		// overflow
		{-0x1.b9p+06, -0x1.e24p+09},

		// Infinity * 0 => NaN
		{math.Inf(1), 0},
		{0, math.Inf(1)},

		// Infinity * anything => Infinity
		{math.Inf(1), 1},
		{math.Inf(1), -1},
		{math.Inf(-1), 1},
		{math.Inf(-1), -1},

		// anything * Infinity => Infinity
		{1, math.Inf(1)},
		{-1, math.Inf(1)},
		{1, math.Inf(-1)},
		{-1, math.Inf(-1)},

		// NaN * anything = NaN
		// anything * NaN = NaN
		{math.NaN(), 1},
		{math.NaN(), math.Inf(1)},
		{1, math.NaN()},
		{math.Inf(1), math.NaN()},
		{math.NaN(), math.NaN()},
	}
	for _, tt := range tests {
		fa := FromFloat64(tt.a)
		if !fa.IsNaN() && fa.Float64() != tt.a {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.a)
		}
		fb := FromFloat64(tt.b)
		if !fb.IsNaN() && fb.Float64() != tt.b {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.b)
		}
		fr := FromFloat64(tt.a * tt.b)
		fc := fa.Mul(fb)
		if fc != fr {
			t.Errorf("%x * %x: expected %x (0x%04x), got %x (0x%04x)", tt.a, tt.b, fr.Float64(), fr, fc.Float64(), fc)
		}
	}
}

func BenchmarkMul(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(fa.Mul(fb))
	}
}

func BenchmarkMul2(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		fc := fa.Float64() * fb.Float64()
		runtime.KeepAlive(FromFloat64(fc))
	}
}

func TestMulQuick(t *testing.T) {
	f := func(a, b uint16) uint16 {
		fa := Float16(a)
		fb := Float16(b)
		fc := fa.Mul(fb)
		if fc.IsNaN() {
			return NaN().Bits()
		}
		return fc.Bits()
	}

	g := func(a, b uint16) uint16 {
		fa := Float16(a).Float64()
		fb := Float16(b).Float64()
		fc := fa * fb // This calculation does not cause any rounding.
		return FromFloat64(fc).Bits()
	}

	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 100,
	}); err != nil {
		var checkErr *quick.CheckEqualError
		if errors.As(err, &checkErr) {
			a := checkErr.In[0].(uint16)
			b := checkErr.In[1].(uint16)
			c1 := checkErr.Out1[0].(uint16)
			c2 := checkErr.Out2[0].(uint16)

			fa := FromBits(a).Float64()
			fb := FromBits(b).Float64()
			fc1 := FromBits(c1).Float64()
			fc2 := FromBits(c2).Float64()

			t.Errorf("%x * %x: got %x, expected %x", fa, fb, fc1, fc2)
		}
		t.Error(err)
	}
}

func TestQuo(t *testing.T) {
	tests := []struct {
		a, b, r float64
	}{
		{1, 2, 0.5},
		{1, 3, 0x1.554p-02},
		{1, 5, 0x1.998p-03},
		{0x1.46p-11, 0x1.13cp+02, 0x1.2ecp-13},
		{0x1p+15, 0.5, math.Inf(1)},
		{0x1p-14, 2, 0x1p-15},
		{0x1p-24, 0x1p-24, 1},

		{math.NaN(), 1, math.NaN()},
		{1, math.NaN(), math.NaN()},

		{math.Inf(1), 1, math.Inf(1)},
		{math.Inf(1), 0, math.Inf(1)},
		{math.Inf(1), math.Inf(1), math.NaN()},
		{1, math.Inf(1), 0},
		{1, math.Inf(-1), negZero},
		{1, 0, math.Inf(1)},

		// 0 / 0
		{negZero, 0, math.Inf(-1)},
		{0, negZero, math.Inf(-1)},
		{0, 0, math.NaN()},
		{negZero, negZero, math.NaN()},
	}

	for _, tt := range tests {
		fa := FromFloat64(tt.a)
		if !fa.IsNaN() && fa.Float64() != tt.a {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.a)
		}
		fb := FromFloat64(tt.b)
		if !fb.IsNaN() && fb.Float64() != tt.b {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.b)
		}
		fr := FromFloat64(tt.r)
		if !fr.IsNaN() && fr.Float64() != tt.r {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.b)
		}
		fc := fa.Quo(fb)
		if fc != fr {
			t.Errorf("%x / %x: expected %x (0x%04x), got %x (0x%04x)", tt.a, tt.b, fr.Float64(), fr, fc.Float64(), fc)
		}
	}
}

// func TestQuoQuick(t *testing.T) {
// 	f := func(a, b uint16) uint16 {
// 		fa := Float16(a)
// 		fb := Float16(b)
// 		fc := fa.Quo(fb)
// 		if fc.IsNaN() {
// 			return NaN().Bits()
// 		}
// 		return fc.Bits()
// 	}

// 	g := func(a, b uint16) uint16 {
// 		fa := Float16(a).Float64()
// 		if math.IsNaN(fa) {
// 			return NaN().Bits()
// 		}
// 		bigA := new(big.Float).SetFloat64(fa)
// 		fb := Float16(b).Float64()
// 		if math.IsNaN(fb) {
// 			return NaN().Bits()
// 		}
// 		bigB := new(big.Float).SetFloat64(fb)
// 		fc := new(big.Float).SetPrec(11).Quo(bigA, bigB)
// 		f64, _ := fc.Float64()
// 		return FromFloat64(f64).Bits()
// 	}

// 	if err := quick.CheckEqual(f, g, &quick.Config{
// 		MaxCountScale: 100,
// 	}); err != nil {
// 		var checkErr *quick.CheckEqualError
// 		if errors.As(err, &checkErr) {
// 			a := checkErr.In[0].(uint16)
// 			b := checkErr.In[1].(uint16)
// 			c1 := checkErr.Out1[0].(uint16)
// 			c2 := checkErr.Out2[0].(uint16)

// 			fa := FromBits(a).Float64()
// 			fb := FromBits(b).Float64()
// 			fc1 := FromBits(c1).Float64()
// 			fc2 := FromBits(c2).Float64()

// 			t.Errorf("%x / %x: got %x, expected %x", fa, fb, fc1, fc2)
// 		}
// 		t.Error(err)
// 	}
// }

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b float64
	}{
		// normal + normal => normal
		{1, 1},
		{1, 2},
		{2, -1},
		{1, 0x1p-11},
		{1 + 0x1p-10, 0x1p-11},
		{0x1p+15, 0x1p+15}, // overflow

		// subnormal + subnormal => subnormal
		{0, 0},
		{negZero, 0},
		{negZero, negZero},
		{0, negZero},
		{0x1p-24, 0x1p-24},

		// infinity + anything => infinity
		{math.Inf(-1), 0x1p+14},
		{math.Inf(1), math.Inf(1)},

		// anything + infinity => infinity
		{0x1p+14, math.Inf(-1)},

		// infinity - infinity => infinity
		{math.Inf(1), math.Inf(-1)},
		{math.Inf(-1), math.Inf(1)},

		// NaN + anything = NaN
		{math.NaN(), -0x1p+14},

		// anything + NaN = NaN
		{-0x1p+14, math.NaN()},
	}
	for _, tt := range tests {
		fa := FromFloat64(tt.a)
		if !fa.IsNaN() && fa.Float64() != tt.a {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.a)
		}
		fb := FromFloat64(tt.b)
		if !fb.IsNaN() && fb.Float64() != tt.b {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.b)
		}
		fr := FromFloat64(tt.a + tt.b)
		fc := fa.Add(fb)
		if fc != fr {
			t.Errorf("%x + %x: expected %x (0x%04x), got %x (0x%04x)", tt.a, tt.b, fr.Float64(), fr, fc.Float64(), fc)
		}

		fr = FromFloat64(tt.b - tt.a)
		fc = fb.Sub(fa)
		if fc != fr {
			t.Errorf("%x - %x: expected %x (0x%04x), got %x (0x%04x)", tt.b, tt.a, fr.Float64(), fr, fc.Float64(), fc)
		}
	}
}

func TestAddQuick(t *testing.T) {
	f := func(a, b uint16) uint16 {
		fa := Float16(a)
		fb := Float16(b)
		fc := fa.Add(fb)
		if fc.IsNaN() {
			return NaN().Bits()
		}
		return fc.Bits()
	}

	g := func(a, b uint16) uint16 {
		fa := Float16(a).Float64()
		fb := Float16(b).Float64()
		fc := fa + fb // This calculation does not cause any rounding.
		return FromFloat64(fc).Bits()
	}

	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 100,
	}); err != nil {
		var checkErr *quick.CheckEqualError
		if errors.As(err, &checkErr) {
			a := checkErr.In[0].(uint16)
			b := checkErr.In[1].(uint16)
			c1 := checkErr.Out1[0].(uint16)
			c2 := checkErr.Out2[0].(uint16)

			fa := FromBits(a).Float64()
			fb := FromBits(b).Float64()
			fc1 := FromBits(c1).Float64()
			fc2 := FromBits(c2).Float64()

			t.Errorf("%x + %x: got %x, expected %x", fa, fb, fc1, fc2)
		}
		t.Error(err)
	}
}

func BenchmarkAdd(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(fa.Add(fb))
	}
}

func BenchmarkAdd2(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		fc := fa.Float64() + fb.Float64()
		runtime.KeepAlive(FromFloat64(fc))
	}
}

func TestCmp(t *testing.T) {
	tests := []struct {
		a, b float64
	}{
		// positive numbers
		{1, 1},
		{1, 2},
		{2, 1},

		// negative numbers
		{-1, -1},
		{-1, -2},
		{-2, -1},

		// positive and negative numbers
		{-1, 1},
		{1, -1},
		{-2, 1},
		{2, -1},

		// infinity
		{math.Inf(1), 1},
		{math.Inf(-1), -1},
		{math.Inf(1), math.Inf(1)},
		{math.Inf(1), math.Inf(-1)},
		{math.Inf(-1), math.Inf(1)},
		{math.Inf(-1), math.Inf(-1)},

		// a NaN is considered less than any non-NaN
		{math.NaN(), 0},
		{math.NaN(), 1},
		{math.NaN(), math.Inf(1)},
		{math.NaN(), math.Inf(-1)},

		// negative zero
		{0, 0},
		{negZero, 0},
		{0, negZero},
		{negZero, negZero},
	}
	for _, tt := range tests {
		fa := FromFloat64(tt.a)
		if !fa.IsNaN() && fa.Float64() != tt.a {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.a)
		}
		fb := FromFloat64(tt.b)
		if !fb.IsNaN() && fb.Float64() != tt.b {
			t.Errorf("%x + %x: invalid test case: converting %x to float16 loss data", tt.a, tt.b, tt.b)
		}
		fr := cmp.Compare(tt.a, tt.b)
		fc := fa.Cmp(fb)
		if fc != fr {
			t.Errorf("%x <=> %x: expected %d, got %d", tt.a, tt.b, fr, fc)
		}
	}
}

func TestCmpQuick(t *testing.T) {
	f := func(a, b uint16) int {
		fa := Float16(a)
		fb := Float16(b)
		return fa.Cmp(fb)
	}

	g := func(a, b uint16) int {
		fa := Float16(a).Float64()
		fb := Float16(b).Float64()
		return cmp.Compare(fa, fb)
	}

	if err := quick.CheckEqual(f, g, &quick.Config{
		MaxCountScale: 100,
	}); err != nil {
		var checkErr *quick.CheckEqualError
		if errors.As(err, &checkErr) {
			a := checkErr.In[0].(uint16)
			b := checkErr.In[1].(uint16)
			c1 := checkErr.Out1[0].(int)
			c2 := checkErr.Out2[0].(int)

			fa := FromBits(a).Float64()
			fb := FromBits(b).Float64()

			t.Errorf("%x + %x: got %d, expected %d", fa, fb, c1, c2)
		}
		t.Error(err)
	}
}

func BenchmarkCmp(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(fa.Cmp(fb))
	}
}

func BenchmarkCmp2(b *testing.B) {
	fa := Float16(0x3c00)
	fb := Float16(0x4000)
	for i := 0; i < b.N; i++ {
		c := cmp.Compare(fa.Float64(), fb.Float64())
		runtime.KeepAlive(c)
	}
}
