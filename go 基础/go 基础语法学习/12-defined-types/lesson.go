/*
章节：12-自定义类型与别名
	本章掌握清单：
	1. 定义类型：相同表示，不同身份
	2. 显式转换与未类型化常量
	3. 类型别名：同一类型的另一个名字
	4. 新类型的方法与转换
	5. iota 枚举与合法值校验
	6. 转换不会代替范围检查
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./12-defined-types
// 运行练习：go run ./12-defined-types exam
// 阅读顺序：上方集中声明领域类型和方法；从 main 的赋值与调用开始，再回看对应类型。
// TS 对照：type UserID = string 在 TS 中只是别名；Go 的 type UserID string 是新类型。
// Go 的定义类型可减少误传参数，但转换本身并不会完成业务校验。

package main

import (
	"fmt"
	"os"
)

// 创建两个不同的定义类型，底层类型都为 string。
type UserID string
type OrderID string

func (id UserID) IsEmpty() bool {
	return id == ""
}

func loadProfile(id UserID) string {
	return "profile:" + string(id)
}

// 新名字 RawID 和 string 指代同一类型，没有建立新的类型身份。
type RawID = string

// 别名也能指向我们在当前包定义的类型，并保留该类型的方法。
type LegacyUserID = UserID

// “别名绝对不能定义方法”不准确。
// 此别名指向同包定义的非泛型 UserID，因此允许把方法写在该别名接收者上；
// 方法实际仍属于同一个 UserID 类型，并没有创造一份只属于别名的方法集。
func (id LegacyUserID) APIValue() string {
	return string(id)
}

// 若写 func (id RawID) ... 则不行：RawID 指向内建 string，string 不是本包定义的类型。
// 对其他包类型建立别名也不会让你获得向那个外部类型追加方法的权限。
// 接收者基类型必须符合“本包定义、非指针/接口基类型”等规则；泛型别名还有额外限制。

type StatusCode int

const (
	StatusOK       StatusCode = 200
	StatusNotFound StatusCode = 404
)

func (code StatusCode) IsSuccess() bool {
	return code >= 200 && code < 300
}

// 新定义类型保留底层表示，但不会自动继承原定义类型的方法。
type InternalCode StatusCode

// iota 在同一个 const 块中按常量声明行从 0 递增。省略表达式时沿用上一行。
// 保留 0 表示未知状态，能让尚未初始化的值更容易识别。
type JobState uint8

const (
	JobUnknown JobState = iota
	JobQueued
	JobRunning
	JobFinished
)

func (state JobState) IsKnown() bool {
	switch state {
	case JobQueued, JobRunning, JobFinished:
		return true
	default:
		return false
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：定义类型：相同表示，不同身份
	// ═══════════════════════════════════════════════════════
	// type UserID string 创建新类型，能让编译器发现把订单编号当用户编号传入的错误。
	// 底层都是 string 并不表示可以直接互相赋值。

	userID := UserID("user-7")
	orderID := OrderID("order-7")
	fmt.Println("合法用户查询：", loadProfile(userID))
	fmt.Printf("不同类型：%T / %T\n", userID, orderID) // main.UserID / main.OrderID
	// loadProfile(orderID) 无法编译，两个定义类型不能直接互相赋值。
	// 不要随意 UserID(orderID) 绕过边界；编译通过并不代表业务语义正确。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：显式转换与未类型化常量
	// ═══════════════════════════════════════════════════════
	// 变量 plain 已是 string，需要 UserID(plain)；字符串常量能适配 UserID。
	// 这里的转换只改变类型，不校验编号格式。

	// 有类型的 string 变量需要显式转换；无类型字符串常量可直接赋给 UserID。
	plain := "user-8"
	converted := UserID(plain)
	var literal UserID = "user-9"
	fmt.Println("显式转换与常量：", converted, literal)
	fmt.Println("空 ID：", UserID("").IsEmpty()) // true
	// UserID("!!!") 也能成功转换，校验格式需要额外函数。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：类型别名：同一类型的另一个名字
	// ═══════════════════════════════════════════════════════
	// type RawID = string 中的等号决定它是别名；去掉等号就变成定义类型。
	// LegacyUserID 与 UserID 身份相同，方法也属于同一类型。

	var alias RawID = plain
	var sameString string = alias                      // 同一类型，不需要转换。
	fmt.Printf("string 别名：%T %s\n", alias, sameString) // string user-8
	var legacy LegacyUserID = userID
	fmt.Println("本地类型别名保留方法：", legacy.IsEmpty(), legacy.APIValue())
	fmt.Println("原类型也具有该方法：", userID.APIValue())

	// ═══════════════════════════════════════════════════════
	// 第 4 节：新类型的方法与转换
	// ═══════════════════════════════════════════════════════
	// StatusCode 有 IsSuccess 方法，InternalCode 是另外定义的类型，不会自动得到它。
	// 把 internal 显式转回 StatusCode，才能调用原类型的方法。

	code := StatusOK
	fmt.Println("状态码成功：", code.IsSuccess())            // true
	fmt.Println("404 成功：", StatusNotFound.IsSuccess()) // false
	internal := InternalCode(code)
	fmt.Println("转回原定义类型后调用方法：", StatusCode(internal).IsSuccess())
	// internal.IsSuccess() 无法编译：定义新类型不会继承 StatusCode 的方法。

	// 常见选择：
	// - 领域上不同的值（UserID/OrderID、金额/数量）用定义类型表达差异。
	// - 迁移公开 API 时保留旧名字，可能用别名维护类型兼容性。
	// - 只想缩短局部变量名，不需要为此建立新的类型。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：iota 枚举与合法值校验
	// ═══════════════════════════════════════════════════════
	// 枚举常量提供易读名字，但变量仍可保存未命名的数值。
	// 本例把 0 留给未知状态，并用 IsKnown 判断业务允许的状态。

	// iota 给一组状态起名，并不会限制变量只能取这些命名常量。
	fmt.Println("任务状态枚举：", JobUnknown, JobQueued, JobRunning, JobFinished) // 0 1 2 3
	fmt.Println("命名状态是否合法：", JobRunning.IsKnown())                         // true
	fmt.Println("未命名值是否合法：", JobState(99).IsKnown())                       // false
	// Go 的这类枚举不是封闭联合类型，不会自动要求 switch 穷尽所有状态。
	// 存入数据库或公开 API 的数值最好显式约定，随意插入 iota 声明会改变后续编号。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：转换不会代替范围检查
	// ═══════════════════════════════════════════════════════
	// 常量 uint8(300) 会被编译器拒绝，变量窄化转换却会按规则截断。
	// 处理外部输入时先检查原始数值，再转成较窄的领域类型。

	// 数值转换也不是校验：变量的窄化转换会截断高位，而不是自动报错。
	wide := 300
	fmt.Println("变量转为 uint8：", uint8(wide)) // 44
	// uint8(300) 这个常量转换却无法编译，因为编译器知道常量超出了 uint8 范围。
	// 若外部输入要变成 JobState，应在原始较宽类型上检查范围和允许值，再转换。
	decimal := 3.9
	fmt.Println("浮点转整数舍去小数：", int(decimal)) // 3，不是四舍五入。
	// 字符串转换也需区分：string(65) 表示 Unicode 码点对应的 "A"，不是文本 "65"。
	// 数字的十进制文本用 strconv.Itoa/FormatInt；文本转数字用 Atoi/ParseInt 并检查错误。
}

// 后续再学：类型断言与转换的区别（13 章）、泛型约束与泛型别名（19 章）；
// 不必为每个原始类型都包一层，新类型应表达明确的业务差别。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 比较 type RawID = string 与 type UserID string，分别尝试赋给 string 变量。
// 把 JobState(99) 换成 JobUnknown，解释零值为什么也未通过 IsKnown。
