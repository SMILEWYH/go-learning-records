/*
章节：23-目录操作
	本章掌握清单：
	1. MkdirTemp / MkdirAll：创建目录树
	2. ReadDir：列出当前一层
	3. WalkDir：递归、错误与跳过子树
	4. 软链接与 Stat / Lstat
	5. 硬链接与同一文件身份
	6. 路径文字与实际访问边界
	7. Rename、Remove 与 RemoveAll
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./23-directories
// 运行练习：go run ./23-directories exam
// 阅读顺序：先读 demoDirectories 的第 1–3 节，再随调用进入 demoPathOperations 的第 4–7 节。
// TS 对照：MkdirAll 类似 fs.mkdir({recursive:true})；ReadDir 类似 fs.readdir。
// 用 filepath 处理操作系统文件路径；URL 路径使用 path 包，两者用途不同。
// 后续再学：文件监听、跨平台权限、文件系统挂载与更复杂的资源沙箱。

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func demoPathOperations(root string) (err error) {

	// ═══════════════════════════════════════════════════════
	// 第 4 节：软链接与 Stat / Lstat
	// ═══════════════════════════════════════════════════════
	// 先看 demoPathOperations：软链接保存路径，Lstat 看链接本身，Stat 看最终目标。
	// 删除链接名称不会删除它所指向的 README.txt。

	original := filepath.Join(root, "README.txt")
	symbolicLink := filepath.Join(root, "readme-link.txt")
	// 软链接存目标路径，目标可以不存在；相对目标按链接自身所在目录解释。
	if err := os.Symlink("README.txt", symbolicLink); err != nil {
		fmt.Println("当前环境无法创建软链接，跳过演示：", err)
	} else {
		linkInfo, err := os.Lstat(symbolicLink) // 查看链接本身，不跟随末尾链接。
		if err != nil {
			return err
		}
		targetInfo, err := os.Stat(symbolicLink) // 跟随链接，查看目标。
		if err != nil {
			return err
		}
		fmt.Println("Lstat 看到软链接：", linkInfo.Mode()&os.ModeSymlink != 0,
			"Stat 看到普通文件：", targetInfo.Mode().IsRegular())
		if err := os.Remove(symbolicLink); err != nil { // 删除链接，不删除链接目标。
			return err
		}
	}

	// ═══════════════════════════════════════════════════════
	// 第 5 节：硬链接与同一文件身份
	// ═══════════════════════════════════════════════════════
	// 硬链接给同一个文件增加一个名称，因此 SameFile 为 true。
	// 它不是复制一份数据；删除其中一个名字后，另一个名字仍可访问文件。

	hardLink := filepath.Join(root, "readme-hard.txt")
	// 硬链接是同一底层文件的另一名称，不是内容副本；常要求位于同一文件系统。
	if err := os.Link(original, hardLink); err != nil {
		fmt.Println("当前环境无法创建硬链接，跳过演示：", err)
	} else {
		originalInfo, err := os.Stat(original)
		if err != nil {
			return err
		}
		linkedInfo, err := os.Stat(hardLink)
		if err != nil {
			return err
		}
		fmt.Println("硬链接指向同一文件：", os.SameFile(originalInfo, linkedInfo))
		if err := os.Remove(hardLink); err != nil {
			return err
		}
	}

	// ═══════════════════════════════════════════════════════
	// 第 6 节：路径文字与实际访问边界
	// ═══════════════════════════════════════════════════════
	// Join、Clean、IsLocal 处理路径表达式，不能保证链接解析后的实际位置。
	// os.OpenRoot 返回受根目录约束的访问入口，再通过它读取相对路径。

	// IsLocal 只检查路径文字，不解析实际符号链接；Join/Clean/Rel 也不是访问沙箱。
	fmt.Println("词法本地路径：", filepath.IsLocal("uploads/avatar.png"), filepath.IsLocal("../outside.txt"))
	// Go 1.24 起的 os.Root 提供受根目录限制的文件访问；当前项目工具链支持此 API。
	// 相对软链接可在根内解析，但不允许经 .. 或链接逃到根目录之外。
	limited, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, limited.Close()) }()
	content, err := limited.ReadFile("README.txt")
	if err != nil {
		return err
	}
	fmt.Printf("通过 Root 读取内部文件：%q\n", content)
	_, outsideErr := limited.ReadFile("../outside.txt")
	fmt.Println("Root 拒绝父目录路径：", outsideErr != nil)
	// Root 不是完整进程沙箱：不阻止挂载点等文件系统边界，不限制内容大小，也不替代业务鉴权。
	// js/plan9 等平台存在额外限制；需要跨平台部署时查 go doc os.Root。

	// ═══════════════════════════════════════════════════════
	// 第 7 节：Rename、Remove 与 RemoveAll
	// ═══════════════════════════════════════════════════════
	// Rename 改变名称，Remove 删除单项，RemoveAll 递归删除目录树。
	// 本例只操作自己创建的 staged.txt/final.txt；文件和目录的删除范围要分清。

	staged := filepath.Join(root, "staged.txt")
	final := filepath.Join(root, "final.txt")
	if err := os.WriteFile(staged, []byte("ready"), 0o600); err != nil {
		return err
	}
	if err := os.Rename(staged, final); err != nil {
		return err
	}
	if err := os.Remove(final); err != nil {
		return err
	}
	fmt.Println("重命名后删除本次创建的文件：完成")
	// Remove 删除一个文件或空目录，非空目录会报错；RemoveAll 递归删除整棵目录。
	// RemoveAll 对不存在的路径可返回 nil，不代表此前一定存在并删除了内容。
	// Rename 可能替换已有普通文件；跨文件系统可能失败，非 Unix 平台不承诺原子替换。
	return nil
}

func demoDirectories() (err error) {

	// ═══════════════════════════════════════════════════════
	// 第 1 节：MkdirTemp / MkdirAll：创建目录树
	// ═══════════════════════════════════════════════════════
	// 先创建本次演示的根目录，再用 MkdirAll 补齐 uploads/2026 等层级。
	// RemoveAll 的目标始终是这个临时根目录。

	root, err := os.MkdirTemp("", "go-directories-*")
	if err != nil {
		return err
	}
	// RemoveAll 会递归删除，所以演示只删除本函数自己创建的临时目录。
	defer func() { err = errors.Join(err, os.RemoveAll(root)) }()
	for _, subdir := range []string{filepath.Join("uploads", "2026"), "cache"} {
		// Mkdir 只创建一层；MkdirAll 创建缺少的所有层级，已存在的目录不会报错。
		if err := os.MkdirAll(filepath.Join(root, subdir), 0o700); err != nil {
			return err
		}
	}
	files := []struct {
		path    string
		content string
	}{
		{"README.txt", "demo root\n"},
		{filepath.Join("uploads", "2026", "avatar.txt"), "avatar placeholder\n"},
		{filepath.Join("cache", "temporary.txt"), "skip this subtree\n"},
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(root, file.path), []byte(file.content), 0o600); err != nil {
			return err
		}
	}

	// ═══════════════════════════════════════════════════════
	// 第 2 节：ReadDir：列出当前一层
	// ═══════════════════════════════════════════════════════
	// ReadDir 返回 DirEntry 列表，Name 是这一层的名字，IsDir 判断是否为目录。
	// 它不会自动进入 uploads；要看子目录需要继续读取或递归遍历。

	entries, err := os.ReadDir(root) // 只列出这一层，按文件名排序。
	if err != nil {
		return err
	}
	fmt.Println("根目录内容：")
	for _, entry := range entries {
		fmt.Printf("  %s，目录=%t\n", entry.Name(), entry.IsDir())
		// DirEntry.Info() 可取得大小、权限等，但可能访问文件系统并返回 error。
	}

	// ═══════════════════════════════════════════════════════
	// 第 3 节：WalkDir：递归、错误与跳过子树
	// ═══════════════════════════════════════════════════════
	// WalkDir 对访问到的路径调用回调；先处理 walkErr，再使用 entry。
	// 对 cache 目录返回 SkipDir，才会跳过它下面的全部内容。

	fmt.Println("递归遍历（跳过 cache）：")
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		// 出错时 entry 可能为 nil，必须先处理 walkErr 再使用 entry。
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == "cache" {
			return filepath.SkipDir
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		fmt.Println(" ", filepath.ToSlash(relative)) // 显示时统一成 /，真实路径仍用 filepath。
		return nil
	}); err != nil {
		return err
	}
	// WalkDir 默认不沿符号链接进入目录，可避免遍历时意外循环。
	// Join/Clean 只拼接/规范化路径，不是上传路径安全校验。
	// 若路径来自用户，需要另外验证边界，并考虑符号链接及路径被并发替换的问题。
	return demoPathOperations(root)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := demoDirectories(); err != nil {
		fmt.Fprintln(os.Stderr, "目录演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 去掉 WalkDir 中跳过 cache 的分支，观察多出了哪些路径。
// 解释 Lstat 与 Stat 为什么能对同一个软链接给出不同的文件类型。
