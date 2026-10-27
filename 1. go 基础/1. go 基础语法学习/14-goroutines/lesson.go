/*
章节：14-协程与任务等待
*/

// 在函数前使用 go 关键字, go f() 可以并发在新的 goroutine 中执行 f
// goroutine 由 Go 运行时调度, 并发不代表所有任务一定在同一时刻并行运行

package main

import (
	"fmt"
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
	// ═══════════════════════════════════════════════════════
	// 第 1 节：go 启动任务，WaitGroup 等待完成
	// ═══════════════════════════════════════════════════════
	// 主 goroutine 结束时进程退出, 不会自动等待其他 goroutine 完成
	// bad idea: 用 Sleep 猜任务何时完成不可靠、忙循环轮询一个共享状态还可能产生数据竞争
	// good ldea: WaitGroup
	// WaitGroup 的零值可用, 开始使用后不要复制它, 传递时使用指针
	// 可类比 Ts 的 Promise.all 的等待全部完成, 但 WaitGroup 不替你收集返回值或错误
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
	// 主 goroutine 等 Wait 返回后才读取，因而安全
	// Done 必须与 Add 登记数量一致：少调用会一直等待，多调用可能 panic
	// WaitGroup 只负责等待完成，不传结果、错误，也没有内建取消或超时

	// ═══════════════════════════════════════════════════════
	// 第 2 节：循环中的并发任务
	// ═══════════════════════════════════════════════════════
	// 循环并发任务可能会乱序结束
	// good idea: 声明一个切片, 在循环执行并发任务时, 给各个并发任务分发一个专属的切片下索引
	// 这样无论它是首个结束还是最后结束, 它的结果都可以保存切片的对应下索引中
	buildPage := func(title string) PageResult {
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

	buildPages := func(titles []string) []PageResult {
		results := make([]PageResult, len(titles))
		var wait sync.WaitGroup
		wait.Add(len(titles))
		for index, title := range titles {
			go func(i int, value string) {
				defer wait.Done()
				results[i] = buildPage(value)
			}(index, title)
		}
		wait.Wait()
		return results
	}

	titles := []string{"Go Basics", "API Design", ""}
	results := buildPages(titles)
	for index, result := range results {
		fmt.Printf("页面 %d：%+v\n", index, result)
	}
	fmt.Println("空输入结果数：", len(buildPages(nil))) // 0

	// ═══════════════════════════════════════════════════════
	// 第 3 节：独立实例与共享数据边界
	// ═══════════════════════════════════════════════════════
	// 上面两小节的例子都是可以通过上下文链去找到 WaitGroup 实例进行使用, 而如果无法依靠上下文链了呢?
	// good idea: 那么需要传递 WaitGroup 指针
	buildInto := func(title string, result *PageResult, wait *sync.WaitGroup) {
		defer wait.Done()
		*result = buildPage(title)
	}

	first, second := PageResult{}, PageResult{}
	var separate sync.WaitGroup
	separate.Add(2)
	go buildInto("Admin", &first, &separate)
	go buildInto("Settings", &second, &separate)
	separate.Wait()
	fmt.Println("独立实例：", first.Slug, second.Slug) // admin settings
	// 本例字段仅为 string 和 bool。若两个结构体内的 map、slice、指针共享数据,
	// “顶层是两个实例”仍不能保证没有竞争，需要另外安排所有权或同步。

	// 数据竞争：多个 goroutine 并发访问同一内存位置，至少一个是写入，且没有同步。
	// 不要并发 append 同一个 slice，不要无保护写共享 map，也不要无保护 total++。
	// Wait 只保证读取发生在任务完成之后，不会消除多个任务执行期间彼此的竞争。
	// 下章用 channel 交付结果，后续并发安全章再讲 Mutex。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：WaitGroup.Go 与分批复用
	// ═══════════════════════════════════════════════════════
	// Go 1.25 起, WaitGroup.Go 把"登记 Add、启动 go、完成通知 Done"合成一次调用
	// 仍需先启动任务再 Wait, 函数正常返回后自动完成
	var modern sync.WaitGroup
	var message string
	modern.Go(func() {
		message = "task completed"
	})
	modern.Wait()
	fmt.Println("WaitGroup.Go：", message)
	// API 约定传入的函数不能 panic；它不是异常捕获器，也不会把失败变成返回值
	// 复用同一个 WaitGroup 等下一批任务时，必须等上一批所有 Wait 调用都已返回, 才为新批次调用 Add 或 Go
	modern.Go(func() {
		message = "second batch completed"
	})
	modern.Wait()
	fmt.Println("上一批结束后复用：", message)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：限制并发数并收集错误
	// ═══════════════════════════════════════════════════════
	// 不默认每条数据一个 goroutine: 数据很大时会消耗大量内存与外部服务容量
	// 此处先学固定数量任务分摊工作；第 15 章再用 channel 实现动态领任务的 worker pool。
	type VersionResult struct {
		Input  string
		Number int
		Err    error // 先理解为失败信息；第 18 章再系统学习 error。
	}

	parseVersionsBounded := func(inputs []string, maxWorkers int) ([]VersionResult, bool) {
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

	parsed, accepted := parseVersionsBounded([]string{"3", "invalid", "7"}, 2)
	fmt.Println("并发上限参数有效：", accepted)
	for _, result := range parsed {
		fmt.Printf("版本 %q -> %d，解析失败=%v\n", result.Input, result.Number, result.Err != nil)
	}
	_, accepted = parseVersionsBounded(nil, 0)
	fmt.Println("拒绝无效并发数：", accepted) // false
}

// goroutine 不是越多越快，应根据任务规模、CPU 与外部服务限制选择并发量
