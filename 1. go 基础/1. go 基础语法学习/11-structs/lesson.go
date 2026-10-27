/*
章节:11-结构体与方法
*/

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
	Password string `json:"-"`             // 在 JSON 序列化和反序列化时都忽略这个字段, json:"-" 只约束该编码器, fmt.Printf("%+v", user) 仍可能打印 Password
	internal string // 小写开头, 未导出, 其他包不能直接访问, 所以使用 encoding/json 这个包进行序列化和反序列化时不会被编码, 即使加了 json 标签也一样
}

// (u User) 是方法的接收者, 表示 Rename_1 是属于 User 类型的方法
// - User: 接收者的类型
// - u：方法内部使用的变量名
// 包级结构体的方法必须在包级声明, 不能像第 08 章的匿名函数那样写进 main
func (u User) Rename_1(name string) { // 仅修改副本
	u.Name = name
}

func (u *User) Rename_2(name string) { // 这样才可以修改到结构体的值
	u.Name = name
}

// 用于第九节示例
func (u User) Label() {
	fmt.Println(u.Name)
}

// 结构体没有类似 class 继承, 而是直接在内部声明另一种结构体
type Organization struct {
	Name string
}

type Member struct {
	Organization // 默认字段名是 Organization, 使用方式: Member实例.Organization.Name
	Name         string
}

func (m Member) consoleName() {
	fmt.Println(m.Name)
	fmt.Println(m.Organization.Name)
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节:声明结构体、零值与字段初始化
	// ═══════════════════════════════════════════════════════
	var zero User              // 声明结构体但未初始化, 其字段都是零值
	fmt.Printf("零值：%+v", zero) // 零值：{ID:0 Name: Age:0 Password: internal:}

	// 声明、初始化结构体
	user := User{ID: 7, Name: "小林", Password: "demo-secret", internal: "cache"}

	// ═══════════════════════════════════════════════════════
	// 第 2 节:结构体复制与两种方法接收者
	// ═══════════════════════════════════════════════════════
	// 复制结构体是值复制, 修改复制结构体内字符串不会影响到原来的结构体, 但如果是切片、map 或指针类型, 复制后仍然指向同一个底层数据, 则修改会影响指针指向的数据
	copyOfUser := user
	copyOfUser.Name = "副本"
	fmt.Println("原值 / 副本：", user.Name, copyOfUser.Name) // 小林 副本

	// 现在再试一试两种不同结构体方法修改内部值
	user.Rename_1("无效改名")
	fmt.Println("值接收者执行后：", user.Name) // 小林
	user.Rename_2("小周")
	fmt.Println("指针接收者执行后：", user.Name) // 小周
	// Rename_2 需要一个指向 User 的指针, 那么应该这么调用 (&user).Rename_2("小周") 但有语法糖所以直接写 user.Rename_2("小周") 也是允许的

	// ═══════════════════════════════════════════════════════
	// 第 3 节:嵌入:组合、提升与同名字段
	// ═══════════════════════════════════════════════════════
	// 构造复合字面量时必须写 Organization: { ... }, 不可以直接 { ... }
	member := Member{
		Organization: Organization{
			Name: "平台组",
		},
		Name: "小周",
	}
	fmt.Println("外层 Name: ", member.Name)                // 小周
	fmt.Println("嵌入字段 Name: ", member.Organization.Name) // 平台组
	member.consoleName()
	// 如果 Member 第一层没有 Name 属性, 而其内部 Organization 有, 那么 member.Name 会指向 Organization.Name, 类似 Ts 的作用域链

	// ═══════════════════════════════════════════════════════
	// 第 4 节:构造函数与业务校验
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
	// 第 5 节:结构体能否比较
	// ═══════════════════════════════════════════════════════
	// 所有字段类型都可比较时, 结构体才能使用 == 、 != 来比较结构体值, 也可作为 map 键
	user1 := User{ID: 1}
	user2 := User{ID: 1}
	fmt.Println("结构体逐字段比较：", user1 == user2) // true
	// 带 []string、map 或函数字段的结构体不能直接 ==

	// ═══════════════════════════════════════════════════════
	// 第 6 节:JSON 编码与字段标签 json.Marshal()
	// ═══════════════════════════════════════════════════════
	// Age 为 0 被省略, Password 被永久忽略, internal 不导出
	encoded, err := json.Marshal(user)
	if err != nil {
		fmt.Println("JSON 编码失败：", err)
		return
	}
	// 注意得到的结果是字节切片类型, 所以还需要多一步 string()
	fmt.Println("JSON 序列化: ", string(encoded)) // {"id":7,"name":"小周"}
	// 若需要区分“未提供 age”和“提供 age=0”, 可以使用属性类型为 *int + omitempty 字段

	// ═══════════════════════════════════════════════════════
	// 第 7 节:JSON 解码:缺失、零值与 null  json.Unmarshal()
	// ═══════════════════════════════════════════════════════
	demonstrateJSONInput := func() {
		user := User{Name: "旧名字", Age: 30}
		// 第一个参数为字节切片类型, 第二个参数为结构体指针
		err := json.Unmarshal([]byte(`{"name":"新名字","extra":true}`), &user)
		if err != nil {
			fmt.Println("JSON 解码失败：", err)
			return
		}
		fmt.Println("缺失字段保留已有值：", user.Name, user.Age) // 新名字 30
		// 结构体中没有 extra 字段所以跳过
		// 结构体 age 没有找到对应字段进行覆盖, 所以旧值保留, 每次请求通常用新的目标结构体, 避免无意复用旧数据
	}

	demonstrateJSONInput()

	// ═══════════════════════════════════════════════════════
	// 第 8 节:局部结构体声明 (有空再看看)
	// ═══════════════════════════════════════════════════════
	// 不需要为每个简单响应声明公共类型
	type Input struct {
		Limit *int `json:"limit,omitempty"`
	}

	// 下方示例再用 *int 区分“没有数值”和“数值为 0”
	for _, raw := range []string{`{}`, `{"limit":0}`, `{"limit":null}`} {
		var input Input
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

	// ═══════════════════════════════════════════════════════
	// 第 9 节:方法值与方法表达式
	// ═══════════════════════════════════════════════════════
	demonstrateMethodValues := func() {
		vvyh := User{
			Name: "vvyh",
		}

		// 方法值
		labelSnapshot := vvyh.Label // 保存当时的 User 副本
		rename := user.Rename_2     // 保存 &user, 后续修改同一实例
		rename("wyh")

		labelSnapshot() // "vvyh"
		user.Label()    // "wyh"

		// 方法表达式
		labelFunction := User.Label
		labelFunction(user) // "wyh"
	}

	demonstrateMethodValues()
}
