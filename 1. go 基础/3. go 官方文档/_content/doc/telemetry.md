---
title: Go 遥测
layout: article
breadcrumb: true
date: 2024-02-07:00:00Z
template: true
---

<style>
.DocInfo {
  background-color: var(--color-background-info);
  padding: 1.5rem 2rem 1.5rem 4rem;
  border-left: 0.875rem solid var(--color-border);
  position: relative;
}
.DocInfo:before {
  content: "ⓘ";
  position: absolute;
  top: 1rem;
  left: 1rem;
  font-size: 2rem;
}
</style>

目录：

 [背景](#background)\
 [概述](#overview)\
 [配置](#config)\
 [计数器](#counters)\
 [报告与上传](#reports)\
 [图表](#charts) \
 [遥测提案](#proposals)\
 [IDE 提示](#ide) \
 [常见问题](#faq)

## 背景 {#background}

Go 遥测用于让 Go 工具链程序采集自身性能和使用情况的数据。这里的“Go 工具链”指 Go 团队维护的开发工具，包括 `go` 命令，以及 Go 语言服务器 [`gopls`]、Go 安全工具 [`govulncheck`] 等辅助工具。Go 遥测只供 Go 团队维护的程序及其选定的依赖（例如 [Delve]）使用。

默认情况下，遥测数据只保存在本机。用户可以主动选择向 [telemetry.go.dev] 上传经过批准的数据子集。上传的数据有助于 Go 团队了解工具的使用情况和故障，进而改进 Go 语言及其工具。

在开源软件领域，“遥测”一词常带有负面含义，很多情况下这种评价有其原因。不过，衡量用户体验是现代软件工程的重要组成部分，而 GitHub 问题报告、年度调查等数据来源较为粗略且滞后，无法回答 Go 团队需要了解的所有问题。Go 遥测旨在帮助工具链程序采集有关可靠性、性能和使用情况的有效数据，同时保持用户对 Go 项目所期待的透明度与隐私保护。其设计过程与动机见[遥测系列博客](https://research.swtch.com/telemetry)。遥测与隐私方面的说明见[遥测隐私政策](https://telemetry.go.dev/privacy)。

本页详细介绍 Go 遥测的工作方式。常见问题的简短解答见[常见问题](#faq)。

<div class="DocInfo">
使用 Go 1.23 或更高版本时，要<strong>主动启用</strong>向 Go 团队上传遥测数据，请运行：
<pre>
go telemetry on
</pre>
要完全禁用遥测，包括本地采集，请运行：
<pre>
go telemetry off
</pre>
要恢复为默认的仅本地遥测模式，请运行：
<pre>
go telemetry local
</pre>
Go 1.23 之前也可以使用 <code>golang.org/x/telemetry/cmd/gotelemetry</code> 命令设置，详情见<a href="#config">配置</a>。
</div>

## 概述 {#overview}

Go 遥测使用三种核心数据类型：

- [_计数器_](#counters)是对具名事件进行轻量计数的机制，通过插桩加入工具链程序。启用采集时（[模式](#config)为 **local** 或 **on**），计数器会写入本地文件系统中的内存映射文件。
- [_报告_](#reports)汇总某一周的计数器数据。启用上传时（[模式](#config)为 **on**），[经过批准的计数器](#proposals)报告会上传到 [telemetry.go.dev]，供公众访问。
- [_图表_](#charts)汇总所有用户上传的报告，可以在 [telemetry.go.dev] 查看。

所有本地 Go 遥测数据和配置都存放在 <code>[os.UserConfigDir()](/pkg/os#UserConfigDir)/go/telemetry</code> 目录。下文将其简称为 `<gotelemetry>`。

下图展示了数据流向。

<div class="image">
  <center>
    <img max-width="800px" src="/doc/telemetry/dataflow.png" />
  </center>
</div>

下文将逐一介绍图中的各个组成部分。首先了解控制这些行为的配置。

## 配置 {#config}

Go 遥测行为由一个值控制：遥测_模式_。`mode` 可以是 `local`（默认值）、`on` 或 `off`：

- `mode` 为 `local` 时，数据仅在本机采集和保存，绝不上传到远程服务器。
- `mode` 为 `on` 时，会采集数据，并根据[抽样策略](#uploads)决定是否上传。
- `mode` 为 `off` 时，既不采集，也不上传。

在 Go 1.23 或更高版本中，可以使用以下命令查看和设置遥测模式：

- `go telemetry`：查看当前模式。
- `go telemetry on`：将模式设为 `on`。
- `go telemetry off`：将模式设为 `off`。
- `go telemetry local`：将模式设为 `local`。

也可以通过只读的 Go 环境变量查看遥测配置：

- `go env GOTELEMETRY`：显示遥测模式。
- `go env GOTELEMETRYDIR`：显示存放遥测配置和数据的目录。

[`gotelemetry`](/pkg/golang.org/x/telemetry/cmd/gotelemetry) 命令也可以配置遥测模式、检查本地遥测数据。使用以下命令安装：

```
go install golang.org/x/telemetry/cmd/gotelemetry@latest
```

`gotelemetry` 命令行工具的完整用法见其[包文档](/pkg/golang.org/x/telemetry/cmd/gotelemetry)。

## 计数器 {#counters}

如前所述，Go 遥测通过_计数器_插桩。计数器分为两类：基本计数器和调用栈计数器。

### 基本计数器 {#basic-counters}

_基本计数器_是一个可以递增的值，其名称描述所统计的事件。例如，`gopls/client:vscode` 记录 VS Code 发起 `gopls` 会话的次数。还可以有 `gopls/client:neovim`、`gopls/client:eglot` 等计数器，用于记录不同编辑器或语言客户端的会话。如果一周内使用了多个编辑器，可能得到如下数据：

    gopls/client:vscode 8
    gopls/client:neovim 5
    gopls/client:eglot  2

当计数器按这种方式关联时，我们有时将 `:` 前面的部分称为_图表名称_（本例为 `gopls/client`），后面的部分称为_桶名称_（`vscode`）。[图表](#charts)一节将介绍这种约定的作用。

基本计数器也能表示_直方图_。例如，{{raw
`<code>gopls/completion/latency:&lt;50ms</code>`}} 计数器记录自动补全耗时小于 50 毫秒的次数。

{{raw `
<pre>
gopls/completion/latency:&lt;10ms
gopls/completion/latency:&lt;50ms
gopls/completion/latency:&lt;100ms
...
</pre>
`}}

这种记录直方图数据的方式是一种约定，桶名称 {{raw `<code>&lt;50ms</code>`}} 本身并无特殊含义。这类计数器常用于衡量性能。

### 调用栈计数器 {#stack-counters}

_调用栈计数器_在递增计数时，还会记录 Go 工具链程序当前的调用栈。例如，`crash/crash` 记录工具链程序崩溃时的调用栈：

    crash/crash
    golang.org/x/tools/gopls/internal/golang.hoverBuiltin:+22
    golang.org/x/tools/gopls/internal/golang.Hover:+94
    golang.org/x/tools/gopls/internal/server.Hover:+42
    ...

调用栈计数器通常用于统计程序不变量被破坏的事件。最常见的是崩溃，另一个例子是 `gopls/bug`：它统计开发者预先标记的异常情况，例如被恢复的 panic 或按理“不可能发生”的错误。调用栈计数器只包含 Go 工具链程序内部的函数名和行号，不包含用户输入的信息，例如用户源码的名称或内容。

调用栈计数器有助于定位其他渠道未能报告的罕见或棘手缺陷。引入 `gopls/bug` 之后，我们发现了[数十处](https://github.com/golang/go/issues?q=label%3Agopls%2Ftelemetry-wins)实际执行到了的“不可达”代码。追查这些异常，帮助我们发现并修复了许多影响用户、却不易察觉或难以报告的缺陷。尤其在预发布测试阶段，调用栈计数器让我们能够借助自动化更高效地改进产品。

### 计数器文件 {#counter-files}

所有计数器数据都写入 `<gotelemetry>/local` 目录，文件名遵循以下格式：

```
[program name]@[program version]-[go version]-[GOOS]-[GOARCH]-[date].v1.count
```

- **程序名**为 [debug.BuildInfo] 报告的程序包路径的最后一段。
- **程序版本**和 **Go 版本**也由 [debug.BuildInfo] 报告。
- **GOOS** 和 **GOARCH** 分别由 [`runtime.GOOS`](/pkg/runtime#GOOS) 和 [`runtime.GOARCH`](/pkg/runtime#GOARCH) 报告。
- **日期**为计数器文件的创建日期，格式为 `YYYY-MM-DD`。

这些文件会被映射到每个正在运行的已插桩程序实例的内存中。借助内存映射，即使程序立即崩溃，或者多个工具实例同时运行，计数器也能安全记录。

## 报告与上传 {#reports}

计数器数据大约每周汇总一次，在 `<gotelemetry>/local` 中生成名为 `<date>.json` 的报告。报告累加上一周的计数，并按与计数器文件相同的程序标识分组：程序名、程序版本、Go 版本、GOOS 和 GOARCH。

可以通过 [`gotelemetry view`](/pkg/golang.org/x/telemetry/cmd/gotelemetry) 命令，以图表形式查看本地报告。下面是 `gopls/completion/latency` 计数器的汇总示例：

<div class="image">
  <center>
    <img max-width="800px" src="/doc/telemetry/gopls-latency.png" />
  </center>
</div>

### 上传 {#uploads}

如果启用了遥测上传，每周生成报告时还会创建一份仅包含[上传配置](https://telemetry.go.dev/config)中列出的计数器子集的报告。这些计数器必须经过下一节介绍的公开评审流程批准。上传成功后，已上传报告的副本会保存在 `<gotelemetry>/upload` 目录。

当选择上传遥测数据的用户足够多时，上传流程会随机跳过一部分报告，在保持统计显著性的同时减少采集量、增强隐私保护。

## 图表 {#charts}

[telemetry.go.dev] 除了接收上传，还会公开这些数据。每天，上传的报告会被处理成两类输出，均可在 [telemetry.go.dev] 首页查看。

- _合并报告（merged）_汇总当天收到的所有上传数据中的计数器。
- _图表（charts）_依据提案流程产生的 [chart config] 绘制数据。回顾[计数器](#counters)一节，`foo:bar` 这样的名称可拆分为图表名 `foo` 和桶名 `bar`。每张图表会将具有相同图表名的计数器聚合到对应的桶中。

> 译注：原文中的 `[chart config]` 缺少链接定义，它指的是图表配置文件 [config.txt](https://go.googlesource.com/telemetry/+/refs/heads/master/internal/chartconfig/config.txt)。原引用保留以便对照。

图表采用 [chartconfig] 包定义的格式。例如，`gopls/client` 的图表配置如下：

    title: Editor Distribution
    counter: gopls/client:{vscode,vscodium,vscode-insiders,code-server,eglot,govim,neovim,coc.nvim,sublimetext,other}
    description: measure editor distribution for gopls users.
    type: partition
    issue: https://go.dev/issue/61038
    issue: https://go.dev/issue/62214 # add vscode-insiders
    program: golang.org/x/tools/gopls
    version: v0.13.0 # temporarily back-version to demonstrate config generation.

这份配置描述要生成的图表，列出需要聚合的计数器集合，并指定适用的程序版本。此外，[提案流程](#proposals)要求图表关联一项已获批准的提案。以下是该配置生成的图表：

<div class="image">
  <center>
    <img src="/doc/telemetry/gopls-clients.png" />
  </center>
</div>

## 遥测提案流程 {#proposals}

对上传配置或 [telemetry.go.dev] 图表集合的修改，必须经过_遥测提案流程_，以确保配置变更透明。

在这个流程中，上传配置与图表配置实际上并不分离。上传配置本身就是通过希望在 telemetry.go.dev 展示的聚合结果来表达的，遵循“只采集我们确实希望_查看_的数据”这一原则。

提案流程如下：

1. 提案者创建一个 CL，修改 [chartconfig] 包的 [config.txt]，加入所需的新计数器聚合方式。
2. 提案者提交 [proposal]，请求合并这个 CL。
3. 问题讨论结束后，由 Go 团队成员批准或拒绝提案。
4. 自动化流程重新生成上传配置，允许上传新图表所需的计数器。相关程序发布新版本后，该流程也会定期将新版本加入上传配置。

新图表必须不携带用户敏感信息，同时具备实用性和可行性，才可能获批。实用性意味着图表应服务于明确目的，并能据此采取行动；可行性意味着所需数据能可靠采集，测量结果具有统计显著性。为证明可行，提案者可能需要先在目标程序中加入计数器，并在本地采集数据。

所有相关提案都可以在 GitHub 的[提案项目](https://github.com/orgs/golang/projects/29)中查看。

## IDE 提示 {#ide}

为了让遥测回答我们关心的问题，主动选择上传的用户不必特别多：约 16,000 名参与者，就能在所需粒度上得到具有统计显著性的测量结果。但建立这样的样本仍有成本：我们需要询问大量 Go 开发者是否愿意参与。

即使现在有很多用户主动参与（例如读过 Go 博客之后），这些人也可能更多是有经验的 Go 开发者，而且这种样本偏差会随时间扩大。另外，用户更换电脑后，需要再次主动选择参与。在遥测系列博客中，这被称为主动选择模式的[“动员成本”](https://research.swtch.com/telemetry-opt-in#campaign)。

为了持续补充参与用户，Go 语言服务器 [`gopls`] 支持显示提示，询问用户是否愿意加入 Go 遥测。在 VS Code 中，提示如下：

<div class="image">
  <center>
    <img width="600px" src="/doc/telemetry/prompt.png" />
  </center>
</div>

如果用户选择“Yes”，遥测[模式](#config)就会设为 `on`，相当于运行了 [`gotelemetry on`](/pkg/golang.org/x/telemetry/cmd/gotelemetry)。这样既方便用户参与，也能持续覆盖规模较大、分布多样的 Go 开发者样本。

## 常见问题 {#faq}

**问：如何启用或禁用 Go 遥测？**

答：使用 `gotelemetry` 命令，可通过 `go install
golang.org/x/telemetry/cmd/gotelemetry@latest` 安装。运行 `gotelemetry off` 可以禁用全部功能，包括本地采集；运行 `gotelemetry on` 可以启用全部功能，包括向 [telemetry.go.dev] 上传经过批准的计数器。详情见[配置](#config)。

**问：本地数据保存在哪里？**

答：保存在 <code>[os.UserConfigDir()](/pkg/os#UserConfigDir)/go/telemetry</code> 目录中。

**问：选择参与后，多久上传一次数据？**

答：大约每周一次。

**问：选择参与后，会上传哪些数据？**

答：只有[上传配置](https://telemetry.go.dev/config)中列出的计数器才能上传。该配置由 [chart config] 生成，后者通常更易阅读。

**问：如何将计数器加入上传配置？**

答：通过[公开提案流程](#proposals)。

**问：在哪里可以查看已上传的遥测数据？**

答：可以在 [telemetry.go.dev] 以图表或合并汇总的形式查看。

**问：Go 遥测的源码在哪里？**

答：位于 [golang.org/x/telemetry](/pkg/golang.org/x/telemetry)。

[`gopls`]: /pkg/golang.org/x/tools/gopls
[`govulncheck`]: /pkg/golang.org/x/vuln/cmd/govulncheck
[Delve]: /pkg/github.com/go-delve/delve#section-readme
[debug.BuildInfo]: /pkg/runtime/debug#BuildInfo
[proposal]: /issue/new?assignees=&labels=Telemetry-Proposal&projects=golang%2F29&template=12-telemetry.yml&title=x%2Ftelemetry%2Fconfig%3A+proposal+title
[telemetry.go.dev]: https://telemetry.go.dev
[chartconfig]: /pkg/golang.org/x/telemetry/internal/chartconfig
[config.txt]: https://go.googlesource.com/telemetry/+/refs/heads/master/internal/chartconfig/config.txt
