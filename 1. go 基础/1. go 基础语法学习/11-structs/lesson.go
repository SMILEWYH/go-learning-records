/*
章节：11-结构体与方法
	本章掌握清单：
	1. 声明结构体、零值与字段初始化
	2. 结构体复制与两种方法接收者
	3. 嵌入：组合、提升与同名字段
	4. JSON 编码与字段标签
	5. 构造函数与业务校验
	6. 结构体能否比较
	7. JSON 解码：缺失、零值与 null
	8. 方法值与方法表达式（选读）
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./11-structs
// 运行练习：go run ./11-structs exam
// 阅读顺序：上方 User、Organization、Member 是示例的数据模型和方法，先看 main，再回看定义。
// TS 对照：struct 是明确的数据布局；方法接收者类似显式 this，但没有 class 继承。
// 大写开头表示可被其他包访问；并不是所有结构体或字段都必须大写。

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSON 只编码导出的字段。标签中的冒号后没有空格，标签使用反引号包裹。
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age,omitempty"`
	Password string `json:"-"`
	internal string // 未导出的字段不会被 encoding/json 编码。
}

// 值接收者 u 是副本，适合不需要修改接收者的操作。
func (u User) Label() string {
	return fmt.Sprintf("%d:%s", u.ID, u.Name)
}

func (u User) RenameCopy(name string) {
	u.Name = name // 仅修改副本。
}

// 指针接收者修改原值；u.Name 等价于 (*u).Name。
func (u *User) Rename(name string) {
	if u == nil {
		return
	}
	u.Name = name
}

// 此处同时展示两种接收者以便比较。实际工程中同一类型的方法一般保持一致；
// 尤其当它持有锁或表示可变实体时，更应避免无意复制。
// 即使是值接收者，内部 slice、map、指针字段仍可能共享底层数据。

type Organization struct {
	Name string
}

func (o Organization) OrganizationLabel() string {
	return "组织：" + o.Name
}

type Member struct {
	Organization // 嵌入字段：字段名是 Organization。
	Name         string
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：声明结构体、零值与字段初始化
	// ═══════════════════════════════════════════════════════
	// User 把一个用户的相关数据放在同一个值中，字段可有不同类型。
	// User{ID: 7, Name: "小林"} 按字段名初始化，未写的字段保留零值。

	// 零值结构体的每个字段也都是零值，不需要构造函数才能使用。
	var zero User
	fmt.Printf("零值：%+v\n", zero)
	user := User{ID: 7, Name: "小林", Password: "demo-secret", internal: "cache"}
	fmt.Println("用户标签：", user.Label())

	// ═══════════════════════════════════════════════════════
	// 第 2 节：结构体复制与两种方法接收者
	// ═══════════════════════════════════════════════════════
	// func (u User) 声明值接收者，func (u *User) 声明指针接收者。
	// RenameCopy 修改副本，Rename 修改原对象；对照输出“小林”和“小周”。

	// 复制结构体，再修改字符串字段，不影响原来的结构体。
	copyOfUser := user
	copyOfUser.Name = "副本"
	fmt.Println("原值 / 副本：", user.Name, copyOfUser.Name) // 小林 副本
	user.RenameCopy("无效改名")
	fmt.Println("值接收者执行后：", user.Name) // 小林
	user.Rename("小周")
	fmt.Println("指针接收者执行后：", user.Name) // 小周
	// user 可取地址，编译器允许 user.Rename(...) 简写 (&user).Rename(...)。
	// 这项调用便利不会改变方法集规则：第 13 章接口里会区分 User 和 *User。
	// map 下标表达式通常不可取地址；修改 map[string]User 的值需取出、修改、写回。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：嵌入：组合、提升与同名字段
	// ═══════════════════════════════════════════════════════
	// Member 中未另起名字的 Organization 是嵌入字段，能简写访问其方法。
	// 完整路径 member.Organization.Name 仍然存在，也能解决同名字段遮蔽。

	// 嵌入是组合。Member 拥有 Organization 字段，不是继承其类型身份。
	member := Member{Organization: Organization{Name: "平台组"}, Name: "小周"}
	fmt.Println("外层 Name：", member.Name)                // 小周
	fmt.Println("嵌入字段 Name：", member.Organization.Name) // 平台组
	fmt.Println("提升的方法：", member.OrganizationLabel())
	// 外层 Name 遮蔽了嵌入层的同名字段，仍可写完整路径访问内部字段。
	// 若同一深度两个嵌入类型都有 Name，而外层没有 Name，member.Name 会有歧义。
	// 构造复合字面量时必须写 Organization: ...，不能直接设置提升的内部字段。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：JSON 编码与字段标签
	// ═══════════════════════════════════════════════════════
	// 标签是附在字段上的说明，encoding/json 按 json 标签决定输出名称和省略规则。
	// 先预测哪些字段会出现，再看编码后的字符串。

	// 转成 API JSON：Age 为 0 被省略，Password 被永久忽略，internal 不导出。
	encoded, err := json.Marshal(user)
	// Marshal 返回数据和错误。本章先采用最简单的非 nil 检查，错误章再系统介绍。
	if err != nil {
		fmt.Println("JSON 编码失败：", err)
		return
	}
	fmt.Println("API 响应：", string(encoded)) // {"id":7,"name":"小周"}
	// omitempty 对此处 int 使用 0 作为空值，对字符串使用 ""；并非通用校验器。
	// 若 API 需要区分“未提供 age”和“提供 age=0”，可以使用 *int + omitempty。
	// json:"-" 只约束该编码器，fmt.Printf("%+v", user) 仍可能打印 Password。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：构造函数与业务校验
	// ═══════════════════════════════════════════════════════
	// NewUser 是我们自己写的普通函数，负责去空格、检查 ID 和姓名。
	// 创建结构体只保证字段类型正确，业务是否合法还需要显式判断。

	// NewXxx 是普通函数的命名约定，不是语言内建构造语法，也不会被自动调用。
	created, ok := NewUser(9, "  阿青  ")
	fmt.Println("构造并校验：", created.Label(), ok) // 9:阿青 true
	_, ok = NewUser(0, "")
	fmt.Println("拒绝无效数据：", ok) // false
	// 仍然可以直接写 User{}，因此需要明确哪些 API 要求先校验，不能假设所有实例合法。
	// 工程中可用未导出字段限制外部直接写字段，但零值能否使用仍需自己设计。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：结构体能否比较
	// ═══════════════════════════════════════════════════════
	// 比较能力由全部字段类型共同决定。
	// User 的字段都可比较，所以 == 能逐字段比较；加入切片字段后就不能这样写。

	// 所有字段类型都可比较时，结构体才能使用 == / !=，也可作为 map 键。
	fmt.Println("结构体逐字段比较：", User{ID: 1} == User{ID: 1}) // true
	// 带 []string、map 或函数字段的结构体不能直接 ==，要按业务定义如何比较。
	// 带接口字段时还需要注意接口中的实际值是否可比较，第 13 章会展开。

	// ═══════════════════════════════════════════════════════
	// 第 7 节：JSON 解码：缺失、零值与 null
	// ═══════════════════════════════════════════════════════
	// 传 &user 才能填充变量；复用旧目标时，没出现的字段仍保留旧值。
	// 下方示例再用 *int 区分“没有数值”和“数值为 0”。

	demonstrateJSONInput()

	// ═══════════════════════════════════════════════════════
	// 第 8 节：方法值与方法表达式（选读）
	// ═══════════════════════════════════════════════════════
	// user.Label 是已经选定接收者的方法值，User.Label 则需要把接收者作为参数传入。
	// 先掌握直接调用，再看把方法保存为回调时复制了什么。

	demonstrateMethodValues(user)
}

func NewUser(id int, name string) (User, bool) {
	name = strings.TrimSpace(name)
	if id <= 0 || name == "" {
		return User{}, false
	}
	return User{ID: id, Name: name}, true
	// 第 18 章学习 error 后，可返回更具体的校验失败原因。
}

func demonstrateJSONInput() {
	// Unmarshal 必须接收可写目标的指针；传入结构体值无法填充调用方变量。
	user := User{Name: "旧名字", Age: 30}
	if err := json.Unmarshal([]byte(`{"name":"新名字","extra":true}`), &user); err != nil {
		fmt.Println("JSON 解码失败：", err)
		return
	}
	fmt.Println("缺失字段保留已有值：", user.Name, user.Age) // 新名字 30
	// 缺失字段不会自动清零。每次请求通常用新的目标结构体，避免无意复用旧数据。
	// 默认忽略未知 extra 字段；需要严格拒绝时再学习 Decoder.DisallowUnknownFields。
	// 成功解码只说明类型能匹配，不代表 ID、年龄等业务值已经通过校验。

	for _, raw := range []string{`{}`, `{"limit":0}`, `{"limit":null}`} {
		// 匿名结构体适合局部、一次性的数据形状，不需要为每个简单响应声明公共类型。
		var input struct {
			Limit *int `json:"limit,omitempty"`
		}
		if err := json.Unmarshal([]byte(raw), &input); err != nil {
			fmt.Println("分页输入失败：", err)
			continue
		}
		encoded, err := json.Marshal(input)
		if err != nil {
			fmt.Println("分页输出失败：", err)
			continue
		}
		fmt.Printf("输入 %s -> 指针为 nil=%v -> 输出 %s\n", raw, input.Limit == nil, encoded)
	}
	// 新目标上：缺失与 null 都得到 nil；0 得到非 nil 指针，omitempty 仍保留 limit:0。
	// omitempty 只影响编码，不改变解码规则；移除该标签后，nil 指针编码为 JSON null。
	// 单个 *int 无法区分“缺失”与“显式 null”；PATCH 需要三态时需另设计存在性信息。
	// 将 JSON null 解到普通 int 不会把它自动设为 0，而是保持原值；类型选择需符合 API 约定。
}

func demonstrateMethodValues(user User) {
	// 选读：把方法当回调时，方法值会在创建时求值并保存接收者。
	labelSnapshot := user.Label // 值接收者：保存当时的 User 副本。
	rename := user.Rename       // 指针接收者：保存 &user，后续修改同一实例。
	rename("阿禾")
	fmt.Println("方法值捕获的副本 / 当前值：", labelSnapshot(), user.Label())
	// 方法表达式把接收者变成第一个显式参数，类型是 func(User) string。
	labelFunction := User.Label
	fmt.Println("方法表达式：", labelFunction(user))
}

// 后续再学：接口方法集（13 章）、错误链（18 章）、严格 JSON 校验和自定义序列化；
// 方法值/方法表达式先能读懂即可，不必为简单业务刻意使用。

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 Age 改为 18，再改回 0，观察 omitempty 的输出变化。
// 用相同名字分别调用 RenameCopy 与 Rename，说明各自操作哪个 User。
