package main

import (
	"embed"
	"fmt"
	"go/ast"
	"io"
	"os"
	"path"
	"testing"
)

func TestTagParsers(t *testing.T) {
	type testCase struct {
		tag    string
		expect EnvDocItem
		fail   bool
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
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out EnvDocItem
			field := &ast.Field{
				Tag: &ast.BasicLit{Value: c.tag},
			}

			ok := parseField(field, &out)
			if ok != !c.fail {
				t.Error("parseTag returned false")
			}
			if out != c.expect {
				t.Errorf("parseTag of `%s` returned wrong result: %+v; expected: %+v", c.tag, out, c.expect)
			}
		})
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
					Name:     "Foo",
					typeName: "Foo",
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
					Name:     "Bar",
					typeName: "Bar",
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
	} {
		scopes := c.expectScopes
		if scopes == nil {
			scopes = []EnvScope{
				{
					Name:     c.typeName,
					typeName: c.typeName,
					Vars:     c.expect,
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
		insp := newInspector(typeName, all, lineN)
		scopes, err := insp.inspectFile(sourceFile)
		if err != nil {
			t.Fatal("Inspector failed", err)
		}
		if len(scopes) != len(expect) {
			t.Fatalf("inspector found %d scopes; expected %d", len(scopes), len(expect))
		}
		skipScopesCheck := len(expect) == 1 && expect[0].typeName == ""
		for i, s := range scopes {
			e := expect[i]
			if !skipScopesCheck {
				if s.Name != e.Name {
					t.Fatalf("[%d]scope: expect name %q; expected %q", i, e.Name, s.Name)
				}
				if s.typeName != e.typeName {
					t.Fatalf("[%d]scope: expect type name %q; expected %q", i, e.typeName, s.typeName)
				}
				if len(s.Vars) != len(e.Vars) {
					t.Fatalf("[%d]scope: expect %d vars; expected %d", i, len(e.Vars), len(s.Vars))
				}
			}
			for j, v := range s.Vars {
				ev := e.Vars[j]
				if v.Name != ev.Name {
					t.Fatalf("[%d]scope: var[%d]: expect name %q; expected %q", i, j, ev.Name, v.Name)
				}
				if v.Doc != ev.Doc {
					t.Fatalf("[%d]scope: var[%d]: expect doc %q; expected %q", i, j, ev.Doc, v.Doc)
				}
				if v.Opts != ev.Opts {
					t.Fatalf("[%d]scope: var[%d]: expect opts %+v; expected %+v", i, j, ev.Opts, v.Opts)
				}
			}

		}
	}
}
