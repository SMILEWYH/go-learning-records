/*
章节：30-编译与部署
	本章掌握清单：
	1. 构建产物、版本与目标平台
	2. 运行配置：环境变量、默认值与校验
	3. 健康检查与版本路由
	4. 启动服务与优雅退出
	5. 命令行分发与信号取消
	6. 程序入口、错误退出与部署自查
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./30-deployment
// 运行练习：go run ./30-deployment exam
// 阅读顺序：按第 1–6 节阅读构建、配置、路由、退出与入口；默认运行只打印说明，server 模式才监听端口。
// TS 对照：Go build 常直接生成可执行程序；纯 Go 服务运行时通常不需要安装 Node.js 或 Go。
// 运行服务：go run ./30-deployment server；查看版本：go run ./30-deployment version
// 详细命令见本章 lesson.go 末尾的部署附录。默认讲解只打印说明，不启动常驻服务或执行部署。
// 后续再学：systemd/容器、TLS 反向代理、日志与指标、CI/CD、滚动更新与回滚。

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// ═══════════════════════════════════════════════════════
// 第 1 节：构建产物、版本与目标平台
// ═══════════════════════════════════════════════════════
// go run 临时构建并执行；go build 生成可执行产物，-o 只指定输出位置和名字。
// 从 go 基础语法学习 根目录可运行：
//   go build -o /tmp/go-lesson-service ./30-deployment
//   /tmp/go-lesson-service version
// -X main.version=v0.1.0 给下面的 string 变量注入版本，默认值 dev 用于本地学习。
// 例如：go build -ldflags "-X main.version=v0.1.0" -o /tmp/go-lesson-service ./30-deployment
// 交叉编译还需设置 GOOS、GOARCH；Linux 产物一般不能直接在 macOS 上执行。

// 构建时可用 -ldflags '-X main.version=v0.1.0' 注入 string 变量。
var version = "dev"

// ═══════════════════════════════════════════════════════
// 第 2 节：运行配置：环境变量、默认值与校验
// ═══════════════════════════════════════════════════════
// defaultAddress 读取 PORT，未提供时用 8080，随后检查是否为 1..65535 的整数。
// 端口与主机通过 net.JoinHostPort 组成地址；本例只监听本机回环地址。
// 下方 server 分支先读 -addr，非空时直接使用；否则才进入这里解析 PORT。

func defaultAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return "", fmt.Errorf("PORT 必须为 1..65535，实际为 %q", port)
	}
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(number)), nil
}

// ═══════════════════════════════════════════════════════
// 第 3 节：健康检查与版本路由
// ═══════════════════════════════════════════════════════
// /healthz 返回 ok，说明进程能处理请求；/version 返回当前构建版本。
// 部署后可用 curl http://127.0.0.1:8080/healthz 核对响应，再核对 /version。

func deploymentMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(writer, "ok")
	})
	mux.HandleFunc("GET /version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(writer, version)
	})
	return mux
	// /healthz 表示进程能响应；真实 readiness 检查还需按业务判断数据库等依赖是否就绪。
}

// ═══════════════════════════════════════════════════════
// 第 4 节：启动服务与优雅退出
// ═══════════════════════════════════════════════════════
// 先 Listen 确认端口可用，再启动 Serve；主流程等待服务错误或取消通知。
// 取消时创建独立的 5 秒 Shutdown context，让正在处理的请求有时间结束。
// 若到期仍未完成，就 Close 强制关闭；最后读取 finished，确认 Serve 已返回。

func serveUntilCanceled(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	server := &http.Server{
		Handler: deploymentMux(), ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second,
	}
	finished := make(chan error, 1)
	go func() { finished <- server.Serve(listener) }()
	fmt.Printf("服务已启动：http://%s，版本=%s；Ctrl+C 触发优雅退出\n", listener.Addr(), version)
	select {
	case err := <-finished:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		fmt.Println("收到退出信号：停止接受新连接，等待正在处理的 HTTP 请求")
	}
	// 原 ctx 已取消，不能拿它直接做 Shutdown 的等待期限；新建独立的有期限 context。
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdownErr := server.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		// 到期还没完成时强制关闭，确保进程能够结束。
		shutdownErr = errors.Join(shutdownErr, server.Close())
	}
	serveErr := <-finished
	if errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil // 主动关停的正常结果，不能作为异常启动失败处理。
	}
	return errors.Join(shutdownErr, serveErr)
	// Shutdown 等待普通 HTTP 请求，不会自动关闭 WebSocket 或任意后台 goroutine。
	// 那些资源需要各自的取消与等待机制；直接 os.Exit 会跳过 defer。
}

// ═══════════════════════════════════════════════════════
// 第 5 节：命令行分发与信号取消
// ═══════════════════════════════════════════════════════
// 无参数打印构建说明，version 查看版本，server 启动服务并解析 -addr。
// signal.NotifyContext 把 Ctrl+C 或 SIGTERM 转成 ctx.Done 通知，第 4 节据此收尾。
// 完整平台与部署命令见同目录 lesson.go 末尾的部署附录。

func runDeployment(args []string) error {
	if len(args) == 0 {
		fmt.Println("当前构建目标：", runtime.GOOS, runtime.GOARCH, "；Go：", runtime.Version())
		fmt.Println("版本：", version)
		fmt.Println("1. go build 构建当前包；本项目根包是导航，要构建服务请指定章节目录。")
		fmt.Println("2. go build file.go 只编译列出的文件，本项目章节应构建整个包。")
		fmt.Println("3. -o 只指定文件名；把文件叫 main.exe 不会自动生成 Windows 程序。")
		fmt.Println("4. macOS -> Linux：设置 GOOS=linux 和匹配服务器 CPU 的 GOARCH。")
		fmt.Println("5. 纯 Go 示例可用 CGO_ENABLED=0；依赖 C 的项目需要另外处理交叉编译工具链。")
		fmt.Println("6. 阅读 30-deployment/lesson.go 末尾的部署附录，按命令构建、运行并检查 /healthz。")
		return nil
	}
	switch args[0] {
	case "version":
		fmt.Printf("version=%s go=%s target=%s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return nil
	case "server":
		flags := flag.NewFlagSet("server", flag.ContinueOnError)
		address := flags.String("addr", "", "监听地址，优先于 PORT，例如 127.0.0.1:9090")
		if err := flags.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return err
		}
		if flags.NArg() != 0 {
			return errors.New("server 只接受 -addr 参数")
		}
		if *address == "" {
			var err error
			*address, err = defaultAddress()
			if err != nil {
				return err
			}
		}
		// 配置优先级：非空 -addr > PORT > 默认 127.0.0.1:8080。
		// Unix 部署管理器通常发 SIGTERM；本地 Ctrl+C 对应 Interrupt。
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return serveUntilCanceled(ctx, *address)
	default:
		return errors.New("用法：不带参数，或 server [-addr 地址]，或 version，或 exam")
	}
}

// ═══════════════════════════════════════════════════════
// 第 6 节：程序入口、错误退出与部署自查
// ═══════════════════════════════════════════════════════
// main 将错误写入标准错误流，并以非零退出码告诉外部运行器本次失败。
// 实际关闭资源的 defer 位于被调用函数内，会在错误返回 main 之前执行。
// 检查顺序：构建成功 -> 目标平台匹配 -> 配置正确 -> 健康路由可访问 -> 退出能完成清理。

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := runDeployment(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "部署示例失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 先构建并运行 version，再用 -ldflags 注入不同版本，核对输出确实来自新产物。
// 启动 server，访问 /healthz 后按 Ctrl+C，沿取消信号追踪 Shutdown 与 Serve 的返回。

/*
附录：完整构建与部署命令

# 编译与部署：从 macOS 生成 Linux 程序

本章所有命令默认在项目根目录的 zsh / Bash 中执行。默认 `go run ./30-deployment` 只打印讲解；加 `server` 才会常驻。

## 1. 区分三种构建命令

| 命令 | 含义 |
| --- | --- |
| `go build` | 构建当前包。当前包是 `main` 时生成可执行文件；其他包只做编译检查。 |
| `go build xxx.go` | 仅构建明确列出的 Go 文件，不会自动补同目录其他文件。 |
| `go build -o main.exe xxx.go` | 指定输出名称；`.exe` 后缀不会改变目标操作系统。 |
| `go build -o ./bin/go-server ./30-deployment` | 构建整个章节，包含 `lesson.go` 和 `exercises.go`，排除 `_test.go`。 |

这个项目的根 `main.go` 是章节导航，在根目录单独 `go build` 得到的是导航程序。每章又是独立模块，根目录的 `go build ./...` 不会递归构建其他章节模块。单独构建 `lesson.go` 会缺少 `runExercises` 等同包函数。

`go run` 临时编译并运行，适合学习；发布时用 `go build` 留下确定的产物。`go build` 不会运行测试。

```sh
# 构建前检查；章节模块由 go.work 连接
go test ./... ./[0-9]* ./example
go vet ./... ./[0-9]* ./example
mkdir -p ./bin

# 本机产物：-trimpath 移除记录在程序中的本机源码路径
go build -trimpath -ldflags '-X main.version=v0.1.0' -o ./bin/go-server ./30-deployment
./bin/go-server version
./bin/go-server server -addr 127.0.0.1:8080
```

保持服务终端开启，在另一个终端执行：

```sh
curl -i http://127.0.0.1:8080/healthz
curl -i http://127.0.0.1:8080/version
```

预期分别为 `200` + `ok`、`200` + `v0.1.0`。回服务终端按 Ctrl+C，看到退出提示后进程结束。检查实际构建产物也能防止误把根目录导航程序当作服务发布。

## 2. 在 macOS 交叉编译

`GOOS` 选择操作系统，`GOARCH` 选择目标 CPU。目标由 Linux 服务器决定，与这台 Mac 是 Intel 还是 Apple Silicon 无关。

| 服务器 `uname -m` | Go 的目标架构 |
| --- | --- |
| `x86_64` | `GOARCH=amd64` |
| `aarch64` / `arm64` | `GOARCH=arm64` |

```sh
# 查看本机默认配置，以及工具链支持的目标组合
go env GOOS GOARCH CGO_ENABLED
go tool dist list

# Linux x86-64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-X main.version=v0.1.0' -o ./bin/go-server-linux-amd64 ./30-deployment

# Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags '-X main.version=v0.1.0' -o ./bin/go-server-linux-arm64 ./30-deployment

# 真正生成 Windows 程序需要设置 GOOS；只改扩展名没有作用
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/go-server-windows-amd64.exe ./30-deployment

# macOS 可以查看文件格式，但不能直接原生运行 Linux ELF 程序
file ./bin/go-server-linux-amd64 ./bin/go-server-linux-arm64
go version -m ./bin/go-server-linux-amd64
```

这几个环境变量只作用于当前这一条命令，不会永久改变本机 Go 配置。不要为了交叉编译而把全局 `GOOS` 改成 Linux，影响后续本机运行与测试。

本章是纯 Go 示例，关闭 CGO 即可跨平台编译。依赖 C 库的项目可能无法关闭 CGO，需要目标平台的 C 编译器、头文件和库。**不能由“本章构建成功”推断任意 Go 项目都只需这三个变量。**

纯 Go 产物通常无需在目标机安装 Go；配置文件、证书、模板等外部资源仍需按应用需要提供。操作系统最低版本与 CPU 特性也要匹配。

## 3. 到 Linux 上验证

下面只是手动操作模板。将占位主机替换为你自己的服务器；学习这章无需购买服务器，也不会自动上传或修改远端。

```sh
# 在 Mac 上：以 x86-64 目标为例
scp ./bin/go-server-linux-amd64 user@your-server:/tmp/go-server

# 在对应 Linux 服务器上执行
chmod +x /tmp/go-server
/tmp/go-server version
PORT=8080 /tmp/go-server server

# 在 Linux 的另一个终端检查
curl -i http://127.0.0.1:8080/healthz
```

配置优先级为非空 `server -addr ...` > `PORT` > 默认 `127.0.0.1:8080`。`PORT=9090 ./bin/go-server server` 只为这次运行设置环境变量；无效端口会报错并非零退出。`server -h` 可查看参数。

`127.0.0.1` 只接受本机连接。部署到容器或让其他机器访问时，可根据实际网络设置使用 `server -addr 0.0.0.0:8080`；访问方使用服务器真实 IP 或域名，`0.0.0.0` 是监听地址。还需配好端口映射、防火墙与反向代理；公网入口通常由代理提供 HTTPS。

这个示例只有健康检查和版本接口，不需要数据库。构建产物不会因为终端退出而自动转为后台常驻服务；正式部署应由 systemd、容器等管理启动、重启和日志。发布新版本后检查健康状态，再切流量；保留上一版本以便回滚。

## 4. 服务如何退出

本章 `serveUntilCanceled` 展示完整顺序：

1. 捕获 Ctrl+C 或 SIGTERM。
2. 调用 `Shutdown`，停止接受新连接，等待正在处理的 HTTP 请求。
3. 用新的 5 秒 context 限制等待；超时后 `Close` 强制收尾。
4. 等待 `Serve` 返回，完成清理再结束进程。

不能直接把已取消的信号 context 传给 `Shutdown`，否则它没有等待余地。`http.ErrServerClosed` 是主动关停的正常返回值。WebSocket 已升级连接、数据库连接、后台任务需要各自关闭和等待；可结合第 14–17、29 章继续练习。

## 5. 常见错误

| 现象 | 优先检查 |
| --- | --- |
| `exec format error` | GOOS/GOARCH 是否与目标机匹配。 |
| `permission denied` | 文件是否有执行权限，所在挂载点是否允许执行。 |
| `address already in use` | 端口已被其他进程占用；修改 PORT 或 -addr。 |
| 只打印章节目录后退出 | 构建了根导航包，应指定第 30 章目录。 |
| 本机可访问，其他机器访问失败 | 是否只监听回环地址、端口映射与网络规则是否正确。 |
| 二进制存在却提示加载器缺失 | 检查是否启用了 CGO、目标系统是否提供所需动态加载器与库。 |

参考：[Go build 命令](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)、[构建环境变量](https://pkg.go.dev/cmd/go#hdr-Environment_variables)、[Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)。

*/
