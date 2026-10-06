package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDecodeBoundary(t *testing.T) {
	for _, tc := range []struct {
		body, media string
		status      int
	}{
		{`{"name":"ok"}`, "application/json", 0},
		{`{"name":"ok"} {}`, "application/json", 400},
		{`{"name":"ok","admin":true}`, "application/json", 400},
		{`{"name":1}`, "application/json", 400}, {`[]`, "application/json", 400},
		{`{}`, "text/plain", 415}, {"", "application/json", 400},
		{`{"name":"` + strings.Repeat("a", 16384) + `"}`, "application/json", 413},
		{`{"name":"ok"}` + strings.Repeat(" ", 16384), "application/json", 413},
	} {
		t.Run(strconv.Itoa(tc.status)+tc.media, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", tc.media)
			var value struct {
				Name string `json:"name"`
			}
			err := Decode(httptest.NewRecorder(), r, &value)
			if tc.status == 0 {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var known *Error
			if !errors.As(err, &known) || known.Status != tc.status {
				t.Fatalf("got %v want status %d", err, tc.status)
			}
		})
	}
}
func TestPagination(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	page, err := Paginate(httptest.NewRequest("GET", "/?page=2&size=2", nil), items)
	if err != nil || page.Total != 5 || len(page.Items) != 2 || page.Items[0] != 3 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	page.Items[0] = 99
	if items[2] != 3 {
		t.Fatal("page aliases original")
	}
	for _, query := range []string{"page=0", "size=101", "page=a", "size=", "page=1&page=2"} {
		if _, err := Paginate(httptest.NewRequest("GET", "/?"+query, nil), items); err == nil {
			t.Fatal(query)
		}
	}
	huge := strconv.Itoa(int(^uint(0) >> 1))
	page, err = Paginate(httptest.NewRequest("GET", "/?page="+huge+"&size=100", nil), items)
	if err != nil || len(page.Items) != 0 || page.Items == nil {
		t.Fatalf("huge page=%+v err=%v", page, err)
	}
}
func TestServeShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, "127.0.0.1:0", http.NewServeMux()) }()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown timed out")
	}
}
func TestAddressPriority(t *testing.T) {
	t.Setenv("PORT", "invalid")
	if addr, err := Address("8091", "127.0.0.1:9000"); err != nil || addr != "127.0.0.1:9000" {
		t.Fatal(addr, err)
	}
	if _, err := Address("8091", ""); err == nil {
		t.Fatal("invalid PORT accepted")
	}
}
