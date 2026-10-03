/*
章节：16-多路选择与超时取消
	本章掌握清单：
	1. select 与 default：等待还是立即返回
	2. 关闭的通道与 nil：结束一路输入
	3. time.After：一次等待的超时
	4. Timer：整个操作共用总时限
	5. context：发出取消并等待退出
	6. 父子 context 的期限与传播
	7. Ticker：周期节拍与停止
	8. select 求值时机（易错点）
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./16-select-timeout
// 运行练习：go run ./16-select-timeout exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：select 类似在多项可完成操作中选择一项；超时像 Promise.race 的超时分支，
// 仅停止当前等待，不会自动终止另一个任务。取消机制更接近 AbortSignal 的协作方式。

package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：select 与 default：等待还是立即返回
	// ═══════════════════════════════════════════════════════
	// select 选择一个当前能执行的通道操作；没有可执行项时通常阻塞。
	// default 将等待改成“现在不能完成就走这里”，适合尝试操作。

	// 没有 case 可以执行且存在 default 时，立即执行 default。
	empty := make(chan string)
	select {
	case message := <-empty:
		fmt.Println(message)
	default:
		fmt.Println("当前没有待处理消息")
	}
	// 去掉 default 就会等待。for + select + default 可能形成消耗 CPU 的忙循环。
	// 多个 case 同时就绪时会伪随机选择一个，case 的书写顺序不构成优先级。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：关闭的通道与 nil：结束一路输入
	// ═══════════════════════════════════════════════════════
	// 关闭通道会持续可读，因此不能忽略 ok=false，否则会反复收到零值。
	// 把变量设为 nil 禁用该路，外层条件在两路都结束后退出。

	// 关闭且读空的通道会一直可读，必须识别 ok=false。
	left, right := make(chan int, 1), make(chan int, 1)
	left <- 10
	right <- 20
	close(left)
	close(right)
	fmt.Println("两个流的总和：", mergeSum(left, right)) // 30，接收先后不固定。
	// 把通道变量设为 nil 可禁用对应 select case，但不会关闭原通道。
	// 单纯写 break 只会退出 select，不能退出外层 for；可用 return 或带标签的 break。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：time.After：一次等待的超时
	// ═══════════════════════════════════════════════════════
	// time.After 返回通道，计时到期后接收操作就绪。
	// 超时分支只决定当前等待何时结束，不会替你关闭 empty 或终止其他任务。

	// 一次性等待可用 time.After。此处没有发送者，所以必然走超时分支。
	// 示例没有为不存在的结果启动 goroutine，因此超时后没有遗留任务。
	select {
	case <-empty:
		fmt.Println("收到消息")
	case <-time.After(3 * time.Millisecond):
		fmt.Println("一次等待已超时")
	}
	// 时间表示最早可触发的期限，实际调度可能更晚，不能用它保证精确执行时间。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Timer：整个操作共用总时限
	// ═══════════════════════════════════════════════════════
	// 在循环外创建 timer，表示预算从开始一直累计。
	// 如果每轮都新建 time.After，收到数据就进入新一轮等待，含义会变成每轮时限。

	// 整个循环共用一个截止时间：timer 建在循环外，避免每次收到值都重置时限。
	// nil 输入会禁用接收分支，因此这次演示一定由 timer 结束。
	values, timedOut := collectUntilDeadline(nil, 3*time.Millisecond)
	fmt.Println("总时限结束，结果数/是否超时：", len(values), timedOut) // 0 true
	// time.After 也能在循环外创建后复用；不是“for 必须使用 NewTimer”。
	// NewTimer 的好处是拿到对象，可显式 Stop；需要 Reset 时要理解版本对应的语义。
	// 本例不复用或 Reset timer，不涉及旧版 Go 定时器通道的排空差异。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：context：发出取消并等待退出
	// ═══════════════════════════════════════════════════════
	// ctx.Done() 是取消信号，stopped 是任务结束信号；两者表达不同阶段。
	// 任务在发送和取消之间 select，所以下游不接收时也能结束。

	// 用 context 明确通知任务取消，并等待任务真正退出。
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel() // 及时释放关联资源；即使期限已到，重复 cancel 也安全。
	output := make(chan string)
	stopped := make(chan struct{})
	go deliverWhenPossible(ctx, output, stopped)
	// 故意不接收 output：模拟下游已不再消费。发送者必须能通过 ctx.Done 退出。
	<-ctx.Done()
	<-stopped                          // 取消信号与任务完成不是同一件事，明确等到完成通知。
	fmt.Println("任务已响应取消：", ctx.Err()) // context deadline exceeded
	// WithTimeout 到期会关闭 Done 通道；Err 可区分取消和截止时间到达。
	// 它不会强杀 goroutine，也不能打断一个完全不检查 ctx 的 CPU 循环或阻塞调用。
	// 真实请求应把 ctx 传给支持 context 的 I/O API，并让所有可能阻塞的通道操作可取消。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：父子 context 的期限与传播
	// ═══════════════════════════════════════════════════════
	// 子任务可以更早结束，却不能超过父任务的截止时间。
	// 取消子任务不取消父任务，取消父任务会通知关联的子任务。

	demonstrateContextTree()

	// ═══════════════════════════════════════════════════════
	// 第 7 节：Ticker：周期节拍与停止
	// ═══════════════════════════════════════════════════════
	// ticker.C 提供周期事件，循环负责决定处理多少次或何时退出。
	// Stop 停止计时但不关闭通道，因此还需要显式结束条件。

	demonstrateTicker()

	// ═══════════════════════════════════════════════════════
	// 第 8 节：select 求值时机（易错点）
	// ═══════════════════════════════════════════════════════
	// 选中哪个 case 之前，发送的值表达式就已经执行。
	// 本例 output 为 nil，发送永远不就绪，buildPayload 却仍然调用一次。

	demonstrateSelectEvaluation()
}

func demonstrateContextTree() {
	// Background 用作入口的根；真实 HTTP 业务通常从请求的 r.Context() 继续派生。
	// 服务函数优先接收首个参数 ctx context.Context，继续往数据库、请求和任务传递。
	// 在内部重新使用 Background 会切断上游取消；也不要把 ctx 当普通可选配置袋。
	parent, cancelParent := context.WithTimeout(context.Background(), time.Hour)
	defer cancelParent()
	child, cancelChild := context.WithTimeout(parent, 2*time.Hour)
	defer cancelChild()
	parentDeadline, _ := parent.Deadline()
	childDeadline, _ := child.Deadline()
	fmt.Println("子任务不能延长父期限：", childDeadline.Equal(parentDeadline)) // true

	cancelChild()
	<-child.Done()
	fmt.Println("取消子任务 / 父任务仍有效：", child.Err() == context.Canceled, parent.Err() == nil)
	// 子取消不反向取消父；父取消则会传播给所有仍关联的子任务。
	sibling, cancelSibling := context.WithCancel(parent)
	defer cancelSibling()
	cancelParent()
	<-sibling.Done()
	fmt.Println("父取消传播给子任务：", sibling.Err() == context.Canceled) // true
	// 主动 cancel 的 Err 为 context.Canceled；期限先到达则为 DeadlineExceeded。
	// CancelFunc 可以重复调用，但并不等待任务退出，仍需完成信号或 WaitGroup。
	// 同一个 context 可以安全传给多个 goroutine；里面附带的数据仍需自己保证访问安全。
}

func demonstrateTicker() {
	// Timer 只触发一次；Ticker 按周期提供节拍，适合周期检查，不是精确的任务调度器。
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for count := 1; count <= 2; count++ {
		<-ticker.C
		fmt.Println("周期检查次数：", count)
	}
	// 本例用固定次数结束；后台周期任务应 select 等待 ticker.C 或 ctx.Done。
	// 消费慢时 ticker 可能调整或丢弃节拍，不能假设每次间隔都精确到指定毫秒。
	// Stop 不会关闭 C，所以不能靠 Stop 让 for range ticker.C 自然退出。
	// NewTicker 的间隔必须大于 0；NewTimer 的非正时长则会使定时器立即到期。
}

func demonstrateSelectEvaluation() {
	var output chan int // nil，发送 case 永远不会就绪。
	prepared := 0
	buildPayload := func() int {
		prepared++
		return 42
	}
	select {
	case output <- buildPayload():
		fmt.Println("发送成功")
	default:
		fmt.Println("选择 default，但准备函数已执行：", prepared) // 1
	}
	// 进入 select 时，所有 case 的通道表达式及发送值表达式会按源码顺序求值一次，
	// 然后才选择可执行分支；未选中的 case 也可能已经产生求值副作用。
	// 不要把不可取消的慢计算藏在发送值表达式里，指望另一个取消 case 把它打断。
	// 接收赋值左侧的复杂表达式则只在对应 case 被选中后求值。
}

func mergeSum(left, right <-chan int) int {
	total := 0
	for left != nil || right != nil {
		select {
		case value, ok := <-left:
			if !ok {
				left = nil // 禁用已完成的流，避免不断收到零值。
				continue
			}
			total += value
		case value, ok := <-right:
			if !ok {
				right = nil
				continue
			}
			total += value
		}
	}
	// 所有 case 都为 nil 且没有 default 时会永远阻塞，所以外层条件必须能结束。
	return total
}

func collectUntilDeadline(input <-chan int, budget time.Duration) ([]int, bool) {
	timer := time.NewTimer(budget)
	defer timer.Stop()
	var values []int
	for {
		select {
		case value, ok := <-input:
			if !ok {
				return values, false // 流正常关闭，无需等到超时。
			}
			values = append(values, value)
		case <-timer.C:
			return values, true
		}
	}
	// 调用方若会提前退出，还需负责通知生产者停止；单独接收超时无法替它取消生产。
}

func deliverWhenPossible(ctx context.Context, output chan<- string, stopped chan<- struct{}) {
	defer close(stopped) // 此任务是 stopped 的唯一关闭者。
	select {
	case output <- "preview-ready":
		return
	case <-ctx.Done():
		return
	}
	// 若发送和取消同时就绪，任一分支都可能执行；不要假定取消天然有最高优先级。
}

// 后续再学：在 HTTP/数据库 API 中传递请求 context、统一收集并发错误；
// 复杂任务树和重试策略应先明确所有权与取消协议，不靠额外 Sleep 保证正确。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 比较“每轮等 3ms”和“整个收集过程最多等 3ms”，指出 timer 应创建在哪里。
// 找出取消通知与任务完成通知两处接收，说明为什么只等前者还不够。
