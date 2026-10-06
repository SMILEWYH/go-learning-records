// 章节：17-并发安全与锁
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
	"sync"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 练习 1：库存扣减。多个订单会同时购买同一个 SKU。
// 要求：定义 Inventory，使用 Mutex 保护库存；检查库存和扣减必须在同一把锁内。
// 建议：func (i *Inventory) Reserve(count int) bool；func (i *Inventory) Remaining() int。
// 初始库存 10，让 20 个 goroutine 各扣 1，WaitGroup 等待结束：恰好 10 次成功，剩余 0。
// 边界：扣减 0 或负数失败，库存不变；库存不足失败，库存不能变为负数。
// 提示：不要无锁地共享“成功次数”，可使用 channel 收集每个 goroutine 的结果。
func exercise1() {
	// TODO：自己定义类型、实现方法，并打印常规与边界用例结果。
}

// 练习 2：用户权限缓存。缓存的每个用户对应一个 []string 权限列表。
// 要求：用 map + RWMutex；Set 和 Get 都复制切片，避免调用者从锁外修改内部数据。
// 建议：func (c *PermissionCache) Set(userID string, permissions []string)
//
//	func (c *PermissionCache) Get(userID string) ([]string, bool)
//
// 常规：存入 ["read"] 后读出 ["read"], true；同时多次 Set/Get 不出现竞争。
// 边界：修改 Set 的输入切片或 Get 的返回切片，缓存仍为 ["read"]；不存在返回 nil, false。
func exercise2() {
	// TODO：实现并用 go run -race ./17-concurrency-safety exam 检查。
}

// 练习 3：Webhook 事件去重。这里只实现进程内演示，重启后记录不保留。
// 要求：使用 sync.Map.LoadOrStore 实现“记录事件 ID 并判断是不是第一次”。
// 建议：func (d *EventDeduplicator) FirstSeen(eventID string) bool。
// 常规：20 个 goroutine 同时传 "evt-42"，恰好一次得到 true；新 ID "evt-43" 得到 true。
// 边界：空 ID 返回 false 且不存入；重复 ID 返回 false。不能把 Load 与 Store 拆开。
// 思考：如果事件处理失败，应如何设计重试？本题只判断首次出现，不代表业务处理成功。
func exercise3() {
	// TODO：定义去重器，汇总并打印并发调用结果。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type Inventory struct {
	mu    sync.Mutex
	stock int
}

func (i *Inventory) Reserve(count int) bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Reserve；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func (i *Inventory) Remaining() int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Remaining；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type PermissionCache struct {
	mu     sync.RWMutex
	values map[string][]string
}

func (c *PermissionCache) Set(userID string, permissions []string) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Set；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func (c *PermissionCache) Get(userID string) ([]string, bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Get；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type EventDeduplicator struct{ seen sync.Map }

func (d *EventDeduplicator) FirstSeen(eventID string) bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：FirstSeen；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
