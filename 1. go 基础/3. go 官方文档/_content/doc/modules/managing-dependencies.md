<!--{
  "Title": "管理依赖"
}-->

当代码使用外部包时，这些以模块形式分发的包就成为项目的依赖。随着项目演进，你可能需要升级或替换它们。Go 提供依赖管理工具，帮助你在引入外部依赖的同时，维护应用的安全性。

本文介绍如何管理代码中的依赖。大多数操作可以通过 Go 工具完成，也会介绍一些其他实用的依赖管理任务。

**另请参阅**

* 如果你刚开始通过模块使用依赖，请先阅读[入门教程](/doc/tutorial/getting-started)。
* 使用 `go` 命令管理依赖，有助于保持依赖要求一致，并确保 go.mod 文件内容有效。命令参考见 [go 命令](/cmd/go/)。也可以在命令行输入 `go help` *command-name*，例如 `go help mod tidy`。
* 用来修改依赖的 Go 命令会更新 go.mod 文件。其内容说明见 [go.mod 文件参考](/doc/modules/gomod-ref)。
* 让编辑器或 IDE 支持 Go 模块，可以简化管理工作。详情参见[编辑器插件与 IDE](/doc/editors.html)。
* 本文不介绍如何开发、发布模块并管理版本以供他人使用。这些内容参见[开发与发布模块](developing)。

## 使用与管理依赖的工作流程 {#workflow}

