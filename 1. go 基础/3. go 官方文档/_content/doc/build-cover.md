---
title: 集成测试的覆盖率分析
layout: article
template: true
---

目录：

 [概述](#overview)\
 [构建用于覆盖率分析的二进制文件](#building)\
 [运行插桩后的二进制文件](#running)\
 [处理覆盖率数据文件](#working)\
 [常见问题](#FAQ)\
 [相关资源](#resources)\
 [术语表](#glossary)


从 Go 1.20 开始，Go 支持从应用程序以及集成测试中采集覆盖率数据。集成测试是针对 Go 程序执行的规模更大、更加复杂的测试。

# 概述 {#overview}

Go 提供了简单易用的包级单元测试覆盖率采集功能，可以通过 `go test -coverprofile=... <pkg_target>` 命令使用。从 Go 1.20 开始，还可以采集规模更大的[集成测试](#glos-integration-test)的覆盖率：这类测试更复杂、开销更大，通常会多次运行同一个应用程序二进制文件。

对于单元测试，采集覆盖率并生成报告需要两步：先运行 `go test -coverprofile=...`，再调用 `go tool cover {-func,-html}` 生成报告。

对于集成测试，需要以下三步：[构建](#building)、[运行](#running)（可能多次调用构建所得的二进制文件），最后[生成报告](#reporting)。下文分别介绍。

# 构建用于覆盖率分析的二进制文件 {#building}

要构建能够采集覆盖率数据的应用程序，请在针对目标程序调用 `go build` 时传入 `-cover` 标志。[下文](#packageselection)给出了 `go build -cover` 的调用示例。运行生成的二进制文件时，再通过环境变量指定覆盖率数据的保存位置（参见下一节[运行](#running)）。

## 如何选择要插桩的包 {#packageselection}

执行 `go build -cover` 时，Go 命令默认只对主模块中的包进行覆盖率插桩；其他参与构建的包（go.mod 中列出的依赖以及 Go 标准库中的包）默认不包含在内。

例如，下面这个简单程序包含一个 main 包、主模块中的本地包 `greetings`，以及从模块外导入的若干包，包括 `rsc.io/quote` 和 `fmt` 等（[完整程序](/play/p/VSQJN8xkkf-?v=gotip)）。

```
$ cat go.mod
module mydomain.com

go 1.20

require rsc.io/quote v1.5.2

require (
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
	rsc.io/sampler v1.3.0 // indirect
)

$ cat myprogram.go
package main

import (
	"fmt"
	"mydomain.com/greetings"
	"rsc.io/quote"
)

func main() {
	fmt.Printf("I say %q and %q\n", quote.Hello(), greetings.Goodbye())
}
$ cat greetings/greetings.go
package greetings

func Goodbye() string {
	return "see ya"
}
$ go build -cover -o myprogram.exe .
$
```

如果使用 `-cover` 命令行标志构建并运行这个程序，覆盖率数据只会包含两个包：`main` 和 `mydomain.com/greetings`；其他依赖包不会包含在内。

如果希望更精确地控制哪些包参与覆盖率采集，可以在构建时使用 `-coverpkg` 标志。例如：

```
$ go build -cover -o myprogramMorePkgs.exe -coverpkg=io,mydomain.com,rsc.io/quote .
$
```

上述构建选择了 `mydomain.com` 中的主包，以及 `rsc.io/quote` 和 `io` 包来采集覆盖率。由于没有明确列出 `mydomain.com/greetings`，即使它位于主模块中，也不会被纳入覆盖率数据。

# 运行插桩后的二进制文件 {#running}

使用 `-cover` 构建的二进制文件会在执行结束时，将覆盖率数据文件写入环境变量 `GOCOVERDIR` 指定的目录。例如：

```
$ go build -cover -o myprogram.exe myprogram.go
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
I say "Hello, world." and "see ya"
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$
```

注意写入 `somedata` 目录的两个文件：这些二进制文件包含覆盖率结果。如何从这些数据文件生成便于阅读的结果，请参见下文[生成报告](#reporting)。

如果没有设置 `GOCOVERDIR` 环境变量，经过覆盖率插桩的程序仍然可以正常执行，但会发出警告。例如：

```
$ ./myprogram.exe
warning: GOCOVERDIR not set, no coverage data emitted
I say "Hello, world." and "see ya"
$
```

## 需要多次运行程序的测试 {#tests-involving-multiple-runs}

集成测试经常需要多次运行程序。如果程序使用 `-cover` 构建，每次运行都会生成一个新的数据文件。例如：

```
$ mkdir somedata2
$ GOCOVERDIR=somedata2 ./myprogram.exe          // first run
I say "Hello, world." and "see ya"
$ GOCOVERDIR=somedata2 ./myprogram.exe -flag    // second run
I say "Hello, world." and "see ya"
$ ls somedata2
covcounters.890814fca98ac3a4d41b9bd2a7ec9f7f.2456041.1670259309405583534
covcounters.890814fca98ac3a4d41b9bd2a7ec9f7f.2456047.1670259309410891043
covmeta.890814fca98ac3a4d41b9bd2a7ec9f7f
$
```

覆盖率输出文件分为两类：元数据文件保存各次运行中不变的信息，例如源文件名和函数名；计数器数据文件记录程序中哪些部分实际执行过。

在上例中，第一次运行生成了两个文件（计数器数据和元数据），第二次运行则只生成计数器数据文件。因为元数据不会随运行次数改变，所以只需写入一次。

# 处理覆盖率数据文件 {#working}

Go 1.20 引入了新工具 `covdata`，用于读取和处理 `GOCOVERDIR` 目录中的覆盖率数据文件。

Go 的 `covdata` 工具有多种工作模式，一般调用形式如下：

```
$ go tool covdata <mode> -i=<dir1,dir2,...> ...flags...
```

其中，`-i` 标志指定要读取的目录列表。每个目录都存放了运行覆盖率插桩程序时通过 `GOCOVERDIR` 输出的数据。

## 生成覆盖率报告 {#reporting}

本节介绍如何使用 `go tool covdata`，从覆盖率数据文件生成便于阅读的报告。

### 报告语句覆盖率 {#reporting-percent-statements-covered}

要报告每个插桩包的“已覆盖语句百分比”，请使用 `go tool covdata percent -i=<directory>` 命令。以上文[运行](#running)一节中的程序为例：

```
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata percent -i=somedata
	main	coverage: 100.0% of statements
	mydomain.com/greetings	coverage: 100.0% of statements
$
```

这里的语句覆盖率与 `go test -cover` 报告的百分比含义相同。

## 转换为旧版文本格式 {#converting-to-legacy-text-format}

可以使用 covdata 的 `textfmt` 子命令，将二进制覆盖率数据转换为 `go test -coverprofile=<outfile>` 生成的旧版文本格式。随后可以将该文本文件交给 `go tool cover -func` 或 `go tool cover -html`，生成其他形式的报告。例如：

```
$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata textfmt -i=somedata -o profile.txt
$ cat profile.txt
mode: set
mydomain.com/myprogram.go:10.13,12.2 1 1
mydomain.com/greetings/greetings.go:3.23,5.2 1 1
$ go tool cover -func=profile.txt
mydomain.com/greetings/greetings.go:3:	Goodbye		100.0%
mydomain.com/myprogram.go:10:		main		100.0%
total:					(statements)	100.0%
$
```

## 合并数据 {#merging}

`go tool covdata` 的 `merge` 子命令可以合并多个数据目录中的覆盖率数据。

例如，一个程序同时在 macOS 和 Windows 上运行。开发者可能希望将两个操作系统上各次运行的覆盖率数据合并为一份数据集，以生成跨平台的覆盖率汇总。例如：

```
$ ls windows_datadir
covcounters.f3833f80c91d8229544b25a855285890.1025623.1667481441036838252
covcounters.f3833f80c91d8229544b25a855285890.1025628.1667481441042785007
covmeta.f3833f80c91d8229544b25a855285890
$ ls macos_datadir
covcounters.b245ad845b5068d116a4e25033b429fb.1025358.1667481440551734165
covcounters.b245ad845b5068d116a4e25033b429fb.1025364.1667481440557770197
covmeta.b245ad845b5068d116a4e25033b429fb
$ ls macos_datadir
$ mkdir merged
$ go tool covdata merge -i=windows_datadir,macos_datadir -o merged
$
```

上述合并操作会读取指定输入目录中的数据，将其合并后写入 `merged` 目录，生成一组新的数据文件。

## 选择包 {#package-selection}

大多数 `go tool covdata` 命令支持使用 `-pkg` 标志选择要处理的包。`-pkg` 参数的格式与 Go 命令的 `-coverpkg` 标志相同。例如：

```

$ ls somedata
covcounters.c6de772f99010ef5925877a7b05db4cc.2424989.1670252383678349347
covmeta.c6de772f99010ef5925877a7b05db4cc
$ go tool covdata percent -i=somedata -pkg=mydomain.com/greetings
	mydomain.com/greetings	coverage: 100.0% of statements
$ go tool covdata percent -i=somedata -pkg=nonexistentpackage
$
```

使用 `-pkg` 标志可以为报告选出所关注的包子集。

#

## 常见问题 {#FAQ}

1. [如何对 `go.mod` 文件中涉及的所有导入包进行覆盖率插桩？](#gomodselect)
2. [可以在 GOPATH/GO111MODULE=off 模式下使用 `go build -cover` 吗？](#gopathmode)
3. [程序发生 panic 时还会写入覆盖率数据吗？](#panicprof)
4. [`-coverpkg=main` 会选择我的 main 包来采集覆盖率吗？](#mainpkg)


#### 如何对 `go.mod` 文件中涉及的所有导入包进行覆盖率插桩？ {#gomodselect}

默认情况下，`go build -cover` 会对主模块中的所有包进行覆盖率插桩，但不会对主模块之外的导入包插桩，例如标准库包或 `go.mod` 中列出的依赖。要将所有非标准库依赖都纳入插桩，可以把 `go list` 的输出传给 `-coverpkg`。以下仍使用前面介绍的[示例程序](/play/p/VSQJN8xkkf-?v=gotip)：

```
$ go list -f '{{"{{if not .Standard}}{{.ImportPath}}{{end}}"}}' -deps . | paste -sd "," > pkgs.txt
$ go build -o myprogram.exe -coverpkg=`cat pkgs.txt` .
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
$ go tool covdata percent -i=somedata
	golang.org/x/text/internal/tag	coverage: 78.4% of statements
	golang.org/x/text/language	coverage: 35.5% of statements
	mydomain.com	coverage: 100.0% of statements
	mydomain.com/greetings	coverage: 100.0% of statements
	rsc.io/quote	coverage: 25.0% of statements
	rsc.io/sampler	coverage: 86.7% of statements
$
```

#### 可以在 GO111MODULE=off 模式下使用 `go build -cover` 吗？ {#gopathmode}

可以，`go build -cover` 支持 `GO111MODULE=off`。在这种模式下构建程序时，默认只会对命令行中明确指定为构建目标的包进行覆盖率插桩。要包含其他包，请使用 `-coverpkg` 标志。

#### 程序发生 panic 时还会写入覆盖率数据吗？ {#panicprof}

使用 `go build -cover` 构建的程序，只有在调用 `os.Exit()` 或从 `main.main` 正常返回时，才会在执行结束时写出完整的覆盖率数据。如果程序因未恢复的 panic 而终止，或者遇到致命异常（例如段错误、除零等），本次运行中已执行语句的覆盖率数据就会丢失。

#### `-coverpkg=main` 会选择我的 main 包来采集覆盖率吗？ {#mainpkg}

`-coverpkg` 标志接受的是导入路径列表，而不是包名列表。如果要对 `main` 包进行覆盖率插桩，请通过导入路径指定它，不要使用包名。例如（使用[这个示例程序](/play/p/VSQJN8xkkf-?v=gotip)）：

```
$ go list -m
mydomain.com
$ go build -coverpkg=main -o oops.exe .
warning: no packages being built depend on matches for pattern main
$ go build -coverpkg=mydomain.com -o myprogram.exe .
$ mkdir somedata
$ GOCOVERDIR=somedata ./myprogram.exe
I say "Hello, world." and "see ya"
$ go tool covdata percent -i=somedata
	mydomain.com	coverage: 100.0% of statements
$
```

## 相关资源 {#resources}

- **介绍 Go 1.2 单元测试覆盖率的博客文章**：
  - 单元测试覆盖率分析在 Go 1.2 中引入，详见[这篇博客](/blog/cover)。
- **文档**：
  - [`cmd/go`](https://pkg.go.dev/cmd/go) 包文档介绍了与覆盖率相关的构建和测试标志。
- **技术细节**：
  - [设计草案](/design/51430-revamp-code-coverage)
  - [提案](/issue/51430)

## 术语表 {#glossary}

<a id="glos-unit-test"></a>
**单元测试（unit test）：** 使用 Go 的 `testing` 包，在某个 Go 包对应的 `*_test.go` 文件中编写的测试。

<a id="glos-integration-test"></a>
**集成测试（integration test）：** 针对应用程序或二进制文件开展的更全面、开销更大的测试。通常先构建一个或一组程序，再由测试框架控制程序，使用多种输入和场景进行一系列运行。该测试框架可以基于 Go 的 `testing` 包，也可以采用其他实现。

