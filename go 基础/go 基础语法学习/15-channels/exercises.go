// 章节：15-通道与通信
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
练习 1：通过单个结果通道交付购物车小计。

背景：页面发起一个后台计算任务，通过 channel 取得总价。

要求：
  - 主 goroutine 创建无缓冲 chan int，再启动计算 goroutine。
  - 工作函数接收 []int，求和后恰好发送一个总价，发送后不再做其他工作。
  - 工作函数接收 chan<- int 参数；调用方恰好接收一次。
  - 本题结果数量已知，不必关闭通道；不要对它用 range 等待关闭。
  - 输入切片在计算期间不再被调用方修改，且总价不会溢出。

建议签名：func sendSubtotal(prices []int, output chan<- int)

手动核对：
  - [100,250,50]：收到 400。
  - [0]：收到 0；nil 或空输入也必须发送一次 0，不能提前丢失结果。
  - 不在主 goroutine 中直接调用阻塞发送函数后再接收。
*/
func exercise1() {
	// TODO：明确“启动一个任务、发送一次、接收一次”的协议，再手写实现。
}

/*
练习 2：单生产者的导出行流。

背景：报表任务逐行产生 CSV 文本，消费方不知道共有多少有效行。
本题的名字不含逗号、双引号或换行；真实通用 CSV 应使用 encoding/csv 处理转义。

要求：
  - 启动一个生产者，遍历 names，跳过空字符串，发送 "user,"+name。
  - 生产者拥有关闭权，所有发送完成后关闭通道；空输入也要关闭。
  - 消费者通过 <-chan string 参数配合 range 收集 []string。
  - 用容量为 2 的缓冲通道；不要靠缓冲大小保证流程最终结束。
  - 生产者不再修改已发送的数据；消费者不关闭通道。

建议签名：
func exportNames(names []string, output chan<- string)
func collectRows(input <-chan string) []string

手动核对：
  - ["lin","","zoe"]：得到 ["user,lin","user,zoe"]，顺序一致。
  - 输入超过两项也能完成，说明接收与发送正确协作。
  - nil 或全空字符串：长度为 0，range 正常退出。
*/
func exercise2() {
	// TODO：让唯一生产者负责 close，消费方持续接收直到关闭。
}

/*
练习 3：聚合多项后台校验结果。

背景：管理后台同时校验多个用户名，不希望共享追加结果 slice。

要求：
  - 定义 ValidationResult{Index int; Valid bool}，每个输入启动一个任务。
  - 规则：ASCII 用户名字节长度至少 3 才有效，每个任务恰好发送一条结果。
  - 用 WaitGroup 先 Add；任务 defer Done，任务自身不关闭共享通道。
  - 单独协调 goroutine 在 Wait 后 close；主 goroutine 立刻开始 range 接收。
  - 由主 goroutine 根据 Index 写入等长 []bool，返回时保持输入顺序。
  - 空输入时协调者也要关闭通道；不准靠 Sleep 或无限增大缓冲区绕过同步。

建议签名：func collectValidation(names []string) []bool

手动核对：
  - ["lin","","zoe","ab"]：返回 [true,false,true,false]。
  - nil：结果长度为 0，程序退出，不残留等待接收的 goroutine。
  - 使用 go run -race ./15-channels exam 检查已实现的代码。
*/
func exercise3() {
	// TODO：先写清发送者、接收者、关闭者各自职责，再实现任务聚合。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func sendSubtotal(prices []int, output chan<- int) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：sendSubtotal；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func exportNames(names []string, output chan<- string) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：exportNames；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func collectRows(input <-chan string) []string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：collectRows；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type ValidationResult struct {
	Index int
	Valid bool
}

func collectValidation(names []string) []bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：collectValidation；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
