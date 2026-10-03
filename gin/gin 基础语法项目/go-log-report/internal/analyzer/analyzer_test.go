package analyzer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fixture(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
}

const good = `{"time":"2026-09-27T09:00:00+08:00","service":"orders","status":200,"duration_ms":100}`

func TestReportCountsAndWorkerDeterminism(t *testing.T) {
	dir := t.TempDir()
	fixture(t, dir, "a.jsonl", good+"\n"+strings.Replace(strings.Replace(good, `"status":200`, `"status":500`, 1), `"duration_ms":100`, `"duration_ms":300`, 1)+"\nbad\n")
	fixture(t, dir, "b.jsonl", strings.Replace(good, `"status":200`, `"status":404`, 1)) // 最后一行没有换行。
	first, err := Analyze(context.Background(), dir, Options{Workers: 1})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Analyze(context.Background(), dir, Options{Workers: 4})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("worker count changed output")
	}
	if first.Files != 2 || first.ValidLines != 3 || first.InvalidLines != 1 || len(first.Services) != 1 || len(first.Issues) != 1 || first.Issues[0].Line != 3 {
		t.Fatalf("report=%+v", first)
	}
	stat := first.Services[0]
	if stat.Requests != 3 || stat.ServerErrors != 1 || stat.ClientErrors != 1 || stat.TotalMS != 500 || stat.MaxMS != 300 || stat.AverageMS != float64(500)/3 || stat.ErrorRate != float64(1)/3 {
		t.Fatalf("stat=%+v", stat)
	}
	if _, err := Analyze(context.Background(), dir, Options{Workers: 2, Strict: true}); err == nil {
		t.Fatal("strict accepted invalid lines")
	}
}
func TestValidation(t *testing.T) {
	for _, body := range []string{"null", "", good + good, strings.Replace(good, `,"duration_ms":100`, "", 1), strings.Replace(good, `"duration_ms":100`, `"duration_ms":-1`, 1), strings.Replace(good, `"status":200`, `"status":600`, 1), strings.Replace(good, `"service":"orders"`, `"service":"Bad Name"`, 1), strings.Replace(good, `"duration_ms":100`, `"duration_ms":null`, 1), strings.Replace(good, `"status":200`, `"extra":1,"status":200`, 1)} {
		if _, err := parseLine([]byte(body)); err == nil {
			t.Fatal("accepted", body)
		}
	}
	if _, err := parseLine([]byte(strings.Replace(good, `"duration_ms":100`, `"duration_ms":0`, 1))); err != nil {
		t.Fatal(err)
	}
}
func TestCancellationAndFileErrors(t *testing.T) {
	dir := t.TempDir()
	fixture(t, dir, "a.jsonl", good)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Analyze(ctx, dir, Options{Workers: 2}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	fixture(t, dir, "too-long.jsonl", strings.Repeat("a", 70*1024))
	if _, err := Analyze(context.Background(), dir, Options{Workers: 2}); err == nil {
		t.Fatal("oversized line accepted")
	}
	if _, err := Analyze(context.Background(), dir, Options{Workers: 0}); err == nil {
		t.Fatal("invalid workers")
	}
	if _, err := Analyze(context.Background(), t.TempDir(), Options{Workers: 2}); err == nil {
		t.Fatal("empty directory accepted")
	}
}
func TestEmptyFileAndIgnoredEntries(t *testing.T) {
	dir := t.TempDir()
	fixture(t, dir, "empty.jsonl", "")
	fixture(t, dir, "notes.txt", "not a log")
	if err := os.Mkdir(filepath.Join(dir, "nested.jsonl"), 0700); err != nil {
		t.Fatal(err)
	}
	fixture(t, filepath.Join(dir, "nested.jsonl"), "child.jsonl", good)
	outside := filepath.Join(t.TempDir(), "outside.jsonl")
	if err := os.WriteFile(outside, []byte(good), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "link.jsonl")); err != nil {
		t.Log("symlink unavailable:", err)
	}
	report, err := Analyze(context.Background(), dir, Options{Workers: 2})
	if err != nil || report.Files != 1 || report.ValidLines != 0 || report.Services == nil {
		t.Fatal(report, err)
	}
}
func TestIssueLimit(t *testing.T) {
	dir := t.TempDir()
	fixture(t, dir, "bad.jsonl", strings.Repeat("bad\n", 120))
	report, err := Analyze(context.Background(), dir, Options{Workers: 1})
	if err != nil || report.InvalidLines != 120 || len(report.Issues) != 100 || !report.IssuesTruncated {
		t.Fatal(report, err)
	}
}
