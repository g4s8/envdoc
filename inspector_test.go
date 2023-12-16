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
		tag    string
		expect docItem
		fail   bool
	}
	for i, c := range []testCase{
		{tag: "", fail: true},
		{tag: " ", fail: true},
		{tag: `env:"FOO"`, expect: docItem{envName: "FOO"}},
		{tag: ` env:FOO `, fail: true},
		{tag: `json:"bar"   env:"FOO"   qwe:"baz"`, expect: docItem{envName: "FOO"}},
		{tag: `env:"SECRET,file"`, expect: docItem{envName: "SECRET", flags: docItemFlagFromFile}},
		{
			tag:    `env:"PASSWORD,file"           envDefault:"/tmp/password"   json:"password"`,
			expect: docItem{envName: "PASSWORD", flags: docItemFlagFromFile, envDefault: "/tmp/password"},
		},
		{
			tag:    `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`,
			expect: docItem{envName: "CERTIFICATE", flags: docItemFlagFromFile | docItemFlagExpand, envDefault: "${CERTIFICATE_FILE}"},
		},
		{
			tag:    `env:"SECRET_KEY,required" json:"secret_key"`,
			expect: docItem{envName: "SECRET_KEY", flags: docItemFlagRequired},
		},
		{
			tag:    `json:"secret_val" env:"SECRET_VAL,notEmpty"`,
			expect: docItem{envName: "SECRET_VAL", flags: docItemFlagNonEmpty | docItemFlagRequired},
		},
		{
			tag: `fooo:"1" env:"JUST_A_MESS,required,notEmpty,file,expand" json:"just_a_mess" envDefault:"${JUST_A_MESS_FILE}" bar:"2"`,
			expect: docItem{
				envName:    "JUST_A_MESS",
				flags:      docItemFlagRequired | docItemFlagNonEmpty | docItemFlagFromFile | docItemFlagExpand,
				envDefault: "${JUST_A_MESS_FILE}",
			},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var out docItem
			ok := parseTag(c.tag, &out)
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
	// flags      docItemFlags
	// envDefault string
	for _, c := range []struct {
		name     string
		typeName string
		goLine   int
		expect   []docItem
	}{
		{
			name:   "example_generate.go",
			goLine: 3,
			expect: []docItem{
				{
					envName: "FOO",
					doc:     "Foo stub",
				},
			},
		},
		{
			name:     "example_tags.go",
			typeName: "Type1",
			expect: []docItem{
				{
					envName: "SECRET",
					doc:     "Secret is a secret value that is read from a file.",
					flags:   docItemFlagFromFile,
				},
				{
					envName:    "PASSWORD",
					doc:        "Password is a password that is read from a file.",
					flags:      docItemFlagFromFile,
					envDefault: "/tmp/password",
				},
				{
					envName:    "CERTIFICATE",
					doc:        "Certificate is a certificate that is read from a file.",
					flags:      docItemFlagFromFile | docItemFlagExpand,
					envDefault: "${CERTIFICATE_FILE}",
				},
				{
					envName: "SECRET_KEY",
					doc:     "Key is a secret key.",
					flags:   docItemFlagRequired,
				},
				{
					envName: "SECRET_VAL",
					doc:     "SecretVal is a secret value.",
					flags:   docItemFlagNonEmpty | docItemFlagRequired,
				},
			},
		},
		{
			name:     "example_type.go",
			typeName: "Type1",
			expect: []docItem{
				{
					envName: "FOO",
					doc:     "Foo stub",
				},
			},
		},
	} {
		t.Run(c.name, inspectorTester(c.name, c.typeName, c.goLine, c.expect))
	}
}

type testOutput struct {
	items []docItem
}

func (o *testOutput) writeItem(i docItem) {
	o.items = append(o.items, i)
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

func inspectorTester(name string, typeName string, lineN int, expect []docItem) func(*testing.T) {
	return func(t *testing.T) {
		var out testOutput
		sourceFile := path.Join(t.TempDir(), "tmp.go")
		if err := copyTestFile(path.Join("testdata", name), sourceFile); err != nil {
			t.Fatal("Copy test file data", err)
		}
		insp := newInspector(typeName, &out, lineN)
		if err := insp.inspectFile(sourceFile); err != nil {
			t.Fatal("Inspector failed", err)
		}
		if len(out.items) != len(expect) {
			t.Errorf("inspector found %d items; expected %d", len(out.items), len(expect))
		}
		for i := range out.items {
			if out.items[i] != expect[i] {
				t.Errorf("inspector found item[%d] %+v; expected %+v", i,
					out.items[i], expect[i])
			}
		}
	}
}
