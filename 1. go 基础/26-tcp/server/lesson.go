/*
章节:26-TCP网络编程 - 服务端
*/

/*
服务端的小节顺序是:
1. TCP 连接与字节流的介绍
2. 启动服务端: 通过 net.Listen 监听地址
3. 接受客户端连接: 通过 Accept 获取连接
4. 处理连接中的消息: 循环读取、转换大写、回复
5. 消息读写细节: 解释 readMessage、writeAll 等函数
6. 连接 API 复习
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节:TCP 连接与字节流
// ═══════════════════════════════════════════════════════
// TCP 连接的整个生命周期: 先建立连接 -> 再双向传字节 -> 最后关闭连接
// TCP 传递的是字节流, 不保留消息边界, 双方需要自行约定如何识别每条完整消息
// 本章节约定以换行符 \n 作为每条消息的结束标记, 消息正文中不能包含换行符
/*
1. TCP 是传输层协议, 向应用提供可靠、有序、全双工的字节流

2. 常见三次握手 (建立连接时发生)
   客户端 -> 服务端: 服务端收到后, 确认可以收到客户端的消息
   服务端 -> 客户端: 客户端收到回复后, 确认自己的消息可以到达服务端, 也能收到服务端的消息
   客户端 -> 服务端: 服务端收到确认后, 确认自己的回复可以达到客户端

3. 常见四次挥手:  (结束连接时发生)
   客户端 -> 服务端: 客户端通知服务端, 自己的数据已发送完毕, 之后不会再发送, 但仍可以接收服务端的数据
   服务端 -> 客户端: 服务端确认收到客户端的结束通知, 此时服务端仍可能继续发送剩余数据
   服务端 -> 客户端: 服务端的数据也发送完毕, 通知客户端, 之后不会再发送
   客户端 -> 服务端: 客户端确认收到服务端的结束通知, 让服务端知道自己的结束通知已到达客户端

4. TCP 不保留应用消息边界: 一次 Write 不对应一次 Read
   一条消息可被拆到多次 Read, 多条消息也可能被一次 Read 读到, 俗称拆包/粘包
   所以必须约定分隔符、固定长度或长度前缀
   HTTP、RPC、WebSocket 都会在自己的层次处理消息格式
*/

// ═══════════════════════════════════════════════════════
// 第 2 节: 启动服务端: Listen 监听 TCP 地址
// ═══════════════════════════════════════════════════════
// 从这里开始阅读服务端流程: 先 Listen, 再 Accept, 最后处理连接中的消息
// 双终端练习时, 两边必须使用相同的 host:port
func main() {
	address := "127.0.0.1:8081"

	if err := runTCPServer(address); err != nil {
		fmt.Println("TCP 客户端失败:", err)
	}
}

func runTCPServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("TCP 服务监听:", listener.Addr(), "; Ctrl+C 结束")
	return serveTCP(listener)
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 接受客户端连接: Accept 循环与任务等待
// ═══════════════════════════════════════════════════════
func serveTCP(listener net.Listener) error {
	var workers sync.WaitGroup
	defer workers.Wait()
	for {
		// 等待并接受一个客户端连接, 没有连接时阻塞代码往下执行
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		workers.Add(1)
		go func() {
			defer workers.Done()
			if err := handleTCP(conn); err != nil {
				fmt.Fprintln(os.Stderr, "TCP 连接结束：", err)
			}
		}()
	}
	// Close(listener) 只停止接收新连接, 已有连接仍需各自关闭
}

// ═══════════════════════════════════════════════════════
// 第 4 节: 处理一个连接: 循环读取并回复
// ═══════════════════════════════════════════════════════
func handleTCP(conn net.Conn) error {
	defer conn.Close() // 每个连接由处理它的函数关闭, 别把 defer 堆在无限 Accept 循环里
	reader := newMessageReader(conn)
	for { // 内层循环:在同一个连接上处理多条消息
		if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
			return err
		}
		message, err := readMessage(reader)
		if errors.Is(err, io.EOF) {
			return nil // 对端正常关闭发送方向, 且没有半条消息
		}
		if err != nil {
			return err
		}
		if err := writeAll(conn, []byte(strings.ToUpper(message)+"\n")); err != nil {
			return err
		}
	}
}

// ═══════════════════════════════════════════════════════
// 第 5 节: 消息读写细节: 按换行读取, 写完剩余字节
// ═══════════════════════════════════════════════════════
// newMessageReader 为一个连接创建一个持续使用的缓冲读取器
// readMessage 负责按换行切消息
// writeAll 负责处理尚未写完的字节
// 注意分层: net.Conn 提供字节传输, readMessage 识别消息边界, writeAll 处理字节写入

const messageLimit = 1024 // 本章一条消息最多 1024 字节, 包括末尾的 \n
const ioTimeout = 5 * time.Second

func newMessageReader(reader io.Reader) *bufio.Reader {
	return bufio.NewReaderSize(reader, messageLimit)
}

func readMessage(reader *bufio.Reader) (string, error) {
	// ReadSlice 会跨多次底层 Read 查找换行, 并把多出来的字节留给下一条消息
	// 本例用固定缓冲区限制长度;不能每次读消息都重新创建 Reader, 否则可能丢失预读数据
	line, err := reader.ReadSlice('\n')
	if errors.Is(err, bufio.ErrBufferFull) {
		return "", errors.New("消息超过 1024 字节或缺少换行")
	}
	if errors.Is(err, io.EOF) && len(line) > 0 {
		return "", io.ErrUnexpectedEOF // 本协议要求换行, 半条消息不能当作成功
	}
	if err != nil {
		return "", err
	}
	return string(line[:len(line)-1]), nil
}

func writeAll(writer io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := writer.Write(data)
		data = data[n:]
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
	// 只重发尚未写入的部分出现 error 后不能盲目从头重发整条业务消息
}

// ═══════════════════════════════════════════════════════
// 第 6 节: 连接 API 复习
// ═══════════════════════════════════════════════════════
// ResolveTCPAddr -> ListenTCP 可得到 *net.TCPListener;每一步都要检查 err
// net.Listen("tcp", address) 是更简洁的入口;AcceptTCP/DialTCP 可取得 *net.TCPConn
// Read([]byte) (n, err):先处理 n>0 的字节, 再判断 err, 二者可能同时存在
// io.EOF 是错误值, 不是 io.EOF() 函数;用 errors.Is(err, io.EOF) 判断
// Write([]byte) (n, err):写入字节数不等于对方业务处理成功;不必无条件套无限循环
// Close() 释放连接;LocalAddr() 是自己这一端, RemoteAddr() 才是另一端
// DialTimeout 只约束建立连接;SetReadDeadline/SetWriteDeadline 约束已建立连接的 I/O
// SetDeadline 是绝对截止时间, 持续会话需要刷新;time.Time{} 可清除 deadline
// 服务端在同一个 goroutine 里无限 Accept 后再写客户端代码, 客户端部分就永远执行不到
