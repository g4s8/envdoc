package main

import (
	"io/fs"
	"testing"
	"time"
)

type fakeFileInfo struct {
	name string
}

func (fi fakeFileInfo) Name() string {
	return fi.name
}

func (fi fakeFileInfo) Size() int64 {
	panic("Size() not implemented")
}

func (fi fakeFileInfo) Mode() fs.FileMode {
	panic("Mode() not implemented")
}

func (fi fakeFileInfo) ModTime() time.Time {
	panic("ModTime() not implemented")
}

func (fi fakeFileInfo) IsDir() bool {
	panic("IsDir() not implemented")
}

func (fi fakeFileInfo) Sys() interface{} {
	panic("Sys() not implemented")
}

func globTesteer(matcher func(string) bool, targets map[string]bool) func(*testing.T) {
	return func(t *testing.T) {
		for target, expected := range targets {
			if matcher(target) != expected {
				t.Errorf("unexpected result for %q: got %v, want %v", target, !expected, expected)
			}
		}
	}
}

var globTestTargets = map[string]bool{
	"main.go":         true,
	"main_test.go":    true,
	"utils.go":        true,
	"utils_test.java": false,
	"file.txt":        false,
	"test.go.txt":     false,
	"cfg/Config.go":   true,
}

func TestGlobMatcher(t *testing.T) {
	m, err := newGlobMatcher("*.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("match", globTesteer(m, globTestTargets))

	t.Run("error", func(t *testing.T) {
		_, err := newGlobMatcher("[")
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}

func TestGlobFileMatcher(t *testing.T) {
	m, err := newGlobFileMatcher("*.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fileWrapper := func(name string) bool {
		fi := fakeFileInfo{name}
		return m(fi)
	}

	t.Run("match", globTesteer(fileWrapper, globTestTargets))
	t.Run("error", func(t *testing.T) {
		_, err := newGlobFileMatcher("[")
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
	})
}

func TestCamelToSnake(t *testing.T) {
	tests := map[string]string{
		"CamelCase":         "CAMEL_CASE",
		"camelCase":         "CAMEL_CASE",
		"camel":             "CAMEL",
		"Camel":             "CAMEL",
		"camel_case":        "CAMEL_CASE",
		"camel_case_":       "CAMEL_CASE_",
		"camel_case__":      "CAMEL_CASE__",
		"camelCase_":        "CAMEL_CASE_",
		"camelCase__":       "CAMEL_CASE__",
		"camel_case__snake": "CAMEL_CASE__SNAKE",
		"":                  "",
		" ":                 " ",
		"_":                 "_",
		"_A_":               "_A_",
		"ABBRFoo":           "ABBR_FOO",
		"FOO_BAR":           "FOO_BAR",
		"ЮниКод":            "ЮНИ_КОД",
		"ՅունիԿոդ":          "ՅՈՒՆԻ_ԿՈԴ", 
	}

	for input, expected := range tests {
		if got := camelToSnake(input); got != expected {
			t.Errorf("unexpected result for %q: got %q, want %q", input, got, expected)
		}
	}
}
