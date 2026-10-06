// 章节：24-单元测试
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
	"sync"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 本章每题都要自己写业务 helper；exercises_test.go 已提供测试和参考答案，可继续增加测试。
// exercises_test.go 只是教学示例，已有测试通过不代表本章练习已经完成。
// 每题至少包含常规和边界用例；写入文件时只能使用 t.TempDir 下的路径。

// 练习 1：检验后端分页参数。
// 建议：func validatePagination(page, size int) error。
// 规则：page>=1，size 为 1..100。
// 在 exercises_test.go 编写 TestValidatePagination，使用表驱动 + t.Run。
// 常规：(1,20)、(3,100) -> nil；边界：(0,20)、(1,0)、(1,101)、(-1,10) -> error。
// 故意改错一次业务边界，确认测试能够失败，然后恢复正确实现。
func exercise1() {
	// TODO：实现 helper，并新增真正的 _test.go 测试，不用 Println 充当断言。
}

// 练习 2：验证配置文件保存的副作用。
// 建议：func saveFeatureFlag(path string, enabled bool) error。
// 规则：输出恰好为 "enabled=true\n" 或 "enabled=false\n"，覆盖旧内容。
// 在 exercises_test.go 编写 TestSaveFeatureFlag，每个子测试使用 t.TempDir。
// 常规：先保存 false 再保存 true，读回恰好为 "enabled=true\n"，没有旧尾部。
// 边界：父目录不存在 -> error；普通已有文件被正确覆盖；不要用 chmod 制造不稳定的权限测试。
// 用 t.Fatal 检查前置文件操作错误，再断言内容。
func exercise2() {
	// TODO：实现保存函数及独立文件测试，禁止修改仓库文件作为测试样本。
}

// 练习 3：用并发测试验证接口请求计数器。
// 建议：type RequestCounter struct { /* Mutex 和计数值 */ }
//
//	func (c *RequestCounter) Inc()；func (c *RequestCounter) Value() int。
//
// 编写 TestRequestCounter，20 个 goroutine 各递增 100 次，WaitGroup 等待后应为 2000。
// 边界：新计数器为 0；多个读者与写者并行不会触发数据竞争；不要用 sleep 猜测完成时机。
// 测试完成后在主测试 goroutine 断言，运行 go test -race ./24-unit-testing。
// 不要再定义第二个 TestMain：同一个包已有 exercises_test.go 中的 TestMain。
func exercise3() {
	// TODO：实现自己的计数器及 tests；保留初始版讲解测试。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func validatePagination(page, size int) error {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：validatePagination；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func saveFeatureFlag(path string, enabled bool) error {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：saveFeatureFlag；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type RequestCounter struct {
	mu    sync.Mutex
	value int
}

func (c *RequestCounter) Inc() {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Inc；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func (c *RequestCounter) Value() int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Value；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
