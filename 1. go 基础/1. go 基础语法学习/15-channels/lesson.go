/*
章节：15-通道与通信
	本章掌握清单：
	1. 无缓冲 channel：发送与接收配对
	2. 缓冲、关闭与接收的 ok
	3. 多个发送者与统一关闭
	4. 发送值与数据所有权
	5. worker pool：有限任务队列
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./15-channels
// 运行练习：go run ./15-channels exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：channel 可以传递一串有类型的值，并同步发送方和接收方；不是 Promise。
// 关闭表示“不会再发送”，不表示失败，也不是销毁资源或强制终止 goroutine。

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type TaskResult struct {
	TaskID string
	Count  int
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：无缓冲 channel：发送与接收配对
	// ═══════════════════════════════════════════════════════
	// make(chan string) 创建通道；ch <- value 发送，<-ch 接收。
	// 两个方向配对完成数据交付；另一个 finished 通道用于确认发送后的工作也完成。

	// 无缓冲通道没有队列，发送与接收要配对才能交付。
	result := make(chan string)
	finished := make(chan struct{}) // 只用于完成通知，不携带业务数据。
	go func() {
		result <- "profile-ready"
		// 发送之后还可以有工作；收到结果不代表本 goroutine 的所有工作已完成。
		close(finished)
	}()
	fmt.Println("无缓冲交付：", <-result)
	<-finished // 明确等待发送之后的完成信号。
	// result 不必关闭：此处已知只接收一次，没有接收者要靠 close 结束 range。
	// 若在同一 goroutine 里先 result <- x 再 <-result，会卡在发送处。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：缓冲、关闭与接收的 ok
	// ═══════════════════════════════════════════════════════
	// make(chan string, 2) 允许先排队两个值；第三次发送需要有人取走数据。
	// close 后仍能读完已有值，只有关闭且排空后才得到零值和 ok=false。

	// 缓冲区可先容纳有限个值；满时发送仍会阻塞，空时接收仍会阻塞。
	queue := make(chan string, 2)
	queue <- "compile"
	queue <- "deploy"
	fmt.Println("队列长度/容量：", len(queue), cap(queue)) // 2 2
	// 不能仅凭 len(queue) 判断下一步一定不会阻塞：有其他 goroutine 时状态会变化。
	close(queue) // 关闭后，已缓冲的数据仍然可以被读完。
	for job := range queue {
		fmt.Println("缓冲任务：", job) // compile，然后 deploy。
	}
	value, ok := <-queue
	fmt.Printf("关闭并读空后：%q, ok=%v\n", value, ok) // "", false
	// ok=false 只表示关闭且已读空，不表示“暂时没有值”。
	// 空但未关闭的通道会等待；nil 通道的发送与接收会永远等待。
	// closed 通道重复 close、向其发送，以及 close(nil) 都会 panic。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：多个发送者与统一关闭
	// ═══════════════════════════════════════════════════════
	// 每个任务只发送，协调者等全部任务结束后关闭，消费者一直读取到关闭。
	// 这三个角色共同保证不会出现“还有任务发送，通道却已关闭”。

	// 多个任务共享输出通道，各任务负责发送，协调者统一关闭。
	stream := make(chan TaskResult)
	var wait sync.WaitGroup
	wait.Add(2) // 先登记所有任务，再启动工作及等待者。
	go func() {
		defer wait.Done()
		produceResult("users", 2, stream)
	}()
	go func() {
		defer wait.Done()
		produceResult("orders", 3, stream)
	}()
	go func() {
		wait.Wait()   // 等全部发送者完成，保证关闭之后不会再发送。
		close(stream) // 只有此协调者负责关闭一次。
	}()

	// 主 goroutine 必须同时接收，不能在这里先 Wait 再读无缓冲 stream。
	// 否则发送者卡在发送处无法 Done，主 goroutine 又卡在 Wait，形成互相等待。
	counts := collectResults(stream)
	fmt.Println("用户数 / 订单数：", counts["users"], counts["orders"]) // 2 3
	// goroutine 启动顺序不保证结果到达顺序。这里按固定键打印，方便核对。
	// collectResults 读到 close 才结束，说明发送者已经全部完成。
	// 多个生产者汇入一个结果流常称为 fan-in；同一通道的多个接收者会分摊数据，
	// 每个值只交给其中一个接收者，并不会广播给所有接收者。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：发送值与数据所有权
	// ═══════════════════════════════════════════════════════
	// 通道会复制发送的值，但复制 slice 或指针不会自动复制内部数据。
	// 若发送后还继续修改共享对象，双方必须另外同步或约定移交使用权。

	// 通道发送的是值，指针、slice、map 和结构体内部引用不会被自动深拷贝。
	// 将 *TaskResult 发出去后，继续无同步地修改它仍可能和接收方产生数据竞争。
	// 常用约定是发送后移交使用权，或发送完整的独立副本。
	// 缓冲容量不是修复所有死锁的办法：双方仍需约定发送数量、结束方式和取消方式。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：worker pool：有限任务队列
	// ═══════════════════════════════════════════════════════
	// 阅读顺序：建立队列 -> 启动固定 worker -> 投递任务 -> 收集结果 -> 统一结束。
	// 结果携带原索引，所以不同任务即使乱序完成，也能还原输入顺序。

	// 用固定数量 worker 从队列持续领任务，避免每条数据都新建一个 goroutine。
	previews, valid := buildPreviews([]string{" first ", "second", " third"}, 2)
	fmt.Println("固定 worker 的预览结果：", previews, valid) // [FIRST SECOND THIRD] true
	previews, valid = buildPreviews(nil, 2)
	fmt.Println("空任务直接结束：", len(previews), valid) // 0 true
	_, valid = buildPreviews([]string{"first"}, 0)
	fmt.Println("拒绝无效 worker 数：", valid) // false
}

type previewJob struct {
	index int
	text  string
}

type previewResult struct {
	index int
	text  string
}

func buildPreviews(inputs []string, maxWorkers int) ([]string, bool) {
	if maxWorkers < 1 {
		return nil, false
	}
	previews := make([]string, len(inputs))
	if len(inputs) == 0 {
		return previews, true
	}
	jobs := make(chan previewJob, 1)
	results := make(chan previewResult, 1)
	var wait sync.WaitGroup
	for worker := 0; worker < min(maxWorkers, len(inputs)); worker++ {
		wait.Go(func() {
			for job := range jobs {
				text := strings.ToUpper(strings.TrimSpace(job.text))
				results <- previewResult{index: job.index, text: text}
			}
			// worker 不是唯一发送者，不能自行 close(results)。
		})
	}
	go func() {
		defer close(jobs) // 唯一生产者拥有 jobs 的关闭权。
		for index, text := range inputs {
			jobs <- previewJob{index: index, text: text}
		}
	}()
	go func() {
		wait.Wait()
		close(results) // worker 全结束之后统一关闭结果流。
	}()
	for result := range results {
		previews[result.index] = result.text // 只由当前接收者写入，不并发 append。
	}
	return previews, true
	// 关闭顺序：生产者完成 -> 关闭 jobs -> worker 排空 jobs 并退出 -> 关闭 results。
	// 顺序完整后，消费者才结束 range；每个任务的结果用索引恢复输入顺序。
}

// 队列容量有限时，消费变慢会让 results 填满，worker 等待，接着 jobs 填满、生产者等待。
// 这种向上游传递的等待称为“背压”，它能限制排队工作，而不是让内存无限增长。
// 本例协议要求消费方读到 results 关闭，不能中途直接返回，否则发送者可能永远挂起。
// 如业务允许提前结束，需要第 16 章的取消信号，让投递 jobs、领取任务和发送 results
// 等每个可能阻塞的位置都可退出；不要通过接收方擅自 close 来替代取消协议。

// chan<- T：函数只允许向通道发送，不能读取。
// 它仍然能被 close，因此“谁有权关闭”除了类型，还需通过明确协议来约束。
func produceResult(taskID string, count int, output chan<- TaskResult) {
	output <- TaskResult{TaskID: taskID, Count: count}
	// 不关闭：这个通道还有其他发送者，关闭权属于外部协调者。
}

// <-chan T：函数只能接收，不能发送，也不能 close。
// 双向 chan T 可以传给单向参数，而只收/只发通道不能随意变回双向通道。
func collectResults(input <-chan TaskResult) map[string]int {
	counts := make(map[string]int)
	for result := range input {
		counts[result.TaskID] = result.Count
	}
	return counts
}

// 后续再学：提前取消、超时与多路等待（16 章）、共享对象的同步（17 章）；
// worker pool 是限制资源的一种模式，小任务不必为了并发而额外制造队列。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 queue 容量改成 1，解释为什么原来的第二次发送会阻塞，再恢复容量。
// 沿 buildPreviews 找到 jobs 和 results 各自唯一的关闭者，并说明关闭先后顺序。
