// 章节：20-文件读取
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 文件练习统一约定：用 os.MkdirTemp 创建自己的目录，defer os.RemoveAll 清理；
// 只在该目录内创建/读取测试文件，不读取项目源码或用户真实配置。

// 练习 1：加载一个服务的 JSON 配置。
// 建议：type AppConfig struct { Name string `json:"name"`; Port int `json:"port"` }
//
//	func readConfig(path string) (AppConfig, error)。
//
// 要求：os.ReadFile + json.Unmarshal；校验 Name 非空、Port 为 1..65535；错误保留上下文。
// 常规：{"name":"api","port":8080} -> 对应结构体, nil。
// 边界：缺失文件、非法 JSON、port=0、空 name 都返回 error，不能 panic。
func exercise1() {
	// TODO：创建临时配置文件，实现解析并验证常规及错误输入。
}

// 练习 2：统计结构化日志中以 "ERROR " 开头的行。
// 建议：func countErrorLines(path string) (int, error)。
// 要求：用 Scanner，设置 64 KiB 缓冲上限，并检查 Scanner.Err 和关闭错误。
// 常规："INFO start\nERROR db\nERROR cache" -> 2, nil；最后一行无换行也要处理。
// 边界：空文件 -> 0, nil；"INFO ERROR text" 不计数；一行 70 KiB -> error。
// 出错时统一返回 0, error，避免调用者把不完整统计当成最终结果。
func exercise2() {
	// TODO：在临时目录创建几种日志，用明确预期验证。
}

// 练习 3：给 HTTP 请求体读取设置大小限制（本题不需要启动服务）。
// 建议：func readLimited(reader io.Reader, maxBytes int64) ([]byte, error)。
// 要求：maxBytes 为 1..1 MiB；最多尝试读取 maxBytes+1 字节，超限返回 error。
// 可用 io.LimitReader + io.ReadAll；不要先无上限读取完整内容再判断大小。
// 常规：strings.NewReader("abc"), 3 -> "abc", nil。
// 边界："abcd", 3 -> error；空内容, 3 -> 长度 0；非法上限 -> error。
// 扩展：实现一个最后一次 Read 同时返回数据和 io.EOF 的 reader，确认不会丢数据。
func exercise3() {
	// TODO：实现与文件无关的 Reader 函数，用内存 reader 验证。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type AppConfig struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func readConfig(path string) (AppConfig, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：readConfig；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func countErrorLines(path string) (count int, err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：countErrorLines；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func readLimited(reader io.Reader, maxBytes int64) ([]byte, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：readLimited；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
