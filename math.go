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
		frac := fracA * fracB
		frac >>= shift16 - (exp + bias16) + 1
		return sign | Float16(frac)
	}

	exp = expA + expB + bias16
	frac += (1<<(shift-1) - 1) + ((frac >> shift) & 1)
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
