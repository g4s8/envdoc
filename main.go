//go:build !coverage

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
	"github.com/g4s8/envdoc/render"
)

func main() {
	var cfg Config
	if err := cfg.Load(); err != nil {
		fatal("Failed to load config: %v", err)
	}
	if cfg.Debug {
		debug.Config.Enabled = true
		cfg.fprint(os.Stdout)
	}

	parser := ast.NewParser(cfg.FileGlob, cfg.TypeGlob,
		ast.WithDebug(cfg.Debug),
		ast.WithExecConfig(cfg.ExecFile, cfg.ExecLine))
	converter := NewConverter(cfg.Target, ConverterOpts{
		EnvPrefix:       cfg.EnvPrefix,
		TagName:         cfg.TagName,
		TagDefault:      cfg.TagDefault,
		RequiredIfNoDef: cfg.RequiredIfNoDef,
		UseFieldNames:   cfg.FieldNames,
		CustomTemplate:  cfg.TemplateFile,
	})
	renderer := render.NewRenderer(cfg.OutFormat, cfg.Title, cfg.NoStyles)
	gen := NewGenerator(parser, converter, renderer)

	out, err := os.Create(cfg.OutFile)
	if err != nil {
		fatal("Failed to open output file: %v", err)
	}
	buf := bufio.NewWriter(out)
	defer func() {
		if err := out.Close(); err != nil {
			fatal("Failed to close output file: %v", err)
		}
	}()

	if err := gen.Generate(cfg.Dir, buf); err != nil {
		fatal("Failed to generate: %v", err)
	}
	if err := buf.Flush(); err != nil {
		fatal("Failed to flush output: %v", err)
	}
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
