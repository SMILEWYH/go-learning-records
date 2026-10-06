// 章节：16-多路选择与超时取消
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
练习 1：非阻塞的后台通知队列。

背景：可丢弃的界面刷新通知不应该让业务请求一直等待队列空位。

要求：
  - 用 select 的发送 case 与 default，实现可以立刻入队则返回 true。
  - 队列已满或队列为 nil 时立即返回 false；不启动额外 goroutine。
  - 不使用 len(queue)<cap(queue) 作为能否发送的判断。
  - 本题约定队列不会被关闭；向已关闭通道发送仍会 panic，default 无法兜底。
  - 在练习函数创建容量 1 的队列，演示入队、满载拒绝、接收后再次入队。

建议签名：func tryEnqueue(queue chan<- string, message string) bool

手动核对：
  - 空队列写入 A：true；再写 B：false；接收得到 A。
  - 接收后再次写 B：true，接收得到 B。
  - nil 队列：false；不要用这个例子推断对关闭通道也安全。
*/
func exercise1() {
	// TODO：只使用一次 select，不靠额外缓冲、Sleep 或后台任务逃避阻塞。
}

/*
练习 2：合并两路审计事件并正常退出。

背景：登录事件流与配置变更流分别结束，消费方需持续处理仍然存活的那一路。

要求：
  - 输入两个只读 string 通道，用 select 收集所有事件。
  - 每条流接收时检查 ok，关闭且读空后把对应局部通道变量设为 nil。
  - 两条流都结束后返回结果；初始即 nil 的流视为已结束。
  - 消费者不关闭输入，生产者各自负责关闭；无须保证两条流之间的总顺序。
  - 在练习中用足够容量的缓冲通道填充固定数据后关闭，无需后台生产者。

建议签名：func mergeEvents(login, changes <-chan string) []string

手动核对：
  - 登录流 [login-A,login-B]，变更流 [change-X]：三条恰好各出现一次。
  - 同一条流内部顺序保持 login-A 在 login-B 前，跨流交错允许不同。
  - 一路 nil、另一路有两条并关闭：得到那两条；两路 nil：长度 0。
  - 关闭但没有数据的通道：不收集零值，不忙循环。
*/
func exercise2() {
	// TODO：写完正常输入后重点检查“一路提前关闭”和“两路都结束”。
}

/*
练习 3：超时后能够真正收尾的预览任务。

背景：用户离开预览页面后，后台不应永远卡在无人接收的结果发送上。

要求：
  - 增加 context、time 导入，定义工作函数与完成信号 stopped。
  - 工作函数创建 time.NewTimer(delay)，defer Stop；先 select 等 timer 或取消。
  - 定时器触发后，再用 select 在“发送 preview-ready”和“ctx.Done”之间选择。
  - 所有退出路径均通过 defer close(stopped) 报告完成，只有工作函数关闭它。
  - 调用方拥有 cancel，defer cancel；工作函数不关闭调用方的结果通道。
  - 任何成功、取消、超时路径，调用方最后都等待 stopped，确认工作函数已退出。

建议签名：
func previewWorker(ctx context.Context, delay time.Duration, output chan<- string, stopped chan<- struct{})

手动核对：
  - 成功：Background 派生可取消 ctx，delay=0；接收 preview-ready 后等 stopped。
  - 取消：在启动前主动 cancel；无接收者也应结束，等 stopped 后不读结果。
  - 超时：WithTimeout 设约 10ms，delay=0，但不接收 output；等 ctx.Done，再等 stopped。
  - 超时场景不接收结果，仍必须退出；不要依赖精确毫秒数或 select 分支优先级。
  - 用自己的话解释：为什么仅在调用方写 time.After 不能保证工作函数退出？
*/
func exercise3() {
	// TODO：为“等待准备”和“发送结果”两个阻塞位置都加上取消分支。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func tryEnqueue(queue chan<- string, message string) bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：tryEnqueue；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func mergeEvents(login, changes <-chan string) []string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：mergeEvents；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func previewWorker(ctx context.Context, delay time.Duration, output chan<- string, stopped chan<- struct{}) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：previewWorker；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
