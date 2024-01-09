package main

import (
	"embed"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"os"
	"path"
	"testing"
)

func TestTagParsers(t *testing.T) {
	type testCase struct {
		tag           string
		names         []string
		useFieldNames bool
		expect        EnvDocItem
		expectList    []EnvDocItem
		fail          bool
	}
	for i, c := range []testCase{
		{tag: "", fail: true},
		{tag: " ", fail: true},
		{tag: `env:"FOO"`, expect: EnvDocItem{Name: "FOO"}},
		{tag: ` env:FOO `, fail: true},
		{tag: `json:"bar"   env:"FOO"   qwe:"baz"`, expect: EnvDocItem{Name: "FOO"}},
		{tag: `env:"SECRET,file"`, expect: EnvDocItem{Name: "SECRET", Opts: EnvVarOptions{FromFile: true}}},
		{
			tag:    `env:"PASSWORD,file"           envDefault:"/tmp/password"   json:"password"`,
			expect: EnvDocItem{Name: "PASSWORD", Opts: EnvVarOptions{FromFile: true, Default: "/tmp/password"}},
		},
		{
			tag: `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`,
			expect: EnvDocItem{
				Name: "CERTIFICATE", Opts: EnvVarOptions{
					FromFile: true, Expand: true, Default: "${CERTIFICATE_FILE}",
				},
			},
		},
		{
			tag:    `env:"SECRET_KEY,required" json:"secret_key"`,
			expect: EnvDocItem{Name: "SECRET_KEY", Opts: EnvVarOptions{Required: true}},
		},
		{
			tag:    `json:"secret_val" env:"SECRET_VAL,notEmpty"`,
			expect: EnvDocItem{Name: "SECRET_VAL", Opts: EnvVarOptions{Required: true, NonEmpty: true}},
		},
		{
			tag: `fooo:"1" env:"JUST_A_MESS,required,notEmpty,file,expand" json:"just_a_mess" envDefault:"${JUST_A_MESS_FILE}" bar:"2"`,
			expect: EnvDocItem{
				Name: "JUST_A_MESS",
				Opts: EnvVarOptions{
					Required: true, NonEmpty: true, FromFile: true, Expand: true,
					Default: "${JUST_A_MESS_FILE}",
				},
			},
		},
		{
			tag: `env:"WORDS" envSeparator:";"`,
			expect: EnvDocItem{
				Name: "WORDS",
				Opts: EnvVarOptions{Separator: ";"},
			},
		},
		{
			names: []string{"Foo", "BarBaz"},
			expectList: []EnvDocItem{
				{Name: "FOO"},
				{Name: "BAR_BAZ"},
			},
			useFieldNames: true,
		},
		{
			names: []string{"Foo"},
			tag:   `env:",required"`,
			expectList: []EnvDocItem{
				{Name: "FOO", Opts: EnvVarOptions{Required: true}},
			},
			useFieldNames: true,
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			fieldNames := make([]*ast.Ident, len(c.names))
			for i, name := range c.names {
				fieldNames[i] = &ast.Ident{Name: name}
			}
			var tag *ast.BasicLit
			if c.tag != "" {
				tag = &ast.BasicLit{Value: c.tag}
			}
			field := &ast.Field{
				Tag:   tag,
				Names: fieldNames,
			}

			i := inspector{
				useFieldNames: c.useFieldNames,
			}

			expect := c.expectList
			if len(expect) == 0 && c.expect.Name != "" {
				expect = []EnvDocItem{c.expect}
			}

			actual := i.parseField(field)
			if c.fail {
				if actual != nil {
					t.Errorf("expected nil, got %#v", actual)
				}
				return
			}
			if len(expect) != len(actual) {
				t.Errorf("expected %d items, got %d", len(expect), len(actual))
			}
			for i, e := range expect {
				a := actual[i]
				if e.Name != a.name {
					t.Errorf("expected[%d] name %q, got %q", i, e.Name, a.name)
				}
				if e.Doc != a.doc {
					t.Errorf("expected[%d] doc %q, got %q", i, e.Doc, a.doc)
				}
				if e.Opts != a.opts {
					t.Errorf("expected[%d] opts %#v, got %#v", i, e.Opts, a.opts)
				}
			}
		})
	}
}

func TestInspectorError(t *testing.T) {
	sourceFile := path.Join(t.TempDir(), "tmp.go")
	if err := copyTestFile(path.Join("testdata", "type.go"), sourceFile); err != nil {
		t.Fatal("Copy test file data", err)
	}
	insp := newInspector("", true, 0, false)
	targetErr := errors.New("target error")
	insp.err = targetErr
	_, err := insp.inspectFile(sourceFile)
	if err != targetErr {
		t.Errorf("Expected error %q, got %q", targetErr, err)
	}
}

//go:embed testdata
var testdata embed.FS

