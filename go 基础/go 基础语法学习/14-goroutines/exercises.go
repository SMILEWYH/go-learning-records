// 章节：14-协程与任务等待
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
练习 1：并发拼装仪表盘摘要。

背景：页面需要用户名称和未读消息数，它们可以由两个独立任务获取。
本题用固定数据模拟服务响应，不访问网络，也不使用 Sleep 模拟等待。

要求：
  - 定义 DashboardSummary{UserName string; Unread int}。
  - 启动两个 goroutine：一个填写 UserName，一个填写 Unread。
  - 增加 sync 导入，用 WaitGroup 在启动前 Add(2)，各任务 defer Done()。
  - 每个任务只写自己的字段，主 goroutine 只在 Wait 后读取整个结果。
  - Wait 返回后才返回摘要；不在后台任务里打印以免输出顺序不稳定。

建议签名：func buildDashboard(name string, unread int) DashboardSummary

手动核对：
  - ("Lin",3)：返回 UserName=Lin、Unread=3。
  - ("",0)：返回字段零值，两个任务仍须完成，不挂起。
  - 连续调用十次，结果一致；可以用 -race 运行检查本次执行。
*/
func exercise1() {
	// TODO：只使用两个独立写入位置；不要添加共享的“已完成数量”变量。
}

/*
练习 2：保持输入顺序的批量价格计算。

背景：批量导入商品时，可以分别计算每个商品加固定费用后的价格。

要求：
  - 为每个输入价格启动一个 goroutine，结果为 price+fee。
  - 输入不修改，预先创建等长结果 slice，每个任务只写自己的索引。
  - 不并发 append，不在任务里修改同一个累加变量。
  - 调用前登记任务数量，全部 Done 后返回结果。
  - 本题输入规模很小且数值不会溢出；真实大批量任务需限制并发数。

建议签名：func addFees(prices []int, fee int) []int

手动核对：
  - [100,250,0]、fee=20：返回 [120,270,20]，顺序不变。
  - [100]、fee=0：返回 [100]。
  - nil 或 []：返回长度 0 的结果，Wait 不阻塞。
*/
func exercise2() {
	// TODO：用参数向匿名 goroutine 显式传递索引和价格。
}

/*
练习 3：批量校验注册资料，提前返回也报告完成。

背景：校验任务遇到空用户名时会直接拒绝，但整个批次仍需等待全部任务。

要求：
  - 输入 []string，每项一个 goroutine，输出等长 []bool。
  - 用户名非空且字节长度至少 3 为 true，否则 false；本题输入限定 ASCII。
  - 空字符串分支必须使用提前 return，仍保证执行 defer Done()。
  - 每项只写自己结果槽，不共享追加错误列表。
  - 返回后在主 goroutine 顺序统计有效用户数。

建议签名：func validateUsernames(names []string) []bool

手动核对：
  - ["lin","","ab","zoe"]：返回 [true,false,false,true]，有效数 2。
  - 全为空字符串：全部 false，程序正常结束。
  - nil：长度为 0，有效数 0。
  - 使用 go run -race ./14-goroutines exam，确认实现未发现数据竞争。
*/
func exercise3() {
	// TODO：先定义每个任务的 Done 协议，再写各条校验分支。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type DashboardSummary struct {
	UserName string
	Unread   int
}

func buildDashboard(name string, unread int) DashboardSummary {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：buildDashboard；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func addFees(prices []int, fee int) []int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：addFees；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func validateUsernames(names []string) []bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：validateUsernames；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
