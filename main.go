//go:build !coverage

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/g4s8/envdoc/debug"
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

	p := NewParser(cfg.Dir, cfg.FileGlob, cfg.TypeGlob,
		withDebug(cfg.Debug),
		withExecConfig(cfg.ExecFile, cfg.ExecLine))
	files, err := p.Parse()
	if err != nil {
		fatal("Failed to parse: %v", err)
	}

	res := ResolveAllTypes(files)
	if cfg.Debug {
		res.fprint(os.Stdout)
	}

	conv := NewConverter(cfg.EnvPrefix, cfg.FieldNames)
	scopes := conv.ScopesFromFiles(res, files)
	printScopesTree(scopes)

	r := NewRenderer(cfg.OutFormat, cfg.EnvPrefix, cfg.NoStyles)
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
	if err := r.Render(scopes, buf); err != nil {
		fatal("Failed to render: %v", err)
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
