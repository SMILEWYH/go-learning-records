// 章节：23-目录操作
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 本章所有文件和目录都放在 os.MkdirTemp 创建的根目录内，结束后清理。
// 路径参数只使用自己创建的样本，不对用户目录执行 RemoveAll。

// 练习 1：准备按年月归档的上传目录。
// 建议：func ensureUploadDirectory(root string, year, month int) (string, error)。
// 要求：year 为 2000..2100、month 为 1..12；得到 root/uploads/YYYY/MM；使用 filepath.Join/MkdirAll。
// 常规：2026, 9 -> 目录 uploads/2026/09；重复调用成功，已有文件保留。
// 边界：month=0、13 或非法年份 -> error 且不创建子目录；目标路径上已有普通文件 -> error。
// root 仅由本题的临时目录提供，年份月份均验证后再转换成路径段。
func exercise1() {
	// TODO：实现目录准备并验证重复调用与非法输入。
}

// 练习 2：列出可在前端展示的本地图片文件。
// 建议：func listImages(root string) ([]string, error)。
// 要求：只读根目录一层；返回普通文件的文件名；支持 .png/.jpg/.webp，扩展名忽略大小写；按名称排序。
// 常规：a.png、b.JPG、readme.txt -> ["a.png","b.JPG"]。
// 边界：目录名 fake.png、符号链接 link.png 都不计入；空目录 -> 长度 0；根目录不存在 -> error。
// 通过 DirEntry.Type/Info 判断类型并处理 Info 错误，不递归进入子目录。
func exercise2() {
	// TODO：创建临时样本目录和文件，验证筛选与排序。
}

// 练习 3：统计构建产物体积。
// 建议：func artifactBytes(root string) (int64, error)。
// 要求：WalkDir 递归累加普通文件的 Size；跳过名为 cache 的子目录；忽略符号链接。
// 常规：a.txt 为 5 字节、nested/b.txt 为 3 字节、cache/x.txt 为 99 字节 -> 8。
// 边界：空目录 -> 0；路径不存在或获取 Info 失败 -> error，出错时不能报告成功统计。
// 回调先检查传入 error，再访问 DirEntry；遍历起点 root 必须是目录。
func exercise3() {
	// TODO：实现统计函数并核对嵌套、跳过目录、空目录及错误输入。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func ensureUploadDirectory(root string, year, month int) (string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：ensureUploadDirectory；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func listImages(root string) ([]string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：listImages；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func artifactBytes(root string) (int64, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：artifactBytes；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

// ── 自动检查启动入口：填写习题时无需修改 ──
// runExercises 只启动本章的检查入口；检查规则、测试和参考答案都在 exercises_test.go。
// Go 的普通运行不会编译 _test.go，因此通过 go test 进入 TestExam。
func runExercises() {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法定位本章源码，检查未执行")
		return
	}
	dir := filepath.Dir(source)
	if !filepath.IsAbs(dir) { // 兼容以 -trimpath 构建后在课程根目录或章节目录运行。
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("无法定位章节：", err)
			return
		}
		dir = cwd
		if _, err := os.Stat(filepath.Join(dir, "exercises.go")); err != nil {
			dir = filepath.Join(cwd, filepath.Base(filepath.Dir(source)))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "test", "-json", "-count=1", "-run", "^TestExam$", "-timeout", "5m", ".")
	command.Dir = dir
	command.Env = append(os.Environ(), "GO_LESSON_EXAM=1")
	command.WaitDelay = 2 * time.Second
	output, runErr := command.Output()
	decoder := json.NewDecoder(bytes.NewReader(output))
	var diagnostics strings.Builder
	passed := false
	for {
		var event struct{ Action, Test, Output string }
		if err := decoder.Decode(&event); err != nil {
			break
		}
		diagnostics.WriteString(event.Output)
		if event.Test == "TestExam" {
			if event.Action == "pass" {
				passed = true
			}
			if event.Action == "output" && !strings.HasPrefix(event.Output, "=== RUN") && !strings.HasPrefix(event.Output, "--- PASS") {
				fmt.Print(event.Output)
			}
		}
	}
	if runErr != nil || !passed {
		fmt.Println("练习检查未正常结束：", runErr)
		fmt.Print(diagnostics.String())
		if failure, ok := runErr.(*exec.ExitError); ok {
			fmt.Print(string(failure.Stderr))
		}
	}
}
