// Package filedata 对应第 20、21、23 章：完整读取、严格 JSON、同目录临时文件替换。
package filedata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const MaxBytes = 16 << 20

func Load(path string, target any) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	data, err := io.ReadAll(io.LimitReader(f, MaxBytes+1))
	if err != nil {
		return err
	}
	if len(data) > MaxBytes {
		return errors.New("数据文件超过 16 MiB")
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return errors.New("数据文件不能是 null")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("解析数据文件: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("数据文件必须只包含一个 JSON 值")
	}
	return nil
}

// Save 成功返回后调用方才更新内存；失败时旧目标仍可读。
// 面向本机单进程的小数据文件。文件 Sync 不等于掉电时目录项也一定持久化。
func Save(path string, value any) (err error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if len(data) > MaxBytes {
		return errors.New("数据文件超过 16 MiB，请归档历史数据")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".saving-*")
	if err != nil {
		return err
	}
	name := f.Name()
	closed := false
	defer func() {
		if !closed {
			err = errors.Join(err, f.Close())
		}
		if e := os.Remove(name); e != nil && !errors.Is(e, os.ErrNotExist) {
			err = errors.Join(err, e)
		}
	}()
	if n, e := f.Write(data); e != nil {
		return e
	} else if n != len(data) {
		return io.ErrShortWrite
	}
	if err = f.Sync(); err != nil {
		return err
	}
	err = f.Close()
	closed = true
	if err != nil {
		return err
	}
	return os.Rename(name, path)
}
