/*
章节：18-错误与异常处理
	本章掌握清单：
	1. error 接口与显式检查
	2. 包装错误与 errors.Is
	3. 自定义错误与 errors.As
	4. Unwrap 与 Join：一层包装和多个原因
	5. 带类型的 nil error
	6. panic / recover 与任务边界
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./18-errors
// 运行练习：go run ./18-errors exam
// 阅读顺序：顶部定义错误类型与辅助函数，main 逐个演示如何处理它们，按调用名回看即可。
// TS 对照：TS 常用 throw/try/catch；Go 中可预期的失败通常返回 (结果, error)，
// 调用者用 if err != nil 处理。panic 不替代正常业务错误。
// 后续再学：HTTP/数据库错误映射、超时重试策略、统一日志与监控；不要用重试掩盖所有错误。

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrUserNotFound = errors.New("用户不存在")

type ValidationError struct {
	Field   string
	Message string
}

// 实现 Error() string 就实现了内置 error 接口。
func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func parseQuantity(raw string) (int, error) {
	quantity, err := strconv.Atoi(raw)
	if err != nil {
		// %w 保留原错误链；%v 只保留文字，不支持继续按错误类型检查。
		return 0, fmt.Errorf("解析购买数量 %q: %w", raw, err)
	}
	if quantity <= 0 {
		return 0, &ValidationError{Field: "quantity", Message: "必须大于 0"}
	}
	return quantity, nil
}

func loadUser(id string) (string, error) {
	if id != "u-1" {
		return "", fmt.Errorf("查询用户 %q: %w", id, ErrUserNotFound)
	}
	return "小林", nil
}

// runTaskBoundary 模拟后台任务边界：兜住某个任务中的意外 panic，
// 把失败交还给调度者，由它记录日志、告警并决定重试。
// 只有已经确认任务彼此隔离、进程状态仍可信的边界才适合继续运行。
// recover 必须由同一个 goroutine 内的 deferred 函数直接调用。
// 父 goroutine 的 defer 无法接住另一个 goroutine 的 panic。
// recover 不能处理 os.Exit、进程被终止或 runtime 的致命错误（如并发 map 写）。
// 它只恢复当前展开的 panic，不会撤销之前已发生的文件写入或业务副作用。
func runTaskBoundary(task func() error) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("后台任务发生意外 panic: %v", value)
			// 真实服务在这里记录堆栈（runtime/debug.Stack）等诊断信息。
			// 给客户端返回通用失败信息，避免把内部堆栈发给客户端。
		}
	}()
	return task()
	// recover 后返回本函数调用者，不会跳回 panic 那一行继续执行。
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：error 接口与显式检查
	// ═══════════════════════════════════════════════════════
	// error 的核心方法是 Error() string，nil 表示没有错误。
	// parseQuantity 先解析文本，再检查数量是否大于 0，调用方先处理 err 再使用 quantity。

	quantity, err := parseQuantity("3")
	if err != nil {
		fmt.Println("解析失败：", err)
		return
	}
	fmt.Println("购买数量：", quantity)

	// ═══════════════════════════════════════════════════════
	// 第 2 节：包装错误与 errors.Is
	// ═══════════════════════════════════════════════════════
	// ErrUserNotFound 是用来识别一种失败的固定错误值。
	// loadUser 用 %w 添加查询上下文，Is 仍能在包装里面找到这个原因。

	_, err = loadUser("missing")
	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("可以映射为 HTTP 404：", err)
	}
	// errors.Is 沿错误链查找某个错误；不要通过比较错误字符串识别错误。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：自定义错误与 errors.As
	// ═══════════════════════════════════════════════════════
	// Is 询问是否包含某个原因；As 查找指定错误类型，并取出它的字段。
	// validationErr 是 *ValidationError，传 &validationErr 是让 As 能填入匹配到的指针。

	_, err = parseQuantity("0")
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		fmt.Println("可以映射为 HTTP 400，字段：", validationErr.Field)
	}
	// errors.As 沿错误链查找某种类型，把匹配结果放入目标变量。
	_, err = parseQuantity("abc")
	var numberErr *strconv.NumError
	if errors.As(err, &numberErr) {
		fmt.Println("底层数字解析错误：", numberErr.Num)
	}

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Unwrap 与 Join：一层包装和多个原因
	// ═══════════════════════════════════════════════════════
	// Unwrap 取单个包装的下一层；Join 可以同时保留读取失败和关闭失败。
	// 业务判断通常直接用 Is/As，这样不用自己逐层拆开。

	// 单个包装错误可以用 Unwrap 查看直接内层；业务判断仍优先使用 Is/As。
	wrapped := fmt.Errorf("生成用户页面: %w", ErrUserNotFound)
	fmt.Println("直接内层错误：", errors.Unwrap(wrapped))
	// 多个独立步骤都失败时，Join 保留各分支；常用于“处理失败 + 清理也失败”。
	readFailure := errors.New("读取失败")
	closeFailure := errors.New("关闭失败")
	joined := errors.Join(readFailure, closeFailure)
	fmt.Println("聚合错误能识别两种原因：", errors.Is(joined, readFailure), errors.Is(joined, closeFailure))
	fmt.Println("全部为 nil 时 Join 也为 nil：", errors.Join(nil, nil) == nil)
	// Is/As 会查找多个分支；Unwrap 仅针对 Unwrap() error，不展开 Join 的 []error。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：带类型的 nil error
	// ═══════════════════════════════════════════════════════
	// pointerError 是 nil 指针，装入 error 后仍带有 *ValidationError 类型信息。
	// 成功时直接 return nil，才能让调用方的 err != nil 判断符合预期。

	var pointerError *ValidationError
	var interfaceError error = pointerError
	fmt.Println("指针为 nil：", pointerError == nil, "error 接口为 nil：", interfaceError == nil)
	// error 接口包含动态类型和值；这里的动态类型仍是 *ValidationError，所以不等于 nil。
	// 成功分支应直接 return nil，不要把带类型的 nil 指针转换为 error 返回。
	// 对这样的 nil 指针调用 Error()，若方法解引用接收者仍可能 panic。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：panic / recover 与任务边界
	// ═══════════════════════════════════════════════════════
	// runTaskBoundary 在执行任务的同一个 goroutine 内登记 defer，再调用任务。
	// 恢复后把错误返回调度者，不会回到 panic 位置继续执行。

	result := make(chan error, 1)
	go func() {
		// recovery 边界必须在真正执行任务的 goroutine 中。
		result <- runTaskBoundary(func() error {
			panic("演示：第三方任务出现未预料到的状态")
		})
	}()
	fmt.Println("调度者收到失败：", <-result)
	fmt.Println("调用边界之后的代码仍能继续")
	// 无法初始化关键配置时，可让 main 输出错误并退出；不必为了退出而 panic。
	// 通常底层加上下文后返回，上层统一决定记录日志、回应请求或退出，避免每层重复记录。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 parseQuantity 的输入依次改成 "3"、"0"、"abc"，区分成功、业务校验失败、解析失败。
// 把 loadUser 中的 %w 临时改成 %v，预测 errors.Is 是否还能找到原错误。
