// 25-reflection：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestExam 是 exam 模式入口。普通 go test 仍直接运行下面的测试函数。
func TestExam(t *testing.T) {
	if os.Getenv("GO_LESSON_EXAM") != "1" {
		t.Skip("通过 go run . exam 查看逐题检查结果")
	}
	examRun([]examExercise{
		{Test: "TestExercise1", Functions: []string{"dbColumns"}},
		{Test: "TestExercise2", Functions: []string{"updateField"}},
		{Test: "TestExercise3", Functions: []string{"buildUpdate"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	for _, model := range []any{Profile{}, &Profile{}} {
		got, err := dbColumns(model)
		exerciseOK(t, err)
		exerciseEqual(t, got, []string{"id", "name", "age", "active"})
	}
	for _, model := range []any{nil, (*Profile)(nil), 123, new(*Profile)} {
		_, err := dbColumns(model)
		exerciseError(t, err)
	}
	got, err := dbColumns(struct{ Value string }{})
	exerciseOK(t, err)
	exerciseEqual(t, len(got), 0)
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	profile := Profile{Name: "old", Age: 10}
	exerciseOK(t, updateField(&profile, "Name", "new"))
	exerciseEqual(t, profile.Name, "new")
	for _, c := range []struct {
		target any
		field  string
		next   any
	}{{profile, "Name", "bad"}, {nil, "Name", "bad"}, {(*Profile)(nil), "Name", "bad"}, {&profile, "missing", "bad"}, {&profile, "secret", "bad"}, {&profile, "Age", int64(20)}, {&profile, "Name", nil}} {
		before := profile
		exerciseError(t, updateField(c.target, c.field, c.next))
		exerciseEqual(t, profile, before)
	}
	type embedded struct{ Profile }
	target := embedded{Profile: Profile{Name: "old"}}
	exerciseError(t, updateField(&target, "Name", "bad"))
	exerciseEqual(t, target.Name, "old")
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	model := struct {
		Name string `db:"name"`
	}{"O'Reilly"}
	query, args, err := buildUpdate("users", model, 7)
	exerciseOK(t, err)
	exerciseEqual(t, query, "UPDATE `users` SET `name` = ? WHERE `id` = ?")
	exerciseEqual(t, args, []any{"O'Reilly", int64(7)})
	if strings.Contains(query, model.Name) {
		t.Fatal("用户数据被拼入 SQL")
	}
	for _, c := range []struct {
		table string
		model any
		id    int64
	}{{"users", model, 0}, {"users", struct{}{}, 1}, {"users;DROP", model, 1}, {"users", nil, 1}, {"users", (*Profile)(nil), 1}, {"users", struct {
		A string `db:"x"`
		B string `db:"X"`
	}{}, 1}, {"users", struct {
		A string `db:"bad-name"`
	}{}, 1}, {"users", struct {
		A []int `db:"a"`
	}{}, 1}, {"users", struct {
		ID int `db:"id"`
	}{}, 1}} {
		_, _, err := buildUpdate(c.table, c.model, c.id)
		exerciseError(t, err)
	}
	query, args, err = buildUpdate("users", &Profile{ID: 99, Name: "Lin", Age: 20, Active: true}, 7)
	exerciseOK(t, err)
	exerciseEqual(t, query, "UPDATE `users` SET `name` = ?, `age` = ?, `active` = ? WHERE `id` = ?")
	exerciseEqual(t, args, []any{"Lin", int64(20), true, int64(7)})
}

// 以下是测试辅助函数；不属于需要填写的练习答案。
// 未完成的 helper 使用 panic 提醒，这里将它转成明确的测试失败。
func exerciseGuard(t *testing.T) {
	t.Helper()
	if value := recover(); value != nil {
		t.Fatalf("练习执行失败：%v", value)
	}
}
func exerciseEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("得到 %#v；期望 %#v", got, want)
	}
}
func exerciseOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
}
func exerciseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("期望返回错误，实际得到 nil")
	}
}

