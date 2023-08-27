package float16

import (
	"math/bits"
)

// Mul returns the IEEE 754 binary64 product of a and b.
func (a Float16) Mul(b Float16) Float16 {
	signA := a & signMask16
	expA := int((a>>shift16)&mask16) - bias16
	signB := b & signMask16
	expB := int((b>>shift16)&mask16) - bias16

	if expA == mask16-bias16 {
		if a&fracMask16 == 0 {
			// a is infinity
			if expB == mask16-bias16 && b&fracMask16 != 0 {
				// b is NaN, the result is NaN
				return b
			} else if expB == -bias16 && b&fracMask16 == 0 {
				// b is zero, the result is NaN
				return Float16(uvnan)
			} else {
				// otherwise the result is infinity
				return a ^ signB
			}
		} else {
			// a is NaN
			return a
		}
	}

	if expB == mask16-bias16 {
		if b&fracMask16 == 0 {
			// b is infinity
			if expA == -bias16 && a&fracMask16 == 0 {
				// a is zero, the result is NaN
				return Float16(uvnan)
			} else {
				// NaN check is done above
				// so a is not zero nor NaN. the result is infinity
				return b ^ signA
			}
		} else {
			// b is NaN
			return b
		}
	}

	sign := signA ^ signB

	var fracA uint32
	if expA == -bias16 {
		// a is subnormal
		fracA = uint32(a & fracMask16)
		l := bits.Len32(fracA)
		fracA <<= shift16 - l + 1
		expA = -(bias16 + shift16) + l
	} else {
		// a is normal
		fracA = uint32(a&fracMask16) | (1 << shift16)
	}

	var fracB uint32
	if expB == -bias16 {
		// b is subnormal
		fracB = uint32(b & fracMask16)
		l := bits.Len32(fracB)
		fracB <<= shift16 - l + 1
		expB = -(bias16 + shift16) + l
	} else {
		// b is normal
		fracB = uint32(b&fracMask16) | (1 << shift16)
	}

	exp := expA + expB
	frac := fracA * fracB
	shift := bits.Len32(frac) - (shift16 + 1)
	exp += shift - shift16

	if exp < -(bias16 + shift16) {
		// underflow
		return sign
	} else if exp <= -bias16 {
		// the result is subnormal
		shift := shift16 - (expA + expB + bias16) + 1
		frac += (1<<(shift-1) - 1) + ((frac >> shift) & 1) // round to nearest even
		frac >>= shift
		return sign | Float16(frac)
	}

	exp = expA + expB + bias16
	frac += (1<<(shift-1) - 1) + ((frac >> shift) & 1) // round to nearest even
	shift = bits.Len32(frac) - (shift16 + 1)
	exp += shift - shift16
	if exp >= mask16 {
		// overflow
		return sign | (mask16 << shift16)
	}
	frac >>= shift
	frac &= fracMask16
	return sign | Float16(exp<<shift16) | Float16(frac)
}

// Add returns the IEEE 754 binary64 sum of a and b.
func (a Float16) Add(b Float16) Float16 {
	if a.IsNaN() || b.IsNaN() {
		return uvnan
	}
	if a == 0x8000 { // a is negative zero
		return b
	}
	if a == uvinf {
		if b == uvneginf {
			// +inf + -inf = NaN
			return NaN()
		}
		return uvinf // +inf + anything = +inf
	}
	if a == uvneginf {
		if b == uvinf {
			// -inf + +inf = NaN
			return NaN()
		}
		return uvneginf // -inf + anything = -inf
	}

	fixA := a.fix24()
	fixB := b.fix24()
	return (fixA + fixB).Float16()
}

// Sub returns the IEEE 754 binary64 difference of a and b.
func (a Float16) Sub(b Float16) Float16 {
	return a.Add(b ^ signMask16)
}

// fix24 is a fixed-point number with 24 bits of precision.
type fix24 int64

const fix24inf = fix24(1 << 61)

func (f Float16) fix24() fix24 {
	var ret fix24
	exp := uint32(f>>shift16) & mask16
	frac := uint32(f & fracMask16)
	if exp == 0 {
		// subnormal number
		ret = fix24(frac)
	} else if exp == mask16 {
		// infinity or NaN
		ret = fix24inf
	} else {
		// normal number
		ret = fix24(frac|(1<<shift16)) << (exp - 1)
	}
	sign := uint32(f & signMask16)
	if sign != 0 {
		ret = -ret
	}
	return ret
}

func (f fix24) Float16() Float16 {
	if f == 0 {
		return 0
	}

	var sign uint16
	if f < 0 {
		sign = signMask16
		f = -f
	}
	l := bits.Len64(uint64(f))
	if l <= shift16 {
		// subnormal number
		return Float16(sign | uint16(f))
	}
	shift := l - shift16 - 1
	if shift > 0 {
		f += (1<<(shift-1) - 1) + ((f >> shift) & 1) // round to nearest even
		l = bits.Len64(uint64(f))
	}

	exp := uint16(l) - shift16
	if exp >= mask16 {
		// overflow
		return Float16(sign | (mask16 << shift16))
	}
	frac := uint16(f>>(exp-1)) & fracMask16
	return Float16(sign | (exp << shift16) | frac)
}
