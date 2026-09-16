package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/g4s8/envdoc/linter"
)

func main() {
	analyzer := linter.NewAnlyzer(true)
	unitchecker.Main(analyzer)
}
