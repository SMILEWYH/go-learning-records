---
title: "Go、向后兼容性与 GODEBUG"
layout: article
---

<!--
This document is kept in the Go repo, not x/website,
because it documents the full list of known GODEBUG settings,
which are tied to a specific release.
-->

## 引言 {#intro}

重视向后兼容性是 Go 的主要优势之一。不过，有时无法保持完全兼容。如果代码依赖错误行为，包括不安全的行为，那么修复缺陷就会破坏这些代码。新特性也可能产生类似影响：例如，HTTP 客户端启用 HTTP/2 后，连接到存在 HTTP/2 实现缺陷的服务器的程序可能无法正常工作。这类变更不可避免，也[符合 Go 1 兼容性规则](/doc/go1compat)。即便如此，Go 仍提供名为 GODEBUG 的机制，减少开发者使用新工具链编译旧代码时受到的影响。

一项 GODEBUG 设置是一个 `key=value` 键值对，用于控制 Go 程序某些部分的执行行为。环境变量 `GODEBUG` 可以包含以逗号分隔的设置列表。例如，如果 Go 程序的运行环境包含：

	GODEBUG=http2client=0,http2server=0

那么，该程序的 HTTP 客户端和 HTTP 服务器默认都会禁用 HTTP/2。`GODEBUG` 环境变量中无法识别的设置会被忽略。也可以为特定程序设置默认 `GODEBUG` 值，详见下文。

准备引入符合 Go 1 兼容性规则、但仍可能破坏现有程序的变更时，我们首先会尽可能让更多现有程序继续正常工作。对于剩余程序，则定义一项新的 GODEBUG 设置，允许它们选择恢复旧行为。如果实现这种设置不可行，也可能不提供，但这种情况应当极为罕见。

为兼容性添加的 GODEBUG 设置，至少会保留两年，即四个 Go 版本。某些设置，例如 `http2client` 和 `http2server`，会保留更长时间，甚至无限期保留。

在可行的情况下，每项 GODEBUG 设置都关联一个 [runtime/metrics](/pkg/runtime/metrics/) 计数器，名称为 `/godebug/non-default-behavior/<name>:events`，统计程序因该设置使用非默认值而改变行为的次数。例如，设置 `GODEBUG=http2client=0` 后，`/godebug/non-default-behavior/http2client:events` 会统计程序中配置为不支持 HTTP/2 的 HTTP Transport 数量。

## GODEBUG 的默认值 {#default}

如果某项 GODEBUG 设置没有出现在环境变量中，其值由三个来源依次决定：首先采用构建程序所用 Go 工具链的默认值，然后根据 `go.mod` 中声明的 Go 版本进行调整，最后由程序中显式的 `//go:debug` 指令覆盖。

