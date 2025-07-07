package ast

import "testing"

func TestImportSpec(t *testing.T) {
	s := &ImportSpec{
		Name: "alias",
		Path: "github.com/example/package",
	}
	if s.PathName() != "package" {
		t.Errorf("expected package name 'package', got %q", s.PathName())
	}
}

func TestFieldTypeRef(t *testing.T) {
	type testCase struct {
		ref FieldTypeRef
		str string
	}

	cases := []testCase{
		{FieldTypeRef{Name: "MyType", Pkg: "", Kind: FieldTypeIdent}, "MyType"},
		{FieldTypeRef{Name: "MyType", Pkg: "mypkg", Kind: FieldTypeSelector}, "mypkg.MyType"},
		{FieldTypeRef{Name: "MyType", Kind: FieldTypePtr}, "*MyType"},
		{FieldTypeRef{Name: "MyType", Kind: FieldTypeArray}, "[]MyType"},
		{FieldTypeRef{Name: "MyType", Kind: FieldTypeMap}, "map[string]MyType"},
		{FieldTypeRef{Name: "MyType", Kind: FieldTypeStruct}, "struct"},
	}
	for _, tc := range cases {
		ref := tc.ref
		if ref.String() != tc.str {
			t.Errorf("expected %q, got %q", tc.str, ref.String())
		}
	}
}
