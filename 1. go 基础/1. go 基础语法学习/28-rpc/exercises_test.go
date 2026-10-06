// 28-rpc：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 练习 2 的 TestExercise2 本身就是参考实现；练习 3 补充 batchQuotes(client, quantities) 入口，调用者负责连接 deadline 和关闭。
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
	"net"
	"net/rpc"
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
		{Test: "TestExercise1", Functions: []string{"ShippingService.Calculate"}},
		// 本题实现位于测试文件；现有测试已提供可运行实现，按实际测试结果判定。
		{Test: "TestExercise2", Functions: []string{"TestExercise2"}, File: "exercises_test.go"},
		{Test: "TestExercise3", Functions: []string{"batchQuotes"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	var probe ShippingResponse
	exerciseError(t, (&ShippingService{}).Calculate(ShippingRequest{SubtotalCents: -1}, &probe))
	server := rpc.NewServer()
	exerciseOK(t, server.Register(&ShippingService{}))
	client := exerciseRPCClient(t, server)
	for _, c := range []struct {
		amount  int64
		express bool
		fee     int64
	}{{9999, false, 600}, {10000, false, 0}, {10000, true, 1000}, {0, true, 1600}} {
		var reply ShippingResponse
		exerciseOK(t, client.Call("ShippingService.Calculate", ShippingRequest{c.amount, c.express}, &reply))
		exerciseEqual(t, reply.FeeCents, c.fee)
	}
	var reply ShippingResponse
	exerciseError(t, client.Call("ShippingService.Calculate", ShippingRequest{-1, false}, &reply))
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	// 本题要求核对通信错误；这段测试同时给出可运行的参考写法。
	server, err := newRPCServer()
	exerciseOK(t, err)
	client := exerciseRPCClient(t, server)
	var success QuoteResponse
	exerciseOK(t, client.Call("PriceService.Quote", QuoteRequest{Quantity: 2, UnitPriceCents: 500}, &success))
	exerciseEqual(t, success.TotalCents, int64(1000))
	var rejected QuoteResponse
	err = client.Call("PriceService.Quote", QuoteRequest{Quantity: 0, UnitPriceCents: 500}, &rejected)
	var remote rpc.ServerError
	if !errors.As(err, &remote) {
		t.Fatalf("期望远程业务错误：%v", err)
	}
	var missing QuoteResponse
	exerciseError(t, client.Call("PriceService.Missing", QuoteRequest{}, &missing))
	exerciseOK(t, client.Close())
	var closed QuoteResponse
	exerciseError(t, client.Call("PriceService.Quote", QuoteRequest{Quantity: 1, UnitPriceCents: 500}, &closed))
	exerciseEqual(t, closed.TotalCents, int64(0))
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	server, err := newRPCServer()
	exerciseOK(t, err)
	client := exerciseRPCClient(t, server)
	for _, quantities := range [][]int{{1, 2, 3}, {1, 0, 3}, nil} {
		replies, errs := batchQuotes(client, quantities)
		exerciseEqual(t, len(replies), len(quantities))
		exerciseEqual(t, len(errs), len(quantities))
		for i, q := range quantities {
			if q == 0 {
				exerciseError(t, errs[i])
			} else {
				exerciseOK(t, errs[i])
				exerciseEqual(t, replies[i].TotalCents, int64(q*500))
			}
		}
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

// 使用真实本机 TCP 连接，设置 deadline，并在清理时等待服务退出。
func exerciseRPCClient(t *testing.T, server *rpc.Server) *rpc.Client {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	exerciseOK(t, err)
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
		server.ServeConn(conn)
	}()
	conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	if err != nil {
		listener.Close()
		exerciseWait(t, done)
		t.Fatal(err)
	}
	exerciseOK(t, conn.SetDeadline(time.Now().Add(3*time.Second)))
	client := rpc.NewClient(conn)
	t.Cleanup(func() { client.Close(); listener.Close(); exerciseWait(t, done) })
	return client
}

// ── 讲解配套测试：验证 lesson.go 中的示例 ──
func TestQuoteLimits(t *testing.T) {
	service := &PriceService{}
	for _, test := range []struct {
		quantity int
		price    int64
		valid    bool
		total    int64
	}{
		{3, 1990, true, 5970}, {1, 0, true, 0}, {1000, 100_000_000, true, 100_000_000_000},
		{0, 500, false, 0}, {1001, 500, false, 0}, {1, -1, false, 0}, {1, 100_000_001, false, 0},
	} {
		var reply QuoteResponse
		err := service.Quote(QuoteRequest{Quantity: test.quantity, UnitPriceCents: test.price}, &reply)
		if (err == nil) != test.valid || reply.TotalCents != test.total {
			t.Errorf("quantity=%d price=%d reply=%+v err=%v", test.quantity, test.price, reply, err)
		}
	}
}

func TestRPCExchange(t *testing.T) {
	// 验证 gob 序列化、注册名、同步/异步回复，以及业务错误穿过真实 TCP 连接。
	if err := demoRPC(); err != nil {
		t.Fatal(err)
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
	"net/rpc"
)

type ShippingRequest struct {
	SubtotalCents int64
	Express       bool
}
type ShippingResponse struct{ FeeCents int64 }
type ShippingService struct{}

func (*ShippingService) Calculate(args ShippingRequest, reply *ShippingResponse) error {
	if args.SubtotalCents < 0 {
		return fmt.Errorf("金额不能为负数")
	}
	fee := int64(600)
	if args.SubtotalCents >= 10000 {
		fee = 0
	}
	if args.Express {
		fee += 1000
	}
	reply.FeeCents = fee
	return nil
}

// 练习 2 的答案为 TestExercise2。net/rpc 只传递错误文本，客户端得到 rpc.ServerError，原始具体 error 类型无法保留。
func batchQuotes(client *rpc.Client, quantities []int) ([]QuoteResponse, []error) {
	replies := make([]QuoteResponse, len(quantities))
	calls := make([]*rpc.Call, len(quantities))
	errs := make([]error, len(quantities))
	for i, quantity := range quantities {
		calls[i] = client.Go("PriceService.Quote", QuoteRequest{Quantity: quantity, UnitPriceCents: 500}, &replies[i], make(chan *rpc.Call, 1))
	}
	for i, call := range calls {
		finished := <-call.Done
		errs[i] = finished.Error
	}
	return replies, errs
}

*/
