package utils

import "testing"

func TestBase62Encode(t *testing.T) {
	// Test cases for Base62Encode
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "000000"},
		{1, "000001"},
		{62, "000010"},
		{123456789, "08m0Kx"},
	}

	for _, test := range tests {
		result := Base62Encode(test.input)
		if result != test.expected {
			t.Errorf("Base62Encode(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}

	for i := 0; i < 1_000_000; i++ {
		decode := Base62Decode(Base62Encode(int64(i)))
		if decode != int64(i) {
			t.Errorf("Base62Decode(Base62Encode(%d)) = %d; expected %d", i, decode, i)
		}
	}

}
