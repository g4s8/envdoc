package main

import (
	"io"
	"strings"
	"unicode"
)

func closeWith(closer io.Closer, handler func(error)) {
	err := closer.Close()
	if err != nil {
		handler(err)
	}
}

func camelToSnake(s string) string {
	const underscore = '_'
	var result strings.Builder
	result.Grow(len(s) + 5)

	var prev rune
	for i, r := range s {
		if i > 0 && prev != underscore && r != underscore && unicode.IsUpper(r) {
			result.WriteRune(underscore)
		}
		result.WriteRune(unicode.ToUpper(r))
		prev = r
	}

	return result.String()
}
