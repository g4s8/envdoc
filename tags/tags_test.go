package tags

import (
	"fmt"
	"slices"
	"testing"
)

func tagShouldErr(t *testing.T, tag string) {
	t.Helper()
	if len(ParseFieldTag(tag)) != 0 {
		t.Errorf("expected empty, got %v", tag)
	}
}

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

	t.Run("error", func(t *testing.T) {
		tagShouldErr(t, `envPASSWORD`)
		tagShouldErr(t, `env:"PASSWORD`)
		tagShouldErr(t, `env:PASSWORD"`)
	})
}

func TestFieldTagValues(t *testing.T) {
	tests := []struct {
		tag, key string
		expect   []string
		err      bool
	}{
		{
			tag:    `env:"PASSWORD,required,file"`,
			key:    "env",
			expect: []string{"PASSWORD", "required", "file"},
		},
		{
			tag:    `envDefault:"/tmp/password"`,
			key:    "envDefault",
			expect: []string{"/tmp/password"},
		},
		{
			tag:    `json:"password"`,
			key:    "json",
			expect: []string{"password"},
		},
		{
			tag:    `envDefault:"GET, POST, PUT, PATCH, DELETE, OPTIONS"`,
			key:    "envDefault",
			expect: []string{"GET", " POST", " PUT", " PATCH", " DELETE", " OPTIONS"},
		},
		{
			tag: `jsonpassword`,
			key: "json",
			err: true,
		},
		{
			tag: `json:"password`,
			key: "env",
			err: true,
		},
		{
			tag: `env:PASSWORD"`,
			key: "env",
			err: true,
		},
		{
			tag: `env:"PASSWORD`,
			key: "env",
			err: true,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			vals := fieldTagValues(test.tag, test.key)
			if test.err {
				if vals != nil {
					t.Errorf("expected nil, got %v", vals)
				}
				return
			}

			if !slices.Equal(vals, test.expect) {
				t.Errorf("expected %v, got %v", test.expect, vals)
			}
		})
	}
}
