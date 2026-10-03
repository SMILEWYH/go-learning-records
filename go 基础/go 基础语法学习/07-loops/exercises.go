// 章节：07-循环遍历
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

// 练习 1：聚合一批接口响应状态。
// 背景：后台监控要统计一批请求中每个状态码的出现次数及服务端错误总数。
// 要求：
//  1. 输入为 []int，用 range 遍历。
//  2. 用 map[int]int 累加每种状态码的次数；map 必须先初始化。
//  3. 单独累计 500 <= code < 600 的数量。
//
// 常规：[200, 200, 404, 503, 503, 502] -> 200:2、404:1、503:2、502:1，服务端错误 3 次。
// 边界：空切片或 nil 切片 -> map 长度 0，服务端错误 0 次。
// 本题约定输入都是有效的 100..599 状态码；不要求打印 map 时具有固定顺序。
// 建议实现位置：exercise1()；后续可抽取 countStatuses(codes []int) (map[int]int, int)。
// 提示：尚未出现的 key 读取结果是 int 零值 0，可以直接参与累加。
func exercise1() {
	// TODO：遍历状态码，输出每个状态码的次数和服务端错误总数。
}

// 练习 2：给批量提交的用户 ID 去重。
// 背景：前端合并多次勾选时，同一个 ID 可能重复出现，也可能混入空字符串。
// 要求：
//  1. 忽略空字符串，使用 continue 跳过。
//  2. 用 map[string]bool 记录已经出现的 ID。
//  3. 用切片收集结果，保留每个 ID 第一次出现的顺序。
//
// 常规：["u-02", "", "u-01", "u-02", "u-03", "u-01"] -> ["u-02", "u-01", "u-03"]。
// 边界：[] 或 ["", ""] -> 长度为 0 的结果；["u-01", "u-01"] -> ["u-01"]。
// 本题不裁剪空格，所以 " u-01 " 与 "u-01" 被视为不同 ID。
// 建议实现位置：exercise2()；后续可抽取 uniqueIDs(ids []string) []string。
// 提示：从原切片的遍历顺序构造结果，不能最后 range map 来还原首次出现顺序。
func exercise2() {
	// TODO：用 seen map 和结果切片完成去重。
}

// 练习 3：模拟有上限的任务重试。
// 背景：后端请求失败时可以重试，但必须限制次数，并在成功后立刻停止。
// 要求：输入 outcomes []bool 和 maxAttempts int，每个 bool 代表一次尝试是否成功。
//  1. 最多读取 maxAttempts 项，也不能超出 outcomes 的长度。
//  2. 第一次读到 true 就停止，并记录 succeeded=true。
//  3. 输出真正执行的 attempts 数和 succeeded；本题不访问网络、不延时。
//
// 常规：[false, false, true, false]、maxAttempts=4 -> attempts=3、succeeded=true。
// 达到上限：[false, false, true]、maxAttempts=2 -> attempts=2、succeeded=false。
// 边界：[]、maxAttempts=3 -> 0,false；[true]、maxAttempts=0 或 -1 -> 0,false。
// 建议实现位置：exercise3()；后续可抽取 simulateRetry(outcomes []bool, maxAttempts int) (int, bool)。
// 提示：使用有界 for，把“是否还有输入”和“是否还有尝试次数”同时写入条件。
func exercise3() {
	// TODO：实现有终止条件的模拟循环，输出真实尝试次数。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func countStatuses(codes []int) (map[int]int, int) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：countStatuses；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func uniqueIDs(ids []string) []string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：uniqueIDs；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func simulateRetry(outcomes []bool, maxAttempts int) (int, bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：simulateRetry；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
