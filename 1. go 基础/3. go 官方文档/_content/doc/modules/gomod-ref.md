<!--{
  "Title": "go.mod 文件参考"
}-->

每个 Go 模块都由一个 go.mod 文件定义。该文件描述模块的属性，包括它对其他模块和 Go 版本的依赖。

这些属性包括：

* 当前模块的**模块路径**。该路径应指向 Go 工具可以下载模块的位置，例如模块代码所在的仓库。模块路径与版本号一起构成模块版本的唯一标识，它也是模块内所有包的包路径前缀。有关 Go 如何定位模块，详见 <a href="/ref/mod#vcs-find">Go 模块参考手册</a>。
* 当前模块所需的最低 **Go 版本**。
* 当前模块**依赖的其他模块**及其最低版本列表。
* 可选指令：用其他模块版本或本地目录**替换（replace）**依赖模块，**排除（exclude）**某个依赖模块的特定版本，或在匹配包模式时**忽略（ignore）**模块中的指定目录。

运行 [`go mod init` 命令](/ref/mod#go-mod-init)时，Go 会生成 go.mod 文件。下面的示例创建一个 go.mod 文件，并将模块路径设为 example/mymodule：

```
$ go mod init example/mymodule
```

使用 `go` 命令管理依赖，可以确保 go.mod 中的依赖要求保持一致，文件内容合法。相关命令包括 [`go get`](/ref/mod#go-get)、[`go mod tidy`](/ref/mod#go-mod-tidy) 和 [`go mod edit`](/ref/mod#go-mod-edit)。

`go` 命令的参考文档见 [go 命令](/cmd/go/)。也可以在命令行输入 `go help` _command-name_ 获取帮助，例如 `go help mod tidy`。

**另请参阅**

* 使用 Go 工具管理依赖时，工具会修改 go.mod 文件。详情请参阅[管理依赖](/doc/modules/managing-dependencies)。
* go.mod 文件的更多细节和约束见 [Go 模块参考手册](/ref/mod#go-mod-file)。

## 示例 {#example}

go.mod 文件包含如下示例中的指令。本文后续章节会逐一介绍。

```
module example.com/mymodule

go 1.14

require (
    example.com/othermodule v1.2.3
    example.com/thismodule v1.2.3
    example.com/thatmodule v1.2.3
)

replace example.com/thatmodule => ../thatmodule
exclude example.com/thismodule v1.3.0
```

## module {#module}

声明模块路径。模块路径与模块版本号组合后，是模块版本的唯一标识。模块路径也是模块中所有包的导入路径前缀。

详情请参阅 Go 模块参考手册中的 [`module` 指令](/ref/mod#go-mod-file-module)。

### 语法 {#module-syntax}

<pre>module <var>module-path</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>模块路径，通常是 Go 工具可以下载模块的仓库位置。对于 v2 及更高版本的模块，路径必须以主版本号结尾，例如 <code>/v2</code>。</dd>
</dl>

### 示例 {#module-examples}

以下示例使用 `example.com` 代指可下载模块的仓库域名。

* v0 或 v1 模块的声明：
  ```
  module example.com/mymodule
  ```
* v2 模块的模块路径：
  ```
  module example.com/mymodule/v2
  ```

### 说明 {#module-notes}

模块路径必须能够唯一标识你的模块。对大多数模块而言，该路径是 `go` 命令可以找到代码的 URL，或者重定向到代码位置的 URL。对于永远不会被直接下载的模块，模块路径也可以只是由你控制、能够确保唯一性的名称。`example/` 前缀保留用于此类示例。

详情请参阅[管理依赖](/doc/modules/managing-dependencies#naming_module)。

实际使用中，模块路径通常由源代码仓库的域名和模块代码在仓库内的路径组成。`go` 命令会依照这种结构下载模块版本，为模块使用者解析依赖。

即使最初没有计划让其他代码使用你的模块，采用仓库路径仍是推荐做法，可以避免日后发布时再重命名模块。

如果一开始还不知道模块最终的仓库位置，可以先使用安全的替代名称，例如你拥有的域名或控制的名称（如公司名），再加上由模块名或源码目录推导出的路径。详情请参阅[管理依赖](/doc/modules/managing-dependencies#naming_module)。

例如，如果在 `stringtools` 目录中开发，临时模块路径可以是 `<company-name>/stringtools`，如下所示，其中 _company-name_ 是你的公司名：

```
go mod init <company-name>/stringtools
```

## go {#go}

表明该模块按照指令指定的 Go 版本语义编写。

详情请参阅 Go 模块参考手册中的 [`go` 指令](/ref/mod#go-mod-file-go)。

### 语法 {#go-syntax}

<pre>go <var>minimum-go-version</var></pre>

<dl>
    <dt>minimum-go-version</dt>
    <dd>编译本模块中的包所需的最低 Go 版本。</dd>
</dl>

### 示例 {#go-examples}

* 模块要求 Go 1.14 或更高版本：
  ```
  go 1.14
  ```

### 说明 {#go-notes}

`go` 指令设置使用该模块所需的最低 Go 版本。在 Go 1.21 之前，该指令只是建议；如今它是一项强制要求：Go 工具链会拒绝使用声明了更新 Go 版本的模块。

`go` 指令也是决定运行哪个 Go 工具链的依据之一。详情请参阅[《Go 工具链》](/doc/toolchain)。

`go` 指令会影响新语言特性的使用：

* 对模块中的包，编译器会拒绝使用在 `go` 指令指定版本之后才引入的语言特性。例如，模块声明了 `go 1.12`，其中的包就不能使用 Go 1.13 引入的 `1_000_000` 这类数值字面量。
* 如果旧版本 Go 构建模块中的某个包时遇到编译错误，错误信息会指出该模块是为更新版本的 Go 编写的。例如，模块声明了 `go 1.13`，包中使用了数值字面量 `1_000_000`；使用 Go 1.12 构建该包时，编译器会提示代码面向 Go 1.13。

`go` 指令还会影响 `go` 命令的行为：

* 声明 `go 1.14` 或更高版本时，可以自动启用 [vendor 模式](/ref/mod#vendoring)。如果存在 `vendor/modules.txt`，且其内容与 `go.mod` 一致，就无需显式指定 `-mod=vendor` 标志。
* 声明 `go 1.16` 或更高版本时，`all` 包模式仅匹配[主模块](/ref/mod#glos-main-module)中的包及其测试直接或间接导入的包。这与模块机制引入以来 [`go mod vendor`](/ref/mod#go-mod-vendor) 所保留的包集合一致。在更早版本中，`all` 还包括主模块所导入包的测试，以及这些测试进一步导入的包的测试，依此类推。
* 声明 `go 1.17` 或更高版本时：
   * 对于主模块中的包或测试直接或间接导入的每一个包，提供该包的模块都会在 `go.mod` 中有显式的 [`require` 指令](/ref/mod#go-mod-file-require)。（在 `go 1.16` 及更早版本中，只有当不记录间接依赖会导致[最小版本选择](/ref/mod#minimal-version-selection)选出不同版本时，才会记录该间接依赖。）这些额外信息使[模块图裁剪](/ref/mod#graph-pruning)和[模块懒加载](/ref/mod#lazy-loading)成为可能。
   * 由于 `// indirect` 依赖可能比先前的 `go` 版本多得多，间接依赖会记录在 `go.mod` 的独立块中。
   * `go mod vendor` 不再复制 vendor 依赖中的 `go.mod` 和 `go.sum` 文件，这样在 `vendor` 子目录内调用 `go` 命令时，才能识别正确的主模块。
   * `go mod vendor` 会将各依赖的 `go.mod` 文件中声明的 `go` 版本记录到 `vendor/modules.txt`。
* 声明 `go 1.21` 或更高版本时：
   * `go` 行声明使用该模块所需的最低 Go 版本。
   * `go` 行的版本必须大于或等于所有依赖的 `go` 行版本。
   * `go` 命令不再尝试维持与前一个 Go 旧版本的兼容性。
   * `go` 命令会更谨慎地在 `go.sum` 中保留 `go.mod` 文件的校验和。
<!-- If you update this list, also update /ref/mod#go-mod-file-go. -->

一个 `go.mod` 文件最多只能包含一条 `go` 指令。如果没有该指令，大多数命令会使用当前 Go 版本补上。

## toolchain {#toolchain}

声明建议与该模块一起使用的 Go 工具链。只有当该模块是主模块，且默认工具链比建议的工具链旧时，才会生效。

详情请参阅[《Go 工具链》](/doc/toolchain)以及 Go 模块参考手册中的 [`toolchain` 指令](/ref/mod/#go-mod-file-toolchain)。

### 语法 {#toolchain-syntax}

<pre>toolchain <var>toolchain-name</var></pre>

<dl>
    <dt>toolchain-name</dt>
    <dd>建议使用的 Go 工具链名称。对于 Go 版本 <i>V</i>，标准工具链名称采用 <code>go<i>V</i></code> 的形式，例如 <code>go1.21.0</code> 和 <code>go1.18rc1</code>。特殊值 <code>default</code> 会禁用自动工具链切换。</dd>
</dl>

### 示例 {#toolchain-examples}

* 建议使用 Go 1.21.0 或更新版本：
    ```
    toolchain go1.21.0
    ```

### 说明 {#toolchain-notes}

`toolchain` 行如何影响 Go 工具链选择，详见[《Go 工具链》](/doc/toolchain)。

## godebug {#godebug}

指定应用于本模块各个 main 包的默认 [GODEBUG](/doc/godebug) 设置。这些设置会覆盖工具链的默认值，而 main 包中显式的 `//go:debug` 行又会覆盖这些设置。

### 语法 {#godebug-syntax}

<pre>godebug <var>debug-key</var>=<var>debug-value</var></pre>

<dl>
    <dt>debug-key</dt>
    <dd>要应用的设置名称。可用设置及其引入版本列表见 <a href="/doc/godebug#history">GODEBUG 历史</a>。</dd>
    <dt>debug-value</dt>
    <dd>该设置的取值。除非另有说明，<code>0</code> 表示禁用相应行为，<code>1</code> 表示启用。</dd>
</dl>

### 示例 {#godebug-examples}

* 使用 Go 1.23 引入的 `asynctimerchan=0` 新行为：
  ```
  godebug asynctimerchan=0
  ```
* 使用 Go 1.21 的默认 GODEBUG 设置，但保留旧的 `panicnil=1` 行为：
  ```
  godebug (
      default=go1.21
      panicnil=1
  )
  ```

### 说明 {#godebug-notes}

GODEBUG 设置仅适用于当前模块中 main 包和测试二进制文件的构建。模块被用作依赖时，这些设置不会生效。

有关向后兼容性的更多信息，请参阅[《Go、向后兼容性与 GODEBUG》](/doc/godebug)。

## require {#require}

将一个模块声明为当前模块的依赖，并指定所需的最低版本。

详情请参阅 Go 模块参考手册中的 [`require` 指令](/ref/mod#go-mod-file-require)。

### 语法 {#require-syntax}

<pre>require <var>module-path</var> <var>module-version</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>模块路径，通常由源码仓库的域名和模块名称拼接而成。对于 v2 及更高版本的模块，该值必须以主版本号结尾，例如 <code>/v2</code>。</dd>
    <dt>module-version</dt>
    <dd>模块版本，可以是 v1.2.3 这样的正式版本号，也可以是 Go 生成的伪版本号，例如 v0.0.0-20200921210052-fa0125251cc4。</dd>
</dl>

### 示例 {#require-examples}

* 依赖已发布的 v1.2.3 版本：
    ```
    require example.com/othermodule v1.2.3
    ```
* 使用 Go 工具生成的伪版本号，依赖仓库中尚未打标签的版本：
    ```
    require example.com/othermodule v0.0.0-20200921210052-fa0125251cc4
    ```

### 说明 {#require-notes}

运行 `go get` 等 `go` 命令时，Go 会为提供已导入包的各个模块加入 `require` 指令。若模块在仓库中尚未打标签，Go 会在命令运行时为其生成并指定伪版本号。

通过 [`replace` 指令](#replace)，可以让 Go 从模块仓库之外的位置获取所需模块。

有关版本号，详见[模块版本编号](/doc/modules/version-numbers)。

有关依赖管理，请参阅：

* [添加依赖](/doc/modules/managing-dependencies#adding_dependency)
* [获取特定版本的依赖](/doc/modules/managing-dependencies#getting_version)
* [查找可用更新](/doc/modules/managing-dependencies#discovering_updates)
* [升级或降级依赖](/doc/modules/managing-dependencies#upgrading)
* [同步代码依赖](/doc/modules/managing-dependencies#synchronizing)

## tool {#tool}

将一个包添加为当前模块的依赖，并使当前工作目录位于本模块内时，可以通过 `go tool` 运行该包。

### 语法 {#tool-syntax}

<pre>tool <var>package-path</var></pre>

<dl>
    <dt>package-path</dt>
    <dd>工具的包路径，由包含工具的模块路径与实现工具的包在模块内的路径拼接而成，后者可以为空。</dd>
</dl>

### 示例 {#tool-examples}

* 声明在当前模块中实现的工具：
    ```
    module example.com/mymodule

    tool example.com/mymodule/cmd/mytool
    ```
* 声明在另一个模块中实现的工具：
    ```
    module example.com/mymodule

    tool example.com/atool/cmd/atool

    require example.com/atool v1.2.3
    ```

### 说明 {#tool-notes}

可以通过完整的包路径，或在没有歧义时使用路径的最后一段，通过 `go tool` 运行模块中声明的工具。在上面的第一个示例中，可以运行 `go tool mytool` 或 `go tool example.com/mymodule/cmd/mytool`。

在工作区模式下，可以使用 `go tool` 运行任意工作区模块中声明的工具。

工具使用与模块本身相同的模块图进行构建。需要通过 [`require` 指令](#require)选择实现该工具的模块版本。[`replace` 指令](#replace)和 [`exclude` 指令](#exclude)也会应用于工具及其依赖。

详情请参阅[工具依赖](/doc/modules/managing-dependencies#tools)。

## replace {#replace}

用其他模块版本或本地目录，替换某个模块特定版本（或所有版本）的内容。Go 工具在解析依赖时会使用替换后的路径。

详情请参阅 Go 模块参考手册中的 [`replace` 指令](/ref/mod#go-mod-file-replace)。

### 语法 {#replace-syntax}

<pre>replace <var>module-path</var> <var>[module-version]</var> => <var>replacement-path</var> <var>[replacement-version]</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>要替换的模块的模块路径。</dd>
    <dt>module-version</dt>
    <dd>可选，要替换的特定版本。省略该版本号时，模块的所有版本都会被箭头右侧的内容替换。</dd>
    <dt>replacement-path</dt>
    <dd>Go 应从中查找所需模块的路径，可以是模块路径，也可以是替代模块在本地文件系统中的目录路径。如果是模块路径，必须指定 <em>replacement-version</em>；如果是本地路径，则不能指定 <em>replacement-version</em>。</dd>
    <dt>replacement-version</dt>
    <dd>替代模块的版本。只有当 <em>replacement-path</em> 是模块路径而不是本地目录时，才可以指定替代版本。</dd>
</dl>

### 示例 {#replace-examples}

* 使用模块仓库的 fork 替换

  以下示例将 example.com/othermodule 的任意版本替换为指定的代码 fork。

  ```
  require example.com/othermodule v1.2.3

  replace example.com/othermodule => example.com/myfork/othermodule v1.2.3-fixed
  ```

  将一个模块路径替换为另一个路径时，不要修改原模块中各个包的导入语句。

  有关使用模块代码 fork 的更多信息，请参阅[从自己的仓库 fork 引入外部模块代码](/doc/modules/managing-dependencies#external_fork)。

* 替换为不同的版本号

  以下示例指定使用 v1.2.3，替代该模块的任何其他版本。

  ```
  require example.com/othermodule v1.2.2

  replace example.com/othermodule => example.com/othermodule v1.2.3
  ```

  以下示例将模块的 v1.2.5 版本替换为同一模块的 v1.2.3 版本。

  ```
  replace example.com/othermodule v1.2.5 => example.com/othermodule v1.2.3
  ```

* 使用本地代码替换

  以下示例指定使用本地目录替换该模块的所有版本。

  ```
  require example.com/othermodule v1.2.3

  replace example.com/othermodule => ../othermodule
  ```

  以下示例指定本地目录仅替换 v1.2.5 版本。

  ```
  require example.com/othermodule v1.2.5

  replace example.com/othermodule v1.2.5 => ../othermodule
  ```

  有关使用模块代码本地副本的更多信息，请参阅[引入本地目录中的模块代码](/doc/modules/managing-dependencies#local_directory)。

### 说明 {#replace-notes}

当希望 Go 从另一个路径查找模块源代码时，可以通过 `replace` 指令暂时将模块路径替换为其他值。这相当于将 Go 对模块的查找重定向到替代位置。无需将包的导入路径改成替代路径。

`exclude` 和 `replace` 指令用于在构建当前模块时控制依赖解析。其他模块依赖当前模块时，当前模块中的这些指令会被忽略。

`replace` 指令适用于以下场景：

* 正在开发新模块，代码还没有放入仓库，希望让调用方使用本地版本进行测试。
* 发现依赖存在问题，已克隆其仓库，正在利用本地仓库测试修复。

注意，仅有 `replace` 指令并不会将模块加入[模块图](/ref/mod#glos-module-graph)。主模块或某个依赖的 `go.mod` 中，还必须有引用被替换模块版本的 [`require` 指令](#require)。如果没有要替换的具体版本，可以使用一个虚构的版本号，如下所示。但这会导致依赖你模块的其他模块无法正常构建，因为 `replace` 指令只在主模块中生效。

```
require example.com/mod v0.0.0-replace

replace example.com/mod v0.0.0-replace => ./mod
```

有关替换依赖模块的更多信息，包括如何使用 Go 工具完成修改，请参阅：

* [从自己的仓库 fork 引入外部模块代码](/doc/modules/managing-dependencies#external_fork)
* [引入本地目录中的模块代码](/doc/modules/managing-dependencies#local_directory)

有关版本号，详见[模块版本编号](/doc/modules/version-numbers)。

## exclude {#exclude}

指定要从当前模块的依赖图中排除的模块及其版本。

详情请参阅 Go 模块参考手册中的 [`exclude` 指令](/ref/mod#go-mod-file-exclude)。

### 语法 {#exclude-syntax}

<pre>exclude <var>module-path</var> <var>module-version</var></pre>

<dl>
    <dt>module-path</dt>
    <dd>要排除的模块的模块路径。</dd>
    <dt>module-version</dt>
    <dd>要排除的具体版本。</dd>
</dl>

### 示例 {#exclude-example}

* 排除 example.com/theirmodule 的 v1.3.0 版本

  ```
  exclude example.com/theirmodule v1.3.0
  ```

### 说明 {#exclude-notes}

如果某个间接依赖模块的特定版本由于某种原因无法加载，可以使用 `exclude` 将其排除。例如，排除校验和无效的模块版本。

`exclude` 和 `replace` 指令用于在构建当前模块（也就是正在构建的主模块）时控制依赖解析。其他模块依赖当前模块时，当前模块中的这些指令会被忽略。

可以使用 [`go mod edit`](/ref/mod#go-mod-edit) 命令排除模块版本，例如：

```
go mod edit -exclude=example.com/theirmodule@v1.3.0
```

有关版本号，详见[模块版本编号](/doc/modules/version-numbers)。

## retract {#retract}

表明不应再依赖此 go.mod 所定义模块的某个版本或某个版本范围。当版本过早发布，或发布后发现严重问题时，可以使用 `retract` 指令撤回版本。

详情请参阅 Go 模块参考手册中的 [`retract` 指令](/ref/mod#go-mod-file-retract)。

### 语法 {#retract-syntax}

<pre>
retract <var>version</var> // <var>rationale</var>
retract [<var>version-low</var>,<var>version-high</var>] // <var>rationale</var>
</pre>

<dl>
  <dt>version</dt>
  <dd>要撤回的单个版本。</dd>
  <dt>version-low</dt>
  <dd>要撤回的版本范围的下界。</dd>
  <dt>version-high</dt>
  <dd>要撤回的版本范围的上界。范围同时包含 <var>version-low</var> 和 <var>version-high</var>。</dd>
  <dt>rationale</dt>
  <dd>可选注释，用于说明撤回原因，可能会显示在面向用户的提示中。</dd>
</dl>

### 示例 {#retract-example}

* 撤回单个版本

  ```
  retract v1.1.0 // Published accidentally.
  ```

* 撤回一个范围内的版本

  ```
  retract [v1.0.0,v1.0.5] // Build broken on some platforms.
  ```

### 说明 {#retract-notes}

通过 `retract` 指令表明不应继续使用模块的某个旧版本。用户运行 `go get`、`go mod tidy` 等命令时，不会自动升级到已撤回的版本；运行 `go list -m -u` 时，也不会将已撤回版本列为可用更新。

已撤回版本仍应保持可获取状态，以便已经依赖它们的用户继续构建包。即使从源码仓库中删除了已撤回版本，它也可能仍保存在 [proxy.golang.org](https://proxy.golang.org) 等镜像中。依赖已撤回版本的用户对相关模块运行 `go get` 或 `go list -m -u` 时，可能收到提示。

`go` 命令通过读取模块最新版本的 `go.mod` 中的 `retract` 指令，发现被撤回的版本。最新版本按以下优先顺序确定：

1. 最高的正式版本（如果存在）。
2. 最高的预发布版本（如果存在）。
3. 仓库默认分支最新提交对应的伪版本。

加入撤回声明时，通常都需要打一个更高版本的新标签，命令才能在模块最新版本中发现它。

也可以发布一个仅用于声明撤回信息的版本。这种情况下，新版本还可以撤回自身。

例如，误打了 `v1.0.0` 标签后，可以打 `v1.0.1` 标签，并包含以下指令：

```
retract v1.0.0 // Published accidentally.
retract v1.0.1 // Contains retraction only.
```

遗憾的是，版本一旦发布，就不能再修改。如果之后在另一个提交上重新打 `v1.0.0` 标签，`go` 命令可能检测到它与 `go.sum` 或[校验和数据库](/ref/mod#checksum-database)中的校验和不匹配。

已撤回的模块版本通常不会出现在 `go list -m -versions` 的输出中，但可以使用 `-retracted` 标志显示它们。详情请参阅 Go 模块参考手册中的 [`go list -m`](/ref/mod#go-list-m)。

## ignore {#ignore}

指定模块内的目录路径，`go` 命令匹配包模式时应忽略这些目录。

详情请参阅 Go 模块参考手册中的 [`ignore` 指令](/ref/mod#go-mod-file-ignore)。

### 语法 {#ignore-syntax}

<pre>
ignore <var>path</var>
</pre>

<dl>
  <dt>path</dt>
  <dd>要忽略的路径，使用斜杠分隔。若路径以 <code>./</code> 开头，则相对于模块根目录解释；否则，模块内任意层级下名称匹配的目录都会被忽略。</dd>
</dl>

### 示例 {#ignore-examples}

* 忽略相对于模块根目录的一个本地目录

  ```
  ignore ./node_modules
  ```

* 忽略模块内任意位置名为 `generated` 的所有目录

  ```
  ignore generated
  ```

* 通过块形式忽略多个路径

  ```
  ignore (
      static
      content/html
      ./third_party/javascript
  )
  ```

### 说明 {#ignore-notes}

使用 `./...` 等通配包模式时，`ignore` 指令可以防止 `go` 命令匹配生成目录或非 Go 目录中的包。被忽略的目录及其内容不会参与包模式匹配，但仍保留在模块文件树中。

`ignore` 指令仅在主模块的 `go.mod` 文件中生效，在依赖模块的 `go.mod` 中无效。
