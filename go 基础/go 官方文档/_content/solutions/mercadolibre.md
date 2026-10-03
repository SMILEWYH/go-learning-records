---
title: "MercadoLibre 与 Go 共同成长"
company: MercadoLibre
logoSrc: mercadolibre_light.svg
logoSrcDark: mercadolibre_dark.svg
heroImgSrc: go_mercadolibre_case_study_logo.png
carouselImgSrc: go_mercadolibre_case_study.png
date: 2019-11-10T16:26:31-04:00
series: Case Studies
quote: Go 提供简洁、高效的代码，可以随着 MercadoLibre 电商业务的增长轻松扩展；工程师用更少的代码服务不断增加的用户，开发者生产力也随之提升。
template: true
---

{{pullquote `
  author: Eric Kohan
  title: 软件工程经理
  company: MercadoLibre
  quote: |
    我认为，**Go 语言之旅是迄今为止我见过最好的编程语言入门教程**。它非常简单，又能让你对这门语言大约 80% 的内容有一个相当全面的了解。想让开发者学习 Go 并快速将它用于生产环境时，我们会建议他们从 Go 语言之旅开始。
`}}

## Go 帮助一体化生态系统吸引开发者、扩展电商规模 {#go-helps-integrated-ecosystem-attract-developers-and-scale-ecommerce}

MercadoLibre, Inc. 运营着拉丁美洲最大的在线商务生态系统，业务遍及 18 个国家。公司成立于 1999 年，总部位于阿根廷，选择使用 Go 帮助其扩展生态系统并实现现代化。Go 提供简洁、高效的代码，可以随着电商业务的增长轻松扩展；工程师用更少的代码服务不断增加的用户，开发者生产力也随之提升。

### MercadoLibre 借助 Go 扩展规模 {#mercadolibre-taps-go-for-scale}

早在 2015 年，MercadoLibre 内部就越来越明显地感到，基于 Groovy 和 Grails 的现有 API 框架已经接近极限，需要换一个平台才能继续扩展。公司的平台当时呈指数级增长，而且至今仍在增长，这给开发者带来了大量额外工作：Groovy 和 Grails 都要求开发者作出很多决策，而 Groovy 又是一门动态编程语言。这种组合不适合快速扩张，因为在这样资源消耗很大的环境中，公司需要非常有经验的开发者进行开发和调优，才能达到期望的性能。测试执行缓慢，构建和部署也很慢。因此，代码效率和可扩展性变得与开发速度同样重要。

### Go 提高系统效率 {#go-improves-system-efficiency}

核心 API 团队的工作体现了 Go 对网络效率的贡献。该团队构建并维护公司微服务体系中最核心、规模最大的 API。他们创建的用户 API 被 MercadoLibre Marketplace、MercadoPago 金融科技平台、公司的配送和物流解决方案，以及其他托管解决方案使用。这些系统要求很高的服务水平：用户 API 平均每分钟接收 800 万至 1000 万个请求。团队使用 Go，以低于 10 毫秒的单次请求耗时提供服务。

API 团队还部署 Docker 容器。原文将 Docker 描述为同样用 Go 编写的软件即服务（SaaS）产品；团队利用它虚拟化开发环境，并通过 Docker Engine 便捷地部署微服务。这套系统支持规模更大、对业务至关重要的 API，**用 Go 每分钟处理超过 2000 万个请求。**

其中一个 API 充分利用了 Go 的并发原语，高效地对来自多个服务的 ID 进行多路复用。团队只用了几行 Go 代码就完成了这项工作。该 API 的成功，让核心 API 团队决定将越来越多的微服务迁移到 Go。最终，MercadoLibre 的成本效益和系统响应时间都得到了改善。

### Go 与可扩展性 {#go-for-scalability}

过去，公司技术栈的大部分建立在 Grails 和 Groovy 之上，后端使用关系型数据库。然而，这个具有多层结构的大型框架很快遇到了扩展性问题。

将旧架构迁移到 Go，并采用一个非常轻薄的新 API 框架，简化了中间层，带来了显著的性能收益。例如，一个大型 Go 服务如今能够在**每台机器仅使用 20 MB 内存的情况下处理 70000 个请求。**

{{backgroundquote `
  author: Eric Kohan
  title: 软件工程经理
  company: MercadoLibre
  quote: |
    Go 对我们来说非常出色。它功能强大，又很容易学习；在后端基础设施的可扩展性方面，给我们带来了很大帮助。
`}}

**Go 让 MercadoLibre 将该服务使用的服务器数量减少到原来的八分之一**，从 32 台降到了 4 台，而且每台服务器需要的计算资源也更少：原来使用 4 个 CPU 核心，现在只需 2 个。使用 Go 后，公司**省去了 88% 的服务器，并将剩余服务器的 CPU 核心数减半**，大幅节省了成本。

在开发者与云服务提供商之间，MercadoLibre 使用一个名为 Fury 的平台。它是一个平台即服务工具，可以以不依赖特定云厂商的方式构建、部署、监控和管理服务。因此，任何想用 Go 创建新服务的团队，都能使用经过验证的多种服务模板，快速在 GitHub 上建立仓库，其中包含初始代码、服务的 Docker 镜像和部署流水线。这套系统让工程师专注于构建创新服务，避免新项目初始化过程中的繁琐步骤，同时有效统一了构建和部署流水线。

