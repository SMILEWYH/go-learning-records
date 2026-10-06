// 演示以真正 HTTP 客户端访问临时服务；使用第 27 章的 httptest.NewServer。
package main

import (
	"errors"
	"example.com/go-room-booking/internal/app"
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
	dir, err := os.MkdirTemp("", "go-room-booking-demo-*")
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

	if err := call("GET", "/rooms", "", 200); err != nil {
		return err
	}
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(15 * time.Minute)
	body := fmt.Sprintf(`{"room_id":"room-a","organizer":"小林","title":"需求评审","attendees":3,"start":%q,"end":%q}`, start.Format(time.RFC3339), start.Add(time.Hour).Format(time.RFC3339))
	if err := call("POST", "/bookings", body, 201); err != nil {
		return err
	}
	if err := call("POST", "/bookings", body, 409); err != nil {
		return err
	}
	if err := call("POST", "/bookings/booking-000001/cancel", "", 200); err != nil {
		return err
	}
	if err := call("POST", "/bookings", body, 201); err != nil {
		return err
	}
	if err := call("GET", "/bookings?status=active", "", 200); err != nil {
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
