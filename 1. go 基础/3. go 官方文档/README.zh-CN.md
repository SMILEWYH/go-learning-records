# Go 官方网站文档 · 中文学习译本

本项目基于下载的 `golang.org/x/website` 源码翻译，沿用官方网站的 Go 服务和页面结构。
这是个人中文学习译本，不是 Go 项目发布的官方中文版。本次翻译当前技术文档；博客、历史演讲、历史版本说明和早期更新记录保留英文。

## 本地运行

仓库的 `go.mod` 要求 Go 1.26 或更新版本。本机已使用 Go 1.27.1 验证。
运行网站不需要 Node.js、npm 或 Docker。启动脚本使用 `localcontent` 构建标签，直接读取项目中的 `_content`，避免启动前将全部文档和图片打包进程序；博客订阅和 Go 语言之旅在首次访问相应地址时加载。默认不带该标签的上游构建仍使用内嵌内容。

在项目目录打开终端：

```sh
./run.sh
```

等待终端显示“中文文档已启动”后，打开 **[中文文档目录](http://localhost:6060/go.dev/doc/)**。
macOS 也可以双击 `start.command`。终端中的服务需要保持运行，按 `Ctrl+C` 停止。

常用入口：

- [入门与进阶教程](http://localhost:6060/go.dev/doc/tutorial/)
- [Go 语言之旅](http://localhost:6060/tour/)
- [Effective Go](http://localhost:6060/go.dev/doc/effective_go)
- [管理依赖](http://localhost:6060/go.dev/doc/modules/managing-dependencies)
- [常用中文参考文档](http://localhost:6060/go.dev/doc/reference-zh)
- [标准库目录](http://localhost:6060/go.dev/pkg/)（180 个公开包的中文参考，包含 builtin）
- [翻译范围与术语](http://localhost:6060/go.dev/doc/translation)

主站沿用上游的开发路由约定：使用 `localhost`，页面路径中带 `/go.dev/`。
主站直接访问 `127.0.0.1:6060` 会被上游路由重定向到线上网站，请使用上面的地址。

首次启动需要联网下载 Go 依赖。脚本为当前进程设置国内模块代理，不修改全局 Go 配置；如已有 `GOPROXY` 环境变量，优先使用你的设置。可自行覆盖：

```sh
GOPROXY=https://proxy.golang.org,direct ./run.sh
```

如果 6060 端口已被占用：

```sh
./run.sh -http=localhost:6061
```

然后访问 `http://localhost:6061/go.dev/doc/`，终端会在端口绑定成功后打印实际访问地址。

不用脚本也可以启动（适用于 Windows 等环境）：

```sh
go run -tags=localcontent ./cmd/golangorg -http=localhost:6060 -content=_content -tip=false -wiki=false -gopls=false -localdocs
```

`-localdocs` 优先显示项目内的中文标准库与命令快照；内部实现包、其他平台视图及未导出的工具内容从本机 `GOROOT` 渲染。页面顶部的英文对照入口以及显式平台/显示模式查询使用本机英文文档。第三方包仍链接到 pkg.go.dev。

新增内容从 **[中文参考文档阅读入口](http://localhost:6060/go.dev/doc/reference-zh)** 进入：180 个公开标准库包、完整 go 命令手册、21 项独立工具文档、完整语言规范与 Go 内存模型，以及汇编指南和 GODEBUG 参考。

## 在本机运行 Tour 练习

主站内的 Playground/Tour 运行按钮沿用上游在线编译服务。需要在本机编译练习时，在另一个终端运行：

```sh
./run.sh tour
```

访问 **[本地 Go 语言之旅](http://127.0.0.1:3999/tour/)**。
它监听回环地址，示例由本机 Go 编译器运行。依赖准备完成后，常规课程可离线学习；引用外部包的示例仍可能需要下载依赖。
课程加载后会保存在内存中：主站在首次访问语言之旅时加载，本机练习服务在启动时加载。修改 Tour 内容后，需要重启对应服务。
普通 `_content` 文档修改后刷新页面即可。

## 翻译与维护

当前技术文档与课程、180 个公开标准库包（包含 builtin）的 API 正文及字段和示例注释、完整 go 命令手册、完整语言规范与内存模型均提供中文译文。另补齐 21 项独立工具文档、汇编指南、GODEBUG 参考，以及本仓库的网站说明、企业案例正文、应用场景和案例卡片。另有 arena 实验包的概述译文，默认构建视图未导出其 API 声明。relnote 的上游导出页面只有目录标题，未额外编写不存在的手册。

Go 语言之旅按 7 组课程计数；网站案例卡片按条目计数，外链文章不计为全文译文。33 份版本说明及早期更新记录保留英文。

- [逐篇翻译清单](translations/STATUS.md)记录完成、保留英文的版本说明、重定向和代码示例，未完成内容不会标成已完成。
- [机器可读清单](translations/manifest.json)包含源文件路径与英文原文 SHA-256。
- `translations/en/` 保存文档英文原文，`translations/ui-en/` 保存本轮界面模板原文，方便对照和今后同步上游更新。
- `translations/completed.txt` 只列出已完成译文的文档与页面条目。
- `translations/partial.json` 用于登记部分译文，目前为空；语言规范和 go 命令手册已全文翻译。
- `translations/reference-sources.json` 记录 Go 工具链、语言版本、平台和来源。
- `translations/prose-blocks.json` 登记以 `<pre>` 排版的自然语言说明表及审核哈希；这些说明可以翻译，实际程序代码仍保持一致。
- `translations/comment-blocks.json` 登记允许翻译自然语言注释的代码块及前后哈希；`scripts/go_comment_structure.go` 比较 Go 词法单元、编译指令和示例输出断言。
- [术语说明](http://localhost:6060/go.dev/doc/translation#terminology)统一模块、包、切片、接收者、类型断言、通道等用语。

程序逻辑、命令、包路径、API 标识符、字符串和终端输出保留原样；已译标准库页面中的字段说明、接口注释和示例自然语言注释采用中文。本机 GOROOT 源码不修改。发现原文示例问题时，用“译注”单独说明。
部分图示内嵌文字和外部资源仍是英文。新增参考文档原文来自 **Go 1.27.1**，标准库采用 **linux/amd64** 视图，译文保存在 `_content/pkg/`、`_content/cmd/` 和 `_content/ref/`。内部实现包、其他平台专属 API、需额外构建配置的包（如 runtime/cgo、BoringCrypto、SIMD）、外部子仓库及第三方包文档不在本地中文快照范围内。包的声明签名和示例程序逻辑保持原样，正文、参数表和自然语言注释采用中文。

升级本机 Go 不会自动更新中文快照。源码链接和英文对照使用本机工具链版本，可能与快照不同；页面顶部会标明版本。同步新版前，应在独立目录导出新原文并逐项比较，保留现有 `translations/en/` 与校验记录，不要直接用新原文覆盖旧快照。
JSON v2 等教程可能涉及实验性功能；以当前工具链的支持情况和教程中的说明为准。

下载 Go、外部链接、在线 Playground、Wiki 和 gopls 的完整外部仓库内容仍可能需要联网。本地启动脚本关闭了 Wiki、gopls 与 tip 的自动仓库拉取。

更新清单和校验译文：

```sh
python3 scripts/translation_synopses.py
python3 scripts/translation_status.py
python3 scripts/check_translations.py
python3 scripts/check_translations.py --base-url http://localhost:6060/go.dev
go test scripts/go_comment_structure.go scripts/go_comment_structure_test.go
go test -tags=localcontent ./internal/web ./internal/tour ./cmd/golangorg ./internal/dl ./internal/play
```

结构校验检查英文原文哈希、围栏和缩进代码块、可运行示例、自然语言注释及参数说明表的审核哈希、Go 词法单元和输出断言、参考链接、部分译文范围（如有）、Markdown 标题 ID、显式锚点和 Tour 指令；它不能代替译文语义审校。
原项目的开发与部署说明见 [README.en.md](README.en.md)，许可证继续使用仓库中的 [LICENSE](LICENSE)。
