//go:build !coverage

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
	"github.com/g4s8/envdoc/edit"
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
	if err := cfg.Validate(); err != nil {
		fatal("Invalid config: %v", err)
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
	})
	renderer := render.NewRenderer(cfg.OutFormat, cfg.NoStyles, cfg.Edit)
	gen := NewGenerator(parser, converter, renderer)

	// Branch based on mode
	if cfg.Edit {
		generateEdit(gen, cfg)
	} else {
		generateNormal(gen, cfg)
	}
}

// generateEdit generates documentation in edit mode, replacing content between markers
func generateEdit(gen *Generator, cfg Config) {
	var buf bytes.Buffer
	if err := gen.Generate(cfg.Dir, &buf); err != nil {
		fatal("Failed to generate: %v", err)
	}

	editor := edit.NewEditor(cfg.OutFile)
	if err := editor.ReplaceSection(buf.Bytes()); err != nil {
		fatal("Failed to edit file: %v", err)
	}

	if cfg.Debug {
		fmt.Fprintf(os.Stderr, "Successfully updated %s\n", cfg.OutFile)
	}
}

// generateNormal generates documentation in normal mode, writing to output file
func generateNormal(gen *Generator, cfg Config) {
	out, err := os.Create(cfg.OutFile)
	if err != nil {
		fatal("Failed to open output file: %v", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			fatal("Failed to close output file: %v", err)
		}
	}()

	buf := bufio.NewWriter(out)
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
