// 27-http：本章检查与参考答案，所有检查逻辑均在本文件内。
// exam 入口 -> 练习测试 -> 讲解测试（如有）-> 检查辅助代码 -> 文件末尾参考答案。
// 只检查一题：go test -run "^TestExercise1$" -v -count=1 -timeout 30s
//
// 本章练习测试与参考答案。只测一题：go test -run "^TestExercise1$" -v -timeout 30s
// 未完成的练习会失败；参考答案在文件末尾注释中，不参与运行。
// 练习 1 使用 productHandler，练习 2 使用 newUserPatchMux 作为测试入口。练习 3 的取消案例检查请求 context；fetchMessage 仍保留题目原签名。
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
		{Test: "TestExercise1", Functions: []string{"productHandler"}},
		{Test: "TestExercise2", Functions: []string{"newUserPatchMux"}},
		{Test: "TestExercise3", Functions: []string{"fetchMessage"}},
	})
}

// ── 本章练习的验收测试 ──
func TestExercise1(t *testing.T) {
	defer exerciseGuard(t)
	for _, c := range []struct {
		method, url        string
		status, page, size int
	}{{"GET", "/products", 200, 1, 20}, {"GET", "/products?page=2&size=10", 200, 2, 10}, {"GET", "/products?page=1&size=100", 200, 1, 100}, {"GET", "/products?page=abc", 400, 0, 0}, {"GET", "/products?page=0", 400, 0, 0}, {"GET", "/products?size=101", 400, 0, 0}, {"GET", "/products?size=0", 400, 0, 0}, {"POST", "/products", 405, 0, 0}} {
		rec := httptest.NewRecorder()
		productHandler(rec, httptest.NewRequest(c.method, c.url, nil))
		exerciseEqual(t, rec.Code, c.status)
		if c.status == 200 {
			if !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
				t.Fatal("缺少 JSON Content-Type")
			}
			var got map[string]int
			exerciseOK(t, json.Unmarshal(rec.Body.Bytes(), &got))
			exerciseEqual(t, got, map[string]int{"page": c.page, "size": c.size})
		}
	}
}

func TestExercise2(t *testing.T) {
	defer exerciseGuard(t)
	server := httptest.NewServer(newUserPatchMux())
	defer server.Close()
	client := server.Client()
	client.Timeout = time.Second
	for _, c := range []struct {
		path, body, media string
		status            int
	}{{"/users/7", `{"name":"小周"}`, "application/json", 200}, {"/users/abc", `{"name":"x"}`, "application/json", 400}, {"/users/0", `{"name":"x"}`, "application/json", 400}, {"/users/7", `{"name":""}`, "application/json", 400}, {"/users/7", `{"name":"x","extra":1}`, "application/json", 400}, {"/users/7", `{"name":"x"} {}`, "application/json", 400}, {"/users/7", `{"name":"x"}`, "text/plain", 415}, {"/users/7", `{"name":"` + strings.Repeat("x", 1100) + `"}`, "application/json", 413}, {"/users/7", `{"name":"x"}` + strings.Repeat(" ", 1024), "application/json", 413}, {"/users/7", `{`, "application/json", 400}} {
		request, err := http.NewRequest(http.MethodPatch, server.URL+c.path, strings.NewReader(c.body))
		exerciseOK(t, err)
		request.Header.Set("Content-Type", c.media)
		response, err := client.Do(request)
		exerciseOK(t, err)
		data, readErr := io.ReadAll(response.Body)
		closeErr := response.Body.Close()
		exerciseOK(t, readErr)
		exerciseOK(t, closeErr)
		exerciseEqual(t, response.StatusCode, c.status)
		if c.status == 200 {
			var got map[string]any
			exerciseOK(t, json.Unmarshal(data, &got))
			exerciseEqual(t, got, map[string]any{"id": float64(7), "name": "小周"})
		}
	}
}

