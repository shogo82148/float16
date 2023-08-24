package float16

import "testing"

func TestMul(t *testing.T) {
	tests := []struct {
		a, b, r Float16
	}{
		{0x3c00, 0x3c00, 0x3c00}, // 1*1 = 1
		{0x3c00, 0x4000, 0x4000}, // 1*2 = 2
	}
	for _, tt := range tests {
		r := tt.a.Mul(tt.b)
		if r != tt.r {
			t.Errorf("%x*%x: expected %x, got %x", tt.a, tt.b, tt.r, r)
		}
	}
}
