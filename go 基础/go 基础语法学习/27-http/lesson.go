/*
章节：27-HTTP网络编程
	本章掌握清单：
	1. HTTP 请求与响应、方法与状态码
	2. 响应模型与 JSON 输出顺序
	3. 路由与输入：路径、查询、JSON 和表单
	4. 响应读取、状态检查与关闭 Body
	5. 客户端：复用 client 并设置超时
	6. 本机自动演示与服务启动
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./27-http
// 运行练习：go run ./27-http exam
// 阅读顺序：按协议、响应、路由、客户端的顺序阅读；main -> runHTTP -> demoHTTP 是默认运行路径。
// TS 对照：HandleFunc 类似路由回调，Request 类似 req；Go 客户端类似 fetch，但要关闭 Body。
// 双终端：go run ./27-http server；go run ./27-http client
// 可在 server/client 后传地址，默认 127.0.0.1:8080。
// 后续再学：中间件、认证授权、CORS、TLS、流式响应、连接池与优雅退出（30 章）。

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节：HTTP 请求与响应、方法与状态码
// ═══════════════════════════════════════════════════════
// 一次请求表达要做什么，响应通过状态码和正文说明结果。
// 先用下面的协议说明建立概念，再把请求中的路径、查询和正文对应到 Go API。

/*
TCP 与 HTTP：
  TCP 负责传输字节；HTTP 是应用层协议，规定方法、URL、头、状态码、正文等语义。
  HTTP/1.1、HTTP/2 通常基于 TCP；HTTP/3 使用基于 UDP 的 QUIC，不能说 HTTP 永远基于 TCP。
  HTTPS 给 HTTP 通信加密并验证对端身份；HTTP/1.1/2 常用 TLS 加 TCP。
  HTTP 是请求/响应模型；一条连接可服务多个请求，不是每次请求必定重新三次握手。
  “无状态”指请求语义本身不替应用保存会话；登录态仍可用 Cookie/Token 等实现。

常用方法（服务器具体支持哪些方法由路由决定）：
  GET 获取资源；HEAD 类似 GET 但不返回正文；POST 提交/创建；PUT 整体替换；
  PATCH 局部修改；DELETE 删除；OPTIONS 查询能力，也常用于浏览器跨域预检。
  幂等指相同请求重复执行的预期效果相同，不要求每次响应正文/状态码相同。
  GET/HEAD 应只读取；PUT/DELETE 在协议语义上幂等，POST/PATCH 通常不保证幂等。
  2xx 成功，3xx 重定向，4xx 请求问题，5xx 服务端问题。
  常见：200、201、204、400、401、403、404、405、413、415、500。
*/

// ═══════════════════════════════════════════════════════
// 第 2 节：响应模型与 JSON 输出顺序
// ═══════════════════════════════════════════════════════
// CreateUserRequest 表达客户端输入，UserResponse 表达返回数据。
// writeJSON 先编码，再设置响应头、状态码、正文，避免编码失败后已发出成功状态。

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func writeJSON(writer http.ResponseWriter, status int, data any) {
	// 先编码，避免编码失败时已经发出 201 等成功状态码。
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, "响应编码失败", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status) // 必须在 Write 前；首次 Write 默认隐式发送 200。
	if _, err := writer.Write(append(body, '\n')); err != nil {
		// 头已发出后无法改写成 500，记录传输错误即可。
		fmt.Fprintln(os.Stderr, "响应写入失败：", err)
	}
}

func bodyError(writer http.ResponseWriter, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		http.Error(writer, "请求体超过 1 KiB", http.StatusRequestEntityTooLarge)
		return
	}
	http.Error(writer, "请求体必须是一个合法 JSON 对象", http.StatusBadRequest)
}

// ═══════════════════════════════════════════════════════
// 第 3 节：路由与输入：路径、查询、JSON 和表单
// ═══════════════════════════════════════════════════════
// newHTTPMux 注册 GET /hello/{name}、POST /users、POST /form 三个路由。
// 路径参数取名字，查询参数决定语言，JSON 或表单正文承载提交内容。
// 读 POST /users 时按媒体类型、大小、语法、尾部内容、业务字段的顺序检查。

func newHTTPMux() *http.ServeMux {
	mux := http.NewServeMux() // 独立路由表，避免多个示例共用全局 DefaultServeMux。
	// Go 1.22+ 支持“方法 路径”模式和 {name} 路径参数。
	mux.HandleFunc("GET /hello/{name}", func(writer http.ResponseWriter, request *http.Request) {
		name := request.PathValue("name")
		language := request.URL.Query().Get("lang")
		greeting := "Hello, " + name
		if language == "zh" {
			greeting = "你好，" + name
		}
		writer.Header().Set("X-Lesson", "27")
		writeJSON(writer, http.StatusOK, map[string]string{"message": greeting})
	})
	mux.HandleFunc("POST /users", func(writer http.ResponseWriter, request *http.Request) {
		mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			http.Error(writer, "请使用 application/json", http.StatusUnsupportedMediaType)
			return
		}
		request.Body = http.MaxBytesReader(writer, request.Body, 1024)
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		var input CreateUserRequest
		if err := decoder.Decode(&input); err != nil {
			bodyError(writer, err)
			return
		}
		// Decode 一次只消费一个 JSON 值；还要拒绝尾部的第二个值和其他垃圾内容。
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			bodyError(writer, err)
			return
		}
		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" {
			http.Error(writer, "name 不能为空", http.StatusBadRequest)
			return
		}
		// 本例固定返回 id=1，只模拟创建响应，不存储用户，不连接数据库。
		writeJSON(writer, http.StatusCreated, UserResponse{ID: 1, Name: input.Name})
	})
	mux.HandleFunc("POST /form", func(writer http.ResponseWriter, request *http.Request) {
		request.Body = http.MaxBytesReader(writer, request.Body, 1024)
		if err := request.ParseForm(); err != nil {
			http.Error(writer, "表单解析失败", http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"name": request.PostForm.Get("name")})
	})
	return mux
	// Request 常用字段：Method、URL.Path、Header、Body、Host、RemoteAddr。
	// 常用方法：URL.Query().Get、PathValue、ParseForm/PostForm、Cookie、Context。
	// FormValue 会合并查询与表单且忽略解析错误；需要校验时先 ParseForm 再读取指定来源。
	// JSON 正文要用 json.Decoder；ParseForm 不会解析 JSON。
	// request.Context() 可传播客户端断开/取消；耗时业务需主动检查 ctx。
	// ResponseWriter 提供 Header、WriteHeader、Write；不是 *ResponseWriter。
	// 服务端在 handler 返回后关闭请求 Body；客户端的响应 Body 由调用方关闭。
	// handler 会并发运行，共享 map/计数器需要锁；本例没有可变共享状态。
}

