/*
章节：17-并发安全与锁
*/

// TS 对照：前端事件循环不代表所有环境都只有一条执行线程；Go 的多个
// goroutine 可以真正并行。共享变量的“读、修改、写”必须整体受到保护。

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 带锁的类型和方法保留在包级；main 按“互斥锁、读写锁、并发容器”逐节使用。
// Mutex 的零值可直接使用。锁使用后不能复制，所以方法用指针接收者，
// Counter 也应通过指针传递；不要把含锁的结构体按值放进切片再复制。
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta // 整个更新过程位于临界区内。
}

func (c *Counter) Value() int {
	c.mu.Lock() // 读也要遵守同一套锁规则。
	defer c.mu.Unlock()
	return c.value
}

// RWMutex 允许多个读者同时持有 RLock，写者需要独占 Lock。
// 它不保证比 Mutex 快；优先选择清楚、正确的方案，再根据测量优化。
// 不要持有 RLock 再申请 Lock 来“升级”：它会等待包括自己的读锁释放。
// 需要先释放读锁，取得写锁后重新检查条件，因为中途状态可能已经变化。
type Settings struct {
	mu     sync.RWMutex
	values map[string]string
}

func (s *Settings) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.values[key]
	return value, ok
}

func (s *Settings) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.values == nil {
		s.values = make(map[string]string)
	}
	s.values[key] = value
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节：Mutex：保护完整的读改写过程
	// ═══════════════════════════════════════════════════════
	// balance.Add 内先 Lock，再修改 value，最后 Unlock；同一时刻只有一个任务能修改它。
	// 上方 Counter.Value 读取时也使用同一把锁，才能与写入遵守同一套同步规则。

	var balance Counter
	var wg sync.WaitGroup
	for _, delta := range []int{1, -1} {
		wg.Add(1) // 启动前登记；不要在 goroutine 内才 Add。
		go func(delta int) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				balance.Add(delta)
			}
		}(delta)
	}
	wg.Wait()
	fmt.Println("加减后的余额：", balance.Value()) // 0

	// ═══════════════════════════════════════════════════════
	// 第 2 节：RWMutex：区分读操作与写操作
	// ═══════════════════════════════════════════════════════
	// Settings.Get 使用读锁，Set 使用写锁；持有读锁期间不能修改 map。
	// 是否使用读写锁应根据访问模式和测量决定，先保证所有访问都有保护。

	var settings Settings
	settings.Set("theme", "dark")
	theme, exists := settings.Get("theme")
	fmt.Println("读取配置：", theme, exists)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：sync.Map 与条目自身的锁
	// ═══════════════════════════════════════════════════════
	// 普通 map 加锁通常更容易保持类型清晰、维护多个字段之间的约束。
	// sync.Map 适用于特定并发访问模式；Load/Store 各自安全，
	// 但 Load -> 修改 -> Store 这一串操作不会自动变成事务。
	var requestCounts sync.Map // key: string, value: *Counter
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// LoadOrStore 把“没有就存入”合为一个原子操作。
			// 参数仍会先求值，不能用它保证昂贵初始化函数只执行一次。
			actual, _ := requestCounts.LoadOrStore("/api/profile", &Counter{})
			entry := actual.(*Counter)
			for j := 0; j < 250; j++ {
				// map 只保护条目，不保护条目所指向的对象；对象需要自己的锁。
				entry.Add(1)
			}
		}()
	}
	wg.Wait()
	actual, ok := requestCounts.Load("/api/profile")
	if ok {
		fmt.Println("接口访问次数：", actual.(*Counter).Value()) // 1000
	}
	// sync.Map.Range 可遍历，但并发更新时不是一致性快照。
	requestCounts.Range(func(key, value any) bool {
		fmt.Printf("访问统计 %s：%d\n", key, value.(*Counter).Value())
		return true // false 停止遍历；遍历顺序不固定。
	})
	requestCounts.Delete("/api/profile")
	_, exists = requestCounts.Load("/api/profile")
	fmt.Println("删除统计后仍存在：", exists)
	// 另一种思路是让一个 goroutine 独占状态，通过 channel 接受操作请求。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Once 与 atomic：一次初始化和独立计数
	// ═══════════════════════════════════════════════════════
	// once.Do 让 endpoint 初始化一次；requests.Add 让单个计数更新不可被拆开。
	// 多个操作之间有关联时，不能把“每步原子”当成“整体原子”。

	// Once 适合初始化后只读的共享配置；不是会自动重试的任务执行器。
	// Do 的函数即使 panic，也会被视为已经执行；不要在同一个 Once 的函数内再次 Do。
	demoOnceAndAtomic := func() {
		var once sync.Once
		var endpoint string
		var requests atomic.Int64 // 零值可用；使用后同样不能复制。
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				once.Do(func() { endpoint = "/api/v1" })
				// Do 返回后能看到初始化写入的值；后续不能再无锁修改 endpoint。
				if endpoint != "" {
					requests.Add(1)
				}
			}()
		}
		wg.Wait()
		fmt.Println("一次初始化后的路由：", endpoint, "请求数：", requests.Load())
		// atomic 适合独立计数等单个状态；Load 后判断再 Add 仍是多个操作。
		// “余额够才扣款，同时写账本”需要一个整体协议，通常用锁更清楚。
	}

	demoOnceAndAtomic()

	// ═══════════════════════════════════════════════════════
	// 第 5 节：避免死锁与理解竞争检测
	// ═══════════════════════════════════════════════════════
	// 避免死锁的实用规则：
	// 1. Mutex 不可重入：持锁的函数不要再调用会取得同一把锁的方法。
	// 2. 必须同时取得多把锁时，所有调用路径遵守相同顺序。
	// 3. 锁内只处理受保护状态；尽量不要做网络调用、磁盘操作或不受控回调。
	// 锁粒度太大可能串行化全部请求；太小又可能无法保护完整业务不变量。
	// race detector 检查竞争，不是死锁检测器；没有竞争也可能永久等待。
}
