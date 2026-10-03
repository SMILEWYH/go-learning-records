// 演示以真正 HTTP 客户端访问临时服务；使用第 27 章的 httptest.NewServer。
package main

import (
	"errors"
	"example.com/go-shop-orders/internal/app"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func run() (err error) {
	dir, err := os.MkdirTemp("", "go-shop-orders-demo-*")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	service, err := app.Open(filepath.Join(dir, "state.json"))
	if err != nil {
		return err
	}
	server := httptest.NewServer(app.Handler(service, "demo"))
	defer server.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	defer client.CloseIdleConnections()
	call := func(method, path, body string, want int) error {
		request, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if err != nil {
			return err
		}
		if body != "" {
			request.Header.Set("Content-Type", "application/json")
		}
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		data, err := io.ReadAll(io.LimitReader(response.Body, 65537))
		if err != nil {
			return err
		}
		fmt.Printf("\n%s %s → %d\n%s", method, path, response.StatusCode, data)
		if response.StatusCode != want {
			return fmt.Errorf("期望 %d，实际 %d", want, response.StatusCode)
		}
		return nil
	}

	if err := call("GET", "/products", "", 200); err != nil {
		return err
	}
	body := `{"request_id":"demo-001","customer":"小林","items":[{"product_id":"p-1001","quantity":2}]}`
	if err := call("POST", "/orders", body, 201); err != nil {
		return err
	}
	if err := call("POST", "/orders", body, 201); err != nil {
		return err
	}
	if err := call("PATCH", "/orders/ord-000001/status", `{"status":"canceled"}`, 200); err != nil {
		return err
	}
	if err := call("PATCH", "/orders/ord-000001/status", `{"status":"shipped"}`, 409); err != nil {
		return err
	}
	if err := call("GET", "/products", "", 200); err != nil {
		return err
	}

	fmt.Println("\n演示通过；临时服务和数据自动清理。")
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
