//go:build ignore

package main

// 练习 1:定义应用层 JSON 消息
// 要求:定义 Message{Type string, Text string}, 带 json 标签;客户端发 {"type":"chat","text":"hi"}
// 服务端读 JSON, chat 类型返回 {"type":"ack","text":"hi"};未知类型或空 text 返回 type=error
// 每条消息单独 WriteJSON/ReadJSON;保留 ReadLimit、deadline 和 Close
// 常规:chat/hi -> ack/hi;边界:未知 type、空文本 -> 可识别错误消息, 服务不 panic

// 练习 2:验证握手的 Origin 策略
// 要求:用 httptest.Server 和 Dialer 测试同源、异源、不带 Origin 三种情况
// 默认策略下:Origin 的主机等于 Host -> 101;https://other.example -> 403;不带 Origin -> 可建立
// 自行创建 websocket.Upgrader, 保留默认 CheckOrigin;升级成功后发送 welcome, 不要全部放行来源
// 成功握手后读 welcome, 再正常关闭;失败时检查 HTTP 响应状态并关闭 Body
// 说明为什么“无 Origin 可连”意味着这项检查不能替代用户鉴权

// 练习 3:把多个生产者的消息交给唯一写协程
// 要求:两个 goroutine 分别生成 2 条消息, 发送到容量为 4 的 channel;唯一 writer 写入同一连接
// 等生产者都结束后, 由协调者关闭 channel, writer 消费完再发关闭帧
// reader 持续读取 welcome 和 4 条回显;分别保存并按内容核对, 不依赖不同生产者的先后顺序
// 边界:生产者数量为 0 也能结束;读写出错通过取消信号让所有 goroutine 退出
// 每个阻塞发送/等待都要考虑取消, 关闭前确保 writer 不再发送;用 go run -race exercises.go 验证

func main() {
}
