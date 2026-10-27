/*
章节:10-初始化与延迟执行
*/

// TS 对照:包初始化类似模块顶层初始化;defer 常承担 try/finally 的清理职责
// defer 属于“当前函数”, 不是“当前代码块”, 也不是异步调度任务

package main

import (
	"bytes"
	"fmt"
)

// Go 程序启动顺序:初始化依赖包 → 初始化 main 包的包级变量 → 执行 main 包的 init() → 执行 main()
// 每个包的初始化顺序:先初始化其依赖包, 再初始化自身包级变量, 最后执行自身的 init()
// 假设 main 导入 service, service 又导入 config : config 先于 service, service 先于 main

var startupSteps = []string{"包级变量已初始化"}

// 包内可以存在多个 init 初始化函数, 并且它们按照从上到下顺序执行
// 并且它们必须没有参数、没有返回值, 也不可以直接调用或当值传递
// init 适合简单注册;需要配置、错误处理的数据库连接更适合在 main 中显式建立
func init() {
	startupSteps = append(startupSteps, "第一个 init 初始化函数")
}

func init() {
	startupSteps = append(startupSteps, "第二个 init 初始换函数")
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节:包初始化与 init
	// ═══════════════════════════════════════════════════════
	// "包级变量已初始化"、"第一个 init 初始化函数"、"第二个 init 初始换函数"
	fmt.Println("初始化顺序：", startupSteps)

	// ═══════════════════════════════════════════════════════
	// 第 2 节:defer 的登记与逆序执行
	// ═══════════════════════════════════════════════════════
	// 作用类似 React.useEffect 的清理函数用来清理副作用
	// defer 在单个函数体中可以存在多个, 可以写在函数体内的不同位置, 但只有执行到时才会注册, 并且按压入栈的顺序执行 (后进先出)
	func1 := func() {
		fmt.Println("获取 DOM 元素")
		defer fmt.Println("second: 最后执行的清理")
		defer fmt.Println("frist: 给 DOM 元素解绑点击事件")
		fmt.Println("给 DOM 元素绑定点击事件")
	}

	func1()

	// ═══════════════════════════════════════════════════════
	// 第 3 节:返回值与 defer 的执行时机
	// ═══════════════════════════════════════════════════════
	func3 := func() (code int) {
		// defer 会在 return 把结果返回给调用方之前执行
		// 有名返回变量可被闭包修改, 日常代码尽量避免靠这种技巧隐藏修改
		defer func() {
			code++
		}()
		return 200
	}

	func3()

	// ═══════════════════════════════════════════════════════
	// 第 4 节:defer 属于函数, 不属于 {} 代码块
	// ═══════════════════════════════════════════════════════
	func4 := func() {
		for i := 0; i < 3; i++ {
			defer fmt.Println("函数退出后才执行: ", i)
		}
		// 随后输出 2、1、0
		// 更合适的做法通常是把"业务逻辑、清理"封装为独立的小函数放到循环体内消费
	}

	func4()

	// ═══════════════════════════════════════════════════════
	// 第 5 节:defer 普通调用与闭包保存变量、指针快照的规则
	// ═══════════════════════════════════════════════════════
	func5 := func() {
		status := "pending"

		// defer 普通调用在登记时就保存当前变量的快照 "pending", 但注意这时候还未执行
		defer fmt.Println("登记时的字符串：", status)

		// defer 闭包在登记时不会保存当前变量的快照, 等执行时候才去获取变量 "done"
		defer func() {
			fmt.Println("退出时的字符串：", status)
		}()

		// 显式给 defer 闭包传参也遵循“参数立即求值”, 所以这里会保存变量的快照 "pending", 但注意这时候还未执行
		defer func(snapshot string) {
			fmt.Println("闭包参数快照：", snapshot)
		}(status)

		status = "done"
		// 当前只举例了普通类型变量, 试一试推论下指针会怎么样?
	}

	func5()

	func6 := func() {
		first := bytes.NewBufferString("first")
		current := first

		// defer 闭包在登记时不会保存当前指针的快照, 等执行时候才去获取指针
		defer func() {
			fmt.Println("defer 保存的接收者 / 后来的接收者：", first.String(), current.String())
		}()

		// defer 普通调用在登记时就保存当前指针的快照(指向 first), 但注意这时候还未执行
		defer current.WriteString(":finished")

		current = bytes.NewBufferString("second")
	}

	func6()
}
