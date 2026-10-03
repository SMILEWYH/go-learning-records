---
title: "Go 文档注释"
layout: article
date: 2022-06-01T00:00:00Z
template: true
---

目录：

 [包](#package)\
 [命令](#cmd)\
 [类型](#type)\
 [函数](#func)\
 [常量](#const)\
 [变量](#var)\
 [语法](#syntax)\
 [常见错误与陷阱](#mistakes)

“文档注释”是紧挨顶层 package、const、func、type 和 var 声明之前的注释，中间没有空行。每个导出名称（以大写字母开头的名称）都应有文档注释。

[go/doc](/pkg/go/doc) 和 [go/doc/comment](/pkg/go/doc/comment) 包提供从 Go 源代码中提取文档的能力，许多工具都利用了这一功能。[`go` `doc`](/cmd/go#hdr-Show_documentation_for_package_or_symbol) 命令查找并输出指定包或符号的文档注释。这里的符号指顶层的 const、func、type 或 var。[pkg.go.dev](https://pkg.go.dev/) 网站会在许可证允许的情况下，展示公开 Go 包的文档。提供该网站服务的程序是 [golang.org/x/pkgsite/cmd/pkgsite](https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite)，它也能在本地运行，用于查看私有模块的文档，或在离线时浏览文档。语言服务器 [gopls](https://pkg.go.dev/golang.org/x/tools/gopls) 则会在 IDE 中编辑 Go 源文件时提供文档。

下文介绍如何编写 Go 文档注释。

## 包 {#package}

每个包都应有介绍该包的包注释。它提供与整个包相关的信息，并说明使用者可以期待包提供什么功能。特别是对于较大的包，简要概述最重要的 API，并按需链接到其他文档注释，通常很有帮助。

如果包很简单，包注释也可以很简短。例如：

	// Package path implements utility routines for manipulating slash-separated
	// paths.
	//
	// The path package should only be used for paths separated by forward
	// slashes, such as the paths in URLs. This package does not deal with
	// Windows paths with drive letters or backslashes; to manipulate
	// operating system paths, use the [path/filepath] package.
	package path


`[path/filepath]` 中的方括号会创建一个[文档链接](#links)。

如例所示，Go 文档注释使用完整句子。对于英文包注释，[第一句](/pkg/go/doc/#Package.Synopsis)应以“Package <name>”开头。

对于由多个文件组成的包，包注释应只放在一个源文件中。如果多个文件都有包注释，它们会被拼接起来，形成整个包的一段较长注释。

## 命令 {#cmd}

命令的包注释与普通包类似，但描述的是程序的行为，而不是包中的 Go 符号。第一句通常以程序名称开头，因为位于句首，英文名称的首字母应大写。例如，下面是 [gofmt](/cmd/gofmt) 包注释的节选：

	/*
	Gofmt formats Go programs.
	It uses tabs for indentation and blanks for alignment.
	Alignment assumes that an editor is using a fixed-width font.

	Without an explicit path, it processes the standard input. Given a file,
	it operates on that file; given a directory, it operates on all .go files in
	that directory, recursively. (Files starting with a period are ignored.)
	By default, gofmt prints the reformatted sources to standard output.

	Usage:

		gofmt [flags] [path ...]

	The flags are:

		-d
			Do not print reformatted sources to standard output.
			If a file's formatting is different than gofmt's, print diffs
			to standard output.
		-w
			Do not print reformatted sources to standard output.
			If a file's formatting is different from gofmt's, overwrite it
			with gofmt's version. If an error occurred during overwriting,
			the original file is restored from an automatic backup.

	When gofmt reads from standard input, it accepts either a full Go program
	or a program fragment. A program fragment must be a syntactically
	valid declaration list, statement list, or expression. When formatting
	such a fragment, gofmt preserves leading indentation as well as leading
	and trailing spaces, so that individual sections of a Go program can be
	formatted by piping them through gofmt.
	*/
	package main


注释开头使用了[语义换行](https://rhodesmill.org/brandon/2012/one-sentence-per-line/)：每个新句子或较长短语单独占一行。这样在代码和注释发生变化时，差异更容易阅读。后面的段落没有遵循这个约定，而是手动按行宽换行。选择最适合自己代码库的方式即可。无论哪种方式，`go` `doc` 和 `pkgsite` 输出文档注释时都会重新折行。例如：

	$ go doc gofmt
	Gofmt formats Go programs. It uses tabs for indentation and blanks for
	alignment. Alignment assumes that an editor is using a fixed-width font.

	Without an explicit path, it processes the standard input. Given a file, it
	operates on that file; given a directory, it operates on all .go files in that
	directory, recursively. (Files starting with a period are ignored.) By default,
	gofmt prints the reformatted sources to standard output.

	Usage:

		gofmt [flags] [path ...]

	The flags are:

		-d
			Do not print reformatted sources to standard output.
			If a file's formatting is different than gofmt's, print diffs
			to standard output.
	...


带缩进的行会被当作预格式化文本，不会重新折行，在 HTML 和 Markdown 中会使用代码字体展示。详细规则见下方的[语法](#syntax)一节。

## 类型 {#type}

类型的文档注释应说明该类型的每个实例表示什么，或提供什么能力。如果 API 简单，注释也可以很短。例如：

	package zip

	// A Reader serves content from a ZIP archive.
	type Reader struct {
		...
	}


默认情况下，程序员应认为一个类型的实例只能在同一时刻由一个 goroutine 安全使用。如果某个类型提供更强的保证，文档注释应明确说明。例如：

	package regexp

	// Regexp is the representation of a compiled regular expression.
	// A Regexp is safe for concurrent use by multiple goroutines,
	// except for configuration methods, such as Longest.
	type Regexp struct {
		...
	}


Go 类型还应尽量让零值具有实用意义。如果这一意义并不显然，就应该写入文档。例如：

	package bytes

	// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	// The zero value for Buffer is an empty buffer ready to use.
	type Buffer struct {
		...
	}


对于包含导出字段的结构体，类型文档注释或各字段的注释应解释每个导出字段的含义。例如，下面的类型直接在文档注释中解释字段：

{{raw `
	package io

	// A LimitedReader reads from R but limits the amount of
	// data returned to just N bytes. Each call to Read
	// updates N to reflect the new amount remaining.
	// Read returns EOF when N <= 0.
	type LimitedReader struct {
		R   Reader // underlying reader
		N   int64  // max bytes remaining
	}
`}}

相比之下，下面的类型将说明放在各字段自己的注释中：

{{raw `
	package comment

	// A Printer is a doc comment printer.
	// The fields in the struct can be filled in before calling
	// any of the printing methods
	// in order to customize the details of the printing process.
	type Printer struct {
		// HeadingLevel is the nesting level used for
		// HTML and Markdown headings.
		// If HeadingLevel is zero, it defaults to level 3,
		// meaning to use <h3> and ###.
		HeadingLevel int
		...
	}
`}}

与上文的包和下文的函数一样，类型的文档注释应以完整句子开头，点明所声明符号的名称。明确的主语通常能让表述更清楚，也便于在网页或命令行中搜索。例如：

	$ go doc -all regexp | grep pairs
	pairs within the input string: result[2*n:2*n+2] identifies the indexes
	    FindReaderSubmatchIndex returns a slice holding the index pairs identifying
	    FindStringSubmatchIndex returns a slice holding the index pairs identifying
	    FindSubmatchIndex returns a slice holding the index pairs identifying the
	$


## 函数 {#func}

函数的文档注释应说明函数返回什么；如果调用目的是产生某种作用，则应说明它做什么。注释中可以直接引用具名参数和结果，无需反引号等特殊语法。因此，通常会避免使用 `a` 这样容易与普通英文单词混淆的名称。例如：

	package strconv

	// Quote returns a double-quoted Go string literal representing s.
	// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
	// for control characters and non-printable characters as defined by IsPrint.
	func Quote(s string) string {
		...
	}


再如：

	package os

	// Exit causes the current program to exit with the given status code.
	// Conventionally, code zero indicates success, non-zero an error.
	// The program terminates immediately; deferred functions are not run.
	//
	// For portability, the status code should be in the range [0, 125].
	func Exit(code int) {
		...
	}


描述返回布尔值的函数时，英文文档注释通常使用“reports whether”（报告是否），无需再加“or not”。例如：

	package strings

	// HasPrefix reports whether the string s begins with prefix.
	func HasPrefix(s, prefix string) bool


如果文档注释需要解释多个返回结果，即使函数体内不用这些名称，为结果命名也能使注释更易理解。例如：

	package io

	// Copy copies from src to dst until either EOF is reached
	// on src or an error occurs. It returns the total number of bytes
	// written and the first error encountered while copying, if any.
	//
	// A successful Copy returns err == nil, not err == EOF.
	// Because Copy is defined to read from src until EOF, it does
	// not treat an EOF from Read as an error to be reported.
	func Copy(dst Writer, src Reader) (n int64, err error) {
		...
	}


反过来，如果文档注释不需要引用结果名称，代码中通常也省略它们，避免显得杂乱，例如上面的 `Quote`。

这些规则同时适用于普通函数和方法。对于同一类型的方法，使用一致的接收者名称，可以避免列出所有方法时出现不必要的差异：

	$ go doc bytes.Buffer
	package bytes // import "bytes"

	type Buffer struct {
		// Has unexported fields.
	}
	    A Buffer is a variable-sized buffer of bytes with Read and Write methods.
	    The zero value for Buffer is an empty buffer ready to use.

	func NewBuffer(buf []byte) *Buffer
	func NewBufferString(s string) *Buffer
	func (b *Buffer) Bytes() []byte
	func (b *Buffer) Cap() int
	func (b *Buffer) Grow(n int)
	func (b *Buffer) Len() int
	func (b *Buffer) Next(n int) []byte
	func (b *Buffer) Read(p []byte) (n int, err error)
	func (b *Buffer) ReadByte() (byte, error)
	...


这个例子还展示了另一点：返回类型 `T` 或指针 `*T` 的包级函数，即使同时返回一个错误，也会与类型 `T` 及其方法一起展示，因为文档工具会将它们视为 `T` 的构造函数。

默认情况下，程序员可以认为包级函数允许从多个 goroutine 安全调用，无需在文档中专门声明。

另一方面，正如上一节所说，对类型实例的任何使用，包括调用方法，通常都假定同一时刻仅限一个 goroutine。如果类型的文档注释没有说明哪些方法可以安全地并发使用，就应在相应方法的注释中写明。例如：

	package sql

	// Close returns the connection to the connection pool.
	// All operations after a Close will return with ErrConnDone.
	// Close is safe to call concurrently with other operations and will
	// block until all other operations finish. It may be useful to first
	// cancel any used context and then call Close directly after.
	func (c *Conn) Close() error {
		...
	}


函数和方法的文档注释应关注操作返回什么、做什么，提供调用方需要知道的细节。特殊情况尤其值得记录。例如：

{{raw `
	package math

	// Sqrt returns the square root of x.
	//
	// Special cases are:
	//
	//	Sqrt(+Inf) = +Inf
	//	Sqrt(±0) = ±0
	//	Sqrt(x < 0) = NaN
	//	Sqrt(NaN) = NaN
	func Sqrt(x float64) float64 {
		...
	}
`}}

文档注释不应解释当前实现使用何种算法之类的内部细节，这些内容更适合写在函数体内的注释中。如果渐近时间或空间复杂度对调用方特别重要，可以在文档中说明。例如：

	package sort

	// Sort sorts data in ascending order as determined by the Less method.
	// It makes one call to data.Len to determine n and O(n*log(n)) calls to
	// data.Less and data.Swap. The sort is not guaranteed to be stable.
	func Sort(data Interface) {
		...
	}


这段文档没有指定使用哪种排序算法，因此未来可以更方便地替换实现算法。

## 常量 {#const}

Go 的声明语法允许将声明分组。这时可以用一段文档注释介绍一组相关常量，而各个常量只需简短的行尾注释。例如：

	package scanner // import "text/scanner"

	// The result of Scan is one of these tokens or a Unicode character.
	const (
		EOF = -(iota + 1)
		Ident
		Int
		Float
		Char
		...
	)


有时整组常量不需要额外的文档注释。例如：

	package unicode // import "unicode"

	const (
		MaxRune         = '\U0010FFFF' // maximum valid Unicode code point.
		ReplacementChar = '\uFFFD'     // represents invalid code points.
		MaxASCII        = '\u007F'     // maximum ASCII value.
		MaxLatin1       = '\u00FF'     // maximum Latin-1 value.
	)


另一方面，未分组的常量通常应配有以完整句子开头的文档注释。例如：

	package unicode

	// Version is the Unicode edition from which the tables are derived.
	const Version = "13.0.0"


有类型常量会显示在其类型声明附近，因此通常省略常量组的文档注释，使用类型本身的注释说明。例如：

	package syntax

	// An Op is a single regular expression operator.
	type Op uint8

	const (
		OpNoMatch        Op = 1 + iota // matches no strings
		OpEmptyMatch                   // matches empty string
		OpLiteral                      // matches Runes sequence
		OpCharClass                    // matches Runes interpreted as range pair list
		OpAnyCharNotNL                 // matches any character except newline
		...
	)


HTML 展示效果见 [pkg.go.dev/regexp/syntax#Op](https://pkg.go.dev/regexp/syntax#Op)。

## 变量 {#var}

变量遵循与常量相同的约定。例如，下面是一组变量：

	package fs

	// Generic file system errors.
	// Errors returned by file systems can be tested against these errors
	// using errors.Is.
	var (
		ErrInvalid    = errInvalid()    // "invalid argument"
		ErrPermission = errPermission() // "permission denied"
		ErrExist      = errExist()      // "file already exists"
		ErrNotExist   = errNotExist()   // "file does not exist"
		ErrClosed     = errClosed()     // "file already closed"
	)


下面是单个变量：

	package unicode

	// Scripts is the set of Unicode script tables.
	var Scripts = map[string]*RangeTable{
		"Adlam":                  Adlam,
		"Ahom":                   Ahom,
		"Anatolian_Hieroglyphs":  Anatolian_Hieroglyphs,
		"Arabic":                 Arabic,
		"Armenian":               Armenian,
		...
	}


## 语法 {#syntax}

Go 文档注释采用简单语法，支持段落、标题、链接、列表和预格式化代码块。为了让注释在源文件中保持轻量和可读，不支持字体变化、原始 HTML 等复杂功能。熟悉 Markdown 的读者可以将它理解为 Markdown 的一个简化子集。

标准格式化工具 [gofmt](/cmd/gofmt) 会重新格式化文档注释，让这些元素采用统一格式。Gofmt 兼顾可读性与用户对源代码中注释写法的控制，但也会调整排版，使注释的语义更清楚，类似于将普通源代码中的 `1+2 * 3` 格式化为 `1 + 2*3`。

Gofmt 会移除文档注释开头和结尾的空行。如果注释的所有行都以相同的空格和制表符序列开头，也会移除这个共同前缀。

### 段落 {#paragraphs}

段落是一段连续的、没有缩进的非空行。前面已经有许多段落示例。

两个连续的反引号（\`，U+0060）会被解释为 Unicode 左双引号（“，U+201C）；两个连续的单引号（\'，U+0027）会被解释为 Unicode 右双引号（”，U+201D）。

Gofmt 会保留段落中的换行，不会重新折行，因此可以使用前面介绍的[语义换行](https://rhodesmill.org/brandon/2012/one-sentence-per-line/)。段落之间连续的多个空行会被合并为一个。Gofmt 还会把连续的反引号或单引号转换为对应的 Unicode 引号。

#### 备注 {#notes}

备注是一种形式为 `MARKER(uid): body` 的特殊注释。MARKER 应由至少两个大写 `[A-Z]` 字母组成，用于标识备注类型；uid 至少包含一个字符，通常是能够提供更多信息的人的用户名。uid 后面的 `:` 可以省略。

pkg.go.dev 会收集这些备注，并在独立章节中展示。

例如：

	// TODO(user1): refactor to use standard library context
	// BUG(user2): not cleaned up
	var ctx context.Context


#### 弃用说明 {#deprecations}

以 `Deprecated: ` 开头的段落会被视为弃用说明。有些工具会在使用已弃用标识符时发出提示，[pkg.go.dev](https://pkg.go.dev) 则默认隐藏它们的文档。

弃用标记之后应说明弃用的相关信息，并在适用时建议替代方案。这一段不一定要位于文档注释的最后。

例如：

	// Package rc4 implements the RC4 stream cipher.
	//
	// Deprecated: RC4 is cryptographically broken and should not be used
	// except for compatibility with legacy systems.
	//
	// This package is frozen and no new functionality will be added.
	package rc4

	// Reset zeros the key data and makes the Cipher unusable.
	//
	// Deprecated: Reset can't guarantee that the key will be entirely removed from
	// the process's memory.
	func (c *Cipher) Reset()


### 标题 {#headings}

标题行以井号（U+0023）开头，后接一个空格和标题文字。要被识别为标题，该行不能缩进，并且必须用空行与前后的段落隔开。

例如：

	// Package strconv implements conversions to and from string representations
	// of basic data types.
	//
	// # Numeric Conversions
	//
	// The most common numeric conversions are [Atoi] (string to int) and [Itoa] (int to string).
	...
	package strconv


以下则不是标题：

	// #This is not a heading, because there is no space.
	//
	// # This is not a heading,
	// # because it is multiple lines.
	//
	// # This is not a heading,
	// because it is also multiple lines.
	//
	// The next paragraph is not a heading, because there is no additional text:
	//
	// #
	//
	// In the middle of a span of non-blank lines,
	// # this is not a heading either.
	//
	//     # This is not a heading, because it is indented.


`#` 标题语法在 Go 1.19 中加入。此前，满足某些条件的单行段落会被隐式识别为标题，其中最主要的条件是不以标点结束。

Gofmt 会把早期版本[视为隐式标题的行](https://github.com/golang/proposal/blob/master/design/51082-godocfmt.md#headings)改为使用 `#` 的标题。如果这种格式化不符合本意，也就是原本不想将该行作为标题，最简单的处理方式是在末尾加上句号、冒号等标点，或将其拆成两行，使它成为普通段落。

### 链接 {#links}

一段没有缩进的非空行，如果每行都采用“[Text]: URL”的形式，就会定义链接目标。在同一段文档注释的其他文本中，“[Text]”表示使用指定文字链接到 URL，对应的 HTML 为 \<a href="URL">Text\</a>。例如：

	// Package json implements encoding and decoding of JSON as defined in
	// [RFC 7159]. The mapping between JSON and Go values is described
	// in the documentation for the Marshal and Unmarshal functions.
	//
	// For an introduction to this package, see the article
	// “[JSON and Go].”
	//
	// [RFC 7159]: https://tools.ietf.org/html/rfc7159
	// [JSON and Go]: https://golang.org/doc/articles/json_and_go.html
	package json


将 URL 集中放在独立部分，可以尽量减少对正文阅读的干扰。这种格式也大致对应 Markdown 的[简短引用链接格式](https://spec.commonmark.org/0.30/#shortcut-reference-link)，只是没有可选的标题文字。

如果没有相应的 URL 定义，那么除下一节介绍的符号文档链接外，“[Text]”不会变成超链接，显示时也会保留方括号。每段文档注释独立处理：一段注释中的链接目标定义不会影响其他注释。

虽然链接目标定义可以穿插在普通段落之间，但 gofmt 会将它们统一移到文档注释末尾，最多分成两个块：先放正文中引用过的目标，再放正文中_没有_引用过的目标。单独列出未使用目标，便于发现并修正拼写错误，或删除不再需要的定义。

纯文本中被识别为 URL 的内容，在 HTML 展示时会自动变成链接。

### 符号文档链接 {#doclinks}

符号文档链接使用“[Name1]”或“[Name1.Name2]”引用当前包中的导出标识符，使用“[pkg]”、“[pkg.Name1]”或“[pkg.Name1.Name2]”引用其他包或其中的标识符。

例如：

	package bytes

	// ReadFrom reads data from r until EOF and appends it to the buffer, growing
	// the buffer as needed. The return value n is the number of bytes read. Any
	// error except [io.EOF] encountered during the read is also returned. If the
	// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
	func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
		...
	}


符号链接方括号中的文字可以带一个前导星号，方便引用指针类型，例如 \[\*bytes.Buffer\]。

引用其他包时，“pkg”可以是完整导入路径，也可以是已有导入的推定包名。推定包名要么是重命名导入中的标识符，要么是 [goimports 推定的名称](https://pkg.go.dev/golang.org/x/tools/internal/imports#ImportPathToAssumedName)。如果推定不正确，goimports 会补上导入别名，因此这个规则基本适用于所有 Go 代码。例如，当前包导入 encoding/json 后，可以写“[json.Decoder]”代替“[encoding/json.Decoder]”，链接到 encoding/json 的 Decoder 文档。如果同一个包的不同源文件使用相同名称导入不同的包，这种简写就存在歧义，不能使用。

只有以域名（包含点的路径片段）开头，或者属于标准库包（如“[os]”、“[encoding/json]”）时，“pkg”才会被视为完整导入路径。例如，`[os.File]` 和 `[example.com/sys.File]` 都是文档链接，尽管后者会是无效链接；而 `[os/sys.File]` 不是，因为标准库没有 os/sys 包。

为了避免与 map、泛型及数组类型冲突，文档链接前后必须是标点、空格、制表符，或行首、行尾。例如，“map[ast.Expr]TypeAndValue”中不包含文档链接。

### 列表 {#lists}

列表由连续的缩进行和空行组成，第一条缩进行以无序列表标记或有序列表标记开头。如果没有这种标记，这些行通常会成为下一节介绍的代码块。

无序列表标记可以是星号、加号、减号或 Unicode 项目符号（*、+、-、•；分别为 U+002A、U+002B、U+002D、U+2022），后接空格或制表符，再接正文。在无序列表中，每个以列表标记开头的行都会开始一个新条目。

例如：

	package url

	// PublicSuffixList provides the public suffix of a domain. For example:
	//   - the public suffix of "example.com" is "com",
	//   - the public suffix of "foo1.foo2.foo3.co.uk" is "co.uk", and
	//   - the public suffix of "bar.pvt.k12.ma.us" is "pvt.k12.ma.us".
	//
	// Implementations of PublicSuffixList must be safe for concurrent use by
	// multiple goroutines.
	//
	// An implementation that always returns "" is valid and may be useful for
	// testing but it is not secure: it means that the HTTP server for foo.com can
	// set a cookie for bar.com.
	//
	// A public suffix list implementation is in the package
	// golang.org/x/net/publicsuffix.
	type PublicSuffixList interface {
		...
	}


有序列表标记由任意长度的十进制数字组成，后接句点或右括号，再接空格或制表符及正文。每个以数字列表标记开头的行都会开始一个新条目。条目编号会保持原样，不会重新编号。

例如：

	package path

	// Clean returns the shortest path name equivalent to path
	// by purely lexical processing. It applies the following rules
	// iteratively until no further processing can be done:
	//
	//  1. Replace multiple slashes with a single slash.
	//  2. Eliminate each . path name element (the current directory).
	//  3. Eliminate each inner .. path name element (the parent directory)
	//     along with the non-.. element that precedes it.
	//  4. Eliminate .. elements that begin a rooted path:
	//     that is, replace "/.." by "/" at the beginning of a path.
	//
	// The returned path ends in a slash only if it is the root "/".
	//
	// If the result of this process is an empty string, Clean
	// returns the string ".".
	//
	// See also Rob Pike, “[Lexical File Names in Plan 9].”
	//
	// [Lexical File Names in Plan 9]: https://9p.io/sys/doc/lexnames.html
	func Clean(path string) string {
		...
	}


列表条目只能包含段落，不支持代码块或嵌套列表。这样可以避免计算空格数量的微妙规则，也避免缩进不一致时需要判断一个制表符算几个空格。

Gofmt 会将无序列表统一为减号标记，标记前缩进两个空格，续行缩进四个空格。

对于有序列表，gofmt 会在数字前保留一个空格，数字后使用句点，续行同样缩进四个空格。

列表与前一个段落之间可以没有空行；如果存在，gofmt 会保留。列表与后面的段落或标题之间，gofmt 会插入一个空行。

### 代码块 {#code}

代码块由连续的缩进行和空行组成，但不以无序或有序列表标记开头。它会以预格式化文本展示，在 HTML 中对应 \<pre> 块。

代码块经常包含 Go 代码。例如：

{{raw `
	package sort

	// Search uses binary search...
	//
	// As a more whimsical example, this program guesses your number:
	//
	//	func GuessingGame() {
	//		var s string
	//		fmt.Printf("Pick an integer from 0 to 100.\n")
	//		answer := sort.Search(100, func(i int) bool {
	//			fmt.Printf("Is your number <= %d? ", i)
	//			fmt.Scanf("%s", &s)
	//			return s != "" && s[0] == 'y'
	//		})
	//		fmt.Printf("Your number is %d.\n", answer)
	//	}
	func Search(n int, f func(int) bool) int {
		...
	}
`}}

当然，代码块也常用于展示代码以外的预格式化文本。例如：

{{raw `
	package path

	// Match reports whether name matches the shell pattern.
	// The pattern syntax is:
	//
	//	pattern:
	//		{ term }
	//	term:
	//		'*'         matches any sequence of non-/ characters
	//		'?'         matches any single non-/ character
	//		'[' [ '^' ] { character-range } ']'
	//		            character class (must be non-empty)
	//		c           matches character c (c != '*', '?', '\\', '[')
	//		'\\' c      matches character c
	//
	//	character-range:
	//		c           matches character c (c != '\\', '-', ']')
	//		'\\' c      matches character c
	//		lo '-' hi   matches character c for lo <= c <= hi
	//
	// Match requires pattern to match all of name, not just a substring.
	// The only possible returned error is [ErrBadPattern], when pattern
	// is malformed.
	func Match(pattern, name string) (matched bool, err error) {
		...
	}
`}}

Gofmt 会将代码块所有行的共同缩进替换为一个制表符，并在每个代码块前后插入空行，使它与周围段落明显区分。

### 指令 {#directives}

`//go:generate` 这样的指令注释不属于文档注释，生成文档时会省略。Gofmt 会将指令注释移到文档注释末尾，并在前面加一个空行。例如：

	package regexp

	// An Op is a single regular expression operator.
	//
	//go:generate stringer -type Op -trimprefix Op
	type Op uint8


指令注释是以匹配正则表达式 `//(line |extern |export |[a-z0-9]+:[a-z0-9])` 的内容开头的行。

工具可以使用 `//toolname:directive arguments` 的形式定义自己的指令注释。工具指令匹配正则表达式 `//([a-z0-9]+):([a-z0-9]\PZ*)($|\pZ+)(.*)`，第一组为工具名称，第二组为指令名称。可选参数与指令名称之间用一个或多个 Unicode 空白字符分隔。各工具可以自定义参数语法，但常见约定是使用空格分隔参数，每个参数可以是普通单词，也可以是双引号或反引号包围的 Go 字符串。工具名称 `go` 保留给 Go 工具链使用。

[`go/ast.ParseDirective`](/pkg/go/ast#ParseDirective) 函数及其相关类型用于解析工具指令语法。

## 常见错误与陷阱 {#mistakes}

文档注释中连续的缩进行和空行会被渲染为代码块，这一规则从 Go 最早期就存在。遗憾的是，过去 gofmt 不支持文档注释格式化，导致许多已有注释使用了缩进，却并不是想创建代码块。

例如，下面这个未缩进的列表，一直会被 godoc 解释为三行段落，后接一行代码块：

	package http

	// cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	// 1) On Read error or close, the stop func is called.
	// 2) On Read failure, if reqDidTimeout is true, the error is wrapped and
	//    marked as net.Error that hit its timeout.
	type cancelTimerBody struct {
		...
	}


在 `go` `doc` 中，它一直显示为：

	cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	1) On Read error or close, the stop func is called. 2) On Read failure,
	if reqDidTimeout is true, the error is wrapped and

	    marked as net.Error that hit its timeout.


类似地，下面注释中的命令被解释为一行段落，后接一行代码块：

	package smtp

	// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:
	//
	// go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
	//     --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
	var localhostCert = []byte(`...`)


在 `go` `doc` 中显示为：

	localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:

	go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \

	    --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h


下面这个注释则是两行段落（第二行是“{”），后接六行缩进代码块，再接一行段落“}”：

	// On the wire, the JSON will look something like this:
	// {
	//	"kind":"MyAPIObject",
	//	"apiVersion":"v1",
	//	"myPlugin": {
	//		"kind":"PluginA",
	//		"aOption":"foo",
	//	},
	// }


在 `go` `doc` 中显示为：

	On the wire, the JSON will look something like this: {

	    "kind":"MyAPIObject",
	    "apiVersion":"v1",
	    "myPlugin": {
	    	"kind":"PluginA",
	    	"aOption":"foo",
	    },

	}


另一个常见错误是未缩进的 Go 函数定义或语句块，同样由“{”和“}”包围。

Go 1.19 的 gofmt 增加文档注释格式化功能后，会在代码块周围加上空行，使这些问题更容易被发现。

2022 年的一项分析发现，公开 Go 模块中只有 3% 的文档注释会被 Go 1.19 开发版本的 gofmt 改写。在这些被改写的注释中，约 87% 保持了人们阅读原注释时理解的结构；约 6% 则遇到了上述问题，包括未缩进的列表、多行 shell 命令和大括号包围的代码块。

基于这项分析，Go 1.19 的 gofmt 采用了一些启发式规则，将部分未缩进的行合并到相邻的缩进列表或代码块中。经过这些调整，上面的示例会被格式化为：

	// cancelTimerBody is an io.ReadCloser that wraps rc with two features:
	//  1. On Read error or close, the stop func is called.
	//  2. On Read failure, if reqDidTimeout is true, the error is wrapped and
	//     marked as net.Error that hit its timeout.

	// localhostCert is a PEM-encoded TLS cert generated from src/crypto/tls:
	//
	//	go run generate_cert.go --rsa-bits 1024 --host 127.0.0.1,::1,example.com \
	//	    --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

	// On the wire, the JSON will look something like this:
	//
	//	{
	//		"kind":"MyAPIObject",
	//		"apiVersion":"v1",
	//		"myPlugin": {
	//			"kind":"PluginA",
	//			"aOption":"foo",
	//		},
	//	}


这种格式化既让含义更清楚，也使文档注释在旧版本 Go 中正确展示。如果启发式规则判断有误，可以插入空行，明确分隔段落文本与其他内容，从而覆盖自动判断。

即使有这些启发式规则，仍有一些已有注释需要手动调整。最常见的问题是：原本没有缩进的正文换行后，续行却加了缩进。例如：

	// TODO Revisit this design. It may make sense to walk those nodes
	//      only once.

	// According to the document:
	// "The alignment factor (in bytes) that is used to align the raw data of sections in
	//  the image file. The value should be a power of 2 between 512 and 64 K, inclusive."


这两个例子的最后一行都因为缩进而变成代码块。修复方式是去掉这些行的缩进。

另一个常见问题是：列表或代码块本来有缩进，但换行后的续行没有保持缩进。例如：

	// Uses of this error model include:
	//
	//   - Partial errors. If a service needs to return partial errors to the
	// client,
	//     it may embed the `Status` in the normal response to indicate the
	// partial
	//     errors.
	//
	//   - Workflow errors. A typical workflow has multiple steps. Each step
	// may
	//     have a `Status` message for error reporting.


修复方式是为续行添加适当缩进。

Go 文档注释不支持嵌套列表，因此 gofmt 会将：

	// Here is a list:
	//
	//  - Item 1.
	//    * Subitem 1.
	//    * Subitem 2.
	//  - Item 2.
	//  - Item 3.


重新格式化为：

	// Here is a list:
	//
	//  - Item 1.
	//  - Subitem 1.
	//  - Subitem 2.
	//  - Item 2.
	//  - Item 3.


最佳解决方案通常是重写文本，避免嵌套列表，这也往往能改善文档质量。另一种可能的变通方式是混用列表标记，因为无序列表标记不会在有序列表中开始新条目，反过来也一样。例如：

	// Here is a list:
	//
	//  1. Item 1.
	//
	//     - Subitem 1.
	//
	//     - Subitem 2.
	//
	//  2. Item 2.
	//
	//  3. Item 3.

