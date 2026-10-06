package main

import (
	"context"
	"encoding/json"
	"errors"
	"example.com/go-log-report/internal/analyzer"
	"example.com/go-log-report/internal/api"
	"example.com/go-log-report/internal/filedata"
	"example.com/go-log-report/internal/web"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var version = "dev"

func run() error {
	flags := flag.NewFlagSet("go-log-report", flag.ContinueOnError)
	serve := flags.Bool("serve", false, "启动 Gin HTTP 服务")
	addr := flags.String("addr", "", "HTTP 监听地址；优先于 PORT，默认 8093")
	input := flags.String("input", "./samples", "日志目录，只读取第一层 .jsonl 普通文件")
	output := flags.String("out", "", "可选 .json 报告路径；不填则输出至终端")
	workers := flags.Int("workers", 4, "并发文件数，1..16")
	timeout := flags.Duration("timeout", 10*time.Second, "分析时间上限，例如 5s")
	strict := flags.Bool("strict", false, "任一无效日志行导致整个批次失败")
	showVersion := flags.Bool("version", false, "显示构建版本")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("不接受位置参数，请使用 -h 查看帮助")
	}
	if *showVersion {
		fmt.Println(version)
		return nil
	}
	if *timeout <= 0 || *timeout > 10*time.Minute {
		return errors.New("timeout 必须大于 0 且不超过 10m")
	}
	if *workers < 1 || *workers > 16 {
		return errors.New("workers 必须为 1..16")
	}
	if *serve {
		if *output != "" {
			return errors.New("HTTP 模式直接返回报告，请移除 -out")
		}
		info, err := os.Stat(*input)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("input 必须是日志目录")
		}
		address, err := web.Address("8093", *addr)
		if err != nil {
			return err
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		handler := api.Handler(api.Config{Input: *input, Workers: *workers, Strict: *strict, Timeout: *timeout}, version)
		return web.Serve(ctx, address, handler)
	}
	if *addr != "" {
		return errors.New("-addr 需要配合 -serve")
	}
	if *output != "" {
		if strings.ToLower(filepath.Ext(*output)) != ".json" {
			return errors.New("报告必须使用 .json 扩展名，避免覆盖输入 .jsonl")
		}
		// 目标不能借硬链接指向任何日志输入。符号链接由 Rename 替换其目录项。
		target, e := os.Stat(*output)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			return e
		}
		if e == nil {
			entries, err := os.ReadDir(*input)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if strings.ToLower(filepath.Ext(entry.Name())) == ".jsonl" {
					info, err := entry.Info()
					if err != nil {
						return err
					}
					if os.SameFile(info, target) {
						return errors.New("报告目标与输入日志是同一个文件")
					}
				}
			}
		}
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(signalCtx, *timeout)
	defer cancel()
	report, err := analyzer.Analyze(ctx, *input, analyzer.Options{Workers: *workers, Strict: *strict})
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if *output != "" {
		if err := filedata.Save(*output, report); err != nil {
			return err
		}
		fmt.Printf("报告已保存到 %s：%d 个文件，%d 行有效，%d 行无效\n", *output, report.Files, report.ValidLines, report.InvalidLines)
		return nil
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "分析失败:", err)
		os.Exit(1)
	}
}
