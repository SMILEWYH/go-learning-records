// 章节：25-反射与ORM原理
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

// 练习 1：提取 API 模型的字段清单。
// 建议：func dbColumns(model any) ([]string, error)。只处理结构体或一层非 nil 指针。
// 要求：按声明顺序提取已导出字段的 db 标签；缺失、空标签、"-" 均跳过。
// 常规：Profile{} 与 &Profile{} -> ["id","name","age","active"]。
// 边界：nil、(*Profile)(nil)、123 -> error；没有标签的结构体 -> 长度 0。
// 在注释中说明为什么 Type.Field 与 Value.Field 不能互换。
func exercise1() {
	// TODO：检查输入再读取结构体元数据。
}

// 练习 2：按字段名修改配置。
// 建议：func updateField(target any, fieldName string, next any) error。
// 要求：接收非 nil 结构体指针，只匹配该结构体直接声明的字段；拒绝未知/未导出字段和不兼容类型。
// 常规：&Profile{Name:"old"}, "Name", "new" -> Name 为 new。
// 边界：传结构体值、nil、"missing"、"secret"、给 Age 传 int64(20) -> error 且原值不变。
// 提示：先 Kind/IsNil，再 Elem/CanSet/AssignableTo；不要仅比较 Kind。
func exercise2() {
	// TODO：自己封装 helper，覆盖不能写和类型不符的情况。
}

// 练习 3：为小 ORM 添加按 id 更新的语句生成器。
// 建议：func buildUpdate(table string, model any, id int64) (string, []any, error)。
// 要求：参考讲解中的字段规则，但跳过 db:"id"；id 必须大于 0；不能省略 WHERE。
// 只生成 SQL 和 args，不执行数据库操作；检查表名、列名、重复列、无可更新字段及不支持类型。
// 常规：模型只有 Name string `db:"name"`，值为 O'Reilly，id=7：
// UPDATE `users` SET `name` = ? WHERE `id` = ?；参数依次为 "O'Reilly"、int64(7)。
// 边界：id=0、空结构体、非法表名、nil 模型 -> error；用户输入不得出现在 SQL 字符串中。
func exercise3() {
	// TODO：实现 SQL 与参数分离，并在 exercises_test.go 中断言边界。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func dbColumns(model any) ([]string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：dbColumns；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func updateField(target any, fieldName string, next any) error {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：updateField；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func buildUpdate(table string, model any, id int64) (string, []any, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：buildUpdate；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
