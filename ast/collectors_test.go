package ast

import (
	"testing"

	"github.com/g4s8/envdoc/testutils"
)

func testGlob(t *testing.T, name string) func(string) bool {
	t.Helper()
	return func(s string) bool {
		ok := s == name
		t.Logf("test glob: %s == %s = %t", s, name, ok)
		return ok
	}
}

func TestRootCollector(t *testing.T) {
	t.Run("glob", func(t *testing.T) {
		col := NewRootCollector("./base",
			WithFileGlob(testGlob(t, "./first.go")),
			WithTypeGlob(testGlob(t, "TestType")))
		col.onFile(&FileSpec{
			Name: "./base/first.go",
		})
		col.onFile(&FileSpec{
			Name: "./base/second.go",
		})
		files := col.Files()
		testutils.AssertFatal(t, len(files) == 2, "unexpected files count: %d", len(files))
		testutils.AssertError(t, files[0].Name == "./first.go", "[0]unexpected file name: %s", files[0].Name)
		testutils.AssertError(t, files[1].Name == "./second.go", "[1]unexpected file name: %s", files[1].Name)
		testutils.AssertError(t, files[0].Export == true, "[0]unexpected export flag: %t", files[0].Export)
		testutils.AssertError(t, files[1].Export == false, "[1]unexpected export flag: %t", files[1].Export)
		current := col.currentFile()
		testutils.AssertError(t, current.Name == "./second.go", "unexpected current file: %s", current.Name)

		col.onType(&TypeSpec{
			Name: "TestType",
		})
		col.onType(&TypeSpec{
			Name: "AnotherType",
		})
		types := current.Types
		testutils.AssertFatal(t, len(types) == 2, "unexpected types count: %d", len(types))
		testutils.AssertError(t, types[0].Name == "TestType", "[0]unexpected type name: %s", types[0].Name)
		testutils.AssertError(t, types[1].Name == "AnotherType", "[1]unexpected type name: %s", types[1].Name)
		testutils.AssertError(t, types[0].Export == true, "[0]unexpected export flag: %t", types[0].Export)
		testutils.AssertError(t, types[1].Export == false, "[1]unexpected export flag: %t", types[1].Export)
	})
	t.Run("comment", func(t *testing.T) {
		col := NewRootCollector("./base",
			WithGoGenDecl(1, "./first.go"))
		col.onFile(&FileSpec{
			Name: "./base/first.go",
		})
		col.setComment(&CommentSpec{
			Line: 1,
		})
		current := col.currentFile()
		testutils.AssertFatal(t, current.Name == "./first.go", "unexpected current file: %s", current.Name)
		testutils.AssertError(t, col.pendingType == true, "unexpected pending type flag: %t", col.pendingType)

		col.onType(&TypeSpec{
			Name: "TestType",
		})
		col.onType(&TypeSpec{
			Name: "AnotherType",
		})
		types := current.Types
		testutils.AssertFatal(t, len(types) == 2, "unexpected types count: %d", len(types))
		testutils.AssertError(t, types[0].Name == "TestType", "[0]unexpected type name: %s", types[0].Name)
		testutils.AssertError(t, types[1].Name == "AnotherType", "[1]unexpected type name: %s", types[1].Name)
		testutils.AssertError(t, types[0].Export == true, "[0]unexpected export flag: %t", types[0].Export)
		testutils.AssertError(t, types[1].Export == false, "[1]unexpected export flag: %t", types[1].Export)
	})
}
