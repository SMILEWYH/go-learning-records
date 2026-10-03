/*
章节：28-RPC远程调用
	本章掌握清单：
	1. RPC：把本地调用拆成远程通信
	2. 服务方法：请求、reply 指针与业务校验
	3. 服务注册与连接处理
	4. 客户端同步调用、异步调用与远程错误
	5. 自动演示与资源归属
	6. 运行模式与远程调用的失败边界
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./28-rpc
// 运行练习：go run ./28-rpc exam
// 阅读顺序：按概念、类型、服务、客户端的顺序阅读；默认 main -> runRPC -> demoRPC 完成全部通信后退出。
// TS 对照：调用形式像 await api.quote(input)，但这里把方法名和参数编码后交给远端执行。
// 双终端：go run ./28-rpc server；go run ./28-rpc client
// 可在 server/client 后传地址，默认 127.0.0.1:8082。
// 后续再学：gRPC/Protobuf、请求级取消、服务发现、认证、重试和幂等。

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
// 第 1 节：RPC：把本地调用拆成远程通信
// ═══════════════════════════════════════════════════════
// 调用形式简洁，但参数和结果都要经过编码与网络。
// 先看调用链与注册条件，再对应 Quote 的请求结构体、结果指针和 error。

/*
RPC（Remote Procedure Call）：像调用本地方法一样请求远程服务。
调用链：客户端参数 -> 编码 -> 传输 -> 服务端解码 -> 执行方法 -> 编码结果 -> 客户端解码。
HTTP 是协议，RPC 是调用方式；RPC 可使用 HTTP 或其他传输，两者不是互斥关系。
本例用标准库 net/rpc + TCP + gob 编码，两端均使用 Go，无需另装框架。
net/rpc 已冻结、不再增加新功能，适合学习原理；它不等于 gRPC，也不是浏览器 fetch API。

用 Register 注册的方法需要满足：
  - 服务类型和方法导出，例如 PriceService.Quote。
  - 两个业务参数的类型是导出类型或内置类型；第二个参数必须是指针。
  - 返回值只有 error；参数和结果还必须能被所选编码器序列化。
  - 常见签名：func (service *Service) Method(args Request, reply *Response) error。
  - gob 只编码结构体的导出字段，字段小写可能导致数据丢失或编码失败。
*/

// ═══════════════════════════════════════════════════════
// 第 2 节：服务方法：请求、reply 指针与业务校验
// ═══════════════════════════════════════════════════════
// QuoteRequest 是输入，QuoteResponse 保存结果；reply 指针让方法写回计算值。
// 数量和单价先校验范围，再相乘；返回 nil 才表示服务端业务成功。

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
	reply.TotalCents = int64(args.Quantity) * args.UnitPriceCents // 整数分避免常见小数金额误差。
	return nil
	// 限制乘数范围，也避免 int64 乘法溢出。业务可预期失败通过 error 返回。
	// 方法可能并发执行；若在 service 中加 map/计数器，需像第 17 章那样保护共享状态。
}

// ═══════════════════════════════════════════════════════
// 第 3 节：服务注册与连接处理
// ═══════════════════════════════════════════════════════
// Register 将 PriceService 的可用方法登记到服务表，客户端通过字符串方法名查找。
// ServeConn 接管单条连接的编解码和分发，外层 Accept 负责接待不同连接。

func newRPCServer() (*rpc.Server, error) {
	server := rpc.NewServer() // 独立注册表；不污染默认全局 rpc.Server。
	if err := server.Register(&PriceService{}); err != nil {
		return nil, err
	}
	return server, nil
}

func serveRPC(server *rpc.Server, listener net.Listener, once bool) error {
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
		// 学习示例限制每条连接的总 I/O 时间为 5 秒，避免坏客户端永久占着连接。
		// 长期服务需要另外设计空闲超时及请求级限制，不应照搬这个固定会话时长。
		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			conn.Close()
			return err
		}
		if once {
			server.ServeConn(conn) // 内部读请求、调方法、写响应；结束时关闭连接。
			return nil
		}
		workers.Add(1)
		go func() {
			defer workers.Done()
			server.ServeConn(conn)
		}()
	}
}

// ═══════════════════════════════════════════════════════
// 第 4 节：客户端同步调用、异步调用与远程错误
// ═══════════════════════════════════════════════════════
// Call 等待本次完成；Go 返回 *rpc.Call，结果要等 Done 到达后读取。
// 方法名必须写 PriceService.Quote；reply 传地址，客户端才能填入解码结果。
// 远端返回的错误文本不会保留本进程中的原始错误实例。

