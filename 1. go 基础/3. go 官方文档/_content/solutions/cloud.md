---
title: "Go 与云计算及网络服务"
linkTitle: "云计算与网络服务"
description: "依托主流云服务商丰富的工具和 API 生态，使用 Go 构建服务比以往更容易。"
date: 2019-10-04T15:26:31-04:00
series: Use Cases
template: true
icon:
  file: cloud-green.svg
  alt: "云图标"
iconDark:
  file: cloud-white.svg
  alt: "云图标"
---

## 概述 {#overview .sectionHeading}

<div class="UseCase-halfColumn">
    <h3>Go 帮助企业构建和扩展云计算系统</h3>
    <p>随着应用和处理任务迁移到云端，并发成为一个重大问题。云计算系统本质上需要共享和扩展资源。协调共享资源的访问会影响每个在云中处理任务的应用，因此需要“专门面向高可靠并发应用开发”的编程语言。</p>
  </div>

{{quote `
  author: Ruchi Malik
  title: "Choozle 开发者"
  link: https://builtin.com/software-engineering-perspectives/golang-advantages
  quote: |
    Go 让企业能够轻松扩展规模。这非常重要，因为随着工程团队扩大，每项服务都可以由不同的小组负责。
`}}

## 主要优势 {#key-benefits .sectionHeading}

### 缓解开发周期与服务器性能之间的取舍 {#address-tradeoff-between-development-cycle-time-and-server-performance}

Go 正是为解决大规模应用、微服务和云开发中的这些并发需求而诞生的。事实上，云原生计算基金会超过 75% 的项目都用 Go 编写。

Go 的快速构建支持迭代开发，同时减少内存和 CPU 占用，降低了开发效率与运行性能之间的取舍压力。用 Go 构建的服务器几乎可以立即启动，在按量付费和 Serverless 部署中运行成本也更低。

### 应对现代云计算的挑战，提供符合语言惯用方式的标准 API {#address-challenges-with-the-modern-cloud-delivering-standard-idiomatic-apis}

Go 为现代云开发中的许多挑战提供了解决方案：符合语言惯用方式的标准 API，以及充分利用多核处理器的内置并发能力。低延迟和无需繁复调参的特点，让 Go 在性能与生产力之间取得良好平衡，赋予工程团队选择和行动的能力。

## 应用场景 {#use-case .sectionHeading}

### 使用 Go 进行云计算开发 {#use-go-for-cloud-computing}

构建服务时，Go 的优势尤其明显。它的速度和内置并发支持可以带来快速、高效的服务；静态类型、完善的工具，以及对简洁性和可读性的重视，则有助于构建可靠、易维护的代码。

