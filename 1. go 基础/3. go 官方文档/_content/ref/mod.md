<!--{
  "Template": true,
  "Title": "Go 模块参考"
}-->
<!-- TODO(golang.org/issue/33637): Write focused "guide" articles on specific
module topics and tasks. Link to those instead of the blog, which will probably
not be updated over time. -->

## 引言 {#introduction}

模块是 Go 管理依赖的方式。

本文是 Go 模块系统的详细参考手册。有关创建 Go 项目的入门说明，请参阅[如何编写 Go 代码](/doc/code.html)。有关使用模块、将项目迁移到模块以及其他主题，请参阅以[使用 Go 模块](/blog/using-go-modules)为开篇的博客系列。

## 模块、包与版本 {#modules-overview}

<dfn>模块</dfn>是一组统一发布、管理版本和分发的包。模块可以直接从版本控制仓库下载，也可以从模块代理服务器下载。

模块通过[模块路径](#glos-module-path)标识，该路径与模块的依赖信息一起声明在 [`go.mod` 文件](#go-mod-file)中。<dfn>模块根目录</dfn>是包含 `go.mod` 文件的目录。<dfn>主模块</dfn>是包含执行 `go` 命令所在目录的模块。

模块中的每个<dfn>包</dfn>，都是同一目录下一起编译的一组源文件。<dfn>包路径</dfn>由模块路径与包所在子目录相对于模块根目录的路径拼接而成。例如，模块 `"golang.org/x/net"` 的 `"html"` 目录中有一个包，其包路径为 `"golang.org/x/net/html"`。

### 模块路径 {#module-path}

<dfn>模块路径</dfn>是模块的规范名称，通过模块 [`go.mod` 文件](#glos-go-mod-file)中的 [`module` 指令](#go-mod-file-module)声明。模块路径是该模块内所有包路径的前缀。

模块路径应同时说明模块的用途以及获取位置。通常，它由仓库根路径、仓库内目录（通常为空）和主版本后缀组成；主版本后缀仅用于主版本号为 2 或更高的版本。

* <dfn>仓库根路径</dfn>是模块路径中对应于版本控制仓库根目录的部分，该仓库用于开发模块。大多数模块定义在仓库根目录中，所以它通常就是完整模块路径。例如，`golang.org/x/net` 是同名模块的仓库根路径。有关 `go` 命令如何根据模块路径发起 HTTP 请求并定位仓库，请参阅[根据模块路径查找仓库](#vcs-find)。
* 如果模块不在仓库根目录中，<dfn>模块子目录</dfn>就是模块路径中指定目录的部分，不包含主版本后缀。它同时也是语义化版本标签的前缀。例如，模块 `golang.org/x/tools/gopls` 位于 `gopls` 子目录中，所属仓库的根路径为 `golang.org/x/tools`，因此模块子目录为 `gopls`。请参阅[将版本映射到提交](#vcs-version)和[仓库中的模块目录](#vcs-dir)。
* 如果模块发布的主版本号为 2 或更高，模块路径必须以 `/v2` 这样的[主版本后缀](#major-version-suffixes)结尾。该后缀可以是实际子目录名的一部分，也可以不是。例如，路径为 `golang.org/x/repo/sub/v2` 的模块，可以位于仓库中的 `/sub` 或 `/sub/v2` 子目录，该仓库为 `golang.org/x/repo`。

如果模块可能被其他模块依赖，就必须遵循这些规则，确保 `go` 命令能够找到并下载它。模块路径中允许使用的字符还受到一些[词法限制](#go-mod-file-ident)。

如果模块永远不会作为其他模块的依赖下载，可以使用任意合法包路径作为模块路径，但必须避免与其依赖或 Go 标准库可能使用的路径冲突。Go 标准库的包路径第一段不包含点号，`go` 命令不会尝试从网络服务器解析这种路径。`example` 和 `test` 路径专门保留给用户，标准库不会使用它们，适合用于独立模块，例如教程、示例代码，以及测试过程中创建和操作的模块。

### 版本 {#versions}

<dfn>版本</dfn>标识模块的一个不可变快照，可以是[正式版本](#glos-release-version)，也可以是[预发布版本](#glos-pre-release-version)。每个版本以字母 `v` 开头，后面跟语义化版本号。有关格式、含义和比较规则，请参阅[语义化版本 2.0.0](https://semver.org/spec/v2.0.0.html)。

语义化版本由三个以点号分隔的非负整数组成，从左到右分别为主版本号、次版本号和补丁版本号。补丁版本号后可以添加以连字符开头的预发布字符串；预发布字符串或补丁版本号之后，还可以添加以加号开头的构建元数据字符串。例如，`v0.0.0`、`v1.12.134`、`v8.0.5-pre` 和 `v2.0.9+meta` 都是有效版本。

版本的各个部分说明该版本是否稳定，以及是否与之前的版本兼容。

* 当模块的公共接口或文档承诺的功能发生不向后兼容的更改，例如移除某个包时，必须递增[主版本号](#glos-major-version)，并将次版本号和补丁版本号归零。
* 发生向后兼容的更改，例如新增函数时，必须递增[次版本号](#glos-minor-version)，并将补丁版本号归零。
* 发生不影响模块公共接口的更改，例如修复缺陷或优化时，必须递增[补丁版本号](#glos-patch-version)。
* 预发布后缀表示这是[预发布版本](#glos-pre-release-version)。预发布版本排在对应正式版本之前，例如 `v1.2.3-pre` 小于 `v1.2.3`。
* 比较版本时忽略构建元数据后缀。go 命令接受包含构建元数据的版本，并将它们转换为伪版本，以维持版本间的全序关系。
  * 特殊后缀 `+incompatible` 表示主版本号为 2 或更高、且在迁移到模块之前发布的版本，详见[兼容非模块仓库](#non-module-compat)。
  * 使用 Go 1.24 或更高版本工具链，在有效的本地版本控制系统（VCS）仓库中构建二进制文件时，如果工作目录包含未提交的更改，会在其版本信息中附加特殊后缀 `+dirty`。

主版本号为 0，或包含预发布后缀的版本，被视为不稳定版本，不受兼容性要求约束。例如，`v0.2.0` 可能不兼容 `v0.1.0`，`v1.5.0-beta` 也可能不兼容 `v1.5.0`。

Go 可以通过不符合上述约定的标签、分支或修订版本访问版本控制系统中的模块。不过，在主模块中，`go` 命令会自动将不符合标准的修订名称转换为规范版本，`go` 命令也会在此过程中移除构建元数据后缀，`+incompatible` 除外。这可能生成一个[伪版本](#glos-pseudo-version)：一种将版本控制系统中的修订标识符（例如 Git 提交哈希）和时间戳编码进去的预发布版本。例如，`go get
golang.org/x/net@daa7c041` 会将提交哈希 `daa7c041` 转换为伪版本 `v0.0.0-20191109021931-daa7c04131f5`。主模块之外必须使用规范版本；如果 `go.mod` 文件中出现 `master` 这样的非规范版本，`go` 命令会报错。

### 伪版本 {#pseudo-versions}

<dfn>伪版本</dfn>是一种具有特殊格式的[预发布](#glos-pre-release-version)[版本](#glos-version)，其中编码了版本控制仓库内某次特定修订的信息。例如，`v0.0.0-20191109021931-daa7c04131f5` 就是伪版本。

伪版本可以引用尚无[语义化版本标签](#glos-semantic-version-tag)的修订。例如，在开发分支上创建版本标签之前，可以用伪版本测试某个提交。

每个伪版本都包含三部分：

* 基础版本前缀（`vX.0.0` 或 `vX.Y.Z-0`）：从该修订之前的语义化版本标签推导而来；如果没有这种标签，就使用 `vX.0.0`。
* 时间戳（`yyyymmddhhmmss`）：修订创建时的 UTC 时间。在 Git 中使用提交时间，而不是作者时间。
* 修订标识符（`abcdefabcdef`）：提交哈希的前 12 个字符；对于 Subversion，则是补零后的修订号。

根据基础版本的不同，伪版本有三种形式。它们保证伪版本大于其基础版本，但小于下一个带标签的版本。

* `vX.0.0-yyyymmddhhmmss-abcdefabcdef`：用于没有已知基础版本的情况。与所有版本一样，主版本号 `X` 必须与模块的[主版本后缀](#glos-major-version-suffix)匹配。
* `vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef`：用于基础版本为 `vX.Y.Z-pre` 这样的预发布版本的情况。
* `vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef`：用于基础版本为 `vX.Y.Z` 这样的正式版本的情况。例如，基础版本为 `v1.2.3` 时，伪版本可以是 `v1.2.4-0.20191109021931-daa7c04131f5`。

使用不同基础版本时，多个伪版本可以引用同一个提交。在已经生成伪版本之后，才为更早的提交补上较低版本标签，就会自然出现这种情况。

这些形式赋予伪版本两个实用特性：

* 具有已知基础版本的伪版本，大于该基础版本，但小于后续版本的其他预发布版本。
* 具有相同基础版本前缀的伪版本，按时间先后排序。

`go` 命令会进行多项检查，确保模块作者能够控制伪版本与其他版本的相对顺序，并确保伪版本引用的修订确实属于模块的提交历史。

* 如果指定了基础版本，就必须存在对应的语义化版本标签，而且它必须是伪版本所指修订的祖先。这可以防止开发者用大于所有标签版本的伪版本，例如 `v1.999.999-99999999999999-daa7c04131f5`，绕过[最小版本选择](#glos-minimal-version-selection)。
* 时间戳必须与修订时间戳一致，防止攻击者用无限多个只有时间戳不同的伪版本淹没[模块代理](#glos-module-proxy)，也防止模块使用者改变版本间的相对顺序。
* 该修订必须是模块仓库某个分支或标签的祖先，以防攻击者引用尚未批准的更改或拉取请求。

无须手动编写伪版本。许多命令接受提交哈希或分支名，并自动转换为伪版本；如果存在对应标签，也可能转换为标签版本。例如：

```
go get example.com/mod@master
go list -m -json example.com/mod@abcd1234
```

### 主版本后缀 {#major-version-suffixes}

从主版本 2 开始，模块路径必须带有与主版本号匹配的<dfn>主版本后缀</dfn>，例如 `/v2`。如果模块在 `v1.0.0` 时的路径为 `example.com/mod`，那么在 `v2.0.0` 时就必须使用 `example.com/mod/v2`。

主版本后缀实现了[<dfn>导入兼容性规则</dfn>](https://research.swtch.com/vgo-import)：

> 如果旧包和新包具有相同的导入路径，新包必须向后兼容旧包。

按照定义，模块的新主版本中的包，与前一个主版本中的对应包不向后兼容。因此，从 `v2` 开始，包需要新的导入路径。这通过在模块路径后添加主版本后缀实现。由于模块路径是模块内每个包导入路径的前缀，添加主版本后缀就为每个不兼容版本提供了不同的导入路径。

主版本 `v0` 和 `v1` 不允许带主版本后缀。从 `v0` 升级到 `v1` 时无须更改模块路径，因为 `v0` 本来就是不稳定版本，不保证兼容。此外，大多数模块的 `v1` 都向后兼容最后一个 `v0` 版本；发布 `v1` 表示开始承诺兼容性，而不是说明它相对 `v0` 有不兼容更改。

一个特例是，以 `gopkg.in/` 开头的模块路径始终必须带主版本后缀，包括 `v0` 和 `v1`。该后缀以点号而不是斜杠开头，例如 `gopkg.in/yaml.v2`。

主版本后缀允许同一模块的多个主版本共存于一次构建中，这有时是解决[菱形依赖问题](https://research.swtch.com/vgo-import#dependency_story)所必需的。通常，传递依赖要求同一模块的两个不同版本时，会使用较高版本。但如果二者不兼容，任何一个都无法满足所有使用方。不兼容版本必须具有不同的主版本号，也就必须因主版本后缀而具有不同的模块路径。这样就消除了冲突：后缀不同的模块视为独立模块，其中的包也彼此不同，即使它们位于各自模块根目录下相同的相对子目录中。

许多 Go 项目在迁移到模块之前，甚至在模块机制出现之前，就已经发布了不带主版本后缀的 `v2` 或更高版本。这些版本使用 `+incompatible` 构建元数据标记，例如 `v2.0.0+incompatible`。更多信息请参阅[兼容非模块仓库](#non-module-compat)。

### 将包解析到模块 {#resolve-pkg-mod}

`go` 命令通过[包路径](#glos-package-path)加载包时，需要确定由哪个模块提供该包。

`go` 命令首先在[构建列表](#glos-build-list)中查找路径为该包路径前缀的模块。例如，导入 `example.com/a/b`，且构建列表包含 `example.com/a` 时，`go` 命令会检查 `example.com/a` 是否在 `b` 目录中提供该包。目录中至少需要有一个 `.go` 文件才会被视为包，此时不应用[构建约束](/pkg/go/build/#hdr-Build_Constraints)。如果构建列表中恰有一个模块提供该包，就使用它；如果没有，或有两个以上模块提供该包，`go` 命令就会报错。`-mod=mod` 选项会让 `go` 命令尝试查找提供缺失包的新模块，并更新 `go.mod` 和 `go.sum`。[`go get`](#go-get) 和 [`go
mod tidy`](#go-mod-tidy) 会自动执行这一过程。

<!-- NOTE(golang.org/issue/27899): the go command reports an error when two
or more modules provide a package with the same path as above. In the future,
we may try to upgrade one (or all) of the colliding modules.
-->

`go` 命令为某个包路径查找新模块时，会读取 `GOPROXY` 环境变量。它是由代理 URL 或关键字 `direct`、`off` 组成的逗号分隔列表。代理 URL 表示 `go` 命令应通过 [`GOPROXY` 协议](#goproxy-protocol)访问[模块代理](#glos-module-proxy)；`direct` 表示 `go` 命令应直接[与版本控制系统通信](#vcs)；`off` 表示不尝试通信。也可以通过 `GOPRIVATE` 和 `GONOPROXY` [环境变量](#environment-variables)控制这一行为。

对于 `GOPROXY` 列表中的每一项，`go` 命令会请求所有可能提供该包的模块路径，即包路径的各级前缀的最新版本。对于每个请求成功的模块路径，`go` 命令会下载其最新版本，并检查其中是否包含目标包。如果有一个或多个模块包含该包，选择路径最长的模块。如果找到了模块，但没有任何模块包含该包，则报错。如果没有找到任何模块，`go` 命令会尝试 `GOPROXY` 中的下一项；没有更多条目时就报错。

例如，假设 `go` 命令要查找提供 `golang.org/x/net/html` 包的模块，且 `GOPROXY` 设置为 `https://corp.example.com,https://proxy.golang.org`，则 `go` 命令可能发起以下请求：

* 向 `https://corp.example.com/` 并行请求：
  * `golang.org/x/net/html` 的最新版本
  * `golang.org/x/net` 的最新版本
  * `golang.org/x` 的最新版本
  * `golang.org` 的最新版本
* 如果向 `https://corp.example.com/` 发出的请求全部返回 404 或 410，则向 `https://proxy.golang.org/` 请求：
  * `golang.org/x/net/html` 的最新版本
  * `golang.org/x/net` 的最新版本
  * `golang.org/x` 的最新版本
  * `golang.org` 的最新版本

找到合适模块后，`go` 命令会把新模块的路径和版本作为一条新的[依赖要求](#go-mod-file-require)，加入主模块的 `go.mod` 文件。这样，以后加载同一个包时，就会使用相同模块的相同版本。如果解析到的包没有被主模块中的包导入，新增的依赖要求会带有 `// indirect` 注释。

## `go.mod` 文件 {#go-mod-file}

模块通过根目录下名为 `go.mod` 的 UTF-8 文本文件定义。`go.mod` 按行组织，每行包含一条指令，由关键字及其参数组成。例如：

```
module example.com/my/thing

go 1.23.0

require example.com/other/thing v1.0.2
require example.com/new/thing/v2 v2.3.4
exclude example.com/old/thing v1.2.3
replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
retract [v1.9.0, v1.9.5]
```

可以把相邻多行开头相同的关键字提取出来，组成一个块，形式类似于 Go 的分组导入。

```
require (
    example.com/new/thing/v2 v2.3.4
    example.com/old/thing v1.2.3
)
```

`go.mod` 的设计兼顾人工阅读和机器写入。`go` 命令提供了多个修改 `go.mod` 文件的子命令，例如 [`go get`](#go-get) 可以升级或降级特定依赖。加载模块图的命令会在必要时[自动更新](#go-mod-file-updates) `go.mod`。[`go mod
edit`](#go-mod-edit) 支持更底层的编辑操作。Go 程序也可以使用 [`golang.org/x/mod/modfile`](https://pkg.go.dev/golang.org/x/mod/modfile?tab=doc) 包，以编程方式进行同样的修改。

[主模块](#glos-main-module)以及通过本地文件路径指定的[替换模块](#go-mod-file-replace)，都必须具有 `go.mod` 文件。不过，没有显式 `go.mod` 的模块仍然可以作为[依赖](#go-mod-file-require)，或以模块路径和版本的形式用作替换项；详见[兼容非模块仓库](#non-module-compat)。

### 词法元素 {#go-mod-file-lexical}

解析 `go.mod` 文件时，内容会被拆分为一系列词法单元，包括空白、注释、标点、关键字、标识符和字符串。

*空白*包括空格（U+0020）、制表符（U+0009）、回车（U+000D）和换行（U+000A）。除换行外，空白字符只用于分隔原本会连在一起的词法单元，不产生其他作用。换行本身是有意义的词法单元。

*注释*以 `//` 开头，延续到行末。不允许使用 `/* */` 注释。

*标点*词法单元包括 `(`、`)` 和 `=>`。

*关键字*用于区分 `go.mod` 文件中的不同指令。允许的关键字包括 `module`、`go`、`require`、`replace`、`exclude` 和 `retract`。

*标识符*是不包含空白的字符序列，例如模块路径或语义化版本。

*字符串*是由引号包围的字符序列，有两种形式：以双引号（`"`，U+0022）开头和结尾的解释型字符串，以及以反引号（<code>&#x60;</code>，U+0060）开头和结尾的原始字符串。解释型字符串可以包含转义序列，由反斜杠（`\`，U+005C）及其后的另一个字符组成。转义后的双引号（`\"`）不会结束字符串。解释型字符串去除引号后的值，是将引号间每个转义序列替换为反斜杠后面的字符后得到的序列，例如 `\"` 替换为 `"`，`\n` 替换为 `n`。原始字符串去除引号后的值则直接是两个反引号之间的字符序列，其中反斜杠没有特殊含义。

在 `go.mod` 语法中，标识符与字符串可以互换使用。

### 模块路径与版本 {#go-mod-file-ident}

`go.mod` 文件中的大多数标识符和字符串，表示模块路径或版本。

模块路径必须满足以下要求：

* 路径由一个或多个以斜杠（`/`，U+002F）分隔的路径段组成，不能以斜杠开头或结尾。
* 每个路径段都是非空字符串，只能包含 ASCII 字母、ASCII 数字以及有限的 ASCII 标点：`-`、`.`、`_` 和 `~`。
* 路径段不能以点号（`.`，U+002E）开头或结尾。
* 路径段在第一个点号之前的部分，不能是不区分大小写的 Windows 保留文件名，例如 `CON`、`com1`、`NuL` 等。
* 路径段在第一个点号之前的部分，不能以波浪号加一个或多个数字结尾，例如 `EXAMPL~1.COM`。

如果模块路径出现在 `require` 指令中且没有被替换，或者出现在 `replace` 指令的右侧，`go` 命令可能需要下载该路径对应的模块，因此还必须满足一些额外要求。

* 第一个路径段（第一个斜杠之前的部分；如果没有斜杠，则为整个路径）按照惯例是域名，只能包含小写 ASCII 字母、ASCII 数字、点号（`.`，U+002E）和连字符（`-`，U+002D）；必须至少包含一个点号，且不能以连字符开头。
* 如果最后一个路径段形如 `/vN`，且 `N` 看起来像数值（由 ASCII 数字和点号组成），那么 `N` 不能有前导零，整个段不能是 `/v1`，并且不能包含点号。
  * 对于以 `gopkg.in/` 开头的路径，此要求替换为遵循 [gopkg.in](https://gopkg.in) 服务自身的约定。

`go.mod` 文件中的版本可以是[规范版本](#glos-canonical-version)，也可以是非规范版本。

规范版本以字母 `v` 开头，后面跟符合[语义化版本 2.0.0](https://semver.org/spec/v2.0.0.html)规范的版本号。更多信息请参阅[版本](#versions)。

大多数其他标识符和字符串都可以用作非规范版本，但为了避免文件系统、仓库和[模块代理](#glos-module-proxy)方面的问题，仍有一些限制。非规范版本仅允许出现在主模块的 `go.mod` 中。`go` 命令自动[更新](#go-mod-file-updates) `go.mod` 文件时，会尝试将每个非规范版本替换为等价的规范版本。

模块路径与版本关联出现时，例如在 `require`、`replace` 和 `exclude` 指令中，路径最后一段必须与版本一致。请参阅[主版本后缀](#major-version-suffixes)。

### 语法 {#go-mod-file-grammar}

下面使用扩展巴科斯范式（EBNF）描述 `go.mod` 语法。有关 EBNF 的说明，请参阅 [Go 语言规范中的记法一节](/ref/spec#Notation)。

```
GoMod = { Directive } .
Directive = ModuleDirective |
            GoDirective |
            ToolDirective |
            IgnoreDirective |
            RequireDirective |
            ExcludeDirective |
            ReplaceDirective |
            RetractDirective .
```

换行、标识符和字符串分别用 `newline`、`ident` 和 `string` 表示。

模块路径和版本分别用 `ModulePath` 和 `Version` 表示。

```
ModulePath = ident | string . /* see restrictions above */
Version = ident | string .    /* see restrictions above */
```

### `module` 指令 {#go-mod-file-module}

`module` 指令定义主模块的[路径](#glos-module-path)。`go.mod` 文件必须且只能包含一条 `module` 指令。

```
ModuleDirective = "module" ( ModulePath | "(" newline ModulePath newline ")" ) newline .
```

示例：

```
module golang.org/x/net
```

#### 弃用 {#go-mod-file-module-deprecation}

可以通过注释块将模块标记为已弃用：注释某个段落的开头必须包含区分大小写的字符串 `Deprecated:`。弃用说明从冒号之后开始，持续到该段落结束。注释可以紧邻 `module` 指令之前，也可以放在该指令之后的同一行。

示例：

```
// Deprecated: use example.com/mod/v2 instead.
module example.com/mod
```

从 Go 1.17 开始，[`go list -m -u`](#go-list-m) 会检查[构建列表](#glos-build-list)中所有模块的弃用信息。[`go get`](#go-get) 会检查构建命令行指定的包所需模块是否已弃用。

`go` 命令获取模块弃用信息时，会加载与 `@latest` [版本查询](#version-queries)匹配的版本中的 `go.mod`，此时不考虑[撤回](#go-mod-file-retract)或[排除](#go-mod-file-exclude)。`go` 命令也从这个 `go.mod` 文件加载[已撤回版本](#glos-retracted-version)列表。

要弃用模块，作者可以添加 `// Deprecated:` 注释，并为新版本打标签。作者也可以在更高版本中修改或移除弃用说明。

弃用声明适用于模块的所有次版本。在此语境下，`v2` 及更高主版本被视为独立模块，因为其[主版本后缀](#glos-major-version-suffix)赋予了不同的模块路径。

弃用说明用于告知用户该模块不再受支持，并提供迁移建议，例如迁移到最新主版本。不能单独弃用某个次版本或补丁版本；这种情况使用 [`retract`](#go-mod-file-retract) 可能更合适。

### `go` 指令 {#go-mod-file-go}

`go` 指令表示模块的代码以某个 Go 版本的语义为基础编写。其版本必须是有效的 [Go 版本](/doc/toolchain#version)，例如 `1.14`、`1.21rc1` 或 `1.23.0`。

`go` 指令设置使用该模块所需的最低 Go 版本。在 Go 1.21 之前，这条指令只是建议；现在它是强制要求：Go 工具链会拒绝使用声明了更高 Go 版本的模块。

`go` 指令也是选择要运行哪个 Go 工具链的依据之一。详细说明请参阅《[Go 工具链](/doc/toolchain)》。

`go` 指令影响新语言特性的使用：

* 对于模块内的包，编译器会拒绝使用晚于 `go` 指令指定版本引入的语言特性。例如，模块声明 `go 1.12` 时，包中不能使用 Go 1.13 引入的 `1_000_000` 这样的数值字面量。
* 较旧的 Go 版本构建模块内的包时，如果遇到编译错误，错误信息会指出模块面向较新的 Go 版本。例如，模块声明 `go 1.13`，并使用数值字面量 `1_000_000`，用 Go 1.12 构建时，编译器就会提示该代码面向 Go 1.13。

`go` 指令也会影响 `go` 命令的行为：

* 在 `go 1.14` 或更高版本下，可以自动启用 [vendor 模式](#vendoring)。如果存在 `vendor/modules.txt`，且内容与 `go.mod` 一致，就不必显式指定 `-mod=vendor`。
* 在 `go 1.16` 或更高版本下，`all` 包模式仅匹配[主模块](#glos-main-module)中的包及测试直接或间接导入的包。它与模块机制引入以来 [`go mod vendor`](#go-mod-vendor) 保留的包集合相同。在更低版本中，`all` 还包括主模块所导入包的测试、这些测试所导入包的测试，依此类推。
* 在 `go 1.17` 或更高版本下：
  * 对于提供主模块包或测试直接、间接导入的任意包的模块，`go.mod` 都会显式记录 [`require` 指令](#go-mod-file-require)。（在 `go
     1.16` 及更低版本中，只有不记录就会导致[最小版本选择](#minimal-version-selection)选出不同版本时，才会加入[间接依赖](#glos-direct-dependency)。）这些额外信息使[模块图剪枝](#graph-pruning)和[模块延迟加载](#lazy-loading)成为可能。
  * 由于 `// indirect` 依赖可能比之前的 `go` 版本多得多，间接依赖会记录在 `go.mod` 的单独分组中。
  * `go mod vendor` 不复制被 vendoring 的依赖中的 `go.mod` 和 `go.sum` 文件，使得在 `vendor` 子目录中运行 `go` 命令时能够识别正确的主模块。
  * `go mod vendor` 会将每个依赖的 `go.mod` 中声明的 `go` 版本记录到 `vendor/modules.txt`。
* 在 `go 1.21` 或更高版本下：
  * `go` 行声明使用模块所必需的最低 Go 版本。
  * `go` 行的版本必须不低于所有依赖的 `go` 行版本。
  * `go` 命令不再尝试保持与前一个较旧 Go 版本的兼容性。
  * `go` 命令会更谨慎地在 `go.sum` 中保留 `go.mod` 文件的校验和。
<!-- If you update this list, also update /doc/modules/gomod-ref#go-notes. -->

`go.mod` 文件最多只能包含一条 `go` 指令。如果没有，大多数命令会添加使用当前 Go 版本的 `go` 指令。

如果没有 `go` 指令，则按 `go 1.16` 处理。

```
GoDirective = "go" GoVersion newline .
GoVersion = string | ident .  /* valid release version; see above */
```

示例：

```
go 1.23.0
```

### `toolchain` 指令 {#go-mod-file-toolchain}

`toolchain` 指令声明建议与模块配合使用的 Go 工具链。建议的工具链版本不能低于 `go` 指令要求的版本。`toolchain` 指令仅在该模块是主模块，且默认工具链版本低于建议版本时生效。

为保证可复现性，`go` 命令每次更新 `go.mod` 中的 `go` 版本时，通常是在执行 `go get` 期间，都会在 `toolchain` 行写入自身的工具链名称。

详细说明请参阅《[Go 工具链](/doc/toolchain)》。

```
ToolchainDirective = "toolchain" ToolchainName newline .
ToolchainName = string | ident .  /* valid toolchain name; see “Go toolchains” */
```

示例：

```
toolchain go1.21.0
```

### `godebug` 指令 {#go-mod-file-godebug}

`godebug` 指令声明一项 [GODEBUG 设置](/doc/godebug)，在该模块作为主模块时应用。可以存在多条此类指令，也可以使用分组形式。主模块指定不存在的 GODEBUG 键会报错。`godebug key=value` 的效果，等同于每个正在编译的 main 包都包含一个写有 `//go:debug key=value` 的源文件。

```
GodebugDirective = "godebug" ( GodebugSpec | "(" newline { GodebugSpec } ")" newline ) .
GodebugSpec = GodebugKey "=" GodebugValue newline.
GodebugKey = GodebugChar { GodebugChar }.
GodebugValue = GodebugChar { GodebugChar }.
GodebugChar = any non-space character except , " ` ' (comma and quotes).
```

示例：

```
godebug default=go1.21
godebug (
	panicnil=1
	asynctimerchan=0
)
```

### `require` 指令 {#go-mod-file-require}

`require` 指令声明某个依赖模块所需的最低版本。对于每个被要求的模块版本，`go` 命令会加载该版本的 `go.mod`，并纳入其中的依赖要求。加载全部要求后，`go` 命令通过[最小版本选择（MVS）](#minimal-version-selection)解析它们，生成[构建列表](#glos-build-list)。

`go` 命令会为部分依赖要求自动添加 `// indirect` 注释。`// indirect` 表示[主模块](#glos-main-module)中的任何包都没有直接导入该依赖模块中的包。

如果 [`go` 指令](#go-mod-file-go)指定 `go 1.16` 或更低版本，当选中的模块版本高于主模块其他依赖已经间接要求的版本时，`go` 命令会添加间接依赖要求。这可能是因为显式升级（`go get -u ./...`），移除了原先要求该版本的其他依赖（`go
mod tidy`），或者某个依赖导入了一个包，却没有在自己的 `go.mod` 中声明相应要求，例如该依赖根本没有 `go.mod` 文件。

在 `go 1.17` 及更高版本下，`go` 命令会为每个提供相关包的模块添加间接依赖要求：相关包包括主模块中的包或测试直接或[间接](#glos-indirect-dependency)导入的包，以及作为参数传给 `go get` 的包。这些更完整的要求支持[模块图剪枝](#graph-pruning)和[模块延迟加载](#lazy-loading)。

```
RequireDirective = "require" ( RequireSpec | "(" newline { RequireSpec } ")" newline ) .
RequireSpec = ModulePath Version newline .
```

示例：

```
require golang.org/x/net v1.2.3

require (
    golang.org/x/crypto v1.4.5 // indirect
    golang.org/x/text v1.6.7
)
```

### `tool` 指令 {#go-mod-file-tool}

从 Go 1.24 开始，`tool` 指令可以将一个包加入当前模块的依赖。当当前工作目录位于该模块内，或位于包含该模块的工作区内时，还可以通过 `go tool` 运行该工具。

如果工具包不在当前模块中，必须通过 `require` 指令指定要使用的工具版本。

`tool` 元模式解析为当前模块 `go.mod` 中定义的工具列表；在工作区模式下，则解析为工作区所有模块中定义的工具的并集。

```
ToolDirective = "tool" ( ToolSpec | "(" newline { ToolSpec } ")" newline ) .
ToolSpec = ModulePath newline .
```

示例：

```
tool golang.org/x/tools/cmd/stringer

tool (
    example.com/module/cmd/a
    example.com/module/cmd/b
)
```

### `ignore` 指令 {#go-mod-file-ignore}

`ignore` 指令使 go 命令在匹配包模式时，忽略指定的、以斜杠分隔的目录路径，以及这些目录递归包含的所有文件和子目录。

如果路径以 `./` 开头，就按相对于模块根目录的路径解释。匹配包模式时，该目录及其递归包含的目录和文件都会被忽略。

否则，模块内任意深度与该路径匹配的目录，以及这些目录递归包含的所有目录和文件，都会被忽略。

```
IgnoreDirective = "ignore" ( IgnoreSpec | "(" newline { IgnoreSpec } ")" newline ) .
IgnoreSpec = RelativeFilePath newline .
RelativeFilePath = /* slash-separated relative file path */ .
```

示例：
```
ignore ./node_modules

ignore (
    static
    content/html
    ./third_party/javascript
)
```

### `exclude` 指令 {#go-mod-file-exclude}

`exclude` 指令禁止 `go` 命令加载某个模块版本。

从 Go 1.16 开始，如果任意 `go.mod` 中的 `require` 指令引用了某个版本，而主模块 `go.mod` 的 `exclude` 指令排除了该版本，那么这条依赖要求会被忽略。这可能使 [`go get`](#go-get) 和 [`go mod tidy`](#go-mod-tidy) 等命令在 `go.mod` 中添加更高版本的新要求，并在适当情况下加上 `// indirect` 注释。

在 Go 1.16 之前，如果 `require` 引用了被排除的版本，`go` 命令会列出模块的可用版本，如 [`go
list -m -versions`](#go-list-m) 所显示的列表，并改为加载下一个更高且未被排除的版本。由于下一个更高版本可能随时间变化，这会造成版本选择不确定。此过程考虑正式版本和预发布版本，但不考虑伪版本。如果没有更高版本，`go` 命令就会报错。

`exclude` 指令只在主模块的 `go.mod` 中生效，其他模块中的该指令会被忽略。详见[最小版本选择](#minimal-version-selection)。

```
ExcludeDirective = "exclude" ( ExcludeSpec | "(" newline { ExcludeSpec } ")" newline ) .
ExcludeSpec = ModulePath Version newline .
```

示例：

```
exclude golang.org/x/net v1.2.3

exclude (
    golang.org/x/crypto v1.4.5
    golang.org/x/text v1.6.7
)
```

### `replace` 指令 {#go-mod-file-replace}

`replace` 指令将模块的某个特定版本或全部版本的内容，替换为其他位置的内容。替换项可以通过另一个模块路径及版本指定，也可以通过当前平台的文件路径指定。

如果箭头（`=>`）左侧指定了版本，则只替换该版本，其他版本仍正常访问。如果左侧省略版本，则替换该模块的所有版本。

如果箭头右侧是绝对路径或相对路径（以 `./` 或 `../` 开头），就将其视为替换模块根目录的本地文件路径，该目录必须包含 `go.mod`。这种情况下必须省略右侧版本。

如果右侧不是本地路径，就必须是有效模块路径，并且必须指定版本。同一个模块版本不能同时出现在构建列表的其他位置。

无论替换项使用本地路径还是模块路径，只要替换模块包含 `go.mod`，其中的 `module` 指令就必须与被替换的模块路径匹配。

`replace` 指令只在主模块的 `go.mod` 中生效，其他模块中的该指令会被忽略。详见[最小版本选择](#minimal-version-selection)。

存在多个主模块时，它们各自的 `go.mod` 都会生效。不同主模块之间不允许出现冲突的 `replace` 指令，必须移除冲突，或通过 [`go.work file` 中的 replace](#go-work-file-replace) 覆盖它们。

注意，单独一条 `replace` 指令不会把模块加入[模块图](#glos-module-graph)。主模块或依赖的 `go.mod` 中还必须有引用被替换模块版本的 [`require` 指令](#go-mod-file-require)，无论它位于主模块的 `go.mod` 还是依赖文件中。如果左侧模块版本没有被要求，`replace` 就不会产生作用。

```
ReplaceDirective = "replace" ( ReplaceSpec | "(" newline { ReplaceSpec } ")" newline ) .
ReplaceSpec = ModulePath [ Version ] "=>" FilePath newline
            | ModulePath [ Version ] "=>" ModulePath Version newline .
FilePath = /* platform-specific relative or absolute file path */
```

示例：

```
replace golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5

replace (
    golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5
    golang.org/x/net => example.com/fork/net v1.4.5
    golang.org/x/net v1.2.3 => ./fork/net
    golang.org/x/net => ./fork/net
)
```

### `retract` 指令 {#go-mod-file-retract}

`retract` 指令表示不应再依赖当前 `go.mod` 定义的模块的某个版本或版本范围。过早发布版本，或发布后发现严重问题时，`retract` 很有用。已撤回版本仍应保留在版本控制仓库和[模块代理](#glos-module-proxy)中，确保依赖它们的构建不会被破坏。*retract* 一词借自学术出版：被撤回的论文仍然可以查阅，但存在问题，不应作为后续工作的基础。

模块版本被撤回后，用户不会通过 [`go get`](#go-get)、[`go mod tidy`](#go-mod-tidy) 等命令自动升级到该版本。依赖已撤回版本的构建仍应正常工作，但用户通过 [`go list
-m -u`](#go-list-m) 检查更新，或通过 [`go get`](#go-get) 更新相关模块时，会收到撤回通知。

要撤回版本，模块作者应在 `go.mod` 中添加 `retract`，然后发布包含该指令的新版本。新版本必须高于其他正式版本或预发布版本，也就是说，在尚未考虑撤回信息时，`@latest` [版本查询](#version-queries)应解析到这个新版本。`go` 命令会从 `go list -m -retracted $modpath@latest` 显示的版本中加载并应用撤回信息，其中 `$modpath` 是模块路径。

[`go list -m
-versions`](#go-list-m) 默认不显示已撤回版本，除非指定 `-retracted`。解析 `@>=v1.2.3` 或 `@latest` 等版本查询时，也会排除已撤回版本。

包含撤回指令的版本可以撤回自身。如果模块最高的正式版本或预发布版本撤回了自己，那么排除已撤回版本后，`@latest` 会解析到更低版本。

例如，模块 `example.com/m` 的作者意外发布了 `v1.0.0`。为防止用户升级到 `v1.0.0`，作者可以在 `go.mod` 中添加两条 `retract` 指令，然后发布包含这些撤回信息的 `v1.0.1` 标签。

```
retract (
    v1.0.0 // Published accidentally.
    v1.0.1 // Contains retractions only.
)
```

用户运行 `go get example.com/m@latest` 时，`go` 命令会从当前最高版本 `v1.0.1` 读取撤回信息。由于 `v1.0.0` 和 `v1.0.1` 都被撤回，`go` 命令会升级，甚至降级，到下一个最高版本，例如 `v0.9.5`。

`retract` 可以指定单个版本，例如 `v1.0.0`；也可以用 `[` 和 `]` 包围带上下界的闭区间，例如 `[v1.1.0, v1.2.0]`。单个版本等价于上下界相同的区间。与其他指令一样，多条 `retract` 可以分组：以行末的 `(` 开始，以单独一行的 `)` 结束。

每条 `retract` 都应附上解释撤回原因的注释，尽管这不是强制要求。`go` 命令可能在撤回警告和 `go list` 输出中显示原因。原因注释可以紧邻 `retract` 之前且不留空行，也可以写在指令之后的同一行。如果注释位于分组之前，就适用于组内所有没有自身注释的 `retract`。原因注释可以跨多行。

```
RetractDirective = "retract" ( RetractSpec | "(" newline { RetractSpec } ")" newline ) .
RetractSpec = ( Version | "[" Version "," Version "]" ) newline .
```

示例：

* 撤回 `v1.0.0` 到 `v1.9.9` 之间的全部版本：

```
retract v1.0.0
retract [v1.0.0, v1.9.9]
retract (
    v1.0.0
    [v1.0.0, v1.9.9]
)
```

* 过早发布 `v1.0.0` 后，恢复到尚未发布正式版本的状态：

```
retract [v0.0.0, v1.0.1] // assuming v1.0.1 contains this retraction.
```

* 撤回整个模块，包括全部伪版本和标签版本：

```
retract [v0.0.0-0, v0.15.2]  // assuming v0.15.2 contains this retraction.
```

`retract` 在 Go 1.16 中引入。Go 1.15 及更低版本遇到[主模块](#glos-main-module) `go.mod` 中的 `retract` 会报错，而依赖模块 `go.mod` 中的 `retract` 则会被忽略。

### 自动更新 {#go-mod-file-updates}

如果 `go.mod` 缺少信息，或未准确反映实际情况，大多数命令都会报错。可以使用 [`go get`](#go-get) 和 [`go mod tidy`](#go-mod-tidy) 修复其中的大多数问题。此外，大多数支持模块的命令，例如 `go build`、`go test`，都可以使用 `-mod=mod` 选项，让 `go` 命令自动修复 `go.mod` 和 `go.sum` 中的问题。

例如，考虑下面这个 `go.mod`：

```
module example.com/M

go 1.23.0

require (
    example.com/A v1
    example.com/B v1.0.0
    example.com/C v1.0.0
    example.com/D v1.2.3
    example.com/E dev
)

exclude example.com/D v1.2.3
```

由 `-mod=mod` 触发的更新，会把非规范版本标识符改写为[规范](#glos-canonical-version)语义化版本。因此，`example.com/A` 的 `v1` 会变成 `v1.0.0`，`example.com/E` 的 `dev` 则会变成 `dev` 分支最新提交的伪版本，例如 `v0.0.0-20180523231146-b3f5c0f6e5f1`。

更新会调整依赖要求以遵守排除规则。例如，被排除的 `example.com/D v1.2.3` 会改为 `example.com/D` 的下一个可用版本，可能是 `v1.2.4` 或 `v1.3.0`。

更新会删除多余或具有误导性的依赖要求。例如，`example.com/A v1.0.0` 自身要求 `example.com/B v1.2.0` 和 `example.com/C
v1.0.0`，那么 `go.mod` 中的 `example.com/B v1.0.0` 会产生误导，因为 `example.com/A` 已要求 `v1.2.0`；而 `example.com/C v1.0.0` 是多余的，因为 `example.com/A` 已经间接要求同一版本。因此，两条要求都会删除。如果主模块中的包直接导入 `example.com/B` 或 `example.com/C` 中的包，则保留对应要求，但更新为实际使用的版本。

最后，更新会按照规范格式重新格式化 `go.mod`，使后续自动修改产生尽量小的差异。如果只需要调整格式，`go` 命令不会更新 `go.mod`。

模块图决定导入语句的含义，因此所有加载包的命令都会使用 `go.mod`，也就可能更新它，包括 `go build`、`go get`、`go install`、`go list`、`go test` 和 `go mod tidy`。

Go 1.15 及更低版本默认启用 `-mod=mod`，因此会自动更新。从 Go 1.16 开始，`go` 命令默认按 `-mod=readonly` 处理：如果需要修改 `go.mod`，`go` 命令会报错并建议修复方法。

## 最小版本选择（MVS） {#minimal-version-selection}

Go 使用称为<dfn>最小版本选择（MVS）</dfn>的算法，选择构建包时使用的一组模块版本。Russ Cox 的[最小版本选择](https://research.swtch.com/vgo-mvs)一文详细介绍了该算法。

从概念上说，MVS 操作的是由 [`go.mod` 文件](#glos-go-mod-file)定义的模块有向图。每个顶点表示一个模块版本，每条边表示通过 [`require`](#go-mod-file-require) 指定的依赖最低版本要求。主模块的 `go.mod` 中的 [`exclude`](#go-mod-file-exclude) 和 [`replace`](#go-mod-file-replace)，以及 `go.work` 中的 [`replace`](#go-work-file-replace)，都可以改变这张图。

MVS 的输出是[构建列表](#glos-build-list)，即构建所使用的模块版本列表。

MVS 从主模块开始遍历模块图；主模块是图中没有版本号的特殊顶点。遍历时记录每个模块被要求的最高版本。遍历结束后，这些最高要求版本构成构建列表：它们是满足所有要求的最低版本。

可以使用 [`go list -m
all`](#go-list-m) 查看构建列表。与其他依赖管理系统不同，构建列表不保存在“锁定”文件中。MVS 的结果是确定的，不会因为依赖发布新版本而改变，因此每个支持模块的命令都会在开始时重新计算构建列表。

考虑下图：主模块要求 A 至少为 1.2、B 至少为 1.2。A 1.2 和 B 1.2 分别要求 C 1.3 和 C 1.4，而 C 1.3 和 C 1.4 都要求 D 1.2。

![已访问版本高亮显示的模块版本图](/doc/mvs/buildlist.svg "MVS 构建列表图")

MVS 会访问并加载蓝色高亮模块版本的 `go.mod`。遍历结束后，返回包含加粗版本的构建列表：A 1.2、B 1.2、C 1.4 和 D 1.2。注意，虽然 B 和 D 存在更高版本，但没有依赖要求它们，因此 MVS 不会选择。

### 替换 {#mvs-replace}

可以通过主模块 `go.mod` 或工作区 `go.work` 中的 [`replace` 指令](#go-mod-file-replace)，替换模块内容，包括其 `go.mod`。`replace` 可以作用于模块的某个版本，也可以作用于全部版本。

替换会改变模块图，因为替换模块的依赖可能与原版本不同。

考虑下例：C 1.4 被替换为 R，R 依赖 D 1.3，而不是 D 1.2。因此 MVS 返回的构建列表包含 A 1.2、B 1.2、C 1.4（由 R 替换）和 D 1.3。

![包含替换项的模块版本图](/doc/mvs/replace.svg "MVS 替换")

### 排除 {#mvs-exclude}

也可以通过主模块 `go.mod` 中的 [`exclude` 指令](#go-mod-file-exclude)，排除模块的特定版本。

排除同样会改变模块图。版本被排除后，会从模块图中移除；对它的依赖要求则重新指向下一个更高版本。

如下图，C 1.3 被排除，MVS 会按照 A 1.2 要求 C 1.4（下一个更高版本），而不是 C 1.3 来处理。

![包含排除项的模块版本图](/doc/mvs/exclude.svg "MVS 排除")

### 升级 {#mvs-upgrade}

[`go get`](#go-get) 可以升级一组模块。执行升级时，`go` 命令会先修改模块图，为已访问版本添加指向升级后版本的边，然后运行 MVS。

如下图，模块 B 可以从 1.2 升级到 1.3，C 从 1.3 升级到 1.4，D 从 1.2 升级到 1.3。

![执行升级后的模块版本图](/doc/mvs/upgrade.svg "MVS 升级")

升级和降级都可能添加或移除间接依赖。在本例中，升级后构建列表新增了 E 1.1 和 F 1.1，因为 B 1.3 要求 E 1.1。

为保留升级结果，`go` 命令会更新 `go.mod` 中的依赖要求，将 B 的要求改为 1.3，并添加带 `// indirect` 注释的 C 1.4 和 D 1.3 要求，因为不添加这些要求就不会选中它们。

### 降级 {#mvs-downgrade}

[`go get`](#go-get) 也可以降级一组模块。执行降级时，`go` 命令从模块图中移除高于目标版本的版本，同时移除依赖这些被移除版本的其他模块版本，因为它们可能与降级后的依赖不兼容。如果主模块要求的版本因降级而被移除，就将要求改为一个尚未被移除的较早版本；如果不存在这样的版本，就删除这条要求。

如下图，假设 C 1.4 存在问题，需要降级到 C 1.3。C 1.4 会从模块图中移除。B 1.2 因为要求 C 1.4 或更高版本，也会被移除。主模块对 B 的要求改为 1.1。

![执行降级后的模块版本图](/doc/mvs/downgrade.svg "MVS 降级")

在参数后添加 `@none`，还可以用 [`go get`](#go-get) 完全移除依赖。它的工作方式类似降级，会从模块图中移除指定模块的所有版本。

## 模块图剪枝 {#graph-pruning}

如果主模块声明 `go 1.17` 或更高版本，那么用于[最小版本选择](#minimal-version-selection)的[模块图](#glos-module-graph)，对于自身 `go.mod` 同样声明 `go 1.17` 或更高版本的每个依赖模块，只包含其*直接*要求；但如果该模块版本还被某个声明 `go 1.16` 或更低版本的*其他*依赖直接或间接要求，则不适用这一规则。也就是说，`go 1.17` 依赖的*传递依赖*会从模块图中*剪除*。

由于 `go 1.17` 的 `go.mod` 为构建模块中任意包或测试所需的每个依赖，都记录了 [require 指令](#go-mod-file-require)，剪枝后的模块图仍包含执行 `go build` 或 `go test` 所需的全部依赖，适用于[主模块](#glos-main-module)显式要求的任意依赖模块中的包。如果某个模块并非构建特定模块内任意包或测试所必需，它就不会影响这些包的运行时行为。因此，被剪除的依赖只会在本来不相关的模块之间引入干扰。

即使某个模块的依赖要求被剪除了，模块本身仍保留在模块图中，也仍会出现在 `go list -m all` 中：它的[选定版本](#glos-selected-version)是明确已知的，也可以从中加载包，例如作为其他模块测试的传递依赖。不过，`go` 命令难以判断这些模块的哪些依赖已被满足，因此 `go
build` 和 `go test` 的参数不能直接指定依赖要求已被剪除的模块中的包。[`go get`](#go-get) 会将每个指定包所属的模块提升为显式依赖，从而允许对该包执行 `go build` 或 `go test`。

Go 1.16 及更早版本不支持模块图剪枝，因此对于每个声明 `go 1.16` 或更低版本的模块，仍会保留其依赖的完整传递闭包，包括间接涉及的 `go 1.17` 依赖。在 `go 1.16` 及更低版本下，`go.mod` 只记录[直接依赖](#glos-direct-dependency)，所以必须加载更大的模块图，才能保证包含所有间接依赖。

[`go mod tidy`](#go-mod-tidy) 为模块生成的 [`go.sum`](#go-sum-files)，默认包含比 [`go` 指令](#go-mod-file-go)指定版本*低一个版本*的 Go 所需的校验和。因此，`go 1.17` 模块包含 Go 1.16 加载完整模块图所需的校验和，而 `go 1.18` 模块只包含 Go 1.17 加载剪枝后模块图所需的校验和。可以用 `-compat` 覆盖这一默认版本，例如更积极地精简 `go 1.17` 模块的 `go.sum`。

更多细节请参阅[设计文档](https://go.googlesource.com/proposal/+/master/design/36460-lazy-module-loading.md)。

### 模块延迟加载 {#lazy-loading}

模块图剪枝所需的更完整依赖要求，还支持模块内工作的另一项优化。如果主模块声明 `go 1.17` 或更高版本，`go` 命令会尽量推迟加载完整模块图，直到确实需要时才加载。它先只读取主模块的 `go.mod`，尝试仅根据这些要求加载待构建的包。如果某个待导入的包不在这些要求中，例如主模块之外某个包的测试依赖，才会按需加载其余模块图。

如果不加载模块图就能找到所有导入包，`go` 命令接下来只加载这些包所属模块的 `go.mod`，并将其要求与主模块要求比较，确保局部一致。版本控制合并、手工编辑以及通过本地路径[替换](#go-mod-file-replace)的模块发生变化，都可能造成不一致。

## 工作区 {#workspaces}

<dfn>工作区</dfn>是磁盘上的一组模块，在运行[最小版本选择（MVS）](#minimal-version-selection)时，它们共同作为主模块。

可以通过 [`go.work` 文件](#go-work-file)声明工作区，其中指定每个成员模块目录的相对路径。如果不存在 `go.work`，工作区就只包含当前目录所属的单个模块。

大多数处理模块的 `go` 子命令，都作用于当前工作区确定的模块集合。`go mod init`、`go mod why`、`go mod edit`、`go mod tidy`、`go mod vendor` 和 `go get` 始终只作用于单个主模块。

命令首先检查 `GOWORK` 环境变量，以判断是否处于工作区上下文。若 `GOWORK` 为 `off`，则使用单模块上下文。若为空或未设置，就依次在当前工作目录和各级父目录中查找 `go.work`。找到时使用其中定义的工作区，否则只使用工作目录所属的模块。如果 `GOWORK` 指向一个存在且以 .work 结尾的文件，也会启用工作区模式；其他值会报错。可以通过 `go env GOWORK` 查看 `go` 命令使用的 `go.work` 文件；如果 `go` 命令不处于工作区模式，`go env GOWORK` 的输出为空。

### `go.work` 文件 {#go-work-file}

工作区由名为 `go.work` 的 UTF-8 文本文件定义。`go.work` 按行组织，每行包含一条由关键字及参数组成的指令。例如：

```
go 1.23.0

use ./my/first/thing
use ./my/second/thing

replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
```

与 `go.mod` 一样，可以把相邻多行共同的开头关键字提取出来，组成一个分组。

```
use (
    ./my/first/thing
    ./my/second/thing
)
```

`go` 命令提供了多个操作 `go.work` 的子命令。[`go work init`](#go-work-init) 创建新的 `go.work`；[`go work use`](#go-work-use) 向 `go.work` 添加模块目录；[`go work edit`](#go-work-edit) 提供底层编辑操作。Go 程序也可以使用 [`golang.org/x/mod/modfile`](https://pkg.go.dev/golang.org/x/mod/modfile?tab=doc) 包，以编程方式进行同样的修改。

go 命令会维护一个 `go.work.sum` 文件，记录工作区所需、但工作区各模块的 go.sum 文件合起来仍未包含的哈希值。

通常不建议将 go.work 提交到版本控制系统，原因有两个：

* 提交的 `go.work` 可能覆盖开发者放在父目录中的个人 `go.work`，使其 `use` 指令失效，造成困惑。
* 提交的 `go.work` 可能让持续集成（CI）系统选择并测试错误的依赖版本。通常应让 CI 不使用 `go.work`，以便测试该模块被其他模块依赖时的实际行为；在后一种情况下，模块内部的 `go.work` 不生效。

不过，某些场景下提交 `go.work` 是合理的。例如，仓库中的模块始终一起开发，不会与外部模块联合开发，开发者可能不需要在工作区中采用其他模块组合。此时，模块作者应确保各个模块得到正确测试和发布。

### 词法元素 {#go-work-file-lexical}

`go.work` 的词法元素与 [`go.mod files` 的定义](#go-mod-file-lexical)完全相同。

### 语法 {#go-work-file-grammar}

下面使用扩展巴科斯范式（EBNF）描述 `go.work` 语法。有关 EBNF 记法，请参阅 [Go 语言规范中的记法一节](/ref/spec#Notation)。

```
GoWork = { Directive } .
Directive = GoDirective |
            ToolchainDirective |
            UseDirective |
            ReplaceDirective .
```

换行、标识符和字符串分别用 `newline`、`ident` 和 `string` 表示。

模块路径和版本分别用 `ModulePath` 和 `Version` 表示，其定义与 [`go.mod files`](#go-mod-file-lexical) 完全相同。

```
ModulePath = ident | string . /* see restrictions above */
Version = ident | string .    /* see restrictions above */
```

### `go` 指令 {#go-work-file-go}

有效的 `go.work` 必须包含 `go` 指令。其版本必须是有效的 Go 发布版本：一个正整数，后接点号和非负整数，例如 `1.18`、`1.19`。

`go` 指令表示 `go.work` 预期配合使用的 go 工具链版本。如果 `go.work` 格式发生变化，未来工具链会依据声明的版本解释文件。

`go.work` 最多包含一条 `go` 指令。

```
GoDirective = "go" GoVersion newline .
GoVersion = string | ident .  /* valid release version; see above */
```

示例：

```
go 1.23.0
```

### `toolchain` 指令 {#go-work-file-toolchain}

`toolchain` 指令声明工作区建议使用的 Go 工具链，仅在默认工具链比建议工具链旧时生效。

详细说明请参阅《[Go 工具链](/doc/toolchain)》。

```
ToolchainDirective = "toolchain" ToolchainName newline .
ToolchainName = string | ident .  /* valid toolchain name; see “Go toolchains” */
```

示例：

```
toolchain go1.21.0
```

### `godebug` 指令 {#go-work-file-godebug}

`godebug` 指令声明在该工作区内应用的一项 [GODEBUG 设置](/doc/godebug)，语法和效果与 [`go.mod` 的 `godebug` 指令](#go-mod-file-godebug)相同。启用工作区时，`go.mod` 中的 `godebug` 指令会被忽略。


### `use` 指令 {#go-work-file-use}

`use` 将磁盘上的一个模块加入工作区的主模块集合。参数是包含模块 `go.mod` 文件的目录的相对路径。`use` 不会同时加入该目录子目录中的其他模块；需要通过单独的 `use` 指令，指定它们各自 `go.mod` 所在目录。

```
UseDirective = "use" ( UseSpec | "(" newline { UseSpec } ")" newline ) .
UseSpec = FilePath newline .
FilePath = /* platform-specific relative or absolute file path */

```

示例：

```
use ./mymod  // example.com/mymod

use (
    ../othermod
    ./subdir/thirdmod
)
```

### `replace` 指令 {#go-work-file-replace}

与 `go.mod` 中的 `replace` 类似，`go.work` 中的 `replace` 会将模块特定版本或全部版本的内容替换为其他位置的内容。`go.work` 中不限定版本的替换，会覆盖 `go.mod` 中限定版本的 `replace`。

`go.work` 中的 `replace` 会覆盖工作区各模块中针对同一模块或模块版本的替换规则。

```
ReplaceDirective = "replace" ( ReplaceSpec | "(" newline { ReplaceSpec } ")" newline ) .
ReplaceSpec = ModulePath [ Version ] "=>" FilePath newline
            | ModulePath [ Version ] "=>" ModulePath Version newline .
FilePath = /* platform-specific relative or absolute file path */
```

示例：

```
replace golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5

replace (
    golang.org/x/net v1.2.3 => example.com/fork/net v1.4.5
    golang.org/x/net => example.com/fork/net v1.4.5
    golang.org/x/net v1.2.3 => ./fork/net
    golang.org/x/net => ./fork/net
)
```

## 兼容非模块仓库 {#non-module-compat}

为确保从 `GOPATH` 平滑迁移到模块，`go` 命令在模块模式下，也能从尚未通过添加 [`go.mod`](#glos-go-mod-file) 迁移到模块的仓库中下载并构建包。

`go` 命令从仓库[直接](#vcs)下载指定模块版本时，会根据模块路径查找仓库 URL，将版本映射到仓库修订，再提取该修订的归档。如果[模块路径](#glos-module-path)与[仓库根路径](#glos-repository-root-path)相同，而仓库根目录没有 `go.mod`，`go` 命令会在模块缓存中合成一个 `go.mod`，其中只包含 [`module` 指令](#go-mod-file-module)。合成的 `go.mod` 没有声明依赖的 [`require` 指令](#go-mod-file-require)，因此依赖它的其他模块可能需要补充带 `// indirect` 注释的 `require`，确保每次构建都获取相同版本的依赖。

`go` 命令从[代理](#communicating-with-proxies)下载模块时，会将 `go.mod` 与其他模块内容分开下载。如果原模块没有 `go.mod`，代理应提供一个合成文件。

### `+incompatible` 版本 {#incompatible-versions}

主版本号为 2 或更高的模块，路径必须带匹配的[主版本后缀](#major-version-suffixes)。例如，发布 `v2.0.0` 时，路径必须以 `/v2` 结尾。这样，即使项目的多个主版本都在同一仓库中开发，`go` 命令仍能把它们视为独立模块。

主版本后缀要求是在 `go` 命令引入模块支持时提出的，此前许多仓库已经发布了主版本 `2` 或更高的标签。为兼容这些仓库，`go` 命令会为主版本号至少为 2 且没有 `go.mod` 的版本添加 `+incompatible`。`+incompatible` 表示该版本与更低主版本仍属于同一个模块，因此 `go` 命令可能自动升级到更高的 `+incompatible` 版本，即使这可能导致构建失败。

例如，考虑下面的依赖要求：

```
require example.com/m v4.1.2+incompatible
```

版本 `v4.1.2+incompatible` 引用的是提供 `example.com/m` 模块的仓库中的[语义化版本标签](#glos-semantic-version-tag) `v4.1.2`。模块必须位于仓库根目录，即[仓库根路径](#glos-module-path)也必须是 `example.com/m`，并且不能存在 `go.mod`。该模块可能还有 `v1.5.2` 等更低主版本，`go` 命令可能从这些版本自动升级到 `v4.1.2+incompatible`。有关升级规则，请参阅[最小版本选择（MVS）](#minimal-version-selection)。

仓库在发布 `v2.0.0` 标签后才迁移到模块，通常应发布一个新的主版本。以上例来说，作者应创建路径为 `example.com/m/v5` 的模块，并发布 `v5.0.0`；同时把模块内包的导入前缀从 `example.com/m` 更新为 `example.com/m/v5`。更详细的示例请参阅 [Go 模块：v2 及后续版本](/blog/v2-go-modules)。

注意，仓库标签不应包含 `+incompatible`；`v4.1.2+incompatible` 这样的标签会被忽略。该后缀只出现在 `go` 命令使用的版本中。版本与标签的区别，详见[将版本映射到提交](#vcs-version)。

`+incompatible` 也可以出现在[伪版本](#glos-pseudo-version)中。例如，`v2.0.1-20200722182040-012345abcdef+incompatible` 可以是有效的伪版本。

### 最小模块兼容性 {#minimal-module-compatibility}

主版本号为 2 或更高的模块，必须在[模块路径](#glos-module-path)中包含[主版本后缀](#glos-major-version-suffix)。它可以位于仓库中的[主版本子目录](#glos-major-version-subdirectory)，也可以不使用该目录。在 `GOPATH` 模式下构建时，这会影响其他包如何导入该模块中的包。

通常在 `GOPATH` 模式下，包所在目录由[仓库根路径](#glos-repository-root-path)与仓库内相对目录拼接而成。例如，仓库根路径为 `example.com/repo`，包位于 `sub` 子目录，就会存储在 `$GOPATH/src/example.com/repo/sub`，导入路径为 `example.com/repo/sub`。

对于带主版本后缀的模块，人们可能期待在 `$GOPATH/src/example.com/repo/v2/sub` 中找到 `example.com/repo/v2/sub` 包。这要求模块在仓库的 `v2` 子目录中开发。`go` 命令支持但不强制这种布局，详见[将版本映射到提交](#vcs-version)。

如果模块*不在*主版本子目录中开发，其 `GOPATH` 目录就不包含主版本后缀，包也可能用不带后缀的路径导入。以上例来说，包位于 `$GOPATH/src/example.com/repo/sub`，导入路径为 `example.com/repo/sub`。

这给希望同时支持模块模式和 `GOPATH` 模式的包带来问题：模块模式要求后缀，而 `GOPATH` 模式不要求。

为解决这一问题，Go 1.11 引入了<dfn>最小模块兼容性</dfn>，并回移到 Go 1.9.7 和 1.10.3。在 `GOPATH` 模式下将导入路径解析为目录时，适用以下规则：

* 对于形如 `$modpath/$vn/$dir` 的导入，其中：
  * `$modpath` 是有效模块路径；
  * `$vn` 是主版本后缀；
  * `$dir` 是可能为空的子目录；
* 如果同时满足以下条件：
  * 包 `$modpath/$vn/$dir` 不在任何相关的 [`vendor` 目录](#glos-vendor-directory)中；
  * 导入文件所在目录，或直到 `$GOPATH/src` 为止的某个父目录中，存在 `go.mod`；
  * 对于任意根目录 `$GOPATH[i]`，都不存在 `$GOPATH[i]/src/$modpath/$vn/$suffix`；
  * 对于某个根目录 `$GOPATH[d]`，文件 `$GOPATH[d]/src/$modpath/go.mod` 存在，并声明模块路径为 `$modpath/$vn`；
* 那么，导入 `$modpath/$vn/$dir` 会解析到 `$GOPATH[d]/src/$modpath/$dir` 目录。

这一规则使已迁移到模块的包，即使没有使用主版本子目录，也能在 `GOPATH` 模式下导入其他已迁移到模块的包。

## 支持模块的命令 {#mod-commands}

大多数 `go` 命令可以运行在*模块模式*或 *`GOPATH` 模式*。模块模式下，`go` 命令利用 `go.mod` 查找带版本的依赖，通常从[模块缓存](#glos-module-cache)加载包，缺少模块时会下载。在 `GOPATH` 模式下，`go` 命令忽略模块机制，在 [`vendor` 目录](#glos-vendor-directory)和 `GOPATH` 中查找依赖。

从 Go 1.16 开始，无论是否存在 `go.mod`，都默认启用模块模式。更早版本仅在当前目录或某个父目录中存在 `go.mod` 时启用模块模式。

可以通过 `GO111MODULE` 环境变量控制模块模式，取值为 `on`、`off` 或 `auto`。

* `GO111MODULE=off`：`go` 命令忽略 `go.mod`，运行于 `GOPATH` 模式。
* `GO111MODULE=on` 或未设置：`go` 命令运行于模块模式，即使没有 `go.mod`。并非所有命令都能在没有 `go.mod` 时工作，详见[模块目录之外的模块命令](#commands-outside)。
* `GO111MODULE=auto`：当前目录或任意父目录存在 `go.mod` 时，`go` 命令使用模块模式。这是 Go 1.15 及更早版本的默认行为。`go mod` 子命令，以及带[版本查询](#version-queries)的 `go install`，即使不存在 `go.mod` 也会使用模块模式。

模块模式下，`GOPATH` 不再决定构建时导入路径的含义，但仍用于保存下载的依赖，位置为 `GOPATH/pkg/mod`，详见[模块缓存](#module-cache)；也用于存放安装的命令，位置为 `GOPATH/bin`，除非设置了 `GOBIN`。

### 构建命令 {#build-commands}

所有加载包信息的命令都支持模块，包括：

* `go build`
* `go fix`
* `go generate`
* `go install`
* `go list`
* `go run`
* `go test`
* `go vet`

在模块模式下，这些命令根据 `go.mod` 解释命令行或 Go 源文件中的导入路径。它们接受以下模块命令通用选项。

* `-mod` 控制是否自动更新 `go.mod`，以及是否使用 `vendor`。
  * `-mod=mod`：让 `go` 命令忽略 vendor 目录，并[自动更新](#go-mod-file-updates) `go.mod`，例如导入包不由任何已知模块提供时。
  * `-mod=readonly`：让 `go` 命令忽略 `vendor`，并在需要更新 `go.mod` 时报错。
  * `-mod=vendor`：让 `go` 命令使用 `vendor`。该模式下，`go` 命令不访问网络或模块缓存。
  * 默认情况下，如果 `go.mod` 的 [`go` 版本](#go-mod-file-go)至少为 `1.14`，且存在 `vendor`，`go` 命令就按 `-mod=vendor` 运行；否则，`go` 命令按 `-mod=readonly` 运行。
  * `go get` 不接受此选项，因为它的目的就是修改依赖，而只有 `-mod=mod` 允许这种操作。
* `-modcacherw` 让 `go` 命令在模块缓存中新建目录时使用读写权限，而不是只读权限。如果始终使用该选项，通常通过设置环境变量 `GOFLAGS=-modcacherw` 或执行 `go env -w GOFLAGS=-modcacherw`，就可以直接用 `rm -r` 等命令删除模块缓存，无须先修改权限。无论是否使用过 `-modcacherw`，都可以通过 [`go clean -modcache`](#go-clean-modcache) 删除缓存。
* `-modfile=file.mod` 让 `go` 命令读取并可能写入替代文件，而不是模块根目录中的 `go.mod`。文件名必须以 `.mod` 结尾。仍然必须存在名为 `go.mod` 的文件，用来确定模块根目录，但不会读取其内容。指定 `-modfile` 后，也会使用替代的 `go.sum` 文件：从 `-modfile` 路径去掉 `.mod`，再添加 `.sum` 得到。

### Vendoring：将依赖纳入项目 {#vendoring}

使用模块时，`go` 命令通常从模块来源下载依赖到模块缓存，再从下载的副本中加载包。<dfn>Vendoring</dfn> 将依赖代码复制到项目中，可以与较旧 Go 版本协作，也可以确保一次构建所需的所有文件都位于同一目录树内。

[`go mod vendor`](#go-mod-vendor) 会在[主模块](#glos-main-module)根目录创建 `vendor`，其中包含构建和测试主模块内所有包所需的包副本。仅由主模块之外的包的测试导入的包，不会包含在内。与 [`go mod tidy`](#go-mod-tidy) 等模块命令一样，构造 `vendor` 时，除 `ignore` 外的[构建约束](#glos-build-constraint)都不予考虑。

`go mod vendor` 还会生成 `vendor/modules.txt`，列出纳入 vendor 的包及其来源模块版本。启用 vendor 模式后，[`go list -m`](#go-list-m) 和 [`go version
-m`](#go-version-m) 等命令以此清单作为模块版本信息来源。`go` 命令读取 `vendor/modules.txt` 时，会检查其中的模块版本与 `go.mod` 是否一致。如果生成 `vendor/modules.txt` 后又修改了 `go.mod`，`go` 命令会报错。应重新运行 `go mod vendor`，更新 `vendor`。

如果主模块根目录存在 `vendor`，且主模块 [`go.mod`](#glos-go-mod-file) 中的 [`go` 版本](#go-mod-file-go)至少为 `1.14`，就会自动使用它。要显式启用 vendor 模式，运行 `go` 命令时指定 `-mod=vendor`；要禁用，则使用 `-mod=readonly` 或 `-mod=mod`。

启用 vendor 模式后，`go build`、`go test` 等[构建命令](#build-commands)从 `vendor` 加载包，不访问网络或本地模块缓存。[`go list -m`](#go-list-m) 只输出 `go.mod` 中列出的模块信息。[`go mod download`](#go-mod-download)、[`go mod tidy`](#go-mod-tidy) 等 `go mod` 命令的行为不受 vendor 模式影响，仍会下载模块并访问缓存。[`go get`](#go-get) 的行为同样不变。

与 [`GOPATH` 模式下的 vendoring](/s/go15vendor) 不同，`go` 命令忽略主模块根目录之外的 vendor 目录。由于不会使用其他模块中的 vendor，`go` 命令生成[模块 zip 文件](#zip-files)时也不会包含这些目录，但请注意已知缺陷 [#31562](/issue/31562) 和 [#37397](/issue/37397)。

### `go get` {#go-get}

用法：

```
go get [-d] [-t] [-u] [-tool] [build flags] [packages]
```

示例：

```
# Upgrade a specific module.
$ go get golang.org/x/net

# Upgrade modules that provide packages imported by packages in the main module.
$ go get -u ./...

# Upgrade or downgrade to a specific version of a module.
$ go get golang.org/x/text@v0.3.2

# Update to the commit on the module's master branch.
$ go get golang.org/x/text@master

# Remove a dependency on a module and downgrade modules that require it
# to versions that don't require it.
$ go get golang.org/x/text@none

# Upgrade the minimum required Go version for the main module.
$ go get go

# Upgrade the suggested Go toolchain, leaving the minimum Go version alone.
$ go get toolchain

# Upgrade to the latest patch release of the suggested Go toolchain.
$ go get toolchain@patch
```

`go get` 更新[主模块](#glos-main-module) [`go.mod`](#go-mod-file) 中的依赖；在下文所述的旧版行为中，还会构建并安装命令行指定的包。

第一步是确定要更新的模块。`go get` 接受包、包模式和模块路径列表作为参数。指定包时，`go get` 更新提供该包的模块。指定包模式，例如 `all` 或包含 `...` 通配符的路径时，`go get` 将模式展开为包集合，再更新提供这些包的模块。如果参数只对应模块、不对应包，例如 `golang.org/x/net` 根目录没有包，`go get` 就只更新模块，不构建包。如果没有参数，`go get` 按指定了 `.` 处理，即当前目录中的包；可以配合 `-u` 更新提供其导入包的模块。

每个参数都可以带<dfn>版本查询后缀</dfn>来指定目标版本，例如 `go get golang.org/x/text@v0.3.0`。后缀由 `@` 和[版本查询](#version-queries)组成，可以指定具体版本（`v0.3.0`）、版本前缀（`v0.3`）、分支或标签（`master`）、修订（`1234abcd`），或特殊查询 `latest`、`upgrade`、`patch`、`none`。未指定版本时，`go get` 使用 `@upgrade`。

`go get` 将参数解析为具体模块和版本后，`go
get` 会在主模块 `go.mod` 中添加、修改或删除 [`require` 指令](#go-mod-file-require)，确保以后仍使用预期版本。注意，`go.mod` 声明的是*最低版本*，添加新依赖后可能自动提高。支持模块的命令如何选择版本和解决冲突，详见[最小版本选择（MVS）](#minimal-version-selection)。

添加、升级或降级命令行指定的模块时，如果目标版本要求其他模块的更高版本，也可能同时升级这些模块。例如，将 `example.com/a` 升级到 `v1.5.0`，而它要求 `example.com/b` 为 `v1.2.0`。如果当前对 `example.com/b` 的要求是 `v1.1.0`，执行 `go get example.com/a@v1.5.0` 也会将 `example.com/b` 升级到 `v1.2.0`。

![go get 升级传递依赖要求](/doc/mvs/get-upgrade.svg)

降级或移除命令行指定的模块时，也可能降级其他模块。延续上例，如果把 `example.com/b` 降级到 `v1.1.0`，`example.com/a` 也会降级到只要求 `example.com/b` 为 `v1.1.0` 或更低的版本。

![go get 降级传递依赖要求](/doc/mvs/get-downgrade.svg)

可以通过版本后缀 `@none` 移除依赖要求，这是一种特殊降级。依赖被移除模块的其他模块会按需降级或移除。即使主模块仍导入该模块中的一个或多个包，也可以移除其要求；这时，下一次构建命令可能重新添加依赖要求。

如果同一个模块被要求使用两个不同版本，无论是命令行显式指定，还是为满足升级、降级而产生，`go get` 都会报错。

`go get` 选出新的版本集合后，会检查新选中的模块版本，以及提供命令行指定包的模块，是否已经[撤回](#glos-retracted-version)或[弃用](#glos-deprecated-module)。`go get` 会为每个发现的已撤回版本或已弃用模块输出警告。可以用 [`go list -m -u
all`](#go-list-m) 检查全部依赖的撤回和弃用情况。

按照这里介绍的旧版安装行为，`go get` 更新 `go.mod` 后，会构建命令行指定的包，并把可执行文件安装到 `GOBIN` 指定的目录。其默认值为 `$GOPATH/bin`；如果未设置 `GOPATH`，则为 `$HOME/go/bin`。

`go get` 支持以下选项：

* `-d` 让 `go get` 不构建或安装包。指定 `-d` 时，`go get` 只管理 `go.mod` 中的依赖。从 Go 1.17 开始，不带 `-d` 使用 `go get` 构建和安装包的方式已弃用；从 Go 1.18 开始，`-d` 始终启用。
* `-u` 让 `go get` 升级那些提供命令行指定包的直接或间接导入包的模块。除非已经要求了更高的预发布版本，否则 `-u` 选中的每个模块都会升级到其最新版本。
* `-u=patch`，注意不是 `-u patch`，也让 `go get` 升级依赖，但 `go get` 只升级到各依赖的最新补丁版本，类似于 `@patch` 查询。
* `-t` 让 `go get` 同时考虑构建命令行指定包的测试所需的模块。组合使用 `-t` 和 `-u` 时，`go get` 也会更新测试依赖。
* 不应再使用 `-insecure`。它允许 `go get` 通过 HTTP 等不安全协议解析自定义导入路径，以及从仓库和模块代理获取内容。应改用提供更细粒度控制的 `GOINSECURE` [环境变量](#environment-variables)。
* `-tool` 让 go 为每个列出的包在 `go.mod` 中添加对应 tool 行。若同时使用 `-tool` 与 `@none`，则移除该行。

从 Go 1.16 开始，推荐用 [`go install`](#go-install) 构建和安装程序。带 `@latest`、`@v1.4.6` 等版本后缀时，`go install` 使用模块模式构建，并忽略当前目录或父目录中可能存在的 `go.mod`。

`go get` 更专注于管理 `go.mod` 中的依赖要求。`-d` 已弃用，并且自 Go 1.18 起始终启用。

### `go install` {#go-install}

用法：

```
go install [build flags] [packages]
```

示例：

```
# Install the latest version of a program,
# ignoring go.mod in the current directory (if any).
$ go install golang.org/x/tools/gopls@latest

# Install a specific version of a program.
$ go install golang.org/x/tools/gopls@v0.6.4

# Install a program at the version selected by the module in the current directory.
$ go install golang.org/x/tools/gopls

# Install all programs in a directory.
$ go install ./cmd/...
```

`go install` 构建并安装命令行路径指定的包。可执行文件，即 `main` 包，安装到 `GOBIN` 指定的目录，默认是 `$GOPATH/bin`；未设置 `GOPATH` 时则为 `$HOME/go/bin`。`$GOROOT` 中的可执行文件安装到 `$GOROOT/bin` 或 `$GOTOOLDIR`，而不是 `$GOBIN`。不可执行的包只构建并缓存，不安装。

从 Go 1.16 开始，如果参数带有 `@latest`、`@v1.0.0` 等版本后缀，`go install` 使用模块模式构建，并忽略当前目录或父目录中可能存在的 `go.mod`。这适合安装可执行程序而不影响主模块依赖。

为消除构建所用模块版本的歧义，只要任意参数带有版本后缀，所有参数就必须满足以下约束：

* 参数必须是包路径或包含 `...` 通配符的包模式，不能是 `fmt` 等标准库包、`std`、`cmd`、`all`、`work`、`tool` 等元模式，也不能是相对或绝对文件路径。注意，不带版本后缀时可以使用 `go install tool`，详见下文。
* 所有参数必须具有相同版本后缀。即使不同查询最终指向同一版本，也不允许混用。
* 所有参数必须指向同一模块、同一版本中的包。
* 包路径参数必须对应 `main` 包；模式参数也只匹配 `main` 包。
* 此时没有任何模块被视为[主模块](#glos-main-module)。
  * 如果指定包所属模块有 `go.mod`，其中不能包含会导致“作为主模块时解释结果不同”的 `replace` 或 `exclude` 指令。
  * 模块不能要求自身的更高版本。
  * 不使用任何模块的 vendor 目录，因为[模块 zip 文件](#zip-files)不包含它们，`go install` 也就不会下载它们。

支持的查询语法详见[版本查询](#version-queries)。Go 1.15 及更早版本不支持在 `go install` 中使用版本查询。

参数没有版本后缀时，`go install` 可能运行在模块模式或 `GOPATH` 模式，具体取决于 `GO111MODULE` 和是否存在 `go.mod`，详见[支持模块的命令](#mod-commands)。启用模块模式时，`go
install` 在主模块上下文中运行，主模块可能与待安装包所属模块不同。模块模式下，可以在模块内执行 `go install tool`，安装该模块的全部工具。

### `go tool` {#go-tool}

用法：

```
go tool [-n] command [args...]
```

示例：

```
$ go tool golang.org/x/tools/cmd/stringer
$ go tool stringer
```

模块模式下，`go tool` 可以构建并运行 `go.mod` 中通过 [`tool` 指令](#go-mod-file-tool)声明的工具。可以使用该工具的完整包路径指定命令；如果工具默认的二进制名称在已安装工具中唯一，也可以使用该名称，即包路径最后一段，去掉主版本后缀。

### `go list -m` {#go-list-m}

用法：

```
go list -m [-u] [-retracted] [-versions] [list flags] [modules]
```

示例：

```
$ go list -m all
$ go list -m -versions example.com/m
$ go list -m -json example.com/m@latest
```

`-m` 使 `go list` 列出模块而非包。该模式下，`go list` 的参数可以是模块、包含 `...` 的模块模式、[版本查询](#version-queries)，或匹配[构建列表](#glos-build-list)全部模块的特殊模式 `all`。不指定参数时，列出[主模块](#glos-main-module)。

列出模块时，`-f` 仍用于指定作用于 Go 结构体的格式模板，只是结构体变成了 `Module`：

```
type Module struct {
    Path       string        // module path
    Version    string        // module version
    Versions   []string      // available module versions
    Replace    *Module       // replaced by this module
    Time       *time.Time    // time version was created
    Update     *Module       // available update (with -u)
    Main       bool          // is this the main module?
    Indirect   bool          // module is only indirectly needed by main module
    Dir        string        // directory holding local copy of files, if any
    GoMod      string        // path to go.mod file describing module, if any
    GoVersion  string        // go version used in module
    Retracted  []string      // retraction information, if any (with -retracted or -u)
    Deprecated string        // deprecation message, if any (with -u)
    Error      *ModuleError  // error loading module
}

type ModuleError struct {
    Err string // the error itself
}
```

默认输出模块路径，然后输出版本和替换信息，如果存在。例如，`go list -m all` 可能输出：

```
example.com/main/module
golang.org/x/net v0.1.0
golang.org/x/text v0.3.0 => /tmp/text
rsc.io/pdf v0.1.1
```

`Module` 的 `String` 方法负责格式化这一行，因此默认格式等价于 {{raw "`-f '{{.String}}'`"}}。

模块被替换时，`Replace` 字段描述替换模块，`Dir` 则设置为替换模块的源代码目录，如果该目录存在。也就是说，`Replace` 非 nil 时，`Dir` 为 `Replace.Dir`，不能通过它访问被替换的原始代码。

`-u` 添加可用升级信息。如果模块最新版本高于当前版本，`list -u` 会将较新模块的信息写入 `Update`。`list -u` 还会显示当前版本是否[已撤回](#glos-retracted-version)、模块是否[已弃用](#go-mod-file-module-deprecation)。`String` 方法会在当前版本之后用方括号显示可升级版本。例如，`go list -m -u all` 可能输出：

```
example.com/main/module
golang.org/x/old v1.9.9 (deprecated)
golang.org/x/net v0.1.0 (retracted) [v0.2.0]
golang.org/x/text v0.3.0 [v0.4.0] => /tmp/text
rsc.io/pdf v0.1.1 [v0.1.2]
```

对于工具而言，`go list -m -u -json all` 通常更便于解析。

`-versions` 让 `list` 将模块的 `Versions` 字段设置为全部已知版本，按语义化版本从低到高排序。它也将默认输出改为模块路径后接空格分隔的版本列表。除非同时指定 `-retracted`，否则列表不包含已撤回版本。

`-retracted` 让 `list` 在 `-versions` 输出中显示已撤回版本，并在解析[版本查询](#version-queries)时考虑它们。例如，`go list -m
-retracted example.com/m@latest` 显示 `example.com/m` 的最高正式版本或预发布版本，即使已经撤回。该版本的 `go.mod` 用于加载 [`retract` 指令](#go-mod-file-retract)和[弃用信息](#go-mod-file-module-deprecation)。`-retracted` 在 Go 1.16 中加入。

模板函数 `module` 接受一个字符串参数，必须是模块路径或查询，并将指定模块以 `Module` 结构体返回。如果出错，结果仍为 `Module`，但 `Error` 字段非 nil。

### `go mod download` {#go-mod-download}

用法：

```
go mod download [-x] [-json] [-reuse=old.json] [modules]
```

示例：

```
$ go mod download
$ go mod download golang.org/x/mod@v0.2.0
```

`go mod download` 将指定模块下载到[模块缓存](#glos-module-cache)。参数可以是选择主模块依赖的模块路径或模块模式，也可以是 `path@version` 形式的[版本查询](#version-queries)。没有参数时，`download` 作用于[主模块](#glos-main-module)的全部依赖。

`go` 命令在正常运行时会按需自动下载模块。`go mod download` 主要用于预先填充模块缓存，或为[模块代理](#glos-module-proxy)准备要提供的数据。

默认情况下，`download` 不向标准输出写入内容，只向标准错误输出进度和错误信息。

`-json` 让 `download` 向标准输出写出一系列 JSON 对象，描述每个下载模块或失败情况，对应以下 Go 结构体：

```
type Module struct {
    Path     string // module path
    Query    string // version query corresponding to this version
    Version  string // module version
    Error    string // error loading module
    Info     string // absolute path to cached .info file
    GoMod    string // absolute path to cached .mod file
    Zip      string // absolute path to cached .zip file
    Dir      string // absolute path to cached source root directory
    Sum      string // checksum for path, version (as in go.sum)
    GoModSum string // checksum for go.mod (as in go.sum)
    Origin   any    // provenance of module
    Reuse    bool   // reuse of old module info is safe
}
```

`-x` 让 `download` 将其执行的命令，即 `download` 运行的命令，输出到标准错误。

-reuse 接受一个文件名，该文件保存先前执行 'go mod download -json' 得到的 JSON 输出。go 命令可以据此判断模块自上次调用后未发生变化，从而避免重新下载。未重新下载的模块会在新输出中将 Reuse 字段设为 true。通常模块缓存会自动实现这种复用；对于不保留模块缓存的系统，-reuse 会很有用。

### `go mod edit` {#go-mod-edit}

用法：

```
go mod edit [editing flags] [-fmt|-print|-json] [go.mod]
```

示例：

```
# Add a replace directive.
$ go mod edit -replace example.com/a@v1.0.0=./a

# Remove a replace directive.
$ go mod edit -dropreplace example.com/a@v1.0.0

# Set the go version, add a requirement, and print the file
# instead of writing it to disk.
$ go mod edit -go=1.14 -require=example.com/m@v1.0.0 -print

# Format the go.mod file.
$ go mod edit -fmt

# Format and print a different .mod file.
$ go mod edit -print tools.mod

# Print a JSON representation of the go.mod file.
$ go mod edit -json
```

`go mod edit` 提供编辑和格式化 `go.mod` 的命令行接口，主要供工具和脚本使用。`go mod edit` 只读取一个 `go.mod`，不查询其他模块的信息。默认情况下，`go mod edit` 读写主模块的 `go.mod`，也可以在编辑选项之后指定其他目标文件。

编辑选项指定一系列编辑操作。

* `-module` 修改模块路径，即 `go.mod` 的 module 行。
* `-go=version` 设置预期的 Go 语言版本。
* `-require=path@version` 和 `-droprequire=path` 添加或删除对指定模块路径及版本的依赖要求。`-require` 会覆盖 `path` 已有的要求。这些选项主要供理解模块图的工具使用。用户应优先使用 `go get
  path@version` 或 `go get path@none`，它们会按需调整 `go.mod` 的其他内容，以满足其他模块施加的约束。详见 [`go
  get`](#go-get)。
* `-exclude=path@version` 和 `-dropexclude=path@version` 添加或删除对指定模块版本的排除。若排除项已存在，`-exclude=path@version` 不产生变化。
* `-replace=old[@v]=new[@v]` 为指定模块路径及版本添加替换。如果 `old@v` 中省略 `@v`，就添加左侧不限定版本的替换，适用于旧模块路径的全部版本。如果 `new@v` 中省略 `@v`，新路径必须是本地模块根目录，而不是模块路径。`-replace` 会覆盖 `old[@v]` 的冗余替换，因此省略 `@v` 会删除针对特定版本的替换。
* `-dropreplace=old[@v]` 删除指定模块路径及版本的替换。提供 `@v` 时，只删除对应版本的替换，已有的左侧不限定版本的替换仍可能生效。省略 `@v` 时，则删除不限定版本的替换。
* `-retract=version` 和 `-dropretract=version` 添加或删除指定版本的撤回，既可以是 `v1.2.3` 这样的单个版本，也可以是 `[v1.1.0,v1.2.0]` 这样的区间。`-retract` 不能为 `retract` 指令添加原因注释。建议提供此类注释，`go list -m -u` 等命令可能显示它们。
* `-tool=path` 和 `-droptool=path` 为指定路径添加或删除 `tool` 指令，但不会向构建图加入所需依赖。用户应优先使用 `go get -tool path` 添加工具，用 `go get -tool path@none` 移除工具。

编辑选项可以重复，修改会按给定顺序应用。

`go mod edit` 还提供控制输出的选项。

* `-fmt` 只重新格式化 `go.mod`，不做其他更改。其他会修改或重写 `go.mod` 的操作，也会自动格式化。因此，只有不指定其他选项时才需要它，例如 `go mod edit -fmt`。
* `-print` 以文本格式输出最终 `go.mod`，而不写回磁盘。
* `-json` 以 JSON 格式输出最终 `go.mod`，而不以文本格式写回磁盘。JSON 对应以下 Go 类型：

```
type Module struct {
    Path    string
    Version string
}

type GoMod struct {
    Module  ModPath
    Go      string
    Require []Require
    Exclude []Module
    Replace []Replace
    Retract []Retract
}

type ModPath struct {
    Path       string
    Deprecated string
}

type Require struct {
    Path     string
    Version  string
    Indirect bool
}

type Replace struct {
    Old Module
    New Module
}

type Retract struct {
    Low       string
    High      string
    Rationale string
}

type Tool struct {
    Path      string
}
```

注意，这只描述 `go.mod` 本身，不包含其间接引用的其他模块。要获取构建可用的全部模块，使用 `go list -m -json all`，详见 [`go list -m`](#go-list-m)。

例如，工具可以解析 `go mod edit -json` 的输出，将 `go.mod` 读取为数据结构，再通过带 `-require`、`-exclude` 等选项的 `go mod edit` 修改它。

工具也可以通过 [`golang.org/x/mod/modfile`](https://pkg.go.dev/golang.org/x/mod/modfile?tab=doc) 包解析、编辑和格式化 `go.mod`。

### `go mod graph` {#go-mod-graph}

用法：

```
go mod graph [-go=version]
```

`go mod graph` 以文本形式输出已应用替换规则的[模块依赖要求图](#glos-module-graph)。例如：

```
example.com/main example.com/a@v1.1.0
example.com/main example.com/b@v1.2.0
example.com/a@v1.1.0 example.com/b@v1.1.1
example.com/a@v1.1.0 example.com/c@v1.3.0
example.com/b@v1.1.0 example.com/c@v1.1.0
example.com/b@v1.2.0 example.com/c@v1.2.0
```

图中的每个顶点代表模块的一个特定版本，每条边代表对依赖最低版本的要求。

`go mod graph` 每行输出一条边。每行有两个空格分隔的字段：一个模块版本及其一个依赖。模块版本以 `path@version` 表示。主模块没有版本，因此不带 `@version` 后缀。

`-go` 让 `go mod graph` 按指定 Go 版本的加载方式报告模块图，而不是使用 `go.mod` 中 [`go` 指令](#go-mod-file-go)指定的版本。

版本选择规则详见[最小版本选择（MVS）](#minimal-version-selection)。要输出选定版本，请参阅 [`go list -m`](#go-list-m)；要了解为何需要某个模块，请参阅 [`go mod why`](#go-mod-why)。

### `go mod init` {#go-mod-init}

用法：

```
go mod init [module-path]
```

示例：

```
go mod init
go mod init example.com/m
```

`go mod init` 在当前目录初始化并写入一个新的 `go.mod`，相当于创建以当前目录为根的新模块。已有 `go.mod` 时不能执行。

`init` 接受一个可选参数，即新模块的[模块路径](#glos-module-path)。如何选择路径，详见[模块路径](#module-path)。省略参数时，`init` 会根据 `.go` 文件中的导入注释，以及当前目录（如果位于 `GOPATH`）尝试推断路径。

### `go mod tidy` {#go-mod-tidy}

用法：

```
go mod tidy [-e] [-v] [-x] [-diff] [-go=version] [-compat=version]
```

`go mod tidy` 确保 `go.mod` 与模块源代码一致。它添加构建当前模块的包及其依赖所需的缺失要求，删除不提供相关包的模块要求，同时补充 `go.sum` 中缺失的条目并移除不再需要的条目。

`-e` 在 Go 1.16 中引入，让 `go mod tidy` 即使在加载包时遇到错误，也尽量继续执行。

`-v` 让 `go mod tidy` 向标准错误输出被移除模块的信息。

`-x` 让 `go mod tidy` 输出 `tidy` 执行的命令。

`-diff` 让 `go mod tidy` 不修改 go.mod 或 go.sum，而是以统一差异格式输出所需更改。如果差异非空，命令以非零状态码退出。

`go mod tidy` 会递归加载[主模块](#glos-main-module)的所有包、所有工具，以及它们导入的所有包，包括测试导入的包，也包括其他模块中的测试。`go mod tidy` 按所有构建标签都启用的方式处理，因此会考虑平台专用文件和需要自定义构建标签的文件，即使通常不会构建它们。唯一例外是 `ignore`：它不会启用，因此带 `// +build ignore` 约束的文件不参与处理。另请注意，`go mod tidy` 不主动考虑主模块中名为 `testdata`、或名称以 `.`、`_` 开头的目录中的包，除非其他包显式导入它们。

加载这些包后，`go mod tidy` 确保每个提供相关包的模块，都在主模块 `go.mod` 中有 `require`；如果主模块声明 `go 1.16` 或更低版本，也可以通过其他已要求模块间接引入。`go mod tidy` 会为缺失模块添加其最新版本要求，`latest` 的定义详见[版本查询](#version-queries)。对于不提供上述包集合中任何包的模块，`go mod tidy` 会移除对应 `require`。

`go mod tidy` 也可能添加或移除 `require` 上的 `// indirect` 注释。`// indirect` 表示模块没有提供被主模块中的包直接导入的包。何时添加 `// indirect` 依赖及注释，详见 [`require` 指令](#go-mod-file-require)。

指定 `-go` 时，`go mod tidy` 将 [`go` 指令](#go-mod-file-go)更新为指定版本，并据此启用或禁用[模块图剪枝](#graph-pruning)和[模块延迟加载](#lazy-loading)，按需添加或删除间接依赖要求。

默认情况下，`go mod tidy` 会检查：使用比 `go` 指令指定版本低一个版本的 Go 加载模块图时，模块的[选定版本](#glos-selected-version)是否保持不变。也可以通过 `-compat` 显式指定兼容性检查版本。

### `go mod vendor` {#go-mod-vendor}

用法：

```
go mod vendor [-e] [-v] [-o]
```

`go mod vendor` 在[主模块](#glos-main-module)根目录创建 `vendor`，其中保存构建和测试主模块包所需的全部包副本。仅被主模块外部包的测试导入的包，不包含在内。与 [`go mod tidy`](#go-mod-tidy) 等模块命令一样，构造 `vendor` 时忽略除 `ignore` 以外的[构建约束](#glos-build-constraint)。

启用 vendor 模式时，`go` 命令从 `vendor` 加载包，不再从来源下载模块到模块缓存后使用下载副本。更多信息请参阅 [Vendoring](#vendoring)。

`go mod vendor` 还会生成 `vendor/modules.txt`，记录 vendor 中的包及其来源模块版本。启用 vendor 模式后，[`go list -m`](#go-list-m) 和 [`go version
-m`](#go-version-m) 等命令使用该清单获取模块版本信息。`go` 命令读取 `vendor/modules.txt` 时，会检查版本与 `go.mod` 是否一致。如果生成 `vendor/modules.txt` 后修改了 `go.mod`，就应重新运行 `go mod vendor`。

注意，`go mod vendor` 会先删除已有 `vendor`，再重新构建。不要直接修改 vendor 中的包。`go` 命令不会检查 `vendor` 中的包是否被修改，不过可以重新运行 `go mod vendor` 并确认没有产生差异，来检查 `vendor` 的完整性。

`-e` 在 Go 1.16 中引入，让 `go mod vendor` 在加载包时遇到错误后仍尽量继续。

`-v` 让 `go mod vendor` 向标准错误输出纳入 vendor 的模块名和包名。

`-o` 在 Go 1.18 中引入，让 `go mod vendor` 将依赖目录树输出到指定目录，而不是 `vendor`。参数可以是绝对路径，也可以是相对于模块根目录的路径。

### `go mod verify` {#go-mod-verify}

用法：

```
go mod verify
```

`go mod verify` 检查[模块缓存](#glos-module-cache)中的[主模块](#glos-main-module)依赖，自下载后是否被修改。它通过 `go mod verify` 计算每个下载模块的 [`.zip` 文件](#zip-files)及解压目录的哈希，与首次下载时记录的哈希比较。`go mod verify` 检查[构建列表](#glos-build-list)中的每个模块，该列表可以用 [`go list -m
all`](#go-list-m) 查看。

如果全部模块均未修改，`go mod verify` 输出“all modules verified”。否则，会报告被修改的模块，并以非零状态退出。

所有支持模块的命令都会检查：主模块 `go.sum` 中的哈希是否与下载到模块缓存的模块哈希一致。如果 `go.sum` 缺少哈希，例如首次使用模块时，`go` 命令会通过[校验和数据库](#checksum-database)验证，除非模块路径匹配 `GOPRIVATE` 或 `GONOSUMDB`。详见[验证模块真实性](#authenticating)。

相比之下，`go mod verify` 检查的是模块 `.zip` 及解压目录是否与首次下载时记录在缓存中的哈希一致，用于发现模块下载并验证*之后*发生的修改。`go mod verify` 不会下载缓存中不存在的模块内容，也不使用 `go.sum` 验证模块内容。不过，`go mod verify` 可能为执行[最小版本选择](#minimal-version-selection)而下载 `go.mod`，会使用 `go.sum` 验证这些文件，并可能为缺失哈希添加 `go.sum` 条目。

### `go mod why` {#go-mod-why}

用法：

```
go mod why [-m] [-vendor] packages...
```

`go mod why` 显示导入图中从主模块到每个指定包的一条最短路径。

输出按段落组织，命令行指定的每个包或模块对应一段，段落之间以空行分隔。每段以 `#` 开头的注释行标明目标包或模块，后续每行一个包，构成导入图中的路径。如果主模块没有引用该包或模块，就只显示一条括号包围的说明。

例如：

```
$ go mod why golang.org/x/text/language golang.org/x/text/encoding
# golang.org/x/text/language
rsc.io/quote
rsc.io/sampler
golang.org/x/text/language

# golang.org/x/text/encoding
(main module does not need package golang.org/x/text/encoding)
```

`-m` 让 `go mod why` 将参数视为模块列表，`go mod why` 会输出通往每个模块内任意包的一条路径。即使指定 `-m`，`go mod why` 查询的仍是包导入图，而不是 [`go mod graph`](#go-mod-graph) 输出的模块图。

`-vendor` 让 `go mod why` 忽略主模块之外的包的测试导入，行为与 [`go mod vendor`](#go-mod-vendor) 相同。默认情况下，`go mod why` 考虑 `all` 模式匹配的包图。从 Go 1.16 开始，对于在 `go.mod` 的 [`go` 指令](#go-mod-file-go)中声明 `go 1.16` 或更高版本的模块，此选项没有作用，因为 `all` 已改为与 `go mod vendor` 使用相同的包集合。

### `go version -m` {#go-version-m}

用法：

```
go version [-m] [-v] [file ...]
```

示例：

```
# Print Go version used to build go.
$ go version

# Print Go version used to build a specific executable.
$ go version ~/go/bin/gopls

# Print Go version and module versions used to build a specific executable.
$ go version -m ~/go/bin/gopls

# Print Go version and module versions used to build executables in a directory.
$ go version -m ~/go/bin/
```

`go version` 报告构建命令行指定的每个可执行文件时使用的 Go 版本。

不指定文件时，`go version` 输出自身的版本信息。

指定目录时，`go version` 递归遍历目录，查找可识别的 Go 二进制文件并报告版本。默认情况下，`go
version` 不报告扫描时遇到的不可识别文件；`-v` 会让它也报告这些文件。

`-m` 让 `go version` 输出可执行文件内嵌的模块版本信息，如果存在。对于每个可执行文件，`go version -m` 输出一个制表符分隔列的表格，例如：

```
$ go version -m ~/go/bin/goimports
/home/jrgopher/go/bin/goimports: go1.14.3
        path    golang.org/x/tools/cmd/goimports
        mod     golang.org/x/tools      v0.0.0-20200518203908-8018eb2c26ba      h1:0Lcy64USfQQL6GAJma8BdHCgeofcchQj+Z7j0SXYAzU=
        dep     golang.org/x/mod        v0.2.0          h1:KU7oHjnv3XNWfa5COkzUifxZmxp1TyI7ImMXqFxLwvQ=
        dep     golang.org/x/xerrors    v0.0.0-20191204190536-9bdfabe68543      h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=
```

表格格式未来可能变化。相同信息也可以通过 [`runtime/debug.ReadBuildInfo`](https://pkg.go.dev/runtime/debug?tab=doc#ReadBuildInfo) 获取。

表格每行的含义由第一列的单词决定。

* **`path`**：构建可执行文件的 `main` 包路径。
* **`mod`**：包含 `main` 包的模块。各列依次是模块路径、版本和校验和。[主模块](#glos-main-module)的版本为 `(devel)`，没有校验和。
* **`dep`**：提供了链接进可执行文件的一个或多个包的模块，格式与 `mod` 相同。
* **`=>`**：上一行模块的[替换项](#go-mod-file-replace)。如果是本地目录，只显示目录路径，不显示版本和校验和；如果是模块版本，则与 `mod`、`dep` 一样显示路径、版本和校验和。被替换的原模块没有校验和。

### `go clean -modcache` {#go-clean-modcache}

用法：

```
go clean [-modcache]
```

`-modcache` 让 [`go
clean`](/cmd/go/#hdr-Remove_object_files_and_cached_files) 删除整个[模块缓存](#glos-module-cache)，包括已解压的带版本依赖源代码。

这通常是删除模块缓存的最佳方式。默认情况下，模块缓存中的大多数文件和目录都是只读的，以防测试或编辑器在[验证真实性](#authenticating)后意外修改文件。但这也会使 `rm -r` 等命令失败，因为删除文件前必须先让其父目录可写。

[`go
build`](/cmd/go/#hdr-Compile_packages_and_dependencies) 等模块命令接受的 `-modcacherw`，会让模块缓存中的新目录可写。要为所有模块命令传入 `-modcacherw`，可以将它加入 `GOFLAGS`。`GOFLAGS` 可通过环境变量设置，也可以用 [`go env
-w`](/cmd/go/#hdr-Print_Go_environment_information) 设置。例如，下面的命令会永久保存该设置：

```
go env -w GOFLAGS=-modcacherw
```

应谨慎使用 `-modcacherw`，避免修改模块缓存中的文件。可以使用 [`go mod verify`](#go-mod-verify) 检查缓存文件的完整性；主模块的 `go.sum` 则记录依赖的预期校验和。

### 版本查询 {#version-queries}

多个命令支持通过*版本查询*指定模块版本，查询位于命令行中的模块或包路径之后，以 `@` 引出。

示例：

```
go get example.com/m@latest
go mod download example.com/m@master
go list -m -json example.com/m@e3702bed2
```

版本查询可以采用以下形式：

* 完整语义化版本，例如 `v1.2.3`，选择具体版本。语法详见[版本](#versions)。
* 语义化版本前缀，例如 `v1` 或 `v1.2`，选择具有该前缀的最高可用版本。
* 语义化版本比较，例如 {{raw "`<v1.2.3` or `>=v1.5.6`"}}，选择满足比较条件且最接近目标的可用版本：`>` 和 `>=` 选择最低版本，{{raw "`<` and `<=`"}} 选择最高版本。
* 底层源代码仓库的修订标识符，例如提交哈希前缀、修订标签或分支名。如果修订有语义化版本标签，选择该版本；否则选择该提交对应的[伪版本](#glos-pseudo-version)。不能用这种形式选择名称会被其他版本查询规则匹配的分支或标签。例如，`v2` 查询选择以 `v2` 开头的最新版本，而不是名为 `v2` 的分支。
* `latest`：选择最高可用正式版本。没有正式版本时，`latest` 选择最高预发布版本；没有任何版本标签时，`latest` 选择仓库默认分支最新提交的伪版本。
* `upgrade`：类似 `latest`，但如果当前已要求的版本高于 `latest` 会选出的版本，例如更高的预发布版本，`upgrade` 就保留当前版本。
* `patch`：选择与当前要求具有相同主版本号和次版本号的最新可用版本。如果当前没有要求版本，`patch` 等同于 `latest`。从 Go 1.16 开始，[`go get`](#go-get) 使用 `patch` 时必须已有当前版本，但 `-u=patch` 没有这一要求。

除查询明确指定的版本或修订之外，所有查询都基于 `go list -m -versions` 报告的可用版本，详见 [`go list
-m`](#go-list-m)。该列表只包含标签版本，不含伪版本。主模块 [`go.mod`](#glos-go-mod-file) 中 [`exclude` 指令](#go-mod-file-exclude)排除的版本不参与选择。同一模块 `latest` 版本的 `go.mod` 中 [`retract` 指令](#go-mod-file-retract)撤回的版本也会忽略，例外是带 `-retracted` 的 [`go list -m`](#go-list-m)，以及加载 `retract` 指令本身的过程。

[正式版本](#glos-release-version)优先于预发布版本。例如，同时存在 `v1.2.2` 和 `v1.2.3-pre` 时，`latest` 选择 `v1.2.2`，即使 `v1.2.3-pre` 更高。{{raw "`<v1.2.4`"}} 也选择 `v1.2.2`，即使 `v1.2.3-pre` 更接近 `v1.2.4`。如果既没有正式版本，也没有预发布版本，`latest`、`upgrade` 和 `patch` 会选择仓库默认分支最新提交的伪版本，其他查询则报错。

### 模块目录之外的模块命令 {#commands-outside}

支持模块的 Go 命令通常在[主模块](#glos-main-module)上下文中运行，该模块由工作目录或父目录的 `go.mod` 定义。有些命令可以在没有 `go.mod` 时使用模块模式，但大多数命令在不存在 `go.mod` 时会改变行为或报错。

有关启用和禁用模块感知模式的信息，请参阅[支持模块的命令](#mod-commands)。

<table class="ModTable">
  <thead>
    <tr>
      <th>命令</th>
      <th>行为</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>
        <code>go build</code><br>
        <code>go doc</code><br>
        <code>go fix</code><br>
        <code>go fmt</code><br>
        <code>go generate</code><br>
        <code>go install</code><br>
        <code>go list</code><br>
        <code>go run</code><br>
        <code>go test</code><br>
        <code>go vet</code>
      </td>
      <td>
        只能加载、导入和构建标准库中的包，以及在命令行中以 <code>.go</code> 文件指定的包。无法构建其他模块中的包，因为此时没有位置可以记录模块依赖要求，以确保构建结果确定一致。
      </td>
    </tr>
    <tr>
      <td><code>go get</code></td>
      <td>
        在支持此行为的旧版 Go 中，包和可执行文件可以照常构建和安装。注意，运行 <code>go get</code> 时如果没有 <code>go.mod</code> 文件，不存在主模块，因此不会应用 <code>replace</code> 和 <code>exclude</code> 指令。
      </td>
    </tr>
    <tr>
      <td><code>go list -m</code></td>
      <td>
        大多数参数都需要显式的<a href="#version-queries">版本查询</a>，除非使用 <code>-versions</code> 标志。
      </td>
    </tr>
    <tr>
      <td><code>go mod download</code></td>
      <td>
        大多数参数都需要显式的<a href="#version-queries">版本查询</a>。
      </td>
    </tr>
    <tr>
      <td><code>go mod edit</code></td>
      <td>必须显式指定文件参数。</td>
    </tr>
    <tr>
      <td>
        <code>go mod graph</code><br>
        <code>go mod tidy</code><br>
        <code>go mod vendor</code><br>
        <code>go mod verify</code><br>
        <code>go mod why</code>
      </td>
      <td>
        这些命令需要 <code>go.mod</code> 文件；如果文件不存在，就会报错。
      </td>
    </tr>
  </tbody>
</table>

### `go work init` {#go-work-init}

用法：

```
go work init [moddirs]
```

init 在当前目录中初始化并写入新的 go.work 文件，从而在当前目录中创建一个新工作区。

go work init 可以接受工作区模块的路径作为可选参数。如果省略参数，则创建一个不包含模块的空工作区。

每个参数路径都会添加到 go.work 文件中的 use 指令中。go.work 文件还会记录当前 Go 版本。

### `go work edit` {#go-work-edit}

用法：

```
go work edit [editing flags] [go.work]
```

`go work edit` 命令提供编辑 `go.work` 的命令行接口，主要供工具或脚本使用。它只读取 `go.work`，不会查询涉及的模块的信息。如果未指定文件，edit 会在当前目录及其父目录中查找 `go.work` 文件。

编辑标志指定一系列编辑操作。
* `-fmt` 标志会重新格式化 go.work 文件，而不进行其他修改。任何使用或重写 `go.work` 文件的修改操作也都会进行格式化。只有在未指定其他标志时才需要此标志，例如 'go work edit `-fmt`'。
* `-use=path` 和 `-dropuse=path` 标志通过添加和删除 use 指令，修改 `go.work` 文件中的模块目录集合。
* `-replace=old[@v]=new[@v]` 标志为给定的模块路径和版本对添加替换。如果省略 `old@v` 中的 `@v`，则添加左侧不含版本的替换，对旧模块路径的所有版本生效。如果省略 `new@v` 中的 `@v`，则新路径应为本地模块根目录，而非模块路径。注意，`-replace` 会覆盖针对 `old[@v]` 的冗余替换，因此省略 `@v` 会移除已有的针对特定版本的替换。
* `-dropreplace=old[@v]` 标志删除给定模块路径和版本对的替换。如果省略 `@v`，则删除左侧不含版本的替换。
* `-go=version` 标志设置预期的 Go 语言版本。

编辑标志可以重复使用。各项更改按给定顺序应用。

`go work edit` 还提供用于控制输出的标志：

* -print 标志以文本格式输出最终的 go.work，而不将其写回 go.work。
* -json 标志以 JSON 格式输出最终的 go.work 文件，而不将其写回 go.work。JSON 输出对应以下 Go 类型：

```
type Module struct {
    Path    string
    Version string
}

type GoWork struct {
    Go        string
    Directory []Directory
    Replace   []Replace
}

type Use struct {
    Path       string
    ModulePath string
}

type Replace struct {
    Old Module
    New Module
}
```

### `go work use` {#go-work-use}

用法：

```
go work use [-r] [moddirs]
```

`go work use` 命令提供将目录添加到 `go.work` 文件的命令行接口，也可以递归添加目录。

对于命令行中指定的每个目录，如果目录在磁盘上存在，就在 `go.work` 文件中为其添加一条 [`use` 指令](#go-work-file-use)。如果目录不存在，则从 `go.work` 文件中移除对应指令，从而更新 `go.work` 记录的模块目录集合。

`-r` 标志会递归搜索参数目录中的模块，use 命令将每个找到的目录都视为指定的参数进行处理。

### `go work sync` {#go-work-sync}

用法：

```
go work sync
```

`go work sync` 命令将工作区的构建列表同步回工作区中的各个模块。

工作区的构建列表是在工作区中执行构建时使用的所有依赖模块（包括传递依赖）的版本集合。`go
work sync` 使用[最小版本选择（MVS）](#glos-minimal-version-selection)算法生成该构建列表，然后将这些版本同步回工作区中通过 `use` 指令指定的各个模块。

计算出工作区构建列表后，工作区中每个模块的 `go.mod` 文件都会被重写，将与该模块相关的依赖升级到工作区构建列表中的版本。注意，[最小版本选择](#glos-minimal-version-selection)保证构建列表中各模块的版本始终不低于任何工作区模块中指定的相应版本。

## 模块代理 {#module-proxy}

### `GOPROXY` 协议 {#goproxy-protocol}

<dfn>模块代理</dfn>是一种 HTTP 服务器，能够响应针对下列路径的 `GET` 请求。这些请求不带查询参数，也不要求特定请求头，因此即使是提供固定文件系统内容的网站（包括 `file://` URL），也可以充当模块代理。

成功的 HTTP 响应必须使用状态码 200（OK）。客户端会跟随重定向（3xx）。状态码为 4xx 和 5xx 的响应被视为错误。错误码 404（Not Found）和 410（Gone）表示代理上没有请求的模块或版本，但其他位置可能存在。错误响应的内容类型应为 `text/plain`，`charset` 应为 `utf-8` 或 `us-ascii`。

可以通过 `GOPROXY` 环境变量配置 `go` 命令访问的代理或源码版本控制服务器；该变量接受一组代理 URL。列表中可以包含关键字 `direct` 或 `off`（详见[环境变量](#environment-variables)）。列表元素可用逗号（`,`）或竖线（`|`）分隔，分隔符决定发生错误时的回退行为。URL 后为逗号时，`go` 命令只有收到 404（Not Found）或 410（Gone）响应，才会尝试后续来源。URL 后为竖线时，`go` 命令遇到任何错误都会尝试后续来源，包括超时等非 HTTP 错误。代理可以利用这一错误处理机制控制对未知模块的访问。例如，代理可以对不在允许列表中的模块返回 403（Forbidden）错误（参见[提供私有模块的私有代理](#private-module-proxy-private)）。

下表列出模块代理必须响应的查询。对于每条路径，`$base` 为代理 URL 的路径部分，`$module` 为模块路径，`$version` 为版本。例如，代理 URL 为 `https://example.com/mod`，客户端要请求模块 `golang.org/x/text` 的 `v0.3.2` 版本的 `go.mod` 文件时，会向 `https://example.com/mod/golang.org/x/text/@v/v0.3.2.mod` 发送 `GET` 请求。

为了避免在不区分大小写的文件系统上提供内容时产生歧义，`$module` 和 `$version` 元素会进行大小写编码：将每个大写字母替换为一个感叹号及其对应的小写字母。这样，模块 `example.com/M` 和 `example.com/m` 就能同时保存在磁盘上，因为前者会编码为 `example.com/!m`。

<table class="ModTable">
  <thead>
    <tr>
      <th>路径</th>
      <th>说明</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>$base/$module/@v/list</code></td>
      <td>
        以纯文本返回给定模块的已知版本列表，每行一个版本。该列表不应包含伪版本。
      </td>
    </tr>
    <tr>
      <td><code>$base/$module/@v/$version.info</code></td>
      <td>
        <p>
          返回模块特定版本的 JSON 格式元数据。响应必须是与以下 Go 数据结构对应的 JSON 对象：
        </p>
        <pre>
type Info struct {
    Version string    // version string
    Time    time.Time // commit time
}
</pre>
        <p>
          <code>Version</code> 字段为必填项，且必须包含有效的<a href="#glos-canonical-version">规范版本</a>（参见<a href="#versions">版本</a>）。请求路径中的 <code>$version</code> 不必与此版本相同，甚至不必是有效版本；此端点可以用于查找分支名或修订标识符对应的版本。但是，如果 <code>$version</code> 是规范版本，且其主版本号与 <code>$module</code> 兼容，那么成功响应中的 <code>Version</code> 字段必须与其相同。
        </p>
        <p>
          <code>Time</code> 字段为可选项。如果存在，则必须是 RFC 3339 格式的字符串，表示该版本的创建时间。
        </p>
        <p>
          未来可能添加更多字段，因此其他名称均予以保留。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>$base/$module/@v/$version.mod</code></td>
      <td>
        返回模块特定版本的 <code>go.mod</code> 文件。如果模块在请求的版本中没有 <code>go.mod</code> 文件，则必须返回一个仅包含 <code>module</code> 语句的文件，该语句使用所请求的模块路径。否则，必须返回未经修改的原始 <code>go.mod</code> 文件。
      </td>
    </tr>
    <tr>
      <td><code>$base/$module/@v/$version.zip</code></td>
      <td>
        返回包含模块特定版本内容的 zip 文件。该 zip 文件的格式要求详见<a href="#zip-files">模块 zip 文件</a>。
      </td>
    </tr>
    <tr>
      <td><code>$base/$module/@latest</code></td>
      <td>
        返回模块最新已知版本的 JSON 格式元数据，格式与 <code>$base/$module/@v/$version.info</code> 相同。此最新版本应为 <code>go</code> 命令在 <code>$base/$module/@v/list</code> 为空或所列版本都不合适时应使用的版本。此端点是可选的，模块代理不必实现它。
      </td>
    </tr>
  </tbody>
</table>

解析模块的最新版本时，`go` 命令会先请求 `$base/$module/@v/list`；如果找不到合适的版本，再请求 `$base/$module/@latest`。`go` 命令依次优先选择：语义版本最高的正式版本、语义版本最高的预发布版本，以及时间上最新的伪版本。在 Go 1.12 及更早版本中，`go` 命令将 `$base/$module/@v/list` 中的伪版本视为预发布版本；从 Go 1.13 开始不再如此。

对于 `$base/$module/$version.mod` 和 `$base/$module/$version.zip` 查询的成功响应，模块代理必须始终提供相同的内容。这些内容会利用 [`go.sum` 文件](#go-sum-files)进行[密码学验证](#authenticating)，默认情况下还会使用[校验和数据库](#checksum-database)。

`go` 命令将从模块代理下载的大部分内容缓存在 `$GOPATH/pkg/mod/cache/download` 中。即使直接从版本控制系统下载，`go` 命令也会生成明确的 `info`、`mod` 和 `zip` 文件，并将它们存储在此目录中，与从代理直接下载时相同。缓存布局与代理的 URL 路径空间一致，因此只要在 `https://example.com/proxy` 上提供 `$GOPATH/pkg/mod/cache/download` 的内容（或将其复制到该位置），用户就能通过将 `GOPROXY` 设置为 `https://example.com/proxy` 来访问缓存的模块版本。

### 与代理通信 {#communicating-with-proxies}

`go` 命令可以从[模块代理](#glos-module-proxy)下载模块源码和元数据。`GOPROXY` [环境变量](#environment-variables)可用于配置 `go` 命令可以连接哪些代理，以及是否可以直接与[版本控制系统](#vcs)通信。下载的模块数据保存在[模块缓存](#glos-module-cache)中。只有在需要缓存中尚不存在的信息时，`go` 命令才会访问代理。

[`GOPROXY` 协议](#goproxy-protocol)一节介绍了可以向 `GOPROXY` 服务器发送哪些请求。了解 `go` 命令何时发出这些请求也很有帮助。例如，`go build` 按以下流程执行：

* 读取 [`go.mod` 文件](#glos-go-mod-file)，并执行[最小版本选择（MVS）](#glos-minimal-version-selection)，以计算[构建列表](#glos-build-list)。
* 读取命令行中指定的包及其导入的包。
* 如果构建列表中的任何模块都不提供某个包，则查找提供该包的模块。在 `go.mod` 中添加对该模块最新版本的依赖要求，然后重新开始上述流程。
* 全部加载完成后，构建各个包。

`go` 命令计算构建列表时，会加载[模块图](#glos-module-graph)中每个模块的 `go.mod` 文件。如果缓存中不存在某个 `go.mod` 文件，`go` 命令会通过 `$module/@v/$version.mod` 请求从代理下载它（其中 `$module` 为模块路径，`$version` 为版本）。这些请求可以使用 `curl` 等工具进行测试。例如，以下命令下载 `golang.org/x/mod` 的 `v0.2.0` 版本的 `go.mod` 文件：

```
$ curl https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.mod
module golang.org/x/mod

go 1.12

require (
    golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
    golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e
    golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898
)
```

要加载一个包，`go` 命令需要提供该包的模块的源码。模块源码以 `.zip` 文件分发，并解压到模块缓存中。如果缓存中没有模块的 `.zip` 文件，`go` 命令会通过 `$module/@v/$version.zip` 请求下载它。

```
$ curl -O https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.zip
$ unzip -l v0.2.0.zip | head
Archive:  v0.2.0.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
     1479  00-00-1980 00:00   golang.org/x/mod@v0.2.0/LICENSE
     1303  00-00-1980 00:00   golang.org/x/mod@v0.2.0/PATENTS
      559  00-00-1980 00:00   golang.org/x/mod@v0.2.0/README
       21  00-00-1980 00:00   golang.org/x/mod@v0.2.0/codereview.cfg
      214  00-00-1980 00:00   golang.org/x/mod@v0.2.0/go.mod
     1476  00-00-1980 00:00   golang.org/x/mod@v0.2.0/go.sum
     5224  00-00-1980 00:00   golang.org/x/mod@v0.2.0/gosumcheck/main.go
```

注意，即使 `go.mod` 文件通常包含在 `.zip` 文件中，`.mod` 和 `.zip` 请求仍然是分开的。`go` 命令可能需要下载许多不同模块的 `go.mod` 文件，而 `.mod` 文件远小于 `.zip` 文件。此外，如果某个 Go 项目没有 `go.mod` 文件，代理会提供一个仅包含 [`module` 指令](#go-mod-file-module)的合成 `go.mod` 文件。从[版本控制系统](#vcs)下载时，合成的 `go.mod` 文件由 `go` 命令生成。

如果 `go` 命令需要加载的包不由构建列表中的任何模块提供，它会尝试查找提供该包的新模块。[将包解析到模块](#resolve-pkg-mod)一节介绍了这一过程。概括而言，`go` 命令会请求所有可能包含该包的模块路径的最新版本信息。例如，对于包 `golang.org/x/net/html`，`go` 命令会尝试查找模块 `golang.org/x/net/html`、`golang.org/x/net`、`golang.org/x/` 和 `golang.org` 的最新版本。实际上只有 `golang.org/x/net` 存在并提供该包，因此 `go` 命令使用该模块的最新版本。如果多个模块都提供该包，`go` 命令会使用路径最长的模块。

`go` 命令请求模块的最新版本时，会先发送 `$module/@v/list` 请求。如果列表为空，或返回的版本都不可用，则发送 `$module/@latest` 请求。选定版本后，`go` 命令发送 `$module/@v/$version.info` 请求以获取元数据，随后可能发送 `$module/@v/$version.mod` 和 `$module/@v/$version.zip` 请求，以加载 `go.mod` 文件和源码。

```
$ curl https://proxy.golang.org/golang.org/x/mod/@v/list
v0.1.0
v0.2.0

$ curl https://proxy.golang.org/golang.org/x/mod/@v/v0.2.0.info
{"Version":"v0.2.0","Time":"2020-01-02T17:33:45Z"}
```

下载 `.mod` 或 `.zip` 文件后，`go` 命令计算其密码学哈希值，并检查它是否与主模块 `go.sum` 文件中的哈希值匹配。如果 `go.sum` 中没有该哈希值，默认情况下，`go` 命令会从[校验和数据库](#checksum-database)获取它。如果计算出的哈希值不匹配，`go` 命令会报告安全错误，且不将该文件安装到模块缓存中。`GOPRIVATE` 和 `GONOSUMDB` [环境变量](#environment-variables)可用于禁止针对特定模块访问校验和数据库。也可以将 `GOSUMDB` 环境变量设置为 `off`，以完全禁止访问校验和数据库。更多信息参见[验证模块](#authenticating)。注意，版本列表以及 `.info` 请求返回的版本元数据不经过验证，可能随时间变化。

### 直接通过代理提供模块 {#serving-from-proxy}

大多数模块在版本控制仓库中开发，并通过该仓库提供。在[直连模式](#glos-direct-mode)下，`go` 命令使用版本控制工具下载这类模块（参见[版本控制系统](#vcs)）。也可以直接通过模块代理提供模块。对于希望提供模块却不暴露版本控制服务器的组织，以及使用 `go` 命令不支持的版本控制工具的组织，这种方式很有用。

`go` 命令在直连模式下下载模块时，首先根据模块路径发出 HTTP GET 请求，以查找模块服务器的 URL。它会在 HTML 响应中查找名称为 `go-import` 的 `<meta>` 标签。标签内容必须包含[仓库根路径](#glos-repository-root-path)、版本控制系统和 URL，三者以空格分隔。详见[查找模块路径对应的仓库](#vcs-find)。

如果版本控制系统为 `mod`，则 `go` 命令会通过 [`GOPROXY` 协议](#goproxy-protocol)从给定 URL 下载模块。

例如，假设 `go` 命令要下载模块 `example.com/gopher` 的 `v1.0.0` 版本。它向 `https://example.com/gopher?go-get=1` 发送请求，服务器返回包含以下标签的 HTML 文档：

```
<meta name="go-import" content="example.com/gopher mod https://modproxy.example.com">
```

根据此响应，`go` 命令通过请求 `https://modproxy.example.com/example.com/gopher/@v/v1.0.0.info`、`v1.0.0.mod` 和 `v1.0.0.zip` 来下载模块。

注意，在 GOPATH 模式下，无法使用 `go get` 下载直接由代理提供的模块。

## 版本控制系统 {#vcs}

`go` 命令可以直接从版本控制仓库下载模块源码和元数据。从[代理](#communicating-with-proxies)下载模块通常更快，但在没有可用代理，或代理无法访问模块仓库时（私有仓库经常如此），必须直接连接仓库。目前支持 Git、Subversion、Mercurial、Bazaar 和 Fossil。版本控制工具必须安装在 `PATH` 中的某个目录下，`go` 命令才能使用它。

要从源码仓库而非代理下载特定模块，请设置 `GOPRIVATE` 或 `GONOPROXY` 环境变量。要配置 `go` 命令直接从源码仓库下载所有模块，请将 `GOPROXY` 设置为 `direct`。更多信息参见[环境变量](#environment-variables)。

### 查找模块路径对应的仓库 {#vcs-find}

`go` 命令在 `direct` 模式下下载模块时，首先定位包含该模块的仓库。

如果模块路径的某个路径分量以 VCS 限定符（`.bzr`、`.fossil`、`.git`、`.hg` 或 `.svn`）结尾，`go` 命令将使用该限定符之前的路径作为仓库 URL。例如，对于模块 `example.com/foo.git/bar`，`go` 命令使用 git 下载位于 `example.com/foo` 的仓库，并期望在 `bar` 子目录中找到模块。`go` 命令会根据版本控制工具支持的协议，推测应使用的协议。

如果模块路径不含限定符，`go` 命令会根据模块路径构造 URL，附加 `?go-get=1` 查询字符串，并向其发送 HTTP `GET` 请求。例如，对于模块 `golang.org/x/mod`，`go` 命令可能发送以下请求：

```
https://golang.org/x/mod?go-get=1 (preferred)
http://golang.org/x/mod?go-get=1  (fallback, only with GOINSECURE)
```

`go` 命令会跟随重定向，但忽略其他响应状态码，因此服务器可以返回 404 或其他错误状态。可以设置 `GOINSECURE` 环境变量，允许特定模块回退或重定向到未加密的 HTTP。

服务器必须返回 HTML 文档，包含的 `<meta>` 标签应位于文档的 `<head>` 中。为避免 `go` 命令功能有限的解析器产生歧义，`<meta>` 标签应尽早出现在文档中，尤其应位于任何内嵌 JavaScript 或 CSS 之前。`<meta>` 标签必须采用以下形式：

```
<meta name="go-import" content="root-path vcs repo-url [subdirectory]">
```

`root-path` 是仓库根路径，即模块路径中与仓库根目录对应的部分；如果指定了 `subdirectory` 且使用 Go 1.25 或更新版本，则对应于该子目录（参见下文对 `subdirectory` 的说明）。它必须是请求的模块路径的前缀，或与其完全一致。如果不完全一致，还会向此前缀发起一次请求，以验证两次响应的 `<meta>` 标签相符。

`vcs` 是版本控制系统。它必须是下表列出的工具之一，或关键字 `mod`；后者指示 `go` 命令使用 [`GOPROXY` 协议](#goproxy-protocol)从给定 URL 下载模块。详见[直接通过代理提供模块](#serving-from-proxy)。

`repo-url` 是仓库的 URL，必须包含协议方案，且不包含 .vcs 限定符。只有模块路径匹配 `GOINSECURE` 环境变量时，才能使用不安全的协议（例如 `http://` 和 `git://`）。

如果存在，`subdirectory` 表示 `root-path` 对应的仓库子目录，以斜杠分隔，用来覆盖默认的仓库根目录。只有 Go 1.25 及更新版本才识别提供 `subdirectory` 的 `go-import` 元标签。较早版本的 Go 尝试获取和解析模块时会忽略此元标签；如果无法从其他位置解析模块，就会解析失败。

<table id="vcs-support" class="ModTable">
  <thead>
    <tr>
      <th>名称</th>
      <th>命令</th>
      <th>GOVCS 默认设置</th>
      <th>安全协议方案</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Bazaar</td>
      <td><code>bzr</code></td>
      <td>仅私有模块</td>
      <td><code>https</code>, <code>bzr+ssh</code></td>
    </tr>
    <tr>
      <td>Fossil</td>
      <td><code>fossil</code></td>
      <td>仅私有模块</td>
      <td><code>https</code></td>
    </tr>
    <tr>
      <td>Git</td>
      <td><code>git</code></td>
      <td>公开和私有模块</td>
      <td><code>https</code>, <code>git+ssh</code>, <code>ssh</code></td>
    </tr>
    <tr>
      <td>Mercurial</td>
      <td><code>hg</code></td>
      <td>公开和私有模块</td>
      <td><code>https</code>, <code>ssh</code></td>
    </tr>
    <tr>
      <td>Subversion</td>
      <td><code>svn</code></td>
      <td>仅私有模块</td>
      <td><code>https</code>, <code>svn+ssh</code></td>
    </tr>
  </tbody>
</table>

仍以 `golang.org/x/mod` 为例。`go` 命令向 `https://golang.org/x/mod?go-get=1` 发送请求，服务器返回包含以下标签的 HTML 文档：

```
<meta name="go-import" content="golang.org/x/mod git https://go.googlesource.com/mod">
```

根据此响应，`go` 命令将使用远程 URL 为 `https://go.googlesource.com/mod` 的 Git 仓库。

GitHub 和其他常见托管服务会响应所有仓库的 `?go-get=1` 查询，因此托管在这些网站上的模块通常不需要进行服务器配置。

找到仓库 URL 后，`go` 命令会将仓库克隆到模块缓存中。一般来说，`go` 命令会尽量避免从仓库获取不需要的数据。不过，实际使用的命令因版本控制系统而异，也可能随时间变化。对于 Git，`go` 命令无需下载提交即可列出大多数可用版本。通常它会只获取目标提交而不下载祖先提交，但有时仍有必要下载祖先提交。

### 将版本映射到提交 {#vcs-version}

`go` 命令可以检出仓库中某个模块的特定[规范版本](#glos-canonical-version)，例如 `v1.2.3`、`v2.4.0-beta` 或 `v3.0.0+incompatible`。每个模块版本都应在仓库中有一个<dfn>语义版本标签</dfn>，用来指明该版本应检出的修订。

如果模块定义在仓库根目录中，或定义在根目录下的主版本子目录中，那么每个版本的标签名就等于对应版本。例如，模块 `golang.org/x/text` 定义在其仓库的根目录中，因此版本 `v0.3.2` 在该仓库中的标签为 `v0.3.2`。大多数模块都属于这种情况。

如果模块定义在仓库内的子目录中，即模块路径中的[模块子目录](#glos-module-subdirectory)部分非空，那么每个标签名都必须以模块子目录加一个斜杠作为前缀。例如，模块 `golang.org/x/tools/gopls` 定义在根路径为 `golang.org/x/tools` 的仓库的 `gopls` 子目录中。该模块的 `v0.4.0` 版本在仓库中的标签名必须为 `gopls/v0.4.0`。

语义版本标签的主版本号必须与模块路径中的主版本后缀（如果存在）一致。例如，标签 `v1.0.0` 可以属于模块 `example.com/mod`，但不能属于 `example.com/mod/v2`，后者的标签应类似于 `v2.0.0`。

如果模块没有 `go.mod` 文件且位于仓库根目录中，那么主版本为 `v2` 或更高的标签可以属于没有主版本后缀的模块。这类版本以 `+incompatible` 后缀表示，但版本标签本身不能带此后缀。参见[与非模块仓库的兼容性](#non-module-compat)。

标签一旦创建，就不应删除，也不应改为指向其他修订。版本会经过[验证](#authenticating)，以保证构建安全、可重复。如果标签被修改，客户端下载时可能遇到安全错误。即使删除了标签，其内容仍可能保存在[模块代理](#glos-module-proxy)上。

### 将伪版本映射到提交 {#vcs-pseudo}

`go` 命令可以检出仓库中某个模块的特定修订，该修订以[伪版本](#glos-pseudo-version)编码，例如 `v1.3.2-0.20191109021931-daa7c04131f5`。

伪版本的最后 12 个字符（上例中的 `daa7c04131f5`）指明要检出的仓库修订。其含义取决于版本控制系统。对于 Git 和 Mercurial，它是提交哈希值的前缀；对于 Subversion，它是前面补零的修订号。

检出提交前，`go` 命令会验证时间戳（上例中的 `20191109021931`）与提交时间相符。它还会验证基础版本（上例中 `v1.3.2` 之前的版本 `v1.3.1`）对应的语义版本标签位于该提交的祖先提交上。这些检查确保模块作者可以完全控制伪版本与其他已发布版本之间的比较结果。

更多信息参见[伪版本](#pseudo-versions)。

### 将分支和提交映射到版本 {#vcs-branch}

可以使用[版本查询](#version-queries)检出模块的特定分支、标签或修订。

```
go get example.com/mod@master
```

`go` 命令将这些名称转换为可用于[最小版本选择（MVS）](#minimal-version-selection)的[规范版本](#glos-canonical-version)。MVS 依赖于对版本进行无歧义排序的能力。分支名称和修订取决于可能变化的仓库结构，因此无法在时间推移后仍可靠地进行比较。

如果某个修订有一个或多个类似 `v1.2.3` 的语义版本标签，则使用其中最高的有效版本标签。`go` 命令只考虑可能属于目标模块的语义版本标签；例如，对于 `example.com/mod/v2`，不会考虑标签 `v1.5.2`，因为其主版本号与模块路径的后缀不匹配。

如果某个修订没有有效的语义版本标签，`go` 命令会生成一个[伪版本](#glos-pseudo-version)。如果该修订的祖先提交具有有效的语义版本标签，则以祖先提交中最高的版本作为伪版本的基础版本。参见[伪版本](#pseudo-versions)。

### 仓库中的模块目录 {#vcs-dir}

将模块仓库检出到特定修订后，`go` 命令必须定位包含模块 `go.mod` 文件的目录，也就是模块根目录。

回顾一下，[模块路径](#module-path)由三部分组成：仓库根路径（对应仓库根目录）、模块子目录和主版本后缀（仅用于发布版本为 `v2` 或更高的模块）。

对于大多数模块，模块路径与仓库根路径相同，因此模块根目录就是仓库根目录。

模块有时定义在仓库子目录中。这通常用于大型仓库：其中有多个组件，需要各自独立发布和管理版本。这类模块应位于子目录中，子目录路径与模块路径中仓库根路径之后的部分相同。例如，假设模块 `example.com/monorepo/foo/bar` 位于根路径为 `example.com/monorepo` 的仓库中，那么其 `go.mod` 文件必须位于 `foo/bar` 子目录。

如果模块发布的主版本为 `v2` 或更高，其路径必须带有[主版本后缀](#major-version-suffixes)。带主版本后缀的模块可以定义在两种子目录之一：带该后缀的目录，或不带该后缀的目录。例如，假设上述模块发布新版本，其路径为 `example.com/monorepo/foo/bar/v2`，那么它的 `go.mod` 文件可以位于 `foo/bar`，也可以位于 `foo/bar/v2`。

带有主版本后缀的子目录称为<dfn>主版本子目录</dfn>。它们可用于在同一分支上开发模块的多个主版本。如果多个主版本分别在不同分支上开发，可能就不需要这种目录。不过，主版本子目录有一项重要特性：在 `GOPATH` 模式下，包的导入路径与 `GOPATH/src` 下的目录完全匹配。`go` 命令在 `GOPATH` 模式下提供了最低限度的模块兼容性（参见[与非模块仓库的兼容性](#non-module-compat)），因此，要与在 `GOPATH` 模式下构建的项目兼容，并非总是需要主版本子目录。但不支持这种最低限度模块兼容性的旧工具仍可能遇到问题。

`go` 命令找到模块根目录后，会将目录内容打包成 `.zip` 文件，然后将该 `.zip` 文件解压到模块缓存。关于哪些文件可以包含在 `.zip` 文件中，详见[文件路径和大小限制](#zip-path-size-constraints)。将 `.zip` 文件解压到模块缓存前，会对其内容进行[验证](#authenticating)，验证方式与从代理下载 `.zip` 文件时相同。

模块 zip 文件不包含 `vendor` 目录的内容，也不包含任何嵌套模块（含有 `go.mod` 文件的子目录）。这意味着模块必须避免引用自身目录之外或其他模块中的文件。例如，[`//go:embed`](https://pkg.go.dev/embed#hdr-Directives) 模式不能匹配嵌套模块中的文件。如果有些文件不应包含在模块中，可以利用这一行为变通处理。例如，如果仓库在 `testdata` 目录中提交了大文件，模块作者可以在 `testdata` 中添加一个空的 `go.mod` 文件，使用户无需下载这些文件。当然，这可能降低用户测试其依赖时的覆盖范围。

### LICENSE 文件的特殊处理 {#vcs-license}

当 `go` 命令为不在仓库根目录中的模块创建 `.zip` 文件时，如果模块根目录中（即 `go.mod` 旁边）没有名为 `LICENSE` 的文件，`go` 命令会复制同一修订中仓库根目录下的 `LICENSE` 文件（如果存在）。

这种特殊处理使同一个 `LICENSE` 文件可以适用于仓库中的所有模块。它仅适用于恰好名为 `LICENSE`、不带 `.txt` 等扩展名的文件。遗憾的是，如果扩展这一规则，就会破坏已有模块的密码学校验和；参见[验证模块](#authenticating)。其他工具及 [pkg.go.dev](https://pkg.go.dev) 等网站可能会识别其他名称的文件。

还要注意，`go` 命令创建模块 `.zip` 文件时不会包含符号链接，参见[文件路径和大小限制](#zip-path-size-constraints)。因此，如果仓库根目录中没有 `LICENSE` 文件，作者可以在定义于子目录中的模块里放置许可证文件的副本，以确保这些文件包含在模块 `.zip` 文件中。

### 使用 `GOVCS` 控制版本控制工具 {#vcs-govcs}

`go` 命令能够通过 `git` 等版本控制命令下载模块，这对于允许从任意服务器导入代码的去中心化包生态至关重要。但如果恶意服务器找到方法，让被调用的版本控制命令运行非预期代码，这种能力也会带来潜在的安全问题。

为了兼顾功能和安全，默认情况下，`go` 命令只使用 `git` 和 `hg` 从公开服务器下载代码。从私有服务器下载代码时，则可以使用任何[已知的版本控制系统](#vcs-support)；这里的私有服务器指托管的包匹配 `GOPRIVATE` [环境变量](#environment-variables)的服务器。仅允许 Git 和 Mercurial 的原因是，这两个系统在作为客户端访问不可信服务器时的安全问题得到了最充分的关注。相比之下，Bazaar、Fossil 和 Subversion 主要用于可信、经过身份验证的环境，针对攻击面的审查还不够充分。

这些版本控制命令限制只在直接通过版本控制系统下载代码时适用。从代理下载模块时，`go` 命令改用始终允许的 [`GOPROXY` 协议](#goproxy-protocol)。默认情况下，`go` 命令通过 Go 模块镜像（[proxy.golang.org](https://proxy.golang.org)）获取公开模块，仅在获取私有模块，或镜像拒绝提供某个公开包时（通常是出于法律原因），才回退到版本控制系统。因此，默认情况下，客户端仍然能够访问 Bazaar、Fossil 或 Subversion 仓库提供的公开代码，因为这些下载通过 Go 模块镜像进行；镜像使用专门的沙箱承担运行版本控制命令的安全风险。

可以使用 `GOVCS` 变量更改特定模块允许使用的版本控制系统。无论在模块感知模式还是 GOPATH 模式下构建包，`GOVCS` 变量都生效。使用模块时，模式与模块路径匹配；使用 GOPATH 时，模式与版本控制仓库根目录对应的导入路径匹配。

`GOVCS` 变量的一般形式为以逗号分隔的 `pattern:vcslist` 规则列表。pattern 是一个[通配符模式](/pkg/path#Match)，必须匹配模块路径或导入路径开头的一个或多个分量。vcslist 是以竖线分隔的允许使用的版本控制命令列表，也可以是允许任意已知命令的 `all`，或禁用全部命令的 `off`。注意，即使模块匹配的模式对应的 vcslist 为 `off`，只要源服务器使用 `mod` 方案，该模块仍然可以下载，因为这一方案指示 go 命令使用 [`GOPROXY` 协议](#goproxy-protocol)下载模块。列表中最先匹配的模式生效，即使后面的模式也能匹配。

例如，考虑以下设置：

```
GOVCS=github.com:git,evil.com:off,*:git|hg
```

采用此设置后，模块路径或导入路径以 `github.com/` 开头的代码只能使用 `git`；`evil.com` 上的路径不能使用任何版本控制命令；所有其他路径（`*` 匹配所有内容）只能使用 `git` 或 `hg`。

特殊模式 `public` 和 `private` 分别匹配公开和私有的模块路径或导入路径。如果路径匹配 `GOPRIVATE` 变量，则视为私有路径；否则为公开路径。

如果 `GOVCS` 变量中的任何规则都不匹配某个模块路径或导入路径，`go` 命令就会应用默认规则。用 `GOVCS` 记法可以将默认规则表示为 `public:git|hg,private:all`。

要允许任意包不受限制地使用任意版本控制系统，请使用：

```
GOVCS=*:all
```

要完全禁用版本控制，请使用：

```
GOVCS=*:off
```

可以使用 [`go env -w` 命令](/cmd/go/#hdr-Print_Go_environment_information)设置 `GOVCS` 变量，使其对后续调用的 go 命令生效。

`GOVCS` 在 Go 1.16 中引入。较早版本的 Go 可能对任意模块使用任意已知的版本控制工具。

## 模块 zip 文件 {#zip-files}

模块版本以 `.zip` 文件分发。通常无需直接操作这些文件，因为 `go` 命令会自动根据[模块代理](#glos-module-proxy)和版本控制仓库创建、下载和解压它们。不过，在理解跨平台兼容性限制或实现模块代理时，了解这些文件仍然很有帮助。

[`go mod download`](#go-mod-download) 命令下载一个或多个模块的 zip 文件，然后将其解压到[模块缓存](#glos-module-cache)中。根据 `GOPROXY` 及其他[环境变量](#environment-variables)的设置，`go` 命令可能从代理下载 zip 文件，也可能克隆源码版本控制仓库并据此创建 zip 文件。可以使用 `-json` 标志查找下载的 zip 文件及其解压内容在模块缓存中的位置。

可以使用 [`golang.org/x/mod/zip`](https://pkg.go.dev/golang.org/x/mod/zip?tab=doc) 包，通过程序创建、解压或检查 zip 文件的内容。

### 文件路径和大小限制 {#zip-path-size-constraints}

模块 zip 文件的内容受到多项限制。这些限制保证 zip 文件能够在各种平台上安全、一致地解压。

* 模块 zip 文件的大小最多为 500 MiB，其中所有文件解压后的总大小也不得超过 500 MiB。`go.mod` 文件不得超过 16 MiB，`LICENSE` 文件也不得超过 16 MiB。这些限制用于减轻针对用户、代理及模块生态其他部分的拒绝服务攻击。如果仓库的模块目录树中文件总量超过 500 MiB，应在仅包含构建模块中各个包所需文件的提交上标记模块版本；视频、模型和其他大型资源通常不是构建必需的。
* 模块 zip 文件中每个文件的路径都必须以 `$module@$version/` 为前缀，其中 `$module` 是模块路径，`$version` 是版本，例如 `golang.org/x/mod@v0.3.0/`。模块路径必须有效，版本必须有效且为规范形式，版本还必须与模块路径的主版本后缀匹配。具体定义和限制参见[模块路径和版本](#go-mod-file-ident)。
* 文件模式、时间戳和其他元数据会被忽略。
* 模块 zip 文件可以包含空目录（路径以斜杠结尾的条目），但不会解压这些目录。`go` 命令创建的 zip 文件不包含空目录。
* 创建 zip 文件时会忽略符号链接及其他非常规文件，因为它们无法在不同操作系统和文件系统之间移植，zip 格式中也没有可移植的方式来表示它们。
* 创建 zip 文件时会忽略名为 `vendor` 的目录内的文件，因为主模块之外的 `vendor` 目录永远不会被使用。
* 创建 zip 文件时，会忽略模块根目录以外、包含 `go.mod` 文件的目录内的文件，因为这些文件不属于当前模块。解压 zip 文件时，`go` 命令也会忽略包含 `go.mod` 文件的子目录。
* zip 文件中不能存在两个路径在 Unicode 大小写折叠后相同的文件（参见 [`strings.EqualFold`](https://pkg.go.dev/strings?tab=doc#EqualFold)）。这保证 zip 文件可以在不区分大小写的文件系统中解压而不发生冲突。
* 顶层目录中可以包含 `go.mod` 文件（`$module@$version/go.mod`），也可以不包含。如果存在，文件名必须是 `go.mod`（全部小写）。其他任何目录中都不允许存在名为 `go.mod` 的文件。
* 模块内的文件名和目录名可以由 Unicode 字母、ASCII 数字、ASCII 空格字符（U+0020），以及 ASCII 标点字符 `!#$%&()+,-.=@[]^_{}~` 组成。注意，包路径不能包含上述所有字符。有关区别，参见 [`module.CheckFilePath`](https://pkg.go.dev/golang.org/x/mod/module?tab=doc#CheckFilePath) 和 [`module.CheckImportPath`](https://pkg.go.dev/golang.org/x/mod/module?tab=doc#CheckImportPath)。
* 文件名或目录名在第一个点号之前的部分，不得是不区分大小写的 Windows 保留文件名（例如 `CON`、`com1`、`NuL` 等）。

## 私有模块 {#private-modules}

Go 模块经常在无法从公共互联网访问的版本控制服务器和模块代理上开发和分发。`go` 命令可以从私有来源下载和构建模块，不过通常需要进行一些配置。

可以通过以下环境变量配置对私有模块的访问。详见[环境变量](#environment-variables)。关于如何控制发送给公开服务器的信息，另请参阅[隐私](#private-module-privacy)。

* `GOPROXY`：模块代理 URL 列表。`go` 命令依次尝试从各个服务器下载模块。关键字 `direct` 指示 `go` 命令直接从开发模块的版本控制仓库下载，而不使用代理。
* `GOPRIVATE`：模块路径前缀的通配符模式列表，匹配的模块视为私有模块。它也是 `GONOPROXY` 和 `GONOSUMDB` 的默认值。
* `GONOPROXY`：模块路径前缀的通配符模式列表，匹配的模块不应从代理下载。不论 `GOPROXY` 如何设置，`go` 命令都会从开发这些模块的版本控制仓库下载。
* `GONOSUMDB`：模块路径前缀的通配符模式列表，匹配的模块不应通过公共校验和数据库 [sum.golang.org](https://sum.golang.org) 检查。
* `GOINSECURE`：模块路径前缀的通配符模式列表，匹配的模块可以通过 HTTP 及其他不安全的协议获取。

这些变量可以在开发环境中设置（例如在 `.profile` 文件中），也可以使用 [`go env
-w`](/cmd/go/#hdr-Print_Go_environment_information) 持久化设置。

本节后续内容介绍提供私有模块代理和版本控制仓库访问权限的常见方式。

### 提供全部模块的私有代理 {#private-module-proxy-all}

通过集中式私有代理服务器提供所有模块（公开和私有），可以让管理员获得最多的控制能力，同时让开发者个人所需的配置最少。

要配置 `go` 命令使用这样的服务器，请设置以下环境变量，将 `https://proxy.corp.example.com` 替换为代理 URL，将 `corp.example.com` 替换为模块前缀：

```
GOPROXY=https://proxy.corp.example.com
GONOSUMDB=corp.example.com
```

`GOPROXY` 设置指示 `go` 命令仅从 `https://proxy.corp.example.com` 下载模块；`go` 命令不会连接其他代理或版本控制仓库。

`GONOSUMDB` 设置指示 `go` 命令不使用公共校验和数据库验证路径以 `corp.example.com` 开头的模块。

采用此配置的代理通常需要私有版本控制服务器的读取权限。它还需要访问公共互联网，以下载公开模块的新版本。

已有多种 `GOPROXY` 服务器实现可用于此场景。最简单的实现可以提供[模块缓存](#glos-module-cache)目录中的文件，并使用配置适当的 [`go mod
download`](#go-mod-download) 获取缺失的模块。

### 提供私有模块的私有代理 {#private-module-proxy-private}

私有代理服务器可以只提供私有模块，不提供公开模块。可以配置 `go` 命令：对于私有服务器上没有的模块，回退到公开来源。

要让 `go` 命令以这种方式工作，请设置以下环境变量，将 `https://proxy.corp.example.com` 替换为代理 URL，将 `corp.example.com` 替换为模块前缀：

```
GOPROXY=https://proxy.corp.example.com,https://proxy.golang.org,direct
GONOSUMDB=corp.example.com
```

`GOPROXY` 设置指示 `go` 命令首先尝试从 `https://proxy.corp.example.com` 下载模块。如果该服务器返回 404（Not Found）或 410（Gone），`go` 命令会回退到 `https://proxy.golang.org`，然后再回退到直接连接仓库。

`GONOSUMDB` 设置指示 `go` 命令不使用公共校验和数据库验证路径以 `corp.example.com` 开头的模块。

注意，在此配置下，即使代理不提供公开模块，它仍能控制对公开模块的访问。如果代理返回的错误状态不是 404 或 410，`go` 命令就不会回退到 `GOPROXY` 列表中的后续条目。例如，对于许可证不合要求或存在已知安全漏洞的模块，代理可以返回 403（Forbidden）。

### 直接访问私有模块 {#private-module-proxy-direct}

可以配置 `go` 命令绕过公共代理，直接从版本控制服务器下载私有模块。在无法运行私有代理服务器时，这种方式很有用。

要让 `go` 命令以这种方式工作，请设置 `GOPRIVATE`，将 `corp.example.com` 替换为私有模块前缀：

```
GOPRIVATE=corp.example.com
```

在这种情况下无需修改 `GOPROXY` 变量。它的默认值为 `https://proxy.golang.org,direct`，指示 `go` 命令先尝试从 `https://proxy.golang.org` 下载模块；如果该代理返回 404（Not Found）或 410（Gone），再回退到直接连接。

`GOPRIVATE` 设置指示 `go` 命令，对于路径以 `corp.example.com` 开头的模块，不连接代理或校验和数据库。

可能仍需要一台内部 HTTP 服务器来[将模块路径解析为仓库 URL](#vcs-find)。例如，`go` 命令下载模块 `corp.example.com/mod` 时，会向 `https://corp.example.com/mod?go-get=1` 发送 GET 请求，并在响应中查找仓库 URL。要避免这一步，请确保每个私有模块路径都包含标明仓库根前缀的 VCS 后缀（例如 `.git`）。例如，`go` 命令下载模块 `corp.example.com/repo.git/mod` 时，会直接克隆位于 `https://corp.example.com/repo.git` 或 `ssh://corp.example.com/repo.git` 的 Git 仓库，无需发送额外请求。

开发者需要包含私有模块的仓库的读取权限。这可以通过 `.gitconfig` 等全局 VCS 配置文件进行设置。最好将 VCS 工具配置为无需交互式身份验证提示。默认情况下，`go` 命令调用 Git 时会设置 `GIT_TERMINAL_PROMPT=0`，以禁用交互式提示，但会遵循用户显式指定的设置。

### 向私有代理传递凭据 {#private-module-proxy-auth}

`go` 命令可以使用 HTTP [基本身份验证](https://en.wikipedia.org/wiki/Basic_access_authentication)或 `GOAUTH` 提供的请求头，向代理服务器进行身份验证。

可以在 [`.netrc` 文件](https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html)中指定凭据。例如，包含以下内容的 `.netrc` 文件会配置 `go` 命令使用给定的用户名和密码连接主机 `proxy.corp.example.com`。

```
machine proxy.corp.example.com
login jrgopher
password hunter2
```

文件位置可以通过 `NETRC` 环境变量设置。如果未设置 `NETRC`，`go` 命令会在类 UNIX 平台上读取 `$HOME/.netrc`，在 Windows 上读取 `%USERPROFILE%\_netrc`。

`.netrc` 中的字段以空格、制表符和换行符分隔，因此用户名或密码不能包含这些字符。还要注意，主机名不能是完整 URL，所以无法为同一主机上的不同路径指定不同的用户名和密码。

在 Go 1.24 及更新版本中，`GOAUTH` 可以使用 `.netrc`、在指定工作目录中运行 `git credential fill`，或调用自定义命令。自定义命令可以为一个或多个 HTTPS URL 前缀提供任意 HTTP 请求头，例如 Bearer 令牌。这允许将凭据限定到代理主机上的特定路径，而 `.netrc` 无法做到这一点。有关命令协议和配置选项，参见 `go help goauth`。

`go` 命令仅将通过 `GOAUTH`（包括其默认的 `netrc` 方法）获得的凭据用于 HTTPS 请求。

也可以直接在 `GOPROXY` URL 中指定凭据。例如：

```
GOPROXY=https://jrgopher:hunter2@proxy.corp.example.com
```

使用这种方式时应当谨慎：环境变量可能出现在 shell 历史记录和日志中。

### 向私有仓库传递凭据 {#private-module-repo-auth}

`go` 命令可以直接从版本控制仓库下载模块。如果不使用私有代理，就必须通过这种方式获取私有模块。配置方法参见[直接访问私有模块](#private-module-proxy-direct)。

直接下载模块时，`go` 命令会运行 `git` 等版本控制工具。这些工具各自执行身份验证，因此你可能需要在 `.gitconfig` 等工具专用配置文件中配置凭据。

为确保顺利运行，请确认 `go` 命令使用正确的仓库 URL，且版本控制工具不要求交互式输入密码。除非[查找仓库 URL](#vcs-find)时已指定协议方案，否则 `go` 命令优先使用 `https://` URL，而非 `ssh://` 等其他方案。对于 GitHub 仓库，`go` 命令默认使用 `https://`。

<!-- TODO(golang.org/issue/26134): if this issue is fixed, we can remove the
mention of the special case for GitHub above. -->

对于大多数服务器，可以配置客户端通过 HTTP 进行身份验证。例如，GitHub 支持[将 OAuth 个人访问令牌用作 HTTP 密码](https://docs.github.com/en/free-pro-team@latest/github/extending-github/git-automation-with-oauth-tokens)。可以将 HTTP 密码保存在 `.netrc` 文件中，方法与[向私有代理传递凭据](#private-module-proxy-auth)相同。

也可以将 `https://` URL 重写为其他协议方案。例如，在 `.gitconfig` 中配置：

```
[url "git@github.com:"]
    insteadOf = https://github.com/
```

更多信息参见[为什么“go get”克隆仓库时使用 HTTPS？](/doc/faq#git_https)

### 隐私 {#private-module-privacy}

`go` 命令可以从模块代理服务器和版本控制系统下载模块及元数据。环境变量 `GOPROXY` 控制使用哪些服务器；环境变量 `GOPRIVATE` 和 `GONOPROXY` 控制哪些模块从代理获取。

`GOPROXY` 的默认值为：

```
https://proxy.golang.org,direct
```

采用此设置时，`go` 命令下载模块或模块元数据，会先向 Google 运营的公共模块代理 `proxy.golang.org` 发送请求（[隐私政策](https://proxy.golang.org/privacy)）。每个请求发送的信息详见 [`GOPROXY` 协议](#goproxy-protocol)。`go` 命令不传输个人身份信息，但会传输请求的完整模块路径。如果代理返回 404（Not Found）或 410（Gone）状态，`go` 命令会尝试直接连接提供该模块的版本控制系统。详见[版本控制系统](#vcs)。

可以将 `GOPRIVATE` 或 `GONOPROXY` 环境变量设为通配符模式列表，以匹配不应向任何代理请求的私有模块前缀。例如：

```
GOPRIVATE=*.corp.example.com,*.research.example.com
```

`GOPRIVATE` 仅用作 `GONOPROXY` 和 `GONOSUMDB` 的默认值，因此除非 `GONOSUMDB` 需要使用不同的值，否则无需设置 `GONOPROXY`。模块路径匹配 `GONOPROXY` 时，`go` 命令会对该模块忽略 `GOPROXY`，直接从其版本控制仓库获取模块。在没有代理提供私有模块时，这种方式很有用。参见[直接访问私有模块](#private-module-proxy-direct)。

如果有[提供全部模块的可信代理](#private-module-proxy-all)，就不应设置 `GONOPROXY`。例如，将 `GOPROXY` 设置为单一来源时，`go` 命令不会从其他来源下载模块。这种情况下仍应设置 `GONOSUMDB`。

```
GOPROXY=https://proxy.corp.example.com
GONOSUMDB=*.corp.example.com,*.research.example.com
```

如果有[仅提供私有模块的可信代理](#private-module-proxy-private)，就不应设置 `GONOPROXY`，但必须确保代理返回正确的状态码。例如，考虑以下配置：

```
GOPROXY=https://proxy.corp.example.com,https://proxy.golang.org
GONOSUMDB=*.corp.example.com,*.research.example.com
```

假设开发者因拼写错误，尝试下载一个不存在的模块：

```
go mod download corp.example.com/secret-product/typo@latest
```

`go` 命令首先向 `proxy.corp.example.com` 请求此模块。如果该代理返回 404（Not Found）或 410（Gone），`go` 命令就会回退到 `proxy.golang.org`，从而在请求 URL 中传输 `secret-product` 路径。如果私有代理返回任何其他错误码，`go` 命令会输出错误，并且不会回退到其他来源。

除代理外，`go` 命令还可能连接校验和数据库，以验证 `go.sum` 中尚未列出的模块的密码学哈希值。`GOSUMDB` 环境变量设置校验和数据库的名称、URL 和公钥。`GOSUMDB` 的默认值为 `sum.golang.org`，即 Google 运营的公共校验和数据库（[隐私政策](https://sum.golang.org/privacy)）。每个请求传输的信息详见[校验和数据库](#checksum-database)。与访问代理时一样，`go` 命令不传输个人身份信息，但会传输请求的完整模块路径，而校验和数据库无法计算非公开模块的校验和。

可以将 `GONOSUMDB` 环境变量设置为模式列表，指明哪些模块是私有模块、不应向校验和数据库请求。`GOPRIVATE` 是 `GONOSUMDB` 和 `GONOPROXY` 的默认值，因此除非 `GONOPROXY` 需要使用不同的值，否则无需设置 `GONOSUMDB`。

代理可以[镜像校验和数据库](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#proxying-a-checksum-database)。如果 `GOPROXY` 中的某个代理提供此功能，`go` 命令就不会直接连接校验和数据库。

可以将 `GOSUMDB` 设置为 `off`，以完全禁用校验和数据库。采用此设置时，除非下载的模块已记录在 `go.sum` 中，否则 `go` 命令不会验证它们。参见[验证模块](#authenticating)。

## 模块缓存 {#module-cache}

<dfn>模块缓存</dfn>是 `go` 命令存储已下载模块文件的目录。模块缓存不同于构建缓存，后者保存编译后的包和其他构建产物。

模块缓存的默认位置为 `$GOPATH/pkg/mod`。要使用其他位置，请设置 `GOMODCACHE` [环境变量](#environment-variables)。

模块缓存没有最大容量限制，`go` 命令也不会自动删除其中的内容。

同一台机器上开发的多个 Go 项目可以共享缓存。不论主模块位于何处，`go` 命令都会使用同一个缓存。多个 `go` 命令实例可以同时安全地访问同一模块缓存。

`go` 命令在缓存中以只读权限创建模块源码文件和目录，以防止模块下载后被意外修改。这有一个不便的副作用：很难用 `rm -rf` 等命令删除缓存。可以改用 [`go clean -modcache`](#go-clean-modcache) 删除缓存。或者，使用 `-modcacherw` 标志时，`go` 命令会以可读写权限创建新目录，但这会增加编辑器、测试和其他程序修改模块缓存中文件的风险。[`go mod
verify`](#go-mod-verify) 命令可用于检测主模块的依赖是否被修改。它扫描各个依赖模块解压后的内容，检查其完整性；`go.sum` 中则记录了预期的校验和。

下表说明模块缓存中大多数文件的用途，省略了一些临时文件（锁文件、临时目录）。对于每条路径，`$module` 表示模块路径，`$version` 表示版本。以斜杠（`/`）结尾的路径是目录。模块路径和版本中的大写字母使用感叹号转义（`Azure` 转义为 `!azure`），以避免在不区分大小写的文件系统上发生冲突。

<table class="ModTable">
  <thead>
    <tr>
      <th>路径</th>
      <th>说明</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>$module@$version/</code></td>
      <td>
        包含模块 <code>.zip</code> 文件解压内容的目录，作为已下载模块的根目录。如果原始模块没有 <code>go.mod</code> 文件，此目录也不会包含该文件。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/</code></td>
      <td>
        包含从模块代理下载的文件，以及从<a href="#vcs">版本控制系统</a>生成的文件。此目录的布局遵循 <a href="#goproxy-protocol"><code>GOPROXY</code> 协议</a>，因此可以通过 HTTP 文件服务器提供该目录，或使用 <code>file://</code> URL 引用它，将其用作代理。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/$module/@v/list</code></td>
      <td>
        已知版本的列表（参见 <a href="#goproxy-protocol"><code>GOPROXY</code> 协议</a>）。它可能随时间变化，因此 <code>go</code> 命令通常会重新获取一份，而不复用此文件。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/$module/@v/$version.info</code></td>
      <td>
        版本的 JSON 元数据（参见 <a href="#goproxy-protocol"><code>GOPROXY</code> 协议</a>）。它可能随时间变化，因此 <code>go</code> 命令通常会重新获取一份，而不复用此文件。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/$module/@v/$version.mod</code></td>
      <td>
        此版本的 <code>go.mod</code> 文件（参见 <a href="#goproxy-protocol"><code>GOPROXY</code> 协议</a>）。如果原始模块没有 <code>go.mod</code> 文件，则此处是一个不含依赖要求的合成文件。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/$module/@v/$version.zip</code></td>
      <td>
        模块内容的 zip 压缩文件（参见 <a href="#goproxy-protocol"><code>GOPROXY</code> 协议</a>和<a href="#zip-files">模块 zip 文件</a>）。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/$module/@v/$version.ziphash</code></td>
      <td>
        <code>.zip</code> 文件内各文件的密码学哈希值。注意，计算的并非 <code>.zip</code> 文件本身的哈希值，因此文件顺序、压缩方式、对齐和元数据都不会影响哈希值。使用模块时，<code>go</code> 命令会验证此哈希值与 <a href="#go-sum-files"><code>go.sum</code></a> 中对应行的值一致。<a href="#go-mod-verify"><code>go mod verify</code></a> 命令会检查模块 <code>.zip</code> 文件和解压目录的哈希值是否与这些文件中记录的值一致。
      </td>
    </tr>
    <tr>
      <td><code>cache/download/sumdb/</code></td>
      <td>
        包含从<a href="#checksum-database">校验和数据库</a>（通常为 <code>sum.golang.org</code>）下载的文件的目录。
      </td>
    </tr>
    <tr>
      <td><code>cache/vcs/</code></td>
      <td>
        包含直接从源仓库获取模块时克隆的版本控制仓库。目录名是根据仓库类型和 URL 生成的、以十六进制编码的哈希值。仓库会尽量减少磁盘占用。例如，克隆 Git 仓库时会使用裸仓库，并尽可能采用浅克隆。
      </td>
    </tr>
  </tbody>
</table>

## 验证模块 {#authenticating}

`go` 命令将模块的 [zip 文件](#zip-files)或 [`go.mod` 文件](#go-mod-file)下载到[模块缓存](#module-cache)时，会计算其密码学哈希值，并与已知值比较，以验证文件自首次下载以来没有发生变化。如果下载的文件哈希值不正确，`go` 命令会报告安全错误。

对于 `go.mod` 文件，`go` 命令根据文件内容计算哈希值。对于模块 zip 文件，`go` 命令按确定的顺序，根据压缩包内各文件的名称和内容计算哈希值。文件顺序、压缩方式、对齐和其他元数据不会影响哈希值。哈希算法的实现细节参见 [`golang.org/x/mod/sumdb/dirhash`](https://pkg.go.dev/golang.org/x/mod/sumdb/dirhash?tab=doc)。

`go` 命令将每个哈希值与主模块的 [`go.sum` 文件](#go-sum-files)中对应的行进行比较。如果哈希值与 `go.sum` 中的值不同，`go` 命令会报告安全错误，并删除下载的文件，不将其加入模块缓存。

如果 `go.sum` 文件不存在，或不包含下载文件的哈希值，`go` 命令可以使用[校验和数据库](#checksum-database)验证哈希值；该数据库为公开模块提供全局哈希值来源。哈希值验证通过后，`go` 命令会将其加入 `go.sum`，并将下载的文件加入模块缓存。如果模块是私有的（匹配 `GOPRIVATE` 或 `GONOSUMDB` 环境变量），或校验和数据库已禁用（通过设置 `GOSUMDB=off`），则 `go` 命令直接接受该哈希值，并将文件加入模块缓存，不进行验证。

模块缓存通常由系统上的所有 Go 项目共享，而每个模块可能有自己的 `go.sum` 文件，其中的哈希值可能不同。为了避免必须信任其他模块，`go` 命令每次访问模块缓存中的文件时，都会使用主模块的 `go.sum` 验证哈希值。计算 zip 文件的哈希值开销较大，因此 `go` 命令会检查保存在 zip 文件旁边的预计算哈希值，而不是重新计算文件的哈希值。可以使用 [`go mod
verify`](#go-mod-verify) 命令，检查 zip 文件及其解压目录自加入模块缓存以来是否被修改。

### go.sum 文件 {#go-sum-files}

模块根目录中可以有一个名为 `go.sum` 的文本文件，与 `go.mod` 文件位于同一目录。`go.sum` 文件包含模块的直接依赖和间接依赖的密码学哈希值。`go` 命令将模块的 `.mod` 或 `.zip` 文件下载到[模块缓存](#module-cache)时，会计算哈希值，并检查它与主模块 `go.sum` 文件中对应的哈希值一致。如果模块没有依赖，或所有依赖都通过 [`replace` 指令](#go-mod-file-replace)替换为本地目录，`go.sum` 可以为空或不存在。

`go.sum` 的每一行都有三个以空格分隔的字段：模块路径、版本（可能以 `/go.mod` 结尾）和哈希值。

* 模块路径是哈希值所属模块的名称。
* 版本是哈希值所属模块的版本。如果版本以 `/go.mod` 结尾，则哈希值仅对应模块的 `go.mod` 文件；否则，哈希值对应模块 `.zip` 文件内的各个文件。
* 哈希值列由算法名称（例如 `h1`）和经 base64 编码的密码学哈希值组成，两者以冒号（`:`）分隔。目前只支持 SHA-256（`h1`）哈希算法。如果将来发现 SHA-256 存在漏洞，会增加对其他算法的支持（命名为 `h2`，以此类推）。

`go.sum` 文件可以包含同一模块多个版本的哈希值。为了执行[最小版本选择](#minimal-version-selection)，`go` 命令可能需要加载某个依赖多个版本的 `go.mod` 文件。`go.sum` 也可能包含已不再需要的模块版本的哈希值（例如升级后留下的版本）。[`go mod tidy`](#go-mod-tidy) 会向 `go.sum` 添加缺失的哈希值，并移除不再需要的哈希值。

### 校验和数据库 {#checksum-database}

校验和数据库为 `go.sum` 中的记录行提供全局来源。`go` 命令可以在多种情况下利用它检测代理或源服务器的不当行为。

校验和数据库保证所有公开模块版本在全局范围内一致、可靠。它使不可信代理也能使用，因为代理无法在不被发现的情况下提供错误代码。它还确保特定版本对应的字节内容不会随时间改变，即使模块作者随后修改了仓库中的标签也是如此。

校验和数据库由 Google 运营的 [sum.golang.org](https://sum.golang.org) 提供。它是 `go.sum` 记录行哈希值的[透明日志](https://research.swtch.com/tlog)（或“Merkle 树”），底层使用 [Trillian](https://github.com/google/trillian)。Merkle 树的主要优点是独立审计者可以验证其未被篡改，因此它比普通数据库更值得信任。

`go` 命令使用最初在[提案：保护公共 Go 模块生态](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#checksum-database)中概述的协议，与校验和数据库交互。

下表列出校验和数据库必须响应的查询。对于每条路径，`$base` 为校验和数据库 URL 的路径部分，`$module` 为模块路径，`$version` 为版本。例如，如果校验和数据库 URL 为 `https://sum.golang.org`，客户端请求模块 `golang.org/x/text` 的 `v0.3.2` 版本的记录，就会向 `https://sum.golang.org/lookup/golang.org/x/text@v0.3.2` 发送 `GET` 请求。

为了避免在不区分大小写的文件系统上提供内容时产生歧义，`$module` 和 `$version` 元素会进行[大小写编码](https://pkg.go.dev/golang.org/x/mod/module#EscapePath)：将每个大写字母替换为感叹号及其对应的小写字母。这样，模块 `example.com/M` 和 `example.com/m` 就能同时保存在磁盘上，因为前者会编码为 `example.com/!m`。

路径中以方括号括起来的部分，例如 `[.p/$W]`，表示可选值。

<table class="ModTable">
  <thead>
    <tr>
      <th>路径</th>
      <th>说明</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>$base/latest</code></td>
      <td>
        返回最新日志的、经过签名和编码的树描述。该签名描述采用 <a href="https://pkg.go.dev/golang.org/x/mod/sumdb/note">note</a> 形式，即由一个或多个服务器密钥签名的文本，可以使用服务器公钥进行验证。树描述提供树的大小，以及该大小下树头的哈希值。此编码的说明参见 <code><a href="https://pkg.go.dev/golang.org/x/mod/sumdb/tlog#FormatTree">
        golang.org/x/mod/sumdb/tlog#FormatTree</a></code>。
      </td>
    </tr>
    <tr>
    <tr>
      <td><code>$base/lookup/$module@$version</code></td>
      <td>
        返回 <code>$module</code> 的 <code>$version</code> 版本对应条目的日志记录编号，随后是该记录的数据（即 <code>$module</code> 的 <code>$version</code> 版本的 <code>go.sum</code> 记录行），以及包含该记录的、经过签名和编码的树描述。
      </td>
    </tr>
    <tr>
    <tr>
      <td><code>$base/tile/$H/$L/$K[.p/$W]</code></td>
      <td>
        返回一个<a href="https://research.swtch.com/tlog#serving_tiles">日志分块（tile）</a>，即组成日志某一区段的一组哈希值。每个分块由二维坐标确定：位于分块层级 <code>$L</code>，从左数第 <code>$K</code> 块，分块高度为 <code>$H</code>。可选的 <code>.p/$W</code> 后缀表示仅包含 <code>$W</code> 个哈希值的不完整日志分块。如果找不到所需的不完整分块，客户端必须回退为获取完整分块。
      </td>
    </tr>
    <tr>
    <tr>
      <td><code>$base/tile/$H/data/$K[.p/$W]</code></td>
      <td>
        返回 <code>/tile/$H/0/$K[.p/$W]</code> 中各叶子哈希值对应的记录数据（路径分量使用字面值 <code>data</code>）。
      </td>
    </tr>
    <tr>
  </tbody>
</table>

如果 `go` 命令查询校验和数据库，第一步是通过 `/lookup` 端点获取记录数据。如果日志中尚未记录该模块版本，校验和数据库会在响应前尝试从源服务器获取它。`/lookup` 数据提供该模块版本的校验和及其在日志中的位置，客户端据此知道应获取哪些分块来完成证明。在向主模块的 `go.sum` 文件添加新的 `go.sum` 记录行之前，`go` 命令会验证“包含性”证明（某条记录确实存在于日志中）和“一致性”证明（树未被篡改）。必须先根据签名的树哈希值验证 `/lookup` 数据，并根据客户端按时间保存的签名树哈希记录验证该签名树哈希值，之后才能使用这些数据。

校验和数据库提供的签名树哈希值和新分块保存在模块缓存中，因此 `go` 命令只需获取缺失的分块。

`go` 命令不必直接连接校验和数据库。它可以通过[镜像校验和数据库](https://go.googlesource.com/proposal/+/master/design/25530-sumdb.md#proxying-a-checksum-database)且支持上述协议的模块代理请求模块校验和。对于阻止向组织外部发送请求的企业私有代理，这尤其有用。

`GOSUMDB` 环境变量指定要使用的校验和数据库名称，还可以指定其公钥和 URL，例如：

```
GOSUMDB="sum.golang.org"
GOSUMDB="sum.golang.org+<publickey>"
GOSUMDB="sum.golang.org+<publickey> https://sum.golang.org"
```

`go` 命令已知 `sum.golang.org` 的公钥，也知道名称 `sum.golang.google.cn`（可在中国大陆访问）连接的是 `sum.golang.org` 校验和数据库。使用任何其他数据库时，必须显式提供公钥。默认 URL 为 `https://` 加数据库名称。

`GOSUMDB` 默认为 `sum.golang.org`，即 Google 运营的 Go 校验和数据库。该服务的隐私政策见 https://sum.golang.org/privacy 。

如果将 `GOSUMDB` 设置为 `off`，或在支持 `-insecure` 标志的旧版 Go 中使用该标志调用 `go get`，则不会查询校验和数据库，并会接受所有未识别的模块，代价是放弃所有模块下载可验证、可重复的安全保证。要对特定模块绕过校验和数据库，更好的方式是使用 `GOPRIVATE` 或 `GONOSUMDB` 环境变量。详见[私有模块](#private-modules)。

可以使用 `go env -w` 命令[设置这些变量](/pkg/cmd/go/#hdr-Print_Go_environment_information)，使其对后续调用的 `go` 命令生效。

## 环境变量 {#environment-variables}

可以使用下列环境变量配置 `go` 命令的模块行为。此处仅列出与模块相关的环境变量。`go` 命令识别的全部环境变量参见 [`go help
environment`](/cmd/go/#hdr-Environment_variables)。

<table class="ModTable">
  <thead>
    <tr>
      <th>变量</th>
      <th>说明</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>GO111MODULE</code></td>
      <td>
        <p>
          控制 <code>go</code> 命令以模块感知模式还是 <code>GOPATH</code> 模式运行。支持三种取值：
        </p>
        <ul>
          <li>
            <code>off</code>：<code>go</code> 命令忽略 <code>go.mod</code> 文件，并以 <code>GOPATH</code> 模式运行。
          </li>
          <li>
            <code>on</code>（或未设置）：<code>go</code> 命令以模块感知模式运行，即使不存在 <code>go.mod</code> 文件。
          </li>
          <li>
            <code>auto</code>：如果当前目录或任意父目录中存在 <code>go.mod</code> 文件，<code>go</code> 命令就以模块感知模式运行。在 Go 1.15 及更早版本中，这是默认值。
          </li>
        </ul>
        <p>
          更多信息参见<a href="#mod-commands">支持模块的命令</a>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOMODCACHE</code></td>
      <td>
        <p>
          <code>go</code> 命令存储下载的模块及相关文件的目录。此目录的结构详见<a href="#module-cache">模块缓存</a>。
        </p>
        <p>
          如果未设置 <code>GOMODCACHE</code>，则默认为 <code>$GOPATH/pkg/mod</code>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOINSECURE</code></td>
      <td>
        <p>
          以逗号分隔的模块路径前缀通配符模式列表（语法与 Go 的 <a href="/pkg/path/#Match"><code>path.Match</code></a> 相同），匹配的模块始终允许通过不安全的方式获取。仅适用于直接获取的依赖。
        </p>
        <p>
          旧版 Go 提供的 <code>-insecure</code> 标志用于 <code>go get</code> 时会禁用模块校验和数据库验证，而 <code>GOINSECURE</code> 不会。要禁用此验证，可以使用 <code>GOPRIVATE</code> 或 <code>GONOSUMDB</code>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GONOPROXY</code></td>
      <td>
        <p>
          以逗号分隔的模块路径前缀通配符模式列表（语法与 Go 的 <a href="/pkg/path/#Match"><code>path.Match</code></a> 相同），匹配的模块始终直接从版本控制仓库获取，而不通过模块代理。
        </p>
        <p>
          如果未设置 <code>GONOPROXY</code>，则默认为 <code>GOPRIVATE</code>。参见<a href="#private-module-privacy">隐私</a>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GONOSUMDB</code></td>
      <td>
        <p>
          以逗号分隔的模块路径前缀通配符模式列表（语法与 Go 的 <a href="/pkg/path/#Match"><code>path.Match</code></a> 相同）；对于匹配的模块，<code>go</code> 命令不应通过校验和数据库验证校验和。
        </p>
        <p>
          如果未设置 <code>GONOSUMDB</code>，则默认为 <code>GOPRIVATE</code>。参见<a href="#private-module-privacy">隐私</a>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOPATH</code></td>
      <td>
        <p>
          在 <code>GOPATH</code> 模式下，<code>GOPATH</code> 变量是一组可能包含 Go 代码的目录。
        </p>
        <p>
          在模块感知模式下，<a href="#glos-module-cache">模块缓存</a>保存在 <code>pkg/mod</code> 子目录中，该子目录位于第一个 <code>GOPATH</code> 目录下。缓存之外的模块源码可以存储在任意目录中。
        </p>
        <p>
          如果未设置 <code>GOPATH</code>，则默认为用户主目录下的 <code>go</code> 子目录。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOPRIVATE</code></td>
      <td>
        以逗号分隔的模块路径前缀通配符模式列表（语法与 Go 的 <a href="/pkg/path/#Match"><code>path.Match</code></a> 相同），匹配的模块应视为私有模块。<code>GOPRIVATE</code> 是 <code>GONOPROXY</code> 和 <code>GONOSUMDB</code> 的默认值。参见<a href="#private-module-privacy">隐私</a>。<code>GOPRIVATE</code> 还决定 <code>GOVCS</code> 是否将模块视为私有模块。
      </td>
    </tr>
    <tr>
      <td><code>GOPROXY</code></td>
      <td>
        <p>
          模块代理 URL 列表，以逗号（<code>,</code>）或竖线（<code>|</code>）分隔。<code>go</code> 命令查询模块信息时，依次访问列表中的各个代理，直到收到成功响应或导致终止的错误。代理可以通过返回 404（Not Found）或 410（Gone）状态，表示该服务器上没有此模块。
        </p>
        <p>
          <code>go</code> 命令在发生错误时的回退行为，由 URL 之间的分隔符决定。如果代理 URL 后为逗号，则 <code>go</code> 命令在遇到 404 或 410 错误后回退到下一个 URL；其他任何错误都视为终止错误。如果代理 URL 后为竖线，则 <code>go</code> 命令遇到任何错误都会回退到下一个来源，包括超时等非 HTTP 错误。
        </p>
        <p>
          <code>GOPROXY</code> URL 可以使用 <code>https</code>、<code>http</code> 或 <code>file</code> 协议方案。如果 URL 未指定协议方案，则默认为 <code>https</code>。模块缓存可以直接用作文件代理：
        </p>
        <pre>GOPROXY=file://$(go env GOMODCACHE)/cache/download</pre>
        <p>可以使用两个关键字代替代理 URL：</p>
        <ul>
          <li>
            <code>off</code>：禁止从任何来源下载模块。
          </li>
          <li>
            <code>direct</code>：直接从版本控制仓库下载，不使用模块代理。
          </li>
        </ul>
        <p>
          <code>GOPROXY</code> 默认为 <code>https://proxy.golang.org,direct</code>。采用此配置时，<code>go</code> 命令首先访问 Google 运营的 Go 模块镜像；如果镜像没有该模块，再回退到直接连接。镜像的隐私政策见 <a href="https://proxy.golang.org/privacy">https://proxy.golang.org/privacy</a>。可以设置 <code>GOPRIVATE</code> 和 <code>GONOPROXY</code> 环境变量，防止通过代理下载特定模块。私有代理的配置信息参见<a href="#private-module-privacy">隐私</a>。
        </p>
        <p>
          有关如何使用代理的更多信息，参见<a href="#module-proxy">模块代理</a>和<a href="#resolve-pkg-mod">将包解析到模块</a>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOSUMDB</code></td>
      <td>
        <p>
          指定要使用的校验和数据库名称，还可以指定其公钥和 URL。例如：
        </p>
        <pre>
GOSUMDB="sum.golang.org"
GOSUMDB="sum.golang.org+&lt;publickey&gt;"
GOSUMDB="sum.golang.org+&lt;publickey&gt; https://sum.golang.org"
</pre>
        <p>
          <code>go</code> 命令已知 <code>sum.golang.org</code> 的公钥，也知道名称 <code>sum.golang.google.cn</code>（可在中国大陆访问）连接的是 <code>sum.golang.org</code> 数据库。使用其他任何数据库时，必须显式提供公钥。默认 URL 为 <code>https://</code> 加数据库名称。
        </p>
        <p>
          <code>GOSUMDB</code> 默认为 <code>sum.golang.org</code>，即 Google 运营的 Go 校验和数据库。该服务的隐私政策见 <a href="https://sum.golang.org/privacy">https://sum.golang.org/privacy</a>。
        <p>
        <p>
          如果将 <code>GOSUMDB</code> 设置为 <code>off</code>，或在旧版 Go 中调用 <code>go get</code> 时使用当时支持的 <code>-insecure</code> 标志，则不会查询校验和数据库，并会接受所有未识别的模块，代价是放弃所有模块下载可验证、可重复的安全保证。要对特定模块绕过校验和数据库，更好的方式是使用 <code>GOPRIVATE</code> 或 <code>GONOSUMDB</code> 环境变量。
        </p>
        <p>
          更多信息参见<a href="#authenticating">验证模块</a>和<a href="#private-module-privacy">隐私</a>。
        </p>
      </td>
    </tr>
    <tr>
      <td><code>GOVCS</code></td>
      <td>
        <p>
          控制 <code>go</code> 命令下载公开模块、私有模块（根据模块路径是否匹配 <code>GOPRIVATE</code> 中的模式区分）或匹配某个通配符模式的其他模块时，可以使用哪些版本控制工具。
        </p>
        <p>
          如果未设置 <code>GOVCS</code>，或模块不匹配 <code>GOVCS</code> 中的任何模式，<code>go</code> 命令可以对公开模块使用 <code>git</code> 和 <code>hg</code>，对私有模块使用任意已知的版本控制工具。具体而言，<code>go</code> 命令的行为等同于将 <code>GOVCS</code> 设置为：
        </p>
        <pre>public:git|hg,private:all</pre>
        <p>
          完整说明参见<a href="#vcs-govcs">使用 <code>GOVCS</code> 控制版本控制工具</a>。
        </p>
      </td>
    </tr>
     <tr>
      <td><code>GOWORK</code></td>
      <td>
       <p>
        <code>GOWORK</code> 环境变量指示 <code>go</code> 命令进入工作区模式，并使用指定的 <a href="#go-work-file"><code>go.work</code> 文件</a>定义工作区。如果将 <code>GOWORK</code> 设置为 <code>off</code>，则禁用工作区模式。可以利用这一设置，以单模块模式运行 <code>go</code> 命令：例如，<code>GOWORK=off go build .</code> 在单模块模式下构建 <code>.</code> 包。如果 <code>GOWORK</code> 为空，<code>go</code> 命令会查找 <code>go.work</code> 文件，具体方式参见<a href="#workspaces">工作区</a>一节。
       </p>
      </td>
    </tr>
  </tbody>
</table>

## 术语表 {#glossary}

<a id="glos-build-constraint"></a>
**构建约束（build constraint）：** 决定编译包时是否使用某个 Go 源文件的条件。构建约束可以通过文件名后缀（例如 `foo_linux_amd64.go`），或构建约束注释（例如 `// +build linux,amd64`）表达。参见[构建约束](/pkg/go/build/#hdr-Build_Constraints)。

<a id="glos-build-list"></a>
**构建列表（build list）：** 执行 `go build`、`go list` 或 `go test` 等构建命令时使用的模块版本列表。构建列表根据[主模块](#glos-main-module)的 [`go.mod` 文件](#glos-go-mod-file)及其传递依赖模块的 `go.mod` 文件，通过[最小版本选择](#glos-minimal-version-selection)确定。构建列表包含[模块图](#glos-module-graph)中所有模块的版本，而不只是与特定命令相关的模块。

<a id="glos-canonical-version"></a>
**规范版本（canonical version）：** 格式正确，且不包含 `+incompatible` 之外的构建元数据后缀的[版本](#glos-version)。例如，`v1.2.3` 是规范版本，而 `v1.2.3+meta` 不是。

<a id="glos-current-module"></a>
**当前模块（current module）：** [主模块](#glos-main-module)的同义词。

<a id="glos-deprecated-module"></a>
**已弃用模块（deprecated module）：** 作者已不再支持的模块（在这一语境下，不同主版本视为不同模块）。已弃用模块在其最新版本的 [`go.mod` 文件](#glos-go-mod-file)中，使用[弃用注释](#go-mod-file-module-deprecation)进行标记。

<a id="glos-direct-dependency"></a>
**直接依赖（direct dependency）：** 路径出现在[主模块](#glos-main-module)的包或测试的 `.go` 源文件中的 [`import` 声明](/ref/spec#import_declarations)里的包，或包含这类包的模块。（对比[间接依赖](#glos-indirect-dependency)。）

<a id="glos-direct-mode"></a>
**直连模式（direct mode）：** 一种[环境变量](#environment-variables)配置，使 `go` 命令直接从[版本控制系统](#vcs)下载模块，而不通过[模块代理](#glos-module-proxy)。`GOPROXY=direct` 对所有模块启用这种行为；`GOPRIVATE` 和 `GONOPROXY` 则对匹配模式列表的模块启用这种行为。

<a id="glos-go-mod-file"></a>
**`go.mod` 文件：** 定义模块路径、依赖要求及其他元数据的文件，位于[模块根目录](#glos-module-root-directory)中。参见 [`go.mod` 文件](#go-mod-file)一节。

<a id="glos-go-work-file"></a>
**`go.work` 文件：** 定义[工作区](#workspaces)中使用哪些模块的文件。参见 [`go.work` 文件](#go-work-file)一节。

<a id="glos-import-path"></a>
**导入路径（import path）：** Go 源文件中用于导入包的字符串，与[包路径](#glos-package-path)同义。

<a id="glos-indirect-dependency"></a>
**间接依赖（indirect dependency）：** 被[主模块](#glos-main-module)中的包或测试传递导入，但其路径没有出现在主模块任何 [`import` 声明](/ref/spec#import_declarations)中的包；或者，出现在[模块图](#glos-module-graph)中，但没有提供任何被主模块直接导入的包的模块。（对比[直接依赖](#glos-direct-dependency)。）

<a id="glos-lazy-module-loading"></a>
**模块延迟加载（lazy module loading）：** Go 1.17 引入的改动。对于指定 `go 1.17` 或更高版本的模块，如果命令不需要[模块图](#glos-module-graph)，就避免加载模块图。参见[模块延迟加载](#lazy-loading)。

<a id="glos-main-module"></a>
**主模块（main module）：** 调用 `go` 命令时所在的模块。主模块由当前目录或某个父目录中的 [`go.mod` 文件](#glos-go-mod-file)定义。参见[模块、包和版本](#modules-overview)。

<a id="glos-major-version"></a>
**主版本号（major version）：** 语义版本中的第一个数字（`v1.2.3` 中的 `1`）。发布包含不兼容变更的版本时，必须递增主版本号，并将次版本号和补丁版本号设为 0。主版本号为 0 的语义版本视为不稳定版本。

<a id="glos-major-version-subdirectory"></a>
**主版本子目录（major version subdirectory）：** 版本控制仓库中与模块的[主版本后缀](#glos-major-version-suffix)一致、可用于定义模块的子目录。例如，[根路径](#glos-repository-root-path)为 `example.com/mod` 的仓库中的模块 `example.com/mod/v2`，可以定义在仓库根目录中，也可以定义在主版本子目录 `v2` 中。参见[仓库中的模块目录](#vcs-dir)。

<a id="glos-major-version-suffix"></a>
**主版本后缀（major version suffix）：** 与主版本号匹配的模块路径后缀，例如 `example.com/mod/v2` 中的 `/v2`。`v2.0.0` 及更高版本必须包含主版本后缀，较早版本则不允许包含。参见[主版本后缀](#major-version-suffixes)一节。

<a id="glos-minimal-version-selection"></a>
**最小版本选择（minimal version selection，MVS）：** 用于确定构建中所有模块所用版本的算法。详见[最小版本选择](#minimal-version-selection)一节。

<a id="glos-minor-version"></a>
**次版本号（minor version）：** 语义版本中的第二个数字（`v1.2.3` 中的 `2`）。发布添加了向后兼容的新功能的版本时，必须递增次版本号，并将补丁版本号设为 0。

<a id="glos-module"></a>
**模块（module）：** 一组共同发布、管理版本和分发的包。

<a id="glos-module-cache"></a>
**模块缓存（module cache）：** 存储已下载模块的本地目录，位于 `GOPATH/pkg/mod`。参见[模块缓存](#module-cache)。

<a id="glos-module-graph"></a>
**模块图（module graph）：** 以[主模块](#glos-main-module)为根、表示模块依赖要求的有向图。图中每个顶点是一个模块；每条边表示 `go.mod` 文件中某条 `require` 语句指定的版本（受主模块 `go.mod` 文件中 `replace` 和 `exclude` 语句的影响）。

<a id="glos-module-graph-pruning"></a>
**模块图剪枝（module graph pruning）：** Go 1.17 引入的改动，通过省略指定 `go
1.17` 或更高版本的模块的传递依赖，缩小模块图。参见[模块图剪枝](#graph-pruning)。

<a id="glos-module-path"></a>
**模块路径（module path）：** 标识模块的路径，同时也是模块内各包导入路径的前缀。例如 `"golang.org/x/net"`。

<a id="glos-module-proxy"></a>
**模块代理（module proxy）：** 实现 [`GOPROXY` 协议](#goproxy-protocol)的 Web 服务器。`go` 命令从模块代理下载版本信息、`go.mod` 文件和模块 zip 文件。

<a id="glos-module-root-directory"></a>
**模块根目录（module root directory）：** 包含用于定义模块的 `go.mod` 文件的目录。

<a id="glos-module-subdirectory"></a>
**模块子目录（module subdirectory）：** [模块路径](#glos-module-path)中[仓库根路径](#glos-repository-root-path)之后的部分，指明定义模块的子目录。模块子目录非空时，还用作[语义版本标签](#glos-semantic-version-tag)的前缀。模块子目录不包含[主版本后缀](#glos-major-version-suffix)（如果存在），即使模块位于[主版本子目录](#glos-major-version-subdirectory)中也是如此。参见[模块路径](#module-path)。

<a id="glos-package"></a>
**包（package）：** 位于同一目录、一起编译的一组源文件。参见 Go 语言规范中的[包](/ref/spec#Packages)一节。

<a id="glos-package-path"></a>
**包路径（package path）：** 唯一标识一个包的路径，由[模块路径](#glos-module-path)与模块内的子目录路径拼接而成。例如，`"golang.org/x/net/html"` 是模块 `"golang.org/x/net"` 中 `"html"` 子目录里的包的包路径。与[导入路径](#glos-import-path)同义。

<a id="glos-patch-version"></a>
**补丁版本号（patch version）：** 语义版本中的第三个数字（`v1.2.3` 中的 `3`）。发布未更改模块公开接口的版本时，必须递增补丁版本号。

<a id="glos-pre-release-version"></a>
**预发布版本（pre-release version）：** 补丁版本号之后紧接连字符及一组以点分隔的标识符的版本，例如 `v1.2.3-beta4`。预发布版本视为不稳定版本，不假定其与其他版本兼容。预发布版本排序在对应的正式版本之前：`v1.2.3-pre` 排在 `v1.2.3` 之前。另见[正式版本](#glos-release-version)。

<a id="glos-pseudo-version"></a>
**伪版本（pseudo-version）：** 编码了版本控制系统中的修订标识符（例如 Git 提交哈希值）和时间戳的版本，例如 `v0.0.0-20191109021931-daa7c04131f5`。用于[兼容非模块仓库](#non-module-compat)，以及其他没有可用标签版本的情况。

<a id="glos-release-version"></a>
**正式版本（release version）：** 不带预发布后缀的版本。例如 `v1.2.3`，而非 `v1.2.3-pre`。另见[预发布版本](#glos-pre-release-version)。

<a id="glos-repository-root-path"></a>
**仓库根路径（repository root path）：** [模块路径](#glos-module-path)中与版本控制仓库根目录对应的部分。参见[模块路径](#module-path)。

<a id="glos-retracted-version"></a>
**已撤回版本（retracted version）：** 不应再被依赖的版本，原因可能是发布过早，或发布后发现严重问题。参见 [`retract` 指令](#go-mod-file-retract)。

<a id="glos-semantic-version-tag"></a>
**语义版本标签（semantic version tag）：** 版本控制仓库中将[版本](#glos-version)映射到特定修订的标签。参见[将版本映射到提交](#vcs-version)。

<a id="glos-selected-version"></a>
**选定版本（selected version）：** [最小版本选择](#minimal-version-selection)为给定模块选出的版本。选定版本是[模块图](#glos-module-graph)中该模块路径出现的最高版本。

<a id="glos-vendor-directory"></a>
**vendor 目录（vendor directory）：** 名为 `vendor` 的目录，包含构建主模块中的包所需的其他模块中的包。使用 [`go mod vendor`](#go-mod-vendor) 维护。参见 [vendor 模式](#vendoring)。

<a id="glos-version"></a>
**版本（version）：** 模块不可变快照的标识符，写作字母 `v` 后跟语义版本。参见[版本](#versions)一节。

<a id="glos-workspace"></a>
**工作区（workspace）：** 磁盘上的一组模块，在执行[最小版本选择（MVS）](#minimal-version-selection)时作为主模块使用。参见[工作区](#workspaces)一节。
