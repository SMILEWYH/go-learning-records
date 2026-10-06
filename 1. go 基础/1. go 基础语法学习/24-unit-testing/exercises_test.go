// 24-unit-testing：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 这三个 TestExercise 测试分别演示表驱动、临时目录和并发断言，可继续增加自己的用例；沿用已有 exercises_test.go 的 TestMain。
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
	"sync"
	"testing"
	"time"
)

// TestExam 是 exam 模式入口。普通 go test 仍直接运行下面的测试函数。
func TestExam(t *testing.T) {
	if os.Getenv("GO_LESSON_EXAM") != "1" {
		t.Skip("通过 go run . exam 查看逐题检查结果")
	}
	examRun([]examExercise{
		{Test: "TestExercise1", Functions: []string{"validatePagination"}},
		{Test: "TestExercise2", Functions: []string{"saveFeatureFlag"}},
		{Test: "TestExercise3", Functions: []string{"RequestCounter.Inc", "RequestCounter.Value"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	for _, c := range []struct {
		page, size int
		valid      bool
	}{{1, 20, true}, {3, 100, true}, {1, 1, true}, {0, 20, false}, {1, 0, false}, {1, 101, false}, {-1, 10, false}} {
		t.Run(fmt.Sprintf("page=%d_size=%d", c.page, c.size), func(t *testing.T) {
			defer exerciseGuard(t)
			err := validatePagination(c.page, c.size)
			if c.valid {
				exerciseOK(t, err)
			} else {
				exerciseError(t, err)
			}
		})
	}
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	t.Run("覆盖", func(t *testing.T) {
		defer exerciseGuard(t)
		path := filepath.Join(t.TempDir(), "flag")
		exerciseOK(t, os.WriteFile(path, []byte("old-long-content"), 0600))
		exerciseOK(t, saveFeatureFlag(path, false))
		exerciseFile(t, path, "enabled=false\n")
		exerciseOK(t, saveFeatureFlag(path, true))
		exerciseFile(t, path, "enabled=true\n")
	})
	t.Run("父目录缺失", func(t *testing.T) {
		defer exerciseGuard(t)
		exerciseError(t, saveFeatureFlag(filepath.Join(t.TempDir(), "missing", "flag"), true))
	})
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	var counter RequestCounter
	exerciseEqual(t, counter.Value(), 0)
	var wait sync.WaitGroup
	wait.Add(40)
	for i := 0; i < 20; i++ {
		go func() {
			defer wait.Done()
			for j := 0; j < 100; j++ {
				counter.Inc()
			}
		}()
		go func() {
			defer wait.Done()
			for j := 0; j < 100; j++ {
				_ = counter.Value()
			}
		}()
	}
	wait.Wait()
	exerciseEqual(t, counter.Value(), 2000)
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

func exerciseFile(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	exerciseOK(t, err)
	exerciseEqual(t, string(data), want)
}

// ── 讲解配套测试：验证 lesson.go 中的示例 ──
// TestMain 包装当前包的整次测试运行，用于确实需要的包级初始化/清理。
// 只需要测试独立资源时，优先用 t.TempDir/t.Cleanup；它们更容易隔离测试。
// os.Exit 不执行 defer，所以必须在 os.Exit 之前完成需要的清理。
func TestMain(m *testing.M) {
	fmt.Println("测试包开始")
	code := m.Run()
	fmt.Println("测试包结束")
	os.Exit(code)
}

func TestShippingFee(t *testing.T) {
	t.Parallel() // 纯函数用例不共享可变状态，可以与其他独立测试并行。
	// 表驱动测试：测试数据说明业务输入/预期，一个循环复用断言逻辑。
	tests := []struct {
		name     string
		subtotal int
		express  bool
		want     int
		wantErr  bool
	}{
		{name: "standard", subtotal: 5900, want: 600},
		{name: "before_free_threshold", subtotal: 9999, want: 600},
		{name: "free_standard", subtotal: 10000, want: 0},
		{name: "express_below_threshold", subtotal: 5900, express: true, want: 1600},
		{name: "express_after_threshold", subtotal: 10000, express: true, want: 1000},
		{name: "zero_subtotal", subtotal: 0, want: 600},
		{name: "negative_subtotal", subtotal: -1, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel() // Go 1.22 起，:= 声明的 range 变量每轮独立，当前项目无需 test := test。
			got, err := ShippingFee(test.subtotal, test.express)
			if (err != nil) != test.wantErr {
				t.Fatalf("ShippingFee() error = %v, wantErr = %t", err, test.wantErr)
			}
			if !test.wantErr && got != test.want {
				t.Errorf("ShippingFee() = %d, want %d", got, test.want)
			}
		})
	}
}

// Helper 将失败报告定位到调用方；Cleanup 在该测试及其子测试结束后执行，按后进先出顺序运行。
// 这与 defer 在当前函数返回时就执行不同，尤其适合需要供子测试共享的资源。
func openTestFile(t *testing.T, path string) *os.File {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	t.Cleanup(func() {
		if err := file.Close(); err != nil {
			t.Errorf("close fixture: %v", err)
		}
	})
	return file
}

func TestSaveReceipt(t *testing.T) {
	t.Parallel()
	// TempDir 自动清理，不依赖工作目录、仓库中的文件或特定机器的路径。
	path := filepath.Join(t.TempDir(), "receipt.txt")
	if err := SaveReceipt(path, "order-123"); err != nil {
		t.Fatal(err)
	}
	file := openTestFile(t, path)
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "order-123\n" {
		t.Errorf("content = %q, want %q", content, "order-123\n")
	}

	// 校验失败必须保留原收据；同时验证错误与文件副作用。
	if err := SaveReceipt(path, ""); !errors.Is(err, ErrEmptyReceipt) {
		t.Fatalf("empty receipt: error = %v, want ErrEmptyReceipt", err)
	}
	content, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "order-123\n" {
		t.Errorf("invalid save changed content to %q", content)
	}
}

func TestSaveReceiptMissingDirectory(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "missing", "receipt.txt")
	err := SaveReceipt(path, "order-123")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error = %v, want os.ErrNotExist", err)
	}
	var pathError *os.PathError
	if !errors.As(err, &pathError) {
		t.Fatalf("error type = %T, want *os.PathError", err)
	}
	// 使用 Is/As 断言稳定的错误契约，不依赖操作系统错误字符串的语言或措辞。
}

// 并行测试须各自拥有临时目录和数据；避免修改全局变量、工作目录或共享文件。
// t.Setenv/t.Chdir 会改变进程状态，不能在并行测试或有并行祖先的测试中使用。
// 工作 goroutine 中不要调用 t.Fatal/t.FailNow；先把结果交回测试 goroutine 再断言。

func BenchmarkShippingFee(b *testing.B) {
	for b.Loop() { // 当前版本推荐的循环形式，自动管理计时并防止整段调用被优化掉。
		if _, err := ShippingFee(12000, true); err != nil {
			b.Fatal(err)
		}
	}
	// -benchmem 输出分配信息；测量值受环境、数据和编译器影响，不能凭一次输出下性能结论。
}

func FuzzShippingFee(f *testing.F) {
	for _, subtotal := range []int{-1, 0, 9999, 10000, 12000} {
		f.Add(subtotal)
	}
	f.Fuzz(func(t *testing.T, subtotal int) {
		standard, standardErr := ShippingFee(subtotal, false)
		express, expressErr := ShippingFee(subtotal, true)
		if subtotal < 0 {
			if standardErr == nil || expressErr == nil {
				t.Fatal("negative subtotal must fail for both shipping modes")
			}
			return
		}
		if standardErr != nil || expressErr != nil {
			t.Fatalf("nonnegative subtotal failed: standard=%v express=%v", standardErr, expressErr)
		}
		// 检查不同输入下必须成立的业务性质，而不是只验证“没有 panic”。
		if standard < 0 || express < 0 || express-standard != 1000 {
			t.Fatalf("invalid fees: standard=%d express=%d", standard, express)
		}
		if subtotal >= 10000 && standard != 0 {
			t.Fatalf("eligible order must have free standard shipping: %d", standard)
		}
	})
	// 普通 go test 只检查上述种子；只有显式 -fuzz 才生成更多输入，配合 -fuzztime 限时运行。
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
	"os"
	"sync"
)

func validatePagination(page, size int) error {
	if page < 1 || size < 1 || size > 100 {
		return fmt.Errorf("分页参数非法")
	}
	return nil
}
func saveFeatureFlag(path string, enabled bool) error {
	return os.WriteFile(path, []byte(fmt.Sprintf("enabled=%t\n", enabled)), 0600)
}

type RequestCounter struct {
	mu    sync.Mutex
	value int
}

func (c *RequestCounter) Inc()       { c.mu.Lock(); defer c.mu.Unlock(); c.value++ }
func (c *RequestCounter) Value() int { c.mu.Lock(); defer c.mu.Unlock(); return c.value }

*/
