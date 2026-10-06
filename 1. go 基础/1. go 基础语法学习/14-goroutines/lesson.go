/*
章节：14-协程与任务等待
	本章掌握清单：
	1. go 启动任务，WaitGroup 等待完成
	2. 并发结果按输入索引保存
	3. 独立实例与共享数据边界
	4. WaitGroup.Go 与分批复用
	5. 限制并发数并收集错误
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./14-goroutines
// 运行练习：go run ./14-goroutines exam
// 阅读顺序：PageResult 保存任务结果；main 展示用法，下方 buildPages 和 parseVersionsBounded 展示分工。
// TS 对照：go f() 可以并发执行 f；它不会返回类似 Promise 的结果对象。
// goroutine 由 Go 运行时调度，并发不代表所有任务一定在同一时刻并行运行。

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type PageResult struct {
	Title string
	Slug  string
	Valid bool
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：go 启动任务，WaitGroup 等待完成
	// ═══════════════════════════════════════════════════════
	// go func(){...}() 中最后的 () 是调用；go 让这次调用在新的 goroutine 中运行。
	// 先登记两个任务，再启动它们，最后 Wait，主流程才能读到完成后的结果。

	// 主 goroutine 结束时进程退出，不会自动等待后台任务。
	// 用 Sleep 猜任务何时完成不可靠；忙循环轮询一个共享 bool 还可能产生数据竞争。
	// WaitGroup 的零值可用，开始使用后不要复制它；传递时使用指针。
	var wait sync.WaitGroup
	var profile string
	var unreadCount int

	wait.Add(2) // 在启动任务前登记；尤其不要在 go 函数内部才 Add。
	go func() {
		defer wait.Done() // 即使提前 return，当前任务也会报告一次完成。
		profile = "Lin"
	}()
	go func() {
		defer wait.Done()
		unreadCount = 3
	}()
	wait.Wait() // 只阻塞当前 goroutine；其他 goroutine 可以继续运行。
	fmt.Printf("页面聚合：用户=%s，未读=%d\n", profile, unreadCount)
	// 这里每个变量只有一个写入者，主 goroutine 等 Wait 返回后才读取，因而安全。
	// Done 必须与 Add 登记数量一致：少调用会一直等待，多调用可能 panic。
	// WaitGroup 只负责等待完成，不传结果、错误，也没有内建取消或超时。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：并发结果按输入索引保存
	// ═══════════════════════════════════════════════════════
	// 预先创建固定长度的结果切片，每个任务只写自己的位置。
	// 任务完成顺序可以变化，但最终按索引读取仍对应原输入顺序。

	// 批量独立工作：每个任务仅写固定长度结果切片中自己的一格。
	// 任务可以乱序结束；最后按输入顺序打印，所以示例输出是稳定的。
	titles := []string{"Go Basics", "API Design", ""}
	results := buildPages(titles)
	for index, result := range results {
		fmt.Printf("页面 %d：%+v\n", index, result)
	}
	fmt.Println("空输入结果数：", len(buildPages(nil))) // 0

	// ═══════════════════════════════════════════════════════
	// 第 3 节：独立实例与共享数据边界
	// ═══════════════════════════════════════════════════════
	// 传入两个不同的 PageResult 指针，让每个任务有各自的写入目标。
	// 若结构体内部又指向相同 map，即使外层实例不同，也需要保护共享内容。

	// 相互独立的实例也可交给不同 goroutine 通过指针填写。
	first, second := PageResult{}, PageResult{}
	var separate sync.WaitGroup
	separate.Add(2)
	go buildInto("Admin", &first, &separate)
	go buildInto("Settings", &second, &separate)
	separate.Wait()
	fmt.Println("独立实例：", first.Slug, second.Slug) // admin settings
	// 本例字段仅为 string/bool。若两个结构体内的 map、slice、指针共享数据，
	// “顶层是两个实例”仍不能保证没有竞争，需要另外安排所有权或同步。

	// 数据竞争：多个 goroutine 并发访问同一内存位置，至少一个是写入，且没有同步。
	// 不要并发 append 同一个 slice，不要无保护写共享 map，也不要无保护 total++。
	// Wait 只保证读取发生在任务完成之后，不会消除多个任务执行期间彼此的竞争。
	// 下章用 channel 交付结果，后续并发安全章再讲 Mutex。
	// 可运行 go run -race ./14-goroutines 检查本次执行中是否发现竞争。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：WaitGroup.Go 与分批复用
	// ═══════════════════════════════════════════════════════
	// Go 方法代办 Add、启动 goroutine 和完成通知，Wait 仍由调用方安排。
	// 不要再在回调里调用 Done，否则任务数量会被重复扣减。

	// Go 1.25 起，WaitGroup.Go 把“登记、启动、完成通知”合成一次调用。
	// 仍需先启动任务再 Wait；函数正常返回后自动完成，不要再手动 Add 或 Done。
	var modern sync.WaitGroup
	var message string
	modern.Go(func() {
		message = "task completed"
	})
	modern.Wait()
	fmt.Println("WaitGroup.Go：", message)
	// API 约定传入的函数不能 panic；它不是异常捕获器，也不会把失败变成返回值。
	// 复用同一个 WaitGroup 等下一批任务时，必须等上一批所有 Wait 调用都已返回，
	// 才为新批次调用 Add 或 Go。不要在上一批 Wait 正在收尾时立即复用。
	modern.Go(func() {
		message = "second batch completed"
	})
	modern.Wait()
	fmt.Println("上一批结束后复用：", message)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：限制并发数并收集错误
	// ═══════════════════════════════════════════════════════
	// maxWorkers 限制启动数量，静态分工让每个 worker 处理互不重叠的索引。
	// invalid 的解析错误保存在对应槽位，其他有效输入仍能得到各自结果。

	// 不默认“每条数据一个 goroutine”：数据很大时会消耗大量内存与外部服务容量。
	// 此处先学固定数量任务分摊工作；第 15 章再用 channel 实现动态领任务的 worker pool。
	parsed, accepted := parseVersionsBounded([]string{"3", "invalid", "7"}, 2)
	fmt.Println("并发上限参数有效：", accepted)
	for _, result := range parsed {
		fmt.Printf("版本 %q -> %d，解析失败=%v\n", result.Input, result.Number, result.Err != nil)
	}
	_, accepted = parseVersionsBounded(nil, 0)
	fmt.Println("拒绝无效并发数：", accepted) // false
}

type VersionResult struct {
	Input  string
	Number int
	Err    error // 先理解为失败信息；第 18 章再系统学习 error。
}

func parseVersionsBounded(inputs []string, maxWorkers int) ([]VersionResult, bool) {
	if maxWorkers < 1 {
		return nil, false
	}
	results := make([]VersionResult, len(inputs))
	workers := min(maxWorkers, len(inputs))
	var wait sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		start := worker
		wait.Go(func() {
			// 第 0 个任务处理 0、workers、2*workers……，各任务负责的索引互不重叠。
			for index := start; index < len(inputs); index += workers {
				number, err := strconv.Atoi(inputs[index])
				results[index] = VersionResult{Input: inputs[index], Number: number, Err: err}
			}
		})
	}
	wait.Wait()
	return results, true
	// 空输入不启动任务，Wait 立即返回；每次最多 maxWorkers 个 goroutine。
	// 失败和成功都存入各自结果槽，由调用方在 Wait 后统一处理，而不是静默丢弃错误。
	// 这种静态分配适合耗时相近的独立任务；任务耗时差异大时，用下一章的队列分配更均衡。
}

func buildPages(titles []string) []PageResult {
	results := make([]PageResult, len(titles))
	var wait sync.WaitGroup
	wait.Add(len(titles))
	for index, title := range titles {
		// 显式传参让“每个任务拿到自己的索引和值”更直观，也不依赖循环变量版本差异。
		go func(i int, value string) {
			defer wait.Done()
			results[i] = buildPage(value)
		}(index, title)
	}
	wait.Wait()
	return results
}

func buildPage(title string) PageResult {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return PageResult{Title: title}
	}
	return PageResult{
		Title: title,
		Slug:  strings.ReplaceAll(strings.ToLower(trimmed), " ", "-"),
		Valid: true,
	}
}

func buildInto(title string, result *PageResult, wait *sync.WaitGroup) {
	defer wait.Done()
	*result = buildPage(title)
}

// 后续再学：动态任务队列（15 章）、取消与超时（16 章）、共享状态加锁（17 章）；
// goroutine 不是越多越快，应根据任务规模、CPU 与外部服务限制选择并发量。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 maxWorkers 改成 1、2、10，比较输出顺序和输入顺序是否一致。
// 说明为什么 Wait 后读取结果安全，却不能保证多个任务并发 total++ 安全。
