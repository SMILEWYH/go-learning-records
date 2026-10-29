/*
章节:27-HTTP网络编程 - 服务端
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节: HTTP 请求与响应、方法与状态码
// ═══════════════════════════════════════════════════════

/*
TCP 与 HTTP:
  TCP 负责传输字节;HTTP 是应用层协议, 规定方法、URL、头、状态码、正文等语义
  HTTP/1.1、HTTP/2 通常基于 TCP;HTTP/3 使用基于 UDP 的 QUIC, 不能说 HTTP 永远基于 TCP
  HTTPS 给 HTTP 通信加密并验证对端身份;HTTP/1.1/2 常用 TLS 加 TCP
  HTTP 是请求/响应模型;一条连接可服务多个请求, 不是每次请求必定重新三次握手
  “无状态”指请求语义本身不替应用保存会话;登录态仍可用 Cookie/Token 等实现

常用方法（服务器具体支持哪些方法由路由决定）:
  GET 获取资源;HEAD 类似 GET 但不返回正文;POST 提交/创建;PUT 整体替换;
  PATCH 局部修改;DELETE 删除;OPTIONS 查询能力, 也常用于浏览器跨域预检
  幂等指相同请求重复执行的预期效果相同, 不要求每次响应正文/状态码相同
  GET/HEAD 应只读取;PUT/DELETE 在协议语义上幂等, POST/PATCH 通常不保证幂等
  2xx 成功, 3xx 重定向, 4xx 请求问题, 5xx 服务端问题
  常见:200、201、204、400、401、403、404、405、413、415、500
*/

// ═══════════════════════════════════════════════════════
// 第 2 节: 启动 HTTP 服务: 设置地址, 注册路由, 开始监听
// ═══════════════════════════════════════════════════════
func main() {
	address := "127.0.0.1:8080"
	if err := runHTTPServer(address); err != nil {
		fmt.Println("HTTP 服务端失败:", err)
	}
}

func runHTTPServer(address string) error {
	// 先注册路由, 再 ListenAndServe;后者会阻塞当前 goroutine
	// 简写 http.ListenAndServe(address, mux);传 nil 才使用全局路由表
	server := &http.Server{
		Addr: address, Handler: newHTTPMux(),
		ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second,
	}
	fmt.Println("HTTP 服务:http://" + address + "; Ctrl+C 结束")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 路由与输入: 接收请求并处理路径、查询、JSON 和表单
// ═══════════════════════════════════════════════════════
// newHTTPMux 注册 GET /hello/{name}、POST /users、POST /form 三个路由
// 路径参数取名字, 查询参数决定语言, JSON 或表单正文承载提交内容
// 读 POST /users 时按媒体类型、大小、语法、尾部内容、业务字段的顺序检查

func newHTTPMux() *http.ServeMux {
	// 类似前端路由表:先登记“方法 + 路径 -> 处理函数”, 再启动服务接收请求
	mux := http.NewServeMux() // 独立路由表, 避免多个示例共用全局 DefaultServeMux
	// Go 1.22+ 支持“方法 路径”模式和 {name} 路径参数
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
		// Decode 一次只消费一个 JSON 值;还要拒绝尾部的第二个值和其他垃圾内容
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
		// 本例固定返回 id=1, 只模拟创建响应, 不存储用户, 不连接数据库
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
	// Request 常用字段:Method、URL.Path、Header、Body、Host、RemoteAddr
	// 常用方法:URL.Query().Get、PathValue、ParseForm/PostForm、Cookie、Context
	// FormValue 会合并查询与表单且忽略解析错误;需要校验时先 ParseForm 再读取指定来源
	// JSON 正文要用 json.Decoder;ParseForm 不会解析 JSON
	// request.Context() 可传播客户端断开/取消;耗时业务需主动检查 ctx
	// ResponseWriter 提供 Header、WriteHeader、Write;不是 *ResponseWriter
	// 服务端在 handler 返回后关闭请求 Body;客户端的响应 Body 由调用方关闭
	// handler 会并发运行, 共享 map/计数器需要锁;本例没有可变共享状态
}

// ═══════════════════════════════════════════════════════
// 第 4 节: 响应模型与 JSON 输出顺序
// ═══════════════════════════════════════════════════════
// CreateUserRequest 表达客户端输入, UserResponse 表达返回数据
// writeJSON 先编码, 再设置响应头、状态码、正文, 避免编码失败后已发出成功状态

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func writeJSON(writer http.ResponseWriter, status int, data any) {
	// 先编码, 避免编码失败时已经发出 201 等成功状态码
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, "响应编码失败", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status) // 必须在 Write 前;首次 Write 默认隐式发送 200
	if _, err := writer.Write(append(body, '\n')); err != nil {
		// 头已发出后无法改写成 500, 记录传输错误即可
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
