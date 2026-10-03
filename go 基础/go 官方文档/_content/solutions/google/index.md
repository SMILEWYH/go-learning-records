---
title: 'Google 的 Go 实践'
date: 2020-08-27
company: Google
logoSrc: google.svg
logoSrcDark: google.svg
heroImgSrc: go_core_data_case_study.png
carouselImgSrc: go_google_case_study_carousel.png
series: Case Studies
type: solutions
template: true
description: |-
  Google 是一家技术公司，其使命是整合全球信息，使人人皆可访问并从中受益。

  Go 于 2007 年诞生于 Google，旨在多核联网计算机和大型代码库的时代提升编程效率。自 2009 年公开发布以来，十多年间，Go 在 Google 内部的使用规模显著增长。
quote: Go 于 2007 年诞生于 Google，此后 Google 的各个工程团队陆续采用 Go，构建大规模产品和服务。
---

{{pullquote `
  author: Rob Pike
  quote: |
    Go 始于 2007 年 9 月。当时，Robert Griesemer、Ken Thompson 和我开始讨论一门新语言，希望解决我们和 Google 同事在日常工作中面临的工程挑战。

    2009 年 11 月首次向公众发布 Go 时，我们并不知道它是否会被广泛采用，也不知道它能否影响未来的语言。从 2020 年回望，Go 在这两个方面都取得了成功：它在 Google 内外得到广泛使用；它处理网络并发和软件工程问题的方式，也对其他语言及其工具产生了明显影响。

    Go 的影响范围远远超出了我们的预期。它在行业中的发展十分迅速，也支撑了 Google 的许多项目。
`}}

下面的故事，只是 Google 使用 Go 的众多方式中的一小部分。

### Google Core Data Solutions 团队如何使用 Go {#how-googles-core-data-solutions-team-uses-go}

Google 的使命是“整合全球信息，使人人皆可访问并从中受益”。Google Core Data Solutions 团队正是负责组织这些信息的团队之一。除了其他工作，该团队还维护对全球网页建立索引的服务。这些网页索引服务通过保持搜索结果的及时性和完整性，为 Google 搜索等产品提供支持，而它们使用 Go 编写。

[了解更多](/solutions/google/coredata/)

---

### Chrome 内容优化服务运行在 Go 之上 {#chrome-content-optimization-service-runs-on-go}

提到 Chrome，你可能只会想到用户安装的浏览器。然而在幕后，Chrome 拥有庞大的后端服务体系，其中就包括 Chrome Optimization Guide 服务。这项服务是 Chrome 用户体验策略的重要基础，处在影响用户体验的关键路径上，并且使用 Go 实现。

[了解更多](/solutions/google/chrome/)

---

### Firebase Hosting 团队如何借助 Go 扩展规模 {#how-the-firebase-hosting-team-scaled-with-go}

Firebase Hosting 团队为 Google Cloud 客户提供静态网站托管服务。他们的静态网站托管服务位于全球内容分发网络之后，并为用户提供易用的工具。团队还开发了从上传网站文件、注册域名到跟踪用量的各种功能。

[了解更多](/solutions/google/firebase/)

---

### 驱动 Google 生产环境：Google 站点可靠性工程团队如何使用 Go {#actuating-google-production-how-googles-site-reliability-engineering-team-uses-go}

Google 运行着少量规模极大的服务。这些服务由全球基础设施支撑，涵盖所需的一切：存储系统、负载均衡器、网络、日志、监控等等。然而，这不是一个静止的系统，也不可能是。架构不断演进，新产品和新想法不断出现，新版本需要发布，配置需要推送，数据库模式需要更新，还有更多事情要做。最终，我们每秒都要向系统部署数十次变更。

[了解更多](/solutions/google/sitereliability/)
