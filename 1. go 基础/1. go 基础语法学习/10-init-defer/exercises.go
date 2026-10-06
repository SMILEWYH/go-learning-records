// 章节：10-初始化与延迟执行
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

/*
练习 1：模拟上传接口的清理流程。

背景：请求取得两个临时资源后，校验可能失败，但已经取得的资源必须释放。
本题只用打印模拟资源，不创建真实文件或连接。

要求：
  - 先打印“打开临时文件”，登记“关闭临时文件”。
  - 再打印“取得上传锁”，登记“释放上传锁”。
  - accepted=false 时打印“校验失败”并提前返回。
  - accepted=true 时打印“保存成功”后返回。
  - 两条路径都用 defer 按资源取得的逆序清理。

建议签名：func simulateUpload(accepted bool)

手动核对：
  - true：打开临时文件 -> 取得上传锁 -> 保存成功 -> 释放上传锁 -> 关闭临时文件。
  - false：打开临时文件 -> 取得上传锁 -> 校验失败 -> 释放上传锁 -> 关闭临时文件。
*/
func exercise1() {
	// TODO：定义 helper，分别传 true/false；先写出预期顺序，再运行核对。
}

/*
练习 2：记录请求开始与结束状态。

背景：审计日志既要保留请求刚进入时的状态，也要记录业务处理后的状态。

要求：
  - 定义局部 status，初值为 pending。
  - 用带参数的 defer 保存初始状态，用无参数闭包在退出时读取最终状态。
  - success=true 时改成 completed；false 时改成 rejected。
  - 调整两个 defer 的登记顺序，让输出先是“开始：pending”，再是最终状态。
  - 两条日志都在函数退出时输出；不能直接写死 pending 来伪装快照。

建议签名：func logRequestStatus(success bool)

手动核对：
  - true：开始：pending，然后结束：completed。
  - false：开始：pending，然后结束：rejected。
  - 在纸上说明：为何捕获变量的闭包与带参数的闭包得到不同状态？
*/
func exercise2() {
	// TODO：实现并运行两种结果，观察登记次序和实际输出次序。
}

/*
练习 3：逐个处理导出文件，及时释放每个资源。

背景：导出任务可能处理大量文件，不能把所有关闭操作积累到批次最后。

要求：
  - 用小函数处理单个文件：打印“打开 <name>”，用 defer 打印“关闭 <name>”。
  - name 为空时打印“跳过空文件名”并返回，不取得资源，也不打印打开或关闭。
  - 合法名称打印“处理 <name>”。
  - batchExport 循环调用小函数，每项处理完即关闭，再处理下一项。
  - 不做真实文件 I/O；思考把 defer 直接放在外层 for 内会有什么差别。

建议签名：func exportOne(name string)；func batchExport(names []string)

手动核对：
  - [a.csv,b.csv]：打开 a -> 处理 a -> 关闭 a -> 打开 b -> 处理 b -> 关闭 b。
    上述 a/b 在实际输出中使用完整文件名。
  - ["",a.csv]：先跳过空文件名，再正常处理 a.csv。
  - nil 或空 slice：没有打开、处理、关闭日志。
*/
func exercise3() {
	// TODO：实现两个 helper，核对每个资源都在下一个任务开始前关闭。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func simulateUpload(accepted bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：simulateUpload；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func logRequestStatus(success bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：logRequestStatus；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func exportOne(name string) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：exportOne；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func batchExport(names []string) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：batchExport；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
