package main

import (
	"fmt"
	"io/fs"
	"strings"
	"unicode"

	"github.com/gobwas/glob"
)

func newGlobMatcher(ptn string) (func(string) bool, error) {
	g, err := glob.Compile(ptn)
	if err != nil {
		return nil, fmt.Errorf("inalid glob pattern: %w", err)
	}
	return g.Match, nil
}

func newGlobFileMatcher(ptn string) (func(fs.FileInfo) bool, error) {
	m, err := newGlobMatcher(ptn)
	if err != nil {
		return nil, err
	}
	return func(fi fs.FileInfo) bool {
		return m(fi.Name())
	}, nil
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
