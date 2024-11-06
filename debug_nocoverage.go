//go:build !coverage

package main

import (
	"github.com/g4s8/envdoc/debug"
	"github.com/g4s8/envdoc/types"
)

func printScopesTree(s []*types.EnvScope) {
	if !debug.Config.Enabled {
		return
	}
	debug.Log("Scopes tree:\n")
	for _, scope := range s {
		debug.Logf(" - %q\n", scope.Name)
		for _, item := range scope.Vars {
			printDocItem("  ", item)
		}
	}
}

func printDocItem(prefix string, item *types.EnvDocItem) {
	debug.Logf("%s- %q\n", prefix, item.Name)
	for _, child := range item.Children {
		printDocItem(prefix+"  ", child)
	}
}
