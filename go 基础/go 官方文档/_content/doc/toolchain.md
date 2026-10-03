---
title: "Go 工具链"
layout: article
---

## 简介 {#intro}

从 Go 1.21 开始，Go 发行版由 `go` 命令及其随附的 Go 工具链组成。工具链包含标准库，以及编译器、汇编器和其他工具。
`go` 命令既可以使用随附的工具链，也可以使用在本地 `PATH` 中找到或按需下载的其他版本。

使用哪个 Go 工具链取决于 `GOTOOLCHAIN` 环境设置，以及主模块 `go.mod` 文件或当前工作区 `go.work` 文件中的 `go` 行和 `toolchain` 行。
在不同主模块和工作区之间切换时，使用的工具链版本也可能变化，就像模块依赖的版本一样。

在标准配置下，如果随附工具链的版本不低于主模块或工作区中 `go` 行和 `toolchain` 行指定的版本，`go` 命令就使用随附的工具链。
例如，在声明了 `go 1.21.0` 的主模块中运行 Go 1.21.3 随附的 `go` 命令时，`go` 命令使用的是 Go 1.21.3。
如果 `go` 行或 `toolchain` 行指定的版本比随附工具链更新，`go` 命令就改用更新的工具链。
例如，在声明了 `go 1.21.9` 的主模块中运行 Go 1.21.3 随附的 `go` 命令时，`go` 命令会查找并运行 Go 1.21.9。
`go` 命令先在 PATH 中查找名为 `go1.21.9` 的程序，找不到时就下载并缓存 Go 1.21.9 工具链。
可以禁用这种自动切换工具链的行为；但为准确保证前向兼容性，在禁用后，如果主模块或工作区的 `go` 行要求更新的 Go 版本，`go` 命令将拒绝运行。
也就是说，`go` 行设置了使用模块或工作区所需的最低 Go 版本。

供其他模块依赖的模块，可能需要把最低 Go 版本要求设得较低，而在直接开发该模块时希望使用较新的工具链。
这种情况下，可以在 `go.mod` 或 `go.work` 的 `toolchain` 行中设置首选工具链。
`go` 命令决定使用哪个工具链时，这一行的优先级高于 `go` 行。

可以把 `go` 行和 `toolchain` 行理解为模块对 Go 工具链本身的依赖版本要求，就像 `go.mod` 中的 `require` 行指定了对其他模块的依赖版本要求一样。
`go get` 命令管理 Go 工具链依赖的方式，与管理其他模块依赖的方式相同。
例如，`go get go@latest` 会更新模块，使其要求最新发布的 Go 工具链。

`GOTOOLCHAIN` 环境设置可以强制使用指定的 Go 版本，覆盖 `go` 行和 `toolchain` 行的设置。例如，要使用 Go 1.21rc3 测试一个包：

	GOTOOLCHAIN=go1.21rc3 go test

`GOTOOLCHAIN` 的默认值是 `auto`，会启用前面介绍的工具链切换行为。
另一种形式 `<name>+auto` 会先设置默认工具链，再决定是否进一步切换。
例如，`GOTOOLCHAIN=go1.21.3+auto` 让 `go` 命令默认选择 Go 1.21.3，但在 `go` 行和 `toolchain` 行要求时，仍会使用更新的工具链。
由于可以通过 `go env -w` 更改默认的 `GOTOOLCHAIN` 设置，如果你已经安装了 Go 1.21.0 或更高版本，那么运行

	go env -w GOTOOLCHAIN=go1.21.3+auto

就相当于把安装的 Go 1.21.0 替换成 Go 1.21.3。

本文接下来会详细介绍 Go 工具链的版本规则、选择方式和管理方法。

## Go 版本 {#version}

Go 已发布版本使用“1.*N*.*P*”的版本格式，表示 Go 1.*N* 的第 *P* 次发布。
首次发布的版本是 1.*N*.0，例如“1.21.0”。之后的版本（例如 1.*N*.9）通常称为补丁版本。

Go 1.*N* 的候选发布版本在 1.*N*.0 之前发布，使用“1.*N*rc*R*”的版本格式。
Go 1.*N* 的第一个候选发布版本是 1.*N*rc1，例如 `1.23rc1`。

