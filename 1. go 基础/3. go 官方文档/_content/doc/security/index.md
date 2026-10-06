---
title: 安全
layout: article
---

本页为 Go 开发者提供提高项目安全性的相关资源。

（另请参阅：[Go 开发者安全最佳实践](/security/best-practices)。）

## 查找并修复已知漏洞 {#find-and-fix-known-vulnerabilities}

Go 的漏洞检测旨在提供误报较少、可靠的工具，帮助开发者了解可能影响项目的已知漏洞。要了解 Go 漏洞管理架构，可先阅读[概览与常见问题](/security/vuln)；实际使用时，可以从以下工具入手。

### 使用 govulncheck 扫描代码中的漏洞 {#scan-code-for-vulnerabilities-with-govulncheck}

开发者可以使用 govulncheck 判断已知漏洞是否影响自己的代码，并根据程序实际调用了哪些存在漏洞的函数和方法，确定后续处理的优先级。

- [查看 govulncheck 文档](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [教程：govulncheck 入门](/doc/tutorial/govulncheck)

### 在编辑器中检测漏洞 {#detect-vulnerabilities-from-your-editor}

VS Code Go 扩展可以检查第三方依赖，并提示相关漏洞。

- [用户文档](/security/vuln/editor)
- [下载 VS Code Go](https://marketplace.visualstudio.com/items?itemName=golang.go)
- [教程：VS Code Go 入门](/doc/tutorial/govulncheck-ide)

### 查找可供使用的 Go 模块 {#find-go-modules-to-build-upon}

[Pkg.go.dev](https://pkg.go.dev/) 用于发现、评估和了解 Go 包与模块。查找和评估包时，如果某个版本存在漏洞，页面顶部会[显示提示横幅](https://pkg.go.dev/golang.org/x/text@v0.3.7/language)。此外，版本历史页面还会列出[影响包各个版本的漏洞](https://pkg.go.dev/golang.org/x/text@v0.3.7/language?tab=versions)。

### 浏览漏洞数据库 {#browse-the-vulnerability-database}

Go 漏洞数据库直接从 Go 包维护者处收集数据，也汇集 [MITRE](https://www.cve.org/) 和 [GitHub](https://github.com/) 等外部来源的信息。报告由 Go 安全团队整理维护。

- [浏览 Go 漏洞数据库中的报告](https://pkg.go.dev/vuln/)
- [查看 Go 漏洞数据库文档](/security/vuln/database)
- [向数据库提交已公开的漏洞](/s/vulndb-report-new)

## 报告 Go 项目的安全缺陷 {#report-security-bugs-in-the-go-project}

### [安全政策](/security/policy) {#security-policysecuritypolicy}

关于如何[报告 Go 项目的漏洞](/security/policy#reporting-a-security-bug)，请阅读安全政策。该页面也详细介绍 Go 安全团队跟踪问题、向公众披露问题的流程。过往安全修复见[发布历史](/doc/devel/release)。根据[发布政策](/doc/devel/release#policy)，Go 安全团队为最近两个主要发行版本提供安全修复。

- [常见漏洞报告的分类判定](/doc/security/decisions)

## 使用模糊测试检查意外输入 {#test-unexpected-inputs-with-fuzzing}

Go 原生模糊测试是一种自动化测试方法，通过持续变异程序输入来发现缺陷。从 Go 1.18 起，标准工具链支持模糊测试。[OSS-Fuzz 也支持 Go 原生模糊测试](https://google.github.io/oss-fuzz/getting-started/new-project-guide/go-lang/#native-go-fuzzing-support)。

- [了解模糊测试基础](/security/fuzz)
- [教程：模糊测试入门](/doc/tutorial/fuzz)

## 使用 Go 密码学库保护服务 {#secure-services-with-gos-cryptography-libraries}

Go 密码学库旨在帮助开发者构建安全的应用。参见 [crypto 相关包](https://pkg.go.dev/golang.org/x/crypto)及 [golang.org/x/crypto/](https://pkg.go.dev/golang.org/x/crypto) 的文档。

## 符合 FIPS 140-3 的密码学实现 {#fips-140-3-compliant-cryptography}

Go 密码学库可以在符合 FIPS 140-3 的模式下运行，以适用于有相关要求的环境。详情参见 [FIPS 140-3 合规说明](/doc/security/fips140)。
