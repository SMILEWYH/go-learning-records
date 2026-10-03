/*
章节：19-泛型
	本章掌握清单：
	1. 类型参数、约束与类型推断
	2. 多个类型参数与泛型变换
	3. 泛型结构体、别名与定义类型
	4. 泛型容器与 comparable
	5. 方法约束与运行时接口
	6. 标准库泛型工具与浅拷贝
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./19-generics
// 运行练习：go run ./19-generics exam
// 阅读顺序：顶部集中放泛型函数、约束和类型声明；从 main 的具体调用开始，再回看 [T ...] 的含义。
// TS 对照：Go 的 [T any] 类似 TS 的 <T>；comparable 用于可比较类型，
// ~int64 表示允许底层类型为 int64 的自定义类型。约束是编译期规则。
// 后续再学：更复杂的约束设计与性能分析；优先复用 slices/maps，避免为泛型而泛型。

package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
)

// ~ 接受 type Cents int64 这样的定义类型；没有 ~ 时只接受列出的精确类型。
// 这种带类型集合的接口用于类型参数约束，不能当作普通变量的类型。
type Number interface {
	~int | ~int64 | ~float64
}

// | 是允许类型的并集，不是 TS 中可以随时换类型的变量联合。
// 编译器只允许约束中所有类型都支持的操作；T any 上不能直接使用 + 或比较大小。

func Sum[T Number](values []T) T {
	var total T // 类型参数也有零值。
	for _, value := range values {
		total += value // 约束中每种类型都支持 +，所以这里可以相加。
	}
	return total
}

// 泛型适合“算法相同，数据类型不同”的代码；不必为每个业务方法加泛型。
func MapSlice[T any, R any](values []T, transform func(T) R) []R {
	result := make([]R, 0, len(values))
	for _, value := range values {
		result = append(result, transform(value))
	}
	return result
}

// 泛型结构体可以表达同一 API 响应形状中的不同实体。
type Page[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

// 泛型别名自 Go 1.24 起可用：同一组类型实参下，ResponsePage 与 Page 是同一类型。
type ResponsePage[T any] = Page[T]

// 不带 = 则创建新的定义类型，不继承 Page 的方法，需要显式转换或自行声明方法。
type DefinedPage[T any] Page[T]

func (p Page[T]) IsEmpty() bool {
	return len(p.Items) == 0
}

// 泛型切片和 map 是定义类型；类型参数在实例化后会固定。
type List[T any] []T
type Lookup[K comparable, V any] map[K]V

// map 的 key 必须支持 == / !=；slice、map、func 不能用作 key。
// TS Record<string, T> 是对象类型；Go map[K]V 可以使用整型等可比较 key。
type Cents int64

type User struct {
	ID   string
	Name string
}

type Status string

func (s Status) String() string { return "状态=" + string(s) }

// 方法约束允许泛型代码调用契约中的方法；是否满足约束也取决于指针/值方法集。
// fmt.Stringer 只有方法，因此既可作泛型约束，也可作保存运行时值的普通接口。
func Labels[T fmt.Stringer](items []T) []string {
	return MapSlice(items, func(item T) string { return item.String() })
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：类型参数、约束与类型推断
	// ═══════════════════════════════════════════════════════
	// 读 Sum[T Number](values []T) T：T 是待确定类型，Number 限制可用类型，输入输出共用 T。
	// Sum([]int{...}) 由参数推断 T；Sum[int](nil) 则明确指定 int。

	fmt.Println("整数和：", Sum([]int{1, 2, 3})) // 从实参推断 T 为 int。
	fmt.Println("空切片和：", Sum[int](nil))      // nil 无法推断元素类型，需要显式写 T。
	fmt.Println("订单总分：", Sum([]Cents{1990, 2500}))
	// 金额用最小货币单位的整数表示；这里不涉及浮点金额舍入。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：多个类型参数与泛型变换
	// ═══════════════════════════════════════════════════════
	// MapSlice 的 T 是输入元素类型，R 是结果元素类型。
	// 这里回调把 User 变为 string，所以编译器推断结果是 []string。

	users := []User{{ID: "u-1", Name: "小林"}, {ID: "u-2", Name: "小陈"}}
	names := MapSlice(users, func(user User) string { return user.Name })
	fmt.Println("用户名：", names)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：泛型结构体、别名与定义类型
	// ═══════════════════════════════════════════════════════
	// Page[User] 是把类型参数具体化后的类型，Items 因此是 []User。
	// ResponsePage 有 =，与 Page 身份相同；DefinedPage 没有 =，需要转换才能调用原方法。

	page := Page[User]{Items: users, Total: len(users)}
	fmt.Println("用户页为空：", page.IsEmpty())
	var aliasPage ResponsePage[User] = page // 别名具有相同类型身份，可直接赋值并调用同一方法。
	fmt.Println("别名页为空：", aliasPage.IsEmpty())
	definedPage := DefinedPage[User](page)
	fmt.Println("定义类型显式转回后调用方法：", Page[User](definedPage).IsEmpty())

	// ═══════════════════════════════════════════════════════
	// 第 4 节：泛型容器与 comparable
	// ═══════════════════════════════════════════════════════
	// Lookup[K comparable, V any] 中只要求 K 能比较，V 可为任意固定类型。
	// 这里 K=string、V=User，类型参数一经实例化就不会在运行中变化。

	tags := List[string]{"go", "api"}
	byID := Lookup[string, User]{"u-1": users[0]}
	fmt.Println("标签：", tags, "索引：", byID["u-1"].Name)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：方法约束与运行时接口
	// ═══════════════════════════════════════════════════════
	// Labels 的约束要求 String() string，所以函数体可安全调用这个方法。
	// fmt.Stringer 也能作为普通接口保存运行时值；约束和接口值的使用场景不同。

	fmt.Println("方法约束：", Labels([]Status{"ready", "done"}))
	var display fmt.Stringer = Status("ready")
	fmt.Println("同一方法契约作为运行时接口：", display.String())

	// ═══════════════════════════════════════════════════════
	// 第 6 节：标准库泛型工具与浅拷贝
	// ═══════════════════════════════════════════════════════
	// slices.Sort 修改原切片，因此先 Clone 再排序；maps.Clone 复制一层映射。
	// 泛型不会改变复制规则，嵌套引用数据仍需明确是否共享。

	// 标准库泛型工具优先：Clone 为浅拷贝，Sort 会修改传入切片。
	priorities := []int{3, 1, 2}
	sorted := slices.Clone(priorities)
	slices.Sort(sorted)
	fmt.Println("原顺序与排序后：", priorities, sorted, "包含 2：", slices.Contains(sorted, 2))
	copiedIndex := maps.Clone(byID)
	copiedIndex["u-1"] = User{ID: "u-1", Name: "副本中的名字"}
	fmt.Println("修改浅拷贝索引后原用户名：", byID["u-1"].Name)
	// 若元素内部含 slice/map/指针，浅拷贝仍共享它们所指向的数据。
	// any 是空接口别名；T any 表示任意一种固定类型，不是每个元素可随便变类型。
	// 可以显式令 T=any，此时元素保存为接口值，就可以持有不同的动态类型。
	// []any 与 []string 是不同类型，不能直接相互赋值。
	// Go 方法不能自行新增类型参数；可使用接收者类型已有的 T，或写泛型函数。
	// 推断依赖函数的参数等类型信息；无参数函数通常不能仅靠接收返回值的变量推断 T。
	// 泛型类型的实例化必须写类型实参，如 Page[User]，不能写 Page{...} 让字段推断。
	// comparable 只允许 ==/!=，不代表能排序。接口类型可满足 comparable，
	// 但 map[any] 的动态 key 仍必须可比较：放入 []int 会在运行时 panic。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 给 Sum 传 []string，阅读 Number 约束导致的编译错误后恢复。
// 把 Number 中的 ~int64 改成 int64，解释 Cents 为什么不再满足约束。
