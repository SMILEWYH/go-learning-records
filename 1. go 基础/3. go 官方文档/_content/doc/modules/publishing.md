<!--{
  "Title": "发布模块"
}-->

当你希望其他开发者使用某个模块时，需要发布它，使 Go 工具能够发现它。发布之后，导入模块中各个包的开发者就可以运行 `go get` 等命令，解析对该模块的依赖。

> **注意：** 模块的某个带标签版本发布后，不要再修改该版本。Go 工具会将下载的模块与首次下载的副本进行校验；如果两者不同，就会报告安全错误。请发布新版本，而不要修改已经发布的版本。

**另请参阅**

* 模块开发概览，参见[开发与发布模块](developing)。
* 包含发布步骤在内的整体开发流程，参见[模块发布与版本管理流程](release-workflow)。

## 发布步骤 {#publishing-steps}

按照以下步骤发布模块。

1. 打开命令行，进入本地仓库中模块的根目录。

1. 运行 `go mod tidy`，移除模块在开发过程中积累的、不再需要的依赖。

    ```
    $ go mod tidy
    ```

1. 最后再运行一次 `go test ./...`，确认一切正常。

    该命令会运行你使用 Go 测试框架编写的单元测试。

    ```
    $ go test ./...
    ok      example.com/mymodule       0.015s
    ```

1. 使用 `git tag` 命令，为项目打上新的版本标签。

    版本号应向用户说明本次发布所包含变更的性质。详情参见[模块版本号](version-numbers)。

    ```
    $ git commit -m "mymodule: changes for v0.1.0"
    $ git tag v0.1.0
    ```

1. 将新标签推送到远程 origin 仓库。

    ```
    $ git push origin v0.1.0
    ```

1. 运行 [`go list` 命令](/cmd/go/#hdr-List_packages_or_modules)，促使 Go 更新模块索引，加入你正在发布的模块信息。

    在命令前设置 `GOPROXY` 环境变量，指定一个 Go 模块代理，以确保请求会发送给该代理。

    ```
    $ GOPROXY=proxy.golang.org go list -m example.com/mymodule@v0.1.0
    ```

对该模块感兴趣的开发者，可以像使用其他模块一样，导入其中的包并运行 [`go get` 命令](/ref/mod#go-get)。既可以获取最新版本，也可以显式指定版本，例如：

```
$ go get example.com/mymodule@v0.1.0
```