func TestInspector(t *testing.T) {
	// envName    string // environment variable name
	// doc        string // field documentation text
	// flags      EnvDocItemFlags
	// envDefault string
	for _, c := range []struct {
		name         string
		typeName     string
		goLine       int
		all          bool
		expect       []EnvDocItem
		expectScopes []EnvScope
	}{
		{
			name:   "go_generate.go",
			goLine: 3,
			expect: []EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo stub",
				},
			},
		},
		{
			name:     "tags.go",
			typeName: "Type1",
			expect: []EnvDocItem{
				{
					Name: "SECRET",
					Doc:  "Secret is a secret value that is read from a file.",
					Opts: EnvVarOptions{FromFile: true},
				},
				{
					Name: "PASSWORD",
					Doc:  "Password is a password that is read from a file.",
					Opts: EnvVarOptions{FromFile: true, Default: "/tmp/password"},
				},
				{
					Name: "CERTIFICATE",
					Doc:  "Certificate is a certificate that is read from a file.",
					Opts: EnvVarOptions{
						FromFile: true, Expand: true,
						Default: "${CERTIFICATE_FILE}",
					},
				},
				{
					Name: "SECRET_KEY",
					Doc:  "Key is a secret key.",
					Opts: EnvVarOptions{Required: true},
				},
				{
					Name: "SECRET_VAL",
					Doc:  "SecretVal is a secret value.",
					Opts: EnvVarOptions{Required: true, NonEmpty: true},
				},
			},
		},
		{
			name:     "type.go",
			typeName: "Type1",
			expect: []EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo stub",
				},
			},
		},
		{
			name:     "arrays.go",
			typeName: "Arrays",
			expect: []EnvDocItem{
				{
					Name: "DOT_SEPARATED",
					Doc:  "DotSeparated stub",
					Opts: EnvVarOptions{Separator: "."},
				},
				{
					Name: "COMMA_SEPARATED",
					Doc:  "CommaSeparated stub",
					Opts: EnvVarOptions{Separator: ","},
				},
			},
		},
		{
			name:     "comments.go",
			typeName: "Comments",
			expect: []EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo stub",
				},
				{
					Name: "BAR",
					Doc:  "Bar stub",
				},
			},
		},
		{
			name: "all.go",
			all:  true,
			expectScopes: []EnvScope{
				{
					Name: "Foo",
					Vars: []EnvDocItem{
						{
							Name: "ONE",
							Doc:  "One is a one.",
						},
						{
							Name: "TWO",
							Doc:  "Two is a two.",
						},
					},
				},
				{
					Name: "Bar",
					Vars: []EnvDocItem{
						{
							Name: "THREE",
							Doc:  "Three is a three.",
						},
						{
							Name: "FOUR",
							Doc:  "Four is a four.",
						},
					},
				},
			},
		},
		{
			name:     "envprefix.go",
			typeName: "Settings",
			expect: []EnvDocItem{
				{
					Name: "DB_PORT",
					Doc:  "Port is the port to connect to",
					Opts: EnvVarOptions{Required: true},
				},
				{
					Name: "DEBUG",
					Doc:  "Debug is the debug flag",
				},
			},
		},
	} {
		scopes := c.expectScopes
		if scopes == nil {
			scopes = []EnvScope{
				{
					Name: c.typeName,
					Vars: c.expect,
				},
			}
		}
		t.Run(c.name, inspectorTester(c.name, c.typeName, c.all, c.goLine, scopes))
	}
}

func copyTestFile(name string, dest string) error {
	srcf, err := testdata.Open(name)
	if err != nil {
		return fmt.Errorf("open testdata file: %w", err)
	}
	defer srcf.Close()

	dstf, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}
	return nil
}

func inspectorTester(name string, typeName string, all bool, lineN int, expect []EnvScope) func(*testing.T) {
	return func(t *testing.T) {
		sourceFile := path.Join(t.TempDir(), "tmp.go")
		if err := copyTestFile(path.Join("testdata", name), sourceFile); err != nil {
			t.Fatal("Copy test file data", err)
		}
		insp := newInspector(typeName, all, lineN, false)
		scopes, err := insp.inspectFile(sourceFile)
		if err != nil {
			t.Fatal("Inspector failed", err)
		}
		if len(scopes) != len(expect) {
			t.Fatalf("inspector found %d scopes; expected %d", len(scopes), len(expect))
		}
		skipScopesCheck := len(expect) == 1 && expect[0].Name == ""
		for i, s := range scopes {
			e := expect[i]
			if !skipScopesCheck {
				if s.Name != e.Name {
					t.Fatalf("[%d]scope: expect name %q; was %q", i, e.Name, s.Name)
				}
				if len(s.Vars) != len(e.Vars) {
					t.Fatalf("[%d]scope: expect %d vars; was %d", i, len(e.Vars), len(s.Vars))
				}
			}
			for j, v := range s.Vars {
				ev := e.Vars[j]
				if v.Name != ev.Name {
					t.Fatalf("[%d]scope: var[%d]: expect name %q; was %q", i, j, ev.Name, v.Name)
				}
				if v.Doc != ev.Doc {
					t.Fatalf("[%d]scope: var[%d]: expect doc %q; was %q", i, j, ev.Doc, v.Doc)
				}
				if v.Opts != ev.Opts {
					t.Fatalf("[%d]scope: var[%d]: expect opts %+v; was %+v", i, j, ev.Opts, v.Opts)
				}
			}

		}
	}
}
