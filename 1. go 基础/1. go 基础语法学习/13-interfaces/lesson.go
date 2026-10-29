/*
章节:13-接口与类型断言
*/

// struct: 用来定义结构体（属性 / 数据）, 可以附带方法
// interface: 定义方法签名集, 不包含字段, 不实现方法具体逻辑

// Go 的 interface 是隐式实现的: 只要一个类型实现了接口中所有方法, 它就自动满足该接口, 不需要显式声明 implements
// Go 的 any 可容纳任意值, 但使用前需进行类型断言, 更接近 TS unknown, 而非 TS any

package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Getter interface {
	Get(key string) (string, bool)
}

type Setter interface {
	Set(key, value string)
}

// 嵌入两个接口, 将它们要求的方法集合起来
type Store interface {
	Getter
	Setter
}

type MemoryStore struct {
	values map[string]string
}

// 当你要满足的接口中存在需要指针修改数据的函数时, 建议全部函数的接收者都改为指针
func (store MemoryStore) Get(key string) (string, bool) {
	value, ok := store.values[key]
	return value, ok
}

func (store *MemoryStore) Set(key, value string) {
	// nil 接收者的方法可以被调用, 是否安全取决于方法自己的处理
	// 比如 var test *MemoryStore 声明一个 MemoryStore 类型的指针然后 test.Set()
	if store == nil {
		return
	}
	if store.values == nil {
		store.values = make(map[string]string)
	}
	store.values[key] = value
}

type FixedGetter struct {
	Value string
}

func (getter FixedGetter) Get(key string) (string, bool) {
	return getter.Value, true
}

// 判断类型是否满足接口(方法签名)的要求
var _ Store = &MemoryStore{}
var _ Getter = FixedGetter{}
var _ Getter = (*FixedGetter)(nil) // 类型转换基本写法: 类型(值)

// var _ Store = MemoryStore{} // 不满足要求, 因为 Set 方法的接收者是 MemoryStore 指针

// 二次封装?
// 第一层结构体负责数据增删改查
// 第二次结构体负责第一层结构体数据的业务逻辑
type PreferenceService struct {
	reader Getter
}

