package float16

import (
	"math/bits"
	"strconv"
)

func (x Float16) String() string {
	return x.Format('g', -1)
}

func (x Float16) Format(fmt byte, prec int) string {
	return string(x.Append(make([]byte, 0, 8), fmt, prec))
}

func (x Float16) Append(buf []byte, fmt byte, prec int) []byte {
	switch fmt {
	case 'x', 'X':
		switch {
		case x.IsNaN():
			return append(buf, "NaN"...)
		case x == uvinf:
			return append(buf, "+Inf"...)
		case x == uvneginf:
			return append(buf, "-Inf"...)
		}

		if x&signMask16 != 0 {
			buf = append(buf, '-')
		}
		buf = append(buf, '0', fmt)

		// normalize x
		exp := int((x >> shift16) & mask16)
		frac := uint16(x & fracMask16)
		if exp == 0 && frac == 0 {
			// zero
			buf = append(buf, '0')
			if prec >= 1 {
				buf = append(buf, '.')
				for i := 0; i < prec; i++ {
					buf = append(buf, '0')
				}
			}
			buf = append(buf, fmt-('x'-'p'))
			buf = append(buf, '+', '0', '0')
			return buf
		}
		if exp == 0 {
			// subnormal number
			l := bits.Len32(uint32(frac))
			frac <<= shift16 - l + 1
			exp = -(bias16 + shift16) + l
		} else {
			// normal number
			frac |= 1 << shift16
			exp -= bias16
		}

		switch prec {
		case -1:
			if frac&0x3ff == 0 {
				buf = append(buf, '1')
			} else if frac&0x3f == 0 {
				buf = append(buf, '1', '.')
				buf = append(buf, nibble(fmt, frac>>6))
			} else if frac&0x3 == 0 {
				buf = append(buf, '1', '.')
				buf = append(buf, nibble(fmt, frac>>6))
				buf = append(buf, nibble(fmt, frac>>2))
			} else {
				buf = append(buf, '1', '.')
				buf = append(buf, nibble(fmt, frac>>6))
				buf = append(buf, nibble(fmt, frac>>2))
				buf = append(buf, nibble(fmt, frac<<2))
			}

		case 0:
			// round to nearest even
			frac += 1 << (shift16 - 1)
			if frac >= 1<<(shift16+1) {
				exp++
				frac >>= 1
			}

			buf = append(buf, '1')

		case 1:
			// round to nearest even
			frac += 0x1f + (frac>>6)&1
			if frac >= 1<<(shift16+1) {
				exp++
				frac >>= 1
			}

			buf = append(buf, '1', '.')
			buf = append(buf, nibble(fmt, frac>>6))

		case 2:
			// round to nearest even
			frac += 1 + (frac>>2)&1
			if frac >= 1<<(shift16+1) {
				exp++
				frac >>= 1

			}

			buf = append(buf, '1', '.')
			buf = append(buf, nibble(fmt, frac>>6))
			buf = append(buf, nibble(fmt, frac>>2))

		default:
			if prec < 0 {
				panic("invalid precision")
			}
			buf = append(buf, '1', '.')
			buf = append(buf, nibble(fmt, frac>>6))
			buf = append(buf, nibble(fmt, frac>>2))
			buf = append(buf, nibble(fmt, frac<<2))
			for i := 3; i < prec; i++ {
				buf = append(buf, '0')
			}
		}

		buf = append(buf, fmt-('x'-'p'))
		if exp >= 0 {
			buf = append(buf, '+')
		} else {
			buf = append(buf, '-')
			exp = -exp
		}
		buf = append(buf, byte(exp/10)+'0', byte(exp%10)+'0')
		return buf
	}
	return strconv.AppendFloat(buf, x.Float64(), fmt, prec, 32)
}

func nibble(fmt byte, x uint16) byte {
	x &= 0xf
	if x < 10 {
		return '0' + byte(x)
	}
	return ('A' + byte(x-10)) | (fmt & ('a' - 'A'))
}
