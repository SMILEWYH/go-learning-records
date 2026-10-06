---
title: Go 安全策略
layout: article
breadcrumb: true
---

## 概述 {#overview}

本文介绍 Go 安全团队处理所报告问题的流程，以及报告者可以期待的反馈。

## 报告安全缺陷 {#reporting-a-security-bug}

Go 发行版中的所有安全缺陷都应通过电子邮件报告至 [security@golang.org](mailto:security@golang.org)。邮件会发送给 Go 安全团队。

请将邮件主题写成“Vulnerability: {package name}: {one-line summary}”的形式，即包含包名和一句话摘要。除非确有必要，否则请避免添加附件，以免邮件被判定为垃圾邮件。

如有需要，请在报告中注明希望使用的署名方式。如果之后发布 CVE，将使用该署名。

请尽量简明地描述发现的问题、你认为可被利用的方式，并提供能够复现问题的小型测试用例或程序。

对于常被报告的若干类问题，我们整理了[判定说明](/doc/security/decisions)。

我们会在 7 天内确认收到邮件。如果届时尚未收到回复，请再次发送邮件到 [security@golang.org](mailto:security@golang.org) 跟进，并确保邮件中包含 **vulnerability** 一词。

如果再过 3 天仍未收到确认，邮件可能已被判为垃圾邮件。此时请[在这里提交问题](https://g.co/vulnz)，选择 _“I want to report a technical security or an abuse risk related bug in a Google product (SQLi, XSS, etc.)”_（报告 Google 产品中的技术安全或滥用风险相关缺陷），并将受影响产品填写为 _“Go”_。

从首次确认收到报告起，问题将在 90 天内得到修复或公开披露。

报告安全问题的邮件将被无限期保留，以便追踪修复状态并正确记录发现者署名。

### 由大语言模型生成的报告 {#llm-generated-reports}

请不要未经充分筛选和核实，就提交大语言模型（LLM）生成的报告。

现代大语言模型很擅长发现真实而重要的安全缺陷，但也很容易提出并不存在的问题，或确实存在却[不属于安全漏洞](/doc/security/decisions#non-vuln)的问题。模型输出往往冗长，将重要信息埋在无关文字中，并包含夸大且缺乏依据的结论。

对于模型生成的报告，价值在于筛选和核实。如果一次发送几十个可能的安全问题，依赖我们找出隐藏其中的有效发现，就等于将筛选工作交给我们，最终会妨碍提升 Go 安全性的目标。因此，我们只会为真实发现占比足够高的报告者记录发现者署名。

这里没有一成不变的硬性规则，但一般而言，如果你的大多数报告都未被认定为安全问题，且我们认为你没有提供有价值的筛选和核实工作，那么剩余报告中即使有有效发现，也不会给你署名。我们刻意不为“大多数”规定具体阈值，也不严格定义“筛选和核实”的标准。我们的意图是：善意且提供实质价值的报告者始终能获得署名，而不加甄别地大量转发模型输出的报告者不会。

## 处理通道 {#tracks}

Go 安全团队会根据问题性质，将其归入 PUBLIC、PRIVATE 或 URGENT 通道。所有安全问题都会分配 CVE 编号。

Go 安全团队不为安全问题分配传统的细粒度严重程度标签，例如 CRITICAL、HIGH、MEDIUM、LOW，因为实际严重程度高度依赖用户如何使用受影响的 API 或功能。基于同样原因，我们认为 CVSS 评分体系并不适用于 Go，因此为 Go 安全问题发布 CVE 时也不会分配 CVSS 分数。MITRE 或 NIST 等第三方可能通过 NVD 为这些漏洞分配 CVSS 分数，但我们不认可这些分数能够准确反映漏洞影响。

例如，`encoding/json` 解析器中的资源耗尽问题，其影响取决于被解析的内容。若用户解析的是本地文件系统中的可信 JSON 文件，影响可能很小；若解析的是来自 HTTP 请求体的不可信任意 JSON，影响就可能大得多。

不过，下列处理通道确实反映了安全团队对问题严重程度和影响范围的判断。例如，对大量用户造成中等至显著影响的问题属于 PRIVATE 通道；影响可忽略或较轻，或者仅影响少量用户的问题属于 PUBLIC 通道。

### PUBLIC {#public}

PUBLIC 通道的问题通常只影响少见配置、影响范围很有限，或已经广为人知。

这类问题会被标记为 [`Proposal-Security`](https://github.com/golang/go/labels/Proposal-Security)，通过 [Go 提案评审流程](https://go.googlesource.com/proposal/+/master/README.md#proposal-review)讨论，**公开修复**，并回移到下一个计划中的[维护版本](/wiki/MinorReleases)（大约每月发布）。发布公告会说明这些问题的细节，但不会提前预告。

以往 PUBLIC 问题的示例包括：

- [#44916](/issue/44916)：archive/zip：调用 Reader.Open 时可能触发 panic。
- [#44913](/issue/44913)：encoding/xml：xml.NewTokenDecoder 使用自定义 TokenReader 时可能陷入死循环。
- [#43786](/issue/43786)：crypto/elliptic：P-224 曲线运算错误。
- [#40928](/issue/40928)：net/http/cgi、net/http/fcgi：未指定 Content-Type 时存在跨站脚本（XSS）问题。
- [#40618](/issue/40618)：encoding/binary：ReadUvarint 和 ReadVarint 可能从无效输入中读取无限数量的字节。
- [#36834](/issue/36834)：crypto/x509：Windows 10 上的证书验证绕过。

### PRIVATE {#private}

PRIVATE 通道的问题违反了项目承诺的安全属性。

这类问题会在**下一个计划中的[维护版本](/wiki/MinorReleases)中修复**，在此之前保持非公开状态。

发布前 3 至 7 天，会向 golang-announce 发送预告，说明即将发布的版本包含一项或多项安全修复，问题影响标准库、工具链还是两者，以及各项修复预留的 CVE 编号。

对于[主版本候选发布版](/s/release)中存在的问题，也采用相同流程，在下一个计划中的候选发布版中修复。

以往 PRIVATE 问题的示例包括：

- [#53416](/issue/53416)：path/filepath：Glob 中的栈耗尽。
- [#53616](/issue/53616)：go/parser：所有 Parse* 函数中的栈耗尽。
- [#54658](/issue/54658)：net/http：发送 GOAWAY 之后的服务器错误处理。
- [#56284](/issue/56284)：syscall、os/exec：环境变量中的 NUL 未被清理。

### URGENT {#urgent}

URGENT 通道的问题威胁 Go 生态的完整性，或正在现实环境中遭到利用并造成严重损害。近期没有此类案例，但例如 net/http 中的远程代码执行，或 crypto/tls 中具有实际可行性的密钥恢复攻击，都属于这一类。

这类问题会秘密修复，并**立即触发专门的安全版本发布**，可能不提前预告。

## 将已有问题标记为安全相关 {#flagging-existing-issues-as-security-related}

如果认为某个[已有问题](/issue)与安全相关，请发送邮件到 [security@golang.org](mailto:security@golang.org)，包含问题编号，并简要说明为什么应按本安全策略处理。

## 披露流程 {#disclosure-process}

Go 项目采用以下披露流程：

1. 收到安全报告后，指定主要负责人，由其协调修复和发布流程。
2. 确认问题，并确定受影响的软件列表。
3. 审计代码，查找潜在的类似问题。
4. 与报告者协商后，如确定需要 CVE 编号，由主要负责人申请。
5. 为最近两个主要版本以及 head/master 修订准备修复，并将修复合并到 head/master。
6. 应用修复当天，向 [golang-announce](https://groups.google.com/group/golang-announce)、[golang-dev](https://groups.google.com/group/golang-dev) 和 [golang-nuts](https://groups.google.com/group/golang-nuts) 发送公告。

该流程可能需要一些时间，尤其是在需要与其他项目维护者协调时。我们会尽力及时处理缺陷，同时也必须遵循上述流程，以保证披露方式的一致性。

分配了 CVE 编号的安全问题，会公开列在 [CVEDetails 网站的“Golang”产品条目](https://www.cvedetails.com/vulnerability-list/vendor_id-14185/Golang.html)和[国家漏洞数据库网站](https://web.nvd.nist.gov/view/vuln/search)中。

## 接收安全更新 {#receiving-security-updates}

接收安全公告的最佳方式是订阅 [golang-announce](https://groups.google.com/forum/#!forum/golang-announce) 邮件列表。安全问题相关邮件的主题会带有 `[security]` 前缀。

## 对本策略提出建议 {#comments-on-this-policy}

如果有改进建议，请[提交问题](/issue/new)进行讨论。
