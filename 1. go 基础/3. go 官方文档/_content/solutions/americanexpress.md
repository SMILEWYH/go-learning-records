---
title: "American Express 使用 Go 构建支付与积分奖励系统"
company: American Express
logoSrc: american-express.svg
logoSrcDark: american-express.svg
heroImgSrc: go_amex_case_study_logo.png
carouselImgSrc: go_amex_case_study.png
date: 2019-12-19
series: Case Studies
template: true
quote: Go 为 American Express 的支付网络和积分奖励网络提供了所需的速度与可扩展性。
---

{{pullquote `
  author: Glen Balliet
  title: 客户忠诚度平台工程总监
  company: American Express
  quote: |
    Go 与其他编程语言的不同之处在于认知负担。你可以用更少的代码完成更多工作，因此更容易推理和理解最终写出的代码。

    大多数 Go 代码看起来都很相似，所以即使面对完全陌生的代码库，也能很快上手。
`}}

## Go 改善微服务并提高生产力 {#go-improves-microservices-and-speeds-productivity}

American Express（美国运通）成立于 1850 年，是一家全球综合支付公司，提供签账卡和信用卡产品、商户收单及处理服务、网络服务以及旅行相关服务。

American Express 的支付处理系统在公司漫长的发展历程中逐步建立，并经历了多轮架构演进。每一次升级都必须首先保证支付处理速度，尤其是在交易量极大的情况下；同时，系统之间需要具备韧性，并且全部符合安全和监管标准。Go 为 American Express 的支付网络和积分奖励网络提供了所需的速度与可扩展性。

### 推进 American Express 系统现代化 {#modernizing-american-express-systems}

American Express 深知编程语言领域正在发生重大变化。公司现有系统专为高并发、低延迟设计，但也即将在不久后迁移到新平台。因此，支付平台团队决定花时间研究，哪些语言最适合公司不断变化的需求。

American Express 的支付平台和积分奖励平台团队，是最早开始评估 Go 的团队之一。他们关注微服务、交易路由和负载均衡等场景，需要实现架构现代化。许多开发者熟悉 Go 的能力，希望在高并发、低延迟应用中试用 Go，例如定制的交易负载均衡器。为此，团队开始争取高级管理层的支持，希望将 Go 部署到支付平台上。

American Express 副总裁兼首席工程师 Benjamin Cane 说：“我们希望找到最适合编写快速、高效的支付处理应用的语言。为此，我们在内部开展了一场编程语言比拼，看看哪种语言最符合我们的设计和性能需求。”

### 比较不同语言 {#comparing-languages}

在评估中，Cane 的团队选择用四种不同的编程语言构建同一个微服务，然后比较它们的速度与性能、工具、测试和开发便捷性。

他们选择的服务是 ISO8583 到 JSON 的转换器。ISO8583 是金融交易的国际标准，在 American Express 的支付网络中广泛使用。参评的语言和运行环境是 C++、Go、Java 和 Node.js。除了 Go，其余几种当时都已在公司内部使用。

在速度方面，Go 以每秒 14 万个请求的成绩排名第二，展现出它在后端微服务场景中的优势。

尽管 Go 不是测试中最快的语言，但强大的工具提升了它的综合表现。Go 内置的测试框架、性能剖析能力和基准测试工具给团队留下了深刻印象。Cane 说：“在 Go 中编写有效的测试很容易。基准测试和性能剖析功能，让应用调优变得简单。再加上快速的构建过程，Go 让我们很容易写出经过充分测试和优化的代码。”

最终，团队选择 Go 作为构建高性能微服务的首选语言。工具、测试框架、性能和语言的简洁性，都是促成这一选择的关键因素。

### Go 与基础设施 {#go-for-infrastructure}

Cane 说：“我们的许多服务运行在 Docker 容器中，部署于基于 Kubernetes 的内部云平台。”Kubernetes 是用 Go 编写的开源容器编排系统，它通过主机集群运行容器工作负载，最典型的就是 Docker 容器。Docker 也是用 Go 编写的软件产品，通过操作系统级虚拟化，提供称为容器的可移植软件运行环境。

American Express 还使用 Prometheus 收集应用指标。Prometheus 是用 Go 编写的开源监控与告警工具包，可以收集并聚合实时事件和指标，用于监控和告警。

Kubernetes、Docker 和 Prometheus 这三项 Go 解决方案，共同帮助 American Express 推进基础设施现代化。

### 用 Go 改善性能 {#improving-performance-with-go}

如今，American Express 已有数十名开发者使用 Go 编程，大多数人从事对高可用性和性能有要求的平台开发。

Cane 说：“工具一直是我们旧代码库亟待改善的关键领域。我们发现 Go 拥有出色的工具，以及内置的测试、基准测试和性能剖析框架，因此很容易编写高效、具有韧性的应用。”

{{backgroundquote `
  author: Benjamin Cane
  title: 副总裁兼首席工程师
  company: American Express
  quote: |
    使用 Go 之后，我们的大多数开发者都不愿再回到其他语言。
`}}

American Express 才刚开始看到 Go 的好处。例如，Go 从设计之初就考虑了并发，使用轻量级 goroutine，而不是更重量级的操作系统线程，因此在同一地址空间内创建数十万个 goroutine 是切实可行的。借助 goroutine，American Express 的实时交易处理性能得到了提升。

与其他语言相比，Go 的垃圾回收机制在性能和开发便捷性方面也有明显优势。Cane 说：“我们观察到 Go 的垃圾回收效果远好于其他语言，而垃圾回收对实时交易处理非常重要。其他语言的垃圾回收调优可能非常复杂，使用 Go 则无需做任何调优。”

如需了解更多，请阅读[《American Express 为何选择 Go》](https://americanexpress.io/choosing-go/)，其中更详细地介绍了公司的 Go 采用历程。

### 让你的企业开始使用 Go {#getting-your-enterprise-started-with-go}

正如 American Express 使用 Go 推进支付和积分奖励网络现代化一样，其他数十家大型企业也在采用 Go。

全球有超过一百万名开发者使用 Go，遍布银行与商业、游戏与媒体、科技等行业，所在企业包括 [PayPal](/solutions/paypal)、[Mercado Libre](/solutions/mercadolibre)、Capital One、Dropbox、IBM、Mercado Libre、Monzo、New York Times、Salesforce、Square、Target、Twitch、Uber，当然还有 Google。

想了解 Go 如何像帮助 American Express 一样，帮助你的企业构建可靠、可扩展的软件，请访问 [go.dev](/)。
