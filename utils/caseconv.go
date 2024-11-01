package utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CamelToSnake(s string) string {
	const underscore = '_'
	var result strings.Builder
	result.Grow(len(s) + 5)

	var buf [utf8.UTFMax]byte
	var prev rune
	var pos int
	for i, r := range s {
		pos += utf8.EncodeRune(buf[:], r)
		// read next rune
		var next rune
		if pos < len(s) {
			next, _ = utf8.DecodeRuneInString(s[pos:])
		}
		if i > 0 && prev != underscore && r != underscore && unicode.IsUpper(r) && (unicode.IsLower(next)) {
			result.WriteRune(underscore)
		}
		result.WriteRune(unicode.ToUpper(r))
		prev = r
	}

	return result.String()
}
