---
title: Go 模糊测试
layout: article
breadcrumb: true
---

Go 从 1.18 开始在标准工具链中支持模糊测试。原生 Go 模糊测试也已获得 [OSS-Fuzz 支持](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support)。

**动手体验 [Go 模糊测试教程](/doc/tutorial/fuzz)。**

## 概述 {#overview}

模糊测试（fuzzing）是一种自动化测试方法，通过持续变异程序输入来寻找缺陷。Go 模糊测试以覆盖率为引导，智能探索被测代码，发现失败并报告给用户。它能够触及人工测试容易遗漏的边界情况，因此在发现可利用的安全缺陷和漏洞方面尤其有价值。

下面的[模糊测试](#glos-fuzz-test)示例标出了它的主要组成部分。

<img class="DarkMode-img" alt="模糊测试示例：测试中包含模糊测试目标，在目标之前通过 f.Add 添加语料，目标的参数被标注为模糊测试参数。"
src="/security/fuzz/example-dark.png" style="width: 600px; height:
auto;"/>
<img alt="模糊测试示例：测试中包含模糊测试目标，在目标之前通过 f.Add 添加语料，目标的参数被标注为模糊测试参数。"
src="/security/fuzz/example.png" style="width: 600px; height:
auto;" class="LightMode-img"/>

## 编写模糊测试 {#writing-fuzz-tests}

### 必须遵守的规则 {#requirements}

模糊测试必须满足以下要求：

- 测试函数的名称采用 `FuzzXxx` 形式，只接收一个 `*testing.F` 参数，没有返回值。
- 测试必须位于 \*\_test.go 文件中才能运行。
- [模糊测试目标](#glos-fuzz-target)通过调用 <code>[(\*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz)</code> 指定。传入的目标函数以 `*testing.T` 为第一个参数，后续为模糊测试参数，没有返回值。
- 每个模糊测试必须且只能有一个模糊测试目标。
- 所有[种子语料](#glos-seed-corpus)条目的类型及其顺序，都必须与[模糊测试参数](#glos-fuzzing-arguments)完全一致。这既适用于 <code>[(\*testing.F).Add](https://pkg.go.dev/testing#F.Add)</code> 调用，也适用于测试的 testdata/fuzz 目录中的语料文件。
- 模糊测试参数仅支持以下类型：
  - `string`、`[]byte`
  - `int`、`int8`、`int16`、`int32`/`rune`、`int64`
  - `uint`、`uint8`/`byte`、`uint16`、`uint32`、`uint64`
  - `float32`、`float64`
  - `bool`

### 建议 {#suggestions}

以下建议有助于发挥模糊测试的效果：

- 模糊测试目标应执行快速、结果确定，以便测试引擎高效运行，也便于复现新发现的失败和代码覆盖情况。
- 测试目标会由多个工作进程并行调用，调用顺序不确定。因此，每次调用结束后不应留下影响后续调用的状态，目标的行为也不应依赖全局状态。

## 运行模糊测试 {#running-fuzz-tests}

模糊测试有两种运行模式：作为单元测试运行（默认的 `go test`），或者启用模糊测试引擎（`go test -fuzz=FuzzTestName`）。

默认情况下，模糊测试的运行方式与单元测试相近：将每个[种子语料条目](#glos-seed-corpus)传给测试目标，报告失败后退出。

要启用模糊测试，请为 `go test` 添加 `-fuzz` 标志，提供一个恰好匹配某个模糊测试的正则表达式。默认情况下，开始模糊测试前，会先运行该包中的其他所有测试，避免重复报告现有测试已经能够发现的问题。

注意，运行多久由你决定。如果没有发现错误，模糊测试完全可能一直运行下去。未来将支持通过 OSS-Fuzz 等工具持续运行这些测试，参见[问题 #50192](/issue/50192)。

**注意：**应在支持覆盖率插桩的平台上运行模糊测试（目前为 AMD64 和 ARM64），这样语料集才能在运行中有效增长，覆盖更多代码。

### 命令行输出 {#command-line-output}

模糊测试运行时，[模糊测试引擎](#glos-fuzzing-engine)生成新输入并执行给定的测试目标。默认会一直运行，直到发现[失败输入](#glos-failing-input)，或用户取消进程，例如按 Ctrl+C。

输出大致如下：

```
~ go test -fuzz FuzzFoo
fuzz: elapsed: 0s, gathering baseline coverage: 0/192 completed
fuzz: elapsed: 0s, gathering baseline coverage: 192/192 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 325017 (108336/sec), new interesting: 11 (total: 202)
fuzz: elapsed: 6s, execs: 680218 (118402/sec), new interesting: 12 (total: 203)
fuzz: elapsed: 9s, execs: 1039901 (119895/sec), new interesting: 19 (total: 210)
fuzz: elapsed: 12s, execs: 1386684 (115594/sec), new interesting: 21 (total: 212)
PASS
ok      foo 12.692s
```

最开始几行表示，模糊测试开始前正在收集“基线覆盖率”。

为收集基线覆盖率，引擎会执行[种子语料](#glos-seed-corpus)和[生成语料](#glos-generated-corpus)，确保它们没有产生错误，并了解现有语料已经能够覆盖哪些代码。

后续各行反映当前运行情况：

- elapsed：从进程启动到当前经过的时间。
- execs：已向测试目标传入并执行的输入总数，括号中是自上一条日志以来的平均每秒执行次数。
- new interesting：本次运行中新加入生成语料集的“有价值”输入总数，括号中是整个语料集的总条目数。

要被视为“有价值”，输入必须拓展现有生成语料集尚未触及的代码覆盖。通常，新的有价值输入数量起初增长很快，之后逐渐放缓；发现新分支时，偶尔会出现一阵快速增长。

随着语料中的输入逐渐覆盖更多代码，通常会看到“new interesting”的增长速度下降；引擎发现新代码路径时，也可能再次短暂加速。

### 失败输入 {#failing-input}

模糊测试可能因以下原因失败：

- 代码或测试中触发了 panic。
- 测试目标直接调用 `t.Fail`，或通过 `t.Error`、`t.Fatal` 等方法报告失败。
- 发生无法恢复的错误，例如 `os.Exit` 或栈溢出。
- 测试目标执行时间过长。目前单次执行的超时时间为 1 秒。原因可能是死锁、死循环，也可能是代码本来就需要较长时间。这也是[建议目标执行快速](#suggestions)的原因之一。

发生错误时，引擎会尝试将输入最小化，得到尽可能小、便于阅读且仍能触发错误的值。相关配置见[自定义设置](#custom-settings)。

最小化完成后，会记录错误信息，输出结尾类似：

```
    Failing input written to testdata/fuzz/FuzzFoo/a878c3134fe0404d44eb1e662e5d8d4a24beb05c3d68354903670ff65513ff49
    To re-run:
    go test -run=FuzzFoo/a878c3134fe0404d44eb1e662e5d8d4a24beb05c3d68354903670ff65513ff49
FAIL
exit status 1
FAIL    foo 0.839s
```

引擎已将这个[失败输入](#glos-failing-input)写入该测试的种子语料集，此后默认运行 `go test` 时就会执行它。修复缺陷后，它也就成为回归测试。

接下来应诊断问题、修复缺陷、重新运行 `go test` 验证修复，并随补丁一并提交新增的 testdata 文件，作为回归测试用例。

### 自定义设置 {#custom-settings}

go 命令的默认设置适合大多数模糊测试场景。因此，典型的命令行用法如下：

```
$ go test -fuzz={FuzzTestName}
```

`go` 命令也提供一些模糊测试选项，完整说明见 [`cmd/go` 包文档](https://pkg.go.dev/cmd/go)。

常用选项包括：

- `-fuzztime`：退出前运行测试目标的总时长或迭代次数，默认无限制。
- `-fuzzminimizetime`：每次尝试最小化输入时，运行测试目标的时长或迭代次数，默认 60 秒。设为 `-fuzzminimizetime 0` 可完全禁用最小化。
- `-parallel`：同时运行的模糊测试进程数量，默认为 `$GOMAXPROCS`。目前模糊测试期间设置 -cpu 不起作用。

## 语料文件格式 {#corpus-file-format}

语料文件采用专门的编码格式，[种子语料](#glos-seed-corpus)和[生成语料](#glos-generated-corpus)使用同一种格式。

示例语料文件如下：

```
go test fuzz v1
[]byte("hello\\xbd\\xb2=\\xbc ⌘")
int64(572293)
```

第一行告知模糊测试引擎该文件的编码版本。虽然目前没有计划推出新版本编码格式，但设计上必须支持这种可能性。

之后每一行都是组成该语料条目的值，需要时可以直接复制到 Go 代码中。

上例先是 `[]byte`，然后是 `int64`。它们的类型和顺序必须与模糊测试参数完全一致。对应的测试目标如下：

```
f.Fuzz(func(*testing.T, []byte, int64) {})
```

指定自定义种子语料最简单的方式，是调用 `(*testing.F).Add`。对于上面的示例，可以写成：

```
f.Add([]byte("hello\\xbd\\xb2=\\xbc ⌘"), int64(572293))
```

不过，如果有较大的二进制文件，可能不想将其作为代码复制进测试，而是希望作为独立种子语料条目保存在 testdata/fuzz/{FuzzTestName} 目录中。golang.org/x/tools/cmd/file2fuzz 提供的 [`file2fuzz`](https://pkg.go.dev/golang.org/x/tools/cmd/file2fuzz) 工具，可将这些二进制文件转换为以 `[]byte` 编码的语料文件。

使用该工具：

```
$ go install golang.org/x/tools/cmd/file2fuzz@latest
$ file2fuzz -h
```

## 资源 {#resources}

- **教程**：
  - 通过 [Go 模糊测试教程](/doc/tutorial/fuzz)深入理解新概念。
  - 更简短的入门介绍见[这篇博客文章](/blog/fuzz-beta)。
- **文档**：
  - [`testing`](https://pkg.go.dev//testing#hdr-Fuzzing) 包文档介绍编写模糊测试时使用的 `testing.F` 类型。
  - [`cmd/go`](https://pkg.go.dev/cmd/go) 包文档介绍模糊测试相关标志。
- **技术细节**：
  - [设计草案](/s/draft-fuzzing-design)
  - [提案](/issue/44551)

## 术语表 {#glossary}

<a id="glos-corpus-entry"></a>
**语料条目（corpus entry）：**语料集中可用于模糊测试的一组输入，可以是特定格式的文件，也可以通过调用 <code>[(\*testing.F).Add](https://pkg.go.dev/testing#F.Add)</code> 提供。

<a id="glos-coverage-guidance"></a>
**覆盖率引导（coverage guidance）：**根据代码覆盖率是否扩大，判断哪些语料条目值得保留供后续使用的模糊测试方法。

<a id="glos-failing-input"></a>
**失败输入（failing input）：**传给[模糊测试目标](#glos-fuzz-target)运行时，会引发错误或 panic 的语料条目。

<a id="glos-fuzz-target"></a>
**模糊测试目标（fuzz target）：**模糊测试中针对语料条目和生成值执行的函数，通过将函数传给 <code>[(\*testing.F).Fuzz](https://pkg.go.dev/testing#F.Fuzz)</code> 来指定。

<a id="glos-fuzz-test"></a>
**模糊测试（fuzz test）：**测试文件中形如 `func FuzzXxx(*testing.F)`、可用于执行模糊测试的函数。

<a id="glos-fuzzing"></a>
**模糊测试过程（fuzzing）：**不断变异程序输入以发现缺陷或潜在[漏洞](#glos-vulnerability)的一种自动化测试方法。

<a id="glos-fuzzing-arguments"></a>
**模糊测试参数（fuzzing arguments）：**传递给测试目标、由[变异器](#glos-mutator)改变其值的参数及其类型。

<a id="glos-fuzzing-engine"></a>
**模糊测试引擎（fuzzing engine）：**管理模糊测试的工具，包括维护语料、调用变异器、识别新的覆盖以及报告失败。

<a id="glos-generated-corpus"></a>
**生成语料（generated corpus）：**引擎在模糊测试过程中持续维护、用于保存进展的语料集，存储在 `$GOCACHE`/fuzz 中。这些条目只在启用模糊测试时使用。

<a id="glos-mutator"></a>
**变异器（mutator）：**在将语料条目传给测试目标之前，随机修改其内容的工具。

<a id="glos-package"></a>
**包（package）：**同一目录中一起编译的一组源文件，详见 Go 语言规范的[包](/ref/spec#Packages)一节。

<a id="glos-seed-corpus"></a>
**种子语料（seed corpus）：**用户为模糊测试提供的语料集，用于引导引擎。它由测试中的 f.Add 调用提供的条目，以及包内 testdata/fuzz/{FuzzTestName} 目录中的文件组成。无论是否启用模糊测试，默认运行 `go test` 时都会执行这些条目。

<a id="glos-test-file"></a>
**测试文件（test file）：**名称形如 xxx_test.go 的文件，可以包含测试、基准测试、示例和模糊测试。

<a id="glos-vulnerability"></a>
**漏洞（vulnerability）：**代码中可被攻击者利用的安全薄弱点。

## 反馈 {#feedback}

如遇到问题或有功能建议，请[提交问题](/issue/new?&labels=fuzz)。

如需讨论或提出一般反馈，也可以加入 Gophers Slack 的 [#fuzzing 频道](https://gophers.slack.com/archives/CH5KV1AKE)。