Go 拥有支持服务开发的强大生态系统。[标准库](/pkg/)提供了 HTTP 服务端与客户端、JSON/XML 解析、SQL 数据库以及多种安全和加密功能，满足常见需求。Go 运行时及相关工具还支持[数据竞争检测](/doc/articles/race_detector.html)、[基准测试](/pkg/testing/#hdr-Benchmarks)与性能剖析、代码生成和静态代码分析。

主流云服务商（[GCP](https://cloud.google.com/go/home)、[AWS](https://aws.amazon.com/sdk-for-go/)、[Azure](https://docs.microsoft.com/en-us/azure/go/)）都为其服务提供 Go API。常用开源库还支持 API 工具（[Swagger](https://github.com/go-swagger/go-swagger)）、传输（[Protocol Buffers](https://github.com/golang/protobuf)、[gRPC](https://grpc.io/docs/quickstart/go/)）、监控（[OpenCensus](https://godoc.org/go.opencensus.io)）、对象关系映射（[gORM](https://gorm.io/)）和身份认证（[JWT](https://github.com/dgrijalva/jwt-go)）。开源社区也提供了多种服务框架，例如 [Go Kit](https://gokit.io/)、[Go Micro](https://micro.mu/docs/go-micro.html) 和 [Gizmo](https://github.com/nytimes/gizmo)，帮助你快速起步。

### 云计算领域的 Go 工具 {#go-tools-for-cloud-computing}

{{toolsblurbs `
  - title: Docker
    url: https://www.docker.com/
    iconSrc: /images/logos/docker.svg
    paragraphs:
      - "Docker 是通过容器交付软件的平台即服务。容器将软件、库和配置文件打包，由 Docker Engine 托管，共享一个操作系统内核运行，因此比虚拟机消耗更少的系统资源。"
      - "Docker 支持开发工作流和部署过程，因此云开发者使用它管理 Go 代码，并支持多个平台。"
  - title: Kubernetes
    url: https://kubernetes.io/
    iconSrc: /images/logos/kubernetes.svg
    paragraphs:
      - "Kubernetes 是用 Go 编写的开源容器编排系统，用于自动化部署 Web 应用。如上所述，Web 应用通常通过容器将依赖和配置一起打包。Kubernetes 帮助大规模部署和管理这些容器，让云开发者快速构建、交付和扩展容器化应用，并通过控制容器运行方式的 API 管理不断增长的复杂度。"
`}}

{{projects `
  - company: Google
    url: https://cloud.google.com/go
    logoSrc: google-cloud.svg
    logoSrcDark: google-cloud.svg
    desc: "Google Cloud 在产品和工具生态中广泛使用 Go，包括 Kubernetes、gVisor、Knative、Istio 和 Anthos。Google Cloud 的所有 API 和运行环境都全面支持 Go。"
    ctas:
      - text: "Google Cloud Platform 上的 Go"
        url: https://cloud.google.com/go
  - company: Capital One
    url: https://www.capitalone.com/
    logoSrc: capitalone_light.svg
    logoSrcDark: capitalone_dark.svg
    desc: "Capital One 使用 Go 支撑关键的 Credit Offers API 服务。工程团队也用 Go 构建 Serverless 架构，原因是 Go 的速度与简洁性；他们表示：“没有 Go，我们就不想采用 Serverless。”"
    ctas:
      - text: "信贷优惠 API"
        url: https://medium.com/capital-one-tech/a-serverless-and-go-journey-credit-offers-api-74ef1f9fde7f
  - company: Dropbox
    url: https://www.dropbox.com/
    logoSrc: dropbox.svg
    logoSrcDark: dropbox.svg
    desc: "Dropbox 最初基于 Python 构建，但在 2013 年决定将对性能要求较高的后端迁移到 Go。如今，公司大部分基础设施都用 Go 编写。"
    ctas:
      - text: "Dropbox 的 Go 库"
        url: https://dropbox.tech/infrastructure/open-sourcing-our-go-libraries
  - company: Mercado Libre
    url: https://www.mercadolibre.com.ar/
    logoSrc: mercadolibre_light.svg
    logoSrcDark: mercadolibre_dark.svg
    desc: "MercadoLibre 使用 Go 扩展电商平台。Go 能生成高效的代码，随着电商业务增长轻松扩展，在简化和扩展服务的同时提高开发生产力。"
    ctas:
      - text: "MercadoLibre 与 Go"
        url: /solutions/mercadolibre
  - company: The New York Times
    url: https://www.nytimes.com/
    logoSrc: the-new-york-times-icon.svg
    logoSrcDark: the-new-york-times-icon.svg
    desc: "《纽约时报》采用 Go 来“构建更好的后端服务”。随着 Go 在公司内部的使用扩大，他们需要一套工具包，帮助开发者快速配置和构建微服务 API 以及发布订阅守护进程，并将其开源。"
    ctas:
      - text: "纽约时报：Gizmo"
        url: https://open.nytimes.com/introducing-gizmo-aa7ea463b208
      - text: "Gizmo 的 GitHub 仓库"
        url: https://github.com/nytimes/gizmo
  - company: Twitch
    url: https://www.twitch.tv/
    logoSrc: twitch.svg
    logoSrcDark: twitch.svg
    desc: "Twitch 使用 Go 支撑许多最繁忙的系统，为数百万用户提供直播视频和聊天服务。"
    ctas:
      - text: "Go 向低延迟垃圾回收迈进"
        url: https://blog.twitch.tv/en/2016/07/05/gos-march-to-low-latency-gc-a6fa96f06eb7/
  - company: Uber
    url: https://www.uber.com/
    logoSrc: uber_light.svg
    logoSrcDark: uber_dark.svg
    desc: "Uber 使用 Go 支撑多项关键服务，影响着全球数百万司机和乘客的体验，包括实时分析引擎 AresDB、地理查询微服务 Geofence，以及资源调度器 Peloton。"
    ctas:
      - text: AresDB
        url: https://eng.uber.com/aresdb/
      - text: Geofence
        url: https://eng.uber.com/go-geofence/
      - text: Peloton
        url:  https://eng.uber.com/open-sourcing-peloton/
`}}

## 开始使用 {#get-started .sectionHeading}

### 云计算领域的 Go 书籍 {#go-books-for-cloud-computing}

{{books `
  - title: "用 Go 构建微服务（Building Microservices with Go）"
    url: https://www.amazon.com/Building-Microservices-Go-efficient-microservices/dp/1786468662/
    thumbnail: /images/books/building-microservices-with-go.jpg
  - title: "Go 软件架构实战（Hands-On Software Architecture with Golang）"
    url: https://www.amazon.com/dp/1788622596/ref=cm_sw_r_tw_dp_U_x_-aZWDbS8PD7R4
    thumbnail: /images/books/hands-on-software-architecture-with-golang.jpg
  - title: "用 Go 构建 RESTful Web 服务（Building RESTful Web services with Go）"
    url: https://www.amazon.com/Building-RESTful-Web-services-gracefully-ebook/dp/B072QB8KL1
    thumbnail: /images/books/building-restful-web-services-with-go.jpg
  - title: "精通 Go Web 服务（Mastering Go Web Services）"
    url: https://www.amazon.com/Mastering-Web-Services-Nathan-Kozyra/dp/178398130X
    thumbnail: /images/books/mastering-go-web-services.jpg
`}}

{{libraries `
  - title: "Web 框架"
    viewMoreUrl: https://pkg.go.dev/search?q=web+framework
    items:
      - text: Echo
        url: https://echo.labstack.com/
        desc: "高性能、可扩展、极简的 Go Web 框架"
      - text: Flamingo
        url: https://www.flamingo.me/
        desc: "基于 Go 的快速开源框架，具有简洁、可扩展的架构"
      - text: Gin
        url: https://gin-gonic.com/
        desc: "用 Go 编写的 Web 框架，提供类似 Martini 的 API。"
      - text: Gorilla
        url: https://www.gorillatoolkit.org/
        desc: "Go 编程语言的 Web 工具包。"
  - title: "路由器"
    viewMoreUrl: https://pkg.go.dev/search?q=http%20router
    items:
      - text: net/http
        url: https://pkg.go.dev/net/http
        desc: "标准库中的 HTTP 包"
      - text: julienschmidt/httprouter
        url: https://pkg.go.dev/github.com/julienschmidt/httprouter?tab=overview
        desc: "轻量级高性能 HTTP 请求路由器"
      - text: gorilla/mux
        url: https://pkg.go.dev/github.com/gorilla/mux?tab=overview
        desc: "强大的 HTTP 路由器和 URL 匹配器，用于构建 Go Web 服务器 🦍"
      - text: Chi
        url: https://pkg.go.dev/github.com/go-chi/chi?tab=overview
        desc: "轻量级、符合 Go 惯用方式且可组合的路由器，用于构建 Go HTTP 服务。"
  - title: "模板引擎"
    viewMoreUrl: https://pkg.go.dev/search?q=templates
    items:
      - text: html/template
        url: https://pkg.go.dev/html/template
        desc: "标准库中的 HTML 模板引擎"
      - text: flosch/pongo2
        url: https://pkg.go.dev/github.com/flosch/pongo2?tab=overview
        desc: "语法类似 Django 的模板语言"
  - title: "数据库与驱动"
    viewMoreUrl: https://pkg.go.dev/search?q=database%20OR%20sql
    items:
      - text: database/sql
        url: https://pkg.go.dev/database/sql
        desc: "标准库中的数据库接口，具有 MySQL、Postgres、Oracle、MS SQL、BigQuery 及大多数 SQL 数据库的驱动支持"
      - text: mongo-driver/mongo
        url: https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?tab=overview
        desc: "MongoDB 官方支持的 Go 驱动"
      - text: elastic/go-elasticsearch
        url: https://pkg.go.dev/github.com/elastic/go-elasticsearch/v8?tab=overview
        desc: "Go 的 Elasticsearch 客户端"
      - text: GORM
        url: https://gorm.io/
        desc: "Go 的 ORM 库"
      - text: Bleve
        url: https://blevesearch.com/
        desc: "Go 的全文检索与索引库"
      - text: CockroachDB
        url: https://www.cockroachlabs.com/
        desc: "面向云环境设计的数据库，提供具有韧性、一致性且可大规模扩展的分布式 SQL"
  - title: "Web 库"
    viewMoreUrl: https://pkg.go.dev/search?q=web
    items:
      - text: markbates/goth
        url: https://pkg.go.dev/github.com/markbates/goth?tab=overview
        desc: "Web 应用身份认证"
      - text: jinzhu/gorm
        url: https://pkg.go.dev/github.com/jinzhu/gorm?tab=overview
        desc: "Go 的 ORM 库"
      - text: dgrijalva/jwt-go
        url: https://pkg.go.dev/github.com/dgrijalva/jwt-go?tab=overview
        desc: "JSON Web Token 的 Go 实现"
  - title: "其他项目"
    items:
      - text: gopherjs
        url: https://pkg.go.dev/github.com/gopherjs/gopherjs?tab=overview
        desc: "将 Go 编译为 JavaScript 的编译器，让开发者用 Go 编写可在所有浏览器中运行的前端代码。"
`}}