func (service PreferenceService) Theme() string {
	value, ok := service.reader.Get("theme")
	if !ok {
		return "system"
	}
	return value
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节: 隐式实现接口
	// ═══════════════════════════════════════════════════════
	showValue := func(getter Getter, key string) {
		value, ok := getter.Get(key)
		fmt.Println("读取：", key, value, ok)
	}

	store := &MemoryStore{} // MemoryStore 隐式满足 Store 接口, 根本不需要 implements 声明
	store.Set("theme", "dark")
	showValue(store, "theme") // MemoryStore 已经隐式满足 Store 接口, 那更满足 Store 接口的子集
	// 另一种类型也有 Get 方法, 所以也隐式满足 Get 接口
	showValue(FixedGetter{Value: "test-value"}, "theme")

	// ═══════════════════════════════════════════════════════
	// 第 2 节:类型断言与 comma-ok
	// ═══════════════════════════════════════════════════════
	// 接口值可理解为"动态类型 + 动态值", 使用断言可以检查里面装的实际类型
	var getter Getter = store // 能力被缩小了
	// 这时候 getter 使用不了 Set 方法
	concrete, ok := getter.(*MemoryStore) // 类型断言语法: 具体值, 是否成功 := 接口值.(目标类型)
	if ok {
		concrete.Set("locale", "zh-CN") // 这时候可以使用 Set 方法了, 因为通过类型断言修改了数据类型
	}
	fmt.Println("断言为 *MemoryStore: ", ok) // true
	_, ok = getter.(FixedGetter)
	fmt.Println("断言为 FixedGetter: ", ok) // false, 不会 panic
	// 单返回值 getter.(FixedGetter) 在失败时会 panic; 对未知输入优先用 comma-ok, 也就是双返回值, 使用 if 判断第二个返回值 ok

	// ═══════════════════════════════════════════════════════
	// 第 3 节:any 与类型 switch
	// ═══════════════════════════════════════════════════════
	// any 是空接口 interface{} 的别名, 两者完全等价 type any = interface{}
	// 所以不能直接调用任意字段/方法或参与任意运算
	describe := func(value any) string {
		// 这不是普通 switch, 而是类型 switch
		switch v := value.(type) {
		case string:
			return "文本：" + v
		case int:
			return fmt.Sprintf("整数：%d", v)
		case nil:
			return "未提供"
		default:
			return fmt.Sprintf("暂不支持：%T", v)
		}
	}

	for _, value := range []any{"hello", 42, true, nil} {
		fmt.Println("动态数据：", describe(value))
	}
	// TS 接口会在编译后擦除; Go 接口值在运行时携带动态类型, 支持运行时断言
	// 不要让所有函数都接收 any; 优先用具体类型或表达所需能力的小接口
	// 断言不是转换: any(int(3)).(float64) 不会自动把 3 转为浮点数

	// ═══════════════════════════════════════════════════════
	// 第 4 节: 接口的 nil: 类型和值都要为空
	// ═══════════════════════════════════════════════════════
	// 接口值只有动态类型和动态值都为空时才等于 nil
	var empty Getter // （无动态类型，无动态值）
	var missing *MemoryStore
	var wrapped Getter = missing                 // （*MemoryStore，nil 指针）
	fmt.Println("空接口值：", empty == nil)           // true
	fmt.Println("带 nil 指针的接口值：", wrapped == nil) // false
	fmt.Printf("接口仍携带类型：%T\n", wrapped)          // *main.MemoryStore
	// wrapped 的动态值为 nil, 但动态类型是 *MemoryStore
	// 本例的 Get 特别处理了 nil, 所以安全;不能推断其他实现也安全 (还未修改)
	showValue(wrapped, "theme")
	// 若某工厂函数的返回类型是 Getter, “无结果”时应直接 return nil, (还未修改)
	// 不要先声明 *MemoryStore 的 nil 再把它装进接口返回 (还未修改)

	// ═══════════════════════════════════════════════════════
	// 第 5 节: 按能力断言与依赖注入
	// ═══════════════════════════════════════════════════════
	// 也可以断言为另一个接口: 询问是否支持写能力, 而不是绑定某个具体实现
	if writable, ok := getter.(Setter); ok {
		writable.Set("theme", "light")
	}
	production := PreferenceService{reader: store}
	fixture := PreferenceService{reader: FixedGetter{Value: "dark"}}
	fmt.Println("相同业务逻辑、不同依赖：", production.Theme(), fixture.Theme()) // light dark

	// ═══════════════════════════════════════════════════════
	// 第 6 节:标准库 Reader / Writer (跳过)
	// ═══════════════════════════════════════════════════════
	// 读写接口让 io.Copy 同时适用于文件、网络和内存, 算法不必随数据来源重写
	// 本例 source 供数据, destination 接数据, 参数顺序是先目标、后来源

	demonstrateStandardInterfaces := func() {
		// 标准库也大量使用小接口:Reader 只有 Read([]byte) (int, error),
		// Writer 只有 Write([]byte) (int, error)字符串、文件和网络连接可提供同类能力
		var source io.Reader = strings.NewReader("shared interfaces")
		var buffer bytes.Buffer
		var destination io.Writer = &buffer
		// io.Copy 只依赖读写能力, 不需要知道这里其实是内存字符串和内存缓冲区
		written, err := io.Copy(destination, source)
		if err != nil {
			fmt.Println("复制失败：", err)
			return
		}
		fmt.Println("接口组合完成复制：", written, buffer.String()) // 17 shared interfaces
		// 若自己循环调用 Read, 要先处理 n>0 的数据, 再处理 err;同次调用可能两者都有
		// io.Copy 已处理这些通用细节;第 18、20 章再展开 error 和文件读取
	}

	demonstrateStandardInterfaces()

	// ═══════════════════════════════════════════════════════
	// 第 7 节:接口比较与动态类型
	// ═══════════════════════════════════════════════════════
	// any(1) 和 any(int64(1)) 的数值相同, 动态类型却不同, 所以不相等
	// 判断接口是否 nil 这种行为是安全的, 比较两个内部都装着相同的不可比较类型的接口则可能 panic

	demonstrateInterfaceEquality := func() {
		fmt.Println("接口比较动态类型和值：", any(1) == any(1), any(1) == any(int64(1))) // true false
		var values any = []int{1, 2}
		fmt.Println("包含切片的接口仍可判断 nil: ", values == nil) // false
		// map[any] 的键同样不能装 slice/map/function;需要键时优先设计明确的可比较类型
	}

	demonstrateInterfaceEquality()
}

// 反射不是使用 any 的默认方案, 优先用具体类型、接口或 type switch
