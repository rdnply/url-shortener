package baseconv

import "testing"

var (
	testBases   = []uint{0, 2, 3, 8, 16, 44, 62, 33, 56, 100}
	testNumbers = []uint{0, 4, 233, 255555555, 500000, 100000000, 9999999, 34, 55}
)

func TestBaseConvertion(t *testing.T) {
	for _, base := range testBases {
		baseconv, err := NewBaseConv(base)
		if err != nil && (base != 0 && base <= 62) {
			t.Errorf("base %v must be valid", base)
		}
		if err != nil {
			continue
		}

		for _, n := range testNumbers {
			if baseconv.Decode(baseconv.Encode(n)) != n {
				t.Errorf("incorrect convertion with base: %v", base)
			}
		}
	}
}
