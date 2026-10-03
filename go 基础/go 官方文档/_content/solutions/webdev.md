---
title: "Go 与 Web 开发"
linkTitle: "Web 开发"
description: "凭借出色的内存性能和对多种 IDE 的支持，Go 为快速、可扩展的 Web 应用提供支持。"
date: 2019-10-04T15:26:31-04:00
series: Use Cases
template: true
books:
icon:
  file: webdev-green.svg
  alt: "Web 开发图标"
iconDark:
  file: webdev-white.svg
  alt: "Web 开发图标"
---

## 概述 {#overview .sectionHeading}

### Go 为 Web 应用提供速度、安全性和对开发者友好的工具 {#go-delivers-speed-security-and-developer-friendly-tools-for-web-applications}

Go 旨在让开发者快速构建可扩展、安全的 Web 应用。它自带易用、安全、高性能的 Web 服务器，以及自己的 Web 模板库。Go 很好地支持各种新技术：从 [HTTP/2](https://pkg.go.dev/net/http)，到 [MySQL](https://pkg.go.dev/mod/github.com/go-sql-driver/mysql)、[MongoDB](https://pkg.go.dev/mod/go.mongodb.org/mongo-driver)、[Elasticsearch](https://pkg.go.dev/mod/github.com/elastic/go-elasticsearch/v8) 等数据库，再到包括 [TLS 1.3](https://pkg.go.dev/crypto/tls) 在内的加密标准。得益于出色的可移植性，Go Web 应用可以原生运行在 [Google App Engine](https://cloud.google.com/appengine/) 和 [Google Cloud Run](https://cloud.google.com/run/) 上以便扩展，也可以运行在其他环境、云平台或操作系统中。

## 主要优势 {#key-benefits .sectionHeading}

### 快速进行跨平台部署 {#deploy-across-platforms-in-record-speed}

Go 支持快速跨平台部署，因此受到企业青睐。借助 goroutine、原生编译以及基于 URI 的包命名空间，Go 代码能够编译为一个小型、无外部依赖的二进制文件，运行速度很快。

### 借助 Go 开箱即用的性能轻松扩展 {#leverage-gos-out-of-the-box-performance-to-scale-with-ease}

Hexact Inc. 联合创始人兼首席技术官 Tigran Bayburtsyan 总结了公司转向 Go 的五个主要原因：

- **编译为单个二进制文件**：“Go 通过静态链接，根据操作系统类型和体系结构，将所有依赖库和模块合并到一个二进制文件中。”

- **静态类型系统**：“类型系统对大规模应用非常重要。”

- **性能**：“Go 的并发模型和 CPU 扩展能力带来了更好的性能。需要处理内部请求时，我们使用独立的 goroutine，其资源消耗仅为 Python 线程的十分之一。”

- **无需 Web 框架**：“在大多数情况下，确实不需要任何第三方库。”

- **出色的 IDE 支持和调试能力**：“将所有项目用 Go 重写后，代码量比之前减少了 64%。”


{{projects `
  - company: Caddy
    url: https://caddyserver.com/
    logoSrc: caddy.svg
    logoSrcDark: caddy.svg
    desc: "Caddy 2 是用 Go 编写的开源 Web 服务器，功能强大，适合企业使用，并自动支持 HTTPS。相比用 C 编写的服务器，Caddy 提供更好的内存安全性。由 Go 标准库提供支持的强化 TLS 协议栈，承载了互联网流量中的相当一部分。"
    ctas:
      - text: Caddy 2
        url: https://caddyserver.com/
  - company: Cloudflare
    url: https://www.cloudflare.com/en-gb/
    logoSrc: cloudflare-icon.svg
    logoSrcDark: cloudflare-icon.svg
    desc: "Cloudflare 为数百万个网站、API、SaaS 服务及其他互联网资源加速并提供保护。“Go 是 CloudFlare 服务的核心，涵盖高延迟 HTTP 连接的压缩处理、整个 DNS 基础设施、SSL、负载测试等。”"
    ctas:
      - text: "Cloudflare 与 Go"
        url: https://blog.cloudflare.com/what-weve-been-doing-with-go/
  - company: gov.uk
    url: https://gov.uk/
    logoSrc: govuk_light.svg
    logoSrcDark: govuk_dark.svg
    desc: "Go 的简洁性和安全性非常适合英国政府的 HTTP 基础设施。对出色的 net/http 包进行简短试验后，Web 开发者确信自己选对了方向。“尤其是 Go 的并发模型，让构建高性能 I/O 密集型应用变得极其简单。”"
    ctas:
      - text: "为 gov.uk 构建新的路由器"
        url: https://technology.blog.gov.uk/2013/12/05/building-a-new-router-for-gov-uk/
      - text: "政府机构中的 Go 实践"
        url: https://technology.blog.gov.uk/2014/11/14/using-go-in-government/
  - company: Hugo
    url: https://gohugo.io/
    logoSrc: hugo.svg
    logoSrcDark: hugo.svg
    desc: "Hugo 是用 Go 编写的快速、现代的网站引擎，旨在让网站创建重新充满乐趣。用 Hugo 构建的网站速度极快、安全，并且无需任何依赖即可托管到各种环境。"
    ctas:
      - text: Hugo
        url: https://gohugo.io/
  - company: Mattermost
    url: https://mattermost.com/
    logoSrc: mattermost_light.svg
    logoSrcDark: mattermost_dark.svg
    desc: "Mattermost 是灵活的开源消息平台，支持安全的团队协作，使用 Go 和 React 编写。"
    ctas:
      - text: Mattermost
        url: https://mattermost.com/
  - company: Medium
    url: https://medium.org/
    logoSrc: medium_light.svg
    logoSrcDark: medium_dark.svg
    desc: "Medium 使用 Go 支撑社交图谱、图片服务器及多项辅助服务。“我们发现 Go 很容易构建、打包和部署。我们喜欢它提供类型安全，同时不需要 Java 那样冗长的代码和 JVM 调优。”"
    ctas:
      - text: "Medium 的 Go 服务"
        url: https://medium.engineering/how-medium-goes-social-b7dbefa6d413
  - company: The Economist
    url: https://economist.com/
    logoSrc: economist.svg
    logoSrcDark: economist.svg
    desc: "《经济学人》需要更灵活地向日益多样化的数字渠道交付内容。用 Go 编写的服务成为新系统的关键组件，使其能够提供可扩展、高性能的服务，并快速迭代新产品。“总体而言，我们认为 Go 最适合在分布式云系统中兼顾易用性与效率。”"
    ctas:
      - text: "《经济学人》的 Go 微服务"
        url: https://www.infoq.com/articles/golang-the-economist/
`}}

## 开始使用 {#get-started .sectionHeading}

### Web 开发领域的 Go 书籍 {#go-books-on-web-development}

{{books `
  - title: "使用 Go 开发 Web 应用（Web Development with Go）"
    url: https://www.amazon.com/Web-Development-Go-Building-Scalable-ebook/dp/B01JCOC6Z6
    thumbnail: /images/books/web-development-with-go.jpg
  - title: "Go Web 编程（Go Web Programming）"
    url: https://www.amazon.com/Web-Programming-Sau-Sheong-Chang/dp/1617292567
    thumbnail: /images/books/go-web-programming.jpg
  - title: "Web 开发实例集：使用 Go 构建全栈 Web 应用（Web Development Cookbook）"
    url: https://www.amazon.com/Web-Development-Cookbook-full-stack-applications-ebook/dp/B077TVQ28W
    thumbnail: /images/books/go-web-development-cookbook.jpg
  - title: "用 Go 构建 RESTful Web 服务（Building RESTful Web services with Go）"
    url: https://www.amazon.com/Building-RESTful-Web-services-gracefully-ebook/dp/B072QB8KL1
    thumbnail: /images/books/building-restful-web-services-with-go.jpg
  - title: "精通 Go Web 服务（Mastering Go Web Services）"
    url: https://www.amazon.com/Mastering-Web-Services-Nathan-Kozyra-ebook/dp/B00W5GUKL6
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

### 课程 {#courses}

* [学习用 Go 创建 Web 应用](https://www.usegolang.com)：付费在线课程。

### 项目 {#projects}

* {{pkg "github.com/gopherjs/gopherjs" "gopherjs"}}：将 Go 编译为 JavaScript 的编译器，让开发者用 Go 编写可在所有浏览器中运行的前端代码。
* [Hugo](https://gohugo.io/)：用于构建网站的高速框架，原文称其为全球最快。
* [Mattermost](https://mattermost.com/)：灵活的开源消息平台，支持安全的团队协作。
* [Caddy](https://caddyserver.com/)：用 Go 编写、自动支持 HTTPS 的强大开源 Web 服务器，适合企业使用。
