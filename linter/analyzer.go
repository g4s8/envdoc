package linter

import (
	"golang.org/x/tools/go/analysis"
)

// Option is a linter configuration option.
type Option func(*linter)

// WithEnvName sets custom env tag name for linter.
func WithEnvName(name string) Option {
	return func(l *linter) {
		l.envName = name
	}
}

// WithNoComments disables check for documentation comments.
func WithNoComments() Option {
	return func(l *linter) {
		l.noComments = true
	}
}

// NewAnlyzer creates a new linter analyzer.
func NewAnlyzer(parseFlags bool, opts ...Option) *analysis.Analyzer {
	l := &linter{
		envName: "env",
	}
	for _, opt := range opts {
		opt(l)
	}
	a := &analysis.Analyzer{
		Name: "docenv",
		Doc:  "check that all environment variables are documented",
		Run:  l.run,
	}
	if parseFlags {
		a.Flags.StringVar(&l.envName, "env-name", l.envName, "environment variable tag name")
	}
	return a
}
