//go:build ignore

package main

// 练习 1:读取部署配置
// 建议:func parsePort(raw string) (int, error)
// 要求:空串使用 8080;其他值必须是 1..65535 的十进制整数;不访问真实环境变量
// 常规:""->8080;"9090"->9090;"65535"->65535
// 边界:"0"、"65536"、"abc"、"-1" -> error
// 将 os.Getenv 留在调用处, 使 helper 可用表驱动测试验证

// 练习 2:给构建产物加入版本信息
// 要求:自行定义包级 version、buildTime 两个 string 变量, 通过 -ldflags -X 注入固定的版本和时间
// 在 version 模式输出两者;不注入时应有可读默认值
// 常规:v0.2.0、2026-09-27T12:00:00Z -> 输出一致;默认构建 -> dev 与自己设置的默认时间
// 在本章目录用 go build -o ./bin/exercises exercises.go 构建, 实际运行本机产物验证;说明 Linux 产物为什么不能在 macOS 原生执行

// 练习 3:验证优雅退出会等待正在处理的请求
// 要求:在自己的测试 server 中用 channel 通知“请求已进入”, 并阻塞 handler 等待释放信号
// 发起 Shutdown 后, 确认它尚未完成;释放 handler 后, 请求与 Shutdown 都应成功结束
// 边界:换一个很短的 shutdown context 且不释放请求, Shutdown 应返回 deadline exceeded, 随后 Close
// 所有等待都应有超时兜底;不要用 sleep 猜 handler 是否已开始;测试结束必须清理连接与 goroutine

func main() {
}
