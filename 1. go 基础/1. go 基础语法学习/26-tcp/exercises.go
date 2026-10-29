//go:build ignore

package main

// 练习 1:为 TCP 聊天准备完整消息
// 建议:func encodeLine(message string) ([]byte, error)
// 要求:拒绝正文含 \n 或 \r;UTF-8 字节数加末尾 \n 不超过 1024;成功追加换行
// 常规:"hello" -> []byte("hello\n");"你好" -> 7 字节（包含换行）
// 边界:空串 -> []byte("\n");1023 个 ASCII 字符成功;1024 个失败;正文含换行失败

// 练习 2:用同一个连接交换多条消息
// 要求:在 127.0.0.1:0 启动本地服务, 客户端发送 a、b、c, 依次收到 A、B、C
// 必须复用同一个 conn 和同一个 bufio.Reader;服务端接待完这个客户端后退出
// 边界:把 "a\nb\n" 放进一次 Write, 仍应读取到两条响应;空消息收到空消息
// 使用 deadline 限制等待, 并通过 channel/WaitGroup 等待服务端结束;不要用 sleep 同步

// 练习 3:区分正常结束、半条消息和超时
// 要求:自行实现按换行读取消息的 readMessage, 在 main 中核对三种结果:完整 "ok\n", 空流 EOF, "unfinished" 的 UnexpectedEOF
// 再用 net.Pipe 建立内存连接;对端不发送, 设置短读取 deadline, 检查 os.IsTimeout(err)
// net.Pipe 不是 TCP, 只用它模拟阻塞读;连接两端都需关闭, 不用“超时后继续无限 Read”
// 用 errors.Is 判断 EOF 类错误, 不比较 err.Error() 文本

func main() {
}
