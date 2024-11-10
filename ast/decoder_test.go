package ast

import (
	"testing"

	"github.com/g4s8/envdoc/testutils"
)

func TestFieldSpecDecoder(t *testing.T) {
	t.Run("single name", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Fooooooo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO,required"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:    []string{"FOO"},
			Required: true,
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("multiple names", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo", "Bar"},
			Doc:     "foo doc",
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", false, true)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names: []string{"FOO", "BAR"},
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("default value", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO" default:"bar"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "default", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:   []string{"FOO"},
			Default: "bar",
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("prefix", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("PREFIX_", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names: []string{"PREFIX_FOO"},
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("separator", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeArray},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:     []string{"FOO"},
			Separator: ",",
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("separator from tag", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO" envSeparator:";"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeArray},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:     []string{"FOO"},
			Separator: ";",
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("required if no default", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", true, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:    []string{"FOO"},
			Required: true,
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("expand", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO,expand"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:  []string{"FOO"},
			Expand: true,
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("non-empty", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO,notEmpty"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:    []string{"FOO"},
			Required: true,
			NonEmpty: true,
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("from file", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO,file"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("", "env", "", false, false)
		res, prefix := d.Decode(fs)
		assertEqFieldInfo(t, FieldInfo{
			Names:    []string{"FOO"},
			FromFile: true,
		}, res)
		testutils.AssertError(t, prefix == "", "unexpected prefix: %s", prefix)
	})
	t.Run("field prefix", func(t *testing.T) {
		fs := &FieldSpec{
			Names:   []string{"Foo"},
			Doc:     "foo doc",
			Tag:     `env:"FOO" envPrefix:"BAR_"`,
			TypeRef: FieldTypeRef{Name: "string", Kind: FieldTypeIdent},
		}
		d := NewFieldSpecDecoder("X_", "env", "", false, false)
		_, prefix := d.Decode(fs)
		testutils.AssertError(t, prefix == "X_BAR_", "unexpected prefix: %s", prefix)
	})
}

func assertEqFieldInfo(t *testing.T, expect, actual FieldInfo) {
	t.Helper()

	testutils.AssertFatal(t, len(expect.Names) == len(actual.Names), "unexpected names: %v", actual.Names)
	for i, name := range expect.Names {
		testutils.AssertError(t, name == actual.Names[i], "[%d] unexpected name: %s", i, actual.Names[i])
	}
	testutils.AssertError(t, expect.Required == actual.Required, "required flag mismatch")
	testutils.AssertError(t, expect.Expand == actual.Expand, "expand flag mismatch")
	testutils.AssertError(t, expect.NonEmpty == actual.NonEmpty, "non-empty flag mismatch")
	testutils.AssertError(t, expect.FromFile == actual.FromFile, "from-file flag mismatch")
	testutils.AssertError(t, expect.Default == actual.Default, "default value mismatch")
	testutils.AssertError(t, expect.Separator == actual.Separator, "separator mismatch")
}
