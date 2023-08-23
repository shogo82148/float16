package float16

import (
	"math"
	"math/bits"
)

type Float16 uint16

// Float16frombits returns the floating point number corresponding
// the IEEE 754 binary representation b.
func Float16frombits(b uint16) Float16 {
	return Float16(b)
}

// Float32 returns the float32 representation of f.
func (f Float16) Float32() float32 {
	sign := uint32(f&0x8000) << 16
	exp := uint32(f>>10) & 0x1f
	mant := uint32(f & 0x3ff)

	if exp == 0 {
		// subnormal number
		if mant == 0 {
			exp = 0
		} else {
			l := bits.Len32(mant)
			mant = (mant << (10 - l + 1)) & 0x3ff
			exp = 127 - 24
		}
	} else if exp == 31 {
		// infinity or NaN
		exp = 255
	} else {
		// normal number
		exp += 127 - 15
	}
	return math.Float32frombits(sign | (exp << 23) | (mant << (23 - 10)))
}

// Float64 returns the float64 representation of f.
func (f Float16) Float64() float64 {
	sign := uint64(f&0x8000) << 48
	exp := uint64(f>>10) & 0x1f
	mant := uint64(f & 0x3ff)

	if exp == 0 {
		// subnormal number
		l := bits.Len64(mant)
		if l == 0 {
			exp = 0
		} else {
			mant = (mant << (10 - l + 1)) & 0x3ff
			exp = 1023 - 24
		}
	} else if exp == 31 {
		// infinity or NaN
		exp = 2047
	} else {
		// normal number
		exp += 1023 - 15
	}
	return math.Float64frombits(sign | (exp << 52) | (mant << (52 - 10)))
}
