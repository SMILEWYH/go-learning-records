/*
章节：12-自定义类型与别名
*/

// TS 对照：type UserID = string 在 TS 中只是别名; Go 的 type UserID string 是新类型
// Go 的定义类型可减少误传参数, 但转换本身并不会完成业务校验。

package main

import (
	"fmt"
)

// 创建新类型 UserID、OrderID, 底层类型都为 string。
type UserID string
type OrderID string

// 给 UserID 类型添加函数方法, 不只是结构体才可以添加函数方法
func (id UserID) IsEmpty() bool {
	return id == ""
}

// 别名, 没有建立新的类型身份
type RawID = string

// 别名也能指向我们在当前包定义的新类型, 并保留该类型的方法
type LegacyUserID = UserID

// "别名绝对不能定义方法" 不准确
// 此别名指向同包定义的新类型 UserID，因此允许把方法写在该别名接收者上；
// 方法实际仍属于同一个 UserID 类型，并没有创造一份只属于别名的方法集。
func (id LegacyUserID) APIValue() string {
	return string(id)
}

// 若写 func (id RawID) ... 则不行: RawID 指向内建 string, string 不是本包定义的类型
// 对其他包类型建立别名也不会让你获得向那个外部类型追加方法的权限
// (下面这句话还未全部搞懂)
// 接收者基类型必须符合“本包定义、非指针/接口基类型”等规则；泛型别名还有额外限制

type StatusCode int

const (
	StatusOK       StatusCode = 200
	StatusNotFound StatusCode = 404
)

func (code StatusCode) IsSuccess() bool {
	return code >= 200 && code < 300
}

// 新定义类型保留底层表示，但不会自动继承原定义类型的方法
type InternalCode StatusCode

// iota 在同一个 const 块中按常量声明行从 0 递增, 省略表达式时沿用上一行。
// 保留 0 表示未知状态，能让尚未初始化的值更容易识别
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
	// ═══════════════════════════════════════════════════════
	// 第 1 节：定义类型：相同表示, 不同身份
	// ═══════════════════════════════════════════════════════
	// type UserID string 创建新类型, 能让编译器发现把订单编号当用户编号传入的错误
	// 底层都是 string 并不表示可以直接互相赋值
	loadProfile := func(id UserID) string {
		return "profile:" + string(id)
	}

	// 如何创建新类型的实例:
	userID := UserID("user-7")
	orderID := OrderID("order-7")

	fmt.Println("合法用户查询：", loadProfile(userID))
	// loadProfile(orderID) 无法编译, 两个定义类型不能直接互相赋值
	// 不要随意 UserID(orderID) 绕过边界; 编译通过并不代表业务语义正确。
	fmt.Printf("不同类型：%T / %T\n", userID, orderID) // main.UserID / main.OrderID

	// ═══════════════════════════════════════════════════════
	// 第 2 节：显式转换与未类型化常量
	// ═══════════════════════════════════════════════════════
	// 有类型的 string 变量需要显式转换; 无类型字符串常量可直接赋给 UserID
	plain := "user-8" //
	converted := UserID(plain)
	var literal UserID = "user-9"
	fmt.Println("显式转换与常量：", converted, literal)
	fmt.Println("空 ID：", UserID("").IsEmpty()) // true

	// ═══════════════════════════════════════════════════════
	// 第 3 节：类型别名：同一类型的另一个名字
	// ═══════════════════════════════════════════════════════
	// type RawID = string 中的等号决定它是别名; 去掉等号就变成新定义类型
	// LegacyUserID 则是 UserID 这个新类型的别名
	// 对照 TS：type RawID = string 也是给原类型起别名; Go 去掉 = 则会创建新类型
	var alias RawID = "vvyh"
	var sameString string = alias                      // RawID 只是 string 类型的别名, 所以不需要转换直接赋值即可
	fmt.Printf("string 别名：%T %s\n", alias, sameString) // string vvyh

	var legacy LegacyUserID = UserID("894689")
	fmt.Println("本地类型别名保留方法：", legacy.IsEmpty(), legacy.APIValue())
	fmt.Println("原类型也具有该方法：", userID.APIValue()) // 注意: 对别名的操作也会影响到原类型

	// ═══════════════════════════════════════════════════════
	// 第 4 节：新类型的方法与转换
	// ═══════════════════════════════════════════════════════
	// StatusCode 有 IsSuccess 方法, InternalCode 是基于 StatusCode 定义的新类型, 不会自动得到它
	// 把 internal 显式转回 StatusCode, 才能调用原类型的方法

	code := StatusOK
	fmt.Println("状态码成功：", code.IsSuccess())            // true
	fmt.Println("404 成功：", StatusNotFound.IsSuccess()) // false

	internal := InternalCode(code)
	// internal.IsSuccess() 无法编译: 定义新类型不会继承 StatusCode 的方法
	fmt.Println("转回原定义类型后调用方法：", StatusCode(internal).IsSuccess())

	// 常见选择：
	// - 领域上不同的值（UserID/OrderID、金额/数量）用定义类型表达差异
	// - 迁移公开 API 时保留旧名字，可能用别名维护类型兼容性
	// - 只想缩短局部变量名，不需要为此建立新的类型

	// (可以简单看看)
	// ═══════════════════════════════════════════════════════
	// 第 5 节：iota 枚举与合法值校验
	// ═══════════════════════════════════════════════════════
	// iota 给一组状态起名，并不会限制变量只能取这些命名常量。
	fmt.Println("任务状态枚举：", JobUnknown, JobQueued, JobRunning, JobFinished) // 0 1 2 3
	fmt.Println("命名状态是否合法：", JobRunning.IsKnown())                         // true
	fmt.Println("未命名值是否合法：", JobState(99).IsKnown())                       // false
	// Go 的这类枚举不是封闭联合类型，不会自动要求 switch 穷尽所有状态。
	// 存入数据库或公开 API 的数值最好显式约定，随意插入 iota 声明会改变后续编号。

	// 因为 99 能放进 uint8，Go 就允许这个转换。Go 不会自动认为：
	// “你只列出了 0、1、2、3，所以 99 不允许使用。”
	// 如果业务上不接受 99，就要主动调用你写的检查方法：

	// ═══════════════════════════════════════════════════════
	// 第 6 节：转换不会代替范围检查
	// ═══════════════════════════════════════════════════════
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
