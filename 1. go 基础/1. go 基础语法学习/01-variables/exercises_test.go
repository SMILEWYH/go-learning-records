// 01-variables：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 输出按下方 want 顺序排列；空白数量不限，不加额外标题。
// 前四章直接检查 exercise1/2/3 的输出。请按测试中的顺序输出各组数据，不加提示标题；空白数量不限。语法用法与思考题仍需结合题目自行核对。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		{Test: "TestExercise1", Functions: []string{"exercise1"}},
		{Test: "TestExercise2", Functions: []string{"exercise2"}},
		{Test: "TestExercise3", Functions: []string{"exercise3"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	exerciseOutput(t, exercise1, "profile-api development true 9000 8080\nstring string bool int int\n\"\" false")
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	exerciseOutput(t, exercise2, "0 \"\" false\n3 \"ap-shanghai\" true\n0")
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	exerciseOutput(t, exercise3, "preview\nproduction\nstaging")
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

// 通过临时文件捕获输出，避免大量输出堵塞管道；本文件不能调用 t.Parallel。
func exerciseCapture(t *testing.T, run func()) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "stdout")
	exerciseOK(t, err)
	original := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = original; file.Close() }()
	run()
	os.Stdout = original
	_, err = file.Seek(0, 0)
	exerciseOK(t, err)
	data, err := io.ReadAll(file)
	exerciseOK(t, err)
	return string(data)
}
func exerciseOutput(t *testing.T, run func(), want string) {
	t.Helper()
	got := exerciseCapture(t, run)
	exerciseEqual(t, strings.Fields(got), strings.Fields(want))
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

// ── 自动检查器回归测试 ──
func writeFixture(t *testing.T, dir, name, source string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestUnfinishedDetection(t *testing.T) {
	for _, test := range []struct {
		name, source, function string
		want                   string
	}{
		{"empty", "func answer() {}", "answer", examUnfinished},
		{"comments", "func answer() { /* TODO: implement */ }", "answer", examUnfinished},
		{"placeholder output", `func answer() { fmt.Println("练习 1（待完成）") }`, "answer", examUnfinished},
		{"placeholder panic", `func answer() int { panic("待完成：answer") }`, "answer", examUnfinished},
		{"missing", "func other() {}", "answer", examUnfinished},
		{"implemented with TODO", "func answer() int { /* TODO: more examples */ return 42 }", "answer", examComplete},
		{"bare return is code", "func answer() { return }", "answer", examComplete},
		{"real panic is code", `func answer() { panic("unexpected") }`, "answer", examComplete},
		{"partial implementation", `func answer() int { x := 42; _ = x; panic("待完成：answer") }`, "answer", examComplete},
		{"value method", "type Answer int; func (Answer) Value() int { return 42 }", "Answer.Value", examComplete},
		{"pointer method", "type Answer int; func (*Answer) Value() int { return 42 }", "Answer.Value", examComplete},
		{"generic method", "type Answer[T any] struct{}; func (*Answer[T]) Value() int { return 42 }", "Answer.Value", examComplete},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFixture(t, dir, "exercises.go", "package main\n"+test.source)
			writeFixture(t, dir, "exercises_test.go", "package main\nfunc TestExercise1(t *testing.T) { answer() }")
			calls := 0
			got := examCheck(dir, examExercise{Test: "TestExercise1", Functions: []string{test.function}}, func(_, _ string) error {
				calls++
				return nil
			})
			if got.status != test.want {
				t.Fatalf("got %+v; want %s", got, test.want)
			}
			if (calls == 1) != (test.want == examComplete) {
				t.Fatalf("unexpected number of test calls: %d", calls)
			}
		})
	}
}

func TestReportAllThreeStatesAndContinueAfterFailure(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "exercises.go", `package main
func first() { /* not implemented */ }
func second() int { return -1 }
func third() int { return 42 }
`)
	writeFixture(t, dir, "exercises_test.go", `package main
func TestExercise1(t *testing.T) { first() }
func TestExercise2(t *testing.T) { second() }
func TestExercise3(t *testing.T) { third() }
`)
	var output bytes.Buffer
	var executed []string
	examReport(&output, dir, []examExercise{
		{Test: "TestExercise1", Functions: []string{"first"}},
		{Test: "TestExercise2", Functions: []string{"second"}},
		{Test: "TestExercise3", Functions: []string{"third"}},
	}, func(actualDir, name string) error {
		if actualDir != dir {
			t.Fatalf("test directory = %s, want %s", actualDir, dir)
		}
		executed = append(executed, name)
		if name == "TestExercise2" {
			return errors.New("got -1; want 42")
		}
		return nil
	})
	for _, want := range []string{
		"练习 1（first）：还未完成", "练习 2（second）：错误", "got -1; want 42",
		"练习 3（third）：完成", "汇总：完成 1 题，错误 1 题，还未完成 1 题",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("output missing %q:\n%s", want, output.String())
		}
	}
	if strings.Join(executed, ",") != "TestExercise2,TestExercise3" {
		t.Fatalf("wrong tests executed: %v", executed)
	}
}

