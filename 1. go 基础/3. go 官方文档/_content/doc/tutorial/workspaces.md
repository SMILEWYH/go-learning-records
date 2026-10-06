<!--{
  "Title": "教程：多模块工作区入门",
  "Breadcrumb": true
}-->

本教程介绍 Go 多模块工作区的基础知识。借助多模块工作区，你可以告诉 Go 命令，自己正在同时开发多个模块，并方便地构建和运行这些模块中的代码。

在本教程中，你将在同一个多模块工作区内创建两个模块，跨模块修改代码，并通过构建查看这些修改的效果。

<!-- TODO TOC -->

**注意：**其他教程请参阅[教程列表](/doc/tutorial/index.html)。

## 准备工作 {#prerequisites}

*   **Go。**建议使用最新版本的 Go 学习本教程。安装步骤请参阅[安装 Go](/doc/install)。
*   **代码编辑工具。**任何文本编辑器都可以。
*   **命令行终端。**Go 可以在 Linux 和 Mac 的各种终端，以及 Windows 的 PowerShell 或 cmd 中正常使用。

## 为代码创建模块 {#create_folder}

首先，为即将编写的代码创建一个模块。

1. 打开命令行终端，切换到用户主目录。

   在 Linux 或 Mac 上：

    ```
    $ cd
    ```

   在 Windows 上：

    ```
    C:\> cd %HOMEPATH%
    ```

   本教程后续统一使用 $ 表示命令提示符。这些命令同样适用于 Windows。

2. 在命令行中创建名为 workspace 的目录，用来存放代码。

    ```
    $ mkdir workspace
    $ cd workspace
    ```

3. 初始化模块。

   本示例将创建一个新模块 `hello`，它依赖 golang.org/x/example 模块。

   创建 hello 模块：

   ```
   $ mkdir hello
   $ cd hello
   $ go mod init example.com/hello
   go: creating new go.mod: module example.com/hello
   ```

   使用 `go get` 添加对 golang.org/x/example/hello/reverse 包的依赖。

   ```
   $ go get golang.org/x/example/hello/reverse
   ```

   在 hello 目录中创建 hello.go 文件，内容如下：

   ```
   package main

   import (
       "fmt"

       "golang.org/x/example/hello/reverse"
   )

   func main() {
       fmt.Println(reverse.String("Hello"))
   }
   ```

   现在运行 hello 程序：

   ```
   $ go run .
   olleH
   ```

## 创建工作区 {#create-the-workspace}

这一步将创建 `go.work` 文件，定义一个包含该模块的工作区。

#### 初始化工作区 {#initialize-the-workspace}

在 `workspace` 目录中运行：

   ```
   $ go work init ./hello
   ```

`go work init` 命令会让 `go` 创建一个 `go.work` 文件，将 `./hello` 目录中的模块加入工作区。

`go` 命令生成的 `go.work` 文件如下所示：

   ```
   go 1.18

   use ./hello
   ```

`go.work` 文件的语法与 `go.mod` 类似。

`go` 指令告诉 Go 应使用哪个 Go 版本的语义解释这个文件，与 `go.mod` 文件中的 `go` 指令类似。

`use` 指令告诉 Go，在构建时应将 `hello` 目录中的模块作为主模块。

因此，在 `workspace` 的任意子目录中，这个模块都会生效。

#### 在工作区目录中运行程序 {#run-the-program-in-the-workspace-directory}

在 `workspace` 目录中运行：

   ```
   $ go run ./hello
   olleH
   ```

Go 命令将工作区中的所有模块都作为主模块处理，因此即使不在某个模块目录内，也可以引用该模块中的包。如果在模块和工作区之外运行 `go run` 命令，就会报错，因为 `go` 命令无法确定应使用哪些模块。

接下来，将 `golang.org/x/example/hello` 模块的本地副本添加到工作区中。这个模块位于 `go.googlesource.com/example` Git 仓库的一个子目录中。然后，我们将在 `reverse` 包中添加一个新函数，用来替代对 `String` 的直接调用。

## 下载并修改 `golang.org/x/example/hello` 模块 {#download-and-modify-the-golangorgxexamplehello-module}

   这一步将下载包含 `golang.org/x/example/hello` 模块的 Git 仓库副本，将其加入工作区，再添加一个新函数，供 hello 程序调用。

