package fil

import (
	"math/big"
	"testing"
)

func TestFILtoAttoFIL(t *testing.T) {
	tests := []struct {
		Value    float64
		Expected *big.Int
	}{
		{
			Value:    1.234,
			Expected: big.NewInt(1234000000000000000),
		},
		{
			Value:    1.234000000123,
			Expected: big.NewInt(1234000000123000064),
		},
	}

	for i, test := range tests {
		ret := FILtoAttoFIL(test.Value)
		if ret.Cmp(test.Expected) != 0 {
			t.Errorf("Test %d: got %s, want %s", i, ret, test.Expected)
		}
	}
}

func TestAttoFILToFIL(t *testing.T) {
	tests := []struct {
		Value    *big.Int
		Expected float64
	}{
		{
			Value:    big.NewInt(1234000000000000000),
			Expected: 1.234,
		},
		{
			Value:    big.NewInt(1234000000000000111),
			Expected: 1.234000000000000111,
		},
	}

	for i, test := range tests {
		ret := AttoFILToFIL(test.Value)
		if ret != test.Expected {
			t.Errorf("Test %d: got %f, want %f", i, ret, test.Expected)
		}
	}
}