// ── 讲解配套测试：验证 lesson.go 中的示例 ──
func TestSetValue(t *testing.T) {
	name := "old"
	if err := setValue(&name, "new"); err != nil || name != "new" {
		t.Fatalf("name=%q err=%v", name, err)
	}
	var id UserID = 7
	if err := setValue(&id, int64(8)); err == nil || id != 7 {
		t.Fatalf("相同 Kind 的不同定义类型不应赋值：id=%v err=%v", id, err)
	}
	if err := setValue(&id, UserID(8)); err != nil || id != 8 {
		t.Fatalf("定义类型自身赋值失败：id=%v err=%v", id, err)
	}
	for _, destination := range []any{nil, (*string)(nil), "copy"} {
		if err := setValue(destination, "next"); err == nil {
			t.Errorf("目标 %T 应失败", destination)
		}
	}
	if err := setValue(&name, nil); err == nil || name != "new" {
		t.Fatal("无类型 nil 应失败且原值不变")
	}
	profile := Profile{secret: "private"}
	for _, field := range []string{"secret", "missing"} {
		if err := assignValue(reflect.ValueOf(&profile).Elem().FieldByName(field), "changed"); err == nil {
			t.Errorf("字段 %q 不应允许赋值", field)
		}
	}
}

func TestBuildInsertParameters(t *testing.T) {
	// 输入含 SQL 字符，仍应完整保留为参数，不得改变语句结构。
	name := "O'Reilly'); DROP TABLE users; --"
	model := Profile{ID: 7, Name: name, Age: 20, Active: true, Note: "skip", secret: "private"}
	wantSQL := "INSERT INTO `users` (`id`, `name`, `age`, `active`) VALUES (?, ?, ?, ?)"
	for _, input := range []any{model, &model} {
		query, args, err := buildInsert("users", input)
		if err != nil {
			t.Fatal(err)
		}
		if query != wantSQL || strings.Contains(query, name) {
			t.Fatalf("SQL=%q", query)
		}
		if want := []any{int64(7), name, int64(20), true}; !reflect.DeepEqual(args, want) {
			t.Fatalf("args=%#v want=%#v", args, want)
		}
	}
}

func TestBuildInsertRejectsInvalidModels(t *testing.T) {
	cases := []struct {
		name  string
		table string
		model any
	}{
		{"nil", "users", nil},
		{"nil pointer", "users", (*Profile)(nil)},
		{"scalar", "users", 1},
		{"invalid table", "users; DROP TABLE users", Profile{}},
		{"no columns", "users", struct{ Name string }{}},
		{"invalid column", "users", struct {
			Name string `db:"bad-name"`
		}{}},
		{"duplicate columns", "users", struct {
			A string `db:"name"`
			B string `db:"Name"`
		}{}},
		{"unsupported field", "users", struct {
			Names []string `db:"names"`
		}{}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			query, args, err := buildInsert(test.table, test.model)
			if err == nil || query != "" || args != nil {
				t.Fatalf("失败不得返回可执行的半成品：query=%q args=%v err=%v", query, args, err)
			}
		})
	}
}

// ── 本章独立检查器：学习时可跳过此部分 ──
// examExercise 对应一道练习及其验收测试。方法名使用 Type.Method，省略接收者的 *。
// 一题依赖多个函数时，必须全部填写后才能用该题测试检查整体行为。
type examExercise struct {
	Test      string
	Functions []string
	File      string   // 默认 exercises.go；少数题目实际修改 lesson.go 或测试文件。
	Variables []string // 例如部署练习要求新增的 buildTime。
}

const (
	examUnfinished = "还未完成"
	examIncorrect  = "错误"
	examComplete   = "完成"
)

type examResult struct {
	status string
	detail string
}

type examTestRunner func(dir, name string) error

// examRun 在 go test 的本章工作目录内运行，逐题调用独立测试进程。
func examRun(exercises []examExercise) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法定位本章目录：", err)
		return
	}
	examReport(os.Stdout, dir, exercises, examRunTest)
}

func examReport(out io.Writer, dir string, exercises []examExercise, run examTestRunner) {
	fmt.Fprintf(out, "%s：检查全部练习\n", filepath.Base(dir))
	counts := map[string]int{}
	for i, exercise := range exercises {
		checked := examCheck(dir, exercise, run)
		counts[checked.status]++
		names := append(append([]string{}, exercise.Functions...), exercise.Variables...)
		fmt.Fprintf(out, "练习 %d（%s）：%s\n", i+1, strings.Join(names, "、"), checked.status)
		if checked.detail != "" {
			for _, line := range strings.Split(strings.TrimSpace(checked.detail), "\n") {
				fmt.Fprintf(out, "  %s\n", line)
			}
		}
	}
	fmt.Fprintf(out, "汇总：完成 %d 题，错误 %d 题，还未完成 %d 题\n",
		counts[examComplete], counts[examIncorrect], counts[examUnfinished])
}

