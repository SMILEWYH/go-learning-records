// 章节：06-分支选择
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

// 练习 1：选择 HTTP 方法对应的业务动作。
// 背景：一个简化的 API 路由层需要把请求方法映射成动作名称。
// 要求：使用 switch method，遵循以下规则。
//   - GET / HEAD -> "read"。
//   - POST -> "create"。
//   - PUT / PATCH -> "update"。
//   - DELETE -> "delete"。
//   - 其他值 -> "unsupported"。
//
// 常规："PATCH" -> "update"；"HEAD" -> "read"。
// 边界：""、"get"、"TRACE" -> "unsupported"，本题不自动转为大写。
// 建议实现位置：exercise1()，多值 case 用逗号列出，不需要 fallthrough。
// 学完函数后可抽取 actionForMethod(method string) string。
func exercise1() {
	// TODO：声明 method，通过 switch 输出动作名称。
}

// 练习 2：给 HTTP 状态码分类。
// 背景：监控面板需要统计不同响应类别。
// 要求：使用不带表达式的 switch，按范围输出类别。
//   - 100..199 -> "informational"；200..299 -> "success"。
//   - 300..399 -> "redirect"；400..499 -> "client_error"。
//   - 500..599 -> "server_error"；其余 -> "invalid"。
//
// 常规：204 -> "success"；404 -> "client_error"；503 -> "server_error"。
// 边界：99 -> "invalid"；100 -> "informational"；599 -> "server_error"；600 -> "invalid"。
// 建议实现位置：exercise2()；修改 status 的值逐一检查分界点即可。
// 学完函数后可抽取 classifyStatus(status int) string。
// 提示：写成一边包含、一边不包含的区间，检查有没有遗漏或重叠。
func exercise2() {
	// TODO：用 switch 的布尔 case 判断状态码所在区间。
}

// 练习 3：解释异步导出任务的状态。
// 背景：前端查询后端导出任务时，需要知道状态是否结束，以及下一步提示。
// 要求：根据 state string，同步设置 terminal bool 和 message string。
//   - "queued" / "running" -> false, "稍后查询"。
//   - "succeeded" -> true, "可以下载"。
//   - "failed" / "canceled" -> true, "任务未完成"。
//   - 其他值 -> false, "未知状态"。
//
// 常规："running" -> false, "稍后查询"；"succeeded" -> true, "可以下载"。
// 边界："" -> false, "未知状态"；"canceled" -> true, "任务未完成"。
// 建议实现位置：exercise3()，先声明结果变量，再在 switch 中赋值。
// 学完函数后可抽取 describeJob(state string) (terminal bool, message string)。
// 思考：此题使用 fallthrough 会不会误把未结束的状态标记为结束？
func exercise3() {
	// TODO：用多值 case 输出两个结果。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func actionForMethod(method string) string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：actionForMethod；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func classifyStatus(status int) string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：classifyStatus；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func describeJob(state string) (terminal bool, message string) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：describeJob；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
