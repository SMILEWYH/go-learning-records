//go:build ignore

package main

// 练习 1:聚合一批接口响应状态
// 背景:后台监控要统计一批请求中每个状态码的出现次数及服务端错误总数
// 要求:
//  1. 输入为 []int, 用 range 遍历
//  2. 用 map[int]int 累加每种状态码的次数;map 必须先初始化
//  3. 单独累计 500 <= code < 600 的数量
//
// 常规:[200, 200, 404, 503, 503, 502] -> 200:2、404:1、503:2、502:1, 服务端错误 3 次
// 边界:空切片或 nil 切片 -> map 长度 0, 服务端错误 0 次
// 本题约定输入都是有效的 100..599 状态码;不要求打印 map 时具有固定顺序
// 建议实现位置:main();后续可抽取 countStatuses(codes []int) (map[int]int, int)
// 提示:尚未出现的 key 读取结果是 int 零值 0, 可以直接参与累加

// 练习 2:给批量提交的用户 ID 去重
// 背景:前端合并多次勾选时, 同一个 ID 可能重复出现, 也可能混入空字符串
// 要求:
//  1. 忽略空字符串, 使用 continue 跳过
//  2. 用 map[string]bool 记录已经出现的 ID
//  3. 用切片收集结果, 保留每个 ID 第一次出现的顺序
//
// 常规:["u-02", "", "u-01", "u-02", "u-03", "u-01"] -> ["u-02", "u-01", "u-03"]
// 边界:[] 或 ["", ""] -> 长度为 0 的结果;["u-01", "u-01"] -> ["u-01"]
// 本题不裁剪空格, 所以 " u-01 " 与 "u-01" 被视为不同 ID
// 建议实现位置:main();后续可抽取 uniqueIDs(ids []string) []string
// 提示:从原切片的遍历顺序构造结果, 不能最后 range map 来还原首次出现顺序

// 练习 3:模拟有上限的任务重试
// 背景:后端请求失败时可以重试, 但必须限制次数, 并在成功后立刻停止
// 要求:输入 outcomes []bool 和 maxAttempts int, 每个 bool 代表一次尝试是否成功
//  1. 最多读取 maxAttempts 项, 也不能超出 outcomes 的长度
//  2. 第一次读到 true 就停止, 并记录 succeeded=true
//  3. 输出真正执行的 attempts 数和 succeeded;本题不访问网络、不延时
//
// 常规:[false, false, true, false]、maxAttempts=4 -> attempts=3、succeeded=true
// 达到上限:[false, false, true]、maxAttempts=2 -> attempts=2、succeeded=false
// 边界:[]、maxAttempts=3 -> 0,false;[true]、maxAttempts=0 或 -1 -> 0,false
// 建议实现位置:main();后续可抽取 simulateRetry(outcomes []bool, maxAttempts int) (int, bool)
// 提示:使用有界 for, 把“是否还有输入”和“是否还有尝试次数”同时写入条件

func main() {
}
