/*
章节：26-TCP网络编程
	本章掌握清单：
	1. TCP 连接与字节流
	2. 消息协议：读完整一条，写完剩余字节
	3. 处理一个连接：循环读取并回复
	4. 接受多个连接：Accept 循环与任务等待
	5. 客户端：连接、发送、读取、校验
	6. 自动演示：先监听，再连接，最后收尾
	7. 命令行模式与连接 API 复习
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./26-tcp
// 运行练习：go run ./26-tcp exam
// 阅读顺序：按第 1–7 节从协议读到客户端；main 只分发 exam 和讲解入口，默认演示会自动结束。
// TS 对照：net.Listen/Accept 类似 Node.js net.createServer；Dial 类似 net.createConnection。
// 双终端：go run ./26-tcp server；go run ./26-tcp client
// 可在 server/client 后传地址，默认 127.0.0.1:8081。
// 后续再学：长度前缀协议、TLS、连接数限制、心跳、重连与优雅停止全部连接。

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
// 第 1 节：TCP 连接与字节流
// ═══════════════════════════════════════════════════════
// 先建立连接，再双向传字节，最后关闭。应用还要自己决定怎样从字节中识别完整消息。
// 本例约定每条消息以换行结束，所以 hello 与后面的消息不会混在一起解释。

/*
1. TCP 是传输层协议，向应用提供可靠、有序、全双工的字节流。
   “可靠”不保证对方业务已经处理成功；超时/断线时仍可能不知道结果，业务需另行确认。
   IP 定位主机，端口定位该主机上的服务。127.0.0.1 只表示本机回环。

2. 常见三次握手（由操作系统 TCP 栈完成，应用不手写这些包）：
   客户端 -> 服务端：SYN，seq=x
   服务端 -> 客户端：SYN+ACK，seq=y，ack=x+1
   客户端 -> 服务端：ACK，ack=y+1
   双方同步各自的序号，并确认双向通信。Listen 等待连接，Dial 发起连接。

3. 常见四次挥手：A 发 FIN -> B 回 ACK -> B 发 FIN -> A 回 ACK。
   每个方向分别关闭，所以收到 FIN 后，另一方向仍可能发送剩余数据。
   ACK 和 FIN 有时合并；异常断开还可能用 RST，不能认为每次都有四个独立报文。
   通常主动关闭方进入 TIME_WAIT，以便处理重传并让旧报文消退。
   *net.TCPConn.CloseWrite() 可以只关闭发送方向；Close() 则释放整个连接。

4. TCP 不保留应用消息边界：一次 Write 不对应一次 Read。
   一条消息可被拆到多次 Read，多条消息也可能被一次 Read 读到，俗称拆包/粘包。
   必须约定分隔符、固定长度或长度前缀。本例用换行分隔，正文禁止换行。
   HTTP、RPC、WebSocket 都会在自己的层次处理消息格式。
*/

// ═══════════════════════════════════════════════════════
// 第 2 节：消息协议：读完整一条，写完剩余字节
// ═══════════════════════════════════════════════════════
// newMessageReader 为一个连接创建一个持续使用的缓冲读取器。
// readMessage 负责按换行切消息，writeAll 负责处理尚未写完的字节。
// 注意分层：net.Conn 提供字节传输，这两个函数负责应用的消息边界。

const messageLimit = 1024 // 本章一条消息最多 1024 字节，包括末尾的 \n。
const ioTimeout = 5 * time.Second

func newMessageReader(reader io.Reader) *bufio.Reader {
	return bufio.NewReaderSize(reader, messageLimit)
}

