// 章节：03-基本数据类型
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

// exam 自动检查三道练习；前四章填写 exercise1/2/3，输出要求见 exercises_test.go。

// 练习 1：用整数分计算购物车总额。
// 背景：价格 19.99 元不适合依赖二进制浮点数做精确金额累加。
// 要求：
//  1. 声明单价 priceCents int64、数量 quantity int、运费 shippingCents int64。
//  2. 先把 quantity 转为 int64，再计算商品金额与含运费总额。
//  3. 用整数除法、取余和 %02d 显示两位小数金额，不经过 float64。
//
// 常规：priceCents=1999、quantity=3、shippingCents=600 -> 6597 分，显示 "65.97"。
// 边界：priceCents=0、quantity=0、shippingCents=0 -> 0 分，显示 "0.00"。
// 另一组：priceCents=105、quantity=1、shippingCents=0 -> 显示 "1.05"。
// 本题约定所有值非负且运算不溢出，不需要提前处理校验和异常。
// 建议实现位置：exercise1()；学完函数后可抽取
// totalCents(priceCents int64, quantity int, shippingCents int64) int64。
func exercise1() {
	// TODO：写出类型转换、金额运算与金额显示。
}

// 练习 2：看清 API 昵称的字节数与码点数。
// 背景：数据库字段限制和用户输入长度的定义可能不同。
// 要求：
//  1. 在本文件导入 unicode/utf8。
//  2. 使用 len 与 utf8.RuneCountInString 统计昵称，并用 %q 输出原文本。
//  3. 对三组文本分别手写调用，暂时不要求循环。
//
// 常规："Go语言" -> 8 字节、4 个 rune。
// 边界："" -> 0 字节、0 个 rune。
// emoji："A😀" -> 5 字节、2 个 rune。
// 建议实现位置：exercise2()；后续可抽取 textLength(text string) (bytes, runes int)。
// 思考：为什么 text[0] 不能普遍代表第一个完整汉字？
// 提示：不要对空字符串做下标读取；本题只需要统计，不要求截断字符串。
func exercise2() {
	// TODO：输出每组文本的字节数和 rune 数量。
}

// 练习 3：计算后台任务完成率。
// 背景：管理端需要显示一批后台任务的完成百分比。
// 要求：
//  1. 用 int 保存 completed 和 total，在相除之前转为 float64。
//  2. 百分比乘以 100，用 %.1f 和 %% 显示，例如 "37.5%"。
//  3. 同时输出 completed / total 的整数结果，观察与浮点除法的差异。
//
// 常规：completed=3、total=8 -> "37.5%"；整数除法为 0。
// 边界：completed=0、total=8 -> "0.0%"；completed=8、total=8 -> "100.0%"。
// 本题约定 total > 0 且 0 <= completed <= total；零总数的处理留到条件判断章节。
// 建议实现位置：exercise3()；后续可抽取 progress(completed, total int) float64。
// 思考：先写 float64(completed / total) 为什么无法找回丢失的小数？
func exercise3() {
	// TODO：对比整数与浮点除法，输出格式化百分比。
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
