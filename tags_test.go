package main

import (
	"slices"
	"testing"
)

func TestFieldTags(t *testing.T) {
	const src = `env:"PASSWORD,required,file" envDefault:"/tmp/password" json:"password"`
	tag := ParseFieldTag(src)
	expectAll := map[string][]string{
		"env":        {"PASSWORD", "required", "file"},
		"envDefault": {"/tmp/password"},
		"json":       {"password"},
	}
	for k, v := range expectAll {
		if got := tag.GetAll(k); !slices.Equal(got, v) {
			t.Errorf("%q: expected %q, got %q", k, v, got)
		}
	}

	expectOne := map[string]string{
		"env":        "PASSWORD",
		"envDefault": "/tmp/password",
		"json":       "password",
	}
	for k, v := range expectOne {
		if got, ok := tag.GetFirst(k); !ok || got != v {
			t.Errorf("%q: expected %q, got %q", k, v, got)
		}
	}

	unexpectedKeys := []string{"yaml", "xml"}
	for _, k := range unexpectedKeys {
		if got := tag.GetAll(k); len(got) != 0 {
			t.Errorf("%q: expected empty, got %v", k, got)
		}
		if got, ok := tag.GetFirst(k); ok {
			t.Errorf("%q: expected empty, got %v", k, got)
		}
	}
}