func TestMultipleFunctionsAndAlternateSources(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "exercises.go", "package main\nfunc first() int { return 1 }; func second() {}")
	writeFixture(t, dir, "lesson.go", "package main\nfunc version() { return }")
	writeFixture(t, dir, "exercises_test.go", "package main\nfunc TestExercise1(t *testing.T) { version() }")
	run := func(_, _ string) error { return nil }
	got := examCheck(dir, examExercise{Test: "TestExercise1", Functions: []string{"first", "second"}}, run)
	if got.status != examUnfinished || !strings.Contains(got.detail, "second") {
		t.Fatalf("partially filled group: %+v", got)
	}
	version := examExercise{Test: "TestExercise1", File: "lesson.go", Functions: []string{"version"}, Variables: []string{"buildTime"}}
	if got := examCheck(dir, version, run); got.status != examUnfinished {
		t.Fatalf("missing buildTime: %+v", got)
	}
	writeFixture(t, dir, "lesson.go", "package main\nvar buildTime = \"unknown\"\nfunc version() { return }")
	if got := examCheck(dir, version, run); got.status != examComplete {
		t.Fatalf("version implementation: %+v", got)
	}
	if got := examCheck(dir, examExercise{Test: "TestExercise1", File: "exercises_test.go", Functions: []string{"TestExercise1"}}, run); got.status != examComplete {
		t.Fatalf("test-writing exercise: %+v", got)
	}
}

