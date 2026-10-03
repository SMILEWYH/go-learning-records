/*
章节：08-函数与闭包
	本章掌握清单：
	1. 声明函数：参数与返回类型
	2. 多返回值与空白标识符
	3. 命名返回值
	4. 变参：零个、多个与切片展开
	5. 函数值、签名与 nil
	6. 闭包：保存并修改外层变量
	7. 回调：把变化传入通用流程
	8. 立即调用、延迟调用与变量快照
	9. 循环闭包：独立变量与复用变量
	10. 递归：基础情况与递进过程
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./08-functions
// 运行练习：go run ./08-functions exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：参数写成 name type；多个结果可直接返回，不一定要包装成对象或元组。
// 后续再学：值传递与指针见第 09 章；defer 的求值时机见第 10 章；并发闭包状态见第 14/17 章。

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：声明函数：参数与返回类型
	// ═══════════════════════════════════════════════════════
	// normalizeUsername 的签名是 func normalizeUsername(name string) string。
	// 括号里说明输入，括号后说明输出；调用时只传值，不再写参数类型。

	// 具名函数通常声明在包级，参数类型写在参数名后。
	fmt.Println("规范化用户名：", normalizeUsername("  Alice  "))
	// Go 不支持通过同名函数不同参数来重载，也不支持声明默认参数值。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：多返回值与空白标识符
	// ═══════════════════════════════════════════════════════
	// findUsername 返回 (string, bool)，分别表达内容和是否找到。
	// 找到空昵称与用户不存在都可能返回 ""，第二个结果保留了区别。

	// 多个返回值一次接收，用 _ 明确忽略某个结果。
	users := map[string]string{"u-01": "alice", "u-02": ""}
	name, found := findUsername(users, "u-02")
	fmt.Printf("查询昵称：%q 存在=%t\n", name, found)
	_, missingFound := findUsername(users, "missing")
	fmt.Println("不存在的用户：", missingFound)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：命名返回值
	// ═══════════════════════════════════════════════════════
	// offset、limit 已由返回列表声明，函数体内给它们赋值要用 =。
	// page=3、pageSize=20 时，前两页共占 40 条，所以 offset=40。

	// 命名返回值是在函数进入时就存在的局部变量，初值为零值。
	offset, limit := pagination(3, 20)
	fmt.Println("命名返回值：", offset, limit)
	// 多返回值是调用约定，不是自动生成一个可存进单个变量的对象。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：变参：零个、多个与切片展开
	// ═══════════════════════════════════════════════════════
	// values ...int 表示调用方可以传不同数量的 int；函数内部看到的是 []int。
	// 已有 []int 时写 counts...，而不是直接把切片当成单个 int。

	// 变参参数在函数内部是切片；变参只能有一个且必须是最后一个参数。
	fmt.Println("零个参数：", sumCounts())
	fmt.Println("多个参数：", sumCounts(2, 3, 5))
	counts := []int{4, 6, 8}
	fmt.Println("展开切片：", sumCounts(counts...))
	// counts... 展开不意味着复制底层数组；若函数改动参数元素，调用方也可能受影响。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：函数值、签名与 nil
	// ═══════════════════════════════════════════════════════
	// 变量 combine 保存函数，combine(3, 7) 才真正调用它。
	// optionalTransform == nil 表示尚未配置回调，先补默认函数再调用。

	// 函数也是值，可用 var 或 := 保存匿名函数，不能使用 const。
	add := func(a, b int) int {
		return a + b
	}
	var combine func(int, int) int = add
	fmt.Println("函数变量：", combine(3, 7))
	// 函数类型由参数/结果的个数、顺序、类型和是否变参决定，参数名字不影响类型。
	// func(int, int) int 与 func(a int, b int) int 相同；func(...int) int 与 func([]int) int 不同。
	// Go 不允许在函数体里再写 func named() {} 的具名声明；用匿名函数赋给变量。
	// 函数类型的零值是 nil；调用 nil 函数会 panic，所以使用前应确保已经赋值。
	var optionalTransform func(string) string
	if optionalTransform == nil {
		optionalTransform = strings.TrimSpace
	}
	fmt.Printf("默认回调：%q\n", optionalTransform(" ready "))
	// 函数值只能和 nil 比较，不能用 == 判断两个函数是否相同，也不能直接作为 map 的 key。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：闭包：保存并修改外层变量
	// ═══════════════════════════════════════════════════════
	// newSequence 返回的不是一个序号，而是下次调用时能生成序号的函数。
	// 每次调用 newSequence 都建立独立的 next；A 与 B 不会共用计数。

	// 闭包捕获外层变量；函数返回后，捕获的变量仍可继续存在。
	nextA := newSequence(100)
	nextB := newSequence(0)
	fmt.Println("独立的序列 A：", nextA(), nextA())
	fmt.Println("独立的序列 B：", nextB(), nextB())
	// 闭包可以修改捕获的变量，并非保存创建时的一份不可变快照。
	// 这个计数器用于顺序调用，并未保证并发安全；并发相关约束在后续章节学习。

	// ═══════════════════════════════════════════════════════
	// 第 7 节：回调：把变化传入通用流程
	// ═══════════════════════════════════════════════════════
	// transformNames 管理遍历和结果收集，transform 决定每一项怎样变。
	// 传 normalizeUsername 时没有括号，因为这里需要函数值而不是本次计算结果。

	// 回调函数把“要做的变化”作为参数交给通用处理流程。
	input := []string{" alice ", " BOB "}
	normalized := transformNames(input, normalizeUsername)
	upper := transformNames(input, func(value string) string {
		return strings.ToUpper(strings.TrimSpace(value))
	})
	fmt.Println("原数据：", input, "规范化：", normalized, "大写：", upper)
	// 这里的回调会立即同步执行；传递一个函数并不自动产生异步行为。

	// ═══════════════════════════════════════════════════════
	// 第 8 节：立即调用、延迟调用与变量快照
	// ═══════════════════════════════════════════════════════
	// func() string { ... } 后面的 () 表示现在调用；没有 () 则只是保存函数。
	// later 在调用时读取 stage，所以修改 stage 会影响它之后返回的值。

	// 匿名函数后跟 () 会立即调用；只保存函数值则要等后续调用才读取捕获变量。
	stage := "queued"
	immediate := func() string { return stage }()
	later := func() string { return stage }
	savedStage := stage
	snapshot := func() string { return savedStage }
	stage = "running"
	fmt.Println("立即调用结果：", immediate, "稍后调用结果：", later(), "独立快照：", snapshot())
	// 输出 queued、running、queued。snapshot 捕获的是 savedStage；这里之后不再修改它。

	// ═══════════════════════════════════════════════════════
	// 第 9 节：循环闭包：独立变量与复用变量
	// ═══════════════════════════════════════════════════════
	// callbacks 存的是三个函数，循环结束后才调用。
	// 第一组各自记住一轮 index，第二组共同读取 reused，最终都得到 2。

	// Go 1.22 起，用 := 声明的循环变量每轮独立，闭包按各自那一轮的变量读取。
	callbacks := []func() int{}
	for index := range 3 {
		callbacks = append(callbacks, func() int { return index })
	}
	fmt.Println("每轮独立的变量：", callbacks[0](), callbacks[1](), callbacks[2]()) // 0 1 2。
	var reused int
	callbacks = nil
	for reused = range 3 {
		callbacks = append(callbacks, func() int { return reused })
	}
	fmt.Println("复用外层同一变量：", callbacks[0](), callbacks[1](), callbacks[2]()) // 2 2 2。
	// 旧版本 Go 常见的“闭包总取到最后一个循环值”不能不分版本、语法地照搬。

	// ═══════════════════════════════════════════════════════
	// 第 10 节：递归：基础情况与递进过程
	// ═══════════════════════════════════════════════════════
	// 空切片是基础情况，直接 return；其余调用每次少传一个路径部分。
	// 先递归后打印，所以先出现 workspace，再出现更长的路径。

	// 递归就是函数调用自己，必须有终止条件，并让每次调用接近这个条件。
	// 例子：从短到长打印资源层级，生成面包屑导航使用的路径。
	printBreadcrumb([]string{"workspace", "team", "project"})
	printBreadcrumb(nil) // 空路径直接返回，不打印任何内容。
	// 过深的递归会消耗调用栈；生产代码处理不受控深度时应限制深度或改成循环。
}

func normalizeUsername(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func findUsername(users map[string]string, id string) (string, bool) {
	name, exists := users[id]
	return name, exists
}

func pagination(page, pageSize int) (offset, limit int) {
	// 为了专注命名返回值，这个示例约定 page >= 1、pageSize > 0，且运算不溢出。
	offset = (page - 1) * pageSize
	limit = pageSize
	// 可以只写 return（裸返回），但显式列出值通常更方便阅读。
	return offset, limit
}

func sumCounts(values ...int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

func newSequence(start int) func() int {
	next := start
	return func() int {
		value := next
		next++
		return value
	}
}

func transformNames(values []string, transform func(string) string) []string {
	// 约定 transform 非 nil。结果用独立切片保存，保留输入顺序，不修改输入。
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, transform(value))
	}
	return result
}

func printBreadcrumb(parts []string) {
	if len(parts) == 0 {
		return // 基础情况：没有层级，递归终止。
	}
	printBreadcrumb(parts[:len(parts)-1]) // 每次少一个层级。
	fmt.Println("资源路径：", strings.Join(parts, "/"))
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 连续调用 nextA 三次，并在中间调用 nextB，预测两个序列各自的值。
// 把 printBreadcrumb 的打印移到递归调用前，观察路径输出顺序如何变化。
