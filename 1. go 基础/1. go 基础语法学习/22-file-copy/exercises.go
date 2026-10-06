// 章节：22-文件复制
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；请填写下方受测函数，具体对应关系见 exercises_test.go 的 TestExam。
// exercise1/2/3 保留题目说明，空演示入口不影响已经完成的业务函数。

// 本章所有源文件、目标文件及故障样本都限定在自己创建的临时目录中，结束后清理。
// 不复制真实用户文件。每个函数明确规定覆盖规则，并检查 io.Copy 与 Close 的错误。

// 练习 1：导出文件的首次备份。
// 建议：func backupNew(sourcePath, destinationPath string) (int64, error)。
// 要求：源为普通文件；目标使用 O_CREATE|O_EXCL|O_WRONLY，不覆盖已有文件。
// 常规：源为 5 字节 "hello"，新目标获得同样内容，返回 5, nil。
// 边界：空源成功且返回 0；源缺失返回 error。
// 目标已存在返回可被 errors.Is(err, os.ErrExist) 识别的错误，且原内容不变。
// 注意：先确认源可读取，再创建目标；若复制失败，明确返回错误，调用者不能使用半成品。
func exercise1() {
	// TODO：自己实现复制函数，验证目标存在时不会丢失旧内容。
}

// 练习 2：复制上传流，同时计算 SHA-256，供内容一致性核验使用。
// 建议：func copyAndHash(destination io.Writer, source io.Reader) (int64, string, error)。
// 要求：使用 io.Copy、io.MultiWriter 和 sha256.New，一次遍历完成；返回十六进制摘要。
// 常规：输入 "abc"，写入 3 字节；摘要等于 fmt.Sprintf("%x", sha256.Sum256([]byte("abc")))。
// 边界：空输入写入 0 字节且有合法摘要；目标写入失败应返回 error，此时摘要返回空串。
// 用 bytes.Buffer 演示正常路径，自己写一个始终返回错误的 Writer 演示故障路径。
// 本函数不关闭调用者传入的接口；创建文件的一方负责关闭文件。
func exercise2() {
	// TODO：实现组合 Reader/Writer 的复制流程，不把完整源读入内存。
}

// 练习 3：覆盖更新本地构建产物。
// 建议：func replaceArtifact(sourcePath, destinationPath string) (int64, error)。
// 要求：自己实现，先检查已打开文件的身份，再截断目标，不能直接调用讲解中的 copyFile。
// 常规：用 "new" 覆盖较长的 "old-content" -> 目标恰好为 "new"，返回 3。
// 边界：源与目标同路径 -> error 且源不变；目标是源的硬链接也必须拒绝；源缺失时目标不变。
// 硬链接用 os.Link 在临时目录创建；环境不支持时清楚输出跳过原因。
// 写一句注释：复制中途失败时，本函数为什么不能保证目标保留旧内容？
func exercise3() {
	// TODO：实现身份检查和复制，使用临时文件验证截断与拒绝规则。
}

// 自动测试入口：填写下面的 TODO 函数体，保留签名。
// 类型与字段已提供骨架；exercise1/2/3 仍可用于手动演示。

func backupNew(sourcePath, destinationPath string) (written int64, err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：backupNew；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func copyAndHash(destination io.Writer, source io.Reader) (int64, string, error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：copyAndHash；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
}

func replaceArtifact(sourcePath, destinationPath string) (written int64, err error) {
	// TODO：按上方题目要求实现；保持函数签名不变。
	panic("待完成：replaceArtifact；请实现 exercises.go 中的函数，参考答案在 exercises_test.go 末尾")
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