func examCheck(dir string, exercise examExercise, run examTestRunner) examResult {
	file := exercise.File
	if file == "" {
		file = "exercises.go"
	}
	source, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, file), nil, parser.SkipObjectResolution)
	if err != nil {
		return examResult{examIncorrect, "无法检查源码：" + err.Error()}
	}
	functions := examFunctionsIn(source)
	var missing []string
	for _, name := range exercise.Functions {
		if examIsEmpty(functions[name]) {
			missing = append(missing, name)
		}
	}
	for _, name := range exercise.Variables {
		if !examHasVariable(source, name) {
			missing = append(missing, name)
		}
	}
	if len(missing) != 0 {
		return examResult{examUnfinished, "请填写 " + file + " 中的 " + strings.Join(missing, "、")}
	}
	tests, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, "exercises_test.go"), nil, parser.SkipObjectResolution)
	if err != nil {
		return examResult{examIncorrect, "无法读取练习测试：" + err.Error()}
	}
	if examIsEmpty(examFunctionsIn(tests)[exercise.Test]) {
		return examResult{examIncorrect, "验收测试 " + exercise.Test + " 缺失或为空，不能判定完成"}
	}
	if err := run(dir, exercise.Test); err != nil {
		return examResult{examIncorrect, err.Error()}
	}
	return examResult{status: examComplete}
}

func examFunctionsIn(file *ast.File) map[string]*ast.FuncDecl {
	functions := make(map[string]*ast.FuncDecl)
	for _, declaration := range file.Decls {
		fn, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}
		name := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			name = examReceiverName(fn.Recv.List[0].Type) + "." + name
		}
		functions[name] = fn
	}
	return functions
}

func examReceiverName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.StarExpr:
		return examReceiverName(value.X)
	case *ast.IndexExpr:
		return examReceiverName(value.X)
	case *ast.IndexListExpr:
		return examReceiverName(value.X)
	}
	return ""
}

func examHasVariable(file *ast.File, name string) bool {
	for _, declaration := range file.Decls {
		group, ok := declaration.(*ast.GenDecl)
		if !ok || group.Tok != token.VAR {
			continue
		}
		for _, spec := range group.Specs {
			for _, variable := range spec.(*ast.ValueSpec).Names {
				if variable.Name == name {
					return true
				}
			}
		}
	}
	return false
}

func examIsEmpty(fn *ast.FuncDecl) bool {
	if fn == nil || fn.Body == nil {
		return true
	}
	for _, statement := range fn.Body.List {
		if !examIsPlaceholder(statement) {
			return false
		}
	}
	return true
}

// 只忽略空语句和仓库原有占位调用。不能因注释仍有 TODO，便忽略已经写入的实现。
func examIsPlaceholder(statement ast.Stmt) bool {
	if _, ok := statement.(*ast.EmptyStmt); ok {
		return true
	}
	expression, ok := statement.(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := expression.X.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}
	literal, ok := call.Args[0].(*ast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return false
	}
	text, err := strconv.Unquote(literal.Value)
	if err != nil || !strings.Contains(text, "待完成") {
		return false
	}
	switch function := call.Fun.(type) {
	case *ast.Ident:
		return function.Name == "panic"
	case *ast.SelectorExpr:
		pkg, ok := function.X.(*ast.Ident)
		return ok && pkg.Name == "fmt" && (function.Sel.Name == "Println" || function.Sel.Name == "Print")
	}
	return false
}

func examRunTest(dir, name string) error {
	return examRunTestWithTimeout(dir, name, 30*time.Second, 90*time.Second)
}

func examRunTestWithTimeout(dir, name string, testTimeout, commandTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "test", "-json", "-count=1",
		"-run", "^"+regexp.QuoteMeta(name)+"$", "-timeout", testTimeout.String(), ".")
	command.Dir = dir
	command.WaitDelay = 2 * time.Second
	output := &examTestEvents{name: name}
	var stderr examLimitedText
	command.Stdout = output
	command.Stderr = &stderr
	err := command.Run()
	if ctx.Err() != nil {
		return fmt.Errorf("%s 检查超时（含构建时间）：%v", name, ctx.Err())
	}
	if err != nil {
		detail := strings.TrimSpace(output.diagnostics.String() + "\n" + stderr.String())
		return fmt.Errorf("%s 未通过：%v\n%s", name, err, detail)
	}
	if !output.ran || !output.passed || output.skipped {
		return fmt.Errorf("%s 未完整执行或包含跳过的用例，不能判定完成\n%s", name, output.diagnostics.String())
	}
	return nil
}

// 限制诊断和单行缓冲大小，错误练习不断输出时也不会无限占用内存。
type examLimitedText struct{ bytes.Buffer }

