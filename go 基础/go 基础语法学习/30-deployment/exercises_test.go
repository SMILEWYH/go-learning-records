// 30-deployment：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 练习 2 需要按题目修改 lesson.go 的 version 输出，测试会实际构建并运行临时产物。练习 3 的 TestExercise3 本身就是答案。
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
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
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
		{Test: "TestExercise1", Functions: []string{"parsePort"}},
		{Test: "TestExercise2", Functions: []string{"runDeployment"}, File: "lesson.go", Variables: []string{"buildTime"}},
		// 本题实现位于测试文件；现有测试已提供可运行实现，按实际测试结果判定。
		{Test: "TestExercise3", Functions: []string{"TestExercise3"}, File: "exercises_test.go"},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	for _, c := range []struct {
		raw   string
		want  int
		valid bool
	}{{"", 8080, true}, {"9090", 9090, true}, {"1", 1, true}, {"65535", 65535, true}, {"0", 0, false}, {"65536", 0, false}, {"abc", 0, false}, {"-1", 0, false}} {
		got, err := parsePort(c.raw)
		if c.valid {
			exerciseOK(t, err)
			exerciseEqual(t, got, c.want)
		} else {
			exerciseError(t, err)
		}
	}
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	// 默认构建与注入构建均实际运行 version 模式；所有产物只写入临时目录。
	for _, inject := range []bool{false, true} {
		name := "service"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		binary := filepath.Join(t.TempDir(), name)
		args := []string{"build", "-o", binary}
		if inject {
			args = append(args, "-ldflags", "-X main.version=v0.2.0 -X main.buildTime=2026-09-27T12:00:00Z")
		}
		args = append(args, ".")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		cmd := exec.CommandContext(ctx, "go", args...)
		output, err := cmd.CombinedOutput()
		cancel()
		if err != nil {
			t.Fatalf("构建失败：%v\n%s", err, output)
		}
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		output, err = exec.CommandContext(ctx, binary, "version").CombinedOutput()
		cancel()
		if err != nil {
			t.Fatalf("运行失败：%v\n%s", err, output)
		}
		text := string(output)
		if inject {
			if !strings.Contains(text, "version=v0.2.0") || !strings.Contains(text, "buildTime=2026-09-27T12:00:00Z") {
				t.Fatalf("未正确输出注入的版本与时间：%s", text)
			}
		} else {
			if !strings.Contains(text, "version=dev") || !strings.Contains(text, "buildTime=") {
				t.Fatalf("默认输出应包含 version=dev 和 buildTime：%s", text)
			}
			for _, field := range strings.Fields(text) {
				if strings.HasPrefix(field, "buildTime=") && strings.TrimPrefix(field, "buildTime=") == "" {
					t.Fatal("默认构建时间不能为空")
				}
			}
		}
	}
	// 交叉构建按照 lesson.go 末尾的部署附录 手动验证，不在测试中执行其他平台产物。
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	// 本题答案：验证 HTTP 标准库 Shutdown 的等待行为，使用通道确认请求已经开始。
	for _, expire := range []bool{false, true} {
		t.Run(map[bool]string{false: "等待请求完成", true: "超时返回错误"}[expire], func(t *testing.T) {
			defer exerciseGuard(t)
			entered := make(chan struct{})
			release := make(chan struct{})
			handlerDone := make(chan struct{})
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				close(entered)
				defer close(handlerDone)
				<-release
				_, _ = io.WriteString(w, "done")
			}))
			server.Start()
			// 无论断言在哪一步失败，都先释放 handler 再关闭服务，避免清理阶段卡住。
			released := false
			defer func() {
				if !released {
					close(release)
				}
				server.Close()
			}()
			client := server.Client()
			client.Timeout = 3 * time.Second
			requestDone := make(chan error, 1)
			go func() {
				response, err := client.Get(server.URL)
				if err == nil {
					data, readErr := io.ReadAll(response.Body)
					response.Body.Close()
					err = readErr
					if err == nil && string(data) != "done" {
						err = errors.New("响应不完整")
					}
				}
				requestDone <- err
			}()
			exerciseWait(t, entered)
			duration := 2 * time.Second
			if expire {
				duration = 20 * time.Millisecond
			}
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			shutdownDone := make(chan error, 1)
			started := make(chan struct{})
			server.Config.RegisterOnShutdown(func() { close(started) })
			go func() { shutdownDone <- server.Config.Shutdown(ctx) }()
			exerciseWait(t, started)
			if expire {
				select {
				case err := <-shutdownDone:
					if !errors.Is(err, context.DeadlineExceeded) {
						t.Fatalf("期望超时：%v", err)
					}
				case <-time.After(time.Second):
					t.Fatal("Shutdown 未按时返回")
				}
				close(release)
				released = true
				exerciseWait(t, handlerDone)
				server.CloseClientConnections()
				select {
				case <-requestDone:
				case <-time.After(time.Second):
					t.Fatal("客户端未收尾")
				}
				return
			}
			select {
			case err := <-shutdownDone:
				t.Fatalf("请求尚未释放却已退出：%v", err)
			default:
			}
			close(release)
			released = true
			select {
			case err := <-requestDone:
				exerciseOK(t, err)
			case <-time.After(3 * time.Second):
				t.Fatal("请求未结束")
			}
			select {
			case err := <-shutdownDone:
				exerciseOK(t, err)
			case <-time.After(3 * time.Second):
				t.Fatal("Shutdown 未结束")
			}
			exerciseWait(t, handlerDone)
		})
	}
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

func exerciseWait(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("等待超时，请检查是否缺少发送、关闭或取消")
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
	"strconv"
)

func parsePort(raw string) (int, error) {
	if raw == "" {
		return 8080, nil
	}
	port, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("端口解析失败: %w", err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("端口超出范围")
	}
	return port, nil
}

// 练习 2：在 lesson.go 的 version 声明旁添加：
// var buildTime = "unknown"
// 将 runDeployment 的 case "version" 输出替换为：
// fmt.Printf("version=%s buildTime=%s go=%s target=%s/%s\n", version, buildTime, runtime.Version(), runtime.GOOS, runtime.GOARCH)
// 本机构建：go build -ldflags "-X main.version=v0.2.0 -X main.buildTime=2026-09-27T12:00:00Z" -o bin/service .
// Linux 构建：CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/service-linux .
// Linux 产物使用不同的可执行文件格式和系统调用 ABI，不能在 macOS 上原生执行。
// version/buildTime 是编译时注入；PORT 通常由进程运行时读取，可以不重新构建就更改。
// 练习 3 的可运行答案就是 TestExercise3，使用 channel 确定请求已经进入 handler。

*/
