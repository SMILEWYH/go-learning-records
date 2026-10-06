// 章节：26-TCP网络编程
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 练习 1：为 TCP 聊天准备完整消息。
// 建议：func encodeLine(message string) ([]byte, error)。
// 要求：拒绝正文含 \n 或 \r；UTF-8 字节数加末尾 \n 不超过 1024；成功追加换行。
// 常规："hello" -> []byte("hello\n")；"你好" -> 7 字节（包含换行）。
// 边界：空串 -> []byte("\n")；1023 个 ASCII 字符成功；1024 个失败；正文含换行失败。
func exercise1() {
	// TODO：用 len 计算字节数，并检查协议不允许的分隔符。
}

// 练习 2：用同一个连接交换多条消息。
// 要求：在 127.0.0.1:0 启动本地服务，客户端发送 a、b、c，依次收到 A、B、C。
// 必须复用同一个 conn 和同一个 bufio.Reader；服务端接待完这个客户端后退出。
// 边界：把 "a\nb\n" 放进一次 Write，仍应读取到两条响应；空消息收到空消息。
// 使用 deadline 限制等待，并通过 channel/WaitGroup 等待服务端结束；不要用 sleep 同步。
func exercise2() {
	// TODO：可复用讲解的服务端，自己编写发送与核对流程。
}

// 练习 3：区分正常结束、半条消息和超时。
// 要求：写测试检查 readMessage 的三种结果：完整 "ok\n"，空流 EOF，"unfinished" 的 UnexpectedEOF。
// 再用 net.Pipe 建立内存连接；对端不发送，设置短读取 deadline，检查 os.IsTimeout(err)。
// net.Pipe 不是 TCP，只用它模拟阻塞读；连接两端都需关闭，不用“超时后继续无限 Read”。
// 用 errors.Is 判断 EOF 类错误，不比较 err.Error() 文本。
func exercise3() {
	// TODO：写入 exercises_test.go，确保测试本身不会永久阻塞。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func encodeLine(message string) ([]byte, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：encodeLine；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func exchangeLines(address, wire string, count int) ([]string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：exchangeLines；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