func TestMissingOrEmptyVerificationCannotPass(t *testing.T) {
	for _, body := range []string{"", "func TestExercise1(t *testing.T) { /* no assertions */ }", "not valid Go"} {
		t.Run(body, func(t *testing.T) {
			dir := t.TempDir()
			writeFixture(t, dir, "exercises.go", "package main\nfunc answer() int { return 1 }")
			writeFixture(t, dir, "exercises_test.go", "package main\n"+body)
			got := examCheck(dir, examExercise{Test: "TestExercise1", Functions: []string{"answer"}}, func(_, _ string) error {
				t.Fatal("invalid test should not run")
				return nil
			})
			if got.status != examIncorrect {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestActualGoTestResultsAndIsolation(t *testing.T) {
	// 真实启动 go test，验证 panic/超时只终止当前题，后续题仍可执行。
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	writeFixture(t, dir, "go.mod", "module example.com/examfixture\n\ngo 1.24.0\n")
	writeFixture(t, dir, "exercises.go", "package main\nfunc answer() int { return 42 }")
	writeFixture(t, dir, "exercises_test.go", `package main
import "testing"
func TestPass(t *testing.T) { if answer() != 42 { t.Fatal("wrong answer") } }
func TestFail(t *testing.T) { t.Fatal("expected failure detail") }
func TestPanic(t *testing.T) { panic("unexpected panic") }
func TestSkip(t *testing.T) { t.Skip("examUnfinished test") }
func TestSubSkip(t *testing.T) { t.Run("case", func(t *testing.T) { t.Skip("examUnfinished subtest") }) }
func TestHang(t *testing.T) { select {} }
`)
	for _, test := range []struct{ name, detail string }{
		{"TestFail", "expected failure detail"}, {"TestPanic", "unexpected panic"},
		{"TestSkip", "不能判定完成"}, {"TestSubSkip", "不能判定完成"},
		{"TestMissing", "不能判定完成"}, {"TestHang", "timed out"}, {"TestPass", ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := examRunTestWithTimeout(dir, test.name, 500*time.Millisecond, 30*time.Second)
			if test.detail == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.detail) {
				t.Fatalf("error = %v; want %q", err, test.detail)
			}
		})
	}
	writeFixture(t, dir, "broken_test.go", "package main\nvar invalid int = \"wrong type\"")
	if err := examRunTest(dir, "TestPass"); err == nil {
		t.Fatal("test build error must not pass")
	}
	if err := os.Remove(filepath.Join(dir, "broken_test.go")); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, dir, "main_test.go", "package main\nimport (\"os\"; \"testing\")\nfunc TestMain(m *testing.M) { m.Run(); os.Exit(1) }")
	if err := examRunTest(dir, "TestPass"); err == nil {
		t.Fatal("package failure after a passing test must not pass")
	}
}

func TestEventOutputLimitStillReadsFinalResult(t *testing.T) {
	output := &examTestEvents{name: "TestExercise1"}
	for _, event := range []map[string]string{
		{"Action": "run", "Test": "TestExercise1"},
		{"Action": "output", "Test": "TestExercise1", "Output": strings.Repeat("x", 100000)},
		{"Action": "pass", "Test": "TestExercise1"},
	} {
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, '\n')
		// 模拟子进程分批输出 JSON，不能假设一次 Write 就是一整行。
		for len(data) > 0 {
			n := min(len(data), 53)
			output.Write(data[:n])
			data = data[n:]
		}
	}
	if !output.ran || !output.passed || output.diagnostics.Len() > 32*1024 {
		t.Fatalf("final examResult lost or output unbounded: ran=%v passed=%v bytes=%d", output.ran, output.passed, output.diagnostics.Len())
	}
}

// 验证章节复制到课程之外后，exam 命令仍能报告三种状态。
func TestIndependentExamCommand(t *testing.T) {
	t.Setenv("GOWORK", "off")
	dir := t.TempDir()
	for _, name := range []string{"go.mod", "lesson.go", "exercises.go", "exercises_test.go"} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		writeFixture(t, dir, name, string(data))
	}
	path := filepath.Join(dir, "exercises.go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, path, data, 0)
	if err != nil {
		t.Fatal(err)
	}
	functions := examFunctionsIn(file)
	// 从后往前替换函数体，避免前面的修改改变后续源码位置。
	for _, answer := range []struct{ name, body string }{
		{"exercise3", "/* 留空 */"},
		{"exercise2", `fmt.Println("wrong answer")`},
		{"exercise1", "fmt.Println(" + strconv.Quote("profile-api development true 9000 8080\nstring string bool int int\n\"\" false") + ")"},
	} {
		fn := functions[answer.name]
		start := positions.Position(fn.Body.Lbrace).Offset + 1
		end := positions.Position(fn.Body.Rbrace).Offset
		data = append(append(append([]byte{}, data[:start]...), []byte(answer.body)...), data[end:]...)
	}
	writeFixture(t, dir, "exercises.go", string(data))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "run", ".", "exam")
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("独立运行失败：%v\n%s", err, output)
	}
	for _, want := range []string{"练习 1（exercise1）：完成", "练习 2（exercise2）：错误", "练习 3（exercise3）：还未完成", "汇总：完成 1 题，错误 1 题，还未完成 1 题"} {
		if !strings.Contains(string(output), want) {
			t.Errorf("缺少 %q：\n%s", want, output)
		}
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
)

func exercise1() {
	const defaultPort = 8080
	var serviceName = "profile-api"
	environment := "development"
	debug := true
	port := defaultPort
	port = 9000
	fmt.Println(serviceName, environment, debug, port, defaultPort)
	fmt.Printf("%T %T %T %T %T\n", serviceName, environment, debug, port, defaultPort)
	serviceName = ""
	debug = false
	fmt.Printf("%q %t\n", serviceName, debug)
}
func exercise2() {
	var (
		maxRetries int
		region     string
		enabled    bool
	)
	fmt.Printf("%d %q %t\n", maxRetries, region, enabled)
	maxRetries, region, enabled = 3, "ap-shanghai", true
	fmt.Printf("%d %q %t\n", maxRetries, region, enabled)
	maxRetries = 0
	fmt.Println(maxRetries)
	// enabled 的 false 不能区分未填写与主动关闭，需要额外记录是否存在。
}
func exercise3() {
	environment := "production"
	{
		environment := "preview"
		fmt.Println(environment)
	}
	fmt.Println(environment)
	{
		environment = "staging"
	}
	fmt.Println(environment)
	// := 在内层声明新变量；= 修改外层变量。空字符串也遵守相同规则。
}

*/
