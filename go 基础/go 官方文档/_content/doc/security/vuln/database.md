---
title: Go 漏洞数据库
layout: article
---

[返回 Go 漏洞管理](/security/vuln)

## 概述 {#overview}

Go 漏洞数据库（[https://vuln.go.dev](https://vuln.go.dev)）以[开源漏洞（OSV）数据格式](https://ossf.github.io/osv-schema/)提供 Go 漏洞信息。

也可以在 [pkg.go.dev/vuln](https://pkg.go.dev/vuln) 浏览数据库中的漏洞。

**不要**依赖 x/vulndb Git 仓库的内部内容。仓库中的 YAML 文件采用内部格式维护，该格式可能在没有预告的情况下变化。

## 贡献 {#contributing}

欢迎所有 Go 包维护者[提交](/s/vulndb-report-new)自己项目中已公开漏洞的信息，并[更新](/s/vulndb-report-feedback)已有漏洞信息。

我们希望让报告流程尽可能简单，欢迎[提出建议](/s/vuln-feedback)。

请**不要**使用上述表单报告 Go 标准库或子仓库中的漏洞。Go 项目的漏洞应遵循 [go.dev/security/policy](/security/policy) 中的流程。

## API {#api}

Go 漏洞数据库的规范服务地址为 [https://vuln.go.dev](https://vuln.go.dev)。它是一个 HTTP 服务器，响应对下述端点的 GET 请求。

各端点不使用查询参数，也不要求特定请求头。因此，即使只提供固定文件系统内容的站点（包括 `file://` URL），也可以实现此 API。

每个端点返回 JSON 编码的响应：请求 `.json` 时返回未压缩内容，请求 `.json.gz` 时返回 gzip 压缩内容。

端点如下：

- `/index/db.json[.gz]`

  返回数据库元数据：

  ```json
  {
    // The latest time the database should be considered
    // to have been modified, as an RFC3339-formatted UTC
    // timestamp ending in "Z".
    "modified": string
  }
  ```

  注意，*不应*将修改时间与实际时钟时间比较，例如用来判断缓存是否失效，因为数据库修改到正式上线之间可能存在延迟。

  实际示例见 [/index/db.json](https://vuln.go.dev/index/db.json)。

- `/index/modules.json[.gz]`

  返回包含数据库中各模块元数据的列表：

  ```json
  [ {
    // The module path.
    "path": string,
    // The vulnerabilities that affect this module.
    "vulns":
      [ {
        // The vulnerability ID.
        "id": string,
        // The latest time the vulnerability should be considered
        // to have been modified, as an RFC3339-formatted UTC
        // timestamp ending in "Z".
        "modified": string,
        // (Optional) The module version (in SemVer 2.0.0 format)
        // that contains the latest fix for the vulnerability.
        // If unknown or unavailable, this should be omitted.
        "fixed": string,
      } ]
  } ]
  ```

  实际示例见 [/index/modules.json](https://vuln.go.dev/index/modules.json)。

- `/index/vulns.json[.gz]`

  返回包含数据库中各漏洞元数据的列表：

  ```json
   [ {
       // The vulnerability ID.
       "id": string,
       // The latest time the vulnerability should be considered
       // to have been modified, as an RFC3339-formatted UTC
       // timestamp ending in "Z".
       "modified": string,
       // A list of IDs of the same vulnerability in other databases.
       "aliases": [ string ]
   } ]
  ```

  实际示例见 [/index/vulns.json](https://vuln.go.dev/index/vulns.json)。

- `/ID/$id.json[.gz]`

  返回编号为 `$id` 的漏洞报告，采用 OSV 格式，详见下文[数据格式](#schema)。

  实际示例见 [/ID/GO-2022-0191.json](https://vuln.go.dev/ID/GO-2022-0191.json)。

### 批量下载 {#bulk-download}

为方便下载整个 Go 漏洞数据库，[vuln.go.dev/vulndb.zip](https://vuln.go.dev/vulndb.zip) 提供包含全部索引和 OSV 文件的 ZIP 压缩包。

### 在 `govulncheck` 中使用 {#usage-in-govulncheck}

默认情况下，`govulncheck` 使用 [vuln.go.dev](https://vuln.go.dev) 提供的官方 Go 漏洞数据库。

可以通过 `-db` 标志配置其他漏洞数据库，它接受使用 `http://`、`https://` 或 `file://` 协议的数据库 URL。

要与 `govulncheck` 正常配合，所指定的数据库必须实现上文描述的 API。`govulncheck` 从 HTTP(S) 来源读取时使用压缩的“.json.gz”端点，从文件来源读取时使用“.json”端点。

### 旧版 API {#legacy-api}

官方数据库还包含一些旧版 API 端点，我们计划近期移除对这些端点的支持。如果你仍在依赖旧版 API，并需要更多迁移时间，请[告知我们](/s/govulncheck-feedback)。

## 数据格式 {#schema}

报告采用[开源漏洞（OSV）数据格式](https://ossf.github.io/osv-schema/)。Go 漏洞数据库对以下字段赋予具体含义。

### id {#id}

id 字段是漏洞条目的唯一标识符，字符串格式为 GO-\<YEAR>-\<ENTRYID>。

### affected {#affected}

[affected](https://ossf.github.io/osv-schema/#affected-fields) 字段是一个 JSON 数组，其中的对象描述存在漏洞的模块版本。

#### affected[].package {#affectedpackage}

[affected[].package](https://ossf.github.io/osv-schema/#affectedpackage-field) 字段是标识受影响*模块*的 JSON 对象，包含两个必需字段：

- **ecosystem**：始终为“Go”。
- **name**：Go 模块路径。
  - 标准库中可导入包的名称为 _stdlib_。
  - go 命令对应的名称为 _toolchain_。

#### affected[].ecosystem_specific {#affectedecosystem-specific}

[affected[].ecosystem_specific](https://ossf.github.io/osv-schema/#affectedecosystem_specific-field) 字段是包含漏洞附加信息的 JSON 对象，供 Go 漏洞检测工具使用。

目前，该对象始终只包含一个字段 `imports`。

##### affected[].ecosystem_specific.imports {#affectedecosystem-specificimports}

`affected[].ecosystem_specific.imports` 字段是 JSON 数组，包含受漏洞影响的包和符号。数组中的对象包含以下字段：

- **path**：字符串，表示存在漏洞的包的导入路径。
- **symbols**：字符串数组，列出存在漏洞的符号（函数或方法）名称。
- **goos**：字符串数组，列出这些符号出现的目标操作系统（如果已知）。
- **goarch**：字符串数组，列出这些符号出现的架构（如果已知）。

### database_specific {#database-specific}

`database_specific` 字段包含 Go 漏洞数据库特有的自定义字段。

#### database_specific.url {#database-specificurl}

`database_specific.url` 是字符串，表示 Go 漏洞报告的完整 URL，例如“https://pkg.go.dev/vuln/GO-2023-1621”。

#### database_specific.review_status {#database-specificreview-status}

`database_specific.review_status` 是字符串，表示漏洞报告的审查状态。若该字段不存在，应视为 `REVIEWED`。可能取值为：

- `UNREVIEWED`：根据 CVE、GHSA 等其他来源自动生成的报告。其数据可能不完整，尚未经过 Go 团队核实。
- `REVIEWED`：来自 Go 团队，或根据外部来源生成的报告。Go 团队成员已审查该报告，并按需补充信息。

其他字段的说明请参阅 [OSV 规范](https://ossf.github.io/osv-schema)。

## 关于版本的说明 {#note-on-versions}

我们的工具会尝试按照标准 [Go 模块版本编号](/doc/modules/version-numbers)，将来源安全公告中的模块和版本自动映射为规范的 Go 模块与版本。`govulncheck` 等工具依靠这些标准版本判断 Go 项目是否受到依赖漏洞的影响。

在某些情况下，例如 Go 项目使用自定义版本规则时，可能无法映射到标准 Go 版本。此时，Go 漏洞数据库的报告可能保守地将所有 Go 模块版本标为受影响，以避免 `govulncheck` 等工具因无法识别版本范围而漏报漏洞（假阴性）。不过，将所有版本标为受影响，也可能导致工具把已经修复的模块版本错误地报告为存在漏洞（假阳性）。

如果认为 `govulncheck` 错报或漏报了漏洞，请针对该漏洞报告[建议修改](https://github.com/golang/vulndb/issues/new?assignees=&labels=Needs+Triage%2CSuggested+Edit&template=suggest_edit.yaml&title=x%2Fvulndb%3A+suggestion+regarding+GO-2024-2965&report=GO-XXXX-YYYY)，我们会进行审查。

## 示例 {#examples}

Go 漏洞数据库中的所有漏洞均使用上述 OSV 格式。

以下链接展示了不同类型的 Go 漏洞：

- **Go 标准库漏洞**（GO-2022-0191）：[JSON](https://vuln.go.dev/ID/GO-2022-0191.json)、[HTML](https://pkg.go.dev/vuln/GO-2022-0191)。
- **Go 工具链漏洞**（GO-2022-0189）：[JSON](https://vuln.go.dev/ID/GO-2022-0189.json)、[HTML](https://pkg.go.dev/vuln/GO-2022-0189)。
- **Go 模块中的漏洞**（GO-2020-0015）：[JSON](https://vuln.go.dev/ID/GO-2020-0015.json)、[HTML](https://pkg.go.dev/vuln/GO-2020-0015)。

## 被排除的报告 {#excluded-reports}

Go 漏洞数据库中的报告来自不同来源，由 Go 安全团队整理。对于发现的安全公告（例如 CVE 或 GHSA），我们可能出于多种原因选择不收录。这种情况下，会在 x/vulndb 仓库的 [x/vulndb/data/excluded](https://github.com/golang/vulndb/tree/master/data/excluded) 下创建一份最简报告。

报告可能因以下原因被排除：

- `NOT_GO_CODE`：漏洞不在 Go 包中，却被其他来源标注为 Go 生态安全公告。这类漏洞不会影响任何 Go 包，例如 C++ 库中的漏洞。
- `NOT_IMPORTABLE`：漏洞位于 `main` 包、仅由 `main` 包导入的 `internal/` 包，或其他无法被外部模块导入的位置。
- `EFFECTIVELY_PRIVATE`：漏洞所在的 Go 包虽然能被其他模块导入，但其设计不面向外部使用，也不太可能被定义它的模块之外的代码导入。
- `DEPENDENT_VULNERABILITY`：该漏洞已被数据库中的另一个漏洞完全涵盖。例如，包 A 存在漏洞，包 B 依赖 A，而 A 和 B 分别拥有 CVE 编号，此时可能将 B 的报告标为依赖引入的漏洞，由 A 的报告完整取代。
- `NOT_A_VULNERABILITY`：虽然已分配 CVE 编号或 GHSA，但没有已知的关联漏洞。
- `WITHDRAWN`：来源方已撤回该漏洞。

目前，[vuln.go.dev](https://vuln.go.dev) API 不提供被排除的报告。如果你有具体使用场景，需要通过 API 获取这些信息，请[告知我们](/s/govulncheck-feedback)。
