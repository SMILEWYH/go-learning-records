/*
章节：29-WebSocket双向通信
	本章掌握清单：
	1. HTTP 升级与 WebSocket 消息
	2. 服务端：握手、主动推送与读取循环
	3. 路由与浏览器入口
	4. 客户端：欢迎消息、回显与关闭握手
	5. 自动演示：等待升级后的处理结束
	6. 双终端运行与单读单写规则
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./29-websocket
// 运行练习：go run ./29-websocket exam
// 阅读顺序：先读协议与 handleWebSocket，再读客户端；默认 main -> runWebSocket -> demoWebSocket 会自动收尾。
// TS 对照：浏览器有 new WebSocket(url)，Go 标准库不提供完整的 WebSocket 消息 API。
// 本章使用 github.com/gorilla/websocket v1.5.3，首次运行需下载这个固定版本。
// 双终端：go run ./29-websocket server；go run ./29-websocket client
// 可在 server/client 后传地址，默认 127.0.0.1:8083。
// 后续再学：聊天室广播、心跳与重连、慢客户端背压、鉴权和管理升级后的连接。

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ═══════════════════════════════════════════════════════
// 第 1 节：HTTP 升级与 WebSocket 消息
// ═══════════════════════════════════════════════════════
// 先完成 HTTP 握手，连接升级后使用 WebSocket 消息协议。
// 服务端可以主动发 welcome，客户端不必先发送一条业务消息才能收到它。

/*
名字是 WebSocket，不是 WebSocker。
本例过程：HTTP/1.1 GET /ws + Upgrade 请求 -> 101 Switching Protocols -> WebSocket 帧。
握手后两端都可以主动发送消息，不需要每次推送都等客户端发新的 HTTP 请求。
ws:// 是明文，wss:// 是经 TLS 保护的连接；部署到 HTTPS 页面通常配套使用 wss://。
WebSocket 有消息边界，一条消息也可拆成多个帧；库的 ReadMessage 帮我们组装完整消息。
底层 TCP 的一次 Read 仍不对应一条 WebSocket 消息。不能把 TCP conn.Read 当作 ReadMessage。

消息有文本、二进制两种常见类型；Ping/Pong 用于探活，Close 用于关闭握手。
客户端帧需要掩码等协议细节由库处理；不要把普通 HTTP 长连接当作 WebSocket。
*/

// ═══════════════════════════════════════════════════════
// 第 2 节：服务端：握手、主动推送与读取循环
// ═══════════════════════════════════════════════════════
// Upgrade 成功才获得 conn，之后用 ReadMessage/WriteMessage 收发完整消息。
// 先限制消息大小并推送欢迎消息，再进入循环；读写出错就退出并关闭连接。

func handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	upgrader := websocket.Upgrader{HandshakeTimeout: 3 * time.Second}
	// 保留默认 CheckOrigin：浏览器带 Origin 时，主机部分须与请求 Host 相同。
	// 不要为了“解决跨域”直接 return true；真实跨站前端应配置明确的来源白名单。
	// Origin 检查不是登录认证，非浏览器客户端可以不带 Origin 或伪造它。
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return // 库已经写了 HTTP 错误响应；升级失败时没有 WebSocket 连接可用。
	}
	defer conn.Close()
	conn.SetReadLimit(1024) // 限制完整消息大小，不只是单个帧的大小。
	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return
	}
	// 服务端主动推送欢迎消息，无需等待客户端先发业务消息。
	if err := conn.WriteMessage(websocket.TextMessage, []byte("welcome")); err != nil {
		return
	}
	for {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return
		}
		kind, message, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Fprintln(os.Stderr, "WebSocket 读取结束：", err)
			}
			return // 出错后退出；不能无限循环读取一个已失败的连接。
		}
		if kind != websocket.TextMessage {
			_ = conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "text only"), time.Now().Add(time.Second))
			return
		}
		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, append([]byte("收到："), message...)); err != nil {
			return
		}
	}
	// 默认 Close handler 会响应对端的关闭帧；默认 Ping handler 会回 Pong。
	// 必须持续 Read 才能处理这些控制帧，即使暂时不关心业务消息也要有读取循环。
}

// ═══════════════════════════════════════════════════════
// 第 3 节：路由与浏览器入口
// ═══════════════════════════════════════════════════════
// GET /ws 负责升级，GET /{$} 只匹配根路径并返回使用说明。
// 普通 HTTP 页面和 WebSocket 端点可以共用同一个 HTTP 服务。

func newWebSocketMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", handleWebSocket)
	mux.HandleFunc("GET /{$}", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(writer, "WebSocket 学习服务：在这个页面的浏览器控制台运行 lesson.go 末尾的 JS，或运行 Go client。")
	})
	return mux
}