[GODEBUG 变更记录](#history)给出了各工具链版本的具体默认值。例如，Go 1.21 引入 `panicnil`，控制是否允许 `panic(nil)`。默认值为 `panicnil=0`，使 `panic(nil)` 产生运行时错误；设置 `panicnil=1` 可恢复 Go 1.20 及更早版本的行为。

编译声明了较早 Go 版本的主模块或工作区时，工具链会调整默认值，尽可能匹配该旧版本。例如，Go 1.21 工具链编译程序时，如果主模块的 `go.mod` 或工作区的 `go.work` 声明 `go` `1.20`，程序就默认采用 `panicnil=1`，与 Go 1.20 保持一致，而不是采用 Go 1.21 的默认行为。

有一个例外：为安全修复版本引入的 GODEBUG 设置，会对所有版本应用新行为。

由于这种默认值设置机制是在 Go 1.21 中才引入的，声明 Go 1.20 之前版本的程序，会配置为匹配 Go 1.20，而不是更早的版本。

从 Go 1.23 开始，可以在主模块的 `go.mod` 或工作区的 `go.work` 中列出一个或多个 `godebug` 指令，覆盖这些默认值：

	godebug (
		default=go1.21
		panicnil=1
		asynctimerchan=0
	)

特殊键 `default` 指定从哪个 Go 版本获取未显式设置项的默认值，使 GODEBUG 默认值可以与模块的 Go 语言版本分开设置。本例要求使用 Go 1.21 的语义，同时采用 Go 1.21 之前的 `panic(nil)` 旧行为，以及 Go 1.23 的 `asynctimerchan=0` 新行为。

只有主模块的 `go.mod` 中的 `godebug` 指令会被读取，依赖模块中的指令都会被忽略。`godebug` 中出现无法识别的设置会报错。Go 1.23 之前的工具链完全不认识 `godebug`，因此会拒绝所有此类指令。使用工作区时，工具链忽略各 `go.mod` 中的 `godebug` 指令，改为读取 `go.work`。

`go` 和 `godebug` 指令产生的默认值，适用于构建的所有 main 包。需要更细粒度的控制时，从 Go 1.21 开始，可以在 main 包源文件的顶部、`package` 语句之前添加一个或多个 `//go:debug` 指令。上例可写为：

	//go:debug default=go1.21
	//go:debug panicnil=1
	//go:debug asynctimerchan=0

从 Go 1.21 开始，若 `//go:debug` 指令包含无法识别的 GODEBUG 设置，工具链会认为程序无效。同一设置出现多条 `//go:debug` 指令，也会使程序无效。较早的工具链则完全忽略这些指令。

以下命令会报告将编译进 main 包的默认值：

	go list -f '{{.DefaultGODEBUG}}' my/main/package

只显示与基础 Go 工具链默认值不同的设置。

测试包时，`*_test.go` 文件中的 `//go:debug` 会被视为测试 main 包的指令。其他上下文中的 `//go:debug` 会被工具链忽略，`go` `vet` 会报告这些指令的位置不正确。

## GODEBUG 变更记录 {#history}

本节记录各 Go 主版本中为兼容性新增和移除的 GODEBUG 设置。包或程序还可能定义用于内部调试的其他设置，例如[运行时文档](/pkg/runtime#hdr-Environment_Variables)和 [go 命令文档](/cmd/go#hdr-Build_and_test_caching)中的设置。

### Go 1.27

Go 1.27 移除了 `gotypesalias`，参见 [Go 1.22](#go-122) 一节。

Go 1.27 移除了 `tlsunsafeekm`，参见 [Go 1.22](#go-122) 一节。

Go 1.27 移除了 `tlsrsakex`，参见 [Go 1.22](#go-122) 一节。

Go 1.27 移除了 `tls3des`，参见 [Go 1.23](#go-123) 一节。

Go 1.27 移除了 `tls10server`，参见 [Go 1.22](#go-122) 一节。

Go 1.27 移除了 `x509keypairleaf`，参见 [Go 1.23](#go-123) 一节。

Go 1.27 移除了 `asynctimerchan`，参见 [Go 1.23](#go-123) 一节。

Go 1.27 新增 `htmlmetacontenturlescape`，控制 html/template 是否转义 HTML meta 标签的 content 属性中 `url=` 部分的 URL。默认值 `htmlmetacontenturlescape=1` 会进行转义，设为 `htmlmetacontenturlescape=0` 则禁用。为防止内容注入攻击，该设置及默认行为已回移到 Go 1.25.8 和 Go 1.26.1。*译注：原文后两处设置名漏写了 `url`；这里依据本项目使用的 Go 源码 `internal/godebugs/table.go` 和 `html/template/escape.go` 更正。*

Go 1.27 将 [Go 1.26](#go-126) 引入的 `tracebacklabels` 默认值改为 `1`。预计会无限期保留关闭此功能的设置，以应对 goroutine 标签中可能包含不应出现在栈追踪中的敏感信息的情况。

Go 1.27 新增 `x509sslcertoverrideplatform`，控制在 Windows 和 Darwin 上设置了 `SSL_CERT_FILE` 或 `SSL_CERT_DIR` 时，crypto/x509 是否从磁盘加载根证书。默认值 `x509sslcertoverrideplatform=1` 会遵循这些环境变量并从磁盘加载；设为 `x509sslcertoverrideplatform=0` 则忽略这些变量，改用平台证书存储。计划在 Go 1.31 中移除此设置。

Go 1.27 新增 `fips140ems`。设为 `0` 时，在 FIPS 140-3 模式下不再强制要求扩展主密钥（Extended Master Secret）。默认行为不变。该设置已回移到 Go 1.26.6 和 Go 1.25.13，计划在 Go 1.31 中移除。

### Go 1.26

Go 1.26 新增 `httpcookiemaxnum`，控制 net/http 解析 HTTP 头部时接受的 Cookie 数量上限。如果某个头部中的 Cookie 数量超过该值，解析会提前失败。默认值为 `httpcookiemaxnum=3000`，设为 `httpcookiemaxnum=0` 则不限制数量。为防止拒绝服务攻击，该设置及默认行为已回移到 Go 1.25.2 和 Go 1.24.8。

Go 1.26 新增 `urlmaxqueryparams`，控制 net/url 解析 URL 编码查询字符串时接受的查询参数数量上限。如果超过该值，解析会提前失败。默认值为 `urlmaxqueryparams=10000`，设为 `urlmaxqueryparams=0` 则禁用限制。为防止拒绝服务攻击，该设置及默认行为已回移到 Go 1.25.6 和 Go 1.24.12。

Go 1.26 新增 `urlstrictcolons`，控制 `net/url.Parse` 是否允许主机名在方括号包围的 IPv6 地址之外出现非法冒号。默认值 `urlstrictcolons=1` 会拒绝 `http://localhost:1:2`、`http://::1/` 等 URL。方括号中的 IPv6 地址仍允许包含冒号，例如 `http://[::1]/`。

Go 1.26 又启用了两种后量子密钥交换机制：SecP256r1MLKEM768 和 SecP384r1MLKEM1024。可以使用 [`tlssecpmlkem` 设置](/pkg/crypto/tls/#Config.CurvePreferences)恢复原来的默认行为。

Go 1.26 新增 `tracebacklabels`，控制是否在栈追踪中包含通过 `runtime/pprof` 设置的 goroutine 标签。设为 `tracebacklabels=1` 时，会在运行时栈追踪，以及 debug=2 的 runtime/pprof 栈转储的 goroutine 状态头部中显示这些键值对。该格式将来可能变化，参见 go.dev/issue/76349。

Go 1.26 新增 `cryptocustomrand`，控制大多数 crypto/... API 是否忽略随机源 `io.Reader` 参数。Go 1.26 默认采用 `cryptocustomrand=0`，即忽略随机源参数；设为 `cryptocustomrand=1` 则恢复 Go 1.26 之前的行为。

### Go 1.25

Go 1.25 新增 `decoratemappings`，控制 Go 运行时是否为操作系统的匿名内存映射添加用途说明。这些注记在 /proc/self/maps 和 /proc/self/smaps 中显示为 "[anon: Go: ...]"。该设置仅用于 Linux。Go 1.25 默认采用 `decoratemappings=1`，启用注记；设为 `decoratemappings=0` 则恢复此前行为。此设置在程序启动时确定，启动后修改 `GODEBUG` 环境变量不会改变它。

Go 1.25 新增 `embedfollowsymlinks`，控制 go 命令嵌入文件时，是否跟随指向普通文件的符号链接。默认值 `embedfollowsymlinks=0` 不允许跟随，`embedfollowsymlinks=1` 则允许。

Go 1.25 新增 `containermaxprocs`，控制 Go 运行时设置默认 GOMAXPROCS 时是否考虑 cgroup CPU 限制。默认值 `containermaxprocs=1` 会同时考虑逻辑 CPU 总数、CPU 亲和性和 cgroup 限制；`containermaxprocs=0` 则不考虑 cgroup 限制。此设置仅影响 Linux。

Go 1.25 新增 `updatemaxprocs`，控制 Go 运行时是否根据新的 CPU 亲和性或 cgroup 限制定期更新 GOMAXPROCS。默认值 `updatemaxprocs=1` 启用定期更新，`updatemaxprocs=0` 禁用。

Go 1.25 根据 RFC 9155，在 TLS 1.2 中禁用了 SHA-1 签名算法。设置 `tlssha1=1` 可以恢复原来的默认行为。

Go 1.25 在 crypto/x509.CreateCertificate 中改用 SHA-256 填充缺失的 SubjectKeyId。设置 `x509sha256skid=0` 可以恢复使用 SHA-1。

Go 1.25 修正了运行时内部锁的竞争报告语义，因此移除了 [`runtimecontentionstacks` 设置](/pkg/runtime#hdr-Environment_Variables)。

Go 1.25 从 RC 2 开始，出于防范版本控制系统（VCS）注入攻击的考虑，在检测到多个 VCS 时不再写入构建信息。此行为和设置已回移到 Go 1.24.5 和 Go 1.23.11。设置 `allowmultiplevcs=1` 可以重新启用这一场景下的构建信息写入。

### Go 1.24

Go 1.24 新增 `fips140`，控制 Go 密码模块是否以 FIPS 140-3 模式运行。可选值为：

- "off"：不特别支持 FIPS 140-3 模式，默认值。
- "on"：Go 密码模块以 FIPS 140-3 模式运行。
- "only"：类似 "on"，但未获 FIPS 140-3 批准的密码算法会返回错误或发生 panic。

详情见 [FIPS 140-3 合规性](/doc/security/fips140)。此设置在程序启动时确定，启动后修改 `GODEBUG` 环境变量不会改变它。

Go 1.24 将全局函数 [`math/rand.Seed`](/pkg/math/rand/#Seed) 改为空操作，由 `randseednop` 控制。Go 1.24 默认采用 `randseednop=1`，设为 `randseednop=0` 则恢复此前行为。

Go 1.24 为 `multipathtcp` 增加了新取值。现在可选值为：

- "0"：拨号器和监听器都默认禁用 MPTCP。
- "1"：拨号器和监听器都默认启用 MPTCP。
- "2"：只有监听器默认启用 MPTCP。
- "3"：只有拨号器默认启用 MPTCP。

Go 1.24 的默认值为 multipathtcp="2"，因此监听器默认启用。使用 multipathtcp="0" 可以恢复此前行为。

Go 1.24 改变了 `go test -json` 的行为，使构建错误以 JSON 而不是文本输出。新的 JSON 事件使用新的 `Action` 值区分，但仍可能影响无法稳健处理这些事件的 CI 系统。可以使用 `gotestjsonbuildtext` 控制此行为，`gotestjsonbuildtext=1` 恢复 Go 1.23 的行为。此设置将在未来版本中移除，最早为 Go 1.28。

Go 1.24 要求 [`crypto/rsa`](/pkg/crypto/rsa) 的 RSA 密钥至少为 1024 位。此行为由 `rsa1024min` 控制，`rsa1024min=0` 恢复 Go 1.23 的行为。

Go 1.24 在 [`crypto/subtle`](/pkg/crypto/subtle) 中引入了启用平台专用数据无关时序（DIT）模式的机制。可以通过 `dataindependenttiming` 为整个程序启用该模式。Go 1.24 默认采用 `dataindependenttiming=0`；未设置时，默认行为与 Go 1.23 一致。`dataindependenttiming=1` 会为整个 Go 程序启用 DIT。启用后，从 Go 调用 C 时也会保持 DIT 开启；从 C 调入 Go 代码时，会启用 DIT，如果进入 Go 之前未启用，则在返回 C 前禁用。目前仅影响 arm64 程序，对其他平台不起作用。

Go 1.24 移除了 `x509sha1`。`crypto/x509` 不再支持验证使用 SHA-1 签名算法的证书签名。

Go 1.24 将 [`x509usepolicies` 设置](/pkg/crypto/x509/#CreateCertificate)的默认值从 `0` 改为 `1`。序列化证书时，默认从 [`Certificate.Policies`](/pkg/crypto/x509/#Certificate.Policies) 字段读取策略，而不是从 [`Certificate.PolicyIdentifiers`](/pkg/crypto/x509/#Certificate.PolicyIdentifiers) 读取。

Go 1.24 默认启用了后量子密钥交换机制 X25519MLKEM768。可以通过 [`tlsmlkem` 设置](/pkg/crypto/tls/#Config.CurvePreferences)恢复原来的默认行为。对于无法正确处理较大记录、导致握手超时的缺陷 TLS 服务器，这可能有帮助，参见 [TLS 后量子 TL;DR 故障](https://tldr.fail/)。Go 1.24 同时移除了 X25519Kyber768Draft00 和 Go 1.23 的 `tlskyber` 设置。

Go 1.24 让 [`ParsePKCS1PrivateKey`](/pkg/crypto/x509/#ParsePKCS1PrivateKey) 使用并验证编码私钥中的 CRT 参数。此行为由 `x509rsacrt` 控制，`x509rsacrt=0` 恢复 Go 1.23 的行为。

### Go 1.23

Go 1.23 将 time 包创建的通道改为无缓冲的同步通道，使正确使用 [`Timer.Stop`](/pkg/time/#Timer.Stop) 和 [`Timer.Reset`](/pkg/time/#Timer.Reset) 的返回结果变得更容易。[`asynctimerchan` 设置](/pkg/time/#NewTimer)可以禁用这一变更。此变更没有对应的运行时指标。该设置将在 Go 1.27 中移除。

Go 1.23 改变了 [`os.Lstat`](/pkg/os#Lstat) 和 [`os.Stat`](/pkg/os#Stat) 为重解析点报告的模式位，可通过 `winsymlink` 控制。从 Go 1.23 开始（`winsymlink=1`），挂载点不再设置 [`os.ModeSymlink`](/pkg/os#ModeSymlink)；既不是符号链接、Unix 套接字，也不是数据去重文件的重解析点，现在始终设置 [`os.ModeIrregular`](/pkg/os#ModeIrregular)。因此，[`filepath.EvalSymlinks`](/pkg/path/filepath#EvalSymlinks) 不再解析挂载点，消除了许多不一致和缺陷的来源。在此前版本中（`winsymlink=0`），挂载点被视为符号链接；具有非默认 [`os.ModeType`](/pkg/os#ModeType) 位（例如 [`os.ModeDir`](/pkg/os#ModeDir)）的其他重解析点，不会设置 `ModeIrregular` 位。

Go 1.23 修改了 [`os.Readlink`](/pkg/os#Readlink) 和 [`filepath.EvalSymlinks`](/pkg/path/filepath#EvalSymlinks)，不再尝试将卷规范化为盘符，因为这种转换并不总是可行。此行为由 `winreadlinkvolume` 控制。Go 1.23 默认采用 `winreadlinkvolume=1`，此前版本默认采用 `winreadlinkvolume=0`。

Go 1.23 默认启用了实验性的后量子密钥交换机制 X25519Kyber768Draft00。可以通过 [`tlskyber` 设置](/pkg/crypto/tls/#Config.CurvePreferences)恢复原来的默认行为。对于无法正确处理较大记录、导致握手超时的缺陷 TLS 服务器，这可能有帮助，参见 [TLS 后量子 TL;DR 故障](https://tldr.fail/)。

Go 1.23 修改了 [crypto/x509.ParseCertificate](/pkg/crypto/x509/#ParseCertificate)，使其拒绝负数序列号。可以通过 [`x509negativeserial` 设置](/pkg/crypto/x509/#ParseCertificate)恢复旧行为。

Go 1.23 默认重新启用了 html/template 对 ECMAScript 6 模板字面量的支持。[`jstmpllitinterp` 设置](/pkg/html/template#hdr-Security_Model)不再产生任何效果。

Go 1.23 修改了客户端和服务器在未显式配置时使用的默认 TLS 密码套件，移除了 3DES 套件。可以通过 [`tls3des` 设置](/pkg/crypto/tls/#Config.CipherSuites)恢复原来的默认行为。此设置将在 Go 1.27 中移除。

Go 1.23 修改了 [`tls.X509KeyPair`](/pkg/crypto/tls#X509KeyPair) 和 [`tls.LoadX509KeyPair`](/pkg/crypto/tls#LoadX509KeyPair)，使其填充返回的 [`tls.Certificate`](/pkg/crypto/tls#Certificate) 的 Leaf 字段。此行为由 `x509keypairleaf` 控制。Go 1.23 默认采用 `x509keypairleaf=1`，此前版本默认采用 `x509keypairleaf=0`。此设置将在 Go 1.27 中移除。

Go 1.23 修改了 [`net/http.ServeContent`](/pkg/net/http#ServeContent)、[`net/http.ServeFile`](/pkg/net/http#ServeFile) 和 [`net/http.ServeFS`](/pkg/net/http#ServeFS)，使其在返回错误响应时移除 Cache-Control、Content-Encoding、Etag 和 Last-Modified 头部。此行为由 [`httpservecontentkeepheaders` 设置](/pkg/net/http#ServeContent)控制。设置 `httpservecontentkeepheaders=1` 可恢复 Go 1.23 之前的行为。

### Go 1.22

Go 1.22 为 TLS 握手中可接受的 RSA 密钥大小增加了可配置上限，由 [`tlsmaxrsasize` 设置](/pkg/crypto/tls#Conn.Handshake)控制。默认值 tlsmaxrsasize=8192 将 RSA 密钥限制为不超过 8192 位。为防止拒绝服务攻击，该设置及默认行为已回移到 Go 1.19.13、Go 1.20.8 和 Go 1.21.1。

Go 1.22 将 net/http 客户端或服务器读取到的请求或响应包含空 Content-Length 头部的情况视为错误。此行为由 `httplaxcontentlength` 控制。

Go 1.22 修改了 ServeMux，使其接受扩展的路由模式，并按路径段对模式和请求路径进行反转义。此行为由 [`httpmuxgo121` 设置](/pkg/net/http/#ServeMux)控制。

Go 1.22 为 [go/types](/pkg/go/types) 添加了 [Alias 类型](/pkg/go/types#Alias)，用于显式表示[类型别名](/ref/spec#Type_declarations)。类型检查器是否生成 `Alias`，由 [`gotypesalias` 设置](/pkg/go/types#Alias)控制。Go 1.22 默认采用 `gotypesalias=0`；Go 1.23 将以 `gotypesalias=1` 为默认值。此设置将在 Go 1.27 中移除。

Go 1.22 将服务器和客户端默认支持的最低 TLS 版本改为 TLS 1.2。可以通过 [`tls10server` 设置](/pkg/crypto/tls/#Config)恢复为 TLS 1.0。此设置将在 Go 1.27 中移除。

Go 1.22 修改了客户端和服务器在未显式配置时使用的默认 TLS 密码套件，移除了基于 RSA 密钥交换的套件。可以通过 [`tlsrsakex` 设置](/pkg/crypto/tls/#Config)恢复原来的默认行为。此设置将在 Go 1.27 中移除。

Go 1.22 在连接既不支持 TLS 1.3、也不支持扩展主密钥（Go 1.21 中实现）时，禁用了 [`ConnectionState.ExportKeyingMaterial`](/pkg/crypto/tls/#ConnectionState.ExportKeyingMaterial)。可以通过 [`tlsunsafeekm` 设置](/pkg/crypto/tls/#ConnectionState.ExportKeyingMaterial)重新启用。此设置将在 Go 1.27 中移除。

Go 1.22 改变了运行时在 Linux 上与透明大页的交互方式。某种常见的 Linux 内核默认配置会带来显著的额外内存开销，Go 1.22 不再规避该默认配置。若不调整内核设置，可以通过 [`disablethp` 设置](/pkg/runtime#hdr-Environment_Variables)为 Go 内存禁用透明大页。该行为已回移到 Go 1.21.1，但设置本身从 Go 1.21.6 才可用。未来版本可能移除此设置；受到影响的用户应按照 [GC 指南](/doc/gc-guide#Linux_transparent_huge_pages)调整 Linux 配置，或改用默认完全禁用透明大页的 Linux 发行版。

Go 1.22 将运行时内部锁的竞争纳入 [`mutex` 剖析](/pkg/runtime/pprof#Profile)。这些锁的竞争始终记在 `runtime._LostContendedRuntimeLock` 下。可以通过 [`runtimecontentionstacks` 设置](/pkg/runtime#hdr-Environment_Variables)启用运行时锁的完整栈追踪。这些栈追踪具有非标准语义，详情见设置文档。

Go 1.22 为 [`crypto/x509.Certificate`](/pkg/crypto/x509/#Certificate) 添加了 [`Policies`](/pkg/crypto/x509/#Certificate.Policies) 字段，支持组成部分超过 31 位的证书策略 OID。默认情况下，仅在解析时将策略 OID 填入此字段，序列化时不使用。通过 [`x509usepolicies` 设置](/pkg/crypto/x509/#CreateCertificate)，可以改用此字段序列化较大的 OID，而不是使用原有的 PolicyIdentifiers 字段。

### Go 1.21

Go 1.21 将使用 nil 接口值调用 `panic` 的情况改为运行时错误，由 [`panicnil` 设置](/pkg/builtin/#panic)控制。

Go 1.21 将 html/template 动作出现在 ECMAScript 6 模板字面量内部的情况视为错误，由 [`jstmpllitinterp` 设置](/pkg/html/template#hdr-Security_Model)控制。此行为已回移到 Go 1.19.8+ 和 Go 1.20.3+。

Go 1.21 限制了 MIME 头部和 multipart 表单分部的最大数量，分别由 [`multipartmaxheaders` 和 `multipartmaxparts` 设置](/pkg/mime/multipart#hdr-Limits)控制。此行为已回移到 Go 1.19.8+ 和 Go 1.20.3+。

Go 1.21 增加了多路径 TCP 支持，但只有应用显式请求时才使用。此行为由 [`multipathtcp` 设置](/pkg/net#Dialer.SetMultipathTCP)控制。

目前没有移除上述设置的计划。

### Go 1.20

Go 1.20 增加了拒绝 tar 和 zip 归档中不安全路径的能力，分别由 [`tarinsecurepath` 设置](/pkg/archive/tar/#Reader.Next)和 [`zipinsecurepath` 设置](/pkg/archive/zip/#NewReader)控制。默认值为 `tarinsecurepath=1` 和 `zipinsecurepath=1`，保留较早 Go 版本的行为。未来版本可能将默认值改为 `tarinsecurepath=0` 和 `zipinsecurepath=0`。

Go 1.20 为 [`math/rand`](/pkg/math/rand) 全局随机数生成器引入自动设置种子的行为，由 [`randautoseed` 设置](/pkg/math/rand/#Seed)控制。

Go 1.20 引入了证书验证时使用的备用根证书概念，由 [`x509usefallbackroots` 设置](/pkg/crypto/x509/#SetFallbackRoots)控制。

Go 1.20 从发行版中移除了标准库预安装的 `.a` 文件。安装后，标准库现在像其他模块的包一样进行构建和缓存。[`installgoroot` 设置](/cmd/go#hdr-Compile_and_install_packages_and_dependencies)可以恢复安装和使用预安装 `.a` 文件的行为。

目前没有移除上述设置的计划。

### Go 1.19

Go 1.19 将通过 PATH 查找可执行文件时解析到当前目录中的二进制文件视为错误，由 [`execerrdot` 设置](/pkg/os/exec#hdr-Executables_in_the_current_directory)控制。目前没有移除此设置的计划。

Go 1.19 开始在 DNS 请求中发送 EDNS0 附加头部。据报告，这可能导致某些路由器内置的 DNS 服务器工作异常，例如 CenturyLink Zyxel C3000Z。可以通过 [`netedns0` 设置](/pkg/net#hdr-Name_Resolution)改变这一行为。此设置在 Go 1.21.12、Go 1.22.5、Go 1.23 及后续版本中可用，目前没有移除计划。

### Go 1.18

Go 1.18 移除了对大多数 X.509 证书中 SHA1 的支持，由 [`x509sha1` 设置](/pkg/crypto/x509#InsecureAlgorithmError)控制。此设置已在 Go 1.24 中移除。

### Go 1.10

Go 1.10 修改了构建缓存的工作方式，并新增测试缓存，同时引入 [`gocacheverify`、`gocachehash` 和 `gocachetest` 设置](/cmd/go/#hdr-Build_and_test_caching)。目前没有移除这些设置的计划。

### Go 1.6

Go 1.6 引入了对 HTTP/2 的透明支持，由 [`http2client`、`http2server` 和 `http2debug` 设置](/pkg/net/http/#hdr-HTTP_2)控制。目前没有移除这些设置的计划。

### Go 1.5

Go 1.5 引入了纯 Go DNS 解析器，由 [`netdns` 设置](/pkg/net/#hdr-Name_Resolution)控制。目前没有移除此设置的计划。
