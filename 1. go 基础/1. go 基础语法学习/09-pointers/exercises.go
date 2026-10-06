// 章节：09-指针与值传递
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
练习 1：更新接口中的可选分页参数。

背景：前端未填写 pageSize 时用默认值，但输入 0 也属于需要纠正的值。
本题用 *int 区分“未提供 nil”与“明确提供了一个数值”。

要求：
  - nil 返回默认值 20。
  - 指向的值小于 1 时，将原变量改成 1，并返回 1。
  - 大于 100 时，将原变量改成 100，并返回 100。
  - 其他数值保持原样并返回。
  - 不解引用 nil，也不通过重新指向局部变量假装修改调用方。

建议签名：func normalizePageSize(size *int) int

手动核对：
  - size = 30：返回 30，原变量仍为 30。
  - size = 0：返回 1，原变量变成 1。
  - size = 200：返回 100，原变量变成 100。
  - nil：返回 20，不发生 panic。
*/
func exercise1() {
	// TODO：在文件中定义 helper；在这里创建变量、传入地址并打印修改前后的值。
}

/*
练习 2：为批量导入任务制作独立的重试次数列表。

背景：界面展示原始重试次数，任务准备阶段需要每项加 1，不能修改展示数据。

要求：
  - 创建并返回独立的 []int，每个值为原值加 1。
  - 不修改传入的 slice；返回后修改结果的元素也不能影响原 slice。
  - 输入 nil 或空 slice 时，返回长度为 0 的 slice；不限定是否为 nil。
  - 说明直接执行 result := source 为什么不能隔离元素修改。

建议签名：func nextAttempts(source []int) []int

手动核对：
  - 输入 [0, 2, 4]，返回 [1, 3, 5]，输入仍为 [0, 2, 4]。
  - 把结果第一个元素改成 99，输入第一个元素仍为 0。
  - 输入 [] 或 nil，结果长度为 0。
*/
func exercise2() {
	// TODO：手写复制与转换逻辑，再验证原切片和返回切片是否独立。
}

/*
练习 3：合并用户的界面设置。

背景：后端根据 PATCH 请求更新少量设置，没有提交的键保持原值。

要求：
  - 将 patch 中所有键值写入 current；已存在的键覆盖，不存在的键新增。
  - current 为 nil 时创建 map，返回最终可用的 map；调用方接收返回值。
  - current 非 nil 时原地更新，同一 map 的其他引用能看到更改。
  - patch 为 nil 时保留已有数据；即使两者都为 nil，返回值也可写入。

建议签名：
func applySettings(current, patch map[string]string) map[string]string

手动核对：
  - current={theme:light, locale:zh}，patch={theme:dark}：结果为 {theme:dark, locale:zh}。
  - current=nil，patch={locale:en}：结果包含 locale=en。
  - 两者 nil：结果长度为 0，执行 result["theme"] = "dark" 不会 panic。
  - 打印 map 时不要依赖键的顺序。
*/
func exercise3() {
	// TODO：实现 helper，并在调用方显式接收返回的 map。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func normalizePageSize(size *int) int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：normalizePageSize；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func nextAttempts(source []int) []int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：nextAttempts；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func applySettings(current, patch map[string]string) map[string]string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：applySettings；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
