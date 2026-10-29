/*
章节:28-RPC远程调用 - 客户端
在本章目录执行: go run ./client/lesson.go
先启动服务端, 再在另一个终端启动客户端.
*/

package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节: 启动客户端: 设置 RPC 服务地址
// ═══════════════════════════════════════════════════════
func main() {
	address := "127.0.0.1:8082"
	if err := runRPCClient(address); err != nil {
		fmt.Println("RPC 客户端失败:", err)
		os.Exit(1)
	}
}

// ═══════════════════════════════════════════════════════
// 第 2 节: 建立连接并调用: 同步调用、异步调用与远程错误
// ═══════════════════════════════════════════════════════
// Call 等待本次完成;Go 返回 *rpc.Call, 结果要等 Done 到达后读取
// 方法名必须写 PriceService.Quote;reply 传地址, 客户端才能填入解码结果
// 远端返回的错误文本不会保留本进程中的原始错误实例

func runRPCClient(address string) error {
	// rpc.Dial("tcp", address) 是简写;这里自行 DialTimeout 以限制连接建立时间
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return err
	}
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		conn.Close()
		return err
	}
	client := rpc.NewClient(conn)
	defer client.Close() // Close 同时关闭底层连接, 无需再由其他代码长期持有它
	// TS 常写 const reply = await api.quote(args)
	// net/rpc 把业务结果写进 &reply;Call 自身返回的是 error
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
		make(chan *rpc.Call, 1)) // Done 通道必须有缓冲;传 nil 时库会创建
	completed := <-call.Done
	if completed.Error != nil {
		return completed.Error
	}
	if asyncReply.TotalCents != 1000 {
		return fmt.Errorf("异步报价不符：%d", asyncReply.TotalCents)
	}
	fmt.Println("异步 RPC，总价（分）：", asyncReply.TotalCents) // 1000
	// Go 返回后可做其他工作, 但只有 Done 到来后才能安全读取 reply

	var invalidReply QuoteResponse
	err = client.Call("PriceService.Quote", QuoteRequest{Quantity: 0, UnitPriceCents: 500}, &invalidReply)
	var serverError rpc.ServerError
	if !errors.As(err, &serverError) {
		return fmt.Errorf("期望远程业务错误，实际为：%v", err)
	}
	fmt.Println("远程校验失败：", err)
	// 远程 error 以字符串传回来;不能用 errors.Is 匹配服务端的本地 sentinel 实例
	// 有业务错误时 reply 不发送, 客户端不能把旧 reply 当作本次成功结果
	return nil
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 请求与响应类型: 与服务端保持字段约定一致
// ═══════════════════════════════════════════════════════
// server 和 client 分开运行时仍共享相同的请求、响应类型和方法名约定
// 两个程序分别定义相同的导出字段, 供 gob 编码和解码使用
type QuoteRequest struct {
	Quantity       int
	UnitPriceCents int64
}

type QuoteResponse struct {
	TotalCents int64
}

// ═══════════════════════════════════════════════════════
// 第 4 节: 远程调用的失败边界
// ═══════════════════════════════════════════════════════
// 超时只说明调用方没有按时得到结果, 服务端是否已经执行还要按业务协议确认
// 与本地函数的区别:会断网、超时, 超时后服务端仍可能已执行
// 给 client.Go 的 Done 加 select 超时只让调用方停止等, 并不会取消远端执行
// net/rpc 没有内建的每次调用 context;关闭 client 会影响该连接上其他未完成调用
// 创建订单等非幂等操作不能在失败后盲目重试, 后续应学习业务幂等键与结果查询