如今，**MercadoLibre 大约一半的流量由 Go 应用处理。**

### MercadoLibre 面向开发者使用 Go {#mercadolibre-uses-go-for-developers}

目前，MercadoLibre 基础设施的通用编程语言是 Go 和 Java。每个应用、程序和微服务都有各自的 GitHub 仓库，公司还维护了额外的工具包仓库，用来解决新问题，并让客户端与服务交互。

这些内容丰富、精心维护的 Go 和 Java 工具包，为程序员快速开发新应用提供了有力支持。此外，在拥有超过 2800 名开发者的社区中，公司设有多个内部小组，供不同开发中心和不同国家的同事交流，并就 Go 的部署提供指导。公司还组织内部工作组，为新加入的 Go 开发者提供培训，并为外部开发者举办 Go 聚会，帮助建立更广泛的拉丁美洲 Go 开发者社区。

### Go 成为招聘优势 {#go-as-a-recruiting-tool}

MercadoLibre 对 Go 的推广也成为公司强有力的招聘工具。它是阿根廷最早使用 Go 的公司之一，也可能是拉丁美洲在生产环境中如此广泛使用 Go 的最大企业。公司总部位于布宜诺斯艾利斯，周围有许多初创企业和新兴科技公司。MercadoLibre 采用 Go，影响了整个潘帕斯地区的开发者市场。

{{backgroundquote `
  author: Eric Kohan
  title: 软件工程经理
  company: MercadoLibre
  quote: |
    我们非常认同这门语言的整体理念。我们喜欢 Go 的简洁性，也发现它非常明确的错误处理方式对开发者大有裨益，因为这能带来更安全、更稳定的生产代码。
`}}

如今，布宜诺斯艾利斯的程序员招聘市场竞争非常激烈，就业选择很多。当地对技术的高需求带来了优厚的薪酬和福利，也让程序员能够更从容地挑选雇主。因此，MercadoLibre 与当地所有招聘工程师和程序员的企业一样，努力提供有吸引力的工作环境和良好的职业发展路径。Go 已成为公司的一项重要差异化优势：公司为外部开发者组织 Go 工作坊，让他们前来学习；当他们喜欢上所做的事情和交流的同伴时，就会很快发现 MercadoLibre 是一个有吸引力的工作场所。

### Go 帮助开发者成长 {#go-enabling-developers}

MercadoLibre 选择 Go，是因为它能简洁地处理大规模系统；而这种简洁性，也正是公司开发者喜欢 Go 的原因。

公司还使用 [Go by Example](https://gobyexample.com/) 和 [Effective Go](/doc/effective_go.html) 等网页培训新程序员，并分享具有代表性的内部 Go API，帮助他们加快理解和熟练掌握。开发者先获得学习语言所需的资源，再凭借自己的技能与热情开始编程。

{{backgroundquote `
  author: Federico Martin Roasio
  title: 技术项目负责人
  company: MercadoLibre
  quote: |
    Go 非常适合编写业务逻辑，而我们正是编写这些 API 的团队。
`}}

Go 富有表达力且简洁的语法，让开发者更容易编写能在现代云平台上高效运行的程序。更快的开发速度为公司提高了成本效益，而开发者个人也受益于 Go 易于快速上手的特点。不仅经验丰富的工程师能用 Go 快速构建关键应用，初级工程师也能编写公司在其他语言中只会交给资深开发者的服务。例如，一组每分钟处理近 1000 万个请求的关键用户 API，就是由初级软件工程师开发的；其中许多人对编程的了解，主要来自不久前的大学课程。同样，公司也看到熟悉 Java、.NET 或 Ruby 等其他语言的开发者，能在短短几周内学会 Go 并开始编写生产服务。

使用 Go 后，MercadoLibre 的**构建速度提升到原来的 3 倍**，**测试套件的运行速度提升到原来的 24 倍**。这意味着开发者修改代码后，可以比过去快得多地完成构建和测试。

将测试套件的运行时间从 90 秒缩短到 **Go 中的仅 3 秒**，给开发者带来了巨大帮助：测试更快完成，让他们能够保持专注，不必丢失当前工作的上下文。

基于这些成功经验，MercadoLibre 不仅持续投入程序员培训，也持续开展 Go 培训。公司每年派关键工程负责人参加 GopherCon 和其他 Go 活动；基础设施和安全团队鼓励所有开发团队及时更新 Go 版本；公司还有专门团队开发 *Go-meli-toolkit*，这是一个用于对接 Fury 所提供全部服务的完整 Go 库。

### 让你的企业开始使用 Go {#getting-your-enterprise-started-with-go}

正如 MercadoLibre 从概念验证项目开始引入 Go 一样，其他数十家大型企业也在采用 Go。

全球有超过一百万名开发者使用 Go，遍布银行与商业、游戏与媒体、科技等行业，所在企业包括 [American Express](/solutions/americanexpress)、[PayPal](/solutions/paypal)、Capital One、Dropbox、IBM、Monzo、New York Times、Salesforce、Square、Target、Twitch、Uber，当然还有 Google。

想了解 Go 如何像帮助 MercadoLibre 一样，帮助你的企业构建可靠、可扩展的软件，请访问 [go.dev](/)。
