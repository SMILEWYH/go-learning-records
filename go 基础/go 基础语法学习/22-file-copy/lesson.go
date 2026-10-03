/*
章节：22-文件复制
	本章掌握清单：
	1. 复制前先准备并验证文件
	2. io.Copy：目标在前，来源在后
	3. 同文件检查必须早于截断
	4. Reader / Writer：复制流程不依赖磁盘
	5. 覆盖约定与失败后的状态
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./22-file-copy
// 运行练习：go run ./22-file-copy exam
// 阅读顺序：从 main 进入 demoCopy 看调用，再回到上方 copyFile 逐步阅读文件校验与清理。
// TS 对照：类似 fs.copyFile 或把 readable stream pipe 到 writable stream；
// io.Copy 使用 Reader/Writer 接口，不需要把整个大文件先读入内存。
// 后续再学：可恢复大文件传输、进度/取消、元数据复制、文件系统快照和并发路径替换。

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// copyFile 用于本地、受控路径，覆盖普通目标文件，不保留源权限/时间等元数据。
// 若复制失败，目标可能只有部分数据；需要完整替换语义时应另写临时文件再提交。
// 此处不处理其他进程同时替换路径/修改文件的情况。
func copyFile(sourcePath, destinationPath string) (written int64, err error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return 0, fmt.Errorf("打开源文件: %w", err)
	}
	defer func() { err = errors.Join(err, source.Close()) }()
	sourceInfo, err := source.Stat()
	if err != nil {
		return 0, err
	}
	if !sourceInfo.Mode().IsRegular() {
		return 0, errors.New("源必须是普通文件")
	}

	// 不能先 os.Create：它会立刻截断，若 src/dst 实际是同一文件就已经损坏源内容。
	// 先打开但不截断，比较已打开文件的身份，再清空目标。
	destination, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return 0, fmt.Errorf("打开目标文件: %w", err)
	}
	defer func() { err = errors.Join(err, destination.Close()) }()
	destinationInfo, err := destination.Stat()
	if err != nil {
		return 0, err
	}
	if !destinationInfo.Mode().IsRegular() {
		return 0, errors.New("目标必须是普通文件")
	}
	if os.SameFile(sourceInfo, destinationInfo) {
		// 单纯比较路径字符串不够：硬链接、符号链接也可能指向同一文件。
		return 0, errors.New("源与目标是同一文件")
	}
	if err := destination.Truncate(0); err != nil {
		return 0, fmt.Errorf("清空目标: %w", err)
	}
	written, err = io.Copy(destination, source)
	// written 是已成功写出的字节数，不是字符数；io.Copy 读到正常 EOF 时返回 nil 错误。
	// 出错时 written 仍可能大于 0，不能只根据“已经写出一些内容”就视作复制成功。
	if err != nil {
		return written, fmt.Errorf("复制内容: %w", err)
	}
	return written, nil
}

func demoCopy() (err error) {

	// ═══════════════════════════════════════════════════════
	// 第 1 节：复制前先准备并验证文件
	// ═══════════════════════════════════════════════════════
	// 本例把目标预先写得更长，用于验证复制后不会残留旧尾部。
	// 阅读 copyFile 时按“打开源 -> 验证 -> 打开目标 -> 验证 -> 截断 -> 复制”追踪。

	dir, err := os.MkdirTemp("", "go-copy-*")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	source := filepath.Join(dir, "source.txt")
	destination := filepath.Join(dir, "backup.txt")
	if err := os.WriteFile(source, []byte("report=ready\n"), 0o600); err != nil {
		return err
	}
	if err := os.WriteFile(destination, []byte("old content that is much longer\n"), 0o600); err != nil {
		return err
	}

	// ═══════════════════════════════════════════════════════
	// 第 2 节：io.Copy：目标在前，来源在后
	// ═══════════════════════════════════════════════════════
	// io.Copy(destination, source) 持续把源字节写到目标，正常读完返回 nil 错误。
	// copyFile 另负责文件的打开与关闭；成功字节数应与源内容长度一致。

	n, err := copyFile(source, destination)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(destination)
	if err != nil {
		return err
	}
	fmt.Printf("复制 %d 字节，目标内容：%q\n", n, content)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：同文件检查必须早于截断
	// ═══════════════════════════════════════════════════════
	// 不同路径也可能是同一文件的链接，必须用文件信息确认身份。
	// 如果先 os.Create 或 O_TRUNC，再检查，源数据可能已经被清空。

	_, sameFileErr := copyFile(source, source)
	fmt.Println("同文件复制被拒绝：", sameFileErr)

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Reader / Writer：复制流程不依赖磁盘
	// ═══════════════════════════════════════════════════════
	// 同一个 io.Copy 也能从 strings.Reader 写入 strings.Builder。
	// Go你好 含 4 个码点、8 个 UTF-8 字节，written 统计后者。

	// Reader/Writer 让同一复制流程也适用于内存、压缩流、网络流等。
	var output strings.Builder
	n, err = io.Copy(&output, strings.NewReader("Go你好"))
	if err != nil {
		return err
	}
	fmt.Printf("内存流复制：%d 字节，%d 个 Unicode 码点，内容=%q\n", n, len([]rune(output.String())), output.String())
	// io.Copy 不会关闭两端，创建资源的一方负责关闭。
	// io.CopyN 复制指定字节数，源不够时会报错；io.CopyBuffer 可复用缓冲区，减少重复分配。
	// 对实现 WriterTo/ReaderFrom 的对象，Copy/CopyBuffer 可能使用对象自己的优化路径。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：覆盖约定与失败后的状态
	// ═══════════════════════════════════════════════════════
	// 开始复制前决定：拒绝已有目标、允许直接覆盖，还是先写临时文件再替换。
	// 字节已经写入后发生错误，目标可能只完成一部分，需要明确向调用方报告。

	// 覆盖策略与失败收尾要一起设计：
	// - 目标不能已存在：用 O_CREATE|O_EXCL，避免“先检查存在，再创建”的竞态窗口。
	// - 允许覆盖但不能暴露半成品：参考上一章，同目录临时文件写好后再提交。
	// - 当前示例直接覆盖：失败可能留下不完整目标；不能无条件 Remove 目标，
	//   因为它可能是原先存在的用户文件。只清理由本次调用创建且明确归自己所有的临时文件。
	// io.Copy 只复制内容，不复制权限、所有者、修改时间或扩展属性；它也不会自动创建父目录。
	return nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := demoCopy(); err != nil {
		fmt.Fprintln(os.Stderr, "复制演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 核对 copyFile(source, source) 拒绝操作后源文件仍完整。
// 把内存流文本换成英文与中文，预测 written 为什么按字节而非码点增加。
