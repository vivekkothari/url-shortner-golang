package utils

import "strings"

const base62Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Base62Decode converts a Base62 encoded string to a long integer
func Base62Decode(url string) int64 {
	var decode int64
	for i := 0; i < len(url); i++ {
		decode = decode*62 + int64(strings.IndexByte(base62Charset, url[i]))
	}
	return decode
}

// Base62Encode converts a long integer to a Base62 encoded string
func Base62Encode(value int64) string {
	if value == 0 {
		return "000000" // Ensures a minimum length of 6
	}

	var sb strings.Builder
	for value > 0 {
		sb.WriteByte(base62Charset[value%62])
		value /= 62
	}

	// Pad with '0' to ensure at least 6 characters
	for sb.Len() < 6 {
		sb.WriteByte('0')
	}

	// Reverse the string
	runes := []rune(sb.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
