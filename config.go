package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type OutFormat string

const (
	OutFormatMarkdown OutFormat = "markdown"
	OutFormatHTML     OutFormat = "html"
	OutFormatTxt      OutFormat = "plaintext"
	OutFormatEnv      OutFormat = "dotenv"
)

type Config struct {
	// Dir to search for files
	Dir string
	// FileGlob to filter by file name
	FileGlob string
	// TypeGlob to filter by type name
	TypeGlob string
	// OutFile to write the output to
	OutFile string
	// OutFormat specify the output format
	OutFormat OutFormat
	// EnvPrefix to prefix the env vars with
	EnvPrefix string
	// NoStyles to disable styles for HTML format
	NoStyles bool
	// FieldNames flag enables field names usage intead of `env` tag.
	FieldNames bool

	// TagName sets custom tag name, `env` by default.
	TagName string
	// TagDefault sets default env tag name, `envDefault` by default.
	TagDefault string
	// TagRequiredIfNoDef sets attributes as required if no default value is set.
	RequiredIfNoDef bool

	// ExecLine is the line of go:generate command
	ExecLine int
	// ExecFile is the file of go:generate command
	ExecFile string

	// Debug output enabled
	Debug bool
}

func (c *Config) parseFlags(f *flag.FlagSet) error {
	// input flags
	f.StringVar(&c.Dir, "dir", "", "Dir to search for files, default is the file dir with go:generate command")
	f.StringVar(&c.FileGlob, "files", "", "FileGlob to filter by file name")
	f.StringVar(&c.TypeGlob, "types", "", "Type glob to filter by type name")
	// output flags
	f.StringVar(&c.OutFile, "output", "", "Output file path")
	f.StringVar((*string)(&c.OutFormat), "format", "markdown", "Output format, default `markdown`")
	f.BoolVar(&c.NoStyles, "no-styles", false, "Disable styles for HTML output")
	// app config flags
	f.StringVar(&c.EnvPrefix, "env-prefix", "", "Environment variable prefix")
	f.BoolVar(&c.FieldNames, "field-names", false, "Use field names if tag is not specified")
	f.BoolVar(&c.Debug, "debug", false, "Enable debug output")
	// customization
	f.StringVar(&c.TagName, "tag-name", "env", "Custom tag name")
	f.StringVar(&c.TagDefault, "tag-default", "envDefault", "Default tag name")
	f.BoolVar(&c.RequiredIfNoDef, "required-if-no-def", false, "Set attributes as required if no default value is set")
	// deprecated flags
	var (
		typeName string
		all      bool
	)
	f.StringVar(&typeName, "type", "", "Type name to filter by type name (deprecated: use -types instead)")
	f.BoolVar(&all, "all", false, "Generate documentation for all types in the file (deprecated: use -types='*' instead)")

	// parse
	if err := f.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	// deprecated flags `all`, `type` and new flag `types` can't be used together
	if all && typeName != "" {
		return errors.New("flags -all and -type can't be used together")
	}
	if all && c.TypeGlob != "" {
		return errors.New("flags -all and -types can't be used together")
	}
	if typeName != "" && c.TypeGlob != "" {
		return errors.New("flags -type and -types can't be used together")
	}

	// check for deprecated flags
	var deprecatedWarning strings.Builder
	if typeName != "" {
		deprecatedWarning.WriteString("\t-type flag is deprecated, use -types instead\n")
		c.TypeGlob = typeName
	}
	if all {
		deprecatedWarning.WriteString("\t-all flag is deprecated, use -types='*' instead\n")
		c.TypeGlob = "*"
	}
	if deprecatedWarning.Len() > 0 {
		fmt.Fprintln(os.Stderr, "WARNING! Deprecated flags are used. It will be removed in the next major release.")
		fmt.Fprintln(os.Stderr, deprecatedWarning.String())
	}

	return nil
}

var ErrNotCalledByGoGenerate = errors.New("not called by go generate")

func (c *Config) parseEnv() error {
	inputFileName := os.Getenv("GOFILE")
	if inputFileName == "" {
		return fmt.Errorf("no exec input file specified: %w", ErrNotCalledByGoGenerate)
	}
	c.ExecFile = inputFileName

	if e := os.Getenv("GOLINE"); e != "" {
		i, err := strconv.Atoi(e)
		if err != nil {
			return fmt.Errorf("invalid exec line number specified: %w", err)
		}
		c.ExecLine = i
	} else {
		return fmt.Errorf("no exec line number specified: %w", ErrNotCalledByGoGenerate)
	}

	if e := os.Getenv("DEBUG"); e != "" {
		c.Debug = true
	}

	return nil
}

func (c *Config) normalize() {
	c.TypeGlob = unescapeGlob(c.TypeGlob)
	c.FileGlob = unescapeGlob(c.FileGlob)
}

func (c *Config) setDefaults() {
	if c.FileGlob == "" {
		c.FileGlob = c.ExecFile
	}
	if c.Dir == "" {
		c.Dir = "."
	}
}

func (c *Config) fprint(out io.Writer) {
	fmt.Fprintln(out, "Config:")
	fmt.Fprintf(out, "  Dir: %q\n", c.Dir)
	if c.FileGlob != "" {
		fmt.Fprintf(out, "  FileGlob: %q\n", c.FileGlob)
	}
	if c.TypeGlob != "" {
		fmt.Fprintf(out, "  TypeGlob: %q\n", c.TypeGlob)
	}
	fmt.Fprintf(out, "  OutFile: %q\n", c.OutFile)
	fmt.Fprintf(out, "  OutFormat: %q\n", c.OutFormat)
	if c.EnvPrefix != "" {
		fmt.Fprintf(out, "  EnvPrefix: %q\n", c.EnvPrefix)
	}
	if c.NoStyles {
		fmt.Fprintln(out, "  NoStyles: true")
	}
	fmt.Printf("  ExecFile: %q\n", c.ExecFile)
	fmt.Printf("  ExecLine: %d\n", c.ExecLine)
	if c.FieldNames {
		fmt.Fprintln(out, "  FieldNames: true")
	}
	if c.Debug {
		fmt.Fprintln(out, "  Debug: true")
	}
}

func (c *Config) Load() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := c.parseFlags(fs); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}
	if err := c.parseEnv(); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}
	c.setDefaults()
	c.normalize()
	return nil
}
