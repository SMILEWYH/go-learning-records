/*
章节:17-并发安全与锁
*/

// Go 的多个 goroutine 可以真正并行, 多个 goroutine 可能同时操作同一个变量, 而 count++、value += delta 这样的一行代码, 并不保证整个更新过程不可被其他 goroutine 干扰
// 所以并发更新共享变量时, 需要保护完整的“读、修改、写”过程

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Mutex 互斥锁 (Lock)
// 很霸道, 同一把锁保护的数据, 同一时刻只允许一个 goroutine 执行读写操作
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Add(delta int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta // 整个更新过程位于临界区内
}

func (c *Counter) Value() int {
	c.mu.Lock() // 读也要遵守同一套锁规则
	defer c.mu.Unlock()
	return c.value
}

// RWMutex 读写互斥锁 (RLock、Lock)
// 不是那么霸道, 同一把锁保护的数据, 同一时刻允许多个 goroutine 进行读(但不可以修改受同一把锁保护的数据), 但写操作会被一个 goroutine 独占, 并且写入期间, 其他读写操作都必须等待
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
	// 第 1 节:Mutex:保护完整的读改写过程
	// ═══════════════════════════════════════════════════════
	var balance Counter
	var wg sync.WaitGroup
	for _, delta := range []int{1, -1} {
		wg.Add(1)
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
	// 第 2 节:RWMutex:区分读操作与写操作
	// ═══════════════════════════════════════════════════════

	var settings Settings
	settings.Set("theme", "dark")
	wg.Add(3)

	// 两个读者：Get 使用 RLock，允许多个 goroutine 同时读取
	for reader := 1; reader <= 2; reader++ {
		go func(id int) {
			defer wg.Done()
			theme, _ := settings.Get("theme")
			fmt.Printf("读者 %d 读到主题：%s\n", id, theme)
		}(reader)
	}

	// 一个写者：Set 使用 Lock，写入期间其他读写操作需要等待
	go func() {
		defer wg.Done()
		settings.Set("theme", "light")
		fmt.Println("写者把主题改为 light")
	}()

	// 执行顺序不固定, 两位读者可能读到 dark 或 light
	wg.Wait()
	theme, exists := settings.Get("theme")
	fmt.Println("读取配置：", theme, exists) // light true

	// 这里示例只示范了结构体加锁, 其实只需要声明锁, 然后使用 Lock()、Unlock() 方法把代码包裹起来, 被包裹的代码块就是被加锁了
	// var mu sync.Mutex
	// var names []string

	// // 追加
	// mu.Lock()
	// names = append(names, "张三")
	// mu.Unlock()
	// 需要注意: 各个 goroutine 必须使用同一把锁. 如果每个 goroutine 内部各自创建一把锁，就无法互相协调

	// ═══════════════════════════════════════════════════════
	// 第 3 节:sync.Map 与条目自身的锁
	// ═══════════════════════════════════════════════════════
	// sync.Map 是 Go 标准库 sync 包提供的一种支持并发访问的 map
	// 多个 goroutine 可以同时调用它的方法, 无需为这些方法本身额外加锁

	// Load、Store 是 sync.Map 的方法, 分别用于查询和存储（新增或更新）数据
	// 单次 Load、Store 是并发安全的, 但 Load -> 修改 -> Store 整个过程不自动具备原子性

	// 相比普通 map, sync.Map 的键和值使用 any, 而普通 map 可以直接限定键和值的类型, 所以普通 map 加锁通常更容易保持类型清晰、维护多个字段之间的约束

	// sync.Map 适合的典型场景包括: 写入一次后反复读取，或者不同 goroutine 主要操作不同的键; 因为这两种场景的共同特点是: 对同一份数据的修改的冲突比较少, sync.Map 可以利用这种特点, 减少锁竞争
	// 它只保证自身操作的并发安全, 不会自动保护值所指向的对象内部的数据
	var requestCounts sync.Map // key: string, value: *Counter
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// LoadOrStore 是 sync.Map 的“有就取, 没有就存”方法
			actual, _ := requestCounts.LoadOrStore("/api/profile", &Counter{})
			// 类型断言
			entry := actual.(*Counter)
			for j := 0; j < 250; j++ {
				// map 只保护条目, 不保护条目所指向的对象; 对象需要自己的锁
				entry.Add(1)
			}
		}()
	}
	wg.Wait()
	actual, ok := requestCounts.Load("/api/profile")
	if ok {
		fmt.Println("接口访问次数：", actual.(*Counter).Value()) // 1000
	}
	// sync.Map.Range 可遍历数据, 但并发更新时不是一致性快照
	requestCounts.Range(func(key, value any) bool {
		fmt.Printf("访问统计 %s：%d\n", key, value.(*Counter).Value())
		return true // false 停止遍历; 遍历顺序不固定
	})
	// 删除数据
	requestCounts.Delete("/api/profile")
	_, exists = requestCounts.Load("/api/profile")
	fmt.Println("删除统计后仍存在：", exists) // false 不存在 key 为 "/api/profile" 的数据了

	// ═══════════════════════════════════════════════════════
	// 第 4 节: Once 与 atomic: 一次初始化和独立计数
	// ═══════════════════════════════════════════════════════
	// 这是两个并发工具:
	// 	1. sync.Once 保证多协程并发中某件事(初始化)只执行一次; 一个 goroutine 里可以调用多个 once.Do(), 限制是: 同一个 sync.Once 实例, 只执行一次传给 Do 的函数
	//  2. atomic.Int64 保证整数的单次读写、加减等操作是原子的; 多个操作之间有关联时, 不能把“每步原子”当成“整体原子”
	demoOnceAndAtomic := func() {
		var once sync.Once
		var endpoint string
		var requests atomic.Int64 // 零值可用; 使用后同样不能值复制, 应该传指针

		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 一个 goroutine 获得执行机会, 初始化 endpoint
				// 其他 goroutine 如果此时调用 Do, 会等待这次初始化完成
				// 初始化完成后，再调用这个 once.Do, 不会重复执行初始化
				once.Do(func() { endpoint = "/api/v1" })
				if endpoint != "" {
					requests.Add(1)
				}
			}()
		}
		wg.Wait()
		fmt.Println("一次初始化后的路由：", endpoint, "请求数：", requests.Load()) // /api/v1 请求数： 3
		// atomic 适合独立计数等单个状态;Load 后判断再 Add 仍是多个操作
		// “余额够才扣款, 同时写账本”需要一个整体协议, 通常用锁更清楚
	}

	demoOnceAndAtomic()

	// ═══════════════════════════════════════════════════════
	// 第 5 节:避免死锁与理解竞争检测
	// ═══════════════════════════════════════════════════════
	// 避免死锁的实用规则:
	// 1. Mutex 不可重入: 持锁的函数不要再调用会取得同一把锁的方法
	// 2. 必须同时取得多把锁时, 所有调用路径遵守相同顺序
	// 3. 锁内只处理受保护状态; 尽量不要做网络调用、磁盘操作或不受控回调
	// 锁粒度太大可能串行化全部请求; 太小又可能无法保护完整业务不变量
	// race detector 检查竞争, 不是死锁检测器; 没有竞争也可能永久等待

	// 示例一
	// 	var mu sync.Mutex

	// func outer() {
	// 	mu.Lock()
	// 	defer mu.Unlock()

	// 	inner()
	// }

	// func inner() {
	// 	mu.Lock() // 卡住：outer 还没释放这把锁
	// 	defer mu.Unlock()
	// }

	// 示例二
	// lockA.Lock()
	// lockB.Lock()

	// // 业务逻辑

	// lockB.Unlock()
	// lockA.Unlock()

	// 示例三
	// 	mu.Lock()
	// response := callRemoteAPI() // 假设网络请求耗时 5 秒
	// cache["result"] = response
	// mu.Unlock()
	// 建议改为:
	// 	response := callRemoteAPI()
	// mu.Lock()
	// cache["result"] = response
	// mu.Unlock()

}
