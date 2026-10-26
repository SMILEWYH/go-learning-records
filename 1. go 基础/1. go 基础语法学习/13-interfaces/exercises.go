//go:build ignore

package main

/*
练习 1：替换通知发送方式。

背景：账户服务要发站内信或邮件，但业务代码只需要“发送文本”能力。
本题只返回模拟结果，不连接真实邮件服务。

要求：
  - 定义 Notifier 接口，包含 Send(message string) string。
  - 定义 InboxNotifier 和 EmailNotifier，各自隐式实现接口。
  - Inbox 返回“inbox:<message>”，Email 返回“email:<message>”。
  - notifyAll 遍历 []Notifier，顺序收集结果并返回 []string。
  - 不用 type switch 判断实现类型；通过接口直接调用 Send。

建议签名：func notifyAll(notifiers []Notifier, message string) []string

手动核对：
  - [InboxNotifier{}, EmailNotifier{}]、"welcome"：返回 ["inbox:welcome","email:welcome"]。
  - 空消息：分别返回 "inbox:" 和 "email:"。
  - 空切片：结果长度为 0；本题约定切片内不含 nil 接口。
*/

/*
练习 2：校验来源不统一的分页大小。

背景：某个内部聚合层提供 any 值；它可能是整数、文本、布尔值或 nil。
在进入业务逻辑之前，需要明确允许哪些类型。

要求：
  - int 类型：仅接受 1 到 100。
  - string 类型：本题仅识别 "small"->10、"large"->50，其余拒绝。
  - 其他类型（包括 float64、bool、nil）均拒绝。
  - 成功返回值与 true；失败返回 0 与 false。
  - 用 type switch 完成；不要把类型断言当作数值转换。

建议签名：func parseLimit(value any) (int, bool)

手动核对：
  - 20 -> 20,true；"small" -> 10,true；"large" -> 50,true。
  - 0、101、"20"、float64(20)、true、nil -> 0,false。
  - 在注释中说明为什么 any(float64(20)).(int) 不是安全转换。
*/

/*
练习 3：正确表示“没有启用的功能配置”。

背景：配置工厂返回接口；未启用时，调用方希望用 config == nil 判断。

要求：
  - 定义 Feature 接口，含 Name() string。
  - 定义 FeatureConfig，给 *FeatureConfig 实现 Name，正常返回 "export"。
  - createFeature(false) 直接返回 nil 接口；true 返回有效的 *FeatureConfig。
  - 在 main 函数另建 var missing *FeatureConfig，再赋给 Feature，比较它与 nil。
  - 不调用这个 nil 指针的方法，只解释其接口动态类型和值。

建议签名：func createFeature(enabled bool) Feature

手动核对：
  - createFeature(false) == nil 为 true。
  - createFeature(true) != nil，调用 Name() 返回 export。
  - Feature(missing) == nil 为 false。
  - 写出为什么 FeatureConfig{} 不满足使用指针接收者的 Feature 接口。
*/

func main() {
}
