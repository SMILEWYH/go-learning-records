---
title: Go 开发者安全最佳实践
layout: article
---

[返回 Go 安全](/security)

本页介绍 Go 开发者应优先考虑的安全实践。从使用模糊测试自动发现问题，到检查并发访问中的竞争，这些方法有助于提高代码的安全性和可靠性。

## 扫描源码和二进制文件中的漏洞 {#scan-source-code-and-binaries-for-vulnerabilities}

定期扫描源码和二进制文件，有助于尽早发现潜在安全风险。可以使用基于 [Go 漏洞数据库](https://pkg.go.dev) 的 [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) 扫描漏洞，并分析哪些漏洞会实际影响你的代码。入门参见 [govulncheck 教程](/doc/tutorial/govulncheck)。

govulncheck 也可以集成到 CI/CD 流程中。Go 团队在 GitHub Marketplace 提供了 [govulncheck GitHub Action](https://github.com/marketplace/actions/golang-govulncheck-action)。govulncheck 还支持 `-json` 标志，便于与其他 CI/CD 系统集成。

也可以通过 [Visual Studio Code 的 Go 扩展](/security/vuln/editor)，直接在编辑器中扫描漏洞。入门参见[此教程](/doc/tutorial/govulncheck-ide)。

## 及时更新 Go 版本和依赖 {#keep-your-go-version-and-dependencies-up-to-date}

[及时更新 Go](/doc/install) 可以获得新的语言特性、性能改进，以及已知安全漏洞的修复。较新的 Go 版本也能更好地兼容新版依赖，避免潜在的集成问题。[Go 发布历史](/doc/devel/release)列出了各版本之间的变化。Go 团队在发布周期内通过补丁版本修复安全缺陷，请更新到相应发行系列的最新补丁版本，以获取安全修复。

及时更新第三方依赖，同样关系到安全、性能以及与 Go 生态中最新标准的兼容。不过，未经充分检查就升级到最新版本[也可能带来风险](https://research.swtch.com/npm-colors)，例如引入新缺陷、不兼容变更，甚至恶意代码。因此，虽然获得安全修复和改进很重要，每次依赖更新仍应仔细审查和测试。

## 通过模糊测试发现边界情况中的漏洞 {#test-with-fuzzing-to-uncover-edge-case-exploits}

[模糊测试](/security/fuzz)是一种自动化测试方法，通过覆盖率引导不断变异随机输入、探索代码路径，查找和报告 SQL 注入、缓冲区溢出、拒绝服务、跨站脚本等潜在漏洞。它常能发现开发者遗漏，或认为出现概率太低而未测试的边界情况。入门参见[此教程](/doc/tutorial/fuzz)。

## 使用 Go 竞争检测器检查并发访问 {#check-for-race-conditions-with-gos-race-detector}

当两个或更多 [goroutine](/tour/concurrency/1) 并发访问同一资源，且至少有一次写入，又没有适当同步时，就可能发生数据竞争。这会导致难以预测、难以诊断的问题。Go 内置的[竞争检测器](/doc/articles/race_detector)可以帮助发现这类问题，提高并发程序的安全性和可靠性。不过，它只能发现运行时实际发生的竞争，无法检测未执行到的代码路径。

运行测试或构建应用时，添加 `-race` 标志即可启用，例如 `go test -race`。这会在编译时加入竞争检测，并报告运行期间发现的数据竞争。发现问题时，检测器会[打印报告](/doc/articles/race_detector#report-format)，其中包含冲突访问的调用栈，以及相关 goroutine 创建位置的调用栈。

> 译注：原文在这一节使用 race condition 一词。Go 的 `-race` 主要检测 data race（数据竞争）；“竞态条件”的含义更广，并非所有竞态条件都能由它检测出来。

## 使用 vet 检查可疑代码 {#use-vet-to-examine-suspicious-constructs}

Go 的 [vet 命令](https://pkg.go.dev/cmd/vet)分析源码，指出不一定属于语法错误、却可能导致运行时问题的可疑写法。原文举例包括不可达代码、未使用变量以及 goroutine 的常见误用。在开发早期发现问题，可以保持代码质量、减少调试时间，提高软件可靠性。对整个项目运行：

```
go vet ./...
```

> 译注：Go 中未使用的局部变量通常直接由编译器报错；`go vet` 的具体检查项目应以当前工具链的帮助信息为准。

## 订阅 golang-announce，接收安全版本通知 {#subscribe-to-golang-announce-for-notification-of-security-releases}

包含安全修复的 Go 版本会提前在邮件量较少的 [golang-announce@googlegroups.com](https://groups.google.com/group/golang-announce) 邮件列表中预告。如果希望了解 Go 自身的安全修复何时发布，可以订阅该列表。
