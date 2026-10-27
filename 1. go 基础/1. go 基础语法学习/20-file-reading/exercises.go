//go:build ignore

package main

// 题目约定:使用 os.MkdirTemp 创建自己的目录, defer os.RemoveAll 清理, 只操作其中的练习文件

// 练习 1:加载一个服务的 JSON 配置
// 建议:type AppConfig struct { Name string `json:"name"`; Port int `json:"port"` }
//
//	func readConfig(path string) (AppConfig, error)
//
// 要求:os.ReadFile + json.Unmarshal;校验 Name 非空、Port 为 1..65535;错误保留上下文
// 常规:{"name":"api","port":8080} -> 对应结构体, nil
// 边界:缺失文件、非法 JSON、port=0、空 name 都返回 error, 不能 panic

// 练习 2:统计结构化日志中以 "ERROR " 开头的行
// 建议:func countErrorLines(path string) (int, error)
// 要求:用 Scanner, 设置 64 KiB 缓冲上限, 并检查 Scanner.Err 和关闭错误
// 常规:"INFO start\nERROR db\nERROR cache" -> 2, nil;最后一行无换行也要处理
// 边界:空文件 -> 0, nil;"INFO ERROR text" 不计数;一行 70 KiB -> error
// 出错时统一返回 0, error, 避免调用者把不完整统计当成最终结果

// 练习 3:给 HTTP 请求体读取设置大小限制（本题不需要启动服务）
// 建议:func readLimited(reader io.Reader, maxBytes int64) ([]byte, error)
// 要求:maxBytes 为 1..1 MiB;最多尝试读取 maxBytes+1 字节, 超限返回 error
// 可用 io.LimitReader + io.ReadAll;不要先无上限读取完整内容再判断大小
// 常规:strings.NewReader("abc"), 3 -> "abc", nil
// 边界:"abcd", 3 -> error;空内容, 3 -> 长度 0;非法上限 -> error
// 扩展:实现一个最后一次 Read 同时返回数据和 io.EOF 的 reader, 确认不会丢数据

func main() {
}
