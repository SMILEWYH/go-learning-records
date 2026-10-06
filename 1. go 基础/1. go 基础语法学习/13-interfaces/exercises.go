// 章节：13-接口与类型断言
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
练习 1：替换通知发送方式。

背景：账户服务要发站内信或邮件，但业务代码只需要“发送文本”能力。
本题只返回模拟结果，不连接真实邮件服务。

要求：
  - 定义 Notifier 接口，包含 Send(message string) string。
  - 定义 InboxNotifier 和 EmailNotifier，各自隐式实现接口。
  - Inbox 返回“inbox:<message>”，Email 返回“email:<message>”。
  - notifyAll 遍历 []Notifier，顺序收集结果并返回 []string。
  - 不用 type switch 判断实现类型；通过接口直接调用 Send。

建议签名：func notifyAll(notifiers []Notifier, message string) []string

手动核对：
  - [InboxNotifier{}, EmailNotifier{}]、"welcome"：返回 ["inbox:welcome","email:welcome"]。
  - 空消息：分别返回 "inbox:" 和 "email:"。
  - 空切片：结果长度为 0；本题约定切片内不含 nil 接口。
*/
func exercise1() {
	// TODO：定义接口与两个实现，再通过统一入口调用。
}

/*
练习 2：校验来源不统一的分页大小。

背景：某个内部聚合层提供 any 值；它可能是整数、文本、布尔值或 nil。
在进入业务逻辑之前，需要明确允许哪些类型。

要求：
  - int 类型：仅接受 1 到 100。
  - string 类型：本题仅识别 "small"->10、"large"->50，其余拒绝。
  - 其他类型（包括 float64、bool、nil）均拒绝。
  - 成功返回值与 true；失败返回 0 与 false。
  - 用 type switch 完成；不要把类型断言当作数值转换。

建议签名：func parseLimit(value any) (int, bool)

手动核对：
  - 20 -> 20,true；"small" -> 10,true；"large" -> 50,true。
  - 0、101、"20"、float64(20)、true、nil -> 0,false。
  - 在注释中说明为什么 any(float64(20)).(int) 不是安全转换。
*/
func exercise2() {
	// TODO：实现 helper，逐项打印输入、类型和校验结果。
}

/*
练习 3：正确表示“没有启用的功能配置”。

背景：配置工厂返回接口；未启用时，调用方希望用 config == nil 判断。

要求：
  - 定义 Feature 接口，含 Name() string。
  - 定义 FeatureConfig，给 *FeatureConfig 实现 Name，正常返回 "export"。
  - createFeature(false) 直接返回 nil 接口；true 返回有效的 *FeatureConfig。
  - 在练习函数另建 var missing *FeatureConfig，再赋给 Feature，比较它与 nil。
  - 不调用这个 nil 指针的方法，只解释其接口动态类型和值。

建议签名：func createFeature(enabled bool) Feature

手动核对：
  - createFeature(false) == nil 为 true。
  - createFeature(true) != nil，调用 Name() 返回 export。
  - Feature(missing) == nil 为 false。
  - 写出为什么 FeatureConfig{} 不满足使用指针接收者的 Feature 接口。
*/
func exercise3() {
	// TODO：实现工厂，并对比真正的 nil 接口和包着 nil 指针的接口。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type Notifier interface{ Send(message string) string }
type InboxNotifier struct{}
type EmailNotifier struct{}

func (InboxNotifier) Send(message string) string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Send；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func (EmailNotifier) Send(message string) string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Send；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func notifyAll(notifiers []Notifier, message string) []string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：notifyAll；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func parseLimit(value any) (int, bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：parseLimit；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type Feature interface{ Name() string }
type FeatureConfig struct{}

func (*FeatureConfig) Name() string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Name；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func createFeature(enabled bool) Feature {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：createFeature；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