func readMessage(reader *bufio.Reader) (string, error) {
	// ReadSlice 会跨多次底层 Read 查找换行，并把多出来的字节留给下一条消息。
	// 本例用固定缓冲区限制长度；不能每次读消息都重新创建 Reader，否则可能丢失预读数据。
	line, err := reader.ReadSlice('\n')
	if errors.Is(err, bufio.ErrBufferFull) {
		return "", errors.New("消息超过 1024 字节或缺少换行")
	}
	if errors.Is(err, io.EOF) && len(line) > 0 {
		return "", io.ErrUnexpectedEOF // 本协议要求换行，半条消息不能当作成功。
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
	// 只重发尚未写入的部分。出现 error 后不能盲目从头重发整条业务消息。
}

// ═══════════════════════════════════════════════════════
// 第 3 节：处理一个连接：循环读取并回复
// ═══════════════════════════════════════════════════════
// handleTCP 的参数是已经建立的连接，函数退出时关闭它。
// 每轮先设置期限，再读一条消息，再返回大写版本；EOF 表示对方已结束发送。

func handleTCP(conn net.Conn) error {
	defer conn.Close() // 每个连接由处理它的函数关闭，别把 defer 堆在无限 Accept 循环里。
	reader := newMessageReader(conn)
	for { // 内层循环：在同一个连接上处理多条消息。
		if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
			return err
		}
		message, err := readMessage(reader)
		if errors.Is(err, io.EOF) {
			return nil // 对端正常关闭发送方向，且没有半条消息。
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
// 第 4 节：接受多个连接：Accept 循环与任务等待
// ═══════════════════════════════════════════════════════
// listener 负责接受新连接，conn 代表某个已连接客户端，两者是不同资源。
// 普通服务给每个连接启动处理任务；自动演示的 once=true 只接待一个客户端。

func serveTCP(listener net.Listener, once bool) error {
	var workers sync.WaitGroup
	defer workers.Wait()
	for { // 外层循环：接受不同客户端；Accept 没有连接时阻塞，不是忙轮询。
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		if once {
			return handleTCP(conn) // 自动演示只接待一个客户端，完成后退出。
		}
		workers.Add(1)
		go func() {
			defer workers.Done()
			if err := handleTCP(conn); err != nil {
				fmt.Fprintln(os.Stderr, "TCP 连接结束：", err)
			}
		}()
	}
	// Close(listener) 只停止接收新连接，已有连接仍需各自关闭。
}

// ═══════════════════════════════════════════════════════
// 第 5 节：客户端：连接、发送、读取、校验
// ═══════════════════════════════════════════════════════
// DialTimeout 先连接服务，随后每条消息依次发送、读取响应并检查内容。
// 发送中文时长度仍按 UTF-8 字节计算；只有收到预期大写结果才算本次演示成功。

func runTCPClient(address string) error {
	conn, err := net.DialTimeout("tcp", address, ioTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println("客户端本地地址：", conn.LocalAddr(), "服务端地址：", conn.RemoteAddr())
	reader := newMessageReader(conn)
	for _, message := range []string{"hello tcp", "你好，Go"} {
		if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
			return err
		}
		if err := writeAll(conn, []byte(message+"\n")); err != nil {
			return err
		}
		reply, err := readMessage(reader)
		if err != nil {
			return err
		}
		if reply != strings.ToUpper(message) {
			return fmt.Errorf("响应不符：%q", reply)
		}
		fmt.Printf("发送=%q，收到=%q\n", message, reply)
	}
	return nil
}

// ═══════════════════════════════════════════════════════
// 第 6 节：自动演示：先监听，再连接，最后收尾
// ═══════════════════════════════════════════════════════
// 端口 0 交给系统挑选空闲端口，Addr 返回实际监听地址。
// 客户端结束后关闭 listener，并等待服务端结果，确保演示能结束。

func demoTCP() error {
	listener, err := net.Listen("tcp", "127.0.0.1:0") // 端口 0 让系统分配空闲端口。
	if err != nil {
		return err
	}
	defer listener.Close()
	done := make(chan error, 1)
	go func() { done <- serveTCP(listener, true) }()
	clientErr := runTCPClient(listener.Addr().String())
	// 先 Listen，再启动客户端；无需 sleep 猜服务是否准备好。
	listener.Close() // 客户端建连失败时也能唤醒阻塞的 Accept。
	return errors.Join(clientErr, <-done)
}

// ═══════════════════════════════════════════════════════
// 第 7 节：命令行模式与连接 API 复习
// ═══════════════════════════════════════════════════════
// 无参数调用 demoTCP，server 持续接待连接，client 只执行客户端通信。
// 双终端练习时，两边必须使用相同的 host:port。

func runTCP(args []string) error {
	if len(args) == 0 {
		return demoTCP()
	}
	address := "127.0.0.1:8081"
	if len(args) > 1 {
		address = args[1]
	}
	switch args[0] {
	case "server":
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}
		defer listener.Close()
		fmt.Println("TCP 服务监听：", listener.Addr(), "；Ctrl+C 结束")
		return serveTCP(listener, false)
	case "client":
		return runTCPClient(address)
	default:
		return errors.New("用法：不带参数，或 server/client [地址]，或 exam")
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := runTCP(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "TCP 演示失败：", err)
		os.Exit(1)
	}
}

// 原笔记还可写为 ResolveTCPAddr -> ListenTCP，得到 *net.TCPListener；每一步都要检查 err。
// net.Listen("tcp", address) 是更简洁的入口；AcceptTCP/DialTCP 可取得 *net.TCPConn。
// Read([]byte) (n, err)：先处理 n>0 的字节，再判断 err，二者可能同时存在。
// io.EOF 是错误值，不是 io.EOF() 函数；用 errors.Is(err, io.EOF) 判断。
// Write([]byte) (n, err)：写入字节数不等于对方业务处理成功；不必无条件套无限循环。
// Close() 释放连接；LocalAddr() 是自己这一端，RemoteAddr() 才是另一端。
// DialTimeout 只约束建立连接；SetReadDeadline/SetWriteDeadline 约束已建立连接的 I/O。
// SetDeadline 是绝对截止时间，持续会话需要刷新；time.Time{} 可清除 deadline。
// 服务端在同一个 goroutine 里无限 Accept 后再写客户端代码，客户端部分就永远执行不到。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 运行默认示例，核对英文与中文两条回显；再用双终端观察一条连接承载多条消息。
// 阅读 readMessage 的 EOF 分支，说明“完整消息后结束”与“缺少换行的半条消息”有什么区别。
