// 29-websocket：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 练习 1 的测试入口为 newJSONSocketMux；练习 3 为 sendProducerMessages，groups 每项对应一个生产者。练习 2 的握手测试本身就是答案。
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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
	"slices"
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
		{Test: "TestExercise1", Functions: []string{"newJSONSocketMux"}},
		// 本题实现位于测试文件；现有测试已提供可运行实现，按实际测试结果判定。
		{Test: "TestExercise2", Functions: []string{"TestExercise2"}, File: "exercises_test.go"},
		{Test: "TestExercise3", Functions: []string{"sendProducerMessages"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	server := httptest.NewServer(newJSONSocketMux())
	defer server.Close()
	dialer := websocket.Dialer{HandshakeTimeout: time.Second}
	conn, response, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/ws", nil)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetReadLimit(1024)
	exerciseOK(t, conn.SetReadDeadline(time.Now().Add(2*time.Second)))
	exerciseOK(t, conn.SetWriteDeadline(time.Now().Add(2*time.Second)))
	for _, c := range []struct {
		message Message
		kind    string
	}{{Message{Type: "chat", Text: "hi"}, "ack"}, {Message{Type: "unknown", Text: "hi"}, "error"}, {Message{Type: "chat", Text: ""}, "error"}} {
		exerciseOK(t, conn.WriteJSON(c.message))
		var reply Message
		exerciseOK(t, conn.ReadJSON(&reply))
		exerciseEqual(t, reply.Type, c.kind)
		if c.kind == "ack" {
			exerciseEqual(t, reply.Text, "hi")
		}
	}
	exerciseOK(t, conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second)))
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	// 本题要求写握手测试，以下即答案；被测服务复用 lesson.go 的默认 Origin 策略。
	server := httptest.NewServer(newWebSocketMux())
	defer server.Close()
	for _, origin := range []string{server.URL, "https://other.example", ""} {
		headers := http.Header{}
		if origin != "" {
			headers.Set("Origin", origin)
		}
		dialer := websocket.Dialer{HandshakeTimeout: time.Second}
		conn, response, err := dialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/ws", headers)
		if origin == "https://other.example" {
			exerciseError(t, err)
			if response == nil {
				t.Fatal("握手失败应有 HTTP 响应")
			}
			exerciseEqual(t, response.StatusCode, 403)
			response.Body.Close()
			continue
		}
		exerciseOK(t, err)
		exerciseEqual(t, response.StatusCode, 101)
		exerciseOK(t, conn.SetReadDeadline(time.Now().Add(time.Second)))
		kind, data, err := conn.ReadMessage()
		exerciseOK(t, err)
		exerciseEqual(t, kind, websocket.TextMessage)
		exerciseEqual(t, string(data), "welcome")
		exerciseOK(t, conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second)))
		conn.Close()
	}
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	server := httptest.NewServer(newWebSocketMux())
	defer server.Close()
	address := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	for _, groups := range [][][]string{{{"a-1", "a-2"}, {"b-1", "b-2"}}, nil} {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		got, err := sendProducerMessages(ctx, address, groups)
		cancel()
		exerciseOK(t, err)
		want := []string{"welcome"}
		for _, messages := range groups {
			for _, message := range messages {
				want = append(want, "收到："+message)
			}
		}
		slices.Sort(got)
		slices.Sort(want)
		exerciseEqual(t, got, want)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := sendProducerMessages(ctx, address, [][]string{{"x"}})
	exerciseError(t, err)
	// 对端立即拒绝握手时也必须返回，不能残留生产者或 writer。
	reject := httptest.NewServer(http.NotFoundHandler())
	defer reject.Close()
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = sendProducerMessages(ctx, "ws"+strings.TrimPrefix(reject.URL, "http")+"/ws", [][]string{{"x", "y"}})
	exerciseError(t, err)
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
func TestWebSocketExchange(t *testing.T) {
	// 验证主动 welcome、两条文本消息，以及 1000 正常关闭握手。
	if err := demoWebSocket(); err != nil {
		t.Fatal(err)
	}
}

func TestWebSocketOrigin(t *testing.T) {
	server := httptest.NewServer(newWebSocketMux())
	defer server.Close()
	address := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	for _, test := range []struct {
		name, origin string
		allowed      bool
	}{
		{"same origin", server.URL, true},
		{"different origin", "https://other.example", false},
		{"no origin", "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			header := make(http.Header)
			if test.origin != "" {
				header.Set("Origin", test.origin)
			}
			dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
			conn, response, err := dialer.Dial(address, header)
			if !test.allowed {
				if conn != nil {
					conn.Close()
				}
				if response != nil {
					defer response.Body.Close()
				}
				if err == nil || response == nil || response.StatusCode != http.StatusForbidden {
					t.Fatalf("异源必须拒绝：response=%v err=%v", response, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if response.StatusCode != http.StatusSwitchingProtocols {
				t.Fatalf("握手状态=%d", response.StatusCode)
			}
			readWelcome(t, conn)
			if err := conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second)); err != nil {
				t.Fatal(err)
			}
			_, _, err = conn.ReadMessage()
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				t.Fatalf("关闭失败：%v", err)
			}
		})
	}
}

