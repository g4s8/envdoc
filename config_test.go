package main

import (
	"flag"
	"os"
	"testing"

	"github.com/g4s8/envdoc/testutils"
)

func TestConfig(t *testing.T) {
	t.Run("parse flags", func(t *testing.T) {
		var c Config
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{
			"test",
			"-types", "foo,bar",
			"-files", "*",
			"-output", "out.txt",
			"-format", "plaintext",
			"-env-prefix", "FOO",
			"-no-styles",
			"-field-names",
			"-debug",
			"-tag-name", "xenv",
			"-tag-default", "default",
			"-required-if-no-def",
		}
		if err := c.parseFlags(fs); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutils.AssertError(t, c.TypeGlob == "foo,bar", "unexpected TypeGlob: %q", c.TypeGlob)
		testutils.AssertError(t, c.FileGlob == "*", "unexpected FileGlob: %q", c.FileGlob)
		testutils.AssertError(t, c.OutFile == "out.txt", "unexpected OutFile: %q", c.OutFile)
		testutils.AssertError(t, c.OutFormat == "plaintext", "unexpected OutFormat: %q", c.OutFormat)
		testutils.AssertError(t, c.EnvPrefix == "FOO", "unexpected EnvPrefix: %q", c.EnvPrefix)
		testutils.AssertError(t, c.NoStyles, "unexpected NoStyles: false")
		testutils.AssertError(t, c.FieldNames, "unexpected FieldNames: false")
		testutils.AssertError(t, c.Debug, "unexpected Debug: false")
		testutils.AssertError(t, c.TagName == "xenv", "unexpected TagName: %q", c.TagName)
		testutils.AssertError(t, c.TagDefault == "default", "unexpected TagDefault: %q", c.TagDefault)
		testutils.AssertError(t, c.RequiredIfNoDef, "unexpected RequiredIfNoDef: false")
	})
	t.Run("normalize", func(t *testing.T) {
		var c Config
		c.TypeGlob = `"foo"`
		c.FileGlob = `"*"`
		c.normalize()
		if c.TypeGlob != "foo" {
			t.Errorf("unexpected TypeGlob: %q", c.TypeGlob)
		}
		if c.FileGlob != "*" {
			t.Errorf("unexpected FileGlob: %q", c.FileGlob)
		}
	})
}
