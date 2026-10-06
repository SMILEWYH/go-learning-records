---
title: 中文翻译说明
---

这是基于 Go 官方网站源码制作的简体中文学习译本，不是 Go 项目发布的官方中文版本。
本次翻译范围是仓库中的当前技术文档、学习入口及网站页面；博客、历史演讲、版本说明和早期更新记录保留英文。

## 阅读方式 {#reading}

从[文档目录](/doc/)或[入门教程](/doc/tutorial/)开始，也可以通过 [Go 语言之旅](/tour/)学习语言基础。
译文保留原文的代码、命令、包路径、API 名称和章节锚点。
已译标准库中的字段说明、接口注释与示例自然语言注释采用中文；程序逻辑、字符串、编译指令、输出和输出断言保留原样，方便直接复制和对照。

## 翻译范围与清单 {#coverage}

<!-- translation-coverage:start -->

已完成 **342 项文档、包参考、课程与网站页面条目**的翻译（程序逻辑与输出保留原文，已译标准库的字段和示例注释采用中文）。

已完成清单中登记的翻译范围；Go 语言之旅按 7 组课程计数；网站案例卡片也按独立源文件计数。公开标准库采用 linux/amd64 快照，内部实现包、其他平台视图和外部网站不在此范围内。

版本说明和早期更新记录共 33 份，按本次约定保留英文。

英文原文保存在项目的 `translations/en/`，逐篇状态与原文哈希见 `translations/manifest.json`。标准库与命令参考快照为 Go 1.27.1、linux/amd64。

