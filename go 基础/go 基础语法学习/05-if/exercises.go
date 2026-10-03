// 章节：05-条件判断
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

// 练习 1：规范化接口分页大小。
// 背景：前端可能不传 pageSize，也可能传来过大的数值。
// 要求：用 if / else if / else 把 requestedSize 转为有效的 pageSize。
//   - requestedSize <= 0：采用默认值 20。
//   - requestedSize > 100：限制为 100。
//   - 其余数值：保持原值。
//
// 常规：requestedSize=35 -> 35。
// 边界：0 -> 20；-1 -> 20；100 -> 100；101 -> 100。
// 建议实现位置：exercise1()，修改输入后多次运行，不要求写循环。
// 学完函数后可抽取 normalizePageSize(requestedSize int) int。
// 提示：不要让有效值 100 被误判为非法，也不要照搬 TS 的 requestedSize || 20。
func exercise1() {
	// TODO：根据三条规则计算并打印 pageSize。
}

// 练习 2：做一次明确的权限判断。
// 背景：后台页面是否展示按钮不能替代后端权限校验。
// 要求：输入 loggedIn bool、role string、isOwner bool，输出一个状态码。
//   - 未登录，无论其他条件如何，返回状态码 401。
//   - 已登录，且 role == "admin" 或 isOwner == true，返回 200。
//   - 其他已登录用户，返回 403。
//
// 常规：true、"editor"、true -> 200；true、"viewer"、false -> 403。
// 边界：false、"admin"、true -> 401；true、"admin"、false -> 200。
// 建议实现位置：exercise2()；使用提前 return 时，只结束本题函数。
// 学完函数后可抽取 accessStatus(loggedIn bool, role string, isOwner bool) int。
// 提示：优先判断是否登录；逻辑组合时用括号明确自己的意图。
func exercise2() {
	// TODO：声明输入，根据规则输出 401、403 或 200。
}

// 练习 3：识别配置项的三种状态。
// 背景：自定义页面标题可能“未配置”“明确配置为空”或“配置了文本”。
// 要求：使用 map[string]string 和 comma-ok 查询 "title"，再用 if 判断。
//   - key 不存在：显示默认标题 "Dashboard"。
//   - key 存在，值为空：显示 "(hidden)"，表示用户选择隐藏标题。
//   - key 存在且非空：显示配置值。
//
// 常规：map[string]string{"title":"Orders"} -> "Orders"。
// 边界：空 map -> "Dashboard"；{"title":""} -> "(hidden)"。
// 额外边界：nil map -> "Dashboard"；本题只读 map，不向 nil map 写入。
// 建议实现位置：exercise3()，可以把 value, ok 写在 if 的初始化语句中。
// 学完函数后可抽取 resolveTitle(config map[string]string) string。
func exercise3() {
	// TODO：读取 title，根据存在性与文本内容选择最终标题。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func normalizePageSize(requestedSize int) int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：normalizePageSize；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func accessStatus(loggedIn bool, role string, isOwner bool) int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：accessStatus；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func resolveTitle(config map[string]string) string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：resolveTitle；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
