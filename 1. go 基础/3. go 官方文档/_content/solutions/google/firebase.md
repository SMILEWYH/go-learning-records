---
title: "Firebase Hosting 团队如何借助 Go 扩展规模"
company: Firebase
logoSrc: firebase.svg
logoSrcDark: firebase.svg
heroImgSrc: go_firebase_case_study.png
series: Case Studies
template: true
quote: |
  Firebase 是 Google 的移动开发平台，帮助你快速开发高质量应用并拓展业务。

  Firebase Hosting 团队分享了他们使用 Go 的历程，包括从 Node.js 迁移后端、新 Go 开发者如何轻松上手，以及 Go 如何帮助他们扩展规模。
---

Firebase Hosting 团队为 Google Cloud 客户提供静态网站托管服务。他们的静态网站托管服务位于全球内容分发网络之后，并为用户提供易用的工具。团队还开发了从上传网站文件、注册域名到跟踪用量的各种功能。

加入 Google 之前，Firebase Hosting 的技术栈使用 Node.js 编写。当团队需要与 Google 的其他多项服务互操作时，他们开始使用 Go。团队知道“并发仍将是一项重要需求”，因此决定用 Go 帮助他们轻松、高效地扩展规模。团队的软件工程师 Michael Bleigh 表示，他们“相信 Go 的性能会更好”，而且与当时考虑的其他语言相比，他们“喜欢 Go 更简洁的特点”。

从一个用 Go 编写的小型服务开始，团队分阶段迁移了整个后端。他们逐步确定想要实现的大型功能，在实现过程中用 Go 重写相关部分，并迁移到 Google Cloud 和 Google 内部的集群管理系统。**如今，Firebase Hosting 团队已用 Go 替换了全部后端 Node.js 代码。**

团队最初只有一名工程师具有 Go 开发经验。Bleigh 说：“通过同事间的互相学习，再加上 Go 本身容易上手，如今团队中的每个人都有 Go 开发经验。”他们发现，虽然大多数新成员此前没有用过 Go，但“绝大多数人在几周内就能高效开展工作”。

Bleigh 代表团队表示：“使用 Go 时，很容易看清代码是如何组织的，以及代码在做什么。Go 总体上非常易读、易懂。由于语言中的惯用写法，错误处理、接收者和接口都容易理解。”

随着规模扩大，并发始终是团队关注的重点。软件工程师 Robert Rossney 表示：“Go 让我们很容易将所有复杂的并发逻辑集中在一个地方，并在其他地方通过抽象使用它。”Rossney 还谈到了使用一种在设计时就考虑并发的语言所带来的好处：“Go 中实现并发的方式也很多。我们需要学习各种方法最适合什么场景、如何判断一个问题是否属于并发问题、如何调试——而这些学习需求，正是因为你确实可以用 Go 代码写出这些模式。”

{{backgroundquote `
  author: Robert Rossney
  title: 软件工程师
  quote: |
    总体而言，团队没有因为 Go 而感到受挫的时候。它只是安静地做好自己的事，让你专心工作。
`}}

数十万客户使用 Firebase Hosting 托管网站，这意味着 Go 代码每天要处理数十亿次请求。Bleigh 分享道：“迁移到 Go 之后，我们的客户数量和流量已经多次翻倍，却从未需要精细的性能调优。”使用 Go 后，软件性能和团队效率都有所提升，生产力显著提高。Rossney 说：“总体而言，……团队没有因为 Go 而感到受挫的时候。它只是安静地做好自己的事，让你专心工作。”

除了 Firebase Hosting 团队，Google 的许多工程团队也已在开发过程中采用 Go。你可以继续了解 [Core Data Solutions](/solutions/google/coredata/) 和 [Chrome](/solutions/google/chrome/) 团队如何使用 Go 大规模构建快速、可靠、高效的软件。
