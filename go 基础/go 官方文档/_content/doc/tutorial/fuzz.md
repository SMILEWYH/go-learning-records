<!--{
  "Template": true,
  "Title": "教程：模糊测试入门",
  "HideTOC": true,
  "Breadcrumb": true
}-->

本教程介绍 Go 模糊测试的基础知识。模糊测试会使用随机数据运行测试，尝试发现漏洞或导致程序崩溃的输入。它可以帮助发现的问题包括 SQL 注入、缓冲区溢出、拒绝服务和跨站脚本攻击等。

在本教程中，你将为一个简单函数编写模糊测试，运行 go 命令，并调试和修复代码中的问题。

有关本教程涉及的术语，请参阅 [Go 模糊测试术语表](/security/fuzz/#glossary)。

你将依次完成以下内容：

1. [为代码创建文件夹。](#create_folder)
2. [添加待测试的代码。](#code_to_test)
3. [添加单元测试。](#unit_test)
4. [添加模糊测试。](#fuzz_test)
5. [修复两个缺陷。](#fix_invalid_string_error)
6. [探索更多资源。](#conclusion)

**注意：**其他教程请参阅[教程列表](/doc/tutorial/index.html)。

**注意：**Go 模糊测试目前支持部分内置类型，具体列表见 [Go 模糊测试文档](/security/fuzz/#requirements)，未来将增加对更多内置类型的支持。

## 准备工作 {#prerequisites}

- **Go。**建议使用最新版本的 Go 学习本教程。安装步骤请参阅[安装 Go](/doc/install)。
- **代码编辑工具。**任何文本编辑器都可以。
- **命令行终端。**Go 可以在 Linux 和 Mac 的各种终端，以及 Windows 的 PowerShell 或 cmd 中正常使用。
- **支持模糊测试的环境。**目前，带有覆盖率插桩的 Go 模糊测试仅支持 AMD64 和 ARM64 架构。

## 为代码创建文件夹 {#create_folder}

首先，为即将编写的代码创建一个文件夹。

1. 打开命令行终端，切换到用户主目录。

   在 Linux 或 Mac 上：

   ```
   $ cd
   ```

   在 Windows 上：

   ```
   C:\> cd %HOMEPATH%
   ```

   本教程后续统一使用 $ 表示命令提示符。这些命令同样适用于 Windows。

2. 在命令行中创建名为 fuzz 的目录，用来存放代码。

   ```
   $ mkdir fuzz
   $ cd fuzz
   ```

3. 创建模块，组织代码。

   运行 `go mod init` 命令，并传入新代码的模块路径。

   ```
   $ go mod init example/fuzz
   go: creating new go.mod: module example/fuzz
   ```

   **注意：**对于生产环境的代码，应根据实际需求指定更合适的模块路径。更多说明请参阅[依赖管理](/doc/modules/managing-dependencies)。

接下来，你将添加一个简单的字符串反转函数，稍后对它进行模糊测试。

## 添加待测试的代码 {#code_to_test}

这一步将添加一个反转字符串的函数。

### 编写代码 {#write-the-code}

1.  使用文本编辑器，在 fuzz 目录中创建 main.go 文件。
2.  将以下包声明粘贴到 main.go 文件顶部。

    ```
    package main
    ```

    独立运行的程序使用 `main` 包，而库则使用其他包名。

3.  在包声明下方粘贴以下函数声明。

    ```
    func Reverse(s string) string {
        b := []byte(s)
        for i, j := 0, len(b)-1; i {{raw "<"}} len(b)/2; i, j = i+1, j-1 {
            b[i], b[j] = b[j], b[i]
        }
        return string(b)
    }
    ```

    这个函数接收一个 `string`，按 `byte` 逐字节遍历，最后返回反转后的字符串。

    _注意：_这段代码基于 golang.org/x/example 中的 `stringutil.Reverse` 函数。

4.  在 main.go 顶部的包声明下方，粘贴以下 `main` 函数。它会初始化一个字符串，对其反转并输出结果，然后再反转一次。

    ```
    func main() {
        input := "The quick brown fox jumped over the lazy dog"
        rev := Reverse(input)
        doubleRev := Reverse(rev)
        fmt.Printf("original: %q\n", input)
        fmt.Printf("reversed: %q\n", rev)
        fmt.Printf("reversed again: %q\n", doubleRev)
    }
    ```

    这个函数调用几次 `Reverse`，再将结果输出到命令行。这便于观察代码的行为，也有助于调试。

5.  `main` 函数使用了 fmt 包，因此需要将它导入。

    代码开头应如下所示：

    ```
    package main

    import "fmt"
    ```

### 运行代码 {#run-the-code}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
original: "The quick brown fox jumped over the lazy dog"
reversed: "god yzal eht revo depmuj xof nworb kciuq ehT"
reversed again: "The quick brown fox jumped over the lazy dog"
```

输出依次展示原始字符串、反转后的结果，以及再次反转后的结果；最后一个结果与原始字符串相同。

现在代码能够运行，可以开始测试了。

## 添加单元测试 {#unit_test}

这一步将为 `Reverse` 函数编写一个基本的单元测试。

### 编写代码 {#write-the-code-1}

1. 使用文本编辑器，在 fuzz 目录中创建 reverse_test.go 文件。
2. 将以下代码粘贴到 reverse_test.go 中。

   ```
   package main

   import (
       "testing"
   )

   func TestReverse(t *testing.T) {
       testcases := []struct {
           in, want string
       }{
           {"Hello, world", "dlrow ,olleH"},
           {" ", " "},
           {"!12345", "54321!"},
       }
       for _, tc := range testcases {
           rev := Reverse(tc.in)
           if rev != tc.want {
                   t.Errorf("Reverse: %q, want %q", rev, tc.want)
           }
       }
   }
   ```

   这个简单测试会断言：列出的各个输入字符串都能被正确反转。

### 运行代码 {#run-the-code-1}

使用 `go test` 运行单元测试。

```
$ go test
PASS
ok      example/fuzz  0.013s
```

接下来，你将把这个单元测试改为模糊测试。

## 添加模糊测试 {#fuzz_test}

单元测试存在局限：每个输入都需要开发者手动加入测试。模糊测试的一个优点是能够自动为代码生成输入，从而发现手写测试用例没有覆盖的边界情况。

本节将把单元测试改为模糊测试，让你用更少的工作量获得更多测试输入。

注意，单元测试、基准测试和模糊测试可以放在同一个 *_test.go 文件中。本例为了演示，将直接把单元测试转换为模糊测试。

### 编写代码 {#write-the-code-2}

在文本编辑器中，用以下模糊测试替换 reverse_test.go 中的单元测试。

```
func FuzzReverse(f *testing.F) {
    testcases := []string{"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc)  // Use f.Add to provide a seed corpus
    }
    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)
        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
```

模糊测试也有一定局限。在单元测试中，可以预先确定 `Reverse` 函数的期望输出，再检查实际输出是否符合预期。

例如，单元测试中的 `Reverse("Hello, world")` 用例明确规定返回值应为 `"dlrow ,olleH"`。

进行模糊测试时，由于无法控制所有生成的输入，就无法逐一预先给出期望输出。

不过，仍可以在模糊测试中验证 `Reverse` 函数应满足的一些性质。本例检查两项性质：

1.  将字符串反转两次，应保持原值不变。
2.  如果原字符串是有效的 UTF-8，反转后的字符串也应保持有效。

注意单元测试与模糊测试在语法上的区别：

- 函数以 FuzzXxx 命名，而不是 TestXxx；接收 `*testing.F` 参数，而不是 `*testing.T`。
- 原本可能调用 `t.Run` 的位置，现在使用 `f.Fuzz`。它接收一个模糊测试目标函数，该函数的参数包括 `*testing.T` 和要进行模糊测试的输入类型。原单元测试中的输入通过 `f.Add` 加入种子语料库。

确保已经导入新增的 `unicode/utf8` 包。

```
package main

import (
    "testing"
    "unicode/utf8"
)
```

单元测试已经转换为模糊测试，现在重新运行测试。

### 运行代码 {#run-the-code-2}

1. 先运行模糊测试函数，但不启用随机输入生成，确认种子输入都能通过。

   ```
   $ go test
   PASS
   ok      example/fuzz  0.013s
   ```

   如果文件中还有其他测试，而你只想运行这个模糊测试，也可以使用 `go test -run=FuzzReverse`。

2. 为 `FuzzReverse` 启用模糊测试，检查随机生成的字符串是否会导致失败。在 `go test` 后添加 `-fuzz` 标志，并将其设置为 `Fuzz`。复制以下命令。

    ```
    $ go test -fuzz=Fuzz
    ```

    另一个实用的标志是 `-fuzztime`，用于限制模糊测试的时长。例如，下面的测试若指定 `-fuzztime 10s`，且期间没有提前出现失败，就会在运行 10 秒后退出。其他测试标志请参阅 cmd/go 文档中的[相关章节](https://pkg.go.dev/cmd/go#hdr-Testing_flags)。

   现在运行刚复制的命令。

   ```
   $ go test -fuzz=Fuzz
   fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
   fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
   fuzz: minimizing 38-byte failing input file...
   --- FAIL: FuzzReverse (0.01s)
       --- FAIL: FuzzReverse (0.00s)
           reverse_test.go:20: Reverse produced invalid UTF-8 string "\x9c\xdd"

       Failing input written to testdata/fuzz/FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a
       To re-run:
       go test -run=FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a
   FAIL
   exit status 1
   FAIL    example/fuzz  0.030s
   ```

   模糊测试发现了一次失败，导致问题的输入被写入种子语料库文件。下次运行 `go test` 时，即使不加 `-fuzz` 标志，也会运行这个输入。要查看导致失败的输入，请使用文本编辑器打开 testdata/fuzz/FuzzReverse 目录中的语料文件。你得到的文件可能包含不同的字符串，但格式相同。

   ```
   go test fuzz v1
   string("泃")
   ```

   语料文件的第一行表示编码版本，后面的各行分别表示构成这个语料条目的各个类型的值。由于本例的模糊测试目标只接收一个输入，因此版本行之后只有一个值。

3. 再次运行 `go test`，不加 `-fuzz` 标志。新加入的、导致失败的种子语料条目也会参与测试：

   ```
   $ go test
   --- FAIL: FuzzReverse (0.00s)
       --- FAIL: FuzzReverse/af69258a12129d6cbba438df5d5f25ba0ec050461c116f777e77ea7c9a0d217a (0.00s)
           reverse_test.go:20: Reverse produced invalid string
   FAIL
   exit status 1
   FAIL    example/fuzz  0.016s
   ```

   测试失败了，接下来需要调试。

## 修复无效字符串错误 {#fix_invalid_string_error}

本节将分析测试失败的原因，并修复缺陷。

继续阅读之前，你也可以先思考一下，尝试自行解决这个问题。

### 诊断错误 {#diagnose-the-error}

有多种方式可以调试这个错误。如果使用 VS Code 作为文本编辑器，可以[配置调试器](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)进行调查。

本教程将通过向终端输出有用的调试信息来定位问题。

首先看看 [`utf8.ValidString`](https://pkg.go.dev/unicode/utf8) 的文档：

```
ValidString reports whether s consists entirely of valid UTF-8-encoded runes.
```

这段说明的意思是：ValidString 判断 s 是否完全由有效 UTF-8 编码的 rune 组成。

当前的 `Reverse` 函数按字节反转字符串，这正是问题所在。要保留原始字符串中以 UTF-8 编码的 rune（Unicode 码点），需要按 rune 反转字符串。

为了分析为什么输入（本例中是汉字 `泃`）会导致 `Reverse` 返回无效字符串，可以检查反转后字符串中的 rune 数量。

#### 编写代码 {#write-the-code-3}

在文本编辑器中，将 `FuzzReverse` 内的模糊测试目标函数替换为以下内容。

```
f.Fuzz(func(t *testing.T, orig string) {
    rev := Reverse(orig)
    doubleRev := Reverse(rev)
    t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))
    if orig != doubleRev {
        t.Errorf("Before: %q, after: %q", orig, doubleRev)
    }
    if utf8.ValidString(orig) && !utf8.ValidString(rev) {
        t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
    }
})
```

发生错误时，或者使用 `-v` 运行测试时，新增的 `t.Logf` 语句会向命令行输出信息，帮助你调试这个问题。

#### 运行代码 {#run-the-code-3}

使用 go test 运行测试。

```
$ go test
--- FAIL: FuzzReverse (0.00s)
    --- FAIL: FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0 (0.00s)
        reverse_test.go:16: Number of runes: orig=1, rev=3, doubleRev=1
        reverse_test.go:21: Reverse produced invalid UTF-8 string "\x83\xb3\xe6"
FAIL
exit status 1
FAIL    example/fuzz    0.598s
```

最初的种子语料库中，每个字符串里的字符都只占一个字节。但“泃”这样的字符可能需要多个字节表示，因此按字节反转字符串会破坏多字节字符。

**注意：**如果想深入了解 Go 如何处理字符串，请阅读博客文章 [Go 中的字符串、字节、rune 和字符](/blog/strings)。

理解原因后，就可以修复 `Reverse` 函数中的错误了。

### 修复错误 {#fix-the-error}

要修复 `Reverse` 函数，需要按 rune 遍历字符串，而不是按字节遍历。

#### 编写代码 {#write-the-code-4}

在文本编辑器中，用以下代码替换现有的 Reverse() 函数。

```
func Reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i {{raw "<"}} len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

关键区别在于，`Reverse` 现在遍历字符串中的每个 `rune`，而不是每个 `byte`。注意，这只是一个示例，仍不能正确处理[组合字符](https://en.wikipedia.org/wiki/Combining_character)。

#### 运行代码 {#run-the-code-4}

1. 使用 `go test` 运行测试。

   ```
   $ go test
   PASS
   ok      example/fuzz  0.016s
   ```

   这次测试通过了！

2. 再次运行 `go test -fuzz` 进行模糊测试，看看是否还能发现新的缺陷。

   ```
   $ go test -fuzz=Fuzz
   fuzz: elapsed: 0s, gathering baseline coverage: 0/37 completed
   fuzz: minimizing 506-byte failing input file...
   fuzz: elapsed: 0s, gathering baseline coverage: 5/37 completed
   --- FAIL: FuzzReverse (0.02s)
       --- FAIL: FuzzReverse (0.00s)
           reverse_test.go:33: Before: "\x91", after: "�"

       Failing input written to testdata/fuzz/FuzzReverse/1ffc28f7538e29d79fce69fef20ce5ea72648529a9ca10bea392bcff28cd015c
       To re-run:
       go test -run=FuzzReverse/1ffc28f7538e29d79fce69fef20ce5ea72648529a9ca10bea392bcff28cd015c
   FAIL
   exit status 1
   FAIL    example/fuzz  0.032s
   ```

   可以看到，反转两次后的字符串与原值不同。这次输入本身就是无效的 UTF-8。明明测试的是字符串，为什么会出现这种情况？

   再次进行调试。

## 修复两次反转后的错误 {#fix_double_reverse_error}

本节将调试两次反转后未还原原值的问题，并修复缺陷。

继续阅读之前，你也可以先思考一下，尝试自行解决这个问题。

### 诊断错误 {#diagnose-the-error-1}

和之前一样，有多种方法可以调试这次失败。此时，使用[调试器](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)是一个不错的选择。

本教程将在 `Reverse` 函数中输出有用的调试信息。

仔细观察反转后的字符串，就能发现问题。在 Go 中，[字符串可以理解为只读的字节序列](/blog/strings)，其中可以包含无效的 UTF-8 字节。原始字符串只包含一个字节 `'\x91'`。将输入字符串转换为 `[]rune` 时，Go 会按 UTF-8 解码，并将这个无效字节替换为 Unicode 替换字符 �。显然，再将替换字符对应的 UTF-8 字节与原始输入比较时，两者并不相等。

#### 编写代码 {#write-the-code-5}

1. 在文本编辑器中，用以下代码替换 `Reverse` 函数。

   ```
   func Reverse(s string) string {
       fmt.Printf("input: %q\n", s)
       r := []rune(s)
       fmt.Printf("runes: %q\n", r)
       for i, j := 0, len(r)-1; i {{raw "<"}} len(r)/2; i, j = i+1, j-1 {
           r[i], r[j] = r[j], r[i]
       }
       return string(r)
   }
   ```

   这有助于理解将字符串转换为 rune 切片时发生了什么。

#### 运行代码 {#run-the-code-5}

这次只需要运行失败的测试来检查日志，可以使用 `go test -run`。

要单独运行与 FuzzXxx 对应的 testdata 目录中的某个语料条目，可以将 {FuzzTestName}/{filename} 传给 `-run`，这在调试时很有用。本例应将 `-run` 设置为失败测试对应的名称和完整哈希值。请从自己的终端复制粘贴这个唯一哈希值；它会与下面的示例不同。

```
$ go test -run=FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0
input: "\x91"
runes: ['�']
input: "�"
runes: ['�']
--- FAIL: FuzzReverse (0.00s)
    --- FAIL: FuzzReverse/28f36ef487f23e6c7a81ebdaa9feffe2f2b02b4cddaa6252e87f69863046a5e0 (0.00s)
        reverse_test.go:16: Number of runes: orig=1, rev=1, doubleRev=1
        reverse_test.go:18: Before: "\x91", after: "�"
FAIL
exit status 1
FAIL    example/fuzz    0.145s
```

既然已经知道输入不是有效的 UTF-8，就可以修复 `Reverse` 函数了。

### 修复错误 {#fix-the-error-1}

为解决这个问题，当传给 `Reverse` 的输入不是有效 UTF-8 时，让函数返回错误。

#### 编写代码 {#write-the-code-6}

1. 在文本编辑器中，用以下代码替换现有的 `Reverse` 函数。

   ```
   func Reverse(s string) (string, error) {
       if !utf8.ValidString(s) {
           return s, errors.New("input is not valid UTF-8")
       }
       r := []rune(s)
       for i, j := 0, len(r)-1; i {{raw "<"}} len(r)/2; i, j = i+1, j-1 {
           r[i], r[j] = r[j], r[i]
       }
       return string(r), nil
   }
   ```

   修改后，如果输入字符串包含无效的 UTF-8 字节序列，函数就会返回错误。

1. 由于 Reverse 函数现在还会返回错误，需要修改 `main` 函数，接收并显示这个额外的返回值。用以下代码替换现有的 `main` 函数。

   ```
   func main() {
       input := "The quick brown fox jumped over the lazy dog"
       rev, revErr := Reverse(input)
       doubleRev, doubleRevErr := Reverse(rev)
       fmt.Printf("original: %q\n", input)
       fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
       fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
   }
   ```

    输入字符串是有效 UTF-8，因此这些 `Reverse` 调用返回的错误应为 nil。

1. 导入 errors 和 unicode/utf8 包。main.go 中的 import 声明应如下所示。

   ```
   import (
       "errors"
       "fmt"
       "unicode/utf8"
   )
   ```

1. 修改 reverse_test.go 文件，检查返回的错误，并在发生错误时直接返回，跳过对该输入的后续检查。

   ```
   func FuzzReverse(f *testing.F) {
       testcases := []string {"Hello, world", " ", "!12345"}
       for _, tc := range testcases {
           f.Add(tc)  // Use f.Add to provide a seed corpus
       }
       f.Fuzz(func(t *testing.T, orig string) {
           rev, err1 := Reverse(orig)
           if err1 != nil {
               return
           }
           doubleRev, err2 := Reverse(rev)
           if err2 != nil {
                return
           }
           if orig != doubleRev {
               t.Errorf("Before: %q, after: %q", orig, doubleRev)
           }
           if utf8.ValidString(orig) && !utf8.ValidString(rev) {
               t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
           }
       })
   }
   ```

   除了直接返回，也可以调用 `t.Skip()`，终止对该模糊测试输入的执行。

#### 运行代码 {#run-the-code-6}

1. 使用 go test 运行测试。

   ```
   $ go test
   PASS
   ok      example/fuzz  0.019s
   ```

2. 使用 `go test -fuzz=Fuzz` 进行模糊测试，运行几秒后按 `ctrl-C` 停止。除非传入 `-fuzztime` 标志，否则测试会持续运行，直到遇到导致失败的输入。默认情况下，只要没有失败，它就会一直运行；可以用 `ctrl-C` 中断。

   ```
   $ go test -fuzz=Fuzz
   fuzz: elapsed: 0s, gathering baseline coverage: 0/38 completed
   fuzz: elapsed: 0s, gathering baseline coverage: 38/38 completed, now fuzzing with 4 workers
   fuzz: elapsed: 3s, execs: 86342 (28778/sec), new interesting: 2 (total: 35)
   fuzz: elapsed: 6s, execs: 193490 (35714/sec), new interesting: 4 (total: 37)
   fuzz: elapsed: 9s, execs: 304390 (36961/sec), new interesting: 4 (total: 37)
   ...
   fuzz: elapsed: 3m45s, execs: 7246222 (32357/sec), new interesting: 8 (total: 41)
   ^Cfuzz: elapsed: 3m48s, execs: 7335316 (31648/sec), new interesting: 8 (total: 41)
   PASS
   ok      example/fuzz  228.000s
   ```

3. 使用 `go test -fuzz=Fuzz -fuzztime 30s` 进行模糊测试。如果没有发现失败，它会运行 30 秒后退出。

   ```
   $ go test -fuzz=Fuzz -fuzztime 30s
   fuzz: elapsed: 0s, gathering baseline coverage: 0/5 completed
   fuzz: elapsed: 0s, gathering baseline coverage: 5/5 completed, now fuzzing with 4 workers
   fuzz: elapsed: 3s, execs: 80290 (26763/sec), new interesting: 12 (total: 12)
   fuzz: elapsed: 6s, execs: 210803 (43501/sec), new interesting: 14 (total: 14)
   fuzz: elapsed: 9s, execs: 292882 (27360/sec), new interesting: 14 (total: 14)
   fuzz: elapsed: 12s, execs: 371872 (26329/sec), new interesting: 14 (total: 14)
   fuzz: elapsed: 15s, execs: 517169 (48433/sec), new interesting: 15 (total: 15)
   fuzz: elapsed: 18s, execs: 663276 (48699/sec), new interesting: 15 (total: 15)
   fuzz: elapsed: 21s, execs: 771698 (36143/sec), new interesting: 15 (total: 15)
   fuzz: elapsed: 24s, execs: 924768 (50990/sec), new interesting: 16 (total: 16)
   fuzz: elapsed: 27s, execs: 1082025 (52427/sec), new interesting: 17 (total: 17)
   fuzz: elapsed: 30s, execs: 1172817 (30281/sec), new interesting: 17 (total: 17)
   fuzz: elapsed: 31s, execs: 1172817 (0/sec), new interesting: 17 (total: 17)
   PASS
   ok      example/fuzz  31.025s
   ```

   模糊测试通过了！

   除了 `-fuzz`，`go test` 还增加了其他相关标志，详见[文档](/security/fuzz/#custom-settings)。

   模糊测试输出中的术语说明，请参阅 [Go 模糊测试](/security/fuzz/#command-line-output)。例如，“new interesting”表示能够扩大现有语料库代码覆盖范围的新输入。测试刚开始时，这类输入的数量通常会快速增长；发现新的代码路径时，可能再次出现明显增长，随后逐渐趋于平缓。

## 总结 {#conclusion}

做得不错！你已经初步了解了 Go 模糊测试。

下一步，可以从自己的代码中选择一个适合的函数，试着为它添加模糊测试。如果模糊测试发现了缺陷，可以考虑将它加入[成果展示](/wiki/Fuzzing-trophy-case)。

如果遇到问题或有新的功能建议，请[提交 issue](/issue/new/?&labels=fuzz)。

如果想讨论这个功能或提供一般性反馈，也可以加入 Gophers Slack 的 [#fuzzing 频道](https://gophers.slack.com/archives/CH5KV1AKE)。

更多内容请参阅 [go.dev/security/fuzz](/security/fuzz/#requirements) 文档。

## 完整代码 {#completed-code}

--- main.go ---

```
package main

import (
    "errors"
    "fmt"
    "unicode/utf8"
)

func main() {
    input := "The quick brown fox jumped over the lazy dog"
    rev, revErr := Reverse(input)
    doubleRev, doubleRevErr := Reverse(rev)
    fmt.Printf("original: %q\n", input)
    fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
    fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
}

func Reverse(s string) (string, error) {
    if !utf8.ValidString(s) {
        return s, errors.New("input is not valid UTF-8")
    }
    r := []rune(s)
    for i, j := 0, len(r)-1; i {{raw "<"}} len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r), nil
}
```

--- reverse_test.go ---

```
package main

import (
    "testing"
    "unicode/utf8"
)

func FuzzReverse(f *testing.F) {
    testcases := []string{"Hello, world", " ", "!12345"}
    for _, tc := range testcases {
        f.Add(tc) // Use f.Add to provide a seed corpus
    }
    f.Fuzz(func(t *testing.T, orig string) {
        rev, err1 := Reverse(orig)
        if err1 != nil {
            return
        }
        doubleRev, err2 := Reverse(rev)
        if err2 != nil {
            return
        }
        if orig != doubleRev {
            t.Errorf("Before: %q, after: %q", orig, doubleRev)
        }
        if utf8.ValidString(orig) && !utf8.ValidString(rev) {
            t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
        }
    })
}
```

[返回顶部](#top)
