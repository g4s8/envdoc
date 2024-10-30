package main

import (
	"flag"
	"os"
	"testing"
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
		if c.TypeGlob != "foo,bar" {
			t.Errorf("unexpected TypeGlob: %q", c.TypeGlob)
		}
		if c.FileGlob != "*" {
			t.Errorf("unexpected FileGlob: %q", c.FileGlob)
		}
		if c.OutFile != "out.txt" {
			t.Errorf("unexpected OutFile: %q", c.OutFile)
		}
		if c.OutFormat != "plaintext" {
			t.Errorf("unexpected OutFormat: %q", c.OutFormat)
		}
		if c.EnvPrefix != "FOO" {
			t.Errorf("unexpected EnvPrefix: %q", c.EnvPrefix)
		}
		if !c.NoStyles {
			t.Error("unexpected NoStyles: false")
		}
		if !c.FieldNames {
			t.Error("unexpected FieldNames: false")
		}
		if !c.Debug {
			t.Error("unexpected Debug: false")
		}
		if c.TagName != "xenv" {
			t.Errorf("unexpected TagName: %q", c.TagName)
		}
		if c.TagDefault != "default" {
			t.Errorf("unexpected TagDefault: %q", c.TagDefault)
		}
		if !c.RequiredIfNoDef {
			t.Error("unexpected RequiredIfNoDef: false")
		}
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
