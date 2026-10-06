package api

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/go-log-report/internal/analyzer"
)

func TestReports(t *testing.T) {
	config := Config{Input: "../../samples", Workers: 2, Timeout: 5 * time.Second}
	router := Handler(config, "test")
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/healthz", 200},
		{"GET", "/reports", 200},
		{"GET", "/reports?workers=1&strict=false", 200},
		{"GET", "/reports?strict=true", 422},
		{"GET", "/reports?workers=0", 400},
		{"GET", "/reports?workers=17", 400},
		{"GET", "/reports?workers=x", 400},
		{"GET", "/reports?workers=", 400},
		{"GET", "/reports?workers=1&workers=2", 400},
		{"GET", "/reports?strict=maybe", 400},
		{"GET", "/reports?strict=", 400},
		{"GET", "/reports?input=/etc", 400},
		{"GET", "/reports?strict=%zz", 400},
		{"POST", "/reports", 405},
		{"GET", "/missing", 404},
	} {
		t.Run(tc.path+tc.method, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status || !json.Valid(w.Body.Bytes()) {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 200 && tc.path != "/healthz" {
				var report analyzer.Report
				if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
					t.Fatal(err)
				}
				if report.Files != 2 || report.ValidLines != 6 || report.InvalidLines != 1 {
					t.Fatalf("unexpected report: %+v", report)
				}
			}
		})
	}
}

func TestReportTimeoutCancellationAndMissingInput(t *testing.T) {
	for _, tc := range []struct {
		name, input string
		timeout     time.Duration
		cancel      bool
		status      int
	}{
		{"timeout", "../../samples", time.Nanosecond, false, 504},
		{"cancel", "../../samples", time.Second, true, 408},
		{"missing", filepath.Join(t.TempDir(), "private-dir"), time.Second, false, 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := Handler(Config{Input: tc.input, Workers: 1, Timeout: tc.timeout}, "test")
			r := httptest.NewRequest("GET", "/reports", nil)
			if tc.cancel {
				ctx, cancel := context.WithCancel(r.Context())
				cancel()
				r = r.WithContext(ctx)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestStrictDefaultCanBeOverridden(t *testing.T) {
	router := Handler(Config{Input: "../../samples", Workers: 1, Strict: true, Timeout: time.Second}, "test")
	for path, want := range map[string]int{"/reports": 422, "/reports?strict=false": 200} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != want {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
	}
	// 严格模式下完全有效的日志依然成功。
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "empty.jsonl"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	router = Handler(Config{Input: dir, Workers: 1, Strict: true, Timeout: time.Second}, "test")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/reports", nil))
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
}

func TestBusyRequestAndCancellationReleaseSlot(t *testing.T) {
	entered := make(chan struct{}, 1)
	analyze := func(ctx context.Context, _ string, _ analyzer.Options) (analyzer.Report, error) {
		entered <- struct{}{}
		<-ctx.Done()
		return analyzer.Report{}, ctx.Err()
	}
	router := newHandler(Config{Workers: 1, Timeout: time.Minute}, "test", analyze)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan int, 1)
	go func() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", "/reports", nil).WithContext(ctx))
		done <- w.Code
	}()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("analysis did not start")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/reports", nil))
	if w.Code != 429 || w.Header().Get("Retry-After") != "1" {
		t.Fatal(w.Code, w.Body.String())
	}
	// 正在分析时健康检查仍然能响应。
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/healthz", nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	cancel()
	select {
	case code := <-done:
		if code != 408 {
			t.Fatal(code)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("request did not cancel")
	}
	// 取消后释放名额：新请求能进入分析函数，而不是继续返回 429。
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/reports", nil).WithContext(ctx))
	if w.Code != 408 {
		t.Fatal(w.Code, w.Body.String())
	}
	select {
	case <-entered:
	default:
		t.Fatal("analysis slot was not released")
	}
}
