---
title: "命令行界面（CLI）"
linkTitle: "命令行界面（CLI）"
description: "借助常用开源包和完善的标准库，使用 Go 创建快速、优雅的命令行工具。"
date: 2019-10-04T15:26:31-04:00
series: Use Cases
template: true
icon:
  file: clis-green.svg
  alt: "命令行图标"
iconDark:
  file: clis-white.svg
  alt: "命令行图标"
---

## 概述 {#overview .sectionHeading}

### CLI 开发者青睐 Go 的可移植性、性能和易用性 {#cli-developers-prefer-go-for-portability-performance-and-ease-of-creation}

命令行界面（CLI）是纯文本界面，与图形用户界面（GUI）不同。由于易于自动化并支持远程操作，云计算和基础设施应用主要采用命令行界面。

## 主要优势 {#key-benefits .sectionHeading}

### 借助快速编译，构建启动迅速、可跨系统运行的程序 {#leverage-fast-compile-times-to-build-programs-that-start-quickly-and-run-on-any-system}

CLI 开发者发现，Go 非常适合设计这类应用。它可以快速编译成单个二进制文件，在不同平台上保持一致的风格，并拥有强大的开发者社区。只需一台 Windows 或 Mac 笔记本电脑，开发者就能在几秒钟内，为 Go 支持的数十种架构和操作系统构建程序，无需复杂的构建集群。原文认为，其他编译型语言难以达到这样的可移植构建能力和速度。Go 应用可以构建为一个自包含的二进制文件，因此安装非常简单。

具体而言，**用 Go 编写的程序可以在目标系统上运行，无需预先安装其他库、运行时或依赖**。而且，**Go 程序可以立即启动**——类似于 C 或 C++，原文认为这是其他编程语言难以实现的。

## 应用场景 {#use-case .sectionHeading}

### 使用 Go 构建优雅的命令行工具 {#use-go-for-building-elegant-clis}

{{backgroundquote `
  author: Steve Domino
  title: "Strala 高级工程师兼架构师"
  link: https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056
  quote: |
    我负责构建公司的 CLI 工具时，发现了两个非常出色的项目：Cobra 和 Viper。它们让 CLI 开发变得简单。两者各自都强大、灵活，并且擅长自己的领域；结合使用时，更能让你从容掌控下一个 CLI 项目！
`}}

{{backgroundquote `
  author: Francesc Campoy
  title: "DGraph Labs 产品副总裁、Just For Func 视频制作者"
  link: https://www.youtube.com/watch?v=WvWPGVKLvR4
  quote: |
    Cobra 非常适合编写小型工具，甚至大型工具。与其说它是一个库，不如说是一个框架：调用它的二进制程序生成骨架后，你再往里面添加代码。
`}}

使用 Go 开发 CLI 时，Cobra 和 Viper 是两种广泛使用的工具。

