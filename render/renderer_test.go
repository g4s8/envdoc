package render

import (
	"strings"
	"testing"

	"github.com/g4s8/envdoc/types"
)

func TestRender(t *testing.T) {
	r := NewRenderer(types.OutFormatTxt, false)
	scopes := []*types.EnvScope{
		{
			Name: "scope1",
			Doc:  "scope1 doc",
			Vars: []*types.EnvDocItem{
				{
					Name: "VAR1",
					Doc:  "VAR1 doc",
					Opts: types.EnvVarOptions{
						Required: true,
					},
				},
			},
		},
	}
	var sb strings.Builder
	if err := r.Render(scopes, &sb); err != nil {
		t.Fatalf("Failed to render: %s", err)
	}
	// Environment Variables

	// scope1

	// scope1 doc

	//  * `VAR1` (required): VAR1 doc (required)
	var expectSb strings.Builder
	expectSb.WriteString("Environment Variables\n\n")
	expectSb.WriteString("## scope1\n\n")
	expectSb.WriteString("scope1 doc\n\n")
	expectSb.WriteString(" * `VAR1` (required) - VAR1 doc\n")
	expectSb.WriteString("\n")

	if expect, actual := expectSb.String(), sb.String(); actual != expect {
		t.Logf("Expected:\n%s", expect)
		t.Logf("Got:\n%s", actual)
		t.Fatalf("Unexpected output")
	}
}

func TestRenderCustom(t *testing.T) {
	r := NewRenderer(types.OutFormatTxt, false)
	scopes := []*types.EnvScope{
		{
			Name: "scope1",
			Doc:  "scope1 doc",
			Vars: []*types.EnvDocItem{
				{
					Name: "VAR1",
					Doc:  "VAR1 doc",
					Opts: types.EnvVarOptions{
						Required: true,
					},
				},
			},
		},
	}

	tmplPath := "testdata/customtext.tmpl"

	var sb strings.Builder

	if err := r.RenderCustom(scopes, tmplPath, &sb); err != nil {
		t.Fatalf("Failed to render: %s", err)
	}

	// Environment Variables

	// scope1

	// scope1 doc

	// - `VAR1` (required): VAR1 doc (REQUIRED)

	var expectSb strings.Builder
	expectSb.WriteString("Environment Variables\n\n")
	expectSb.WriteString("## scope1\n\n")
	expectSb.WriteString("scope1 doc\n\n")
	expectSb.WriteString("- `VAR1` (REQUIRED) - VAR1 doc\n")
	expectSb.WriteString("\n")

	if expect, actual := expectSb.String(), sb.String(); actual != expect {
		t.Logf("Expected:\n%s", expect)
		t.Logf("Got:\n%s", actual)
		t.Fatalf("Unexpected output")
	}
}
