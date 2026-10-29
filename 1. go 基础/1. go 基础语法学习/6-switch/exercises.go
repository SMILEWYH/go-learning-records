//go:build ignore

package main

// 练习 1:选择 HTTP 方法对应的业务动作
// 背景:一个简化的 API 路由层需要把请求方法映射成动作名称
// 要求:使用 switch method, 遵循以下规则
//   - GET / HEAD -> "read"
//   - POST -> "create"
//   - PUT / PATCH -> "update"
//   - DELETE -> "delete"
//   - 其他值 -> "unsupported"
//
// 常规:"PATCH" -> "update";"HEAD" -> "read"
// 边界:""、"get"、"TRACE" -> "unsupported", 本题不自动转为大写
// 建议实现位置:main(), 多值 case 用逗号列出, 不需要 fallthrough
// 学完函数后可抽取 actionForMethod(method string) string

// 练习 2:给 HTTP 状态码分类
// 背景:监控面板需要统计不同响应类别
// 要求:使用不带表达式的 switch, 按范围输出类别
//   - 100..199 -> "informational";200..299 -> "success"
//   - 300..399 -> "redirect";400..499 -> "client_error"
//   - 500..599 -> "server_error";其余 -> "invalid"
//
// 常规:204 -> "success";404 -> "client_error";503 -> "server_error"
// 边界:99 -> "invalid";100 -> "informational";599 -> "server_error";600 -> "invalid"
// 建议实现位置:main();修改 status 的值逐一检查分界点即可
// 学完函数后可抽取 classifyStatus(status int) string
// 提示:写成一边包含、一边不包含的区间, 检查有没有遗漏或重叠

// 练习 3:解释异步导出任务的状态
// 背景:前端查询后端导出任务时, 需要知道状态是否结束, 以及下一步提示
// 要求:根据 state string, 同步设置 terminal bool 和 message string
//   - "queued" / "running" -> false, "稍后查询"
//   - "succeeded" -> true, "可以下载"
//   - "failed" / "canceled" -> true, "任务未完成"
//   - 其他值 -> false, "未知状态"
//
// 常规:"running" -> false, "稍后查询";"succeeded" -> true, "可以下载"
// 边界:"" -> false, "未知状态";"canceled" -> true, "任务未完成"
// 建议实现位置:main(), 先声明结果变量, 再在 switch 中赋值
// 学完函数后可抽取 describeJob(state string) (terminal bool, message string)
// 思考:此题使用 fallthrough 会不会误把未结束的状态标记为结束？

func main() {
}
