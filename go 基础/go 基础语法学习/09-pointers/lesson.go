/*
章节：09-指针与值传递
	本章掌握清单：
	1. 值传递与取地址、解引用
	2. 修改指向的值与修改指针本身
	3. nil、new 与返回局部变量地址
	4. 切片传参：共享元素，长度独立
	5. map 传参：共享映射与重新赋值
	6. new 与 make：返回类型和初始化内容
	7. 指针比较与使用场景
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./09-pointers
// 运行练习：go run ./09-pointers exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：TS 的对象赋值会共享对象；Go 的 struct/array 赋值复制整个值。
// Go 的所有参数都按值传递，传入指针时复制的也是指针值。

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：值传递与取地址、解引用
	// ═══════════════════════════════════════════════════════
	// Go 调用函数会复制实参。复制整数得到独立数值；复制指针得到指向同一变量的地址。
	// 读语法时区分：*int 是类型，&stock 是取地址，*p 是访问指针所指的值。

	// 普通值参数是副本，修改参数不会修改调用方变量。
	stock := 10
	changeCopy(stock)
	fmt.Println("普通参数修改后：", stock) // 10

	// &stock 的类型是 *int；指针保存了 stock 的地址。
	// *p 表示读写这个地址对应的 int。Go 不允许普通指针做地址加减运算。
	p := &stock
	*p = 12
	fmt.Println("解引用赋值后：", stock) // 12
	decreaseStock(p, 2)
	fmt.Println("通过指针修改后：", stock) // 10

	// ═══════════════════════════════════════════════════════
	// 第 2 节：修改指向的值与修改指针本身
	// ═══════════════════════════════════════════════════════
	// decreaseStock 用 *stock 修改外部整数；repoint 用 stock = ... 只改局部指针。
	// 两者虽然都接收 *int，但赋值左侧不同，效果也不同。

	// 指针本身也被复制。把形参改指向另一变量，不会替换外面的 p。
	repoint(p)
	fmt.Println("重新赋值指针形参后：", *p) // 10
	// 若目的是得到一个新值，直接返回新值通常比 **int 更清晰。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：nil、new 与返回局部变量地址
	// ═══════════════════════════════════════════════════════
	// nil 表示指针没有目标；指向数值 0 的有效指针则不等于 nil。
	// 因此 *int 能区分“没有设置”与“主动设置为 0”。

	// 指针零值是 nil：没有指向任何有效值。解引用 nil 会 panic。
	var optionalLimit *int
	fmt.Println("指针未设置：", optionalLimit == nil) // true
	optionalLimit = new(int)                    // 分配一个零值 int，返回 *int。
	fmt.Println("new(int) 指向：", *optionalLimit) // 0
	// 指向局部变量的指针可以安全返回；Go 会管理其存储位置和存活时间。
	limit := makeLimit(50)
	fmt.Println("返回的指针：", *limit) // 50

	// ═══════════════════════════════════════════════════════
	// 第 4 节：切片传参：共享元素，长度独立
	// ═══════════════════════════════════════════════════════
	// editSlice 接收到切片描述的副本，可以改共享数组的元素，却不能改调用方变量的长度。
	// 观察函数内 [99 2 3] 与函数外 [99 2]，把“元素变化”和“描述变化”分开判断。

	// slice 被传入时，会复制包含底层数组信息、长度和容量的描述符。
	// 两个 slice 仍可能引用同一个数组，因此通过元素下标修改可以相互看见。
	ids := make([]int, 2, 4)
	ids[0], ids[1] = 1, 2
	editSlice(ids)
	fmt.Println("slice 的元素被修改：", ids)  // [99 2]
	fmt.Println("调用方长度未改变：", len(ids)) // 2
	// editSlice 内的 append 没有扩容，但它只更新了形参的长度。
	// append 还可能分配新数组，所以需接收返回值：ids = append(ids, 3)。
	// 要得到独立元素副本，可用 make + copy；这仍不是对嵌套引用的深拷贝。
	cloned := make([]int, len(ids))
	copy(cloned, ids)
	cloned[0] = 7
	fmt.Println("独立副本 / 原切片：", cloned, ids) // [7 2] [99 2]

	// ═══════════════════════════════════════════════════════
	// 第 5 节：map 传参：共享映射与重新赋值
	// ═══════════════════════════════════════════════════════
	// values["theme"] = ... 修改共享映射；values = ... 则把局部变量指向新的映射。
	// 因此函数内最终打印 system，调用方读取的仍是 dark。

	// map 的值被复制后仍关联同一底层映射；复制 map 不会复制所有键值对。
	settings := map[string]string{"theme": "light"}
	editMap(settings)
	fmt.Println("共享 map 数据：", settings["theme"]) // dark
	// 函数中重新 make 一个 map 只替换形参，不会替换调用方的变量。
	// map 不是无需同步的“安全引用”：后续并发章节会解释共享数据的竞争。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：new 与 make：返回类型和初始化内容
	// ═══════════════════════════════════════════════════════
	// 先预测 new(map[string]int) 中哪一层为 nil，再进入 demonstrateAllocation 核对。
	// new 得到指针；make 初始化 slice、map、channel 这三类值。

	demonstrateAllocation()

	// ═══════════════════════════════════════════════════════
	// 第 7 节：指针比较与使用场景
	// ═══════════════════════════════════════════════════════
	// 相同内容不代表同一个变量，left == right 与 *left == *right 比较的对象不同。
	// 需要修改外部值或区分未设置状态时可使用指针。

	demonstratePointerEquality()
	// 指针不自动等于“更快”：它可能增加共享修改、逃逸和垃圾回收的成本。
	// 返回局部变量地址的安全性由语言保证；具体放在栈还是堆，由编译器分析决定。
	// 不要套用“使用 & 就一定在堆上”或“用了 new 就一定在堆上”的经验规则。
	// 第 11 章再结合方法解释：需要修改接收者时为何常用指针接收者。
}

func demonstrateAllocation() {
	// new(T) 返回 *T，指向 T 的零值。它不等于“把 T 的内部容器全部初始化”。
	mapPointer := new(map[string]int)
	fmt.Println("new(map) 的指针是否为 nil：", mapPointer == nil) // false
	fmt.Println("指针指向的 map 是否为 nil：", *mapPointer == nil)  // true
	// 此时 (*mapPointer)["views"] = 1 会 panic，因为底层 map 仍是 nil。
	*mapPointer = make(map[string]int)
	(*mapPointer)["views"] = 1
	fmt.Println("初始化后的 map：", (*mapPointer)["views"]) // 1

	// make 用于 slice/map/channel，返回已经初始化的该类型值，不是指向它的指针。
	// map/slice 本身可共享底层数据，通常直接传递它们即可，不必再包一层指针。
	items := make([]int, 2)
	fmt.Printf("make 返回 %T，new 返回 %T\n", items, mapPointer)
	// 当前仓库使用 Go 1.27，还支持 new(值表达式)：创建以该值初始化的变量并取地址。
	// new(int) 是零值 0；new(20) 则指向值为 20 的 int，常用于可选字段。
	pageSize := new(20)
	fmt.Println("new(值表达式)：", *pageSize) // 20
}

func demonstratePointerEquality() {
	first, second := 7, 7
	left, alias, right := &first, &first, &second
	fmt.Println("指向同一变量：", left == alias)                    // true
	fmt.Println("值相等但变量不同：", left == right, *left == *right) // false true
	// 指针可与 nil、同类型指针比较；比较的是指向的变量身份，不是内容深度相等。
	// 此示例使用 int。指向不同零大小变量（例如 struct{}）的指针不保证比较为 false。
}

func changeCopy(stock int) {
	stock = 0
}

func decreaseStock(stock *int, amount int) {
	// 接受指针时，要明确是否允许 nil；此处约定 nil 直接跳过。
	if stock == nil {
		return
	}
	*stock -= amount
}

func repoint(stock *int) {
	other := 100
	stock = &other
	// stock 只是调用方指针的副本；现在该副本指向 other。
}

func makeLimit(value int) *int {
	return &value
}

func editSlice(values []int) {
	values[0] = 99
	values = append(values, 3)
	fmt.Println("函数内 slice：", values) // [99 2 3]
}

func editMap(values map[string]string) {
	values["theme"] = "dark"
	values = map[string]string{"theme": "system"}
	fmt.Println("函数内新 map：", values["theme"]) // system
}

// 后续再学：结构体接收者（11 章）、共享数据的并发安全（17 章）；
// unsafe、逃逸分析细节和内存优化不是当前手写业务代码的前置要求。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 repoint 中的 stock = &other 改成 *stock = other，预测调用方的值，再恢复原例。
// 对比 nil 指针和 new(int)，解释为什么后者可以安全解引用。
