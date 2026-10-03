---
title: Go 威胁模型
layout: article
breadcrumb: true
---

## 概述 {#overview}

本文定义 Go 工具链和标准库的通用威胁模型。如果某个包的文档没有单独定义威胁模型，应认为本文的模型适用于该标准库包。

### 威胁模型 {#threat-model}

#### 构建安全 {#build-safety}

构建 Go 代码应当是安全的，不应产生意外执行代码等副作用。

#### 运行时执行 {#runtime-execution}

通常，执行恶意代码本身不被视为属于本模型的安全问题。这里假定熟悉 Go 的用户了解自己正在执行什么代码。

#### 内存安全 {#memory-safety}

在没有使用 `unsafe` 包的情况下，假定运行时保证内存安全。

#### 信任边界 {#trust-boundaries}

对于合理预期会接收任意用户输入的 API，应采取防护，避免 panic 和不受约束的资源消耗。

向 API 传入无意义的数据并得到意外输出，本身不被视为安全问题。

#### 环境控制 {#environment-control}

假定本地系统是安全的。依赖于操作系统已经被攻破的攻击，不属于此模型的考虑范围。例如，攻击者控制文件系统、PATH 等环境变量，或已经能够访问、控制内存，都不在此模型之内。

### 有独立威胁模型的包 {#packages-with-their-own-models}

* [encoding/json](/pkg/encoding/json/#hdr-Security_Considerations)
* [encoding/json/v2](/pkg/encoding/json/v2/#hdr-Security_Considerations)
* [encoding/json/jsontext](/pkg/encoding/json/jsontext/#hdr-Security_Considerations)
* [encoding/gob](/pkg/encoding/gob/#hdr-Security)
* [html/template](/pkg/html/template/#hdr-Security_Model)
* [image](/pkg/image/#hdr-Security_Considerations)
* [debug/pe](/pkg/debug/pe/#hdr-Security)
* [debug/macho](/pkg/debug/macho/#hdr-Security)
* [debug/dwarf](/pkg/debug/dwarf/#hdr-Security)
* [debug/plan9obj](/pkg/debug/plan9obj/#hdr-Security)
* [debug/elf](/pkg/debug/elf/#hdr-Security)
