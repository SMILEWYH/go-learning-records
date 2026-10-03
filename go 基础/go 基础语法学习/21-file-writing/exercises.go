// 章节：21-文件写入
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

// 所有写入练习只操作 os.MkdirTemp 创建的目录，并 defer os.RemoveAll 清理；
// 不使用用户传入的任意路径，不修改仓库配置。每个打开的文件都检查写入和 Close 错误。

// 练习 1：保存开发环境配置。
// 建议：func saveSettings(path string, values map[string]string) error。
// 要求：json.MarshalIndent 编码后用 os.WriteFile 写入，末尾加换行，新文件权限 0o600。
// 常规：写入 {"theme":"dark"} 后能读回合法 JSON。
// 边界：较短内容覆盖较长内容不留旧尾巴；空 map 和 nil 输入都保存为 {}，不是 null。
// 目标父目录缺失时报错；注意先处理 nil map，再进行 JSON 编码。
// 不要求 map 字段的输出顺序；读取后通过 json.Unmarshal 核对内容。
func exercise1() {
	// TODO：在临时目录实现配置保存，验证覆盖行为和失败路径。
}

// 练习 2：追加审计日志，每个事件单独一行 JSON（JSON Lines）。
// 建议：type AuditEvent struct { Actor string `json:"actor"`; Action string `json:"action"` }
//
//	func appendAudit(path string, event AuditEvent) error。
//
// 要求：使用 O_CREATE|O_WRONLY|O_APPEND 和 json.Encoder.Encode；字段均不能为空。
// 常规：顺序追加两个事件，得到两行独立合法 JSON，之前的内容仍在。
// 边界：Actor 内含换行会被 JSON 转义，不增加物理行数；空 Actor/Action 返回 error 且不写入。
// 本题只需顺序调用；多进程共同写文件需要另外设计并发协议。
func exercise2() {
	// TODO：创建临时日志并验证内容、行数、校验失败时文件保持不变。
}

// 练习 3：只创建一次的本地任务标记。
// 建议：func createJobMarker(path, jobID string) error。
// 要求：O_CREATE|O_EXCL|O_WRONLY，权限 0o600；写入 jobID 和换行；jobID 为空先返回 error。
// 常规：首次创建成功，内容为 "job-1\n"。
// 边界：同一路径第二次创建返回能被 errors.Is(err, os.ErrExist) 识别的错误，旧内容不变。
// 创建失败不能随后继续写；若写或关闭失败，应把错误返回给调用者。
func exercise3() {
	// TODO：在临时目录验证首次写入、重复创建及空 jobID。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func saveSettings(path string, values map[string]string) error {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：saveSettings；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type AuditEvent struct {
	Actor  string `json:"actor"`
	Action string `json:"action"`
}

func appendAudit(path string, event AuditEvent) (err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：appendAudit；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func createJobMarker(path, jobID string) (err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：createJobMarker；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
