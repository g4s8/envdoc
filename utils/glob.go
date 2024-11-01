package utils

import (
	"fmt"
	"io/fs"

	"github.com/gobwas/glob"
)

// un-escape -types and -files globs: '*' -> *, "foo" -> foo
// if first and last characters are quotes, remove them.
func UnescapeGlob(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

func NewGlobMatcher(ptn string) (func(string) bool, error) {
	g, err := glob.Compile(ptn)
	if err != nil {
		return nil, fmt.Errorf("inalid glob pattern: %w", err)
	}
	return g.Match, nil
}

func NewGlobFileMatcher(ptn string) (func(fs.FileInfo) bool, error) {
	m, err := NewGlobMatcher(ptn)
	if err != nil {
		return nil, err
	}
	return func(fi fs.FileInfo) bool {
		return m(fi.Name())
	}, nil
}
