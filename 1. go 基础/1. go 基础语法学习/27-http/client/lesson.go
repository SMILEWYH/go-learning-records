/*
章节:27-HTTP网络编程 - 客户端
在本章目录执行: go run ./client/lesson.go
先启动服务端, 再在另一个终端启动客户端.
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节: 启动客户端: 设置 HTTP 服务地址
// ═══════════════════════════════════════════════════════
// TS 对照: Go 客户端类似 fetch, 但要关闭响应 Body
func main() {
	address := "http://127.0.0.1:8080"
	if err := runHTTPClient(address); err != nil {
		fmt.Println("HTTP 客户端失败:", err)
		os.Exit(1)
	}
}

// ═══════════════════════════════════════════════════════
// 第 2 节: 创建客户端并发送请求: GET、JSON、表单与自定义方法
// ═══════════════════════════════════════════════════════
// client.Get/Post/PostForm 是不同请求的便捷入口, 自定义方法使用 NewRequest + Do
// 这里 DELETE /users 没有注册, 预期收到 405, 因此这是验证路由规则的成功演示

func runHTTPClient(baseURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	// http.Get/http.Post 使用默认客户端;这里用 client.Get/Post 设置整体请求超时
	// 客户端和 Transport 应复用, 本函数中的几个请求共用一个 client
	// 类似 await fetch(url):先得到响应, 再读 Body;这里会阻塞当前 goroutine
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
	// PUT/PATCH/DELETE/HEAD/OPTIONS 均可用 NewRequest + Do;net/http 没有 http.Patch 捷径
	request, err := http.NewRequest(http.MethodDelete, baseURL+"/users", nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	response, err = client.Do(request)
	if err != nil {
		return err
	}
	// 与 fetch 类似, HTTP 4xx/5xx 通常不会成为 Do 的 error, 必须检查 StatusCode
	// 这里只注册了 POST /users, 因此 DELETE 得到 405, Allow 头列出允许的方法
	return printResponse(response, http.StatusMethodNotAllowed)
}

// ═══════════════════════════════════════════════════════
// 第 3 节: 读取响应: 检查状态码并关闭 Body
// ═══════════════════════════════════════════════════════
// HTTP 通信完成不代表业务成功, 404 或 500 也可能得到有效 response
// printResponse 同时检查正文大小与预期状态, 并负责关闭这一份响应体

func printResponse(response *http.Response, expectedStatus int) error {
	defer response.Body.Close() // 每次请求消费后立刻关闭, 别堆在很长的循环外
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
	// 对可控的小响应读到 EOF 再 Close, 有助于连接复用;大响应应流式读取或限制大小
}
