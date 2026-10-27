//go:build ignore

package main

// 练习 1:检验后端分页参数
// 建议:func validatePagination(page, size int) error
// 规则:page>=1, size 为 1..100
// 在 main 中建立用例表, 逐项调用并比较错误是否符合预期;失败时输出用例名或 panic
// 常规:(1,20)、(3,100) -> nil;边界:(0,20)、(1,0)、(1,101)、(-1,10) -> error
// 故意改错一次业务边界, 确认测试能够失败, 然后恢复正确实现

// 练习 2:验证配置文件保存的副作用
// 建议:func saveFeatureFlag(path string, enabled bool) error
// 规则:输出恰好为 "enabled=true\n" 或 "enabled=false\n", 覆盖旧内容
// 在 main 中用 os.MkdirTemp 创建独立目录, defer os.RemoveAll 清理, 逐项核对文件内容
// 常规:先保存 false 再保存 true, 读回恰好为 "enabled=true\n", 没有旧尾部
// 边界:父目录不存在 -> error;普通已有文件被正确覆盖;不要用 chmod 制造不稳定的权限测试
// 先检查文件操作错误, 再比较实际内容与预期内容;不能只打印而不核对

// 练习 3:用并发测试验证接口请求计数器
// 建议:type RequestCounter struct { /* Mutex 和计数值 */ }
//
//	func (c *RequestCounter) Inc();func (c *RequestCounter) Value() int
//
// 在 main 中启动 20 个 goroutine, 各递增 100 次, WaitGroup 等待后核对是否为 2000
// 边界:新计数器为 0;多个读者与写者并行不会触发数据竞争;不要用 sleep 猜测完成时机
// 等待结束后在主 goroutine 比较实际值与预期值, 使用 go run -race exercises.go 检查竞争

func main() {
}
