// 章节：08-函数与闭包
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

// 练习 1：把分页参数解析成查询范围。
// 背景：后端分页需要从 page 与 pageSize 得到 offset、limit，并告知调用者参数是否有效。
// 建议在包级实现：
// func parsePagination(page, pageSize int) (offset, limit int, ok bool)
// 要求：
//  1. page 必须在 1..1_000_000 内，pageSize 必须在 1..100 内。
//  2. 不满足条件时返回 0、0、false，不做自动修正。
//  3. 合法时 offset=(page-1)*pageSize，limit=pageSize，ok=true。
//  4. 在 exercise1() 中手写调用，接收并打印所有返回值。
//
// 常规：page=3、pageSize=20 -> 40、20、true。
// 边界：1、100 -> 0、100、true；0、20 -> 0、0、false；1、101 -> 0、0、false。
// 上限：1_000_001、20 -> 0、0、false。此约束也使 offset 在 32 位 int 上安全。
// 提示：命名返回值仍建议显式 return；不要只设置 ok 却留下错误的 offset。
func exercise1() {
	// TODO：在包级写 parsePagination，然后在这里调用常规与边界案例。
}

// 练习 2：用回调筛选一批状态码。
// 背景：同一批日志既可能需要筛选服务端错误，也可能需要筛选成功响应。
// 建议在包级实现：
// func filterStatuses(accept func(int) bool, statuses ...int) []int
// 要求：
//  1. 遍历 statuses，只保留 accept(status) 为 true 的值，顺序保持不变。
//  2. 创建新的结果切片，不修改输入；本题约定回调非 nil。
//  3. 在 exercise2() 中分别写两个匿名回调：500..599 和 200..299。
//  4. 至少一次直接传多个整数，至少一次用 source... 展开已有切片。
//
// 常规：输入 [200, 503, 404, 502, 204]，错误筛选 -> [503, 502]，成功筛选 -> [200, 204]。
// 边界：不传状态码 -> 长度 0；输入 [200, 204] 配错误回调 -> 长度 0。
// nil 切片与非 nil 空切片都可作为本题的空结果，验收时观察 len 而非与 nil 比较。
// 提示：回调只负责规则，filterStatuses 负责遍历与结果收集。
func exercise2() {
	// TODO：在包级写 filterStatuses，并传入两种筛选规则。
}

// 练习 3：创建彼此独立的开发请求编号生成器。
// 背景：本地调试时，希望用 api-100、api-101 这样的编号观察请求顺序。
// 建议在包级实现：
// func makeRequestID(prefix string, start int) func() string
// 要求：
//  1. 返回一个闭包，每次调用生成 prefix + "-" + 当前编号，再把编号加一。
//  2. 每次调用 makeRequestID 都产生独立计数，不使用包级全局计数变量。
//  3. start < 0 时从 0 开始；prefix == "" 时使用 "req"。
//  4. 用 fmt.Sprintf 构造字符串，本题只要求顺序调用，不要求并发安全。
//
// 常规：A := makeRequestID("api", 100)，连续调用得到 "api-100"、"api-101"。
// 独立性：B := makeRequestID("job", 1) -> "job-1"；随后再次调用 A -> "api-102"。
// 边界：makeRequestID("", -3) 连续调用 -> "req-0"、"req-1"。
// 本题只用于学习闭包；这种本地编号不保证跨进程唯一，也不适合代替鉴权令牌。
func exercise3() {
	// TODO：在包级写 makeRequestID，在这里交替调用两组生成器。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func parsePagination(page, pageSize int) (offset, limit int, ok bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：parsePagination；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func filterStatuses(accept func(int) bool, statuses ...int) []int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：filterStatuses；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func makeRequestID(prefix string, start int) func() string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：makeRequestID；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
