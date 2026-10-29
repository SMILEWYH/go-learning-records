/*
章节:26-TCP网络编程 - 客户端
*/

/*
服务端的小节顺序是:
1. 启动客户端：设置服务端地址
2. 建立连接并通信：通过 DialTimeout 连接，再发送消息、读取和校验回复
3. 消息读写细节
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节: 启动客户端: 设置要连接的服务端地址
// ═══════════════════════════════════════════════════════
func main() {
	address := "127.0.0.1:8081" // 服务端接口地址、端口

	if err := runTCPClient(address); err != nil {
		fmt.Println("TCP 客户端失败:", err)
	}
}

// ═══════════════════════════════════════════════════════
// 第 2 节: 建立 TCP 连接: DialTimeout 连接, 随后发送和接收消息
// ═══════════════════════════════════════════════════════
func runTCPClient(address string) error {
	conn, err := net.DialTimeout("tcp", address, ioTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Println("客户端本地地址：", conn.LocalAddr(), "服务端地址：", conn.RemoteAddr())
	reader := newMessageReader(conn)
	for _, message := range []string{"hello tcp", "你好, Go"} {
		// 	设置本轮 tcp 连接的读写截止时间
		if err := conn.SetDeadline(time.Now().Add(ioTimeout)); err != nil {
			return err
		}
		// 发送消息到客户端
		if err := writeAll(conn, []byte(message+"\n")); err != nil {
			return err
		}
		// 读取客户端返回给客户端的消息
		reply, err := readMessage(reader)
		if err != nil {
			return err
		}
		// 判断当前这条消息是否合格
		if reply != strings.ToUpper(message) {
			return fmt.Errorf("响应不符: %q", reply)
		}
		fmt.Printf("发送=%q, 收到=%q\n", message, reply)
	}
	return nil
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 消息读写细节: 按换行读取, 写完剩余字节
// ═══════════════════════════════════════════════════════
const messageLimit = 1024         // 本章一条消息最多 1024 字节, 包括末尾的 \n
const ioTimeout = 5 * time.Second // TCP 建立连接及每轮消息读写使用的超时时长

// newMessageReader 为一个连接创建一个持续使用的缓冲读取器, 为了使用 ReadSlice 方法
func newMessageReader(reader io.Reader) *bufio.Reader {
	return bufio.NewReaderSize(reader, messageLimit)
}

func readMessage(reader *bufio.Reader) (string, error) {
	// 按 \n 划分来读取缓冲区的消息, 这一步操作是会停在读取操作这里等待，不会继续执行后面的代码
	line, err := reader.ReadSlice('\n')
	if errors.Is(err, bufio.ErrBufferFull) {
		return "", errors.New("缓冲区已满, 但未找到换行符")
	}
	// 数据流结束, 缓冲区剩余数据中并没有找到换行符
	if errors.Is(err, io.EOF) && len(line) > 0 {
		return "", io.ErrUnexpectedEOF
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
}