// ═══════════════════════════════════════════════════════
// 第 4 节：客户端：欢迎消息、回显与关闭握手
// ═══════════════════════════════════════════════════════
// Dial 完成握手后先读 welcome，再逐条发送文本并验证回显。
// 结束时先发 Close 控制帧，再有期限地等对端确认，最后释放底层连接。

func runWebSocketClient(address string) error {
	dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
	conn, response, err := dialer.Dial(address, nil)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return err
	}
	defer conn.Close()
	conn.SetReadLimit(2048)
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	kind, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	if kind != websocket.TextMessage || string(message) != "welcome" {
		return fmt.Errorf("欢迎消息不符：%q", message)
	}
	fmt.Println("服务端主动推送：", string(message))
	for _, text := range []string{"hello websocket", "你好，Go"} {
		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return err
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			return err
		}
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return err
		}
		kind, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if kind != websocket.TextMessage || string(message) != "收到："+text {
			return fmt.Errorf("回显不符：%q", message)
		}
		fmt.Println("收到消息：", string(message))
	}
	// Close() 直接关闭底层连接；正常退出应先发关闭帧，并在有期限的等待中接收对方关闭帧。
	if err := conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second)); err != nil {
		return err
	}
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return err
	}
	_, _, err = conn.ReadMessage()
	if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		return fmt.Errorf("关闭握手未正常完成：%v", err)
	}
	fmt.Println("关闭握手完成：1000")
	return nil
}

// ═══════════════════════════════════════════════════════
// 第 5 节：自动演示：等待升级后的处理结束
// ═══════════════════════════════════════════════════════
// HTTP 连接升级后交给 WebSocket 管理，HTTP server 的关闭不能代替 WebSocket 的关闭。
// 本例客户端关闭、服务端有期限退出，再用 WaitGroup 确认 handler 已返回。

func demoWebSocket() error {
	mux := newWebSocketMux()
	var handlers sync.WaitGroup
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.Add(1)
		defer handlers.Done()
		mux.ServeHTTP(writer, request)
	}))
	err := runWebSocketClient("ws" + strings.TrimPrefix(server.URL, "http") + "/ws")
	server.Close()
	// HTTP Server.Close/Shutdown 不会代管已经被 Upgrade 接管的连接。
	// 本例客户端会 Close，服务端又有读取期限，再等待所有 handler 真正返回。
	handlers.Wait()
	return err
}

// ═══════════════════════════════════════════════════════
// 第 6 节：双终端运行与单读单写规则
// ═══════════════════════════════════════════════════════
// server 提供 /ws，client 建立连接；末尾也提供浏览器控制台的连接代码。
// 同一 Gorilla 连接应由一个读取流程和一个写入流程负责，多个发送者可先汇入发送队列。

func runWebSocket(args []string) error {
	if len(args) == 0 {
		return demoWebSocket()
	}
	address := "127.0.0.1:8083"
	if len(args) > 1 {
		address = args[1]
	}
	switch args[0] {
	case "server":
		server := &http.Server{Addr: address, Handler: newWebSocketMux(), ReadHeaderTimeout: 3 * time.Second}
		fmt.Println("WebSocket 服务：ws://" + address + "/ws；Ctrl+C 结束；空闲 5 秒会关闭连接")
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case "client":
		return runWebSocketClient("ws://" + address + "/ws")
	default:
		return errors.New("用法：不带参数，或 server/client [地址]，或 exam")
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := runWebSocket(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "WebSocket 演示失败：", err)
		os.Exit(1)
	}
}

// Gorilla 的一条连接同时最多一个 reader 和一个 writer；不能让多个 goroutine 并发 WriteMessage。
// 广播时可为每个连接设一个发送 channel，由唯一写协程消费；channel 还需容量和慢客户端策略。
// Close/WriteControl 有独立并发约定，可与其他方法并行。JSON 可用 WriteJSON/ReadJSON 编解码。
// 读写超时后应关闭连接；真实长连接通常用 Ping/Pong 刷新存活期限，并实现重连策略。
// 本例没有周期心跳，所以空闲 5 秒即关闭，只适合短时教学会话。
//
// 浏览器练习：先启动 server，打开 http://127.0.0.1:8083，再在同页面控制台一次执行：
// const ws = new WebSocket("ws://127.0.0.1:8083/ws");
// ws.onopen = () => ws.send("来自浏览器");
// ws.onmessage = event => { console.log(event.data); if (event.data !== "welcome") ws.close(1000, "done"); };
// ws.onclose = event => console.log("closed", event.code);
// ws.onerror = event => console.error(event);
// 浏览器不能随意为 WebSocket 构造器添加 Authorization 请求头，鉴权方案需要单独设计。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 先核对默认演示的 welcome、两次回显、关闭码 1000，再尝试末尾浏览器代码。
// 找出“发送关闭帧”和“关闭底层连接”两处调用，解释为什么它们不是同一个动作。