借助 Go 工具，你可以获取并使用所需的包。在 [pkg.go.dev](https://pkg.go.dev) 搜索包，然后通过 `go` 命令获取它们，在自己的代码中导入并调用其中的函数。

最常见的依赖管理步骤如下，详细说明见后续章节。

1. 在 [pkg.go.dev](https://pkg.go.dev) 上[查找合适的包](#locating_packages)。
1. 在代码中[导入需要的包](#locating_packages)。
1. 如果代码尚未属于某个模块，将其放入模块以跟踪依赖。参见[启用依赖跟踪](#enable_tracking)。
1. [将外部包添加为依赖](#adding_dependency)，以便管理。
1. 按需[升级或降级依赖版本](#upgrading)。

## 以模块为单位管理依赖 {#modules}

在 Go 中，依赖以模块为单位管理，每个模块包含你导入的包。这一过程由以下机制支持：

* **去中心化的发布体系**：用于发布模块和获取代码。开发者在自己的仓库中提供模块，并使用版本号发布，供其他人使用。
* **包搜索引擎和文档浏览器**（pkg.go.dev）：用于查找模块。参见[查找并导入合适的包](#locating_packages)。
* **模块版本号约定**：帮助你判断模块的稳定性和向后兼容性保证。参见[模块版本号](version-numbers)。
* **Go 工具**：简化模块源码获取、升级等依赖管理操作，详见本文各节。

## 查找并导入合适的包 {#locating_packages}

你可以在 [pkg.go.dev](https://pkg.go.dev) 上搜索提供所需功能的包。

找到想使用的包后，在页面顶部找到包路径，点击 Copy path 按钮复制到剪贴板，再将路径粘贴到代码的 import 语句中，例如：

```
import "rsc.io/quote"
```

导入包之后，需要启用依赖跟踪，并获取编译所需的包代码。参见[在代码中启用依赖跟踪](#enable_tracking)和[添加依赖](#adding_dependency)。

## 在代码中启用依赖跟踪 {#enable_tracking}

要跟踪和管理依赖，首先需要将代码放入自己的模块中。这会在源码树的根目录创建 go.mod 文件，之后添加的依赖会记录在其中。

使用 [`go mod init` 命令](/ref/mod#go-mod-init)创建模块。例如，在命令行进入代码根目录，运行：

```
$ go mod init example/mymodule
```

`go mod init` 的参数是模块路径。如果可能，模块路径应当指向源码仓库所在的位置。

如果一开始还不知道最终的仓库地址，可以使用安全的临时名称，例如你拥有的域名或其他可控制的名称（如公司名称），再加上由模块名或源码目录名构成的路径。参见[为模块命名](#naming_module)。

使用 Go 工具管理依赖时，工具会自动更新 go.mod 文件，使其中的依赖列表保持最新。

添加依赖时，Go 工具还会创建 go.sum 文件，记录依赖模块的校验和。Go 使用这些校验和验证下载文件的完整性，这对参与同一项目的其他开发者尤其重要。

请将 go.mod、go.sum 和代码一起纳入版本控制。

详情参见 [go.mod 参考](/doc/modules/gomod-ref)。

## 为模块命名 {#naming_module}

使用 `go mod init` 创建模块以跟踪依赖时，需要指定作为模块名称的模块路径。该路径也会成为模块内各个包的导入路径前缀。请确保模块路径不会与其他模块冲突。

模块路径至少应体现某些来源信息，例如公司、作者或所有者名称，也可以进一步说明模块的内容或功能。

模块路径通常采用以下形式：

```
<prefix>/<descriptive-text>
```

* *prefix* 通常是能够描述模块来源等信息的字符串，例如：

    * Go 工具查找模块源码的仓库位置。发布模块时需要提供这样的地址。

        例如 `github.com/<project-name>/`。

        如果将来可能发布模块供其他人使用，建议采用这种方式。详情参见[开发与发布模块](/doc/modules/developing)。

    * 你可以控制的名称。

        如果不使用仓库名，应选择能够确保不会被其他人使用的前缀。公司名称是不错的选择。避免 `widgets`、`utilities` 或 `app` 等通用词。

* *descriptive-text* 可以使用项目名。请记住，描述功能主要是包名的职责；模块路径为这些包名提供命名空间。

**保留的模块路径前缀**

Go 保留以下前缀，供下述用途使用，避免与实际发布的包路径冲突。

- `test`：用于在本地测试其他模块功能的模块。

    在测试过程中创建模块时，可以使用 `test` 前缀。例如，测试程序可能运行 `go mod init test`，然后按特定方式组织模块，以测试 Go 源码分析工具。

- `example`：Go 文档中的示例模块路径前缀，例如教程中为了跟踪依赖而创建的模块。

    如果示例描述的是一个可能发布的模块，Go 文档也会使用 `example.com`。

## 添加依赖 {#adding_dependency}

导入已发布模块中的包后，可以使用 [`go get` 命令](/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them)，将模块添加为受管理的依赖。

该命令会：

* 根据需要，在 go.mod 文件中为构建命令行指定的包所需的模块添加 `require` 指令。`require` 指令记录当前模块所依赖的模块及其最低版本。详情参见 [go.mod 参考](/doc/modules/gomod-ref)。
* 根据需要下载模块源码，以便编译依赖这些模块的包。源码可从 proxy.golang.org 等模块代理下载，也可以直接从版本控制仓库下载，并缓存在本地。

    你可以指定 Go 工具下载模块的位置。参见[指定模块代理服务器](#proxy_server)。

下面是一些示例。

* 添加当前模块中某个包的所有依赖，可运行以下命令（“.”表示当前目录中的包）：

    ```
    $ go get .
    ```

* 添加指定依赖时，将模块路径作为命令参数：

    ```
    $ go get example.com/theirmodule
    ```

命令还会验证每个下载的模块，确保其内容与发布时一致。如果发布后的模块内容发生变化，例如开发者修改了对应提交的内容，Go 工具会报告安全错误。这种验证可以防范模块遭到篡改。

## 获取依赖的指定版本 {#getting_version}

在 `go get` 命令中指定版本，可以获取依赖模块的特定版本。命令会更新 go.mod 中的 `require` 指令；你也可以手动更新该指令。

以下情况可能需要指定版本：

* 试用模块的某个预发布版本。
* 发现当前版本无法满足需求，希望换成已知可用的版本。
* 升级或降级已有依赖。

[`go get` 命令](/ref/mod#go-get)用法示例：

* 获取指定版本时，在模块路径后添加 @ 和版本号：

    ```
    $ go get example.com/theirmodule@v1.3.4
    ```

* 获取最新版本时，在模块路径后添加 `@latest`：

    ```
    $ go get example.com/theirmodule@latest
    ```

以下 go.mod 文件中的 `require` 指令展示如何要求特定版本，详情参见 [go.mod 参考](/doc/modules/gomod-ref)：

```
require example.com/theirmodule v1.3.4
```

## 查看可用更新 {#discovering_updates}

你可以检查当前模块所使用的依赖是否有新版本。使用 `go list` 列出依赖，以及各模块可用的最新版本。发现更新后，可以结合自己的代码试用，再决定是否升级。

命令说明参见 [`go list -m`](/ref/mod#go-list-m)。

例如：

* 列出当前模块的所有依赖模块及其可用的最新版本：

    ```
    $ go list -m -u all
    ```

* 查看某个模块可用的最新版本：

    ```
    $ go list -m -u example.com/theirmodule
    ```

## 升级或降级依赖 {#upgrading}

先用 Go 工具查找可用版本，再将所需版本添加为依赖，即可升级或降级模块。

1. 按照[查看可用更新](#discovering_updates)中的说明，使用 `go list` 查找新版本。
1. 按照[获取依赖的指定版本](#getting_version)中的说明，使用 `go get` 选择特定版本。

## 同步代码中的依赖 {#synchronizing}

你可以确保代码导入的所有包都具有相应的依赖记录，同时移除不再导入的包对应的依赖。

当代码和依赖不断变化时，已记录的依赖及下载的模块可能不再与代码实际需要的包一致，此操作尤其有用。

使用 `go mod tidy` 保持依赖集合整洁。该命令根据代码导入的包修改 go.mod，添加缺失的必要模块，并移除不再提供任何相关包的无用模块。

不带参数即可运行；添加 `-v` 可以显示被移除模块的信息。

```
$ go mod tidy
```

## 使用尚未发布的模块代码开发和测试 {#unpublished}

你可以指定使用尚未发布的依赖模块。这些模块的代码可能位于原仓库、仓库的 fork，或与调用方模块处于同一本地磁盘上的目录。

以下情况可能需要这样做：

* 希望修改外部模块的代码，例如 fork 或克隆模块后修复缺陷，再通过拉取请求提交给模块开发者。
* 正在开发新模块，但尚未发布，`go get` 无法从仓库获取它。

### 依赖本地目录中的模块代码 {#local_directory}

可以指定使用与调用方模块位于同一本地磁盘上的模块代码。适用于：

* 正在开发自己的另一个模块，希望从当前模块调用并测试。
* 正在修复外部模块的缺陷或增加功能，希望从当前模块测试。（也可以使用自己 fork 的外部仓库，参见[依赖自己 fork 的仓库中的模块代码](#external_fork)。）

在 go.mod 中使用 `replace` 指令，替换 `require` 指令指定的模块路径，让 Go 命令使用本地副本。指令说明参见 [go.mod 参考](/doc/modules/gomod-ref)。

下面的 go.mod 中，当前模块依赖 `example.com/theirmodule`，并使用不存在的版本号 `v0.0.0-unpublished`，确保替换生效。`replace` 将其替换为 `../theirmodule`，即与当前模块目录同级的目录。

```
module example.com/mymodule

go 1.23.0

require example.com/theirmodule v0.0.0-unpublished

replace example.com/theirmodule v0.0.0-unpublished => ../theirmodule
```

配置 `require`/`replace` 时，使用 [`go mod edit`](/ref/mod#go-mod-edit) 和 [`go get`](/ref/mod#go-get)，确保文件中的依赖要求一致：

```
$ go mod edit -replace=example.com/theirmodule@v0.0.0-unpublished=../theirmodule
$ go get example.com/theirmodule@v0.0.0-unpublished
```

**注意：** 使用 replace 指令时，被替换的外部模块不会按[添加依赖](#adding_dependency)中所述的方式验证原模块内容。

版本号的更多说明，参见[模块版本号](/doc/modules/version-numbers)。

Go 1.18 引入了[工作区模式](/blog/get-familiar-with-workspaces)，可以同时开发多个模块。参见[教程：多模块工作区入门](/doc/tutorial/workspaces)。

### 依赖自己 fork 的仓库中的模块代码 {#external_fork}

fork 外部模块仓库后，例如为了修复缺陷或添加功能，可以让 Go 工具从你的 fork 获取源码，方便通过自己的代码测试变更。（也可以引用本地目录中的模块代码，参见[依赖本地目录中的模块代码](#local_directory)。）

在 go.mod 中使用 `replace` 指令，将外部模块的原路径替换为 fork 仓库中的路径。这样，编译等操作会使用替换后的源码位置，同时代码中的 `import` 语句仍可保留原模块路径。

指令说明参见 [go.mod 文件参考](gomod-ref)。

下例中，当前模块依赖 `example.com/theirmodule`。`replace` 将它替换为原仓库的 fork：`example.com/myfork/theirmodule`。

```
module example.com/mymodule

go 1.23.0

require example.com/theirmodule v1.2.3

replace example.com/theirmodule v1.2.3 => example.com/myfork/theirmodule v1.2.3-fixed
```

配置 `require`/`replace` 时，应使用 Go 命令保持依赖要求一致。先用 [`go list`](/ref/mod#go-list-m) 查询当前使用的版本，再用 [`go mod edit`](/ref/mod#go-mod-edit) 将依赖替换为 fork：

```
$ go list -m example.com/theirmodule
example.com/theirmodule v1.2.3
$ go mod edit -replace=example.com/theirmodule@v1.2.3=example.com/myfork/theirmodule@v1.2.3-fixed
```

**注意：** 使用 `replace` 指令时，被替换的外部模块不会按[添加依赖](#adding_dependency)中所述的方式验证原模块内容。

版本号的更多说明，参见[模块版本号](/doc/modules/version-numbers)。

## 使用仓库标识获取指定提交 {#repo_identifier}

可以使用 `go get`，将模块仓库中尚未发布为正式版本的某次提交作为依赖。

在命令中通过 `@` 指定所需代码。`go get` 会在 go.mod 中添加该模块的 `require` 指令，并根据提交信息生成伪版本号。

下面以 Git 仓库中的模块为例。

* 获取特定提交时，在模块路径后添加 @<em>commithash</em>：

    ```
    $ go get example.com/theirmodule@4cf76c2
    ```

* 获取特定分支时，在模块路径后添加 @<em>branchname</em>：

    ```
    $ go get example.com/theirmodule@bugfixes
    ```

## 移除依赖 {#removing_dependency}

当代码不再使用某个模块中的任何包时，可以停止将其作为依赖跟踪。

运行 [`go mod tidy` 命令](/ref/mod#go-mod-tidy)，移除所有未使用的模块依赖。该命令也可能添加构建当前模块所需但尚未记录的依赖。

```
$ go mod tidy
```

要移除指定依赖，使用 [`go get` 命令](/ref/mod#go-get)，在模块路径后添加 `@none`：

```
$ go get example.com/theirmodule@none
```

`go get` 还会降级或移除依赖于被移除模块的其他依赖。

## 工具依赖 {#tools}

工具依赖用于管理开发当前模块时使用的、用 Go 编写的开发工具。例如，配合 [`go generate`](/blog/generate) 使用的 [`stringer`](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)，或提交代码前运行的特定代码检查工具、格式化工具。

Go 1.24 及以上版本可以这样添加工具依赖：

```
$ go get -tool golang.org/x/tools/cmd/stringer
```

这会向 go.mod 添加 [`tool` 指令](/ref/mod/#go-mod-file-tool)，并确保包含必要的 require 指令。之后可以把工具导入路径中最后一个[非主版本号部分](/ref/mod#major-version-suffixes)传给 `go tool`，运行该工具：

```
$ go tool stringer
```

如果多个工具的最后一段路径相同，或该名称与 Go 发行版自带工具重名，则必须使用完整包路径：

```
$ go tool golang.org/x/tools/cmd/stringer
```

不带参数运行 `go tool`，可以列出当前可用的全部工具：

```
$ go tool
```

也可以手动向 go.mod 添加 `tool` 指令，但必须确保有对应的 `require` 指令，依赖定义工具的模块。补齐缺失 `require` 指令最简单的方式是：

```
$ go mod tidy
```

工具依赖所需的依赖要求，与[模块图](/ref/mod#glos-module-graph)中的其他依赖要求一样，参与[最小版本选择](/ref/mod#minimal-version-selection)，并遵守 `require`、`replace` 和 `exclude` 指令。由于模块图裁剪，当你依赖的模块自身有工具依赖时，仅用于满足该工具依赖的模块通常不会成为你的模块的依赖。

`tool` [元模式](/cmd/go#hdr-Package_lists_and_patterns)可以同时操作所有工具。例如，`go get tool`（等价于 `go get tool@upgrade`）升级所有工具；`go install tool` 将所有工具安装到 $GOBIN。

在 Go 1.24 之前，可以在模块中添加一个通过[构建约束](/pkg/go/build/#hdr-Build_Constraints)排除在正常构建之外的 Go 文件，并在其中空白导入工具包，以实现类似 `tool` 指令的效果。之后可以使用 `go run` 加完整包路径运行工具。

## 指定模块代理服务器 {#proxy_server}

使用 Go 工具处理模块时，默认从 Google 运营的公共模块镜像 proxy.golang.org 下载，或直接从模块仓库下载。也可以指定其他代理服务器，用于模块下载与校验相关操作。

如果你或团队部署或选择了其他模块代理，就可能需要修改这一设置。例如，有些团队通过自建代理更严格地控制依赖的使用方式。

将 `GOPROXY` 环境变量设置为一个或多个服务器 URL，即可指定代理。Go 工具按给定顺序尝试。默认设置先使用 Google 的公共模块代理，再根据模块路径直接从仓库下载：

```
GOPROXY="https://proxy.golang.org,direct"
```

关于 `GOPROXY` 及其其他取值，参见 [`go` 命令参考](/cmd/go/#hdr-Module_downloading_and_verification)。

多个模块代理 URL 可以用逗号或竖线分隔。

* 使用逗号时，只有当前 URL 返回 HTTP 404 或 410，才尝试下一个 URL。

    ```
    GOPROXY="https://proxy.example.com,https://proxy2.example.com"
    ```

* 使用竖线时，无论遇到哪种 HTTP 错误，都会尝试下一个 URL。

    ```
    GOPROXY="https://proxy.example.com|https://proxy2.example.com"
    ```

Go 模块也经常在不对公网开放的版本控制服务器和模块代理上开发、分发。可以设置 `GOPRIVATE` 环境变量，配置 `go` 命令从私有来源下载和构建模块。

`GOPRIVATE` 或 `GONOPROXY` 可以设置为一组通配模式，匹配私有模块路径前缀，使这些模块不会被发送给任何公共代理。例如：

```
GOPRIVATE=*.corp.example.com,*.research.example.com
```
