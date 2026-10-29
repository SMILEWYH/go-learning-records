/*
章节:28-RPC远程调用 - 服务端
在本章目录执行: go run ./server/lesson.go
先启动服务端, 再在另一个终端启动客户端.
*/

package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节: RPC: 把本地调用拆成远程通信
// ═══════════════════════════════════════════════════════
// 调用形式简洁, 但参数和结果都要经过编码与网络
// 先看调用链与注册条件, 再对应 Quote 的请求结构体、结果指针和 error

/*
RPC（Remote Procedure Call）:像调用本地方法一样请求远程服务
调用链:客户端参数 -> 编码 -> 传输 -> 服务端解码 -> 执行方法 -> 编码结果 -> 客户端解码
HTTP 是协议, RPC 是调用方式;RPC 可使用 HTTP 或其他传输, 两者不是互斥关系
本例用标准库 net/rpc + TCP + gob 编码, 两端均使用 Go, 无需另装框架
net/rpc 已冻结、不再增加新功能, 适合学习原理;它不等于 gRPC, 也不是浏览器 fetch API

用 Register 注册的方法需要满足:
  - 服务类型和方法导出, 例如 PriceService.Quote
  - 两个业务参数的类型是导出类型或内置类型;第二个参数必须是指针
  - 返回值只有 error;参数和结果还必须能被所选编码器序列化
  - 常见签名:func (service *Service) Method(args Request, reply *Response) error
  - gob 只编码结构体的导出字段, 字段小写可能导致数据丢失或编码失败
*/

// ═══════════════════════════════════════════════════════
// 第 2 节: 启动 RPC 服务: 注册服务并监听 TCP 地址
// ═══════════════════════════════════════════════════════
func main() {
	address := "127.0.0.1:8082"
	if err := runRPCServer(address); err != nil {
		fmt.Println("RPC 服务端失败:", err)
		os.Exit(1)
	}
}

func runRPCServer(address string) error {
	// 先创建 RPC 服务并注册方法, 再监听 TCP 地址
	server, err := newRPCServer()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("RPC 服务监听:", listener.Addr(), "; Ctrl+C 结束")
	return serveRPC(server, listener)
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 创建 RPC 服务并注册可调用的方法
// ═══════════════════════════════════════════════════════
// Register 将 PriceService 的可用方法登记到服务表, 客户端通过字符串方法名查找
// ServeConn 接管单条连接的编解码和分发, 外层 Accept 负责接待不同连接

func newRPCServer() (*rpc.Server, error) {
	server := rpc.NewServer() // 独立注册表;不污染默认全局 rpc.Server
	if err := server.Register(&PriceService{}); err != nil {
		return nil, err
	}
	return server, nil
}

// ═══════════════════════════════════════════════════════
// 第 4 节: 接受客户端连接: 交给 RPC 处理请求和响应
// ═══════════════════════════════════════════════════════
func serveRPC(server *rpc.Server, listener net.Listener) error {
	var workers sync.WaitGroup
	defer workers.Wait()
	for {
		conn, err := listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		// 学习示例限制每条连接的总 I/O 时间为 5 秒, 避免坏客户端永久占着连接
		// 长期服务需要另外设计空闲超时及请求级限制, 不应照搬这个固定会话时长
		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			conn.Close()
			return err
		}
		workers.Add(1)
		go func() {
			defer workers.Done()
			server.ServeConn(conn) // 内部读请求、调方法、写响应;结束时关闭连接
		}()
	}
}

// ═══════════════════════════════════════════════════════
// 第 5 节: 服务方法: 请求、reply 指针与业务校验
// ═══════════════════════════════════════════════════════
// QuoteRequest 是输入, QuoteResponse 保存结果;reply 指针让方法写回计算值
// 数量和单价先校验范围, 再相乘;返回 nil 才表示服务端业务成功

type QuoteRequest struct {
	Quantity       int
	UnitPriceCents int64
}

type QuoteResponse struct {
	TotalCents int64
}

type PriceService struct{}

func (service *PriceService) Quote(args QuoteRequest, reply *QuoteResponse) error {
	if args.Quantity < 1 || args.Quantity > 1000 {
		return errors.New("数量必须在 1..1000")
	}
	if args.UnitPriceCents < 0 || args.UnitPriceCents > 100_000_000 {
		return errors.New("单价必须在 0..100000000 分")
	}
	reply.TotalCents = int64(args.Quantity) * args.UnitPriceCents // 整数分避免常见小数金额误差
	return nil
	// 限制乘数范围, 也避免 int64 乘法溢出业务可预期失败通过 error 返回
	// 方法可能并发执行;若在 service 中加 map/计数器, 需像第 17 章那样保护共享状态
}
