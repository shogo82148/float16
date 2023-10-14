package float16

import (
	"log"
	"math/bits"
	"strconv"

	"github.com/shogo82148/int128"
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
	case 'x', 'X':
		return x.appendHex(buf, fmt, prec)
	case 'f':
		return x.appendDec(buf, fmt, prec)
	}

	// TODO: shortest representation
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

func (x Float16) appendDec(buf []byte, fmt byte, prec int) []byte {
	if prec >= 0 {
		const five24 = 59604644775390625 // = 5^24
		ten := int128.Uint128{L: 10}

		var dec24 int128.Uint128
		fix := x.fix24()
		dec24.H, dec24.L = bits.Mul64(uint64(fix), five24)

		// round to nearest even
		if prec < 24 {
			n := int128.Uint128{L: 1}
			for i := 0; i < 24-prec; i++ {
				n = n.Mul(ten)
			}
			n2 := n.Rsh(1)
			div, mod := dec24.DivMod(n)
			if mod.Cmp(n2) > 0 {
				// round up
				dec24 = div.Add(n)
			} else if mod.Cmp(n2) == 0 {
				// round to even
				if div.L&1 != 0 {
					dec24 = div.Add(n)
				}
			}
		}

		// convert to decimal
		var data [24]byte
		for i := 0; i < 24; i++ {
			var mod int128.Uint128
			dec24, mod = dec24.DivMod(ten)
			data[i] = byte(mod.L)
		}

		// convert integer part
		switch {
		case dec24.L >= 10000:
			buf = append(buf, byte((dec24.L/10000)%10)+'0')
			fallthrough
		case dec24.L >= 1000:
			buf = append(buf, byte((dec24.L/1000)%10)+'0')
			fallthrough
		case dec24.L >= 100:
			buf = append(buf, byte((dec24.L/100)%10)+'0')
			fallthrough
		case dec24.L >= 10:
			buf = append(buf, byte((dec24.L/10)%10)+'0')
			fallthrough
		default:
			buf = append(buf, byte(dec24.L%10)+'0')
		}
		if prec == 0 {
			return buf
		}

		buf = append(buf, '.')

		// convert fractional part
		var i int
		for i = 0; i < prec && i < len(data); i++ {
			buf = append(buf, data[23-i]+'0')
		}
		for ; i < prec; i++ {
			buf = append(buf, '0')
		}

		return buf
	}

	var exact, lower, upper int128.Uint128
	exp := int(x >> shift16 & mask16)
	frac := uint64(x & fracMask16)
	if exp == 0 {
		// subnormal number
		exact.L = frac * 2
		lower.L = exact.L - 1
		upper.L = exact.L + 1
	} else {
		// normal number
		exact.L = (frac | (1 << shift16)) << exp
		lower.L = exact.L - (1 << (exp - 1))
		upper.L = exact.L + (1 << (exp - 1))
	}

	const five25 = 298023223876953125 // = 5^25
	exact.H, exact.L = bits.Mul64(exact.L, five25)
	lower.H, lower.L = bits.Mul64(lower.L, five25)
	upper.H, upper.L = bits.Mul64(upper.L, five25)

	log.Print(exact.String())
	log.Print(lower.String())
	log.Print(upper.String())

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
