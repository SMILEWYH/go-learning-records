// 章节：12-自定义类型与别名
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
练习 1：区分租户 ID 和项目 ID。

背景：多租户后台里，两个参数都是字符串，但交换顺序会查询到错误资源。

要求：
  - 定义 TenantID string 与 ProjectID string，不能使用等号别名写法。
  - 构造“/tenants/<tenant>/projects/<project>”形式的路径。
  - 任一 ID 为空时返回空字符串和 false，正常时返回路径和 true。
  - 为聚焦类型，假设非空输入仅含英文字母、数字或短横线，无需 URL 转义。
  - 在注释中写出一次错误类型调用并解释原因，保持提交的代码可编译。

建议签名：func projectPath(tenant TenantID, project ProjectID) (string, bool)

手动核对：
  - TenantID("t1"), ProjectID("p1")：/tenants/t1/projects/p1，true。
  - 任一值为 ""：""，false。
  - 两个 ID 内容同为 x 也仍是不同类型；同样内容不代表同样业务含义。
*/
func exercise1() {
	// TODO：定义新类型及 helper；把普通 string 转换到对应领域类型后调用。
}

/*
练习 2：用有单位的金额类型表示 API 数据。

背景：前端传来的金额单位约定为分，后端不应把“分”误当成“元”。

要求：
  - 定义 Cents int64，并提供方法 YuanText() string，输出精确的两位小数。
  - 使用整数除法和余数拆分元、分；不要通过 float64 计算金额。
  - 提供 parseCents(raw int64) (Cents, bool)，拒绝负数，正常返回对应值。
  - 本题格式化方法只要求处理非负值；所有金额先经过 parseCents。
  - 思考：直接 Cents(-1) 为什么仍能编译，定义类型能否代替业务校验？

建议签名：func (amount Cents) YuanText() string

手动核对：
  - 1990 -> "19.90"；105 -> "1.05"；1 -> "0.01"。
  - 0 -> "0.00"；parseCents(-1) 返回 0、false。
  - 可以使用 fmt.Sprintf 的 %02d 让“分”始终占两位。
*/
func exercise2() {
	// TODO：定义类型、校验函数和方法，避免浮点计算与隐式单位换算。
}

/*
练习 3：迁移 API 字段类型，同时兼容旧名称。

背景：内部准备把任务编号从旧名称 JobKey 改名为 TaskKey，但现有调用不能中断。

要求：
  - 定义 TaskKey string，并提供 IsValid() bool：非空即为合法。
  - 用别名 type JobKey = TaskKey 保留旧名字。
  - 证明两者可以直接赋值，旧名字创建的值也可调用 IsValid。
  - 额外定义 ExternalTaskKey TaskKey，观察它需要转换且不继承 IsValid。
  - 无法编译的示例留在注释里，不要破坏本章构建。

建议签名：func (key TaskKey) IsValid() bool

手动核对：
  - JobKey("job-1") 赋给 TaskKey 变量无需转换，IsValid 返回 true。
  - JobKey("").IsValid() 返回 false。
  - ExternalTaskKey("job-1") 转回 TaskKey 后，IsValid 返回 true。
  - 用一行自己的话解释“新名字”和“新类型”的区别。
*/
func exercise3() {
	// TODO：定义三个名字，用赋值、转换和方法调用展示它们的关系。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

type TenantID string
type ProjectID string

func projectPath(tenant TenantID, project ProjectID) (string, bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：projectPath；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type Cents int64

func (amount Cents) YuanText() string {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：YuanText；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}
func parseCents(raw int64) (Cents, bool) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：parseCents；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

type TaskKey string
type JobKey = TaskKey
type ExternalTaskKey TaskKey

func (key TaskKey) IsValid() bool {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：IsValid；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