func readWelcome(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatal(err)
	}
	kind, message, err := conn.ReadMessage()
	if err != nil || kind != websocket.TextMessage || string(message) != "welcome" {
		t.Fatalf("welcome: kind=%d message=%q err=%v", kind, message, err)
	}
}

func TestWebSocketRejectsUnsupportedMessages(t *testing.T) {
	server := httptest.NewServer(newWebSocketMux())
	defer server.Close()
	address := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	for _, test := range []struct {
		name      string
		kind      int
		body      string
		closeCode int
	}{
		{"binary", websocket.BinaryMessage, "binary", websocket.CloseUnsupportedData},
		{"oversize", websocket.TextMessage, strings.Repeat("a", 1025), websocket.CloseMessageTooBig},
	} {
		t.Run(test.name, func(t *testing.T) {
			dialer := websocket.Dialer{HandshakeTimeout: 3 * time.Second}
			conn, _, err := dialer.Dial(address, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			readWelcome(t, conn)
			if err := conn.WriteMessage(test.kind, []byte(test.body)); err != nil {
				t.Fatal(err)
			}
			_, _, err = conn.ReadMessage()
			if !websocket.IsCloseError(err, test.closeCode) {
				t.Fatalf("期望关闭码 %d，实际 %v", test.closeCode, err)
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
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func newJSONSocketMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{HandshakeTimeout: time.Second}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		conn.SetReadLimit(1024)
		for {
			if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
				return
			}
			var message Message
			if err := conn.ReadJSON(&message); err != nil {
				return
			}
			reply := Message{Type: "ack", Text: message.Text}
			if message.Type != "chat" || message.Text == "" {
				reply = Message{Type: "error", Text: "无效消息"}
			}
			if err := conn.SetWriteDeadline(time.Now().Add(2 * time.Second)); err != nil {
				return
			}
			if err := conn.WriteJSON(reply); err != nil {
				return
			}
		}
	})
	return mux
}

// 练习 2 的答案是 TestExercise2；Origin 校验不能替代鉴权，因为非浏览器客户端可以省略或伪造 Origin。
// groups 每项是一位生产者的消息。返回 welcome 与回显；不要求生产者之间的顺序。
func sendProducerMessages(parent context.Context, address string, groups [][]string) ([]string, error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	dialer := websocket.Dialer{HandshakeTimeout: time.Second}
	conn, response, err := dialer.DialContext(ctx, address, nil)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return nil, err
	}
	defer conn.Close()
	conn.SetReadLimit(1024)
	// 取消时主动关闭连接，唤醒阻塞 Read。单独等待监视者结束，防止任务残留。
	watcherDone := make(chan struct{})
	defer func() { cancel(); <-watcherDone }()
	go func() { defer close(watcherDone); <-ctx.Done(); conn.Close() }()
	queue := make(chan string, 4)
	var producers sync.WaitGroup
	producers.Add(len(groups))
	for _, messages := range groups {
		go func(items []string) {
			defer producers.Done()
			for _, message := range items {
				select {
				case queue <- message:
				case <-ctx.Done():
					return
				}
			}
		}(messages)
	}
	allProduced := make(chan struct{})
	go func() { producers.Wait(); close(queue); close(allProduced) }()
	// writer 等 reader 收齐回显后才发关闭帧；避免服务器回复关闭帧时丢失仍在途的回显。
	repliesRead := make(chan struct{})
	writerDone := make(chan error, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				writerDone <- ctx.Err()
				return
			case message, ok := <-queue:
				if !ok {
					select {
					case <-ctx.Done():
						writerDone <- ctx.Err()
						return
					case <-repliesRead:
					}
					writerDone <- conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"), time.Now().Add(time.Second))
					return
				}
				if err := conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
					cancel()
					writerDone <- err
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					cancel()
					writerDone <- err
					return
				}
			}
		}
	}()
	count := 1
	for _, messages := range groups {
		count += len(messages)
	}
	result := make([]string, 0, count)
	var readErr error
	for len(result) < count {
		if readErr = conn.SetReadDeadline(time.Now().Add(2 * time.Second)); readErr != nil {
			break
		}
		kind, data, err := conn.ReadMessage()
		if err != nil {
			readErr = err
			break
		}
		if kind != websocket.TextMessage {
			readErr = fmt.Errorf("响应不是文本")
			break
		}
		result = append(result, string(data))
	}
	if readErr != nil {
		cancel()
	} else {
		close(repliesRead)
	}
	writeErr := <-writerDone
	cancel()
	<-allProduced
	if readErr != nil {
		return nil, readErr
	}
	if writeErr != nil {
		return nil, writeErr
	}
	return result, nil
}

*/