{{pkg "github.com/spf13/cobra" "Cobra"}} 既是用于创建强大、现代 CLI 应用的库，也提供用 Go 生成应用和命令行应用的程序。许多流行的 Go 应用都使用 Cobra，包括 CoreOS、Delve、Docker、Dropbox、Git Lfs、Hugo、Kubernetes，以及[更多项目](https://pkg.go.dev/github.com/spf13/cobra?tab=importedby)。OpenFaaS 创始人 [Alex Ellis](https://blog.alexellis.io/5-keys-to-a-killer-go-cli/) 表示，集成的命令帮助、自动补全和文档功能，“让为每个命令编写文档变得非常简单”。


{{pkg "github.com/spf13/viper" "Viper"}} 是 Go 应用的完整配置解决方案，用于在应用内部处理各种配置需求和格式。Cobra 和 Viper 在设计上就支持协同使用。

Viper 的配置[支持嵌套结构](https://scene-si.org/2017/04/20/managing-configuration-with-viper/)，让 CLI 开发者能够管理大型应用各个部分的配置。Viper 还提供了轻松构建十二要素应用所需的工具。

[Geudens 建议](https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html)：“如果不想让命令行变得杂乱，或者正在处理不希望出现在历史记录中的敏感数据，使用环境变量是个好办法。你可以用 Viper 来实现。”

{{projects `
  - company: Comcast
    url: https://xfinity.com/
    logoSrc: comcast.svg
    logoSrcDark: comcast.svg
    desc: "Comcast 用 Go 开发 CLI 客户端，在高流量站点中进行发布和订阅。公司还维护一个用 Go 编写的开源客户端库，专门用于 Apache Pulsar。"
    ctas:
      - text: "Apache Pulsar 客户端库"
        url: https://github.com/Comcast/pulsar-client-go
      - text: "Pulsar 命令行客户端"
        url: https://github.com/Comcast/pulsar-client-go/blob/master/cli/main.go
  - company: GitHub
    url: https://github.com/
    logoSrc: github.svg
    logoSrcDark: github.svg
    desc: "GitHub 使用 Go 编写命令行工具，让 GitHub 操作更加便捷；它通过封装 git，为其增加额外的功能和命令。"
    ctas:
      - text: "GitHub 命令行工具"
        url: https://github.com/cli/cli
  - company: Hugo
    url: https://gohugo.io/
    logoSrc: hugo.svg
    logoSrcDark: hugo.svg
    desc: "Hugo 是最受欢迎的 Go CLI 应用之一，为数千个网站提供支持，包括原文所述的本站。Go 带来的易安装性是其流行的原因之一。Hugo 作者 Bjørn Erik Pedersen 写道：“单个二进制文件消除了安装和升级中的大部分麻烦。”"
    ctas:
      - text: "Hugo 网站"
        url: https://gohugo.io/
  - company: Kubernetes
    url: https://kubernetes.com/
    logoSrc: kubernetes.svg
    logoSrcDark: kubernetes.svg
    desc: "Kubernetes 是最流行的 Go CLI 应用之一。其创建者 Joe Beda 表示，编写 Kubernetes 时，“Go 是唯一合理的选择”，并称 Go 在 C++ 等低级语言和 Python 等高级语言之间找到了恰当的平衡。"
    ctas:
      - text: "Kubernetes 与 Go"
        url: https://blog.gopheracademy.com/birthday-bash-2014/kubernetes-go-crazy-delicious/
  - company: MongoDB
    url: https://mongodb.com/
    logoSrc: mongodb.svg
    logoSrcDark: mongodb.svg
    desc: "MongoDB 选择用 Go 实现备份命令行工具，原因包括 Go 的“类 C 语法、强大的标准库、通过 goroutine 解决并发问题的能力，以及轻松的多平台分发”。"
    ctas:
      - text: "MongoDB 备份服务"
        url: https://www.mongodb.com/blog/post/go-agent-go
  - company: Netflix
    url: https://netflix.com/
    logoSrc: netflix.svg
    logoSrcDark: netflix.svg
    desc: "Netflix 使用 Go 构建 CLI 应用 ChaosMonkey。它会随机终止生产环境中的实例，确保工程师实现的服务能够承受实例故障。"
    ctas:
      - text: "Netflix 技术博客文章"
        url: https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef
  - company: Stripe
    url: https://stripe.com/
    logoSrc: stripe.svg
    logoSrcDark: stripe.svg
    desc: "Stripe 使用 Go 编写 Stripe CLI，让开发者直接在终端中构建、测试和管理 Stripe 集成。"
    ctas:
      - text: Stripe CLI
        url: https://github.com/stripe/stripe-cli
  - company: Uber
    url: https://uber.com/
    logoSrc: uber.svg
    logoSrcDark: uber.svg
    desc: "Uber 使用 Go 编写多种 CLI 工具，包括 Jaeger 的命令行 API。Jaeger 是用于监控微服务分布式系统的分布式追踪系统。"
    ctas:
      - text: "Jaeger 的命令行 API"
        url: https://www.jaegertracing.io/docs/1.14/cli/
`}}

## 开始使用 {#get-started .sectionHeading}

### 创建 CLI 的 Go 书籍 {#go-books-for-creating-clis}

{{books `
  - title: "用 Go 构建强大的命令行应用（Powerful Command-Line Applications in Go）"
    url: https://www.amazon.com/Powerful-Command-Line-Applications-Go-Maintainable/dp/168050696X
    thumbnail: /images/books/powerful-command-line-applications-in-go.jpg
  - title: "Go 语言实战（Go in Action）"
    url: https://www.amazon.com/Go-Action-William-Kennedy/dp/1617291781
    thumbnail: /images/books/go-in-action.jpg
  - title: "Go 程序设计语言（The Go Programming Language）"
    url: https://www.gopl.io/
    thumbnail: /images/learn/go-programming-language-book.png
  - title: "Go 编程蓝图（Go Programming Blueprints）"
    url: https://github.com/matryer/goblueprints
    thumbnail: /images/learn/go-programming-blueprints.png
`}}

{{libraries `
  - title: "命令行工具库"
    viewMoreUrl: https://pkg.go.dev/search?q=command%20line%20OR%20CLI
    items:
      - text: spf13/cobra
        url: https://pkg.go.dev/github.com/spf13/cobra?tab=overview
        desc: "用于创建强大、现代 CLI 应用的库，并提供用 Go 生成应用和命令行应用的程序"
      - text: spf13/viper
        url: https://pkg.go.dev/github.com/spf13/viper?tab=overview
        desc: "Go 应用的完整配置解决方案，用于处理应用内部的配置需求和格式"
      - text: urfave/cli
        url: https://pkg.go.dev/github.com/urfave/cli?tab=overview
        desc: "用于创建和组织 Go 命令行应用的极简框架"
      - text: delve
        url: https://pkg.go.dev/github.com/go-delve/delve?tab=overview
        desc: "简单而强大的调试工具，面向习惯在编译型语言中使用源码级调试器的程序员"
      - text: chzyer/readline
        url: https://pkg.go.dev/github.com/chzyer/readline?tab=overview
        desc: "纯 Go 实现，提供 GNU Readline 的大部分功能，采用 MIT 许可证"
      - text: dixonwille/wmenu
        url: https://pkg.go.dev/github.com/dixonwille/wmenu?tab=overview
        desc: "易于使用的 CLI 菜单结构，用于提示用户作出选择"
      - text: spf13/pflag
        url: https://pkg.go.dev/github.com/spf13/pflag?tab=overview
        desc: "可直接替换 Go 的 flag 包，实现 POSIX/GNU 风格的标志参数"
      - text: golang/glog
        url: https://pkg.go.dev/github.com/golang/glog?tab=overview
        desc: "Go 的分级运行日志库"
      - text: go-prompt
        url: https://pkg.go.dev/github.com/c-bata/go-prompt?tab=overview
        desc: "用于构建强大交互式提示界面的库，让使用 Go 开发跨平台命令行工具更容易。"
`}}
