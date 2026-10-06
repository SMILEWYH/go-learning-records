// 章节：28-RPC远程调用
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 练习 1：实现远程运费计算。
// 要求：定义 ShippingService.Calculate，参数含 SubtotalCents int64、Express bool，结果含 FeeCents int64。
// 普通运费 600 分，满 10000 分免普通运费；加急额外 1000 分；负金额返回 error。
// 常规：(9999,false)->600；(10000,false)->0；(10000,true)->1000。
// 用独立 rpc.Server 注册，通过真实本机客户端调用；第二个参数必须是结果指针。
func exercise1() {
	// TODO：定义导出参数类型与方法，使用空闲端口完成调用后清理连接。
}

// 练习 2：区分本地连接失败和远程业务失败。
// 要求：调用 Quote(Quantity=0)，用 errors.As 检查 rpc.ServerError；调用不存在的方法也应报错。
// 再关闭 client，尝试一次新调用，确认失败且不会把之前成功的 reply 当成新结果。
// 每次调用使用独立结果变量；写注释说明服务端的 error 类型为什么不能完整跨网络保留。
func exercise2() {
	// TODO：逐次打印错误类别，并断言正常报价与失败报价的行为。
}

// 练习 3：批量异步报价。
// 要求：同一个 client 用 Go 发起三个请求，每个请求有自己的 reply 和容量至少为 1 的 Done。
// 单价 500，数量分别 1、2、3，完成后按输入顺序输出 [500,1000,1500]。
// 边界：将其中一个数量改成 0，收集该项 error，其余成功结果仍保留。
// 不假定网络完成顺序；不在 Done 前读取 reply；为底层连接设置 deadline 并最终 Close。
func exercise3() {
	// TODO：将请求、reply、call 一一对应，并运行 race detector。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type ShippingRequest struct {
	SubtotalCents int64
	Express       bool
}
type ShippingResponse struct{ FeeCents int64 }
type ShippingService struct{}

func (*ShippingService) Calculate(args ShippingRequest, reply *ShippingResponse) error {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Calculate；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func batchQuotes(client *rpc.Client, quantities []int) ([]QuoteResponse, []error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：batchQuotes；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
