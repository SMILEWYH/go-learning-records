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
	TCP 和 HTTP 是不同层的协议: TCP 是传输层协议, 负责可靠地传输字节, HTTP 是应用层协议, 负责规定这些字节表达什么请求、什么响应
	HTTP 包含了请求方法、HTTP 版本、请求头、请求体、响应头、响应体、状态码
  	HTTP/1.1、HTTP/2 通常基于 TCP; HTTP/3 使用基于 UDP 的 QUIC, 所以不能说 HTTP 永远基于 TCP
  	HTTPS 给 HTTP 通信加密并验证对端身份; HTTP/1.1/2 常用 TLS 加 TCP
 	HTTP 是请求/响应模型; 一条连接可服务多个请求, 不是每次请求必定重新三次握手
	HTTP 协议本身是“无状态”的, 不会自动在多次请求之间保存和关联用户的会话信息; 但应用仍可通过 Cookie、Token 等机制实现登录态

常见方法（服务器具体支持哪些方法由路由决定）:
	GET 请求:
	POST 请求:
	PUT 请求:
	PATCH 请求:
	DELETE 请求:
	OPTIONS 请求:


  GET 获取资源;HEAD 类似 GET 但不返回正文;POST 提交/创建;PUT 整体替换;
  PATCH 局部修改;DELETE 删除;OPTIONS 查询能力, 也常用于浏览器跨域预检
  幂等指相同请求重复执行的预期效果相同, 不要求每次响应正文/状态码相同
  GET/HEAD 应只读取;PUT/DELETE 在协议语义上幂等, POST/PATCH 通常不保证幂等
  2xx 成功, 3xx 重定向, 4xx 请求问题, 5xx 服务端问题
  常见:200、201、204、400、401、403、404、405、413、415、500

常见响应状态码:
	2xx : 请求成功
		200 : 请求成功
		201 : 已创建资源
		204 : 请求成功, 没有响应正文

	3xx : 重定向
		301 : 网站永久迁移到新地址
		302 : 临时跳转到另一个地址
		304 : 资源未修改，可复用已有缓存

	4xx : 请求存在问题，例如缺少身份凭据、资源不存在
		400 : 请求 JSON 格式错误、参数不符合要求
		401 : 缺少有效的身份凭据
		403 : 服务器拒绝处理 (已登录，但没有访问权限)
		404 : 找不到目标资源
		405 : 目标不支持该请求方法

	5xx : 服务端处理请求时发生问题
		500 : 服务端内部错误
		503 : 服务暂时不可用
		504 : 网关等待上游响应超时
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
	server := http.Server{
		Addr:              address,
		Handler:           newHTTPMux(),    // 请求路由表
		ReadHeaderTimeout: 3 * time.Second, // 读取请求头的超时
		ReadTimeout:       5 * time.Second, // 读取整个请求的超时
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second, // 等待下一个 HTTP 请求的空闲超时
	}
	fmt.Println("HTTP 服务:http://" + address + "; Ctrl+C 结束")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 请求路由表
// ═══════════════════════════════════════════════════════
func newHTTPMux() *http.ServeMux {
	mux := http.NewServeMux()
	// Go 1.22+ 支持“方法 路径”模式和 {name} 路径参数
	mux.HandleFunc("GET /hello/{name}", func(writer http.ResponseWriter, request *http.Request) {
		name := request.PathValue("name")           // 动态路径参数
		language := request.URL.Query().Get("lang") // 查询 ? 后的参数
		greeting := "Hello, " + name
		if language == "zh" {
			greeting = "你好，" + name
		}
		writer.Header().Set("X-Lesson", "27") // 自定义响应头
		writeJSON(writer, http.StatusOK, map[string]string{"message": greeting})
	})
	mux.HandleFunc("POST /users", func(writer http.ResponseWriter, request *http.Request) {
		// 检查请求头的 Content-Type 参数
		mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			http.Error(writer, "请使用 application/json", http.StatusUnsupportedMediaType)
			return
		}
		// 限制请求体大小
		request.Body = http.MaxBytesReader(writer, request.Body, 1024)
		// 创建 JSON 解码器
		decoder := json.NewDecoder(request.Body)
		// 设置规则: 遇到未知字段时报错
		decoder.DisallowUnknownFields()
		var input CreateUserRequest
		// 读取正文, 按规则解码并写入 input
		if err := decoder.Decode(&input); err != nil {
			bodyError(writer, err)
			return
		}
		// Decode 一次只消费一个 JSON 值; 还要拒绝尾部的第二个值和其他垃圾内容
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
	writer.WriteHeader(status) // 设置状态码, 必须在 writer.Write() 前; 首次 Write 默认隐式发送 200
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
