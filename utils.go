package main

import (
	"io"
	"math/rand"
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

func fastRandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.Intn(len(letters)*len(letters)) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[(seed+i)%len(letters)]
	}
	return string(b)
}