// ═══════════════════════════════════════════════════════
// 第 4 节：响应读取、状态检查与关闭 Body
// ═══════════════════════════════════════════════════════
// HTTP 通信完成不代表业务成功，404 或 500 也可能得到有效 response。
// printResponse 同时检查正文大小与预期状态，并负责关闭这一份响应体。

func printResponse(response *http.Response, expectedStatus int) error {
	defer response.Body.Close() // 每次请求消费后立刻关闭，别堆在很长的循环外。
	body, err := io.ReadAll(io.LimitReader(response.Body, 4097))
	if err != nil {
		return err
	}
	if len(body) > 4096 {
		return errors.New("示例响应超过 4 KiB")
	}
	fmt.Printf("HTTP %d：%s\n", response.StatusCode, strings.TrimSpace(string(body)))
	if response.StatusCode != expectedStatus {
		return fmt.Errorf("期望状态码 %d，实际 %d", expectedStatus, response.StatusCode)
	}
	return nil
	// 对可控的小响应读到 EOF 再 Close，有助于连接复用；大响应应流式读取或限制大小。
}

// ═══════════════════════════════════════════════════════
// 第 5 节：客户端：复用 client 并设置超时
// ═══════════════════════════════════════════════════════
// client.Get/Post/PostForm 是不同请求的便捷入口，自定义方法使用 NewRequest + Do。
// 这里 DELETE /users 没有注册，预期收到 405，因此这是验证路由规则的成功演示。

func runHTTPClient(baseURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	// http.Get/http.Post 使用默认客户端；这里用 client.Get/Post 设置整体请求超时。
	// 客户端和 Transport 应复用，本函数中的几个请求共用一个 client。
	response, err := client.Get(baseURL + "/hello/Go?lang=zh")
	if err != nil {
		return err
	}
	if err := printResponse(response, http.StatusOK); err != nil {
		return err
	}
	response, err = client.Post(baseURL+"/users", "application/json", strings.NewReader(`{"name":"小林"}`))
	if err != nil {
		return err
	}
	if err := printResponse(response, http.StatusCreated); err != nil {
		return err
	}
	response, err = client.PostForm(baseURL+"/form", url.Values{"name": {"Go student"}})
	if err != nil {
		return err
	}
	if err := printResponse(response, http.StatusOK); err != nil {
		return err
	}
	// PUT/PATCH/DELETE/HEAD/OPTIONS 均可用 NewRequest + Do；net/http 没有 http.Patch 捷径。
	request, err := http.NewRequest(http.MethodDelete, baseURL+"/users", nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	response, err = client.Do(request)
	if err != nil {
		return err
	}
	// 与 fetch 类似，HTTP 4xx/5xx 通常不会成为 Do 的 error，必须检查 StatusCode。
	// 这里只注册了 POST /users，因此 DELETE 得到 405，Allow 头列出允许的方法。
	return printResponse(response, http.StatusMethodNotAllowed)
}

// ═══════════════════════════════════════════════════════
// 第 6 节：本机自动演示与服务启动
// ═══════════════════════════════════════════════════════
// httptest.NewServer 启动可访问的本机 HTTP 服务，URL 中含系统分配的端口。
// server 模式用配置了超时的 http.Server 长期运行；main 再统一处理启动错误。

func demoHTTP() error {
	// httptest.NewServer 启动真实的本机 HTTP 服务，选择空闲端口；也适合写集成测试。
	server := httptest.NewServer(newHTTPMux())
	defer server.Close()
	return runHTTPClient(server.URL)
}

func runHTTP(args []string) error {
	if len(args) == 0 {
		return demoHTTP()
	}
	address := "127.0.0.1:8080"
	if len(args) > 1 {
		address = args[1]
	}
	switch args[0] {
	case "server":
		// 原笔记需改为：先注册路由，再 ListenAndServe；后者会阻塞。
		// 简写 http.ListenAndServe(address, mux)；传 nil 才使用全局路由表。
		server := &http.Server{
			Addr: address, Handler: newHTTPMux(),
			ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 5 * time.Second,
			WriteTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second,
		}
		fmt.Println("HTTP 服务：http://" + address + "；Ctrl+C 结束")
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case "client":
		return runHTTPClient("http://" + address)
	default:
		return errors.New("用法：不带参数，或 server/client [地址]，或 exam")
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := runHTTP(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "HTTP 演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 GET 查询中的 lang=zh 改成 lang=en，核对响应正文和状态码。
// 在 server 模式提交缺少 name、包含未知字段或多个 JSON 值的请求，沿 handler 找到拒绝位置。
