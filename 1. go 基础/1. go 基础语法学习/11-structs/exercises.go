// 章节：11-结构体与方法
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

/*
练习 1：购物车条目与金额计算。

背景：结算页面需要按“分”计算金额，避免用浮点数表示人民币小数。

要求：
  - 定义 CartItem，含 SKU string、UnitPriceCents int、Quantity int。
  - 用值接收者 TotalCents() int 返回单价乘数量，不修改条目。
  - 用指针接收者 SetQuantity(quantity int) bool 更新数量。
  - quantity<0 或接收者为 nil 时返回 false，其他情况更新并返回 true。
  - 为聚焦结构体，假设单价非负且乘积在 int 范围内。

建议签名：
func (item CartItem) TotalCents() int
func (item *CartItem) SetQuantity(quantity int) bool

手动核对：
  - 单价 1990，数量 2：金额 3980；修改数量为 3 后金额 5970。
  - 数量改成 -1：返回 false，原数量保持 3。
  - 数量改成 0：返回 true，金额为 0；nil 指针调用更新方法返回 false。
*/
func exercise1() {
	// TODO：定义结构体和方法，打印修改前后数量及金额。
}

/*
练习 2：组织信息与成员信息的组合。

背景：成员卡片同时展示成员昵称和所在组织名称，两者都可能叫 Name。

要求：
  - 定义 TeamInfo{Name string} 与 TeamMember{TeamInfo; Name string}。
  - 给 TeamInfo 定义 TeamLabel() string，返回“团队：”加团队名。
  - 给 TeamMember 定义 CardTitle() string，返回“成员名 / 团队名”。
  - 在 CardTitle 中明确访问嵌入字段的 Name，观察外层字段的遮蔽。
  - 演示直接调用 member.TeamLabel() 的方法提升。

建议签名：func (member TeamMember) CardTitle() string

手动核对：
  - 成员小林、团队平台组：CardTitle 为“小林 / 平台组”。
  - member.Name 为“小林”，member.TeamInfo.Name 为“平台组”。
  - 成员名为空、团队平台组：结果为“ / 平台组”，不隐式填充默认值。
  - member.TeamLabel() 为“团队：平台组”。
*/
func exercise2() {
	// TODO：使用不同于 lesson.go 的类型名，避免同包重复声明。
}

/*
练习 3：安全地构造登录后的响应数据。

背景：后端需要给前端返回 id、displayName、可选 avatarURL，但不能发送密码摘要。

要求：
  - 定义 LoginResponse，字段为 ID int、DisplayName、AvatarURL、PasswordHash string。
  - 配置正确 JSON tag：id、displayName、avatarURL,omitempty、-。
  - 新增 encoding/json 导入，构造数据后调用 json.Marshal，并检查 err。
  - 输出 string(data)，不要直接依赖字节切片的数字打印形式。
  - 用注释说明：JSON 忽略字段不等于可以安全地直接打印整个结构体。

建议签名：func printLoginResponse(response LoginResponse)

手动核对（按 JSON 内容核对，不依赖对象键序）：
  - ID=1、DisplayName=Lin、AvatarURL=""：仅含 id=1 与 displayName=Lin。
  - AvatarURL=/avatars/1.png：响应包含 avatarURL。
  - PasswordHash 无论是否为空，JSON 中都不能有该字段或它的值。
*/
func exercise3() {
	// TODO：写出结构体和标签，手动对照有头像与无头像的 JSON 输出。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type CartItem struct {
	SKU            string
	UnitPriceCents int
	Quantity       int
}

func (item CartItem) TotalCents() int {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：TotalCents；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func (item *CartItem) SetQuantity(quantity int) bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：SetQuantity；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type TeamInfo struct{ Name string }
type TeamMember struct {
	TeamInfo
	Name string
}

func (team TeamInfo) TeamLabel() string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：TeamLabel；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func (member TeamMember) CardTitle() string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：CardTitle；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type LoginResponse struct {
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	AvatarURL    string `json:"avatarURL,omitempty"`
	PasswordHash string `json:"-"`
}

func printLoginResponse(response LoginResponse) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：printLoginResponse；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

// ── 自动检查启动入口：填写习题时无需修改 ──
// runExercises 只启动本章的检查入口；检查规则、测试和参考答案都在 exercises_test.go。
// Go 的普通运行不会编译 _test.go，因此通过 go test 进入 TestExam。
func runExercises() {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法定位本章源码，检查未执行")
		return
	}
	dir := filepath.Dir(source)
	if !filepath.IsAbs(dir) { // 兼容以 -trimpath 构建后在课程根目录或章节目录运行。
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("无法定位章节：", err)
			return
		}
		dir = cwd
		if _, err := os.Stat(filepath.Join(dir, "exercises.go")); err != nil {
			dir = filepath.Join(cwd, filepath.Base(filepath.Dir(source)))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "test", "-json", "-count=1", "-run", "^TestExam$", "-timeout", "5m", ".")
	command.Dir = dir
	command.Env = append(os.Environ(), "GO_LESSON_EXAM=1")
	command.WaitDelay = 2 * time.Second
	output, runErr := command.Output()
	decoder := json.NewDecoder(bytes.NewReader(output))
	var diagnostics strings.Builder
	passed := false
	for {
		var event struct{ Action, Test, Output string }
		if err := decoder.Decode(&event); err != nil {
			break
		}
		diagnostics.WriteString(event.Output)
		if event.Test == "TestExam" {
			if event.Action == "pass" {
				passed = true
			}
			if event.Action == "output" && !strings.HasPrefix(event.Output, "=== RUN") && !strings.HasPrefix(event.Output, "--- PASS") {
				fmt.Print(event.Output)
			}
		}
	}
	if runErr != nil || !passed {
		fmt.Println("练习检查未正常结束：", runErr)
		fmt.Print(diagnostics.String())
		if failure, ok := runErr.(*exec.ExitError); ok {
			fmt.Print(string(failure.Stderr))
		}
	}
}
