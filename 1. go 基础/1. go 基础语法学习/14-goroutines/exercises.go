//go:build ignore

package main

/*
练习 1：并发拼装仪表盘摘要。

背景：页面需要用户名称和未读消息数，它们可以由两个独立任务获取。
本题用固定数据模拟服务响应，不访问网络，也不使用 Sleep 模拟等待。

要求：
  - 定义 DashboardSummary{UserName string; Unread int}。
  - 启动两个 goroutine：一个填写 UserName，一个填写 Unread。
  - 增加 sync 导入，用 WaitGroup 在启动前 Add(2)，各任务 defer Done()。
  - 每个任务只写自己的字段，主 goroutine 只在 Wait 后读取整个结果。
  - Wait 返回后才返回摘要；不在后台任务里打印以免输出顺序不稳定。

建议签名：func buildDashboard(name string, unread int) DashboardSummary

手动核对：
  - ("Lin",3)：返回 UserName=Lin、Unread=3。
  - ("",0)：返回字段零值，两个任务仍须完成，不挂起。
  - 连续调用十次，结果一致；可以用 -race 运行检查本次执行。
*/

/*
练习 2：保持输入顺序的批量价格计算。

背景：批量导入商品时，可以分别计算每个商品加固定费用后的价格。

要求：
  - 为每个输入价格启动一个 goroutine，结果为 price+fee。
  - 输入不修改，预先创建等长结果 slice，每个任务只写自己的索引。
  - 不并发 append，不在任务里修改同一个累加变量。
  - 调用前登记任务数量，全部 Done 后返回结果。
  - 本题输入规模很小且数值不会溢出；真实大批量任务需限制并发数。

建议签名：func addFees(prices []int, fee int) []int

手动核对：
  - [100,250,0]、fee=20：返回 [120,270,20]，顺序不变。
  - [100]、fee=0：返回 [100]。
  - nil 或 []：返回长度 0 的结果，Wait 不阻塞。
*/

/*
练习 3：批量校验注册资料，提前返回也报告完成。

背景：校验任务遇到空用户名时会直接拒绝，但整个批次仍需等待全部任务。

要求：
  - 输入 []string，每项一个 goroutine，输出等长 []bool。
  - 用户名非空且字节长度至少 3 为 true，否则 false；本题输入限定 ASCII。
  - 空字符串分支必须使用提前 return，仍保证执行 defer Done()。
  - 每项只写自己结果槽，不共享追加错误列表。
  - 返回后在主 goroutine 顺序统计有效用户数。

建议签名：func validateUsernames(names []string) []bool

手动核对：
  - ["lin","","ab","zoe"]：返回 [true,false,false,true]，有效数 2。
  - 全为空字符串：全部 false，程序正常结束。
  - nil：长度为 0，有效数 0。
  - 使用 go run -race exercises.go，确认实现未发现数据竞争。
*/

func main() {
}
