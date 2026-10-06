// 章节：30-编译与部署
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

// 本章仍保留三道手写练习；构建操作按 lesson.go 末尾的部署附录 在终端完成，不从练习函数调用 shell。

// 练习 1：读取部署配置。
// 建议：func parsePort(raw string) (int, error)。
// 要求：空串使用 8080；其他值必须是 1..65535 的十进制整数；不访问真实环境变量。
// 常规：""->8080；"9090"->9090；"65535"->65535。
// 边界："0"、"65536"、"abc"、"-1" -> error。
// 将 os.Getenv 留在调用处，使 helper 可用表驱动测试验证。
func exercise1() {
	// TODO：实现纯函数，说明编译时配置和运行时配置的区别。
}

// 练习 2：给构建产物加入版本信息。
// 要求：参考 version 再定义 buildTime string，通过 -ldflags -X 注入固定的版本和时间。
// 在 version 模式输出两者；不注入时应有可读默认值。
// 常规：v0.2.0、2026-09-27T12:00:00Z -> 输出一致；默认构建 -> dev 与自己设置的默认时间。
// 构建整个章节到 bin，实际运行本机产物验证；说明 Linux 产物为什么不能在 macOS 原生执行。
func exercise2() {
	// TODO：修改版本输出，并按 lesson.go 末尾的部署附录 构建本机与 Linux 两类产物。
}

// 练习 3：验证优雅退出会等待正在处理的请求。
// 要求：在自己的测试 server 中用 channel 通知“请求已进入”，并阻塞 handler 等待释放信号。
// 发起 Shutdown 后，确认它尚未完成；释放 handler 后，请求与 Shutdown 都应成功结束。
// 边界：换一个很短的 shutdown context 且不释放请求，Shutdown 应返回 deadline exceeded，随后 Close。
// 所有等待都应有超时兜底；不要用 sleep 猜 handler 是否已开始；测试结束必须清理连接与 goroutine。
func exercise3() {
	// TODO：在 exercises_test.go 中验证请求的完成顺序，不需要真实 Linux 服务器。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func parsePort(raw string) (int, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：parsePort；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
