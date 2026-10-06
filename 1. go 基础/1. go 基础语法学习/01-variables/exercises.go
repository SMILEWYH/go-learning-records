// 章节：01-变量与常量
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

// 练习 1：整理本地 API 服务配置。
// 背景：前端转到后端后，需要先理解端口、部署环境和调试开关的类型。
// 要求：
//  1. 用 const 声明默认端口 8080，用 var 声明 serviceName = "profile-api"。
//  2. 用 := 声明 environment = "development"、debug = true。
//  3. 创建可修改的 port 变量并从默认端口初始化，随后改为 9000。
//  4. 用 fmt.Println 输出这些值，并用 %T 查看各自的类型。
//
// 验收：常规配置应包含 profile-api、development、true、9000，默认端口仍为 8080。
// 边界：把 serviceName 改为 ""、debug 改为 false，输出应忠实保留空字符串和 false。
// 建议实现位置：直接在 exercise1() 内完成；本题不需要额外的 helper 函数。
// 思考：TS 的 const config = {} 合法，为什么 Go 的 const 不能保存这种配置对象？
func exercise1() {
	// TODO：手写变量声明、重新赋值以及输出。
}

// 练习 2：观察尚未加载的功能开关。
// 背景：后端配置未提供初值时，Go 会先给变量零值，不会产生 undefined。
// 要求：
//  1. 在一组 var (...) 声明中创建 maxRetries int、region string、enabled bool。
//  2. 不赋初值，先输出它们；用 %q 输出 region，让空字符串也看得见。
//  3. 再分别赋值为 3、"ap-shanghai"、true，并再次输出。
//
// 验收：初始输出为 0、""、false；赋值后为 3、"ap-shanghai"、true。
// 边界：将 maxRetries 再次设为 0，确认 0 也是合法值，不能自动理解成“没配置”。
// 建议实现位置：直接在 exercise2() 内完成。
// 思考：仅凭 enabled == false，是否能区分“用户关闭”与“用户没填写”？
// 本题只需说明不能区分，后续章节再学习如何额外记录配置是否存在。
func exercise2() {
	// TODO：输出零值，再输出手动配置后的值。
}

// 练习 3：避免临时配置意外覆盖外部配置。
// 背景：调试某段代码时，常会临时修改当前环境名。
// 要求：
//  1. 在函数里创建 environment := "production"。
//  2. 创建一个独立的 {...} 代码块，在块内用 := 声明同名变量 "preview"。
//  3. 分别输出块内、块外的 environment。
//  4. 再创建一个块，改用 = 将 environment 赋为 "staging"，观察块外的值。
//
// 验收：第一次块内为 preview，块外为 production；第二次块外变成 staging。
// 边界：将内层 := 的值改成 ""，外层仍应为 production。
// 建议实现位置：直接在 exercise3() 内完成。
// 思考：把所有内层 := 换成 =，哪些输出会发生变化？为什么？
func exercise3() {
	// TODO：手写两个代码块，对比 := 声明与 = 赋值。
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
