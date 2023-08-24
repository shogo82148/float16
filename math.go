package float16

// Mul returns the IEEE 754 binary64 product of a and b.
func (a Float16) Mul(b Float16) Float16 {
	signA := a & signMask16
	expA := int((a>>shift16)&mask16) - bias16
	signB := b & signMask16
	expB := int((b>>shift16)&mask16) - bias16

	sign := signA ^ signB
	exp := expA + expB + bias16

	fracA := int(a&fracMask16) | (1 << shift16)
	fracB := int(b&fracMask16) | (1 << shift16)

	frac := fracA * fracB
	frac += (1<<(shift16-1) - 1) + ((frac >> shift16) & 1)
	frac = (frac >> shift16) & fracMask16
	return sign | Float16(exp<<shift16) | Float16(frac)
}
