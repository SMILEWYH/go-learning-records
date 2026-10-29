/*
章节:08-函数与闭包
*/

package main

import (
	"fmt"
	"strings"
)

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节:声明函数:形参与单个返回值的类型声明
	// ═══════════════════════════════════════════════════════
	// Ts 需要使用 :  而 Go 不需要
	var fun1 = func(name string, age int) bool {
		if name == "vvyh" && age == 26 {
			return true
		}
		return false
	}

	var isV = fun1("vvyh", 26)
	fmt.Print(isV)

	// ═══════════════════════════════════════════════════════
	// 第 2 节:多个返回值与空白标识符
	// ═══════════════════════════════════════════════════════
	// Ts 没有这种多个返回值的写法, 但可以使用返回数组和数组结构来代替
	var fun2 = func(name string, age int) (string, int, string) {
		return name, age, "无用消息"
	}

	var name, age, _ = fun2("vvyh", 26)
	fmt.Print(name, age)

	// ═══════════════════════════════════════════════════════
	// 第 3 节:命名返回值
	// ═══════════════════════════════════════════════════════
	// 这样可以直接 return, 不需要在后面显式写出要返回的值
	var fun3 = func(name string, age int) (nextAge int, msg string) {
		nextAge = age + 1
		msg = "无用消息"
		return
	}

	var nextAge, _ = fun3("vvyh", 26)
	fmt.Print(nextAge)

	// ═══════════════════════════════════════════════════════
	// 第 4 节:变参:零个、多个与切片展开
	// ═══════════════════════════════════════════════════════
	// values ...int 像 Ts 的收集参数,  counts... 像 Ts 的数组展开(只不过冒号在前面)
	var fun4 = func(values ...int) int { // 变参参数在函数内部是切片;变参只能有一个且必须是最后一个参数
		total := 0
		for _, value := range values {
			total += value
		}
		return total
	}

	fmt.Println("零个参数：", fun4())
	fmt.Println("多个参数：", fun4(2, 3, 5))
	counts := []int{4, 6, 8}
	fmt.Println("展开切片：", fun4(counts...)) // counts... 展开不意味着复制底层数组;若函数改动参数元素, 调用方也可能受影响

	var fun5 = func(name string, values ...int) int { // 变参前还可以携带多个参数
		total := 0
		for _, value := range values {
			total += value
		}
		return total
	}

	fmt.Println("展开切片：", fun5("vvyh", counts...))

	// ═══════════════════════════════════════════════════════
	// 第 5 节:函数值、签名与 nil
	// ═══════════════════════════════════════════════════════
	// Go 不允许在函数体里再写 func named() {} 的具名声明, 只允许用匿名函数赋给变量
	// 并且函数也是值, 可用 var 或 := 保存匿名函数, 不能使用 const
	fun6 := func(a, b int) int {
		return a + b
	}
	var combine func(int, int) int = fun6 // 再保存到另一个变量中
	fmt.Println("函数变量：", combine(3, 7))

	// 函数类型的零值是 nil;调用 nil 函数会 panic, 所以使用前应确保已经赋值
	var emptyFunc func(string) string
	if emptyFunc == nil {
		emptyFunc = strings.TrimSpace
	}
	fmt.Printf("默认回调：%q\n", emptyFunc(" ready "))
	// 函数值只能和 nil 比较, 不能用 == 判断两个函数是否相同, 也不能直接作为 map 的 key

	// 函数类型由参数和返回值的个数、顺序、类型和是否存在变参决定, 参数名字不影响类型
	// func(int, int) int 与 func(a int, b int) int 相同
	// func(...int) int 与 func([]int) int 不同

	// ═══════════════════════════════════════════════════════
	// 第 6 节:闭包:保存并修改外层变量
	// ═══════════════════════════════════════════════════════
	// 与 Ts 的闭包一样
	func7 := func(start int) func() int {
		var next = start
		return func() int {
			next++
			return next
		}
	}

	nextA := func7(100)
	nextB := func7(0)
	fmt.Println("独立的序列 A：", nextA(), nextA())
	fmt.Println("独立的序列 B：", nextB(), nextB())

	// ═══════════════════════════════════════════════════════
	// 第 7 节:函数可以当另一个函数的实参、返回值
	// ═══════════════════════════════════════════════════════

	// ═══════════════════════════════════════════════════════
	// 第 8 节:func() { ... } 后面携带 () 表示立即执行函数
	// ═══════════════════════════════════════════════════════

	// ═══════════════════════════════════════════════════════
	// 第 9 节:循环闭包:独立变量与复用变量
	// ═══════════════════════════════════════════════════════
	// Go 1.22 起, 用 := 声明的循环变量每轮独立, 闭包按各自那一轮的变量读取
	// 与 Ts 的 let 变量写的 for 循环类似
	callbacks := []func() int{}
	for index := range 3 {
		callbacks = append(callbacks, func() int { return index })
	}
	fmt.Println("每轮独立的变量：", callbacks[0](), callbacks[1](), callbacks[2]()) // 0 1 2

	// 与 Ts 的 var 变量写的 for 循环类似
	var reused int
	callbacks = nil
	for reused = range 3 {
		callbacks = append(callbacks, func() int { return reused })
	}
	fmt.Println("复用外层同一变量：", callbacks[0](), callbacks[1](), callbacks[2]()) // 2 2 2
	// 旧版本 Go 常见的“闭包总取到最后一个循环值”不能不分版本、语法地照搬

	// ═══════════════════════════════════════════════════════
	// 第 10 节:递归:基础情况与递进过程
	// ═══════════════════════════════════════════════════════
	var func8 func(num int) int

	func8 = func(num int) int {
		if num < 1 {
			return 1
		}

		return num * func8(num-1)
	}
}
