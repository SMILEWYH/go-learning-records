---
title: "开发运维与站点可靠性工程"
linkTitle: "开发运维与站点可靠性工程"
description: "快速的构建、精简的语法、自动格式化工具和文档生成器，使 Go 能够同时支持 DevOps 和 SRE 工作。"
date: 2019-10-03T17:16:43-04:00
series: Use Cases
template: true
books:
icon:
  file: devops-green.svg
  alt: "运维图标"
iconDark:
  file: devops-white.svg
  alt: "运维图标"
---

## 概述 {#overview .sectionHeading}

### Go 帮助企业实现自动化和规模扩展 {#go-helps-enterprises-automate-and-scale}

开发运维（DevOps）团队帮助工程组织自动化各项任务，并改进持续集成、持续交付和持续部署（CI/CD）流程。DevOps 能够打破开发环节之间的隔阂，通过工具和自动化改善软件开发、部署和支持工作。

独立 DevOps 顾问 [Silvia Fressard 写道](https://opensource.com/article/18/10/what-site-reliability-engineer)，站点可靠性工程（SRE）诞生于 Google，旨在使公司的“大型站点更加可靠、高效且可扩展”。“他们形成的实践很好地满足了 Google 的需求，因此 Amazon、Netflix 等其他大型科技公司也采用了这些方法。”SRE 要求兼具开发和运维技能，并“[让软件开发者](https://stackify.com/site-reliability-engineering/)负责其应用在生产环境中的持续日常运行”。

从快速构建、精简语法，到安全性和可靠性支持，Go 同时服务于 DevOps 和 SRE 这两个紧密相关的领域。并发和网络功能也使 Go 非常适合编写云部署管理工具：既便于实现自动化，也能随着开发基础设施扩大，兼顾速度和代码可维护性。

DevOps/SRE 团队编写的软件，既包括小脚本、命令行界面（CLI），也包括复杂的自动化系统和服务。Go 的特性在这些场景中都能带来帮助。

## 主要优势 {#key-benefits .sectionHeading}

### 借助完善的标准库和静态类型，轻松编写小脚本 {#easily-build-small-scripts-with-gos-robust-standard-library-and-static-typing}

Go 构建和启动都很快。丰富的标准库涵盖 HTTP、文件 I/O、时间、正则表达式、进程执行（exec）以及 JSON/CSV 格式等常见需求，让 DevOps/SRE 可以直接投入业务逻辑。静态类型系统和显式错误处理，也让小脚本更加健壮。

### 借助快速构建，迅速部署命令行工具 {#quickly-deploy-clis-with-gos-fast-build-times}

每位站点可靠性工程师都写过“一次性”脚本，后来却变成数十名工程师每天使用的 CLI；小型部署自动化脚本也会逐渐演变成发布管理服务。软件范围不可避免地扩大时，Go 能让 DevOps/SRE 团队从容应对。从一开始就使用 Go，有助于为这种发展做好准备。

### 借助较低的内存占用和文档生成器，扩展并维护大型应用 {#scale-and-maintain-larger-applications-with-gos-low-memory-footprint-and-doc-generator}

Go 的垃圾回收器让 DevOps/SRE 团队无需担心内存管理。自动文档生成器 godoc 则让代码自带文档，降低维护成本，并从项目开始就建立最佳实践。

{{projects `
  - company: Docker
    url: https://docker.com/
    logoSrc: docker.svg
    logoSrcDark: docker.svg
    desc: "Docker 是用 Go 编写的软件即服务（SaaS）产品，DevOps/SRE 团队使用它“在大规模环境中实现安全的自动化和部署”，为 CI/CD 工作提供支持。"
    ctas:
      - text: "Docker 的 CI/CD"
        url: https://www.docker.com/solutions/cicd
  - company: Drone
    url: https://github.com/drone
    logoSrc: drone.svg
    logoSrcDark: drone.svg
    desc: "Drone 是基于容器技术、用 Go 编写的持续交付系统。它使用一个简单的 YAML 配置文件（docker-compose 的超集），定义并在 Docker 容器中执行流水线。"
    ctas:
      - text: Drone
        url: https://github.com/drone
  - company: etcd
    url: https://github.com/etcd-io/etcd
    logoSrc: etcd.svg
    logoSrcDark: etcd.svg
    desc: "etcd 是用 Go 编写的强一致性分布式键值存储，为分布式系统或机器集群需要访问的数据提供可靠的存储方式。"
    ctas:
      - text: etcd
        url: https://github.com/etcd-io/etcd
  - company: IBM
    url: https://ibm.com/
    logoSrc: ibm.svg
    logoSrcDark: ibm.svg
    desc: "IBM 的 DevOps 团队通过 Docker、Kubernetes 及其他用 Go 编写的 DevOps 和 CI/CD 工具使用 Go。公司还提供专用的 Go API，用于连接其消息中间件。"
    ctas:
      - text: "用 Go 编写 IBM 应用"
        url: https://developer.ibm.com/messaging/2019/02/05/simplified-ibm-mq-applications-golang/
  - company: Netflix
    url: https://netflix.com/
    logoSrc: netflix.svg
    logoSrcDark: netflix.svg
    desc: "Netflix 使用名为 Rend 的 Go 服务处理大规模数据缓存，管理个性化数据在全球复制的存储。"
    ctas:
      - text: "应用数据缓存"
        url: https://medium.com/netflix-techblog/application-data-caching-using-ssds-5bf25df851ef
      - text: Rend
        url: https://github.com/netflix/rend
  - company: Microsoft
    url: https://microsoft.com/
    logoSrc: microsoft_light.svg
    logoSrcDark: microsoft_dark.svg
    desc: "Microsoft 在 Azure Red Hat OpenShift 服务中使用 Go。该方案向 DevOps 团队提供 OpenShift 集群，让他们在满足监管合规要求的同时专注于应用开发。"
    ctas:
      - text: OpenShift
        url: https://azure.microsoft.com/en-us/services/openshift/
  - company: Terraform
    url: https://terraform.io/
    logoSrc: terraform-icon.svg
    logoSrcDark: terraform-icon.svg
    desc: "Terraform 是用 Go 编写的工具，可安全、高效地构建、变更基础设施并管理其版本。它支持 AWS、IBM Cloud、GCP 和 Microsoft Azure 等多个云服务商。"
    ctas:
      - text: Terraform
        url: https://www.terraform.io/intro/index.html
  - company: Prometheus
    url: https://github.com/prometheus/prometheus
    logoSrc: prometheus.svg
    logoSrcDark: prometheus.svg
    desc: "Prometheus 是最初由 SoundCloud 构建的开源系统监控与告警工具包。其大多数组件使用 Go 编写，因此很容易构建并以静态二进制文件部署。"
    ctas:
      - text: Prometheus
        url: https://github.com/prometheus/prometheus
  - company: YouTube
    url: https://youtube.com/
    logoSrc: youtube.svg
    logoSrcDark: youtube.svg
    desc: "YouTube 使用 Go 编写的 Vitess（现属 PlanetScale）数据库集群系统，通过通用分片机制实现 MySQL 的水平扩展。自 2011 年起，它就是 YouTube 数据库基础设施的核心组件，如今已扩展到数万个 MySQL 节点。"
    ctas:
      - text: Vitess
        url: https://github.com/vitessio/vitess
`}}

## 开始使用 {#get-started .sectionHeading}

### DevOps 与 SRE 领域的 Go 书籍 {#go-books-on-devops--sre}

{{books `
  - title: "Go 网络运维编程（Go Programming for Network Operations）"
    url: https://www.amazon.com/Go-Programming-Network-Operations-Automation-ebook/dp/B07JKKN34L/ref=sr_1_16
    thumbnail: /images/books/go-programming-for-network-operations.jpg
  - title: "Go 编程蓝图（Go Programming Blueprints）"
    url: https://github.com/matryer/goblueprints
    thumbnail: /images/learn/go-programming-blueprints.png
  - title: "Go 语言实战（Go in Action）"
    url: https://www.amazon.com/Go-Action-William-Kennedy/dp/1617291781
    thumbnail: /images/books/go-in-action.jpg
  - title: "Go 程序设计语言（The Go Programming Language）"
    url: https://www.gopl.io/
    thumbnail: /images/learn/go-programming-language-book.png
`}}

{{libraries `
  - title: "监控与追踪"
    viewMoreUrl: https://pkg.go.dev/search?q=tracing
    items:
      - text: open-telemetry/opentelemetry-go
        url: https://pkg.go.dev/go.opentelemetry.io/otel
        desc: "供应商中立的 API 和插桩工具，用于监控及分布式追踪"
      - text: jaegertracing/jaeger-client-go
        url: https://pkg.go.dev/github.com/jaegertracing/jaeger-client-go?tab=overview
        desc: "Uber 开发的开源分布式追踪系统"
      - text: grafana/grafana
        url: https://pkg.go.dev/github.com/grafana/grafana?tab=overview
        desc: "用于监控与可观测性的开源平台"
      - text: istio/istio
        url: https://pkg.go.dev/github.com/istio/istio?tab=overview
        desc: "开源服务网格与可集成平台"
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
  - title: "其他项目"
    items:
      - text: golang-migrate/migrate
        url: https://pkg.go.dev/github.com/golang-migrate/migrate?tab=overview
        desc: "用 Go 编写的数据库迁移工具"
`}}
