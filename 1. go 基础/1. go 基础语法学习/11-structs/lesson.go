/*
章节：11-结构体与方法
*/

// TS 对照：struct 是明确的数据布局；方法接收者类似显式 this，但没有 class 继承。

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 大写开头表示可被其他包访问; 并不是所有结构体或字段都必须大写
type User struct {
	ID       int    `json:"id"` // 在 JSON 序列化和反序列化时对属性名进行改名
	Name     string `json:"name"`
	Age      int    `json:"age,omitempty"` // omitempty 字段表示如果字段的值是空值, 就省略这个字段
	Password string `json:"-"`             // 在 JSON 序列化和反序列化时都忽略这个字段, json:"-" 只约束该编码器，fmt.Printf("%+v", user) 仍可能打印 Password
	internal string // 小写开头, 未导出, 其他包不能直接访问, 所以在 JSON 序列化和反序列化时不会被编码, 即使加了 json 标签也一样
}

// 结构体的方法必须在包级声明, 不能像第 08 章的匿名函数那样写进 main
// (u User) 表示当前函数方式属于 User 结构体, u 则是结构体本身(this)
func (u User) Rename_1(name string) { // 仅修改副本
	u.Name = name
}

func (u *User) Rename_2(name string) { // 这样才可以修改到结构体的值
	u.Name = name
}

func (u User) Label() {
	fmt.Println(u.Name)
}

// 结构体没有类似 class 继承, 而是直接在内部声明另一种结构体
type Organization struct {
	Name string
}

type Member struct {
	Organization // 嵌入字段：字段名是 Organization, Member实例.Organization.Name
	Name         string
}

func (m Member) consoleName() {
	fmt.Println(m.Name)
	fmt.Println(m.Organization.Name)
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节：声明结构体、零值与字段初始化
	// ═══════════════════════════════════════════════════════
	// 零值结构体的每个字段也都是零值，不需要构造函数才能使用。
	var zero User
	fmt.Printf("零值：%+v", zero)
	user := User{ID: 7, Name: "小林", Password: "demo-secret", internal: "cache"}

	// ═══════════════════════════════════════════════════════
	// 第 2 节：结构体复制与两种方法接收者
	// ═══════════════════════════════════════════════════════
	// Go: 复制结构体是值复制, 修改复制结构体的字符串不会影响到原来的结构体
	// Ts: 对象赋值通常仍指向同一个对象, 修改复制对象的字符串会影响到原来的对象
	copyOfUser := user
	copyOfUser.Name = "副本"
	fmt.Println("原值 / 副本：", user.Name, copyOfUser.Name) // 小林 副本

	// 现在再试一试两种不同结构体方法修改内部值
	user.Rename_1("无效改名")
	fmt.Println("值接收者执行后：", user.Name) // 小林
	user.Rename_2("小周")
	fmt.Println("指针接收者执行后：", user.Name) // 小周
	// Rename_2 需要一个指向 User 的指针, 那么应该这么调用 (&user).Rename_2("小周") 但有语法糖所以直接写 user.Rename_2("小周") 也是可以的

	// ═══════════════════════════════════════════════════════
	// 第 3 节：嵌入：组合、提升与同名字段
	// ═══════════════════════════════════════════════════════
	// 构造复合字面量时必须写 Organization: { ... }, 不可以直接 { ... }
	member := Member{Organization: Organization{Name: "平台组"}, Name: "小周"}
	fmt.Println("外层 Name: ", member.Name)                // 小周
	fmt.Println("嵌入字段 Name: ", member.Organization.Name) // 平台组
	member.consoleName()
	// 如果 Member 第一层没有 Name 属性, 而其内部 Organization 有, 那么 member.Name 会指向 Organization.Name, 类似 Ts 的作用域链

	// ═══════════════════════════════════════════════════════
	// 第 4 节：JSON 编码与字段标签
	// ═══════════════════════════════════════════════════════
	// Age 为 0 被省略, Password 被永久忽略, internal 不导出
	encoded, err := json.Marshal(user)
	if err != nil {
		fmt.Println("JSON 编码失败：", err)
		return
	}
	fmt.Println("JSON 序列化: ", string(encoded)) // {"id":7,"name":"小周"}
	// 若需要区分“未提供 age”和“提供 age=0”，可以使用属性类型为 *int + omitempty 字段

	// ═══════════════════════════════════════════════════════
	// 第 5 节：构造函数与业务校验
	// ═══════════════════════════════════════════════════════
	// struct 没有类似 class 的构造函数, 那么我们可以自己写一个工厂函数代替构造函数
	NewUser := func(id int, name string) (User, bool) {
		name = strings.TrimSpace(name)
		if id <= 0 || name == "" {
			return User{}, false
		}
		return User{ID: id, Name: name}, true
	}

	_, ok := NewUser(9, "vvyh")
	fmt.Println("构造并校验：", ok) // true
	_, ok = NewUser(0, "")
	fmt.Println("拒绝无效数据：", ok) // false

	// ═══════════════════════════════════════════════════════
	// 第 6 节：结构体能否比较
	// ═══════════════════════════════════════════════════════
	// 所有字段类型都可比较时, 结构体才能使用 == / !=，也可作为 map 键
	user1 := User{ID: 1}
	user2 := User{ID: 1}
	fmt.Println("结构体逐字段比较：", user1 == user2) // true
	// 带 []string、map 或函数字段的结构体不能直接 ==

	// ═══════════════════════════════════════════════════════
	// 第 7 节：JSON 解码：缺失、零值与 null
	// ═══════════════════════════════════════════════════════
	// 下方示例再用 *int 区分“没有数值”和“数值为 0”。

	demonstrateJSONInput := func() {
		user := User{Name: "旧名字", Age: 30}
		err := json.Unmarshal([]byte(`{"name":"新名字","extra":true}`), &user) // 第一个参数为 JSON 字符串, 第二个参数为结构体指针
		if err != nil {
			fmt.Println("JSON 解码失败：", err)
			return
		}
		fmt.Println("缺失字段保留已有值：", user.Name, user.Age) // 新名字 30
		// 忽略 extra 字段; 结构体 age 旧值保留, 每次请求通常用新的目标结构体，避免无意复用旧数据

		// 小巧思
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

		// 输入 {} -> 指针为 nil=true -> 输出 {}
		// 输入 {"limit":0} -> 指针为 nil=false -> 输出 {"limit":0}
		// 输入 {"limit":null} -> 指针为 nil=true -> 输出 {}
	}

	demonstrateJSONInput()

	// ═══════════════════════════════════════════════════════
	// 第 8 节：方法值与方法表达式（选读）
	// ═══════════════════════════════════════════════════════
	demonstrateMethodValues := func() {
		vvyh := User{
			Name: "vvyh",
		}

		labelSnapshot := vvyh.Label // 保存当时的 User 副本
		rename := user.Rename_2     // 保存 &user，后续修改同一实例
		rename("wyh")

		labelSnapshot() // "vvyh"
		user.Label()    // "wyh"

		labelFunction := User.Label
		labelFunction(user) // "wyh"
	}

	demonstrateMethodValues()
}
