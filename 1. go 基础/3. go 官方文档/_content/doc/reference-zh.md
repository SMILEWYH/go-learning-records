---
title: 常用 Go 参考文档中文阅读入口
---

这份阅读入口面向已学完基础教程、准备开始写项目的读者。下面链接直接打开已翻译的官方参考文档；包文档可通过索引跳转到具体 API，无需按顺序通读每个函数。

标准库和命令文档以 **Go 1.27.1** 为原文快照，标准库采用 **linux/amd64** 的公开 API 视图。语言规范、go 命令手册与内存模型均已全文翻译。博客、历史演讲和历史版本说明继续保留英文。

## 先从常用内容开始 {#next-project}

以下五组内容适合入门后优先查阅：

| 内容 | 中文阅读入口 |
| --- | --- |
| API 字段说明与示例注释 | 已译包中的自然语言注释，重点包括 `net/http` 的 [Client](/pkg/net/http/#Client)、[Server](/pkg/net/http/#Server)、[Request](/pkg/net/http/#Request)、[Response](/pkg/net/http/#Response)、[Transport](/pkg/net/http/#Transport) |
| 常用语句与并发规则 | [全部语句](/ref/spec#Statements)、[接收操作](/ref/spec#Receive_operator)、[关闭通道](/ref/spec#Close)、[求值顺序](/ref/spec#Order_of_evaluation) |
| 五个实用标准库 | [bytes](/pkg/bytes/)、[flag](/pkg/flag/)、[log/slog](/pkg/log/slog/)、[io/fs](/pkg/io/fs/)、[embed](/pkg/embed/) |
| 环境与依赖排查 | [go env](/cmd/go/#hdr-Print_Go_environment_information)、[go list](/cmd/go/#hdr-List_packages_or_modules)、[go clean](/cmd/go/#hdr-Remove_object_files_and_cached_files)、[go mod download](/cmd/go/#hdr-Download_modules_to_local_cache)、[go mod verify](/cmd/go/#hdr-Verify_dependencies_have_expected_content)、[go work 及子命令](/cmd/go/#hdr-Workspace_maintenance) |
| 项目相关包 | [database/sql](/pkg/database/sql/)、[html/template](/pkg/html/template/)、[encoding/csv](/pkg/encoding/csv/)、[os/exec](/pkg/os/exec/)、[testing/synctest](/pkg/testing/synctest/) |

[标准库目录](/pkg/)现有 **180 个公开包**（包含 builtin）的中文参考。下方仍优先列出入门常用包，后面的分类入口用于按需深入。阅读字段注释时，尤其注意零值、`nil`、资源关闭责任、超时范围和并发使用条件。示例中的字符串、SQL、模板语法、输出及 `Output:` 断言保持原样；自然语言注释采用中文。

## 第一阶段：常用标准库 {#stdlib}

| 包与中文文档（按顺序） | 练习方向 |
| --- | --- |
| 1. [fmt](/pkg/fmt/)、[strings](/pkg/strings/)、[strconv](/pkg/strconv/) | 格式化输出、字符串处理、字符串与数字转换 |
| 2. [errors](/pkg/errors/)、[io](/pkg/io/)、[bufio](/pkg/bufio/) | 错误包装与判断、Reader/Writer、缓冲读写 |
| 3. [os](/pkg/os/)、[path/filepath](/pkg/path/filepath/)、[time](/pkg/time/) | 文件操作、文件路径、时间解析与定时器 |
| 4. [encoding/json](/pkg/encoding/json/)、[net/http](/pkg/net/http/)、[net/url](/pkg/net/url/) | JSON 编解码、HTTP 服务与客户端、URL 参数 |
| 5. [testing](/pkg/testing/)、[net/http/httptest](/pkg/net/http/httptest/) | 表驱动测试、基准测试、HTTP 处理器测试 |
| 6. [slices](/pkg/slices/)、[maps](/pkg/maps/) | 泛型切片与映射工具、迭代器 |
| 7. [bytes](/pkg/bytes/)、[flag](/pkg/flag/)、[log/slog](/pkg/log/slog/) | 字节切片与缓冲区、命令行参数、结构化日志 |
| 8. [io/fs](/pkg/io/fs/)、[embed](/pkg/embed/) | 文件系统接口、遍历目录、嵌入静态资源 |

可以先写一个读取 JSON 配置的命令行程序，再给它增加 HTTP 接口和测试。查到一个 API 后，阅读其说明、边界条件和示例，并在自己的项目中实际使用。

## 第二阶段：常用 go 命令及参数 {#commands}

[完整 go 命令手册](/cmd/go/)的全部章节与参数说明均已翻译，包括高级帮助主题。初学时可优先查阅下列命令：

| 使用场景 | 中文章节 |
| --- | --- |
| 创建模块 | [go mod init](/cmd/go/#hdr-Initialize_new_module_in_current_directory) |
| 运行程序 | [go run](/cmd/go/#hdr-Compile_and_run_Go_program) |
| 构建可执行文件 | [go build 与共享构建参数](/cmd/go/#hdr-Compile_packages_and_dependencies) |
| 整理依赖 | [go mod tidy](/cmd/go/#hdr-Add_missing_and_remove_unused_modules) |
| 添加、升级、降级依赖 | [go get](/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them) |
| 安装命令行工具 | [go install](/cmd/go/#hdr-Compile_and_install_packages_and_dependencies) |
| 格式化与静态检查 | [go fmt](/cmd/go/#hdr-Gofmt__reformat__package_sources)、[go vet](/cmd/go/#hdr-Report_likely_mistakes_in_packages) |
| 测试与性能分析 | [go test](/cmd/go/#hdr-Test_packages)、[测试参数](/cmd/go/#hdr-Testing_flags)、[测试函数](/cmd/go/#hdr-Testing_functions) |
| 查询 API | [go doc](/cmd/go/#hdr-Show_documentation_for_package_or_symbol) |
| 理解 `.`、`./...` 和包路径 | [包列表与匹配模式](/cmd/go/#hdr-Package_lists_and_patterns) |
| 查看环境与包信息 | [go env](/cmd/go/#hdr-Print_Go_environment_information)、[go list](/cmd/go/#hdr-List_packages_or_modules) |
| 清理构建产物与缓存 | [go clean](/cmd/go/#hdr-Remove_object_files_and_cached_files) |
| 下载与校验依赖 | [go mod download](/cmd/go/#hdr-Download_modules_to_local_cache)、[go mod verify](/cmd/go/#hdr-Verify_dependencies_have_expected_content) |
| 多模块工作区 | [概览](/cmd/go/#hdr-Workspace_maintenance)、[init](/cmd/go/#hdr-Initialize_workspace_file)、[use](/cmd/go/#hdr-Add_modules_to_workspace_file)、[sync](/cmd/go/#hdr-Sync_workspace_build_list_to_modules)、[edit](/cmd/go/#hdr-Edit_go_work_from_tools_or_scripts)、[vendor](/cmd/go/#hdr-Vendor_workspace_dependencies) |

先关注 `-o`、`-v`、`-race`、`-run`、`-count`、`-cover`、`-bench`、`-benchmem`、`-timeout` 等常用参数，再按需要阅读交叉编译、链接器参数和性能剖析相关选项。终端中的 `go help` 仍显示工具链自带的英文帮助，本项目提供的是 Web 中文译文。

## 第三阶段：查阅语言规范 {#spec}

[Go 语言规范](/ref/spec)已全文翻译。适合先查阅容易混淆的规则，再回到程序验证：

1. [变量](/ref/spec#Variables)、[零值与初始化](/ref/spec#Program_initialization_and_execution)：理解默认值、初始化顺序与程序启动。
2. [切片](/ref/spec#Slice_types)、[映射](/ref/spec#Map_types)、[底层类型](/ref/spec#Underlying_types)：理解常用类型及其性质。
3. [接口](/ref/spec#Interface_types)、[方法集](/ref/spec#Method_sets)、[方法声明](/ref/spec#Method_declarations)：判断类型何时实现接口。
4. [类型断言](/ref/spec#Type_assertions)、[类型转换](/ref/spec#Conversions)：区分运行时类型检查与转换规则。
5. [类型参数声明](/ref/spec#Type_parameter_declarations)：理解类型约束和泛型声明。继续查阅[实例化](/ref/spec#Instantiations)和[类型推导](/ref/spec#Type_inference)，理解泛型函数的调用规则。
6. [defer](/ref/spec#Defer_statements)、[panic/recover](/ref/spec#Handling_panics)：理解延迟调用及异常展开。
7. [append/copy](/ref/spec#Appending_and_copying_slices)、[make](/ref/spec#Making_slices_maps_and_channels)、[new](/ref/spec#Allocation)：查阅内置函数的精确语义。
8. [if](/ref/spec#If_statements)、[for/range](/ref/spec#For_statements)、[switch](/ref/spec#Switch_statements)：理解分支、循环及迭代变量的规则。
9. [go 语句](/ref/spec#Go_statements)、[select](/ref/spec#Select_statements)、[发送](/ref/spec#Send_statements)、[接收](/ref/spec#Receive_operator)、[close](/ref/spec#Close)：理解启动 goroutine、通道阻塞、关闭和多路选择。
10. [求值顺序](/ref/spec#Order_of_evaluation)、[len/cap](/ref/spec#Length_and_capacity)：区分已规定的求值顺序与未规定的部分，查阅长度和容量的边界规则。

词法、类型、表达式、内置函数、全部语句、程序初始化与执行等章节均已完整翻译。规范用于核对精确规则，不必一次通读。

## 第四阶段：并发与内存模型 {#concurrency}

1. [context](/pkg/context/)：先学会传递取消信号、超时和请求范围内的值。
2. [sync](/pkg/sync/)：学习 Mutex、RWMutex、WaitGroup、Once、Cond、Map 等同步工具及使用限制。
3. [Go 内存模型](/ref/mem)：理解 goroutine 之间的可见性、先行发生（happens-before）关系和同步保证。
4. [sync/atomic](/pkg/sync/atomic/)：在理解同步关系后阅读原子操作，再判断项目是否需要它。
5. [testing/synctest](/pkg/testing/synctest/)：使用隔离测试环境、虚拟时钟和 `Wait`，测试并发代码及超时行为。

结合已经翻译的[数据竞争检测器](/doc/articles/race_detector)文档，在练习中运行 `go test -race ./...`，观察并修复共享数据访问问题。

## 第五阶段：结合项目查阅 {#projects}

| 项目需求 | 中文包文档与已有教程 |
| --- | --- |
| 数据库查询、事务与连接池 | [database/sql](/pkg/database/sql/)，配合[访问关系型数据库](/doc/database/)；具体数据库驱动的文档仍由相应项目提供 |
| 服务端生成 HTML | [html/template](/pkg/html/template/)，查阅上下文自动转义与可信内容类型的规则 |
| 导入或导出表格数据 | [encoding/csv](/pkg/encoding/csv/)，查阅引号、分隔符、字段数量和读写错误 |
| 调用外部程序 | [os/exec](/pkg/os/exec/)，查阅参数传递、标准输入输出、上下文取消与 `Wait` |
| 检查并发逻辑 | [testing/synctest](/pkg/testing/synctest/)，配合 [context](/pkg/context/) 和[内存模型](/ref/mem) |

## 更多标准库：按需求查找 {#more-packages}

| 方向 | 中文参考入口 |
| --- | --- |
| 数据结构、排序与迭代 | [cmp](/pkg/cmp/)、[container/heap](/pkg/container/heap/)、[container/list](/pkg/container/list/)、[container/ring](/pkg/container/ring/)、[iter](/pkg/iter/)、[sort](/pkg/sort/) |
| 文本、正则与模板 | [regexp](/pkg/regexp/)、[text/template](/pkg/text/template/)、[text/scanner](/pkg/text/scanner/)、[unicode](/pkg/unicode/)、[unicode/utf8](/pkg/unicode/utf8/) |
| 编解码与压缩 | [encoding/xml](/pkg/encoding/xml/)、[encoding/gob](/pkg/encoding/gob/)、[encoding/binary](/pkg/encoding/binary/)、[encoding/base64](/pkg/encoding/base64/)、[archive/zip](/pkg/archive/zip/)、[compress/gzip](/pkg/compress/gzip/) |
| 网络与服务 | [net](/pkg/net/)、[net/netip](/pkg/net/netip/)、[net/http/httputil](/pkg/net/http/httputil/)、[mime/multipart](/pkg/mime/multipart/)、[database/sql/driver](/pkg/database/sql/driver/) |
| 密码学与证书 | [crypto](/pkg/crypto/)、[crypto/rand](/pkg/crypto/rand/)、[crypto/sha256](/pkg/crypto/sha256/)、[crypto/tls](/pkg/crypto/tls/)、[crypto/x509](/pkg/crypto/x509/) |
| 数值、图像与日志 | [math](/pkg/math/)、[math/big](/pkg/math/big/)、[math/rand/v2](/pkg/math/rand/v2/)、[image](/pkg/image/)、[image/draw](/pkg/image/draw/)、[log](/pkg/log/) |
| 反射、运行时与性能 | [reflect](/pkg/reflect/)、[runtime](/pkg/runtime/)、[runtime/debug](/pkg/runtime/debug/)、[runtime/metrics](/pkg/runtime/metrics/)、[runtime/pprof](/pkg/runtime/pprof/)、[runtime/trace](/pkg/runtime/trace/) |
| Go 源码分析 | [go/ast](/pkg/go/ast/)、[go/parser](/pkg/go/parser/)、[go/token](/pkg/go/token/)、[go/types](/pkg/go/types/)、[go/doc](/pkg/go/doc/) |
| 实验性 JSON v2 | [encoding/json/v2](/pkg/encoding/json/v2/)、[encoding/json/jsontext](/pkg/encoding/json/jsontext/)，配合[迁移指南](/doc/jsonv2-migration)，注意包文档中的实验要求 |

完整列表见[标准库目录](/pkg/)。180 个包来自 linux/amd64 的默认构建列表（导出时未启用 cgo），包含 builtin；内部实现、其他平台专属 API 及需额外构建配置的 runtime/cgo、BoringCrypto、SIMD 等包不在此快照内，目录已标明其英文入口。另有 [arena 实验包概述](/pkg/arena/)；它受实验构建约束控制，本页不包含 API 声明。包文档中的弃用说明、使用限制及安全条件均保留原意；应结合实际任务选择 API。

## 独立工具与底层参考 {#tools}

| 用途 | 中文文档 |
| --- | --- |
| 格式化、检查与源码修复 | [gofmt](/cmd/gofmt/)、[vet](/cmd/vet/)、[fix](/cmd/fix/) |
| 测试与覆盖率 | [test2json](/cmd/test2json/)、[cover](/cmd/cover/)、[covdata](/cmd/covdata/) |
| 性能剖析与追踪 | [pprof](/cmd/pprof/)、[trace](/cmd/trace/)、[preprofile](/cmd/preprofile/) |
| 编译、链接与 C 互操作 | [compile](/cmd/compile/)、[link](/cmd/link/)、[cgo](/cmd/cgo/) |
| 汇编、符号与目标文件 | [asm](/cmd/asm/)、[Go 汇编器快速指南](/doc/asm)、[addr2line](/cmd/addr2line/)、[nm](/cmd/nm/)、[objdump](/cmd/objdump/)、[pack](/cmd/pack/)、[buildid](/cmd/buildid/) |
| 工具链构建与分发 | [dist](/cmd/dist/)、[distpack](/cmd/distpack/)；[relnote](/cmd/relnote/)的上游导出内容仅有目录标题 |
| 兼容性与运行时设置 | [Go、向后兼容性与 GODEBUG](/doc/godebug) |

这些内容适合遇到具体构建、调试或兼容性问题时查阅。还可以阅读已译的[应用场景](/solutions/use-cases)、[企业案例](/solutions/case-studies)和[项目介绍](/project)，了解 Go 的实际应用；外链案例只翻译本站卡片的标题与简介。

## 版本与英文对照 {#version}

标准库 API 和命令的中文内容保存在本项目中，不修改本机 Go 安装目录。包页面顶部的“查看本机 Go 的英文文档”会切换到本机工具链生成的英文参考。指定 `GOOS`、`GOARCH` 或显示模式的查询也使用本机英文文档。

如果以后升级本机 Go，英文参考和源码链接可能已更新，而中文快照仍对应 Go 1.27.1；请以页首版本标记区分。文件权限、路径、系统调用等行为还需结合目标平台阅读。

原文存放于 `translations/en/`。已译标准库页面中的字段说明、接口注释和示例自然语言注释采用中文；API 签名、程序逻辑、字符串、编译指令及输出断言保持原样。点击源码链接看到的本机 Go 源文件仍是英文。参数说明表和格式动词表中的说明文字也译为中文。在线运行和分享示例仍使用上游服务，需要联网。

完整状态和术语见[中文翻译说明](/doc/translation)。
