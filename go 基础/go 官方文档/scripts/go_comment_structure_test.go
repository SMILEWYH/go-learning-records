//go:build ignore

package main

import (
	"reflect"
	"testing"
)

func TestTranslatedComments(t *testing.T) {
	a := "package p\n// English explanation.\nvar url = `https://example.test/*literal*/` // keep value\n"
	b := "package p\n// 中文说明。\n// 允许调整注释换行。\nvar url = `https://example.test/*literal*/` // 保留值\n"
	if !reflect.DeepEqual(inspectComments(a), inspectComments(b)) {
		t.Fatal("natural-language comments must not change the code structure")
	}
}

func TestProtectedCode(t *testing.T) {
	tests := []struct{ name, source, changed string }{
		{"string", `var s = "// unchanged"`, `var s = "// 已更改"`},
		{"operator", "x := 1 + 2", "x := 1 - 2"},
		{"output", "// Output:\n// hello\n", "// Output:\n// 你好\n"},
		{"output marker", "// Output:\n// hello\n", "// 输出：\n// hello\n"},
		{"unordered output", "// Unordered output:\n// a\n// b\n", "// Unordered output:\n// a\n// c\n"},
		{"directive", "//go:build linux\npackage p", "//go:build darwin\npackage p"},
		{"line directive", "//line source.go:7\npackage p", "//line source.go:8\npackage p"},
		{"implicit semicolon", "return x\n+y", "return x+y"},
		{"comment escape", "x := 1 // comment\n", "x := 1 /* 注释\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if reflect.DeepEqual(inspectComments(tc.source), inspectComments(tc.changed)) {
				t.Fatal("unexpectedly accepted a behavioral change")
			}
		})
	}
}
