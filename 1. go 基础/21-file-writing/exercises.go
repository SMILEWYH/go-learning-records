//go:build ignore

package main

// 练习 1:保存开发环境配置
// 建议:func saveSettings(path string, values map[string]string) error
// 要求:json.MarshalIndent 编码后用 os.WriteFile 写入, 末尾加换行, 新文件权限 0o600
// 常规:写入 {"theme":"dark"} 后能读回合法 JSON
// 边界:较短内容覆盖较长内容不留旧尾巴;空 map 和 nil 输入都保存为 {}, 不是 null
// 目标父目录缺失时报错;注意先处理 nil map, 再进行 JSON 编码
// 不要求 map 字段的输出顺序;读取后通过 json.Unmarshal 核对内容

// 练习 2:追加审计日志, 每个事件单独一行 JSON（JSON Lines）
// 建议:type AuditEvent struct { Actor string `json:"actor"`; Action string `json:"action"` }
//
//	func appendAudit(path string, event AuditEvent) error
//
// 要求:使用 O_CREATE|O_WRONLY|O_APPEND 和 json.Encoder.Encode;字段均不能为空
// 常规:顺序追加两个事件, 得到两行独立合法 JSON, 之前的内容仍在
// 边界:Actor 内含换行会被 JSON 转义, 不增加物理行数;空 Actor/Action 返回 error 且不写入
// 本题只需顺序调用;多进程共同写文件需要另外设计并发协议

// 练习 3:只创建一次的本地任务标记
// 建议:func createJobMarker(path, jobID string) error
// 要求:O_CREATE|O_EXCL|O_WRONLY, 权限 0o600;写入 jobID 和换行;jobID 为空先返回 error
// 常规:首次创建成功, 内容为 "job-1\n"
// 边界:同一路径第二次创建返回能被 errors.Is(err, os.ErrExist) 识别的错误, 旧内容不变
// 创建失败不能随后继续写;若写或关闭失败, 应把错误返回给调用者

func main() {
}