| 文档或课程 | 状态 |
| --- | --- |
| [Go 的作者](/AUTHORS) | 已翻译 |
| [品牌与商标使用指南](/brand) | 已翻译 |
| [addr2line 命令](/cmd/addr2line/) | 已翻译 |
| [asm 命令](/cmd/asm/) | 已翻译 |
| [buildid 命令](/cmd/buildid/) | 已翻译 |
| [cgo 命令](/cmd/cgo/) | 已翻译 |
| [compile 命令](/cmd/compile/) | 已翻译 |
| [covdata 命令](/cmd/covdata/) | 已翻译 |
| [cover 命令](/cmd/cover/) | 已翻译 |
| [dist 命令](/cmd/dist/) | 已翻译 |
| [distpack 命令](/cmd/distpack/) | 已翻译 |
| [fix 命令](/cmd/fix/) | 已翻译 |
| [go 命令（完整中文译文）](/cmd/go/) | 已翻译 |
| [gofmt 命令](/cmd/gofmt/) | 已翻译 |
| [link 命令](/cmd/link/) | 已翻译 |
| [nm 命令](/cmd/nm/) | 已翻译 |
| [objdump 命令](/cmd/objdump/) | 已翻译 |
| [pack 命令](/cmd/pack/) | 已翻译 |
| [pprof 命令](/cmd/pprof/) | 已翻译 |
| [preprofile 命令](/cmd/preprofile/) | 已翻译 |
| [relnote 命令](/cmd/relnote/) | 已翻译 |
| [test2json 命令](/cmd/test2json/) | 已翻译 |
| [trace 命令](/cmd/trace/) | 已翻译 |
| [vet 命令](/cmd/vet/) | 已翻译 |
| [Go 社区行为准则](/conduct) | 已翻译 |
| [版权](/copyright) | 已翻译 |
| [关于 go 命令](/doc/articles/go_command) | 已翻译 |
| [技术文章](/doc/articles/) | 已翻译 |
| [数据竞争检测器](/doc/articles/race_detector) | 已翻译 |
| [编写 Web 应用](/doc/articles/wiki/) | 已翻译 |
| [Go 汇编器快速指南](/doc/asm) | 已翻译 |
| [集成测试的覆盖率分析](/doc/build-cover) | 已翻译 |
| [命令文档](/doc/cmd) | 已翻译 |
| [如何编写 Go 代码](/doc/code) | 已翻译 |
| [如何编写代码导读](/doc/codewalk/codewalk) | 已翻译 |
| [Go 中的一等函数](/doc/codewalk/functions) | 已翻译 |
| [生成随机文本：马尔可夫链算法](/doc/codewalk/markov) | 已翻译 |
| [通过通信共享内存](/doc/codewalk/sharemem) | 已翻译 |
| [Go 文档注释](/doc/comment) | 已翻译 |
| [贡献指南](/doc/contribute) | 已翻译 |
| [取消正在进行的操作](/doc/database/cancel-operations) | 已翻译 |
| [执行不返回数据的 SQL 语句](/doc/database/change-data) | 已翻译 |
| [执行事务](/doc/database/execute-transactions) | 已翻译 |
| [访问关系型数据库](/doc/database/) | 已翻译 |
| [管理连接](/doc/database/manage-connections) | 已翻译 |
| [打开数据库句柄](/doc/database/open-handle) | 已翻译 |
| [使用预处理语句](/doc/database/prepared-statements) | 已翻译 |
| [查询数据](/doc/database/querying) | 已翻译 |
| [避免 SQL 注入风险](/doc/database/sql-injection) | 已翻译 |
| [诊断工具](/doc/diagnostics) | 已翻译 |
| [编辑器插件与 IDE](/doc/editors) | 已翻译 |
| [高效 Go 编程（Effective Go）](/doc/effective_go) | 已翻译 |
| [常见问题（FAQ）](/doc/faq) | 已翻译 |
| [Go 垃圾回收器指南](/doc/gc-guide) | 已翻译 |
| [为 gccgo 前端作贡献](/doc/gccgo_contribute) | 已翻译 |
| [使用 GDB 调试 Go 代码](/doc/gdb) | 已翻译 |
| [弃用 go get 安装可执行文件的功能](/doc/go-get-install-deprecation) | 已翻译 |
| [Go 1 与 Go 程序的未来](/doc/go1compat) | 已翻译 |
| [Go、向后兼容性与 GODEBUG](/doc/godebug) | 已翻译 |
| [文档](/doc/) | 已翻译 |
| [配置和使用 gccgo](/doc/install/gccgo) | 已翻译 |
| [从源代码安装 Go](/doc/install/source) | 已翻译 |
| [下载与安装](/doc/install) | 已翻译 |
| [迁移到 encoding/json/v2](/doc/jsonv2-migration) | 已翻译 |
| [管理 Go 安装](/doc/manage-install) | 已翻译 |
| [开发与发布模块](/doc/modules/developing) | 已翻译 |
| [go.mod 文件参考](/doc/modules/gomod-ref) | 已翻译 |
| [组织 Go 模块](/doc/modules/layout) | 已翻译 |
| [开发主版本更新](/doc/modules/major-version) | 已翻译 |
| [管理依赖](/doc/modules/managing-dependencies) | 已翻译 |
| [管理模块源码](/doc/modules/managing-source) | 已翻译 |
| [发布模块](/doc/modules/publishing) | 已翻译 |
| [模块发布与版本管理流程](/doc/modules/release-workflow) | 已翻译 |
| [模块版本号](/doc/modules/version-numbers) | 已翻译 |
| [基于性能剖析的优化（PGO）](/doc/pgo) | 已翻译 |
| [Go 开发者安全最佳实践](/doc/security/best-practices) | 已翻译 |
| [Go 安全问题判定说明](/doc/security/decisions) | 已翻译 |
| [FIPS 140-3 合规支持](/doc/security/fips140) | 已翻译 |
| [Go 模糊测试](/doc/security/fuzz/) | 已翻译 |
| [Go 模糊测试技术细节](/doc/security/fuzz/technical) | 已翻译 |
| [安全](/doc/security/) | 已翻译 |
| [Go 安全策略](/doc/security/policy) | 已翻译 |
| [Go 威胁模型](/doc/security/threat-model) | 已翻译 |
| [Go CNA 策略](/doc/security/vuln/cna) | 已翻译 |
| [Go 漏洞数据库](/doc/security/vuln/database) | 已翻译 |
| [在 IDE 中扫描漏洞](/doc/security/vuln/editor) | 已翻译 |
| [Go 漏洞管理](/doc/security/vuln/) | 已翻译 |
| [Go 遥测](/doc/telemetry) | 已翻译 |
| [Go 工具链](/doc/toolchain) | 已翻译 |
| [添加测试](/doc/tutorial/add-a-test) | 已翻译 |
| [从另一个模块调用代码](/doc/tutorial/call-module-code) | 已翻译 |
| [编译和安装应用程序](/doc/tutorial/compile-install) | 已翻译 |
| [教程：创建 Go 模块](/doc/tutorial/create-module) | 已翻译 |
| [教程：访问关系型数据库](/doc/tutorial/database-access) | 已翻译 |
| [教程：模糊测试入门](/doc/tutorial/fuzz) | 已翻译 |
| [教程：泛型入门](/doc/tutorial/generics) | 已翻译 |
| [教程：Go 入门](/doc/tutorial/getting-started) | 已翻译 |
| [教程：使用 VS Code Go 查找并修复存在漏洞的依赖](/doc/tutorial/govulncheck-ide) | 已翻译 |
| [教程：使用 govulncheck 查找并修复存在漏洞的依赖](/doc/tutorial/govulncheck) | 已翻译 |
| [为多人返回问候语](/doc/tutorial/greetings-multiple-people) | 已翻译 |
| [返回和处理错误](/doc/tutorial/handle-errors) | 已翻译 |
| [教程](/doc/tutorial/) | 已翻译 |
| [教程：处理 JSON](/doc/tutorial/json) | 已翻译 |
| [总结](/doc/tutorial/module-conclusion) | 已翻译 |
| [返回随机问候语](/doc/tutorial/random-greeting) | 已翻译 |
| [教程：使用 Go 和 Gin 开发 RESTful API](/doc/tutorial/web-service-gin) | 已翻译 |
| [教程：多模块工作区入门](/doc/tutorial/workspaces) | 已翻译 |
| [帮助](/help) | 已翻译 |
| [archive/tar 包](/pkg/archive/tar/) | 已翻译 |
| [archive/zip 包](/pkg/archive/zip/) | 已翻译 |
| [arena 实验包概述](/pkg/arena/) | 已翻译 |
| [bufio 包](/pkg/bufio/) | 已翻译 |
| [builtin 包](/pkg/builtin/) | 已翻译 |
| [bytes 包](/pkg/bytes/) | 已翻译 |
| [cmp 包](/pkg/cmp/) | 已翻译 |
| [compress/bzip2 包](/pkg/compress/bzip2/) | 已翻译 |
| [compress/flate 包](/pkg/compress/flate/) | 已翻译 |
| [compress/gzip 包](/pkg/compress/gzip/) | 已翻译 |
| [compress/lzw 包](/pkg/compress/lzw/) | 已翻译 |
| [compress/zlib 包](/pkg/compress/zlib/) | 已翻译 |
| [container/heap 包](/pkg/container/heap/) | 已翻译 |
| [container/list 包](/pkg/container/list/) | 已翻译 |
| [container/ring 包](/pkg/container/ring/) | 已翻译 |
| [context 包](/pkg/context/) | 已翻译 |
| [crypto/aes 包](/pkg/crypto/aes/) | 已翻译 |
| [crypto/cipher 包](/pkg/crypto/cipher/) | 已翻译 |
| [crypto/des 包](/pkg/crypto/des/) | 已翻译 |
| [crypto/dsa 包](/pkg/crypto/dsa/) | 已翻译 |
| [crypto/ecdh 包](/pkg/crypto/ecdh/) | 已翻译 |
| [crypto/ecdsa 包](/pkg/crypto/ecdsa/) | 已翻译 |
| [crypto/ed25519 包](/pkg/crypto/ed25519/) | 已翻译 |
| [crypto/elliptic 包](/pkg/crypto/elliptic/) | 已翻译 |
| [crypto/fips140 包](/pkg/crypto/fips140/) | 已翻译 |
| [crypto/hkdf 包](/pkg/crypto/hkdf/) | 已翻译 |
| [crypto/hmac 包](/pkg/crypto/hmac/) | 已翻译 |
| [crypto/hpke 包](/pkg/crypto/hpke/) | 已翻译 |
| [crypto 包](/pkg/crypto/) | 已翻译 |
| [crypto/md5 包](/pkg/crypto/md5/) | 已翻译 |
| [crypto/mldsa 包](/pkg/crypto/mldsa/) | 已翻译 |
| [crypto/mlkem 包](/pkg/crypto/mlkem/) | 已翻译 |
| [crypto/mlkem/mlkemtest 包](/pkg/crypto/mlkem/mlkemtest/) | 已翻译 |
| [crypto/pbkdf2 包](/pkg/crypto/pbkdf2/) | 已翻译 |
| [crypto/rand 包](/pkg/crypto/rand/) | 已翻译 |
| [crypto/rc4 包](/pkg/crypto/rc4/) | 已翻译 |
| [crypto/rsa 包](/pkg/crypto/rsa/) | 已翻译 |
| [crypto/sha1 包](/pkg/crypto/sha1/) | 已翻译 |
| [crypto/sha256 包](/pkg/crypto/sha256/) | 已翻译 |
| [crypto/sha3 包](/pkg/crypto/sha3/) | 已翻译 |
| [crypto/sha512 包](/pkg/crypto/sha512/) | 已翻译 |
| [crypto/subtle 包](/pkg/crypto/subtle/) | 已翻译 |
| [crypto/tls 包](/pkg/crypto/tls/) | 已翻译 |
| [crypto/x509 包](/pkg/crypto/x509/) | 已翻译 |
| [crypto/x509/pkix 包](/pkg/crypto/x509/pkix/) | 已翻译 |
| [database/sql/driver 包](/pkg/database/sql/driver/) | 已翻译 |
| [database/sql 包](/pkg/database/sql/) | 已翻译 |
| [debug/buildinfo 包](/pkg/debug/buildinfo/) | 已翻译 |
| [debug/dwarf 包](/pkg/debug/dwarf/) | 已翻译 |
| [debug/elf 包](/pkg/debug/elf/) | 已翻译 |
| [debug/gosym 包](/pkg/debug/gosym/) | 已翻译 |
| [debug/macho 包](/pkg/debug/macho/) | 已翻译 |
| [debug/pe 包](/pkg/debug/pe/) | 已翻译 |
| [debug/plan9obj 包](/pkg/debug/plan9obj/) | 已翻译 |
| [embed 包](/pkg/embed/) | 已翻译 |
| [encoding/ascii85 包](/pkg/encoding/ascii85/) | 已翻译 |
| [encoding/asn1 包](/pkg/encoding/asn1/) | 已翻译 |
| [encoding/base32 包](/pkg/encoding/base32/) | 已翻译 |
| [encoding/base64 包](/pkg/encoding/base64/) | 已翻译 |
| [encoding/binary 包](/pkg/encoding/binary/) | 已翻译 |
| [encoding/csv 包](/pkg/encoding/csv/) | 已翻译 |
| [encoding/gob 包](/pkg/encoding/gob/) | 已翻译 |
| [encoding/hex 包](/pkg/encoding/hex/) | 已翻译 |
| [encoding 包](/pkg/encoding/) | 已翻译 |
| [encoding/json 包](/pkg/encoding/json/) | 已翻译 |
| [encoding/json/jsontext 包](/pkg/encoding/json/jsontext/) | 已翻译 |
| [encoding/json/v2 包](/pkg/encoding/json/v2/) | 已翻译 |
| [encoding/pem 包](/pkg/encoding/pem/) | 已翻译 |
| [encoding/xml 包](/pkg/encoding/xml/) | 已翻译 |
| [errors 包](/pkg/errors/) | 已翻译 |
| [expvar 包](/pkg/expvar/) | 已翻译 |
| [flag 包](/pkg/flag/) | 已翻译 |
| [fmt 包](/pkg/fmt/) | 已翻译 |
| [go/ast 包](/pkg/go/ast/) | 已翻译 |
| [go/build/constraint 包](/pkg/go/build/constraint/) | 已翻译 |
| [go/build 包](/pkg/go/build/) | 已翻译 |
| [go/constant 包](/pkg/go/constant/) | 已翻译 |
| [go/doc/comment 包](/pkg/go/doc/comment/) | 已翻译 |
| [go/doc 包](/pkg/go/doc/) | 已翻译 |
| [go/format 包](/pkg/go/format/) | 已翻译 |
| [go/importer 包](/pkg/go/importer/) | 已翻译 |
| [go/parser 包](/pkg/go/parser/) | 已翻译 |
| [go/printer 包](/pkg/go/printer/) | 已翻译 |
| [go/scanner 包](/pkg/go/scanner/) | 已翻译 |
| [go/token 包](/pkg/go/token/) | 已翻译 |
| [go/types 包](/pkg/go/types/) | 已翻译 |
| [go/version 包](/pkg/go/version/) | 已翻译 |
| [hash/adler32 包](/pkg/hash/adler32/) | 已翻译 |
| [hash/crc32 包](/pkg/hash/crc32/) | 已翻译 |
| [hash/crc64 包](/pkg/hash/crc64/) | 已翻译 |
| [hash/fnv 包](/pkg/hash/fnv/) | 已翻译 |
| [hash 包](/pkg/hash/) | 已翻译 |
| [hash/maphash 包](/pkg/hash/maphash/) | 已翻译 |
| [html 包](/pkg/html/) | 已翻译 |
| [html/template 包](/pkg/html/template/) | 已翻译 |
| [image/color 包](/pkg/image/color/) | 已翻译 |
| [image/color/palette 包](/pkg/image/color/palette/) | 已翻译 |
| [image/draw 包](/pkg/image/draw/) | 已翻译 |
| [image/gif 包](/pkg/image/gif/) | 已翻译 |
| [image 包](/pkg/image/) | 已翻译 |
| [image/jpeg 包](/pkg/image/jpeg/) | 已翻译 |
| [image/png 包](/pkg/image/png/) | 已翻译 |
| [index/suffixarray 包](/pkg/index/suffixarray/) | 已翻译 |
| [io/fs 包](/pkg/io/fs/) | 已翻译 |
| [io 包](/pkg/io/) | 已翻译 |
| [io/ioutil 包](/pkg/io/ioutil/) | 已翻译 |
| [iter 包](/pkg/iter/) | 已翻译 |
| [log 包](/pkg/log/) | 已翻译 |
| [log/slog 包](/pkg/log/slog/) | 已翻译 |
| [log/syslog 包](/pkg/log/syslog/) | 已翻译 |
| [maps 包](/pkg/maps/) | 已翻译 |
| [math/big 包](/pkg/math/big/) | 已翻译 |
| [math/bits 包](/pkg/math/bits/) | 已翻译 |
| [math/cmplx 包](/pkg/math/cmplx/) | 已翻译 |
| [math 包](/pkg/math/) | 已翻译 |
| [math/rand 包](/pkg/math/rand/) | 已翻译 |
| [math/rand/v2 包](/pkg/math/rand/v2/) | 已翻译 |
| [mime 包](/pkg/mime/) | 已翻译 |
| [mime/multipart 包](/pkg/mime/multipart/) | 已翻译 |
| [mime/quotedprintable 包](/pkg/mime/quotedprintable/) | 已翻译 |
| [net/http/cgi 包](/pkg/net/http/cgi/) | 已翻译 |
| [net/http/cookiejar 包](/pkg/net/http/cookiejar/) | 已翻译 |
| [net/http/fcgi 包](/pkg/net/http/fcgi/) | 已翻译 |
| [net/http/httptest 包](/pkg/net/http/httptest/) | 已翻译 |
| [net/http/httptrace 包](/pkg/net/http/httptrace/) | 已翻译 |
| [net/http/httputil 包](/pkg/net/http/httputil/) | 已翻译 |
| [net/http 包](/pkg/net/http/) | 已翻译 |
| [net/http/pprof 包](/pkg/net/http/pprof/) | 已翻译 |
| [net 包](/pkg/net/) | 已翻译 |
| [net/mail 包](/pkg/net/mail/) | 已翻译 |
| [net/netip 包](/pkg/net/netip/) | 已翻译 |
| [net/rpc 包](/pkg/net/rpc/) | 已翻译 |
| [net/rpc/jsonrpc 包](/pkg/net/rpc/jsonrpc/) | 已翻译 |
| [net/smtp 包](/pkg/net/smtp/) | 已翻译 |
| [net/textproto 包](/pkg/net/textproto/) | 已翻译 |
| [net/url 包](/pkg/net/url/) | 已翻译 |
| [os/exec 包](/pkg/os/exec/) | 已翻译 |
| [os 包](/pkg/os/) | 已翻译 |
| [os/signal 包](/pkg/os/signal/) | 已翻译 |
| [os/user 包](/pkg/os/user/) | 已翻译 |
| [path/filepath 包](/pkg/path/filepath/) | 已翻译 |
| [path 包](/pkg/path/) | 已翻译 |
| [plugin 包](/pkg/plugin/) | 已翻译 |
| [reflect 包](/pkg/reflect/) | 已翻译 |
| [regexp 包](/pkg/regexp/) | 已翻译 |
| [regexp/syntax 包](/pkg/regexp/syntax/) | 已翻译 |
| [runtime/coverage 包](/pkg/runtime/coverage/) | 已翻译 |
| [runtime/debug 包](/pkg/runtime/debug/) | 已翻译 |
| [runtime 包](/pkg/runtime/) | 已翻译 |
| [runtime/metrics 包](/pkg/runtime/metrics/) | 已翻译 |
| [runtime/pprof 包](/pkg/runtime/pprof/) | 已翻译 |
| [runtime/race 包](/pkg/runtime/race/) | 已翻译 |
| [runtime/trace 包](/pkg/runtime/trace/) | 已翻译 |
| [slices 包](/pkg/slices/) | 已翻译 |
| [sort 包](/pkg/sort/) | 已翻译 |
| [strconv 包](/pkg/strconv/) | 已翻译 |
| [strings 包](/pkg/strings/) | 已翻译 |
| [structs 包](/pkg/structs/) | 已翻译 |
| [sync/atomic 包](/pkg/sync/atomic/) | 已翻译 |
| [sync 包](/pkg/sync/) | 已翻译 |
| [syscall 包](/pkg/syscall/) | 已翻译 |
| [testing/cryptotest 包](/pkg/testing/cryptotest/) | 已翻译 |
| [testing/fstest 包](/pkg/testing/fstest/) | 已翻译 |
| [testing 包](/pkg/testing/) | 已翻译 |
| [testing/iotest 包](/pkg/testing/iotest/) | 已翻译 |
| [testing/quick 包](/pkg/testing/quick/) | 已翻译 |
| [testing/slogtest 包](/pkg/testing/slogtest/) | 已翻译 |
| [testing/synctest 包](/pkg/testing/synctest/) | 已翻译 |
| [text/scanner 包](/pkg/text/scanner/) | 已翻译 |
| [text/tabwriter 包](/pkg/text/tabwriter/) | 已翻译 |
| [text/template 包](/pkg/text/template/) | 已翻译 |
| [text/template/parse 包](/pkg/text/template/parse/) | 已翻译 |
| [time 包](/pkg/time/) | 已翻译 |
| [time/tzdata 包](/pkg/time/tzdata/) | 已翻译 |
| [unicode 包](/pkg/unicode/) | 已翻译 |
| [unicode/utf16 包](/pkg/unicode/utf16/) | 已翻译 |
| [unicode/utf8 包](/pkg/unicode/utf8/) | 已翻译 |
| [unique 包](/pkg/unique/) | 已翻译 |
| [unsafe 包](/pkg/unsafe/) | 已翻译 |
| [uuid 包](/pkg/uuid/) | 已翻译 |
| [weak 包](/pkg/weak/) | 已翻译 |
| [Go 项目](/project) | 已翻译 |
| [Go 可复现构建报告](/rebuild) | 已翻译 |
| [Go 内存模型](/ref/mem) | 已翻译 |
| [Go 模块参考](/ref/mod) | 已翻译 |
| [Go 语言规范](/ref/spec) | 已翻译 |
| [Allegro：用 Go 编写支持数百万条目的高速缓存服务](/solutions/allegro) | 已翻译 |
| [American Express 使用 Go 构建支付与积分奖励系统](/solutions/americanexpress) | 已翻译 |
| [Armut Labs 如何使用 Go](/solutions/armut) | 已翻译 |
| [Bitly：为什么我们用 Go 编写一切](/solutions/bitly) | 已翻译 |
| [Go 在字节跳动的大规模实践](/solutions/bytedance) | 已翻译 |
| [Capital One：Serverless 与 Go 之旅](/solutions/capital-one) | 已翻译 |
| [案例研究](/solutions/case-studies) | 已翻译 |
| [命令行界面（CLI）](/solutions/clis) | 已翻译 |
| [Go 与云计算及网络服务](/solutions/cloud) | 已翻译 |
| [Go 中的平滑升级](/solutions/cloudflare) | 已翻译 |
| [Cockroach Labs：为什么选择用 Go 构建数据库](/solutions/cockroachlabs) | 已翻译 |
| [Curve 如何借助 Go 取得领先](/solutions/curve) | 已翻译 |
| [开发运维与站点可靠性工程](/solutions/devops) | 已翻译 |
| [Dropbox：开源我们的 Go 库](/solutions/dropbox) | 已翻译 |
| [Facebook 如何用 Go 构建实体框架](/solutions/facebook) | 已翻译 |
| [Chrome 内容优化服务运行在 Go 之上](/solutions/google/chrome) | 已翻译 |
| [Google Core Data Solutions 团队如何使用 Go](/solutions/google/coredata) | 已翻译 |
| [Firebase Hosting 团队如何借助 Go 扩展规模](/solutions/google/firebase) | 已翻译 |
| [Google 的 Go 实践](/solutions/google/) | 已翻译 |
| [驱动 Google 生产环境：Google 站点可靠性工程团队如何使用 Go](/solutions/google/sitereliability) | 已翻译 |
| [Bigslice：用 Go 编写的集群计算系统](/solutions/grail) | 已翻译 |
| [为什么选择 Go](/solutions/) | 已翻译 |
| [MercadoLibre 与 Go 共同成长](/solutions/mercadolibre) | 已翻译 |
| [Microsoft 如何拥抱 Go](/solutions/microsoft) | 已翻译 |
| [Monzo：用 Go、微服务和容器构建银行](/solutions/monzo) | 已翻译 |
| [Netflix：使用 SSD 缓存应用数据](/solutions/netflix) | 已翻译 |
| [PayPal 借助 Go 实现现代化与规模扩展](/solutions/paypal) | 已翻译 |
| [Riot Games：使用 Go 开发和运维游戏](/solutions/riotgames) | 已翻译 |
| [Salesforce：从 Python/C 转向 Go](/solutions/salesforce) | 已翻译 |
| [深入了解 SIXT 的 Go 实践](/solutions/sixt) | 已翻译 |
| [Stream：为什么我们从 Python 转向 Go](/solutions/stream) | 已翻译 |
| [Trivago：为什么我们选择 Go](/solutions/trivago) | 已翻译 |
| [Twitch：Go 向低延迟垃圾回收迈进](/solutions/twitch) | 已翻译 |
| [Uber：用 Go 编写的 GPU 加速分析引擎](/solutions/uber) | 已翻译 |
| [应用场景](/solutions/use-cases) | 已翻译 |
| [Go 与 Web 开发](/solutions/webdev) | 已翻译 |
| [Wildlife Studios 如何用 Go 构建后端系统](/solutions/wildlifestudios) | 已翻译 |
| [X：实时处理每天 50 亿次会话](/solutions/x) | 已翻译 |
| [服务条款](/tos) | 已翻译 |
| [包、变量与函数](/tour/basics/1) | 已翻译 |
| [并发](/tour/concurrency/1) | 已翻译 |
| [流程控制：for、if、else、switch 和 defer](/tour/flowcontrol/1) | 已翻译 |
| [泛型](/tour/generics/1) | 已翻译 |
| [方法与接口](/tour/methods/1) | 已翻译 |
| [更多类型：结构体、切片与映射](/tour/moretypes/1) | 已翻译 |
| [欢迎！](/tour/welcome/1) | 已翻译 |
| [注释](/wiki/Comments) | 已翻译 |
| [Wiki 首页](/wiki/) | 已翻译 |

