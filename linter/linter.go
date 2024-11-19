package linter

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type linter struct {
	envName    string
	noComments bool
}

func (l *linter) run(pass *analysis.Pass) (interface{}, error) {
	tagStr := fmt.Sprintf("%s:", l.envName)
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			field, ok := n.(*ast.Field)
			if !ok {
				return true
			}

			if field.Tag == nil {
				return true
			}

			tag := field.Tag.Value
			if !strings.Contains(tag, tagStr) {
				return true
			}

			if !checkFieldDoc(field, l.noComments) {
				names := fieldNames(field)
				pass.Reportf(field.Pos(),
					"field `%s` with `%s` tag should have a documentation comment",
					names, l.envName)
			}

			return true
		})
	}
	return nil, nil
}

func fieldNames(f *ast.Field) string {
	var names []string
	for _, name := range f.Names {
		names = append(names, name.Name)
	}
	return strings.Join(names, ", ")
}

func checkFieldDoc(f *ast.Field, noComments bool) bool {
	if f.Doc != nil && strings.TrimSpace(f.Doc.Text()) != "" {
		return true
	}

	if noComments {
		return false
	} else if f.Comment != nil && strings.TrimSpace(f.Comment.Text()) != "" {
		return true
	}

	return false
}
