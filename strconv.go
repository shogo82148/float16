package float16

import (
	"strconv"
)

func (x Float16) String() string {
	return x.Format('g', -1)
}

func (x Float16) Format(fmt byte, prec int) string {
	return string(x.Append(make([]byte, 0, 8), fmt, prec))
}

func (x Float16) Append(buf []byte, fmt byte, prec int) []byte {
	switch {
	case x.IsNaN():
		return append(buf, "NaN"...)
	case x == uvinf:
		return append(buf, "+Inf"...)
	case x == uvneginf:
		return append(buf, "-Inf"...)
	}

	switch fmt {
	case 'b':
		return x.appendBin(buf)
	case 'f':
		return x.appendFloat(buf, fmt, prec)
	case 'x', 'X':
		return x.appendHex(buf, fmt, prec)
	}
	return strconv.AppendFloat(buf, x.Float64(), fmt, prec, 32)
}

func (x Float16) appendBin(buf []byte) []byte {
	if x&signMask16 != 0 {
		buf = append(buf, '-')
	}
	exp := int(x>>shift16&mask16) - bias16
	frac := x & fracMask16

	if exp == -bias16 {
		exp++
	} else {
		frac |= 1 << shift16
	}
	exp -= shift16

	switch {
	case frac >= 10000:
		buf = append(buf, byte((frac/10000)%10)+'0')
		fallthrough
	case frac >= 1000:
		buf = append(buf, byte((frac/1000)%10)+'0')
		fallthrough
	case frac >= 100:
		buf = append(buf, byte((frac/100)%10)+'0')
		fallthrough
	case frac >= 10:
		buf = append(buf, byte((frac/10)%10)+'0')
		fallthrough
	default:
		buf = append(buf, byte(frac%10)+'0')
	}

	buf = append(buf, 'p')
	if exp >= 0 {
		buf = append(buf, '+')
	} else {
		buf = append(buf, '-')
		exp = -exp
	}

	switch {
	case exp >= 10:
		buf = append(buf, byte(exp/10)+'0')
		fallthrough
	default:
		buf = append(buf, byte(exp%10)+'0')
	}
	return buf
}

func (x Float16) appendFloat(buf []byte, fmt byte, prec int) []byte {
	f := x.fix24()
	if x&signMask16 != 0 {
		buf = append(buf, '-')
		f = -f
	}
	i := int(f >> 24)
	f &= 0xffffff

	switch {
	case i >= 10000:
		buf = append(buf, byte((i/10000)%10)+'0')
		fallthrough
	case i >= 1000:
		buf = append(buf, byte((i/1000)%10)+'0')
		fallthrough
	case i >= 100:
		buf = append(buf, byte((i/100)%10)+'0')
		fallthrough
	case i >= 10:
		buf = append(buf, byte((i/10)%10)+'0')
		fallthrough
	default:
		buf = append(buf, byte(i%10)+'0')
	}

	return buf
}

func (x Float16) appendHex(buf []byte, fmt byte, prec int) []byte {
	// normalize x
	sign, exp, frac := x.split()
	if sign != 0 {
		buf = append(buf, '-')
	}
	buf = append(buf, '0', fmt) // 0x or 0X
	if x&^signMask16 == 0 {
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

func nibble(fmt byte, x uint16) byte {
	x &= 0xf
	if x < 10 {
		return '0' + byte(x)
	}
	return ('A' + byte(x-10)) | (fmt & ('a' - 'A'))
}
