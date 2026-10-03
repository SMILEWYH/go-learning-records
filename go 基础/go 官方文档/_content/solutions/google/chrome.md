---
title: "Chrome 内容优化服务运行在 Go 之上"
company: Chrome
logoSrc: chrome.svg
logoSrcDark: chrome.svg
heroImgSrc: go_chrome_case_study.png
series: Case Studies
template: true
quote: |
  Google Chrome 内置 Google 的智能技术，是比以往更简洁、更安全、更快速的网页浏览器。

  在本案例中，Chrome Optimization Guide 团队分享了他们如何尝试 Go、快速上手，以及未来继续使用 Go 的计划。
---

提到 Chrome，你可能只会想到用户安装的浏览器。然而在幕后，Chrome 拥有庞大的后端服务体系，其中就包括 Chrome Optimization Guide 服务。这项服务是 Chrome 用户体验策略的重要基础，处在影响用户体验的关键路径上，并且使用 Go 实现。

Chrome Optimization Guide 服务旨在将 Google 的能力带入 Chrome：它向已安装的浏览器提供提示，说明页面加载时可以执行哪些优化，以及何时应用这些优化最有效。它由实时服务器和批量日志分析系统共同组成。

所有使用 Chrome 精简模式的用户，都通过以下机制从该服务接收数据：推送数据块，为所在地区的知名网站提供优化提示；向 Google 服务器发起请求，获取特定用户经常访问的主机的提示；对于设备上尚无提示的页面加载，则按需获取。如果 Chrome Optimization Guide 服务突然消失，用户可能会发现页面加载速度和浏览网页时的数据消耗量发生明显变化。

{{backgroundquote `
  author: Sophie Chang
  title: 软件工程师
  quote: |
    既然 Go 已经在我们的实践中取得成功，我们计划继续在合适的场景中使用它。
`}}

Chrome 工程团队开始构建这项服务时，只有少数成员熟悉 Go。团队大多数人更熟悉 C++，但他们发现，搭建 C++ 服务器所需的复杂样板代码负担过重。团队表示，Go 的“简洁性、易于快速上手和生态系统，让我们很有动力去学习它”，而且“我们的探索得到了回报”。数百万用户依靠这项服务获得更好的 Chrome 体验，因此选择 Go 绝不是一个小决定。基于迄今的经验，团队还表示：“既然 Go 已经在我们的实践中取得成功，我们计划继续在合适的场景中使用它。”

除了 Chrome Optimization Guide 团队，Google 的许多工程团队也已在开发过程中采用 Go。你可以继续了解 [Core Data Solutions](/solutions/google/coredata/) 和 [Firebase Hosting](/solutions/google/firebase/) 团队如何使用 Go 大规模构建快速、可靠、高效的软件。

*编者注：Go 团队感谢 Sophie Chang 为本文作出的贡献。*
