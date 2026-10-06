<!--{
  "Title": "弃用 go get 安装可执行文件的功能",
  "Path": "/doc/go-get-install-deprecation",
  "Breadcrumb": true
}-->

## 概述 {#overview}

从 Go 1.17 开始，使用 `go get` 安装可执行文件的方式已被弃用，请改用 `go install`。

从 Go 1.18 开始，`go get` 不再构建包，只用于添加、更新或移除 `go.mod` 中的依赖。具体而言，`go get` 的行为始终等同于启用了 `-d` 标志。

## 替代方式 {#what-to-use-instead}

要在当前模块的上下文中安装可执行文件，请使用不带版本后缀的 `go install`，如下所示。它会应用当前目录或父目录中 `go.mod` 文件所声明的版本要求和其他指令。

```
go install example.com/cmd
```

要在忽略当前模块的情况下安装可执行文件，请为 `go install` **添加**[版本后缀](/ref/mod#version-queries)，例如 `@v1.2.3` 或 `@latest`，如下所示。带版本后缀时，`go install` 不会读取或更新当前目录或父目录中的 `go.mod` 文件。

```
# Install a specific version.
go install example.com/cmd@v1.2.3

# Install the highest available version.
go install example.com/cmd@latest
```

为避免歧义，`go install` 带版本后缀时，所有参数都必须指向同一模块、同一版本中的 `main` 包。如果该模块包含 `go.mod` 文件，其中不得包含会导致该模块在作为主模块时被作不同解释的 `replace`、`exclude` 等指令。该模块的 `vendor` 目录也不会被使用。

详情请参阅 [`go install`](/ref/mod#go-install)。

## 为什么做出这项调整 {#why-this-is-happening}

自模块机制引入以来，`go get` 既用于更新 `go.mod` 中的依赖，也用于安装命令。这两种职责混在一起，常常带来困惑和不便：大多数时候，开发者只是想更新依赖，或者安装命令，并不希望同时进行两项操作。

从 Go 1.16 开始，`go install` 可以安装命令行指定版本的命令，同时忽略当前目录中的 `go.mod` 文件（如果存在）。因此，现在大多数情况下都应使用 `go install` 安装命令。

由于构建和安装命令的功能与 `go install` 重复，`go get` 的这项能力已被弃用。移除后，`go get` 默认不再编译或链接包，执行速度会更快。更新无法针对当前平台构建的包时，`go get` 也不会因此报错。

完整讨论请参阅提案 [#40276](/issue/40276)。