“1.*N*”这种格式称为“语言版本”，表示实现该版本 Go 语言和标准库的整个 Go 发行版本系列。

将 Go 版本中 *N* 之后的内容全部截去，就得到对应的语言版本：1.21、1.21rc2 和 1.21.3 实现的都是语言版本 1.21。

Go 1.21.0 和 Go 1.21rc1 等已发布工具链，通过 `go version` 和 [`runtime.Version`](/pkg/runtime/#Version) 返回具体版本，例如 `go1.21.0` 或 `go1.21rc1`。
从 Go 开发仓库构建的未发布工具链（仍处于开发阶段）则只返回语言版本，例如 `go1.21`。

任意两个 Go 版本都可以比较大小，确定一个版本小于、大于还是等于另一个版本。
如果语言版本不同，比较结果就由语言版本决定：1.21.9 < 1.22。
同一语言版本内，从小到大的顺序是：语言版本本身、按 *R* 排序的候选发布版本、按 *P* 排序的正式发布版本。

例如，1.21 < 1.21rc1 < 1.21rc2 < 1.21.0 < 1.21.1 < 1.21.2。

在 Go 1.21 之前，Go 工具链首次发布的版本是 1.*N*，而不是 1.*N*.0。
因此，对于 *N* < 21 的情况，排序规则会进行调整，将 1.*N* 排在候选发布版本之后。

例如，1.20rc1 < 1.20rc2 < 1.20rc3 < 1.20 < 1.20.1。

较早的 Go 版本还发布过 beta 测试版，例如 1.18beta2。
在版本排序中，beta 测试版紧接在候选发布版本之前。

例如，1.18beta1 < 1.18beta2 < 1.18rc1 < 1.18 < 1.18.1。

<!-- Unpublished note: the download page also lists Go 1.9.2rc2, which does not respect
this version syntax. That was created as a test of some potential release automation
before Go 1.9.2 but is not considered a “real” toolchain. -->

## Go 工具链名称 {#name}

标准 Go 工具链的名称为 <code>go<i>V</i></code>，其中 *V* 是表示 beta 测试版、候选发布版本或正式发布版本的 Go 版本。
例如，`go1.21rc1` 和 `go1.21.0` 是工具链名称；
`go1.21` 和 `go1.22` 则不是（首次发布的版本是 `go1.21.0` 和 `go1.22.0`），但 `go1.20` 和 `go1.19` 是。

非标准工具链使用 <code>go<i>V</i>-<i>suffix</i></code> 形式的名称，后缀可以任意指定。

比较工具链时，会比较名称中包含的版本 <code><i>V</i></code>：去掉开头的 `go`，并丢弃以 `-` 开头的后缀。
例如，在版本排序中，`go1.21.0` 和 `go1.21.0-custom` 视为相等。

## 模块和工作区配置 {#config}

Go 模块和工作区在各自的 `go.mod` 或 `go.work` 文件中指定与版本有关的配置。

`go` 行声明使用模块或工作区所需的最低 Go 版本。
出于兼容性考虑，如果 `go.mod` 文件省略了 `go` 行，就认为该模块隐含了 `go 1.16` 行；
如果 `go.work` 文件省略了 `go` 行，就认为该工作区隐含了 `go 1.18` 行。

`toolchain` 行声明建议在模块或工作区中使用的工具链。
如后面的“[Go 工具链的选择](#select)”所述，如果默认工具链的版本低于建议的版本，`go` 命令在该模块或工作区中执行操作时可能会运行指定的工具链。
如果省略 `toolchain` 行，就认为模块或工作区隐含了 <code>toolchain go<i>V</i></code> 行，其中 *V* 是 `go` 行中的 Go 版本。

例如，声明了 `go 1.21.0` 而没有 `toolchain` 行的 `go.mod` 文件，会被视为包含了 `toolchain go1.21.0` 行。

如果模块或工作区声明的最低 Go 版本高于工具链自身版本，Go 工具链会拒绝加载它。

例如，Go 1.21.2 会拒绝加载包含 `go 1.21.3` 或 `go 1.22` 行的模块或工作区。

模块的 `go` 行声明的版本，必须大于或等于每个 `require` 语句所列模块声明的 `go` 版本。
工作区的 `go` 行声明的版本，必须大于或等于每个 `use` 语句所列模块声明的 `go` 版本。

例如，如果模块 *M* 依赖模块 *D*，而 *D* 的 `go.mod` 声明了 `go 1.22.0`，那么 *M* 的 `go.mod` 就不能声明 `go 1.21.3`。

每个模块的 `go` 行设置了编译器编译该模块中的包时，所强制采用的语言版本。
可以通过[构建约束](/cmd/go#hdr-Build_constraints)为单个文件更改语言版本：
如果文件包含构建约束，而且该约束所要求的最低版本至少为 `go1.21`，那么编译该文件时使用的语言版本就是这个最低版本。

例如，模块中的代码如果使用 Go 1.21 语言版本，`go.mod` 文件中的 `go` 行就应写为 `go 1.21` 或 `go 1.21.3` 等形式。
如果某个源文件只能用更新的 Go 工具链编译，在其中添加 `//go:build go1.22`，既能保证只有 Go 1.22 及更新的工具链会编译该文件，也会将该文件的语言版本设置为 Go 1.22。

修改 `go` 行和 `toolchain` 行最方便、安全的方式是使用 `go get`；详见[下面专门介绍 `go get` 的章节](#get)。

在 Go 1.21 之前，Go 工具链把 `go` 行视为建议性要求：如果构建成功，就认为一切正常；否则，会输出可能存在版本不匹配的提示。
从 Go 1.21 起，`go` 行改为强制性要求。
这一行为也部分回移到了较早的语言版本：从 Go 1.19.13 开始的 Go 1.19 版本，以及从 Go 1.20.8 开始的 Go 1.20 版本，会拒绝加载声明 Go 1.22 或更高版本的工作区或模块。

在 Go 1.21 之前，工具链并不要求模块或工作区的 `go` 行版本大于或等于其各个依赖模块要求的 `go` 版本。

## `GOTOOLCHAIN` 设置 {#GOTOOLCHAIN}

`go` 命令根据 `GOTOOLCHAIN` 设置选择要使用的 Go 工具链。
查找 `GOTOOLCHAIN` 设置时，`go` 命令遵循适用于所有 Go 环境设置的标准规则：

 - 如果进程环境中的 `GOTOOLCHAIN` 设为非空值（通过 [`os.Getenv`](/pkg/os/#Getenv) 查询），`go` 命令就使用该值。

 - 否则，如果用户的默认环境设置文件中设置了 `GOTOOLCHAIN`（该文件通过 [`go env -w` 和 `go env -u`](/cmd/go/#hdr-Print_Go_environment_information) 管理），`go` 命令就使用该值。

 - 否则，如果随附 Go 工具链的默认环境设置文件（`$GOROOT/go.env`）中设置了 `GOTOOLCHAIN`，`go` 命令就使用该值。

在标准 Go 工具链中，`$GOROOT/go.env` 文件将默认值设置为 `GOTOOLCHAIN=auto`，但重新打包的 Go 工具链可能会更改这个值。

如果 `$GOROOT/go.env` 文件不存在，或者其中没有设置默认值，`go` 命令就假定使用 `GOTOOLCHAIN=local`。

运行 `go env GOTOOLCHAIN` 可以输出 `GOTOOLCHAIN` 设置。

## Go 工具链的选择 {#select}

`go` 命令在启动时选择要使用的 Go 工具链。
它会查看 `GOTOOLCHAIN` 设置，其格式为 `<name>`、`<name>+auto` 或 `<name>+path`。
`GOTOOLCHAIN=auto` 是 `GOTOOLCHAIN=local+auto` 的简写；同样，`GOTOOLCHAIN=path` 是 `GOTOOLCHAIN=local+path` 的简写。
`<name>` 指定默认 Go 工具链：`local` 表示随附的 Go 工具链（即与当前运行的 `go` 命令一同发布的工具链）；其他情况下，`<name>` 必须是具体的 Go 工具链名称，例如 `go1.21.0`。
`go` 命令优先运行默认的 Go 工具链。
如前所述，从 Go 1.21 开始，如果工作区或模块要求更新的 Go 版本，Go 工具链会拒绝运行，报告错误并退出。

当 `GOTOOLCHAIN` 设为 `local` 时，`go` 命令始终运行随附的 Go 工具链。

当 `GOTOOLCHAIN` 设为 `<name>`（例如 `GOTOOLCHAIN=go1.21.0`）时，`go` 命令始终运行指定的 Go 工具链。
如果在系统 PATH 中找到了对应名称的二进制程序，`go` 命令就使用它；否则，`go` 命令会下载并验证 Go 工具链后使用。

当 `GOTOOLCHAIN` 设为 `<name>+auto` 或 `<name>+path`（或其简写 `auto`、`path`）时，`go` 命令会按需选择并运行更新的 Go 版本。
具体来说，它会查看当前工作区 `go.work` 文件中的 `toolchain` 行和 `go` 行；如果不在工作区中，则查看主模块的 `go.mod` 文件。
如果 `go.work` 或 `go.mod` 文件包含 `toolchain <tname>` 行，而且 `<tname>` 比默认 Go 工具链更新，`go` 命令就改为运行 `<tname>`。
如果文件包含 `toolchain default` 行，`go` 命令就运行默认 Go 工具链，并禁用任何升级到比 `<name>` 更新版本的尝试。
否则，如果文件包含 `go <version>` 行，而且 `<version>` 比默认 Go 工具链更新，`go` 命令就改为运行 `go<version>`。

为了运行随附工具链之外的 Go 工具链，`go` 命令会在进程的可执行文件搜索路径（Unix 和 Plan 9 上的 `$PATH`、Windows 上的 `%PATH%`）中查找指定名称（例如 `go1.21.3`）的程序，并运行该程序。
如果找不到这样的程序，`go` 命令就[下载并运行指定的 Go 工具链](#download)。
使用 `<name>+path` 形式的 `GOTOOLCHAIN` 会禁用下载这一后备方案，使 `go` 命令在搜索可执行文件路径后就停止。

运行 `go version` 会输出所选 Go 工具链的版本，具体做法是运行所选工具链实现的 `go version` 命令。

运行 `GOTOOLCHAIN=local go version` 会输出随附 Go 工具链的版本。

从 Go 1.24 开始，运行 `go` 命令时可以在 `GODEBUG` 环境变量中添加 `toolchaintrace=1`，跟踪 `go` 命令选择工具链的过程。

## Go 工具链的切换 {#switch}

对于大多数命令，由于[配置要求](#config)中对版本大小关系的规定，工作区 `go.work` 或主模块 `go.mod` 中的 `go` 行，所指定的版本至少与所有依赖模块的 `go` 行一样新。
因此，启动时选择的 Go 工具链已足够新，可以完成该命令。

有些命令在执行过程中会引入新的模块版本：
`go get` 会向主模块添加新的模块依赖；
`go work use` 会向工作区添加新的本地模块；
`go work sync` 会将工作区与本地模块重新同步，这些模块可能在工作区创建之后已被更新；
`go install package@version` 和 `go run package@version` 实际上会在空的主模块中运行，并将 `package@version` 添加为新依赖。
这些命令都可能遇到这样的模块：其 `go.mod` 中的 `go` 行要求的 Go 版本，比当前运行的版本更新。

如果命令遇到要求更新 Go 版本的模块，而 `GOTOOLCHAIN` 又允许运行其他工具链（设为 `auto` 或 `path` 形式），`go` 命令就会选择并切换到合适的更新工具链，继续执行当前命令。

只要 `go` 命令在启动时选择工具链之后又切换了工具链，就会输出一条消息说明原因。例如：

	go: module example.com/widget@v1.2.3 requires go >= 1.24rc1; switching to go 1.27.9

如示例所示，`go` 命令可能切换到比新发现的最低要求更高的工具链版本。
通常，`go` 命令的目标是切换到仍受支持的 Go 工具链。

为了选择工具链，`go` 命令先获取可用工具链列表。
对于 `auto` 形式，`go` 命令会下载可用工具链列表。
对于 `path` 形式，`go` 命令会扫描 PATH，查找名称符合有效工具链格式的可执行文件，并使用找到的所有工具链组成的列表。
`go` 命令从列表中确定至多三个候选版本：

 - 尚未正式发布的 Go 语言版本的最新候选发布版本（1.*N*₃rc*R*₃）；
 - 最近正式发布的 Go 语言版本的最新补丁版本（1.*N*₂.*P*₂）；
 - 前一个 Go 语言版本的最新补丁版本（1.*N*₁.*P*₁）。

根据 Go 的[发布政策](/doc/devel/release#policy)，这些都是仍受支持的 Go 版本。
随后，按照[最小版本选择](https://research.swtch.com/vgo-mvs)的原则，`go` 命令会保守地选用满足新要求的候选版本中_最小_（最旧）的版本。

例如，假设 `example.com/widget@v1.2.3` 要求 Go 1.24rc1 或更高版本。
`go` 命令获取可用工具链列表后，发现最近两个正式 Go 版本的最新补丁版本是 Go 1.28.3 和 Go 1.27.9，另外还有候选发布版本 Go 1.29rc2。
这种情况下，`go` 命令会选择 Go 1.27.9。
如果 `widget` 要求 Go 1.28 或更高版本，`go` 命令就会选择 Go 1.28.3，因为 Go 1.27.9 太旧。
如果 `widget` 要求 Go 1.29 或更高版本，`go` 命令就会选择 Go 1.29rc2，因为 Go 1.27.9 和 Go 1.28.3 都太旧。

命令在引入要求更新 Go 版本的新模块版本时，会更新当前工作区 `go.work` 文件或主模块 `go.mod` 文件的 `go` 行，写入新的最低 Go 版本要求。
为保证[可重复性](https://research.swtch.com/vgo-principles#repeatability)，任何更新 `go` 行的命令也都会更新 `toolchain` 行，记录自身的工具链名称。
下一次在该工作区或模块中运行 `go` 命令时，就会在[选择工具链](#select)时使用更新后的 `toolchain` 行。

例如，`go get example.com/widget@v1.2.3` 可能会输出前面展示的切换提示，并切换到 Go 1.27.9。
Go 1.27.9 会完成 `go get` 操作，并将 `toolchain` 行更新为 `toolchain go1.27.9`。
下次在该模块或工作区中运行 `go` 命令时，会在启动阶段就选择 `go1.27.9`，因此不会输出切换消息。

一般来说，如果将同一个 `go` 命令运行两次，第一次输出了切换消息，第二次就不会再输出，因为第一次也已更新 `go.work` 或 `go.mod`，让命令可以在启动时选择正确的工具链。
例外是 `go install package@version` 和 `go run package@version` 这两种形式：它们运行时不处于任何工作区或主模块中，无法写入 `toolchain` 行。
因此，它们每次需要切换到更新工具链时，都会输出切换消息。

## 下载工具链 {#download}

使用 `GOTOOLCHAIN=auto` 或 `GOTOOLCHAIN=<name>+auto` 时，Go 命令会按需下载更新的工具链。
这些工具链打包成了特殊的模块，模块路径为 `golang.org/toolchain`，版本为 <code>v0.0.1-go<i>VERSION</i>.<i>GOOS</i>-<i>GOARCH</i></code>。
工具链的下载方式与其他模块相同，这意味着可以通过设置 `GOPROXY` 为工具链下载使用代理，并由 Go 校验和数据库验证其校验和。
由于具体使用的工具链既取决于系统自身的默认工具链，也取决于本地操作系统和架构（GOOS 和 GOARCH），将工具链模块的校验和写入 `go.sum` 并不现实。
因此，如果设置了 `GOSUMDB=off`，工具链下载会因无法验证而失败。
`GOPRIVATE` 和 `GONOSUMDB` 的匹配模式不适用于工具链下载。

## 使用 `go get` 管理模块的 Go 版本要求 {#get}

通常，`go` 命令将 `go` 行和 `toolchain` 行视为主模块对工具链的带版本依赖声明。
`go get` 命令可以管理这两行，就像管理用于指定模块依赖版本的 `require` 行一样。

例如，`go get go@1.22.1 toolchain@1.24rc1` 会修改主模块的 `go.mod` 文件，使其包含 `go 1.22.1` 和 `toolchain go1.24rc1`。

`go` 命令知道，`go` 依赖要求 `toolchain` 依赖的 Go 版本大于或等于自身版本。

接着上面的例子，如果之后运行 `go get go@1.25.0`，工具链也会更新为 `go1.25.0`。
当工具链与 `go` 行完全一致时，就可以省略它，使用隐含值，因此这次 `go get` 会删除 `toolchain` 行。

降级时，同样的约束会反向生效：如果 `go.mod` 原本包含 `go 1.22.1` 和 `toolchain go1.24rc1`，
那么 `go get toolchain@go1.22.9` 只会更新 `toolchain` 行；而 `go get toolchain@go1.21.3` 还会将 `go` 行降级为 `go 1.21.3`。
最终文件中只留下 `go 1.21.3`，没有 `toolchain` 行。

特殊形式 `toolchain@none` 表示删除所有 `toolchain` 行，例如 `go get toolchain@none` 或 `go get go@1.25.0 toolchain@none`。

`go` 命令同时支持 `go` 和 `toolchain` 依赖的版本格式与版本查询。

例如，就像 `go get example.com/widget@v1.2` 使用 `example.com/widget` 最新的 `v1.2` 版本（可能是 `v1.2.3`）一样，
`go get go@1.22` 会使用 Go 1.22 语言版本中最新的可用发行版本（可能是 `1.22rc3`，也可能是 `1.22.3`）。
`go get toolchain@go1.22` 也遵循同样的规则。

`go get` 和 `go mod tidy` 命令会维护 `go` 行，确保其版本大于或等于任何必需依赖模块的 `go` 行版本。

例如，如果主模块声明了 `go 1.22.1`，而我们运行 `go get example.com/widget@v1.2.3`，且该依赖声明了 `go 1.24rc1`，
那么 `go get` 会将主模块的 `go` 行更新为 `go 1.24rc1`。

接着这个例子，如果之后运行 `go get go@1.22.1`，它就会把 `example.com/widget` 降级到与 Go 1.22.1 兼容的版本，或者直接删除该依赖要求，
就像降级 `example.com/widget` 的其他依赖时一样。

在 Go 1.21 之前，要将模块更新为新的 Go 版本（比如 Go 1.22），推荐的做法是运行 `go mod tidy -go=1.22`，
确保在更新 `go` 行的同时，也完成 Go 1.22 对 `go.mod` 所需的调整。
这种形式仍然有效，但现在更推荐使用简单的 `go get go@1.22`。

如果在工作区根目录之内的模块目录中运行 `go get`，`go get` 大体上会忽略工作区；
但如果操作后工作区的 `go` 行版本会过低，它也会更新 `go.work` 文件，提高 `go` 行的版本。

## 使用 `go work` 管理工作区的 Go 版本要求 {#work}

如上一节所述，在工作区根目录之内运行 `go get` 时，会按需更新 `go.work` 文件的 `go` 行，确保其版本大于或等于该根目录内的所有模块。
但工作区也可以引用根目录之外的模块；在这些目录中运行 `go get`，可能导致工作区配置失效，也就是 `go.work` 声明的 `go` 版本低于一个或多个 `use` 指令所引用模块的版本。

用于添加新 `use` 指令的 `go work use` 命令，也会检查 `go.work` 文件中的 `go` 版本，确保它满足所有现有 `use` 指令的要求。
如果工作区的 `go` 版本与其模块已经不同步，可以运行不带参数的 `go work use` 来更新工作区。

`go work init` 和 `go work sync` 命令也会按需更新 `go` 版本。

要删除 `go.work` 文件中的 `toolchain` 行，请使用 `go work edit -toolchain=none`。
