<!--{
  "Title": "教程：使用 govulncheck 查找并修复存在漏洞的依赖",
  "HideTOC": true,
  "Breadcrumb": true
}-->

Govulncheck 是一个能减少无关告警的工具，帮助你查找并修复 Go 项目中存在漏洞的依赖。它会扫描项目依赖中的已知漏洞，再识别代码中是否存在直接或间接调用相关漏洞代码的路径。

本教程将介绍如何使用 govulncheck 扫描一个简单程序中的漏洞，还会说明如何评估漏洞并确定处理优先级，以便优先修复最重要的问题。

要进一步了解 govulncheck，请参阅 [govulncheck 文档](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)以及这篇介绍 Go [漏洞管理的博客文章](/blog/vuln)。我们也欢迎你[提供反馈](/s/govulncheck-feedback)。

## 准备工作 {#prerequisites}

- **Go。**建议使用最新版本的 Go 学习本教程。安装步骤请参阅[安装 Go](/doc/install)。
- **代码编辑器。**任何编辑器都可以。
- **命令行终端。**Go 可以在 Linux 和 Mac 的各种终端，以及 Windows 的 PowerShell 或 cmd 中正常使用。

本教程将依次完成以下步骤：

1. 创建一个包含存在漏洞的依赖的 Go 示例模块。
2. 安装并运行 govulncheck。
3. 评估漏洞。
4. 升级存在漏洞的依赖。

## 创建包含存在漏洞的依赖的 Go 示例模块 {#create-a-sample-go-module-with-a-vulnerable-dependency}

**第 1 步。**首先创建名为 `vuln-tutorial` 的文件夹，并初始化 Go 模块。（如果还不熟悉 Go 模块，请参阅[创建 Go 模块](/doc/tutorial/create-module)。）

例如，在用户主目录中运行：

```
$ mkdir vuln-tutorial
$ cd vuln-tutorial
$ go mod init vuln.tutorial
```

**第 2 步。**在 `vuln-tutorial` 文件夹中创建 `main.go` 文件，并将以下代码复制进去：

```
package main

import (
        "fmt"
        "os"

        "golang.org/x/text/language"
)

func main() {
        for _, arg := range os.Args[1:] {
                tag, err := language.Parse(arg)
                if err != nil {
                        fmt.Printf("%s: error: %v\n", arg, err)
                } else if tag == language.Und {
                        fmt.Printf("%s: undefined\n", arg)
                } else {
                        fmt.Printf("%s: tag %s\n", arg, tag)
                }
        }
}
```

这个示例程序接收一组作为命令行参数传入的语言标签，并为每个标签输出一条消息，说明它是否成功解析、是否未定义，或者解析过程中是否发生错误。

**第 3 步。**运行 `go mod tidy`，将上一步加入 `main.go` 的代码所需的全部依赖写入 `go.mod` 文件。

在 `vuln-tutorial` 文件夹中运行：

```
$ go mod tidy
```

你应该会看到如下输出：

```
go: finding module for package golang.org/x/text/language
go: downloading golang.org/x/text v0.9.0
go: found golang.org/x/text/language in golang.org/x/text v0.9.0
```

**第 4 步。**打开 `go.mod` 文件，确认内容如下：

```
module vuln.tutorial

go 1.20

require golang.org/x/text v0.9.0
```

**第 5 步。**将 `golang.org/x/text` 降级到包含已知漏洞的 v0.3.5 版本。运行：

```
$ go get golang.org/x/text@v0.3.5
```

你应该会看到如下输出：

```
go: downgraded golang.org/x/text v0.9.0 => v0.3.5
```

此时，`go.mod` 文件应为：

```
module vuln.tutorial

go 1.20

require golang.org/x/text v0.3.5
```

现在来看看 govulncheck 如何工作。

## 安装并运行 govulncheck {#install-and-run-govulncheck}

**第 6 步。**使用 `go install` 命令安装 govulncheck：

```
$ go install golang.org/x/vuln/cmd/govulncheck@latest
```

**第 7 步。**在要分析的文件夹中（这里是 `vuln-tutorial`）运行：

```
$ govulncheck ./...
```

你应该会看到如下输出：

