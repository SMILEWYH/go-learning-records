---
title: PayPal 借助 Go 实现现代化与规模扩展
date: 2020-06-01
company: PayPal
logoSrc: paypal.svg
logoSrcDark: paypal.svg
heroImgSrc: go_paypal_case_study_logo.png
carouselImgSrc: go_paypal_case_study.png
series: Case Studies
quote: Go 能够生成简洁、高效的代码，并随着软件部署规模扩大而轻松扩展，因此非常适合帮助 PayPal 实现目标。
template: true
---

{{pullquote `
  author: Bala Natarajan
  title: <span class="NoWrapSpan">工程高级总监，</span>&nbsp;<span class="NoWrapSpan">开发者体验</span>
  company: PayPal
  quote: |
    我们的 NoSQL 和数据库代理在多线程模式下涉及大量系统细节，处理各种情况的代码变得很复杂。Go 提供了通道和 goroutine 来应对复杂性，使我们能够按照需求组织代码。
`}}

## 基于 Go 构建新的代码基础设施 {#new-code-infrastructure-built-on-go}

PayPal 的创立旨在让金融服务更加普惠，帮助个人和企业参与全球经济并蓬勃发展。PayPal 支付平台是这一目标的核心，它结合自研和第三方技术，高效、安全地促成全球数百万商户与消费者之间的交易。随着支付平台规模扩大、复杂度增加，PayPal 希望推进系统现代化，并缩短新应用的上线时间。

Go 能够生成简洁、高效的代码，并随着软件部署规模扩大而轻松扩展，因此非常适合帮助 PayPal 实现目标。

支付处理平台的核心，是 PayPal 用 C++ 开发的一款自研 NoSQL 数据库。然而，代码的复杂性严重影响了开发者演进平台的能力。Go 简洁的代码组织方式、goroutine（轻量级执行线程）和通道（连接并发 goroutine 的通信管道），使其成为 NoSQL 开发团队简化平台、推进现代化的自然选择。

为了验证可行性，一个开发团队花了六个月学习 Go，并从头用 Go 重新实现了 NoSQL 系统；在此过程中，他们还为 PayPal 更广泛地采用 Go 提供了经验。本文撰写时，已有 30% 的集群迁移到了新的 NoSQL 数据库。

## 使用 Go 简化大规模系统 {#using-go-to-simplify-for-scale}

随着 PayPal 平台日益复杂，Go 提供了一种便捷方式，降低大规模软件开发和运行的复杂度。它为 PayPal 带来出色的库、快速的工具，以及并发、垃圾回收和类型安全能力。

Go 让 PayPal 开发者摆脱 C++ 和 Java 开发过程中的繁杂干扰，能够花更多时间阅读代码并进行全局思考。

新的 NoSQL 系统重写成功后，PayPal 内部更多的平台团队和内容团队开始采用 Go。Natarajan 目前负责的团队管理 PayPal 的构建、测试和发布流水线，这些系统全部使用 Go 构建。公司拥有大型构建和测试集群，完全由 Go 基础设施管理，向全公司的开发者提供构建即服务和测试即服务。

  <img
    loading="lazy"
    width="607"
    height="289"
    class=""
    alt="Go Gopher 工厂"
    src="/images/gophers/factory.png">

## 使用 Go 推进 PayPal 系统现代化 {#modernizing-paypal-systems-with-go}

针对 PayPal 所需的分布式计算能力，Go 是更新系统的合适语言。PayPal 需要支持并发与并行、通过编译获得高性能、具有高度可移植性，并让开发者受益于模块化、可组合的开源架构。Go 提供了所有这些能力，乃至更多，帮助 PayPal 实现系统现代化。

安全性和可维护支持能力是 PayPal 的关键关注点。Go 的简洁性和模块化有助于实现这些目标，因此在公司的运维流水线中占据越来越重要的位置。部署 Go 为开发者提供了发挥创造力的平台，使他们能够面向 PayPal 的全球市场，大规模构建简单、高效、可靠的软件。

PayPal 持续使用 Go 更新软件定义网络（SDN）基础设施，在获得更易维护的代码之外，也看到了性能收益。例如，Go 如今已用于路由器、负载均衡器和越来越多的生产系统。

{{backgroundquote `
  author: Bala Natarajan
  title: 工程高级总监
  quote: |
    在我们运行 Go 代码的严格受控环境中，CPU 使用量下降了约 10%，同时代码更简洁、更易维护。
`}}

## Go 提高开发者生产力 {#go-increases-developer-productivity}

作为一家全球运营的企业，PayPal 需要开发团队有效应对两种规模问题：一是生产规模，尤其是与许多其他服务器交互的并发系统，例如云服务；二是开发规模，尤其是由许多程序员协作开发的大型代码库，例如开源开发。

PayPal 使用 Go 应对这些问题。Go 将解释型、动态类型语言的编程便捷性，与静态类型编译语言的效率和安全性结合起来，让公司开发者从中受益。在 PayPal 推进系统现代化的过程中，对网络计算和多核计算的支持至关重要。Go 不仅提供这些支持，而且交付速度很快：在单台计算机上编译一个大型可执行程序，最多只需要几秒钟。

PayPal 目前有超过 100 名 Go 开发者。由于公司已有许多成功运行于生产环境的 Go 实践，未来选择 Go 的开发者会更容易获得语言使用批准。

最重要的是，PayPal 开发者通过 Go 提高了生产力。Go 的并发机制让他们能够轻松编写程序，充分利用多核和联网计算机。Go 可以快速编译为机器码，应用还能享受垃圾回收的便利和运行时反射的能力。

## 缩短 PayPal 的产品上线时间 {#speeding-paypals-time-to-market}

目前，Java 和 Node 是 PayPal 优先支持的语言及运行环境，而 Go 主要用于基础设施。尽管 Go 可能永远不会在某些应用中取代 Node.js，Natarajan 仍在推动将 Go 提升为 PayPal 优先支持的语言。

在他的推动下，PayPal 也在评估迁移到 Google Kubernetes Engine（GKE），以缩短新产品的上线时间。GKE 是面向容器化应用部署的托管环境，具备生产可用性，并带来 Google 在开发者生产力、自动化运维和开源灵活性方面的最新创新。

对于 PayPal，将应用部署到 GKE 后，部署、更新和管理应用与服务会更容易，从而支持快速开发和迭代。此外，也更便于运行机器学习、通用 GPU 计算、高性能计算，以及其他受益于 GKE 所支持的专用硬件加速器的工作负载。

对 PayPal 最重要的是，结合 Go 开发和 GKE，可以轻松扩展以满足需求。Kubernetes 的自动扩缩容能够应对用户对服务需求的增长，在最关键的时候保持服务可用，并在业务低谷期缩容以节省成本。

## 让你的企业开始使用 Go {#getting-your-enterprise-started-with-go}

PayPal 的故事并非个例；其他数十家大型企业也正在发现，Go 如何帮助他们更快交付可靠的软件。全球有超过一百万名开发者使用 Go，遍布银行与商业、游戏与媒体、科技等行业，所在企业包括 [American Express](/solutions/americanexpress)、[Mercado Libre](/solutions/mercadolibre)、Capital One、Dropbox、IBM、Monzo、New York Times、Salesforce、Square、Target、Twitch、Uber，当然还有 Google。

想了解 Go 如何像帮助 PayPal 一样，帮助你的企业构建可靠、可扩展的软件，请访问 [go.dev](/)。
