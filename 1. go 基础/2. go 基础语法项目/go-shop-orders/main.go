package main

import (
	"context"
	"errors"
	"example.com/go-shop-orders/internal/app"
	"example.com/go-shop-orders/internal/web"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var version = "dev"

func run() error {
	flags := flag.NewFlagSet("go-shop-orders", flag.ContinueOnError)
	addr := flags.String("addr", "", "监听地址；优先于 PORT")
	data := flags.String("data", "./data/state.json", "数据文件；相对于当前工作目录")
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
	address, err := web.Address("8091", *addr)
	if err != nil {
		return err
	}
	service, err := app.Open(*data)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return web.Serve(ctx, address, app.Handler(service, version))
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
