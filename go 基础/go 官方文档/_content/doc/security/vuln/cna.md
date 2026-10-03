---
title: Go CNA 策略
layout: article
---

[返回 Go 漏洞管理](/security/vuln)

## 概述 {#overview}

Go CNA 是一个 [CVE 编号分配机构](https://www.cve.org/ProgramOrganization/CNAs)，为 Go 生态中的公开漏洞分配 [CVE 编号](https://www.cve.org/ResourcesSupport/Glossary?activeTerm=glossaryCVEID)，并发布 [CVE 记录](https://www.cve.org/ResourcesSupport/Glossary?activeTerm=glossaryRecord)。它是 Google CNA 下属的 CNA。

## 范围 {#scope}

Go CNA 覆盖 Go 项目中的漏洞，包括 Go [标准库](/pkg)和[子仓库](https://pkg.go.dev/golang.org/x)，以及尚未被其他 CNA 覆盖、位于可导入 Go 模块中的公开漏洞。

该范围明确排除用 Go 编写但不可导入的应用或包中的漏洞，例如 `main` 包中的任何内容。有关被排除报告的更多信息，请参阅 [go.dev/security/vuln/database#excluded-reports](/security/vuln/database#excluded-reports)。

要报告 Go 项目中潜在的新漏洞，请参阅 [go.dev/security/policy](/security/policy)。

## 为公开漏洞申请 CVE 编号 {#requesting-a-cve-id-for-a-public-vulnerability}

**重要：**下方表单会在问题跟踪系统中创建公开问题，因此*绝不能*用于报告 Go 中尚未披露的漏洞。报告未公开问题的方式见[安全策略](/security/policy)。

要为 Go 生态中已有的公开（PUBLIC）漏洞申请 CVE 编号，请[通过此表单提交申请](/s/vulndb-report-new)。

漏洞已被公开披露，或者位于你维护的包中且你已准备好公开披露时，就视为公开漏洞。
