/*
章节：13-接口与类型断言
	本章掌握清单：
	1. 小接口、隐式实现与方法集
	2. 类型断言与 comma-ok
	3. any 与类型 switch
	4. 接口的 nil：类型和值都要为空
	5. 按能力断言与依赖注入
	6. 标准库 Reader / Writer
	7. 接口比较与动态类型
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./13-interfaces
// 运行练习：go run ./13-interfaces exam
// 阅读顺序：上方声明接口和两个实现；先看 main 中如何替换实现，再回看 Get/Set 方法。
// TS 对照：Go 接口在本章描述方法契约，类型无需写 implements。
// Go 的 any 可容纳任意值，但使用前需判断类型，更接近 TS unknown，而非 TS any。

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// 消费者只依赖读取能力，不必知道数据存在缓存、数据库还是测试替身中。
type Getter interface {
	Get(key string) (string, bool)
}

type Setter interface {
	Set(key, value string)
}

// 嵌入两个接口，将它们要求的方法集合起来。
type Store interface {
	Getter
	Setter
}

type MemoryStore struct {
	values map[string]string
}

func (store *MemoryStore) Get(key string) (string, bool) {
	// nil 接收者的方法可以被调用，是否安全取决于方法自己的处理。
	if store == nil {
		return "", false
	}
	value, ok := store.values[key]
	return value, ok
}

func (store *MemoryStore) Set(key, value string) {
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

// 编译期确认类型实现接口；只有方法签名完全匹配才满足要求。
var _ Store = (*MemoryStore)(nil)
var _ Getter = FixedGetter{}
var _ Getter = (*FixedGetter)(nil)

func showValue(getter Getter, key string) {
	value, ok := getter.Get(key)
	fmt.Println("读取：", key, value, ok)
}

// 在使用处定义“小接口”，把依赖传进来，就是最基础的依赖注入，无需专门框架。
// 这个服务只需要读取，没理由要求依赖同时提供写入、删除、连接数据库等能力。
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
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：小接口、隐式实现与方法集
	// ═══════════════════════════════════════════════════════
	// Getter 要求 Get 方法，Setter 要求 Set 方法，Store 把两种能力合在一起。
	// 只看方法签名是否满足，不需要写 implements，也不依赖结构体字段是否相同。

	// *MemoryStore 隐式满足 Store，根本不需要 implements 声明。
	store := &MemoryStore{}
	store.Set("theme", "dark")
	showValue(store, "theme")
	showValue(FixedGetter{Value: "test-value"}, "theme")
	// T 的方法集包含值接收者方法；*T 包含值接收者和指针接收者方法。
	// MemoryStore{} 不满足 Getter；虽然可取地址变量调用方法时允许自动取地址，
	// 赋给接口时不会自动把 MemoryStore 值变成 *MemoryStore。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：类型断言与 comma-ok
	// ═══════════════════════════════════════════════════════
	// getter.(*MemoryStore) 询问接口当前保存的值是不是这个具体类型。
	// 两个结果中的 ok 为 false 时，不要继续把 concrete 当成有效实例使用。

	// 接口可理解为“动态类型 + 动态值”。断言检查里面装的实际类型。
	var getter Getter = store
	concrete, ok := getter.(*MemoryStore)
	if ok {
		concrete.Set("locale", "zh-CN")
	}
	fmt.Println("断言为 *MemoryStore：", ok) // true
	_, ok = getter.(FixedGetter)
	fmt.Println("断言为 FixedGetter：", ok) // false，不会 panic。
	// 单返回值 getter.(FixedGetter) 在失败时会 panic；对未知输入优先用 comma-ok。
	// 断言不是转换：any(int(3)).(float64) 不会自动把 3 转为浮点数。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：any 与类型 switch
	// ═══════════════════════════════════════════════════════
	// any 可以存放不同类型，但每次使用仍需明确当前是哪一种。
	// describe 中的 type switch 让各分支可以按具体类型处理值。

	// any 是 interface{} 的别名，不能直接调用任意字段/方法或参与任意运算。
	for _, value := range []any{"hello", 42, true, nil} {
		fmt.Println("动态数据：", describe(value))
	}
	// TS 接口会在编译后擦除；Go 接口值在运行时携带动态类型，支持运行时断言。
	// 不要让所有函数都接收 any；优先用具体类型或表达所需能力的小接口。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：接口的 nil：类型和值都要为空
	// ═══════════════════════════════════════════════════════
	// empty 没有动态类型和值；wrapped 保留 *MemoryStore 类型，即使里面的指针为 nil。
	// 这也是“返回 nil 指针却得到非 nil error”的根本原因，第 18 章会再演示。

	// 接口只有“动态类型和动态值都为空”时才等于 nil。
	var empty Getter
	var missing *MemoryStore
	var wrapped Getter = missing
	fmt.Println("空接口值：", empty == nil)           // true
	fmt.Println("带 nil 指针的接口值：", wrapped == nil) // false
	fmt.Printf("接口仍携带类型：%T\n", wrapped)          // *main.MemoryStore
	// wrapped 的动态值为 nil，但动态类型是 *MemoryStore。
	// 本例的 Get 特别处理了 nil，所以安全；不能推断其他实现也安全。
	showValue(wrapped, "theme")
	// 若某工厂函数的返回类型是 Getter，“无结果”时应直接 return nil，
	// 不要先声明 *MemoryStore 的 nil 再把它装进接口返回。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：按能力断言与依赖注入
	// ═══════════════════════════════════════════════════════
	// PreferenceService 只知道 Getter，所以既可接内存存储，也可接固定值的测试实现。
	// 断言 Setter 是额外询问写能力，不需要把所有消费者都绑到 MemoryStore。

	// 也可以断言为另一个接口：询问是否支持写能力，而不是绑定某个具体实现。
	if writable, ok := getter.(Setter); ok {
		writable.Set("theme", "light")
	}
	production := PreferenceService{reader: store}
	fixture := PreferenceService{reader: FixedGetter{Value: "dark"}}
	fmt.Println("相同业务逻辑、不同依赖：", production.Theme(), fixture.Theme()) // light dark
	// 测试时传入简单实现即可控制输入；把任意实现都放进 any 会丢掉有用的能力约束。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：标准库 Reader / Writer
	// ═══════════════════════════════════════════════════════
	// 读写接口让 io.Copy 同时适用于文件、网络和内存，算法不必随数据来源重写。
	// 本例 source 供数据，destination 接数据，参数顺序是先目标、后来源。

	demonstrateStandardInterfaces()

	// ═══════════════════════════════════════════════════════
	// 第 7 节：接口比较与动态类型
	// ═══════════════════════════════════════════════════════
	// any(1) 和 any(int64(1)) 的数值相同，动态类型却不同，所以不相等。
	// 判断接口是否 nil 是安全的，比较两个内部都装切片的接口则可能 panic。

	demonstrateInterfaceEquality()
}

func demonstrateStandardInterfaces() {
	// 标准库也大量使用小接口：Reader 只有 Read([]byte) (int, error)，
	// Writer 只有 Write([]byte) (int, error)。字符串、文件和网络连接可提供同类能力。
	var source io.Reader = strings.NewReader("shared interfaces")
	var buffer bytes.Buffer
	var destination io.Writer = &buffer
	// io.Copy 只依赖读写能力，不需要知道这里其实是内存字符串和内存缓冲区。
	written, err := io.Copy(destination, source)
	if err != nil {
		fmt.Println("复制失败：", err)
		return
	}
	fmt.Println("接口组合完成复制：", written, buffer.String()) // 17 shared interfaces
	// 若自己循环调用 Read，要先处理 n>0 的数据，再处理 err；同次调用可能两者都有。
	// io.Copy 已处理这些通用细节；第 18、20 章再展开 error 和文件读取。
}

func demonstrateInterfaceEquality() {
	fmt.Println("接口比较动态类型和值：", any(1) == any(1), any(1) == any(int64(1))) // true false
	var values any = []int{1, 2}
	fmt.Println("包含切片的接口仍可判断 nil：", values == nil) // false
	// 接口语法上可以 ==，但若两边装着相同的不可比较类型，运行时会 panic。
	// 例如 values == values 会尝试比较 []int，而 slice 不支持这种比较，所以不能执行。
	// map[any] 的键同样不能装 slice/map/function；需要键时优先设计明确的可比较类型。
}

func describe(value any) string {
	// switch 内的 v 在各 case 中具有对应的具体类型。
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

// 后续再学：error 也是接口（18 章）、类型集合仅用于泛型约束（19 章）；
// 反射不是使用 any 的默认方案，优先用具体类型、接口或 type switch。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 将 PreferenceService 的 reader 换成 FixedGetter，观察业务方法为何不用改。
// 分别打印 missing == nil 与 wrapped == nil，并用动态类型解释差别。
