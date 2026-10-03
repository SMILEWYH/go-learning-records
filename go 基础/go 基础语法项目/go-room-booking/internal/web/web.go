// Package web 只组合第 18、19、27、30 章已经学过的能力。
package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string              { return e.Message }
func Fail(status int, message string) error { return &Error{status, message} }
func Reply(w http.ResponseWriter, status int, value any) {
	body, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "响应编码失败", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err = w.Write(append(body, '\n')); err != nil {
		fmt.Fprintln(os.Stderr, "写响应:", err)
	}
}
func ReplyError(w http.ResponseWriter, err error) {
	var known *Error
	if errors.As(err, &known) {
		Reply(w, known.Status, map[string]string{"error": known.Message})
		return
	}
	fmt.Fprintln(os.Stderr, "内部错误:", err)
	Reply(w, 500, map[string]string{"error": "保存或处理失败，请稍后重试"})
}
func Decode(w http.ResponseWriter, r *http.Request, target any) error {
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return Fail(415, "请使用 application/json")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return decodeError(err)
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 || body[0] != '{' {
		return decodeError(nil)
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return decodeError(err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return decodeError(err)
	}
	return nil
}
func decodeError(err error) error {
	var large *http.MaxBytesError
	if errors.As(err, &large) {
		return Fail(413, "请求体超过 16 KiB")
	}
	return Fail(400, "请求体必须是单个 JSON 对象，字段和类型须符合接口约定")
}

type Page[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}

func Paginate[T any](r *http.Request, items []T) (Page[T], error) {
	page, size := 1, 20
	for key, target := range map[string]*int{"page": &page, "size": &size} {
		if raw, ok := r.URL.Query()[key]; ok {
			if len(raw) != 1 {
				return Page[T]{}, Fail(400, "分页参数不能重复")
			}
			n, err := strconv.Atoi(raw[0])
			if err != nil {
				return Page[T]{}, Fail(400, "分页参数必须为整数")
			}
			*target = n
		}
	}
	if page < 1 || size < 1 || size > 100 {
		return Page[T]{}, Fail(400, "page 至少为 1，size 为 1..100")
	}
	result := Page[T]{Items: []T{}, Total: len(items), Page: page, Size: size}
	// 先用除法判断，再算偏移，避免极大的页码造成乘法溢出。
	if len(items) == 0 || page-1 > (len(items)-1)/size {
		return result, nil
	}
	start := (page - 1) * size
	result.Items = append(result.Items, items[start:min(start+size, len(items))]...)
	return result, nil
}
func Address(defaultPort string, override string) (string, error) {
	if override != "" {
		return override, nil
	}
	raw := os.Getenv("PORT")
	if raw == "" {
		raw = defaultPort
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > 65535 {
		return "", errors.New("PORT 必须是 1..65535")
	}
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(n)), nil
}
func Serve(ctx context.Context, addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	fmt.Println("服务已启动：http://" + listener.Addr().String() + "；Ctrl+C 退出")
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(shutdown)
	if err != nil {
		err = errors.Join(err, server.Close())
	}
	serveErr := <-done
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	fmt.Println("服务已停止")
	return errors.Join(err, serveErr)
}
