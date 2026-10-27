//go:build ignore

package main

// 练习 1:规范化接口分页大小
// 背景:前端可能不传 pageSize, 也可能传来过大的数值
// 要求:用 if / else if / else 把 requestedSize 转为有效的 pageSize
//   - requestedSize <= 0:采用默认值 20
//   - requestedSize > 100:限制为 100
//   - 其余数值:保持原值
//
// 常规:requestedSize=35 -> 35
// 边界:0 -> 20;-1 -> 20;100 -> 100;101 -> 100
// 建议实现位置:main(), 修改输入后多次运行, 不要求写循环
// 学完函数后可抽取 normalizePageSize(requestedSize int) int
// 提示:不要让有效值 100 被误判为非法, 也不要照搬 TS 的 requestedSize || 20

// 练习 2:做一次明确的权限判断
// 背景:后台页面是否展示按钮不能替代后端权限校验
// 要求:输入 loggedIn bool、role string、isOwner bool, 输出一个状态码
//   - 未登录, 无论其他条件如何, 返回状态码 401
//   - 已登录, 且 role == "admin" 或 isOwner == true, 返回 200
//   - 其他已登录用户, 返回 403
//
// 常规:true、"editor"、true -> 200;true、"viewer"、false -> 403
// 边界:false、"admin"、true -> 401;true、"admin"、false -> 200
// 建议实现位置:main();用 if / else 分支输出结果, 避免提前 return 跳过后续题目
// 学完函数后可抽取 accessStatus(loggedIn bool, role string, isOwner bool) int
// 提示:优先判断是否登录;逻辑组合时用括号明确自己的意图

// 练习 3:识别配置项的三种状态
// 背景:自定义页面标题可能“未配置”“明确配置为空”或“配置了文本”
// 要求:使用 map[string]string 和 comma-ok 查询 "title", 再用 if 判断
//   - key 不存在:显示默认标题 "Dashboard"
//   - key 存在, 值为空:显示 "(hidden)", 表示用户选择隐藏标题
//   - key 存在且非空:显示配置值
//
// 常规:map[string]string{"title":"Orders"} -> "Orders"
// 边界:空 map -> "Dashboard";{"title":""} -> "(hidden)"
// 额外边界:nil map -> "Dashboard";本题只读 map, 不向 nil map 写入
// 建议实现位置:main(), 可以把 value, ok 写在 if 的初始化语句中
// 学完函数后可抽取 resolveTitle(config map[string]string) string

func main() {
}