<!-- translation-coverage:end -->

当前中文参考包括 180 个公开标准库包（包含 builtin）的 API 正文、字段及示例注释，完整 go 命令手册与参数说明、完整语言规范和 Go 内存模型，以及 21 项独立工具文档、汇编指南和 GODEBUG 参考。原文从 Go 1.27.1 工具链导出；包参考采用 linux/amd64 视图。另有 arena 实验包的概述译文，默认构建视图未导出其 API 声明。relnote 的上游导出页面只有目录标题。内部实现包、其他平台专属 API、需额外构建配置的包（如 runtime/cgo、BoringCrypto、SIMD）、外部子仓库和第三方包文档仍使用原文。
具体范围、阅读顺序和版本说明见[中文参考文档阅读入口](/doc/reference-zh)。升级本机 Go 后，中文快照不会自动更新；页面顶部标明快照与本机工具链版本。

本仓库中的网站介绍、帮助、行为准则、品牌指南、条款、应用场景、企业案例正文与案例卡片均提供中文。只包含外链的案例卡片翻译标题和简介，外部文章不计为全文翻译。Wiki 仅翻译仓库内的首页与注释页面，外部完整 Wiki、gopls 仓库和其他链接资源仍为原文。
程序逻辑、字符串、终端输出、本机源码和部分图片内嵌文字保留原文。博客、历史演讲、历史版本说明和早期更新记录仍按约定保留英文。
外部网站、下载链接和在线 Playground 仍需要联网。

