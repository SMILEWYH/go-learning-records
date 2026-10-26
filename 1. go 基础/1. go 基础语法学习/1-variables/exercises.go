//go:build ignore

package main

// 练习 1：整理本地 API 服务配置。
// 背景：前端转到后端后，需要先理解端口、部署环境和调试开关的类型。
// 要求：
//  1. 用 const 声明默认端口 8080，用 var 声明 serviceName = "profile-api"。
//  2. 用 := 声明 environment = "development"、debug = true。
//  3. 创建可修改的 port 变量并从默认端口初始化，随后改为 9000。
//  4. 用 fmt.Println 输出这些值，并用 %T 查看各自的类型。
//
// 验收：常规配置应包含 profile-api、development、true、9000，默认端口仍为 8080。
// 边界：把 serviceName 改为 ""、debug 改为 false，输出应忠实保留空字符串和 false。
// 建议实现位置：直接在 main() 内完成；本题不需要额外的 helper 函数。
// 思考：TS 的 const config = {} 合法，为什么 Go 的 const 不能保存这种配置对象？

// 练习 2：观察尚未加载的功能开关。
// 背景：后端配置未提供初值时，Go 会先给变量零值，不会产生 undefined。
// 要求：
//  1. 在一组 var (...) 声明中创建 maxRetries int、region string、enabled bool。
//  2. 不赋初值，先输出它们；用 %q 输出 region，让空字符串也看得见。
//  3. 再分别赋值为 3、"ap-shanghai"、true，并再次输出。
//
// 验收：初始输出为 0、""、false；赋值后为 3、"ap-shanghai"、true。
// 边界：将 maxRetries 再次设为 0，确认 0 也是合法值，不能自动理解成“没配置”。
// 建议实现位置：直接在 main() 内完成。
// 思考：仅凭 enabled == false，是否能区分“用户关闭”与“用户没填写”？
// 本题只需说明不能区分，后续章节再学习如何额外记录配置是否存在。

// 练习 3：避免临时配置意外覆盖外部配置。
// 背景：调试某段代码时，常会临时修改当前环境名。
// 要求：
//  1. 在函数里创建 environment := "production"。
//  2. 创建一个独立的 {...} 代码块，在块内用 := 声明同名变量 "preview"。
//  3. 分别输出块内、块外的 environment。
//  4. 再创建一个块，改用 = 将 environment 赋为 "staging"，观察块外的值。
//
// 验收：第一次块内为 preview，块外为 production；第二次块外变成 staging。
// 边界：将内层 := 的值改成 ""，外层仍应为 production。
// 建议实现位置：直接在 main() 内完成。
// 思考：把所有内层 := 换成 =，哪些输出会发生变化？为什么？

func main() {
}
