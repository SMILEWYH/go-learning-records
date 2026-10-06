// Package analyzer 用第 14–16 章的有界任务队列处理第 20、23 章的日志文件。
package analyzer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

const maxFileBytes = 10 << 20
const maxLines = 100_000

type Event struct {
	Time       string `json:"time"`
	Service    string `json:"service"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
}
type Issue struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Reason string `json:"reason"`
}
type Stats struct {
	Service      string  `json:"service"`
	Requests     int64   `json:"requests"`
	ClientErrors int64   `json:"client_errors"`
	ServerErrors int64   `json:"server_errors"`
	TotalMS      int64   `json:"total_ms"`
	MaxMS        int64   `json:"max_ms"`
	AverageMS    float64 `json:"average_ms"`
	ErrorRate    float64 `json:"server_error_rate"`
}
type Report struct {
	Files           int     `json:"files"`
	ValidLines      int64   `json:"valid_lines"`
	InvalidLines    int64   `json:"invalid_lines"`
	Services        []Stats `json:"services"`
	Issues          []Issue `json:"issues"`
	IssuesTruncated bool    `json:"issues_truncated"`
}
type Options struct {
	Workers int
	Strict  bool
}
type fileResult struct {
	file    string
	valid   int64
	invalid int64
	stats   map[string]Stats
	issues  []Issue
	err     error
}

// 用指针区分 duration_ms 缺失与显式为 0；不接受 null、未知字段或一行多个 JSON。
func parseLine(data []byte) (Event, error) {
	var in struct {
		Time       string `json:"time"`
		Service    string `json:"service"`
		Status     int    `json:"status"`
		DurationMS *int64 `json:"duration_ms"`
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		return Event{}, errors.New("JSON 格式、字段或类型错误")
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return Event{}, errors.New("每行只能包含一个 JSON 对象")
	}
	if _, err := time.Parse(time.RFC3339, in.Time); err != nil {
		return Event{}, errors.New("time 必须是带时区的 RFC3339 时间")
	}
	if len(in.Service) < 1 || len(in.Service) > 64 {
		return Event{}, errors.New("service 长度为 1..64")
	}
	for _, c := range in.Service {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			return Event{}, errors.New("service 仅支持小写字母、数字、-_")
		}
	}
	if in.Status < 100 || in.Status > 599 {
		return Event{}, errors.New("status 必须为 100..599")
	}
	if in.DurationMS == nil || *in.DurationMS < 0 || *in.DurationMS > 3_600_000 {
		return Event{}, errors.New("duration_ms 必填且为 0..3600000 的整数")
	}
	return Event{in.Time, in.Service, in.Status, *in.DurationMS}, nil
}
func readFile(ctx context.Context, root *os.Root, name string) (result fileResult) {
	result = fileResult{file: name, stats: map[string]Stats{}, issues: []Issue{}}
	f, err := root.Open(name)
	if err != nil {
		result.err = err
		return
	}
	defer func() { result.err = errors.Join(result.err, f.Close()) }()
	info, err := f.Stat()
	if err != nil {
		result.err = err
		return
	}
	if !info.Mode().IsRegular() || info.Size() > maxFileBytes {
		result.err = errors.New("输入必须为不超过 10 MiB 的普通文件")
		return
	}
	limited := &io.LimitedReader{R: f, N: maxFileBytes + 1}
	scanner := bufio.NewScanner(limited)
	scanner.Buffer(make([]byte, 4096), 64*1024)
	line := 0
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			result.err = err
			return
		}
		line++
		if line > maxLines {
			result.err = errors.New("单文件超过 100000 行")
			return
		}
		event, err := parseLine(scanner.Bytes())
		if err != nil {
			result.invalid++
			if len(result.issues) < 100 {
				result.issues = append(result.issues, Issue{name, line, err.Error()})
			}
			continue
		}
		result.valid++
		stat := result.stats[event.Service]
		stat.Service = event.Service
		stat.Requests++
		stat.TotalMS += event.DurationMS
		stat.MaxMS = max(stat.MaxMS, event.DurationMS)
		if event.Status >= 400 && event.Status < 500 {
			stat.ClientErrors++
		}
		if event.Status >= 500 {
			stat.ServerErrors++
		}
		result.stats[event.Service] = stat
	}
	if err := scanner.Err(); err != nil {
		result.err = fmt.Errorf("扫描日志: %w", err)
		return
	}
	if limited.N == 0 {
		result.err = errors.New("读取期间文件超过 10 MiB")
		return
	}
	result.err = ctx.Err()
	return
}
func Analyze(ctx context.Context, dir string, options Options) (Report, error) {
	if options.Workers < 1 || options.Workers > 16 {
		return Report{}, errors.New("workers 必须为 1..16")
	}
	if err := ctx.Err(); err != nil {
		return Report{}, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return Report{}, err
	}
	defer root.Close()
	// 只读取根目录这一层；ReadDir 的排序使任务与最终报告可重复。
	directory, err := root.Open(".")
	if err != nil {
		return Report{}, err
	}
	entries, readErr := directory.ReadDir(-1)
	if err := errors.Join(readErr, directory.Close()); err != nil {
		return Report{}, err
	}
	names := []string{}
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 || strings.ToLower(filepath.Ext(entry.Name())) != ".jsonl" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return Report{}, err
		}
		if info.Mode().IsRegular() {
			names = append(names, entry.Name())
		}
	}
	slices.Sort(names)
	if len(names) == 0 {
		return Report{}, errors.New("输入目录没有普通 .jsonl 文件")
	}
	if len(names) > 100 {
		return Report{}, errors.New("单批最多处理 100 个日志文件")
	}
	jobs := make(chan string, options.Workers)
	results := make(chan fileResult, options.Workers)
	var wait sync.WaitGroup
	for i := 0; i < min(options.Workers, len(names)); i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case name, ok := <-jobs:
					if !ok {
						return
					}
					result := readFile(ctx, root, name)
					select {
					case results <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
	var producer sync.WaitGroup
	producer.Add(1)
	go func() {
		defer producer.Done()
		defer close(jobs)
		for _, name := range names {
			select {
			case jobs <- name:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() { wait.Wait(); close(results) }()
	collected := []fileResult{}
	for result := range results {
		collected = append(collected, result)
	}
	// 即使取消，也等投递者和全部读取者退出后，才关闭共享 root。
	producer.Wait()
	if err := ctx.Err(); err != nil {
		return Report{}, err
	}
	slices.SortFunc(collected, func(a, b fileResult) int { return strings.Compare(a.file, b.file) })
	report := Report{Files: len(names), Services: []Stats{}, Issues: []Issue{}}
	totals := map[string]Stats{}
	for _, r := range collected {
		if r.err != nil {
			return Report{}, fmt.Errorf("读取 %s: %w", r.file, r.err)
		}
		report.ValidLines += r.valid
		report.InvalidLines += r.invalid
		for _, issue := range r.issues {
			if len(report.Issues) < 100 {
				report.Issues = append(report.Issues, issue)
			}
		}
		for name, stat := range r.stats {
			sum := totals[name]
			sum.Service = name
			sum.Requests += stat.Requests
			sum.ClientErrors += stat.ClientErrors
			sum.ServerErrors += stat.ServerErrors
			sum.TotalMS += stat.TotalMS
			sum.MaxMS = max(sum.MaxMS, stat.MaxMS)
			totals[name] = sum
		}
	}
	report.IssuesTruncated = report.InvalidLines > int64(len(report.Issues))
	if options.Strict && report.InvalidLines > 0 {
		return Report{}, fmt.Errorf("严格模式拒绝 %d 行无效日志；首处 %s:%d（%s）", report.InvalidLines, report.Issues[0].File, report.Issues[0].Line, report.Issues[0].Reason)
	}
	for _, stat := range totals {
		stat.AverageMS = float64(stat.TotalMS) / float64(stat.Requests)
		stat.ErrorRate = float64(stat.ServerErrors) / float64(stat.Requests)
		report.Services = append(report.Services, stat)
	}
	slices.SortFunc(report.Services, func(a, b Stats) int { return strings.Compare(a.Service, b.Service) })
	return report, nil
}
