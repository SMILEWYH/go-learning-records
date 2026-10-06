---
title: 注释
---

<!--
This is just a placeholder page for enabling a test.
In the deployed site it is overwritten with the content of go.googlesource.com/wiki.
-->

每个包都应有包注释。它应紧挨着该包某个文件中的 ` package ` 语句之前出现（只需出现在一个文件中）。包注释的第一句应以“Package _packagename_”开头，简洁概括包的功能。godoc 会在包列表中使用这句简介。

后续句子或段落可以提供更多细节。句子应正确使用标点。

```go
// Package superman implements methods for saving the world.
//
// Experience has shown that a small number of procedures can prove
// helpful when attempting to save the world.
package superman
```

几乎每个顶层 type、const、var 和 func 声明都应有注释。bar 的注释应采用“_bar_ ……。”这种以标识符开头的形式。除非代码中的 _bar_ 以大写字母开头，否则注释中也不应将其首字母改成大写。

```go
// enterOrbit causes Superman to fly into low Earth orbit, a position
// that presents several possibilities for planet salvation.
func enterOrbit() os.Error {
  ...
}
```

godoc 会将注释中所有缩进的文本渲染为预格式化文本块，方便展示代码示例。

```go
// fight can be used on any enemy and returns whether Superman won.
//
// Examples:
//
//  fight("a random potato")
//  fight(LexLuthor{})
//
func fight(enemy interface{}) bool {
	// This is testing proper escaping in the wiki.
	for i := 0; i < 10; i++ {
		println("fight!")
	}
}
```