```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-18 21:32:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
Your code is affected by 1 vulnerability from 1 module.

Vulnerability #1: GO-2021-0113
  Due to improper index calculation, an incorrectly formatted
  language tag can cause Parse to panic via an out of bounds read.
  If Parse is used to process untrusted user inputs, this may be
  used as a vector for a denial of service attack.

  More info: https://pkg.go.dev/vuln/GO-2021-0113

  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.7

    Call stacks in your code:
      main.go:12:29: vuln.tutorial.main calls golang.org/x/text/language.Parse

=== Informational ===

Found 1 vulnerability in packages that you import, but there are no call
stacks leading to the use of this vulnerability. You may not need to
take any action. See https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
for details.

Vulnerability #1: GO-2022-1059
  An attacker may cause a denial of service by crafting an
  Accept-Language header which ParseAcceptLanguage will take
  significant time to parse.
  More info: https://pkg.go.dev/vuln/GO-2022-1059
  Found in: golang.org/x/text@v0.3.5
  Fixed in: golang.org/x/text@v0.3.8

```

### 理解输出结果 {#interpreting-the-output}

<font size="2">  *注意：如果使用的不是最新版本的 Go，可能还会看到标准库中的其他漏洞。 </font>

代码受到一个漏洞 [GO-2021-0113](https://pkg.go.dev/vuln/GO-2021-0113) 的影响，因为它直接调用了 `golang.org/x/text/language` 中 `Parse` 函数的存在漏洞的版本（v0.3.5）。

另一个漏洞 [GO-2022-1059](https://pkg.go.dev/vuln/GO-2022-1059) 也存在于 `golang.org/x/text` 模块的 v0.3.5 版本中。不过，它被归为“Informational（提示信息）”，因为我们的代码既没有直接调用，也没有间接调用其中存在漏洞的函数。

接下来，评估这些漏洞并确定处理措施。

### 评估漏洞 {#evaluate-vulnerabilities}

a. 评估漏洞。

首先阅读漏洞描述，判断它是否实际影响你的代码及使用场景。如果需要更多信息，请访问“More info”后的链接。

根据描述，使用 `Parse` 处理不可信的用户输入时，GO-2021-0113 可能导致 panic（恐慌）。假设我们希望程序能够处理不可信输入，并且需要防范拒绝服务攻击，那么这个漏洞很可能会影响程序。

GO-2022-1059 很可能不会影响我们的代码，因为代码没有调用该报告中存在漏洞的函数。

b. 确定处理措施。

针对 GO-2021-0113，可以考虑以下方案：

- **方案 1：升级到已修复的版本。**如果已有修复版本，可以通过升级模块来消除这项依赖中的漏洞。
- **方案 2：停止使用存在漏洞的符号。**可以移除代码中对存在漏洞的函数的全部调用，但需要寻找替代实现，或者自行实现相应功能。

在本例中，已有修复版本，而且 `Parse` 函数是程序的核心组成部分。因此，可以将依赖升级到“Fixed in”所标示的已修复版本 v0.3.7。

虽然我们决定暂时降低提示类漏洞 GO-2022-1059 的修复优先级，但它与 GO-2021-0113 位于同一模块中，修复版本是 v0.3.8，所以直接升级到 v0.3.8 就能方便地同时消除这两个漏洞。

## 升级存在漏洞的依赖 {#upgrade-vulnerable-dependencies}

升级存在漏洞的依赖非常简单。

**第 8 步。**将 `golang.org/x/text` 升级到 v0.3.8：

```
$ go get golang.org/x/text@v0.3.8
```

你应该会看到如下输出：

```
go: upgraded golang.org/x/text v0.3.5 => v0.3.8
```

（也可以选择升级到 `latest`，或者 v0.3.8 之后的其他版本。）

**第 9 步。**再次运行 govulncheck：

```
$ govulncheck ./...
```

此时会看到如下输出：

```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-06 19:19:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
No vulnerabilities found.
```

在这个示例中，govulncheck 最终确认未发现漏洞。

定期使用 govulncheck 命令扫描依赖，可以帮助你识别漏洞、确定处理优先级并修复问题，从而保护代码库。
