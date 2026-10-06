// 章节：18-错误与异常处理
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// 练习 1：解析分页参数。HTTP query 中 pageSize 是字符串。
// 建议：func parsePageSize(raw string) (int, error)。
// 要求：空串使用默认值 20；其他输入必须是 1..100 的十进制整数。
// 数字解析失败要用 %w 包装 strconv.Atoi 的错误；超出范围返回 *ValidationError。
// 常规："50" -> 50, nil；"" -> 20, nil。
// 边界："abc" 可被 errors.As 识别为 *strconv.NumError；"0" 和 "101" 为字段 pageSize 的错误。
func exercise1() {
	// TODO：实现解析器，用 errors.As 检查错误类型并打印结果。
}

// 练习 2：订单查询的错误链。定义 ErrOrderNotFound 和内存订单 map。
// 建议：func findOrder(id string) (string, error)；func orderSummary(id string) (string, error)。
// 要求：findOrder 对缺失记录返回哨兵错误；orderSummary 增加上下文并用 %w 保留错误链。
// 常规：map 中有 "order-1": "paid" 时，查询得到 "paid", nil。
// 边界："missing" 和空 ID 都返回能被 errors.Is(err, ErrOrderNotFound) 识别的错误。
// 禁止比较 err.Error() 字符串；不要把正常的查无记录改成 panic。
func exercise2() {
	// TODO：定义查询函数与包装函数，分别验证成功和缺失记录。
}

// 练习 3：后台任务保护边界。
// 建议：func executeSafely(job func() error) (err error)。
// 要求：使用直接放在 defer 函数中的 recover；普通 error 原样返回，意外 panic 转为 error。
// 常规：返回 nil 的任务 -> nil；返回哨兵错误的任务 -> errors.Is 仍为 true。
// 边界：panic("broken") -> 非 nil error；nil job -> 非 nil error，不能悄悄成功。
// 在 goroutine 内调用边界函数并通过 channel 取回错误，等待结束，避免主函数提前退出。
// 写一句注释说明：父 goroutine 的 recover 为什么不能保护子 goroutine？
func exercise3() {
	// TODO：实现自己的边界函数，执行四个用例并核对结果。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func parsePageSize(raw string) (int, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：parsePageSize；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

var ErrOrderNotFound = errors.New("订单不存在")
var exerciseOrders = map[string]string{"order-1": "paid"}

func findOrder(id string) (string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：findOrder；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func orderSummary(id string) (string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：orderSummary；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func executeSafely(job func() error) (err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：executeSafely；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
