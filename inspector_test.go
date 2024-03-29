package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"testing"
)

func TestTagParsers(t *testing.T) {
	type testCase struct {
		tag        string
		expectName string
		expectOpts EnvVarOptions
	}
	for i, c := range []testCase{
		{tag: ""},
		{tag: " "},
		{tag: `env:"FOO"`, expectName: "FOO"},
		{tag: ` env:FOO `},
		{tag: `json:"bar"   env:"FOO"   qwe:"baz"`, expectName: "FOO"},
		{tag: `env:"SECRET,file"`, expectName: "SECRET", expectOpts: EnvVarOptions{FromFile: true}},
		{
			tag:        `env:"PASSWORD,file"           envDefault:"/tmp/password"   json:"password"`,
			expectName: "PASSWORD",
			expectOpts: EnvVarOptions{FromFile: true, Default: "/tmp/password"},
		},
		{
			tag:        `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`,
			expectName: "CERTIFICATE",
			expectOpts: EnvVarOptions{
				FromFile: true, Expand: true, Default: "${CERTIFICATE_FILE}",
			},
		},
		{
			tag:        `env:"SECRET_KEY,required" json:"secret_key"`,
			expectName: "SECRET_KEY",
			expectOpts: EnvVarOptions{Required: true},
		},
		{
			tag:        `json:"secret_val" env:"SECRET_VAL,notEmpty"`,
			expectName: "SECRET_VAL",
			expectOpts: EnvVarOptions{Required: true, NonEmpty: true},
		},
		{
			tag:        `fooo:"1" env:"JUST_A_MESS,required,notEmpty,file,expand" json:"just_a_mess" envDefault:"${JUST_A_MESS_FILE}" bar:"2"`,
			expectName: "JUST_A_MESS",
			expectOpts: EnvVarOptions{
				Required: true, NonEmpty: true, FromFile: true, Expand: true,
				Default: "${JUST_A_MESS_FILE}",
			},
		},
		{
			tag:        `env:"WORDS" envSeparator:";"`,
			expectName: "WORDS",
			expectOpts: EnvVarOptions{Separator: ";"},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			name, opts := parseEnvTag(c.tag)
			if e, a := c.expectName, name; e != a {
				t.Errorf("expected[%d] name %q, got %q", i, e, a)
			}
			if e, a := c.expectOpts, opts; e != a {
				t.Errorf("expected[%d] opts %#v, got %#v", i, e, a)
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
		name          string
		typeName      string
		goLine        int
		all           bool
		expect        []*EnvDocItem
		expectScopes  []*EnvScope
		useFieldNames bool
	}{
		{
			name:   "go_generate.go",
			goLine: 3,
			expect: []*EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo stub",
				},
			},
		},
		{
			name:     "tags.go",
			typeName: "Type1",
			expect: []*EnvDocItem{
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
			expect: []*EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo stub",
				},
			},
		},
		{
			name:     "arrays.go",
			typeName: "Arrays",
			expect: []*EnvDocItem{
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
			expect: []*EnvDocItem{
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
			expectScopes: []*EnvScope{
				{
					Name: "Foo",
					Vars: []*EnvDocItem{
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
					Vars: []*EnvDocItem{
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
			expect: []*EnvDocItem{
				{
					Doc:       "Database is the database settings.",
					debugName: "Database",
					Children: []*EnvDocItem{
						{
							Name: "DB_PORT",
							Doc:  "Port is the port to connect to",
							Opts: EnvVarOptions{Required: true},
						},
						{
							Name: "DB_HOST",
							Doc:  "Host is the host to connect to",
							Opts: EnvVarOptions{Required: true, NonEmpty: true, Default: "localhost"},
						},
						{
							Name: "DB_USER",
							Doc:  "User is the user to connect as",
						},
						{
							Name: "DB_PASSWORD",
							Doc:  "Password is the password to use",
						},
						{
							Name: "DB_DISABLE_TLS",
							Doc:  "DisableTLS is the flag to disable TLS",
						},
					},
				},
				{
					Doc:       "ServerConfig is the server settings.",
					debugName: "Server",
					Children: []*EnvDocItem{
						{
							Name: "SERVER_PORT",
							Doc:  "Port is the port to listen on",
							Opts: EnvVarOptions{Required: true},
						},
						{
							Name: "SERVER_HOST",
							Doc:  "Host is the host to listen on",
							Opts: EnvVarOptions{Required: true, NonEmpty: true, Default: "localhost"},
						},
						{
							Doc:       "TimeoutConfig is the timeout settings.",
							debugName: "Timeout",
							Children: []*EnvDocItem{
								{
									Name: "SERVER_TIMEOUT_READ",
									Doc:  "Read is the read timeout",
									Opts: EnvVarOptions{Default: "30"},
								},
								{
									Name: "SERVER_TIMEOUT_WRITE",
									Doc:  "Write is the write timeout",
									Opts: EnvVarOptions{Default: "30"},
								},
							},
						},
					},
				},
				{
					Name: "DEBUG",
					Doc:  "Debug is the debug flag",
				},
			},
		},
		{
			name:     "anonymous.go",
			typeName: "Config",
			expect: []*EnvDocItem{
				{
					Doc: "Repo is the configuration for the repository.",
					Children: []*EnvDocItem{
						{
							Name: "REPO_CONN",
							Doc:  "Conn is the connection string for the repository.",
							Opts: EnvVarOptions{Required: true, NonEmpty: true},
						},
					},
				},
			},
		},
		{
			name:     "nodocs.go",
			typeName: "Config",
			expect: []*EnvDocItem{
				{
					Children: []*EnvDocItem{
						{
							Name: "REPO_CONN",
							Opts: EnvVarOptions{Required: true, NonEmpty: true},
						},
					},
				},
			},
		},
		{
			name:          "field_names.go",
			typeName:      "FieldNames",
			useFieldNames: true,
			expect: []*EnvDocItem{
				{
					Name: "FOO",
					Doc:  "Foo is a single field.",
				},
				{
					Name: "BAR",
					Doc:  "Bar and Baz are two fields.",
				},
				{
					Name: "BAZ",
					Doc:  "Bar and Baz are two fields.",
				},
				{
					Name: "QUUX",
					Doc:  "Quux is a field with a tag.",
				},
				{
					Name: "FOO_BAR",
					Doc:  "FooBar is a field with a default value.",
					Opts: EnvVarOptions{Default: "quuux"},
				},
				{
					Name: "REQUIRED",
					Doc:  "Required is a required field.",
					Opts: EnvVarOptions{Required: true},
				},
			},
		},
		{
			name:     "embedded.go",
			typeName: "Config",
			expect: []*EnvDocItem{
				{
					Name: "START",
					Doc:  "Start date.",
					Opts: EnvVarOptions{Required: true, NonEmpty: true},
				},
			},
		},
		{
			name:     "typedef.go",
			typeName: "Config",
			expect: []*EnvDocItem{
				{
					Name: "START",
					Doc:  "Start date.",
				},
			},
		},
	} {
		scopes := c.expectScopes
		if scopes == nil {
			scopes = []*EnvScope{
				{
					Name: c.typeName,
					Vars: c.expect,
				},
			}
		}
		t.Run(c.name, inspectorTester(c.name, c.typeName, c.all, c.goLine, c.useFieldNames,
			scopes))
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

func inspectorTester(name string, typeName string, all bool, lineN int, useFieldNames bool,
	expect []*EnvScope,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Logf("inspect name=%q typeName=%q all=%v lineN=%d", name, typeName, all, lineN)
		sourceFile := path.Join(t.TempDir(), "tmp.go")
		if err := copyTestFile(path.Join("testdata", name), sourceFile); err != nil {
			t.Fatal("Copy test file data", err)
		}
		insp := newInspector(typeName, all, lineN, useFieldNames)
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
			for j, actual := range s.Vars {
				expect := e.Vars[j]
				testScopeVar(t, fmt.Sprintf("[%d]scope: var[%d]", i, j), expect, actual)
			}
		}
	}
}

func testScopeVar(t *testing.T, logPrefix string, expect, actual *EnvDocItem) {
	t.Helper()

	if expect.Name != actual.Name {
		t.Fatalf("%s: expect name %q; was %q", logPrefix, expect.Name, actual.Name)
	}
	if expect.Doc != actual.Doc {
		t.Fatalf("%s: expect doc %q; was %q", logPrefix, expect.Doc, actual.Doc)
	}
	if expect.Opts != actual.Opts {
		t.Fatalf("%s: expect opts %+v; was %+v", logPrefix, expect.Opts, actual.Opts)
	}
	if len(expect.Children) != len(actual.Children) {
		t.Fatalf("%s: expect %d children; was %d", logPrefix, len(expect.Children), len(actual.Children))
	}
	for i, c := range expect.Children {
		testScopeVar(t, fmt.Sprintf("%s -> child[%d]", logPrefix, i), c, actual.Children[i])
	}
}
