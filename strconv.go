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
		if x&^signMask16 == 0 {
			if fmt == 'x' {
				return append(buf, '0', 'p', '+', '0', '0')
			} else {
				return append(buf, '0', 'P', '+', '0', '0')
			}
		}
	}
	return strconv.AppendFloat(buf, x.Float64(), fmt, prec, 32)
}
