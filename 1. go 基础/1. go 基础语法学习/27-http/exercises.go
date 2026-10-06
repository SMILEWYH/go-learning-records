// 章节：27-HTTP网络编程
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 练习 1：编写分页查询接口。
// 要求：GET /products?page=1&size=20，两个参数缺省为 1 和 20；page>=1，size 为 1..100。
// 成功返回 application/json 和 {"page":1,"size":20}；非法输入返回 400。
// 常规：?page=2&size=10 -> 200；无查询参数 -> 默认值。
// 边界：page=abc、page=0、size=101 -> 400；POST 同路径 -> 405。
// 不用任何真实产品数据库；用 httptest.NewRecorder/NewRequest 检查 handler。
func exercise1() {
	// TODO：解析查询字符串、验证边界、输出 JSON。
}

// 练习 2：实现局部修改昵称的接口。
// 要求：PATCH /users/{id}，id 为正整数；正文 {"name":"小周"}；成功返回 id/name，状态 200。
// 校验 application/json；限制 1 KiB；拒绝空名称、未知字段和第二个 JSON 值。
// 边界：id=abc、id=0、空名称 -> 400；错误媒体类型 -> 415；正文超限 -> 413。
// 示例不保存状态；用 NewRequest+Do 发请求，响应 Body 必须关闭。
func exercise2() {
	// TODO：构造自己的 mux 和 httptest.Server，同时核对状态码和返回值。
}

// 练习 3：让客户端识别服务端失败。
// 建议：func fetchMessage(client *http.Client, url string) (string, error)。
// 要求：只接受 200；其他状态返回包含状态码的 error；任何成功拿到的响应都关闭 Body。
// 常规：测试服务返回 200 + "ready" -> ready,nil；返回 404/500 -> error。
// 边界：使用已取消的 context 创建请求并通过 Do 发送，errors.Is(err, context.Canceled) 为 true。
// 客户端设置 Timeout；不要用外部网站或真实休眠很久的接口制造失败。
func exercise3() {
	// TODO：使用本机测试服务器覆盖 HTTP 失败和请求取消两个不同层次。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func productHandler(w http.ResponseWriter, r *http.Request) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：productHandler；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func newUserPatchMux() *http.ServeMux {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：newUserPatchMux；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func fetchMessage(client *http.Client, url string) (string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：fetchMessage；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
