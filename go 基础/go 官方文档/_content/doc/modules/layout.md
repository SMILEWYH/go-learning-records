<!--{
  "Title": "组织 Go 模块"
}-->

刚接触 Go 的开发者经常会问：“Go 项目的文件和目录应该如何组织？”本文提供一些指导，帮助回答这个问题。阅读之前，建议先通过[教程](/doc/tutorial/create-module)和[管理模块源码](/doc/modules/managing-source)了解 Go 模块的基础知识。

Go 项目可以包含包、命令行程序，或者同时包含两者。本指南按项目类型介绍。

### 基本包 {#basic-package}

一个基本 Go 包的全部代码都放在项目根目录中。项目只有一个模块，模块中只有一个包。包名与模块名路径的最后一部分一致。对于只需要一个 Go 源文件的简单包，项目结构如下：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
```

*［本文中的文件名和包名都只是示例，可以自行选择。］*

假设该目录上传到 GitHub 仓库 `github.com/someuser/modname`，那么 `go.mod` 文件中的 `module` 行应为 `module github.com/someuser/modname`。

`modname.go` 中通过以下方式声明包：

```
package modname

// ... package code here
```

用户可以在 Go 代码中通过 `import` 导入该包：

```
import "github.com/someuser/modname"
```

一个 Go 包可以拆分为多个文件，但这些文件都应位于同一目录中，例如：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth.go
  auth_test.go
  hash.go
  hash_test.go
```

目录中的所有文件都声明 `package modname`。

### 基本命令 {#basic-command}

基本的可执行程序（或命令行工具）应根据复杂度和代码规模组织。最简单的程序可以只有一个定义了 `func main` 的 Go 文件。更大的程序可以拆分为多个文件，所有文件都声明 `package main`：

```
project-root-directory/
  go.mod
  auth.go
  auth_test.go
  client.go
  main.go
```

这里 `main.go` 包含 `func main`，但这只是惯例。这个“主”文件也可以叫 `modname.go`（选择合适的 `modname`），或其他名称。

假设该目录上传到 GitHub 仓库 `github.com/someuser/modname`，那么 `go.mod` 文件中的 `module` 行应为：

```
module github.com/someuser/modname
```

用户应该可以通过以下命令安装程序：

```
$ go install github.com/someuser/modname@latest
```

### 带辅助包的包或命令 {#package-or-command-with-supporting-packages}

较大的包或命令可以将部分功能拆分到辅助包中。开始时，建议把这些包放进名为 `internal` 的目录。[这样可以防止](https://pkg.go.dev/cmd/go#hdr-Internal_Directories)其他模块依赖我们并不打算对外公开和长期支持的包。由于其他项目无法导入 `internal` 目录中的代码，我们可以自由重构 API、调整组织方式，而不会破坏外部用户的代码。包的项目结构如下：

```
project-root-directory/
  internal/
    auth/
      auth.go
      auth_test.go
    hash/
      hash.go
      hash_test.go
  go.mod
  modname.go
  modname_test.go
```

`modname.go` 声明 `package modname`，`auth.go` 声明 `package auth`，以此类推。`modname.go` 可以这样导入 `auth` 包：

```
import "github.com/someuser/modname/internal/auth"
```

命令程序将辅助包放在 `internal` 目录时，结构也很相似，区别是根目录中的文件声明 `package main`。

### 多个包 {#multiple-packages}

一个模块可以包含多个可导入的包，每个包都有自己的目录，目录可以按层级组织。以下是示例结构：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
    token/
      token.go
      token_test.go
  hash/
    hash.go
  internal/
    trace/
      trace.go
```

再次说明，这里假定 `go.mod` 中的 `module` 行为：

```
module github.com/someuser/modname
```

`modname` 包位于根目录，声明 `package modname`，用户可以这样导入：

```
import "github.com/someuser/modname"
```

子目录中的包可以这样导入：

```
import "github.com/someuser/modname/auth"
import "github.com/someuser/modname/auth/token"
import "github.com/someuser/modname/hash"
```

位于 `internal/trace` 的 `trace` 包不能被本模块以外的代码导入。建议尽量将不需要对外公开的包保留在 `internal` 中。

### 多个命令 {#multiple-commands}

同一仓库中的多个程序通常各有独立目录：

```
project-root-directory/
  go.mod
  internal/
    ... shared internal packages
  prog1/
    main.go
  prog2/
    main.go
```

每个程序目录中的 Go 文件都声明 `package main`。顶层的 `internal` 目录可以存放仓库内各个命令共用的包。

用户可以这样安装这些程序：

```
$ go install github.com/someuser/modname/prog1@latest
$ go install github.com/someuser/modname/prog2@latest
```

一种常见惯例是将仓库中的所有命令放进 `cmd` 目录。对只包含命令的仓库来说，这不是硬性要求；但对于同时包含命令和可导入包的混合仓库，这种方式非常实用，下一节会介绍。

### 同一仓库中的包与命令 {#packages-and-commands-in-the-same-repository}

有时，一个仓库既提供可导入的包，也提供功能相关的可安装命令。示例结构如下：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
  internal/
    ... internal packages
  cmd/
    prog1/
      main.go
    prog2/
      main.go
```

假设模块名为 `github.com/someuser/modname`，用户既可以导入包：

```
import "github.com/someuser/modname"
import "github.com/someuser/modname/auth"
```

也可以安装程序：

```
$ go install github.com/someuser/modname/cmd/prog1@latest
$ go install github.com/someuser/modname/cmd/prog2@latest
```

### 服务端项目 {#server-project}

Go 经常用于编写*服务端程序*。服务端开发涉及协议（REST、gRPC 等）、部署、前端文件、容器化、脚本等诸多方面，因此项目结构差异很大。这里重点介绍 Go 代码部分。

服务端项目通常不需要对外提供可导入的包，因为服务端一般是独立的二进制程序，或一组二进制程序。因此，建议把实现服务端逻辑的 Go 包放在 `internal` 目录中。此外，项目可能有很多包含非 Go 文件的目录，将所有 Go 命令集中到 `cmd` 目录也很合适：

```
project-root-directory/
  go.mod
  internal/
    auth/
      ...
    metrics/
      ...
    model/
      ...
  cmd/
    api-server/
      main.go
    metrics-analyzer/
      main.go
    ...
  ... the project's other directories with non-Go code
```

如果服务端仓库后来出现适合与其他项目共享的包，最好将这些包拆分为独立模块。