## 术语约定 {#terminology}

| 英文 | 本译本用语 |
| --- | --- |
| module / package | 模块 / 包 |
| dependency / toolchain | 依赖 / 工具链 |
| array / slice / map | 数组 / 切片 / 映射（map） |
| struct / interface | 结构体 / 接口 |
| receiver / method set | 接收者 / 方法集 |
| type assertion / type parameter | 类型断言 / 类型参数（类型形参） |
| type argument | 类型实参 |
| type constraint / type set | 类型约束 / 类型集 |
| goroutine | goroutine（轻量级协程），不与操作系统线程混用 |
| channel | 通道（channel） |
| concurrency / parallelism | 并发 / 并行 |
| zero value / underlying type | 零值 / 底层类型 |
| garbage collection / runtime | 垃圾回收 / 运行时 |
| data race / race condition | 数据竞争 / 竞态条件 |
| happens-before / synchronizes-before | 先行发生 / 同步先于 |
| sequenced-before / sequential consistency | 程序顺序先于 / 顺序一致性 |
| atomic operation / memory ordering | 原子操作 / 内存顺序 |
| fuzzing / benchmark | 模糊测试 / 基准测试 |
| backward compatibility | 向后兼容性 |
| major / minor / patch version | 主版本 / 次版本 / 补丁版本 |

`nil`、`panic`、`recover`、`defer` 等代码标识符保留英文，在正文中解释其语义。
原文存在示例问题时，会用“译注”单独说明，不将说明混入原始代码。