1. 克隆仓库。

   在 workspace 目录中运行 `git` 命令，克隆仓库：

   ```
   $ git clone https://go.googlesource.com/example
   Cloning into 'example'...
   remote: Total 165 (delta 27), reused 165 (delta 27)
   Receiving objects: 100% (165/165), 434.18 KiB | 1022.00 KiB/s, done.
   Resolving deltas: 100% (27/27), done.
   ```

2. 将模块添加到工作区。

   Git 仓库刚刚被检出到 `./example` 目录。`golang.org/x/example/hello` 模块的源代码位于 `./example/hello` 中。将它添加到工作区：

   ```
   $ go work use ./example/hello
   ```

   `go work use` 命令会在 go.work 文件中添加一个新模块。文件现在如下所示：

   ```
   go 1.18

   use (
       ./hello
       ./example/hello
   )
   ```

   工作区现在同时包含 `example.com/hello` 和 `golang.org/x/example/hello` 两个模块，后者提供 `golang.org/x/example/hello/reverse` 包。

   这样就能使用我们将在 `reverse` 包本地副本中编写的新代码，而不是模块缓存中通过 `go get` 命令下载的版本。

3. 添加新函数。

   我们将在 `golang.org/x/example/hello/reverse` 包中添加一个新函数，用于反转整数的十进制数字顺序。

   在 `workspace/example/hello/reverse` 目录中创建名为 `int.go` 的文件，内容如下：

   ```
   package reverse

   import "strconv"

   // Int returns the decimal reversal of the integer i.
   func Int(i int) int {
       i, _ = strconv.Atoi(String(strconv.Itoa(i)))
       return i
   }
   ```

4. 修改 hello 程序，使用新函数。

   将 `workspace/hello/hello.go` 修改为如下内容：

   ```
   package main

   import (
       "fmt"

       "golang.org/x/example/hello/reverse"
   )

   func main() {
       fmt.Println(reverse.String("Hello"), reverse.Int(24601))
   }
   ```

#### 在工作区中运行代码 {#run-the-code-in-the-workspace}

   在 workspace 目录中运行：

   ```
   $ go run ./hello
   olleH 10642
   ```

   Go 命令根据命令行中指定的路径，在 `go.work` 文件列出的 `hello` 目录中找到 `example.com/hello` 模块。同样，它也会根据 `go.work` 文件解析 `golang.org/x/example/hello/reverse` 的导入。

   跨多个模块开发时，可以使用 `go.work`，无需逐一添加 [`replace`](/ref/mod#go-mod-file-replace) 指令。

   由于这两个模块位于同一工作区，可以方便地修改其中一个模块，并立即在另一个模块中使用这些修改。

#### 后续步骤 {#future-step}

   要正式发布这些模块，需要先发布 `golang.org/x/example/hello` 模块的一个版本，例如 `v0.1.0`。通常做法是在模块的版本控制仓库中为某个提交打上版本标签。详细说明请参阅[模块发布流程文档](/doc/modules/release-workflow)。发布完成后，就可以在 `hello/go.mod` 中提高对 `golang.org/x/example/hello` 模块的依赖版本：

   ```
   cd hello
   go get golang.org/x/example/hello@v0.1.0
   ```

   这样，`go` 命令在工作区之外也能正确解析这些模块。

## 进一步了解工作区 {#learn-more-about-workspaces}

   除了本教程前面介绍的 `go work init`，`go` 命令还提供了几个操作工作区的子命令：

   - `go work use [-r] [dir]`：如果 `dir` 存在，就在 `go.work` 文件中为它添加一条 `use` 指令；如果参数指定的目录不存在，则移除对应的 `use` 条目。`-r` 标志会递归检查 `dir` 的子目录。
   - `go work edit`：编辑 `go.work` 文件，用法与 `go mod edit` 类似。
   - `go work sync`：将工作区构建列表中的依赖版本同步到工作区的各个模块中。

   有关工作区和 `go.work` 文件的更多说明，请参阅《Go 模块参考》中的[工作区](/ref/mod#workspaces)。
