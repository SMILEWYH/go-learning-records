---
title: "Google Core Data Solutions 团队如何使用 Go"
company: Core Data
logoSrc: google.svg
logoSrcDark: google.svg
heroImgSrc: go_core_data_case_study.png
series: Case Studies
template: true
quote: |
  Google 是一家技术公司，其使命是整合全球信息，使人人皆可访问并从中受益。

  在本案例中，Google Core Data Solutions 团队分享了他们使用 Go 的历程，包括为何决定用 Go 重写网页索引服务、如何利用 Go 内置的并发能力，以及 Go 如何帮助他们改进开发流程。
authors:
  - Prasanna Meda，Core Data Solutions 软件工程师
---

Google 的使命是“整合全球信息，使人人皆可访问并从中受益”。Google Core Data Solutions 团队正是负责组织这些信息的团队之一。除了其他工作，该团队还维护对全球网页建立索引的服务。这些网页索引服务通过保持搜索结果的及时性和完整性，为 Google 搜索等产品提供支持，而它们使用 Go 编写。

2015 年，为了适应 Google 的规模，我们需要重写索引技术栈：将一个用 C++ 编写的单体二进制程序，改造成微服务架构下的多个组件。我们决定用 Go 重写许多索引服务，如今，架构中的大部分服务都由 Go 提供支持。

{{backgroundquote `
  author: Minjae Hwang
  title: 软件工程师
  quote: |
    Go 内置的并发能力天然适合我们，因为团队鼓励工程师使用并发和并行算法。
`}}

选择语言时，我们发现 Go 的多项特性使其尤其适合这一任务。例如，团队鼓励工程师使用并发和并行算法，因此 Go 内置的并发能力天然契合需求。工程师们还发现，“Go 代码更自然”，让他们能够将时间集中在业务逻辑和分析上，而不是内存管理和性能优化上。

使用 Go 编写代码简单得多，因为它有助于减轻开发过程中的认知负担。例如，Core Data Solutions 团队的软件工程师 MinJae Hwang 表示，使用 C++ 时，功能完善的 IDE 可能“显示源代码没有编译错误，但实际上存在错误”；而在 Go 中，“只要 IDE 说代码没有编译错误，代码就总能编译通过”。减少开发过程中的细小阻碍，例如缩短修复编译错误的周期，帮助团队在最初重写时更快交付，也让后续维护成本保持较低水平。

Hwang 还说：“使用 C++ 时，如果想使用更多的包，我必须编写头文件之类的内容。而使用 Go 时，**内置工具让我更容易使用包，开发速度快得多。**”

简洁的语言语法和 Go 工具的支持，让团队中的多位成员感到编写 Go 代码容易得多。我们还发现，Go 的静态类型检查做得很好，而 godoc 命令等 Go 的基础工具，有助于团队形成更规范的文档编写文化。

{{backgroundquote `
  author: Prasanna Meda
  title: 软件工程师
  quote: |
    ……Google 的网页索引系统在一年内完成了架构重构。更令人印象深刻的是，团队中的大多数开发者是在学习 Go 的同时用它进行重写的。
`}}

开发一款在全球被如此广泛使用的产品绝非易事，团队选择 Go 也并非轻松的决定，但这一选择帮助我们更快地前进。最终，Google 的网页索引系统在一年内完成了架构重构。更令人印象深刻的是，团队中的大多数开发者是在学习 Go 的同时用它进行重写的。

除了 Core Data Solutions 团队，Google 的许多工程团队也已在开发过程中采用 Go。你可以继续了解 [Chrome](/solutions/google/chrome/) 和 [Firebase Hosting](/solutions/google/firebase/) 团队如何使用 Go 大规模构建快速、可靠、高效的软件。
