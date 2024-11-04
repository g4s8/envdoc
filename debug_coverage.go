//go:build coverage

package main

import (
	"io"

	"github.com/g4s8/envdoc/types"
)

func (r *TypeResolver) fprint(out io.Writer) {
}

func printScopesTree(s []*types.EnvScope) {
}

func printDocItem(prefix string, item *types.EnvDocItem) {
}
