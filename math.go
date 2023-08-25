package float16

import "math/bits"

// Mul returns the IEEE 754 binary64 product of a and b.
func (a Float16) Mul(b Float16) Float16 {
	signA := a & signMask16
	expA := int((a>>shift16)&mask16) - bias16
	signB := b & signMask16
	expB := int((b>>shift16)&mask16) - bias16

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
		fracB := uint32(b & fracMask16)
		l := bits.Len32(fracB)
		fracB <<= shift16 - l + 1
		expB = -(bias16 + shift16) + l
	} else {
		// b is normal
		fracB = uint32(b&fracMask16) | (1 << shift16)
	}

	exp := expA + expB + bias16
	frac := fracA * fracB
	frac += (1<<(shift16-1) - 1) + ((frac >> shift16) & 1)
	if frac >= 1<<(2*shift16+1) {
		exp++
		frac >>= 1
	}
	frac = (frac >> shift16) & fracMask16
	return sign | Float16(exp<<shift16) | Float16(frac)
}
