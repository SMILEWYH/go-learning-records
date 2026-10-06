/*
章节：06-分支选择
	本章掌握清单：
	1. 值 switch 与多个 case 值
	2. 初始化语句与 case 作用域
	3. 无表达式 switch：按条件分类
	4. break：提前退出当前 switch
	5. fallthrough：直接进入下一分支
	6. 表达式求值顺序与循环中的 switch
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./06-switch
// 运行练习：go run ./06-switch exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：Go 的 case 默认不会贯穿到下一分支，通常不用写 break。
// 后续再学：嵌套循环中的标签 break 见第 07 章；类型 switch 见第 13 章。

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
	// 第 1 节：值 switch 与多个 case 值
	// ═══════════════════════════════════════════════════════
	// switch 先得到待比较的值，再找匹配的 case。
	// PUT 和 PATCH 共用同一段操作，逗号表示多个候选值，不是同时执行两次。

	// 按值选择分支；只执行第一个匹配分支。
	method := "POST"
	switch method {
	case "GET":
		fmt.Println("读取资源")
	case "POST":
		fmt.Println("创建资源")
	case "PUT", "PATCH": // 多个值可以共用一个分支。
		fmt.Println("更新资源")
	case "DELETE":
		fmt.Println("删除资源")
	default:
		fmt.Println("暂不支持的方法")
	}
	// default 可省略；没有匹配且没有 default 时，switch 什么也不做。
	// switch 字符串匹配区分大小写，"get" 与 "GET" 不相等。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：初始化语句与 case 作用域
	// ═══════════════════════════════════════════════════════
	// 分号前定义的 mode 服务于整个 switch；每个 case 内又有独立作用域。
	// 两个 case 各自声明 message 合法，但都不能在 switch 结束后继续使用。

	// 与 if 类似，switch 可以有初始化语句。
	config := map[string]string{"mode": "preview"}
	switch mode := config["mode"]; mode {
	case "production":
		message := "加载生产配置"
		fmt.Println(message)
	case "preview", "development":
		message := "加载开发配置" // 每个 case 都是独立的隐式作用域，可以复用局部变量名。
		fmt.Println(message)
	default:
		fmt.Println("未知模式，使用默认配置")
	}
	// mode 只在这个 switch 内可用。
	// case 内的 message 只属于自己的分支，也不能在整个 switch 之后使用。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：无表达式 switch：按条件分类
	// ═══════════════════════════════════════════════════════
	// switch 后不写值时，case 后写 bool 条件，适合替代一串 else if。
	// 当范围重叠时，先写更具体的条件，否则宽泛条件会先匹配。

	// 省略 switch 后的表达式，等价于 switch true。
	// 它适合范围判断；重叠条件按书写顺序取第一个匹配项。
	status := 503
	switch {
	case status >= 200 && status < 300:
		fmt.Println("HTTP 类别：成功")
	case status >= 400 && status < 500:
		fmt.Println("HTTP 类别：客户端错误")
	case status >= 500 && status < 600:
		fmt.Println("HTTP 类别：服务端错误")
	default:
		fmt.Println("HTTP 类别：其他")
	}

	// ═══════════════════════════════════════════════════════
	// 第 4 节：break：提前退出当前 switch
	// ═══════════════════════════════════════════════════════
	// case 本来会在末尾自动结束；这里的 break 是为了跳过分支中后面的发布动作。
	// 执行完 switch 后，main 仍会继续打印下一行。

	// break 可以提前结束当前 switch 分支，不会结束整个函数。
	action := "publish"
	validated := false
	switch action {
	case "publish":
		if !validated {
			fmt.Println("暂不发布：内容还未验证")
			break
		}
		fmt.Println("执行发布")
	default:
		fmt.Println("没有发布操作")
	}
	fmt.Println("break 后继续执行当前函数")

	// ═══════════════════════════════════════════════════════
	// 第 5 节：fallthrough：直接进入下一分支
	// ═══════════════════════════════════════════════════════
	// n 为 1，却也打印 case 2 的内容，这是 fallthrough 主动要求的结果。
	// 它只前进到下一分支；该分支结束后，不会自动继续贯穿 default。

	// fallthrough 强制执行紧接着的下一分支，不会再判断它的条件。
	// 这里只作语法演示。大部分业务分支用多值 case 或明确的函数调用更清晰。
	n := 1
	switch n {
	case 1:
		fmt.Println("匹配到 case 1")
		fallthrough
	case 2:
		fmt.Println("也执行了 case 2 的内容，尽管 n 不等于 2")
	default:
		fmt.Println("本次不会执行 default")
	}
	// fallthrough 必须是该 case 的最后一个非空语句，不能用于最后一个 case。
	// 类型 switch 会在接口章节介绍；类型 switch 不能使用 fallthrough。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：表达式求值顺序与循环中的 switch
	// ═══════════════════════════════════════════════════════
	// 观察打印痕迹，区分“判断某个 case”和“执行某个 case 的代码”。
	// 匹配 B 之后不再计算 C；带副作用的 case 因此可能根本不会执行。

	// switch 的表达式只计算一次；case 表达式从上到下、从左到右按需计算。
	// case 不一定是常量，也可以是变量或函数结果，但必须能和 switch 表达式比较。
	// 本例的辅助函数只是打印求值痕迹；第 08 章会讲函数声明。
	switch traceSwitchValue("switch 表达式", 2) {
	case traceSwitchValue("case A", 1), traceSwitchValue("case B", 2):
		fmt.Println("在 case B 匹配，后面的 case C 不再计算")
	case traceSwitchValue("case C", 3):
		fmt.Println("本次不会执行")
	}
	// 会依次打印 switch 表达式、case A、case B，没有 case C。
	// 把可能修改状态的函数写进 case 会让行为难读；通常先算出待分类的值，再用纯条件分支。
	// switch 嵌在 for 内时，未加标签的 break 只退出最近的 switch，循环仍会继续。
}

func traceSwitchValue(label string, value int) int {
	fmt.Println("求值：", label)
	return value
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 method 改成 PATCH 和 patch，比较大小写匹配结果。
// 去掉 fallthrough 再运行，说明为何不需要补 break；把 n 改成 2 再预测输出。