func TestExercise3(t *testing.T) {
	defer exerciseGuard(t)
	for _, status := range []int{200, 404, 500} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status); fmt.Fprint(w, "ready") }))
		client := server.Client()
		client.Timeout = time.Second
		got, err := fetchMessage(client, server.URL)
		server.Close()
		if status == 200 {
			exerciseOK(t, err)
			exerciseEqual(t, got, "ready")
		} else {
			exerciseError(t, err)
			if !strings.Contains(err.Error(), fmt.Sprint(status)) {
				t.Fatal("错误应包含状态码")
			}
		}
	}
	for _, status := range []int{200, 500} {
		body := &exerciseTrackedBody{Reader: strings.NewReader("ready")}
		client := &http.Client{Transport: exerciseRoundTripper(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Body: body, Header: make(http.Header)}, nil
		})}
		_, _ = fetchMessage(client, "http://local.test")
		exerciseEqual(t, body.closed, true)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	exerciseOK(t, err)
	client := server.Client()
	client.Timeout = time.Second
	response, err := client.Do(request)
	if response != nil {
		response.Body.Close()
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("期望 context.Canceled：%v", err)
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

type exerciseRoundTripper func(*http.Request) (*http.Response, error)

func (f exerciseRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type exerciseTrackedBody struct {
	io.Reader
	closed bool
}

func (b *exerciseTrackedBody) Close() error { b.closed = true; return nil }

// ── 讲解配套测试：验证 lesson.go 中的示例 ──
func TestCreateUserValidation(t *testing.T) {
	mux := newHTTPMux()
	cases := []struct {
		name, body, contentType string
		status                  int
	}{
		{"valid", `{"name":" 小林 "}`, "application/json; charset=utf-8", 201},
		{"empty name", `{"name":" "}`, "application/json", 400},
		{"unknown field", `{"name":"x","admin":true}`, "application/json", 400},
		{"second object", `{"name":"x"}{"name":"y"}`, "application/json", 400},
		{"trailing garbage", `{"name":"x"}broken`, "application/json", 400},
		{"null", `null`, "application/json", 400},
		{"array", `[]`, "application/json", 400},
		{"empty body", "", "application/json", 400},
		{"wrong type", `{"name":123}`, "application/json", 400},
		{"wrong media", `{"name":"x"}`, "text/plain", 415},
		{"too large", `{"name":"` + strings.Repeat("a", 1024) + `"}`, "application/json", 413},
		{"large trailing space", `{"name":"x"}` + strings.Repeat(" ", 1024), "application/json", 413},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(test.body))
			request.Header.Set("Content-Type", test.contentType)
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status=%d want=%d body=%s", response.Code, test.status, response.Body)
			}
			if test.status == http.StatusCreated {
				var user UserResponse
				if err := json.Unmarshal(response.Body.Bytes(), &user); err != nil {
					t.Fatal(err)
				}
				if user.ID != 1 || user.Name != "小林" || !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
					t.Fatalf("错误的成功响应：%+v headers=%v", user, response.Header())
				}
			}
		})
	}
}

func TestHTTPRoutes(t *testing.T) {
	mux := newHTTPMux()
	for _, test := range []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/hello/Go?lang=zh", 200},
		{http.MethodDelete, "/users", 405},
		{http.MethodGet, "/missing", 404},
	} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(test.method, test.path, nil))
		if response.Code != test.status {
			t.Errorf("%s %s: status=%d", test.method, test.path, response.Code)
		}
		if test.status == 405 && response.Header().Get("Allow") != "POST" {
			t.Errorf("Allow=%q", response.Header().Get("Allow"))
		}
	}
}

func TestHTTPExchange(t *testing.T) {
	if err := demoHTTP(); err != nil {
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
)

func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	page, size := 1, 20
	var err error
	if raw, exists := r.URL.Query()["page"]; exists {
		page, err = strconv.Atoi(raw[0])
		if err != nil {
			http.Error(w, "page 非法", 400)
			return
		}
	}
	if raw, exists := r.URL.Query()["size"]; exists {
		size, err = strconv.Atoi(raw[0])
		if err != nil {
			http.Error(w, "size 非法", 400)
			return
		}
	}
	if page < 1 || size < 1 || size > 100 {
		http.Error(w, "分页非法", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"page": page, "size": size})
}
func newUserPatchMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.Error(w, "id 非法", 400)
			return
		}
		media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || media != "application/json" {
			http.Error(w, "需要 JSON", 415)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		defer r.Body.Close()
		// 先读完受限正文，保证包括尾部空白在内的超限正文统一返回 413。
		data, err := io.ReadAll(r.Body)
		if err != nil {
			var large *http.MaxBytesError
			if errors.As(err, &large) {
				http.Error(w, "正文超限", 413)
			} else {
				http.Error(w, "正文读取失败", 400)
			}
			return
		}
		var body struct {
			Name string `json:"name"`
		}
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&body); err != nil || body.Name == "" {
			http.Error(w, "昵称非法", 400)
			return
		}
		if err = decoder.Decode(new(any)); err != io.EOF {
			http.Error(w, "只接受一个 JSON", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{id, body.Name})
	})
	return mux
}
func fetchMessage(client *http.Client, url string) (string, error) {
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", fmt.Errorf("HTTP 状态 %d", response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

*/
