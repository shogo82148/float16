package float16

import (
	"math"
	"math/bits"
)

const (
	uvnan      = 0x7e01     // "not-a-number"
	uvinf      = 0x7c00     // infinity
	uvneginf   = 0xfc00     // negative infinity
	uvone      = 0x3c00     // one
	mask16     = 0x1f       // mask for exponent
	shift16    = 16 - 5 - 1 // shift for exponent
	bias16     = 15         // bias for exponent
	signMask16 = 1 << 15    // mask for sign bit
	fracMask16 = 1<<shift16 - 1
)

const (
	mask32     = 0xff       // mask for exponent
	shift32    = 32 - 8 - 1 // shift for exponent
	bias32     = 127        // bias for exponent
	signMask32 = 1 << 31    // mask for sign bit
	fracMask32 = 1<<shift32 - 1
)

const (
	mask64     = 0x7ff       // mask for exponent
	shift64    = 64 - 11 - 1 // shift for exponent
	bias64     = 1023        // bias for exponent
	signMask64 = 1 << 63     // mask for sign bit
	fracMask64 = 1<<shift64 - 1
)

// Float16 represents a 16-bit floating point number.
type Float16 uint16

// NaN returns an IEEE 754 “not-a-number” value.
func NaN() Float16 {
	return Float16(uvnan)
}

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf(sign int) Float16 {
	if sign >= 0 {
		return Float16(uvinf)
	} else {
		return Float16(uvneginf)
	}
}

// FromBits returns the floating point number corresponding
// the IEEE 754 binary representation b.
func FromBits(b uint16) Float16 {
	return Float16(b)
}

func FromFloat64(f float64) Float16 {
	b := math.Float64bits(f)
	sign := uint16((b & signMask64) >> (64 - 16))
	exp := (b >> shift64) & mask64
	frac := b & fracMask64
	if exp == mask64 {
		if frac == 0 {
			// infinity or negative infinity
			return Float16(sign | (mask16 << shift16))
		} else {
			// NaN
			return Float16(uvnan)
		}
	}

	// round
	frac += 1 << (shift64 - shift16 - 1)
	if frac > fracMask64 {
		exp += 1
		frac = (frac >> 1) & fracMask64
	}

	// normal number
	exp16 := exp - (bias64 - bias16)
	frac16 := frac >> (shift64 - shift16)
	return Float16(sign | uint16(exp16<<shift16) | uint16(frac16))
}

// Float32 returns the float32 representation of f.
func (f Float16) Float32() float32 {
	sign := uint32(f&signMask16) << (32 - 16)
	exp := uint32(f>>shift16) & mask16
	frac := uint32(f & fracMask16)

	if exp == 0 {
		// subnormal number
		if frac == 0 {
			exp = 0
		} else {
			l := bits.Len32(frac)
			frac = (frac << (shift16 - l + 1)) & fracMask16
			exp = bias32 - (bias16 + shift16 - 1)
		}
	} else if exp == mask16 {
		// infinity or NaN
		exp = mask32
	} else {
		// normal number
		exp += bias32 - bias16
	}
	return math.Float32frombits(sign | (exp << shift32) | (frac << (shift32 - shift16)))
}

// Float64 returns the float64 representation of f.
func (f Float16) Float64() float64 {
	sign := uint64(f&signMask16) << (64 - 16)
	exp := uint64(f>>shift16) & mask16
	frac := uint64(f & fracMask16)

	if exp == 0 {
		// subnormal number
		l := bits.Len64(frac)
		if l == 0 {
			exp = 0
		} else {
			frac = (frac << (shift16 - l + 1)) & fracMask16
			exp = bias64 - (bias16 + shift16 - 1)
		}
	} else if exp == mask16 {
		// infinity or NaN
		exp = mask64
	} else {
		// normal number
		exp += bias64 - bias16
	}
	return math.Float64frombits(sign | (exp << shift64) | (frac << (shift64 - shift16)))
}
