// 章节：02-输入输出
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

// 练习 1：生成一行 HTTP 访问日志。
// 背景：后端日志需要把请求方法、路径、状态码、耗时输出成统一格式。
// 要求：
//  1. 声明 method、path、status、durationMS、cached 五个变量。
//  2. 使用 fmt.Sprintf 生成日志，再用 fmt.Println 输出。
//  3. 耗时固定两位小数，缓存状态使用布尔格式，日志不能包含额外换行。
//
// 常规：GET、/api/users、200、12.345、true ->
// "method=GET path=/api/users status=200 duration_ms=12.35 cached=true"。
// 边界：POST、/api/jobs、202、0.0、false ->
// "method=POST path=/api/jobs status=202 duration_ms=0.00 cached=false"。
// 建议实现位置：exercise1()；学完函数后可抽取
// formatAccessLog(method, path string, status int, durationMS float64, cached bool) string。
func exercise1() {
	// TODO：生成日志字符串并对照题目中的两组预期。
}

// 练习 2：读取本地开发命令的参数。
// 背景：一个简化的开发工具接收“环境名 端口号”，例如 preview 9000。
// 要求：
//  1. 自己在本文件 import 中加入 strings，用 strings.NewReader 模拟输入。
//  2. 用 fmt.Fscan 读取 environment string 和 port int。
//  3. 打印成功读取的项目数 n、err、environment、port。
//  4. 每一组输入都创建新的 Reader 和零值变量，不复用上一组结果。
//
// 常规："preview 9000" -> n=2、err=<nil>、environment="preview"、port=9000。
// 边界："preview" -> n=1、err=EOF、environment="preview"、port=0。
// 错误："preview nope" -> n=1、err 非 nil、environment="preview"、port=0。
// 建议实现位置：exercise2()；此阶段可以分别手写三次扫描，无需提前学习循环。
// 注意：先观察错误，不需要在本题实现错误处理分支，也不要等待真实 stdin。
func exercise2() {
	// TODO：创建内存输入并显示三种扫描结果。
}

// 练习 3：构造可发送的文本响应。
// 背景：健康检查端点可能返回简单的纯文本响应。
// 要求：
//  1. 在 import 中加入 bytes，声明一个 bytes.Buffer。
//  2. 用 fmt.Fprintf 写入 "service=<服务名>\nready=<状态>\n"。
//  3. 打印返回的写入字节数和错误，再用 %q 检查 buffer.String() 的换行。
//
// 常规：服务名 "api"、状态 true -> "service=api\nready=true\n"，写入 23 字节。
// 边界：服务名 ""、状态 false -> "service=\nready=false\n"，写入 21 字节。
// 建议实现位置：exercise3()；可在学完函数后抽取
// formatHealth(service string, ready bool) string。
// 注意：n 统计字节；如果换成中文服务名，n 不等于肉眼看到的字符个数。
func exercise3() {
	// TODO：写入两组响应并核对字符串、字节数与错误。
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
