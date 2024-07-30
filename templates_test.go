package main

import (
	"slices"
	"testing"
)

func TestTplFuncs(t *testing.T) {
	// 	"repeat": strings.Repeat,
	// 	"split":  strings.Split,
	// 	"strAppend": func(arr []string, item string) []string {
	// 		return append(arr, item)
	// 	},
	// 	"join": strings.Join,
	// 	"strSlice": func() []string {
	// 		return make([]string, 0)
	// 	},
	// 	"list": func(args ...any) []any {
	// 		return args
	// 	},
	// 	"sum": func(args ...int) int {
	// 		var sum int
	// 		for _, v := range args {
	// 			sum += v
	// 		}
	// 		return sum
	// 	},
	t.Run("repeat", func(t *testing.T) {
		f := tplFuncs["repeat"].(func(string, int) string)
		if f("a", 3) != "aaa" {
			t.Error("repeat failed")
		}
	})
	t.Run("split", func(t *testing.T) {
		f := tplFuncs["split"].(func(string, string) []string)
		if f("a,b,c", ",") == nil {
			t.Error("split failed")
		}
	})
	t.Run("strAppend", func(t *testing.T) {
		f := tplFuncs["strAppend"].(func([]string, string) []string)
		if !slices.Equal(f([]string{"a"}, "b"), []string{"a", "b"}) {
			t.Error("strAppend failed")
		}
	})
	t.Run("join", func(t *testing.T) {
		f := tplFuncs["join"].(func([]string, string) string)
		if f([]string{"a", "b"}, ",") != "a,b" {
			t.Error("join failed")
		}
	})
	t.Run("strSlice", func(t *testing.T) {
		f := tplFuncs["strSlice"].(func() []string)
		if f() == nil {
			t.Error("strSlice failed")
		}
	})
	t.Run("list", func(t *testing.T) {
		f := tplFuncs["list"].(func(...any) []any)
		lst := f(1, 2, 3)
		for i, v := range lst {
			if v != i+1 {
				t.Error("list failed")
			}
		}
	})
	t.Run("sum", func(t *testing.T) {
		f := tplFuncs["sum"].(func(...int) int)
		if f(1, 2, 3) != 6 {
			t.Error("sum failed")
		}
	})
}
