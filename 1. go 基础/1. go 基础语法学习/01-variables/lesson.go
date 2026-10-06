/*
章节：01-变量与常量
	本章掌握清单：
	1. 四种声明方式
	2. 零值规则
	3. 多变量与批量声明
	4. 多返回值赋值
	5. 常量与 iota（含 untyped constant）
	6. 作用域与遮蔽（含 if 中陷阱）
	7. 包导出与模块概念
*/

// 请从项目根目录运行以下两个命令，两个文件共同组成同一个 main 包。
// 运行讲解：go run ./01-variables
// 运行练习：go run ./01-variables exam

package main

import (
	"fmt"
	"os"
	"strconv"
)

// ──────────────────────────────────────────────────────────
// 下面三行是"包级声明"：写在函数外面，整个包内都能访问。
// 初学者先跳过这里，到第 7 节会详细讲解包级 vs 函数级的区别。
// ──────────────────────────────────────────────────────────
var serviceName = "account-api" // 包级变量，有初始值。
var deploymentRegion string     // 包级变量，无初始值（零值 ""）。
const defaultPort = 8080        // 包级常量。

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：变量声明的四种方式
	// ═══════════════════════════════════════════════════════
	// Go 有四种常见的声明局部变量的写法，适用于不同场景。

	// 方式一：var + 初始值，编译器自动推断类型。
	var retries = 3 // 推断为 int。

	// 方式二：var + 显式类型 + 初始值，明确指定类型。
	var owner string = "web"

	// 方式三：var + 显式类型，不给初始值，变量自动获得该类型的零值。
	var ready bool // 零值是 false。

	// 方式四：:= 短变量声明，最常用，但只能在函数体内使用。
	port := 8081 // 等价于 var port = 8081。

	// 已存在的变量，用 = 赋新值，不是 :=。
	ready = true

	fmt.Println("服务配置：", retries, owner, ready, port)

	// 类型一旦确定，就不能赋给不同类型的值。
	// 例如 port 是 int，不能写 port = "8081"。

	// 局部变量声明后必须被使用，否则编译报错；包级变量没有这个限制。
	// _ 是空白标识符，可以明确丢弃不需要的值。
	_ = deploymentRegion

	// ═══════════════════════════════════════════════════════
	// 第 2 节：零值规则 —— Go 不存在 undefined
	// ═══════════════════════════════════════════════════════
	// 所有变量声明后如果没有赋初值，都会自动获得"零值"。
	// 不会出现像 JS/TS 那样的 undefined 或 null。

	var zeroInt int       // 0
	var zeroFloat float64 // 0.0
	var zeroBool bool     // false
	var zeroString string // ""（空字符串）
	// 其他零值：指针 → nil，切片/map → nil，接口 → nil。
	// 这些类型会在后续章节中学习。

	fmt.Printf("零值演示：int=%d, float64=%g, bool=%t, string=%q\n",
		zeroInt, zeroFloat, zeroBool, zeroString)
	// 注意：0 和 false 既是零值，也是合法的业务值。
	// 仅凭 enabled == false 不能区分"用户主动关闭"还是"用户没有设置"。
	// 在需要区分时，后续章节会学习使用指针或特殊标记来处理。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：多变量声明与批量声明
	// ═══════════════════════════════════════════════════════
	// 一行同时声明多个同类型变量：
	var host, protocol string = "localhost", "http"

	// 用 var (...) 把相关配置分组在一起，更整洁：
	var (
		maxConnections = 100
		debug          = true
	)
	fmt.Println("地址：", protocol, host, "连接上限：", maxConnections, "调试：", debug)

	// := 也可以一行声明多个：
	width, height := 1920, 1080
	fmt.Println("分辨率：", width, "x", height)

	// 多变量赋值先计算右侧所有值，再统一赋给左侧，因此可以直接交换：
	primary, fallback := "node-a", "node-b"
	primary, fallback = fallback, primary
	fmt.Println("主节点与备用节点交换：", primary, fallback)

	// ═══════════════════════════════════════════════════════
	// 第 4 节：多返回值赋值与 := 复用规则（重要！）
	// ═══════════════════════════════════════════════════════
	// Go 的函数经常返回多个值（值 + error），这在后续章节会大量使用。
	// 这里先用 strconv.Atoi 展示如何接收多返回值。

	// 用 := 接收两个返回值：
	num, err := strconv.Atoi("42")
	fmt.Println("解析结果：", num, "错误：", err)

	// := 复用规则：左侧至少有一个新变量，其余可以是已有变量。
	// 下面 err 已存在，ratio 是新变量，所以合法；err 被重新赋值而不是创建新变量。
	ratio, err := strconv.ParseFloat("3.14", 64)
	fmt.Println("解析浮点：", ratio, "错误：", err)

	// 同理，port 已在第 1 节声明，path 是新变量：
	port, path := 8082, "/health"
	fmt.Println("健康检查：", host, port, path)

	// 如果 := 左侧全部是已有变量（没有新变量），编译器会报错。
	// 此时应该用 = 而不是 :=。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：常量、iota 与未类型化常量
	// ═══════════════════════════════════════════════════════
	// 常量用 const 声明，值必须在编译期确定。
	// 不能把切片 []string{}、map、函数等赋给 const。
	const secondsPerMinute = 60
	const cacheTTL = 5 * secondsPerMinute // 常量之间可以运算。
	fmt.Println("默认端口：", defaultPort, "缓存秒数：", cacheTTL)

	// ── 未类型化常量 vs 有类型常量 ──
	// 没有指定类型的常量叫"未类型化常量"（untyped constant），
	// 它在赋值时可以自动适配目标类型，使用更灵活。
	const untypedVal = 100           // 未类型化，可以赋给 int、int64、float64 等。
	const typedVal int = 100         // 有类型（int），只能赋给 int。
	var asInt64 int64 = untypedVal   // ✓ 未类型化常量自动适配 int64。
	var asFloat float64 = untypedVal // ✓ 未类型化常量自动适配 float64。
	// var asInt64b int64 = typedVal  // ✗ 编译错误：int 不能隐式转为 int64。
	fmt.Println("未类型化常量适配：", asInt64, asFloat)

	// ── iota：枚举计数器 ──
	// iota 在 const 声明组内从 0 开始，每一行自增 1。
	// 省略后续表达式时，会复用上一行的表达式。
	const (
		levelInfo    = iota // 0
		levelWarning        // 1
		levelError          // 2
	)
	fmt.Println("日志级别枚举：", levelInfo, levelWarning, levelError)

	// iota 可以配合位移产生位标记：1 << iota → 1, 2, 4, 8...
	// 详细的位运算将在第 03 章学习。

	// 注意：在枚举中间插入新项会改变后续项的数值，
	// 如果数值会持久化（存数据库）或对外传输（API），应明确指定数值。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：作用域与变量遮蔽（shadowing）
	// ═══════════════════════════════════════════════════════
	// 花括号 {} 创建内层作用域。
	// 内层用 := 声明同名变量会"遮蔽"外层变量（创建了一个全新的变量）。

	environment := "production"
	{
		// 这里的 := 创建了一个新的 environment，只在这个块内可见。
		environment := "development"
		fmt.Println("块内环境：", environment)
	}
	fmt.Println("块外环境仍为：", environment) // 仍然是 production。

	// 如果用 = 赋值，则修改的是外层的变量：
	{
		environment = "staging"
	}
	fmt.Println("块内用 = 修改后的外层环境：", environment) // 变成了 staging。

	// ── 常见陷阱：if 语句中的 := 遮蔽 ──
	// if 的初始化语句也会创建新的作用域，这是最容易踩坑的地方。
	// 例如：
	//   var err error
	//   if err := doSomething(); err != nil {
	//       // 这里的 err 是 if 块内的新变量，不是外面的 err！
	//   }
	//   // 外面的 err 仍然是 nil，即使 doSomething() 出了错。
	// 修正方法：在 if 内用 = 而不是 :=，或者用不同的变量名。
	// 实际示例将在第 05 章（if 语句）和第 18 章（错误处理）中学习。

	// ═══════════════════════════════════════════════════════
	// 第 7 节：包级声明、导出规则与 module/package
	// ═══════════════════════════════════════════════════════

	// ── 包级声明 vs 函数级声明 ──
	// 还记得文件顶部写在 func main() 外面的三行声明吗？
	// 那就是"包级声明"。现在回头看看它们和函数内声明的区别：
	//
	//   对比项          │ 包级（函数外）              │ 函数级（函数内）
	//   ──────────────│───────────────────────────│──────────────────────
	//   声明方式        │ 只能用 var / const          │ 可用 var / const / :=
	//   作用范围        │ 整个包内所有文件都能访问       │ 仅当前函数（或代码块）
	//   是否必须使用     │ 不强制（声明了可以不用）       │ 必须使用，否则编译报错
	//   典型用途        │ 配置、全局状态、导出给其他包    │ 函数内的临时计算、局部状态
	//
	// 回顾顶部声明的含义：
	//   var serviceName = "account-api"  → 包级变量，有初始值，类型推断为 string。
	//   var deploymentRegion string      → 包级变量，无初始值，零值 ""。
	//   const defaultPort = 8080         → 包级常量，编译期确定。
	// 包级声明不能使用 :=，因为 := 是"短变量声明"语法，只在函数体内有效。

	// ── 导出规则 ──
	// Go 没有 export 关键字，可见性完全由首字母大小写决定：
	// - 首字母大写（如 Version）→ 导出，其他包可以访问。
	// - 首字母小写（如 serviceName）→ 未导出，仅当前包内可见。
	// 这个规则适用于变量、常量、函数、类型及其字段，文件名不影响可见性。

	// ── module 与 package ──
	// 这两个概念容易混淆，但其实职责不同：
	// - module（模块）：由 go.mod 定义，是依赖管理和版本控制的单位。
	// - package（包）：是源码组织和编译的单位，通常一个目录是一个包。
	// 一个模块可以包含多个包，导入路径 = 模块路径 + 包的相对目录路径。
	// 包名由 package 声明决定，惯例与目录名相同，但不是强制要求。

	// 本项目用 go.work 把各章节作为独立模块加入工作区。
	// 跨包访问写法示意：在其他包中声明 const Version = "v1.0.0"，导入后使用 example.Version。
	// 本章独立运行，不导入课程其他目录；可在练习项目中另建包体验导出规则。
	fmt.Println("服务：", serviceName)
	fmt.Printf("尚未配置的区域：%q\n", deploymentRegion)

	// ── new() 函数（预告） ──
	// Go 还有一个 new(T) 内建函数，它分配一个 T 类型的零值并返回指向它的指针。
	// 例如 p := new(int)，此时 *p == 0。
	// 这是声明变量的另一种方式，但在实际代码中并不常用。
	// 详细内容会在第 09 章（指针）中学习。
	p := new(int)
	fmt.Printf("new(int) → 指针=%p, 值=%d\n", p, *p)

	// 学习方式：运行后试着修改各个变量，观察输出和编译器报错。
	// 比如：把一个局部变量声明后完全不使用，编译器会报错帮你找到问题。
	// 再比如：试试把 := 改成 =，或者把 = 改成 :=，理解两者的区别。
}