func runRPCClient(address string) error {
	// rpc.Dial("tcp", address) 是简写；这里自行 DialTimeout 以限制连接建立时间。
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return err
	}
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		conn.Close()
		return err
	}
	client := rpc.NewClient(conn)
	defer client.Close() // Close 同时关闭底层连接，无需再由其他代码长期持有它。
	var reply QuoteResponse
	if err := client.Call("PriceService.Quote", QuoteRequest{Quantity: 3, UnitPriceCents: 1990}, &reply); err != nil {
		return err
	}
	if reply.TotalCents != 5970 {
		return fmt.Errorf("报价不符：%d", reply.TotalCents)
	}
	fmt.Println("同步 RPC，总价（分）：", reply.TotalCents) // 5970

	var asyncReply QuoteResponse
	call := client.Go("PriceService.Quote", QuoteRequest{Quantity: 2, UnitPriceCents: 500}, &asyncReply,
		make(chan *rpc.Call, 1)) // Done 通道必须有缓冲；传 nil 时库会创建。
	completed := <-call.Done
	if completed.Error != nil {
		return completed.Error
	}
	if asyncReply.TotalCents != 1000 {
		return fmt.Errorf("异步报价不符：%d", asyncReply.TotalCents)
	}
	fmt.Println("异步 RPC，总价（分）：", asyncReply.TotalCents) // 1000
	// Go 返回后可做其他工作，但只有 Done 到来后才能安全读取 reply。

	var invalidReply QuoteResponse
	err = client.Call("PriceService.Quote", QuoteRequest{Quantity: 0, UnitPriceCents: 500}, &invalidReply)
	var serverError rpc.ServerError
	if !errors.As(err, &serverError) {
		return fmt.Errorf("期望远程业务错误，实际为：%v", err)
	}
	fmt.Println("远程校验失败：", err)
	// 远程 error 以字符串传回来；不能用 errors.Is 匹配服务端的本地 sentinel 实例。
	// 有业务错误时 reply 不发送，客户端不能把旧 reply 当作本次成功结果。
	return nil
}

// ═══════════════════════════════════════════════════════
// 第 5 节：自动演示与资源归属
// ═══════════════════════════════════════════════════════
// 先注册服务，再监听空闲端口，然后启动处理并连接客户端。
// 客户端 Close 关闭底层连接，ServeConn 才能结束，主流程继续等待服务端完成。

func demoRPC() error {
	server, err := newRPCServer()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer listener.Close()
	done := make(chan error, 1)
	go func() { done <- serveRPC(server, listener, true) }()
	clientErr := runRPCClient(listener.Addr().String())
	listener.Close()
	return errors.Join(clientErr, <-done)
}

// ═══════════════════════════════════════════════════════
// 第 6 节：运行模式与远程调用的失败边界
// ═══════════════════════════════════════════════════════
// server 和 client 分开运行时仍共享相同的请求、响应类型和方法名约定。
// 超时只说明调用方没有按时得到结果，服务端是否已经执行还要按业务协议确认。

func runRPC(args []string) error {
	if len(args) == 0 {
		return demoRPC()
	}
	address := "127.0.0.1:8082"
	if len(args) > 1 {
		address = args[1]
	}
	switch args[0] {
	case "server":
		server, err := newRPCServer()
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}
		defer listener.Close()
		fmt.Println("RPC 服务监听：", listener.Addr(), "；Ctrl+C 结束")
		return serveRPC(server, listener, false)
	case "client":
		return runRPCClient(address)
	default:
		return errors.New("用法：不带参数，或 server/client [地址]，或 exam")
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := runRPC(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "RPC 演示失败：", err)
		os.Exit(1)
	}
}

// 与本地函数的区别：会断网、超时，超时后服务端仍可能已执行。
// 给 client.Go 的 Done 加 select 超时只让调用方停止等，并不会取消远端执行。
// net/rpc 没有内建的每次调用 context；关闭 client 会影响该连接上其他未完成调用。
// 创建订单等非幂等操作不能在失败后盲目重试，后续应学习业务幂等键与结果查询。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 核算同步例子的 3×1990=5970 和异步例子的 2×500=1000。
// 把数量改成 0，跟踪错误如何从 Quote 返回客户端，解释为什么不能再使用 reply 当成功结果。
