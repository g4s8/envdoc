package main

import (
	"fmt"
	"testing"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/testutils"
	"github.com/g4s8/envdoc/types"
)

func TestFieldDecoder(t *testing.T) {
	type testCase struct {
		target       types.TargetType
		name         string
		opts         FieldDecoderOpts
		spec         *ast.FieldSpec
		expectField  FieldInfo
		expectPrefix string
	}
	for _, test := range []testCase{
		{
			target: types.TargetTypeCaarlos0,
			name:   "single name",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Fooooooo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO,required"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:    []string{"FOO"},
				Required: true,
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "multiple names",
			opts: FieldDecoderOpts{
				TagName:       "env",
				UseFieldNames: true,
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo", "Bar"},
				Doc:     "foo doc",
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"FOO", "BAR"},
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "default value",
			opts: FieldDecoderOpts{
				TagName:    "env",
				TagDefault: "default",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" default:"bar"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:   []string{"FOO"},
				Default: "bar",
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "prefix",
			opts: FieldDecoderOpts{
				TagName:   "env",
				EnvPrefix: "PREFIX_",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"PREFIX_FOO"},
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "separator",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeArray},
			},
			expectField: FieldInfo{
				Names:     []string{"FOO"},
				Separator: ",",
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "separator from tag",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" envSeparator:":"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeArray},
			},
			expectField: FieldInfo{
				Names:     []string{"FOO"},
				Separator: ":",
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "required if no default",
			opts: FieldDecoderOpts{
				TagName:         "env",
				RequiredIfNoDef: true,
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:    []string{"FOO"},
				Required: true,
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "expand",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO,expand"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:  []string{"FOO"},
				Expand: true,
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "non-empty",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO,notEmpty"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:    []string{"FOO"},
				Required: true,
				NonEmpty: true,
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "from file",
			opts: FieldDecoderOpts{
				TagName: "env",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO,file"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:    []string{"FOO"},
				FromFile: true,
			},
		},
		{
			target: types.TargetTypeCaarlos0,
			name:   "field prefix",
			opts: FieldDecoderOpts{
				TagName:    "env",
				TagDefault: "default",
				EnvPrefix:  "X_",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" envPrefix:"BAR_"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"X_FOO"},
			},
			expectPrefix: "X_BAR_",
		},

		{
			target: types.TargetTypeCleanenv,
			name:   "name",
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"FOO"},
			},
		},
		{
			target: types.TargetTypeCleanenv,
			name:   "required",
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" env-required:"true"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:    []string{"FOO"},
				Required: true,
			},
		},
		{
			target: types.TargetTypeCleanenv,
			name:   "default",
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" env-default:"bar,baz"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:   []string{"FOO"},
				Default: "bar,baz",
			},
		},
		{
			target: types.TargetTypeCleanenv,
			name:   "separator",
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" env-separator:":"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names:     []string{"FOO"},
				Separator: ":",
			},
		},
		{
			target: types.TargetTypeCleanenv,
			name:   "prefix",
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" env-prefix:"BAR_"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"FOO"},
			},
			expectPrefix: "BAR_",
		},
		{
			target: types.TargetTypeCleanenv,
			name:   "field prefix",
			opts: FieldDecoderOpts{
				EnvPrefix: "X_",
			},
			spec: &ast.FieldSpec{
				Names:   []string{"Foo"},
				Doc:     "foo doc",
				Tag:     `env:"FOO" env-prefix:"BAR_"`,
				TypeRef: ast.FieldTypeRef{Name: "string", Kind: ast.FieldTypeIdent},
			},
			expectField: FieldInfo{
				Names: []string{"X_FOO"},
			},
			expectPrefix: "X_BAR_",
		},
	} {
		t.Run(fmt.Sprintf("%s_%s", test.target, test.name), func(t *testing.T) {
			d := NewFieldDecoder(test.target, test.opts)
			res, prefix := d.Decode(test.spec)
			assertEqFieldInfo(t, test.expectField, res)
			testutils.AssertError(t, prefix == test.expectPrefix, "expected prefix: %s, got: %s", test.expectPrefix, prefix)
		})
	}
}

func assertEqFieldInfo(t *testing.T, expect, actual FieldInfo) {
	t.Helper()

	testutils.AssertFatal(t, len(expect.Names) == len(actual.Names), "unexpected names: %v", actual.Names)
	for i, name := range expect.Names {
		testutils.AssertError(t, name == actual.Names[i], "[%d] expected name %q got %q", i, name, actual.Names[i])
	}
	testutils.AssertError(t, expect.Required == actual.Required, "required flag mismatch")
	testutils.AssertError(t, expect.Expand == actual.Expand, "expand flag mismatch")
	testutils.AssertError(t, expect.NonEmpty == actual.NonEmpty, "non-empty flag mismatch")
	testutils.AssertError(t, expect.FromFile == actual.FromFile, "from-file flag mismatch")
	testutils.AssertError(t, expect.Default == actual.Default, "default value mismatch")
	testutils.AssertError(t, expect.Separator == actual.Separator, "separator mismatch")
}
