/*
章节：10-初始化与延迟执行
	本章掌握清单：
	1. 包初始化与 init
	2. defer 的登记与逆序执行
	3. 参数快照与闭包读取
	4. 返回值与 defer 的执行时机
	5. defer 属于函数，不属于循环块
	6. 方法接收者也会提前求值
	7. 显式初始化与清理边界
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./10-init-defer
// 运行练习：go run ./10-init-defer exam
// 阅读顺序：先读 main 的七个小节，再按调用名找到下方演示函数；上方 init 对应第 1 节。
// TS 对照：包初始化类似模块顶层初始化；defer 常承担 try/finally 的清理职责。
// defer 属于“当前函数”，不是“当前代码块”，也不是异步调度任务。

package main

import (
	"bytes"
	"fmt"
	"os"
)

// 初始化依赖包，再初始化当前包的包级变量，再执行当前包的 init。
// 之后才调用 main。包只初始化一次；依赖关系不是源文件的“嵌套”。
// 包变量初始化会考虑变量间的依赖，不应仅按肉眼看到的行序理解。
// 假设 main 导入 service，service 又导入 config：config 先于 service，service 先于 main。
// 这由包依赖关系决定，不由 import 行的排列顺序或目录深度决定；包不能循环导入。
// 多条依赖路径导入同一个包，不会让它重复初始化。空白导入 _ 也会触发包初始化。
var startupSteps = []string{"包级变量已初始化"}

func init() {
	startupSteps = append(startupSteps, "第一个 init")
}

func init() {
	startupSteps = append(startupSteps, "第二个 init")
}

// 同一文件中的 init 按出现顺序执行，可有多个，但不能直接调用或当值传递。
// init 必须没有参数、没有返回值。不同文件不应靠文件名维持业务初始化依赖；
// 源文件呈交编译器的次序影响其顺序，跨文件依赖宜写成显式初始化函数。
// init 适合简单注册；需要配置、错误处理的数据库连接更适合在 main 中显式建立。
// 为让 exam 模式只打印题目，这里的 init 仅记录顺序，不直接输出。

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：包初始化与 init
	// ═══════════════════════════════════════════════════════
	// 程序进入 main 前，startupSteps 已记录包变量和两个 init 的执行顺序。
	// 回看上方包级声明；init 自动执行，不能在 main 中手动调用。

	fmt.Println("初始化顺序：", startupSteps)

	// ═══════════════════════════════════════════════════════
	// 第 2 节：defer 的登记与逆序执行
	// ═══════════════════════════════════════════════════════
	// 运行到 defer 时先登记，当前函数返回时才调用。
	// 把登记过程想成依次压入一叠清理动作，最后登记的最先执行。

	demonstrateOrder()

	// ═══════════════════════════════════════════════════════
	// 第 3 节：参数快照与闭包读取
	// ═══════════════════════════════════════════════════════
	// 普通参数在登记时求值；闭包没有传入参数时，可以在执行时读取外层变量。
	// 进入 demonstrateCapture 前，先预测三个输出里哪一个会读到 done。

	demonstrateCapture()

	// ═══════════════════════════════════════════════════════
	// 第 4 节：返回值与 defer 的执行时机
	// ═══════════════════════════════════════════════════════
	// return 200 先把有名返回变量 code 设为 200，随后 defer 将它加一。
	// 调用方最终拿到 201；defer 不是在调用方收到结果之后才执行。

	fmt.Println("有名返回值：", responseCode()) // 201

	// ═══════════════════════════════════════════════════════
	// 第 5 节：defer 属于函数，不属于循环块
	// ═══════════════════════════════════════════════════════
	// 循环结束时，外层函数仍未返回，三次登记的动作都还在等待。
	// 循环开文件时可把单次处理提成函数，让每一轮及时关闭资源。

	demonstrateScope()

	// ═══════════════════════════════════════════════════════
	// 第 6 节：方法接收者也会提前求值
	// ═══════════════════════════════════════════════════════
	// defer current.WriteString(...) 保存登记时的 current 指针。
	// 后来 current 改指向 second，并不会改变已登记调用的目标。

	demonstrateDeferredReceiver()

	// ═══════════════════════════════════════════════════════
	// 第 7 节：显式初始化与清理边界
	// ═══════════════════════════════════════════════════════
	// 把配置作为参数传入，再检查校验结果，调用顺序和失败处理就都可见。
	// 显式初始化适合可能失败的工作，init 可保留给简单注册。

	mode, ok := prepareMode("development")
	fmt.Println("显式初始化得到模式：", mode, ok) // development true
	// 这里显式传入配置、检查结果，方便测试或替换；无需 init 偷读全局环境或启动服务。

	// defer 会在正常返回和本 goroutine 的 panic 栈展开过程中执行。
	// os.Exit 会直接结束进程，不执行 defer；main 返回时也不会等待其他 goroutine。
	// 后续文件读写章会将 defer 用在成功打开资源之后，负责关闭资源。
}

func demonstrateDeferredReceiver() {
	// 先把 bytes.Buffer 看成可追加文本的对象；下一章会系统介绍方法与接收者。
	first := bytes.NewBufferString("first")
	current := first
	defer func() {
		fmt.Println("defer 保存的接收者 / 后来的接收者：", first.String(), current.String())
	}()
	// 除参数外，方法的接收者也在登记时求值。这里保存的是 first 对应的指针。
	defer current.WriteString(":finished")
	current = bytes.NewBufferString("second")
	// 退出时先追加到 first，再打印 first:finished / second。
	// 如果改为 defer func(){ current.WriteString(":finished") }()，闭包会读取后来的 current。
}

func prepareMode(raw string) (string, bool) {
	switch raw {
	case "development", "production":
		return raw, true
	default:
		return "", false
	}
	// bool 是现阶段的简单校验约定；第 18 章会用 error 表达具体失败原因。
}

func demonstrateOrder() {
	fmt.Println("\n开始处理请求")
	defer fmt.Println("最后：记录请求结束")
	defer fmt.Println("其次：释放模拟连接")
	fmt.Println("现在：执行业务")
	// 输出业务 -> 释放连接 -> 记录结束。
	// 准确规则是“实际登记的逆序”，不是源文件中“离 return 最近”。
	// 没有执行到的 defer 不会登记；提前 return 也会执行已经登记的 defer。
}

func demonstrateCapture() {
	fmt.Println("\n参数求值和闭包：")
	status := "pending"

	// 普通 defer 调用在登记时就计算函数值与参数。这里只保存了 "pending"。
	defer fmt.Println("登记时的字符串：", status)

	// 闭包本体到退出时才执行，此处读取的是同一个变量的最终值。
	defer func() {
		fmt.Println("退出时的字符串：", status)
	}()

	// 显式给闭包传参也遵循“参数立即求值”，会得到 pending。
	defer func(snapshot string) {
		fmt.Println("闭包参数快照：", snapshot)
	}(status)
	status = "done"
	// 按 LIFO 输出：闭包参数 pending -> 闭包读取 done -> 普通调用 pending。
	// 若参数是指针或 slice，“立即求值”不意味着深拷贝其指向的数据。
}

func responseCode() (code int) {
	// return 先确定返回值，然后执行 defer，最后把结果交给调用方。
	// 有名返回变量可被闭包修改。日常代码尽量避免靠这种技巧隐藏修改。
	defer func() {
		code++
	}()
	return 200
}

func demonstrateScope() {
	fmt.Println("\n循环中登记 defer：")
	for i := 0; i < 3; i++ {
		defer fmt.Println("函数退出后才执行：", i)
	}
	fmt.Println("循环已经结束，defer 尚未执行")
	// 随后输出 2、1、0。若循环里打开许多文件，资源可能一直累积到函数结束。
	// 更合适的做法通常是把“一次打开、处理、关闭”放入独立的小函数。
}

// 后续再学：方法与接收者（11 章）、panic/recover（18 章）、真实文件资源的关闭（20 章起）。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 画出 demonstrateCapture 的登记顺序，再倒过来写出实际输出。
// 把 status 修改语句移到 defer 前，比较参数快照与闭包输出。
