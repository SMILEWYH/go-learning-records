/*
章节:15-通道与通信
*/

// goroutine 用来并发执行任务, channel 用来在任务之间传递数据、协调执行
// close(channel) 表示“不会再发送”, 不表示失败、销毁资源或强制终止 goroutine

// 本章节终极奥义!
// 在 channel 未关闭的前提下，无缓冲 channel 遵循的规律:
// 每次发送需要有接收方, 每次接收需要有发送方, 否则当前 goroutine 会阻塞
// 在 channel 未关闭的前提下，有缓冲 channel 遵循的规律:
// 每次发送需要有等待的接收方, 或者缓冲区空位, 否则阻塞
// 每次接收需要缓冲区有数据, 或者有发送方提供数据, 否则阻塞
package main

import (
	"fmt"
	"strings"
	"sync"
)

type TaskResult struct {
	TaskID string
	Count  int
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节: 无缓冲 channel 的发送与接收
	// ═══════════════════════════════════════════════════════
	// 无缓冲通道没有数据缓冲区, 每次发送都要有接收一一配对完成, 不然其中一方就会阻塞
	// 发送时,如果没有接收方准备好,就会阻塞在发送语句处
	// 接收时,如果没有发送方准备好,就会阻塞在接收语句处
	// 配对完成后,双方都可以继续执行,但后续代码的执行顺序不确定
	// 无缓冲 channel 可以反复发送和接收消息。“无缓冲”只是表示它不能暂存数据,每次发送都需要有对应的接收
	result := make(chan string) // 创建只能传 string 类型的无缓冲通道
	// var result chan string // 只是声明, chanel 零值为 nil, 该通道目前还不可用
	var resultWait sync.WaitGroup
	resultWait.Go(func() {
		result <- "profile-ready-1" // 把字符串推入到 channel
		fmt.Println("后续业务逻辑")
	})
	resultWait.Go(func() {
		result <- "profile-ready-2" // 把字符串推入到 channel
		fmt.Println("后续业务逻辑")
	})
	answer_1, ok := <-result                  // 获取通道的字符串
	answer_2, ok := <-result                  // 获取通道的字符串
	fmt.Println("无缓冲交付：", answer_1, answer_2) // 注意这里接受消息的顺序不一定
	resultWait.Wait()                         // 先接收再等待, 否则后台会阻塞在发送处, 无法完成任务
	// 注意: 这里的 Wait() 是必须的, 因为 goroutine 中往 channel 推入消息后还有其他业务逻辑执行需要等待, 不然 main goroutine 结束了, 会导致其他 goroutine 任务丢失
	// 注意: 无缓冲通道保证“发送和接收完成配对”, 但不保证接收方先打印

	// ═══════════════════════════════════════════════════════
	// 第 2 节: 缓冲 channel 的发送与接收、关闭、接收的 ok
	// ═══════════════════════════════════════════════════════
	queue := make(chan string, 2) // 创建只能传 string 类型、缓冲区长度为2的通道
	queue <- "compile"
	queue <- "deploy"
	fmt.Println("队列长度/容量：", len(queue), cap(queue)) // 2 2
	// 关闭后, 已缓冲的数据仍然可以被读完
	close(queue)
	for job := range queue {
		fmt.Println("缓冲任务：", job) // compile、deploy
	}
	// channel 缓冲区消息都已被取完
	value, ok := <-queue
	fmt.Printf("关闭并读空后：%q, ok=%v\n", value, ok) // ""、false; 如果在之前没有 for 循环输出, 那么这时候的值为 compile、true
	// ok=false 只表示关闭且已读空, 不表示“暂时没有值”
	// 空但未关闭的通道会等待; nil 通道的发送与接收会永远等待

	// ═══════════════════════════════════════════════════════
	// 第 3 节: 单个(缓冲或无缓冲) channel 的多个发送者与统一关闭
	// ═══════════════════════════════════════════════════════
	// 负责执行业务逻辑和往 channel 推入消息
	produceResult := func(taskID string, count int, output chan<- TaskResult) {
		output <- TaskResult{TaskID: taskID, Count: count}
	}

	// 负责 renge 循环读取 channel 中的消息, 直到 channel close 才会结束循环
	collectResults := func(input <-chan TaskResult) map[string]int {
		counts := make(map[string]int)
		for result := range input {
			counts[result.TaskID] = result.Count
		}
		return counts
	}

	// 这里使用无缓冲 channel 或缓冲 channel 基本都一样, 主要区别在于发送者是否需要等待接收
	stream := make(chan TaskResult)
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer wait.Done()
		// 业务逻辑
		produceResult("users", 2, stream) // channel 直接传递即可, 不需要改为指针类型
		// 其他业务逻辑
	}()
	go func() {
		defer wait.Done()
		// 业务逻辑
		produceResult("orders", 3, stream)
		// 其他业务逻辑
	}()
	// 该 goroutine 不可以使用 wait.Go()
	// 因为会形成死锁, 自己等待 wait.Wait() , 而 wait.Wait() 又要等待这个 wait.Go() 本身执行结束
	go func() { // 这个 goroutine 被称为 协调者
		wait.Wait()
		close(stream)
		// 等待全部发送者完成后,关闭 channel
	}()

	counts := collectResults(stream)
	fmt.Println("用户数 / 订单数：", counts["users"], counts["orders"]) // 2 3

	// ═══════════════════════════════════════════════════════
	// 第 4 节:发送值与数据所有权
	// ═══════════════════════════════════════════════════════
	// 通道发送的是值, 指针、slice、map 和结构体内部引用不会被自动深拷贝
	// 将 *TaskResult 发出去后, 继续无同步地修改它仍可能和接收方产生数据竞争
	// 常用约定是发送后移交使用权, 或发送完整的独立副本
	// 用缓冲通道在同一个 goroutine 内演示复制, 不引入并发读写
	shared := []int{1, 2}
	values := make(chan []int, 1)
	values <- shared
	shared[0] = 9
	fmt.Println("发送切片仍共享底层数组：", <-values) // [9 2]

	independent := append([]int(nil), shared...) // 复制本例的 int 元素
	values <- independent
	shared[0] = 100
	fmt.Println("发送独立副本：", <-values) // [9 2]
	// 缓冲容量不是修复所有死锁的办法:双方仍需约定发送数量、结束方式和取消方式

	// ═══════════════════════════════════════════════════════
	// 第 5 节:worker pool: 有限任务队列
	// ═══════════════════════════════════════════════════════
	// 用固定数量 worker 从队列持续领任务, 避免每条数据都新建一个 goroutine
	type previewJob struct {
		index int
		text  string
	}

	type previewResult struct {
		index int
		text  string
	}

	buildPreviews := func(inputs []string, maxWorkers int) ([]string, bool) {
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
				// worker 不是唯一发送者, 不能自行 close(results)
			})
		}
		go func() {
			defer close(jobs) // 唯一生产者拥有 jobs 的关闭权
			for index, text := range inputs {
				jobs <- previewJob{index: index, text: text}
			}
		}()
		go func() {
			wait.Wait()
			close(results) // worker 全结束之后统一关闭结果流
		}()
		for result := range results {
			previews[result.index] = result.text // 只由当前接收者写入, 不并发 append
		}
		return previews, true
		// 关闭顺序:生产者完成 -> 关闭 jobs -> worker 排空 jobs 并退出 -> 关闭 results
		// 顺序完整后, 消费者才结束 range;每个任务的结果用索引恢复输入顺序
	}

	previews, valid := buildPreviews([]string{" first ", "second", " third"}, 2)
	fmt.Println("固定 worker 的预览结果：", previews, valid) // [FIRST SECOND THIRD] true
	previews, valid = buildPreviews(nil, 2)
	fmt.Println("空任务直接结束：", len(previews), valid) // 0 true
	_, valid = buildPreviews([]string{"first"}, 0)
	fmt.Println("拒绝无效 worker 数：", valid) // false

	// 队列容量有限时, 消费变慢会让 results 填满, worker 等待, 接着 jobs 填满、生产者等待
	// 这种向上游传递的等待称为“背压”, 它能限制排队工作, 而不是让内存无限增长
	// 本例协议要求消费方读到 results 关闭, 不能中途直接返回, 否则发送者可能永远挂起
	// 如业务允许提前结束, 需要第 16 章的取消信号, 让投递 jobs、领取任务和发送 results
	// 等每个可能阻塞的位置都可退出;不要通过接收方擅自 close 来替代取消协议

	// worker pool 是限制资源的一种模式, 小任务不必为了并发而额外制造队列
}