func (text *examLimitedText) Write(data []byte) (int, error) {
	n := len(data)
	remaining := 32*1024 - text.Len()
	if remaining > 0 {
		text.Buffer.Write(data[:min(remaining, n)])
	}
	return n, nil
}

type examTestEvents struct {
	name        string
	line        []byte
	dropping    bool
	ran         bool
	passed      bool
	skipped     bool
	diagnostics examLimitedText
}

func (events *examTestEvents) Write(data []byte) (int, error) {
	n := len(data)
	for _, b := range data {
		if b == '\n' {
			if !events.dropping {
				events.accept(events.line)
			}
			events.line = events.line[:0]
			events.dropping = false
		} else if !events.dropping {
			if len(events.line) < 1024*1024 {
				events.line = append(events.line, b)
			} else {
				events.dropping = true
			}
		}
	}
	return n, nil
}

func (events *examTestEvents) accept(line []byte) {
	var event struct{ Action, Test, Output string }
	if json.Unmarshal(line, &event) != nil {
		events.diagnostics.Write(append(line, '\n'))
		return
	}
	if event.Test == events.name {
		if event.Action == "run" {
			events.ran = true
		}
		if event.Action == "pass" {
			events.passed = true
		}
	}
	if event.Action == "skip" && (event.Test == events.name || strings.HasPrefix(event.Test, events.name+"/")) {
		events.skipped = true
	}
	if event.Output != "" {
		events.diagnostics.Write([]byte(event.Output))
	}
}

/*
参考答案（先自己完成，再展开核对）

将下面同名函数的实现填入 exercises.go，并合并需要的 import；不要重复声明类型或函数。
类型和字段已提供骨架，题目中的语法要求与思考题仍需自行解释。
第 01–04 章替换 exercise1/2/3；其他章节实现对应 helper 后可从练习入口调用演示。

package main

import (
	"fmt"
	"reflect"
	"strings"
)

func dbColumns(model any) ([]string, error) {
	typ := reflect.TypeOf(model)
	if typ == nil {
		return nil, fmt.Errorf("模型为 nil")
	}
	if typ.Kind() == reflect.Pointer {
		if reflect.ValueOf(model).IsNil() {
			return nil, fmt.Errorf("模型指针为 nil")
		}
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("模型不是结构体")
	}
	var columns []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		column := field.Tag.Get("db")
		if field.IsExported() && column != "" && column != "-" {
			columns = append(columns, column)
		}
	}
	return columns, nil
	// Type.Field 提供字段名、标签等元数据；Value.Field 提供某个实例的字段值。
}
func updateField(target any, fieldName string, next any) error {
	value := reflect.ValueOf(target)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf("需要非 nil 指针")
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("需要结构体指针")
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if field.Name != fieldName {
			continue
		}
		destination := value.Field(i)
		source := reflect.ValueOf(next)
		if !field.IsExported() || !destination.CanSet() || !source.IsValid() || !source.Type().AssignableTo(destination.Type()) {
			return fmt.Errorf("字段不可写或类型不兼容")
		}
		destination.Set(source)
		return nil
	}
	return fmt.Errorf("未知字段")
}
func buildUpdate(table string, model any, id int64) (string, []any, error) {
	if id <= 0 || !sqlIdentifier.MatchString(table) {
		return "", nil, fmt.Errorf("表名或 id 非法")
	}
	value := reflect.ValueOf(model)
	if value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "", nil, fmt.Errorf("nil 模型")
		}
		value = value.Elem()
	}
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("需要结构体")
	}
	var sets []string
	var args []any
	seen := make(map[string]bool)
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		column := field.Tag.Get("db")
		if !field.IsExported() || column == "" || column == "-" {
			continue
		}
		key := strings.ToLower(column)
		if !sqlIdentifier.MatchString(column) || seen[key] {
			return "", nil, fmt.Errorf("非法或重复列")
		}
		seen[key] = true
		if key == "id" {
			continue
		}
		data := value.Field(i)
		var arg any
		switch data.Kind() {
		case reflect.String:
			arg = data.String()
		case reflect.Bool:
			arg = data.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			arg = data.Int()
		case reflect.Float32, reflect.Float64:
			arg = data.Float()
		default:
			return "", nil, fmt.Errorf("不支持字段类型 %v", data.Type())
		}
		sets = append(sets, "`"+column+"` = ?")
		args = append(args, arg)
	}
	if len(sets) == 0 {
		return "", nil, fmt.Errorf("无可更新字段")
	}
	args = append(args, id)
	return "UPDATE `" + table + "` SET " + strings.Join(sets, ", ") + " WHERE `id` = ?", args, nil
}

*/
