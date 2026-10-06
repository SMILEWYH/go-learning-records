// 章节：19-泛型
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

// 练习 1：接口结果筛选。写一个通用函数供订单、商品等列表使用。
// 建议：func Filter[T any](items []T, keep func(T) bool) []T。
// 要求：保持原顺序，不修改输入切片；回调约定非 nil。
// 常规：[1,2,3,4] 按偶数筛选 -> [2,4]；再用含 Active 字段的用户切片筛选启用用户。
// 边界：nil 输入、全部不满足条件 -> 长度为 0；筛选后修改结果元素不能改变输入元素。
// 本题元素用 int 或只含值字段的结构体即可，不要求深拷贝嵌套指针。
func exercise1() {
	// TODO：实现 Filter，至少使用两种元素类型核对结果。
}

// 练习 2：把接口返回的数组转换为按 ID 查找的索引。
// 建议：func IndexBy[T any, K comparable](items []T, keyOf func(T) K) map[K]T。
// 要求：返回非 nil 的 map；相同 key 后出现的记录覆盖前面的；回调约定非 nil。
// 常规：[{ID:"u-1",Name:"旧名"},{ID:"u-1",Name:"新名"}] -> key "u-1" 对应 "新名"。
// 边界：空输入 -> 可直接写入的空 map；空字符串 ID 是合法 key，按同样覆盖规则处理。
// 再用整数 ID 验证一次；说明为什么 K 需要 comparable。
func exercise2() {
	// TODO：实现 IndexBy，验证重复 key、空输入及不同 key 类型。
}

// 练习 3：内存列表分页。复用 lesson.go 的 Page[T] 表达返回结果。
// 建议：func Paginate[T any](items []T, page, pageSize int) (Page[T], error)。
// 要求：page 从 1 开始；pageSize 为 1..100；Total 始终为输入总数；结果 Items 复制一份。
// 常规：[1,2,3,4,5], page=2, pageSize=2 -> Items=[3,4], Total=5。
// 边界：越界页和空输入 -> Items 长度 0；page<=0 或非法 pageSize -> error。
// page 极大也应得到空页，避免直接计算 (page-1)*pageSize 导致整数溢出。
func exercise3() {
	// TODO：实现 Paginate，用整数列表与业务结构体列表验证。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func Filter[T any](items []T, keep func(T) bool) []T {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Filter；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func IndexBy[T any, K comparable](items []T, keyOf func(T) K) map[K]T {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：IndexBy；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func Paginate[T any](items []T, page, pageSize int) (Page[T], error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：Paginate；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
