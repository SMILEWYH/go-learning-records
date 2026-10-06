package filedata

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveReloadAndReplace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data", "state.json")
	for _, value := range []string{"long original content", "x"} {
		if err := Save(path, map[string]string{"name": value}); err != nil {
			t.Fatal(err)
		}
		var got map[string]string
		if err := Load(path, &got); err != nil || got["name"] != value {
			t.Fatalf("got=%v err=%v", got, err)
		}
	}
	// JSON 编码失败不得覆盖已有内容。
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := Save(path, make(chan int)); err == nil {
		t.Fatal("channel must fail")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatal("failed save changed old file")
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil || len(entries) != 1 {
		t.Fatal("temporary file leaked", err)
	}
}
func TestLoadRejectsBadData(t *testing.T) {
	for _, data := range []string{"null", "{}{}", `{"extra":1}`, "{"} {
		t.Run(data, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "bad.json")
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			var target struct {
				Name string `json:"name"`
			}
			if err := Load(path, &target); err == nil {
				t.Fatal("expected invalid data")
			}
		})
	}
	var target any
	if err := Load(filepath.Join(t.TempDir(), "absent"), &target); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
func TestSaveFailureKeepsDestination(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "occupied")
	if err := os.Mkdir(target, 0700); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(target, "keep.txt")
	if err := os.WriteFile(keep, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Save(target, []int{1}); err == nil {
		t.Fatal("cannot replace directory")
	}
	data, err := os.ReadFile(keep)
	if err != nil || string(data) != "keep" {
		t.Fatal("old contents changed", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatal("temporary file leaked", err)
	}
}
