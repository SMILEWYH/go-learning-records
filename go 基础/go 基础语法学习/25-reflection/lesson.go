/*
章节：25-反射与ORM原理
	本章掌握清单：
	1. Type、Value 与 Kind
	2. 通过指针修改原变量
	3. 字段元数据、字段值与 tag
	4. 动态方法调用与方法集
	5. ORM 原理：生成 SQL 结构与独立参数
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./25-reflection
// 运行练习：go run ./25-reflection exam
// 阅读顺序：先读 demoReflection 的五个小节，再回看对应 helper；buildInsert 较长，按校验、遍历、组装三步阅读。
// TS 对照：TS 类型大多在编译后擦除；Go reflect 可在运行时读取具体类型、字段和 tag。
// 后续再学：嵌入字段路径、元数据缓存、database/sql、事务和真正的 ORM。

package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type UserID int64

type Profile struct {
	ID     UserID `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Age    int    `db:"age" json:"age"`
	Active bool   `db:"active" json:"active"`
	Note   string `db:"-" json:"-"`
	secret string // 未导出字段：可以查看元数据，但不能随意 Interface 或 Set。
}

func (profile Profile) Greeting(prefix string) string { return prefix + profile.Name }
func (profile *Profile) Rename(name string)           { profile.Name = name }

func describe(value any) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if !v.IsValid() { // ValueOf(nil) 是无效 Value；TypeOf(nil) 返回 nil。
		fmt.Println("nil：没有动态类型，IsValid=false")
		return
	}
	fmt.Printf("Type=%v，Kind=%v，值=", t, t.Kind())
	// Type 区分 int64 和 UserID；Kind 将两者都归入 reflect.Int64。
	// Kind() 是方法，case 使用 reflect.String 等常量。
	switch v.Kind() {
	case reflect.String:
		fmt.Println(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Println(v.Int()) // 返回 int64，包括定义的整数类型。
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Println(v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Println(v.Float())
	case reflect.Bool:
		fmt.Println(v.Bool())
	case reflect.Pointer:
		fmt.Println("指针，IsNil=", v.IsNil())
	default:
		fmt.Println(v.Interface()) // 这里的 v 来自 ValueOf，允许 Interface。
	}
	// String() 不是通用格式化方法；在错误 Kind 上调用 Int/Bool 等通常会 panic。
	// IsNil 只用于 pointer/interface/map/slice/chan/func 等允许为 nil 的 Kind。
}

// assignValue 要求新值可直接赋给目标，不自动把 int 转成 int64 或 UserID。
// CanSet + AssignableTo 比“只比较 Kind”准确；相同 Kind 不代表相同类型。
func assignValue(target reflect.Value, next any) error {
	if !target.IsValid() || !target.CanSet() {
		return errors.New("目标不可修改")
	}
	source := reflect.ValueOf(next)
	if !source.IsValid() {
		return errors.New("本示例不接受无类型 nil")
	}
	if !source.Type().AssignableTo(target.Type()) {
		return fmt.Errorf("不能把 %v 赋给 %v", source.Type(), target.Type())
	}
	switch target.Kind() {
	case reflect.String:
		target.SetString(source.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		target.SetInt(source.Int())
	case reflect.Bool:
		target.SetBool(source.Bool())
	default:
		target.Set(source) // 通用写法；上面的分支只是演示专用 setter。
	}
	// SetString 接收 string，但不要求“必须类型断言”：来源可以本来就是 string，
	// 也可以像这里用 source.String()。若用 next.(string)，应先检查断言的 ok。
	// 如果另行允许数值转换，需考虑 OverflowInt 等；Convert 并不会替你拒绝数值截断。
	return nil
}

func setValue(destination, next any) error {
	v := reflect.ValueOf(destination)
	if !v.IsValid() || v.Kind() != reflect.Pointer || v.IsNil() {
		return errors.New("请传入非 nil 指针")
	}
	return assignValue(v.Elem(), next) // Elem() 解引用，取得指针指向的可修改变量。
}

func showFields(value Profile) {
	t, v := reflect.TypeOf(value), reflect.ValueOf(value)
	for i := 0; i < t.NumField(); i++ {
		field, data := t.Field(i), v.Field(i)
		fmt.Printf("字段=%s，类型=%v，db=%q，json=%q", field.Name, field.Type,
			field.Tag.Get("db"), field.Tag.Get("json"))
		if data.CanInterface() {
			fmt.Printf("，值=%v\n", data.Interface())
		} else {
			fmt.Println("，未导出，跳过读取")
		}
	}
	// Type.Field 返回 StructField 元数据；Value.Field 返回对应字段的值。
	// NumField/Field 只能用于结构体；如果传入指针，要先检查类型再 Elem()。
	// Tag.Get 不区分“缺少标签”和“标签值为空”；需要区分时用 Tag.Lookup。
}

func demonstrateMethods(profile *Profile) error {
	v, t := reflect.ValueOf(profile), reflect.TypeOf(profile)
	for i := 0; i < v.NumMethod(); i++ {
		fmt.Println("可见方法：", t.Method(i).Name, v.Method(i).Type())
	}
	fmt.Println("值的方法数 / 指针的方法数：", reflect.TypeOf(*profile).NumMethod(), t.NumMethod()) // 1 / 2
	// *Profile 的方法集中含 Greeting 和 Rename，Profile 只有 Greeting。
	// 用名称选择方法比硬编码下标清楚；未导出方法不会作为这里的可调用方法列出。
	method := v.MethodByName("Greeting")
	if !method.IsValid() {
		return errors.New("方法不存在")
	}
	if method.Type().NumIn() != 1 || method.Type().In(0) != reflect.TypeOf("") {
		return errors.New("Greeting 参数签名不匹配")
	}
	results := method.Call([]reflect.Value{reflect.ValueOf("你好，")})
	fmt.Println("反射调用：", results[0].String())
	// Value.Method/MethodByName 取得的函数已经绑定接收者，只传业务参数。
	// Type.Method(i).Func 尚未绑定接收者，Call 时还需要把接收者放在第一个参数。
	// 不是 Method(i).方法名()；Call 的参数数量或类型不匹配会 panic。
	return nil
}

var sqlIdentifier = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// buildInsert 只演示“结构体 -> 列名、占位符、参数”，不连接或创建数据库。
// ORM（对象关系映射）把结构体对应到表，把字段对应到列；真实 ORM 还负责执行、查询映射等。
// 本例使用 MySQL 风格的反引号和 ?，只支持一层结构体、显式 db 标签和基础标量。
// 不支持嵌套结构体、时间、NULL、自动主键或事务；标签为 "-"、缺少标签、未导出字段均跳过。
func buildInsert(table string, model any) (string, []any, error) {
	if !sqlIdentifier.MatchString(table) {
		return "", nil, errors.New("表名只能使用字母、数字和下划线，且不能以数字开头")
	}
	v := reflect.ValueOf(model)
	if v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", nil, errors.New("模型指针不能为 nil")
		}
		v = v.Elem()
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return "", nil, errors.New("模型必须是结构体或它的非 nil 指针")
	}
	t := v.Type()
	var columns, placeholders []string
	var args []any
	seen := make(map[string]bool)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column := field.Tag.Get("db")
		if !field.IsExported() || column == "" || column == "-" {
			continue
		}
		if !sqlIdentifier.MatchString(column) || seen[strings.ToLower(column)] {
			return "", nil, fmt.Errorf("非法或重复列名：%q", column)
		}
		seen[strings.ToLower(column)] = true
		data := v.Field(i)
		var argument any
		switch data.Kind() {
		case reflect.String:
			argument = data.String()
		case reflect.Bool:
			argument = data.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			argument = data.Int() // 将 UserID 等定义类型归一为 int64。
		case reflect.Float32, reflect.Float64:
			argument = data.Float()
		default:
			return "", nil, fmt.Errorf("字段 %s 的类型 %v 暂不支持", field.Name, data.Type())
		}
		columns = append(columns, "`"+column+"`")
		placeholders = append(placeholders, "?")
		args = append(args, argument)
	}
	if len(columns) == 0 {
		return "", nil, errors.New("没有可插入的字段")
	}
	query := "INSERT INTO `" + table + "` (" + strings.Join(columns, ", ") + ") VALUES (" +
		strings.Join(placeholders, ", ") + ")"
	return query, args, nil
	// 将来学 database/sql 时调用 db.ExecContext(ctx, query, args...)。
	// 值不拼进 SQL；单引号等由数据库驱动处理。表名/列名不能用 ? 占位，所以单独限制。
	// 合法标识符不等于有访问权限；真实应用还要明确允许访问哪些表。
}

func demoReflection() error {

	// ═══════════════════════════════════════════════════════
	// 第 1 节：Type、Value 与 Kind
	// ═══════════════════════════════════════════════════════
	// describe 同时取得类型和值：Type 保留 UserID 身份，Kind 只把它归为 Int64 类别。
	// nil 没有动态类型，而 (*Profile)(nil) 有指针类型；反射前先分清两者。

	for _, value := range []any{"Go", 42, UserID(7), true, 3.5, nil, (*Profile)(nil)} {
		describe(value)
	}

	// ═══════════════════════════════════════════════════════
	// 第 2 节：通过指针修改原变量
	// ═══════════════════════════════════════════════════════
	// ValueOf(name) 包装的是副本；ValueOf(&name).Elem() 才能取得原变量的可写位置。
	// setValue 先检查指针和 nil，assignValue 再检查 CanSet 与赋值兼容性。

	name := "旧名字"
	fmt.Println("值副本 CanSet：", reflect.ValueOf(name).CanSet()) // false
	if err := setValue(&name, "新名字"); err != nil {
		return err
	}
	fmt.Println("修改原变量：", name)
	fmt.Println("传值失败：", setValue(name, "不会写入"))
	fmt.Println("类型不匹配：", setValue(&name, 123))

	// ═══════════════════════════════════════════════════════
	// 第 3 节：字段元数据、字段值与 tag
	// ═══════════════════════════════════════════════════════
	// showFields 中 Type.Field(i) 提供字段名、类型、tag；Value.Field(i) 提供这一项的值。
	// 读取未导出字段前检查 CanInterface，修改字段前检查 CanSet。

	profile := Profile{ID: 7, Name: "小林", Age: 20, Active: true, secret: "private"}
	showFields(profile)
	field := reflect.ValueOf(&profile).Elem().FieldByName("Name")
	if err := assignValue(field, "小周"); err != nil {
		return err
	}
	fmt.Println("修改结构体字段：", profile.Name)

	// ═══════════════════════════════════════════════════════
	// 第 4 节：动态方法调用与方法集
	// ═══════════════════════════════════════════════════════
	// MethodByName 得到已经绑定接收者的函数，Call 接收 []reflect.Value 参数并返回同类结果列表。
	// 调用前要确认方法存在、参数个数和类型正确，否则运行时会 panic。

	if err := demonstrateMethods(&profile); err != nil {
		return err
	}

	// ═══════════════════════════════════════════════════════
	// 第 5 节：ORM 原理：生成 SQL 结构与独立参数
	// ═══════════════════════════════════════════════════════
	// buildInsert 依次扫描可导出的 db 字段，分别收集列名、占位符和参数。
	// 名字 O'Reilly 始终留在 args 中；SQL 只放 ?，交给未来的数据库驱动绑定。
	// 反射只解决结构检查与值提取，本例尚未执行数据库操作。

	profile.Name = "O'Reilly" // 含单引号仍留在参数里，不能手动拼到 SQL。
	query, args, err := buildInsert("users", profile)
	if err != nil {
		return err
	}
	fmt.Println("SQL：", query)
	fmt.Printf("参数：%#v\n", args)
	// INSERT INTO `users` (`id`, `name`, `age`, `active`) VALUES (?, ?, ?, ?)
	// 参数顺序对应字段声明顺序：7、O'Reilly、20、true。
	return nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := demoReflection(); err != nil {
		fmt.Fprintln(os.Stderr, "反射演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 对比 setValue(name, ...) 与 setValue(&name, ...)，解释 Elem 为什么需要先有指针。
// 给 Profile 某字段的 db 标签改成 "-"，检查 SQL 列数、占位符数量和参数数量是否一起减少。
