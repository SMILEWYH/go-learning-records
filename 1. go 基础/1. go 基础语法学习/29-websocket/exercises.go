// 章节：29-WebSocket双向通信
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 练习 1：定义应用层 JSON 消息。
// 要求：定义 Message{Type string, Text string}，带 json 标签；客户端发 {"type":"chat","text":"hi"}。
// 服务端读 JSON，chat 类型返回 {"type":"ack","text":"hi"}；未知类型或空 text 返回 type=error。
// 每条消息单独 WriteJSON/ReadJSON；保留 ReadLimit、deadline 和 Close。
// 常规：chat/hi -> ack/hi；边界：未知 type、空文本 -> 可识别错误消息，服务不 panic。
func exercise1() {
	// TODO：不要把 JSON 业务 Type 与 WebSocket TextMessage 帧类型混淆。
}

// 练习 2：验证握手的 Origin 策略。
// 要求：用 httptest.Server 和 Dialer 测试同源、异源、不带 Origin 三种情况。
// 默认策略下：Origin 的主机等于 Host -> 101；https://other.example -> 403；不带 Origin -> 可建立。
// 本题服务器用讲解中的默认 CheckOrigin；不要替换成全部放行。
// 成功握手后读 welcome，再正常关闭；失败时检查 HTTP 响应状态并关闭 Body。
// 说明为什么“无 Origin 可连”意味着这项检查不能替代用户鉴权。
func exercise2() {
	// TODO：在自己的测试服务器上验证，不连接外部网站。
}

// 练习 3：把多个生产者的消息交给唯一写协程。
// 要求：两个 goroutine 分别生成 2 条消息，发送到容量为 4 的 channel；唯一 writer 写入同一连接。
// 等生产者都结束后，由协调者关闭 channel，writer 消费完再发关闭帧。
// reader 持续读取 welcome 和 4 条回显；分别保存并按内容核对，不依赖不同生产者的先后顺序。
// 边界：生产者数量为 0 也能结束；读写出错通过取消信号让所有 goroutine 退出。
// 每个阻塞发送/等待都要考虑取消，关闭前确保 writer 不再发送；用 go run -race 验证。
func exercise3() {
	// TODO：结合第 15–17 章实现协调，避免多个 writer 并发操作连接。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func newJSONSocketMux() *http.ServeMux {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：newJSONSocketMux；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func sendProducerMessages(parent context.Context, address string, groups [][]string) ([]string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：sendProducerMessages；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

// ── 自动检查启动入口：填写习题时无需修改 ──
// runExercises 只启动本章的检查入口；检查规则、测试和参考答案都在 exercises_test.go。
// Go 的普通运行不会编译 _test.go，因此通过 go test 进入 TestExam。
func runExercises() {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法定位本章源码，检查未执行")
		return
	}
	dir := filepath.Dir(source)
	if !filepath.IsAbs(dir) { // 兼容以 -trimpath 构建后在课程根目录或章节目录运行。
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("无法定位章节：", err)
			return
		}
		dir = cwd
		if _, err := os.Stat(filepath.Join(dir, "exercises.go")); err != nil {
			dir = filepath.Join(cwd, filepath.Base(filepath.Dir(source)))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "test", "-json", "-count=1", "-run", "^TestExam$", "-timeout", "5m", ".")
	command.Dir = dir
	command.Env = append(os.Environ(), "GO_LESSON_EXAM=1")
	command.WaitDelay = 2 * time.Second
	output, runErr := command.Output()
	decoder := json.NewDecoder(bytes.NewReader(output))
	var diagnostics strings.Builder
	passed := false
	for {
		var event struct{ Action, Test, Output string }
		if err := decoder.Decode(&event); err != nil {
			break
		}
		diagnostics.WriteString(event.Output)
		if event.Test == "TestExam" {
			if event.Action == "pass" {
				passed = true
			}
			if event.Action == "output" && !strings.HasPrefix(event.Output, "=== RUN") && !strings.HasPrefix(event.Output, "--- PASS") {
				fmt.Print(event.Output)
			}
		}
	}
	if runErr != nil || !passed {
		fmt.Println("练习检查未正常结束：", runErr)
		fmt.Print(diagnostics.String())
		if failure, ok := runErr.(*exec.ExitError); ok {
			fmt.Print(string(failure.Stderr))
		}
	}
}
