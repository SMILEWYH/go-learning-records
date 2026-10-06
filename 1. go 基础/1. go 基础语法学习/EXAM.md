# 独立章节与自动检查练习

每章包含三个 Go 文件：

| 文件 | 内容 |
| --- | --- |
| `lesson.go` | 中文讲解与可运行示例；第 30 章末尾含部署命令附录。 |
| `exercises.go` | 三道练习与待填写函数；末尾是启动本章检查的简短入口。 |
| `exercises_test.go` | 题目对应关系、全部检查逻辑、练习与讲解测试，以及末尾注释中的参考答案。 |

各章自带 `go.mod`，第 29 章另有第三方依赖校验文件 `go.sum`。这些是 Go 配置文件。章节不导入课程根目录或其他章节，可以将整个章节目录复制到课程之外，运行 `go run .` 和 `go run . exam`。第 29 章首次运行可能需要下载 WebSocket 依赖。

在 `go 基础语法学习` 目录运行：

```sh
go run ./01-variables exam
go run ./05-if exam
```

也可以进入某章目录后运行 `go run . exam`。不带 `exam` 仍运行讲解。

每道题会显示实际受测的函数名和一个结果：

```text
练习 1（normalizePageSize）：还未完成
练习 2（accessStatus）：错误
练习 3（resolveTitle）：完成
汇总：完成 1 题，错误 1 题，还未完成 1 题
```

| 结果 | 判断方式 |
| --- | --- |
| 还未完成 | 函数体为空、只有注释，或只保留原来的“待完成”打印 / panic。此题不运行测试。 |
| 错误 | 已填写实现，但对应测试失败、panic、超时，或验收测试缺失、为空、被跳过。会附上错误原因。 |
| 完成 | 对应测试实际执行并通过。 |

函数里仍有 `TODO` 注释不会阻止检查。只要已经写入实际代码，就交给测试判定；无效实现、普通 panic 和错误返回值都不会被当作空函数。

## 答案写在哪里

- **01–04 章**：直接填写 `exercises.go` 中的 `exercise1/2/3`，输出顺序见 `exercises_test.go`。
- **05–30 章的业务函数题**：填写 `exercises.go` 中“自动测试入口”下的函数或方法，例如 `normalizePageSize`。`exercise1/2/3` 保留题目和演示入口，不用为了检查再填一份答案。
- **一题包含多个函数**：全部填写后运行该题测试，按整题给出结果；输出会列出仍未填写的函数。
- **26 章第 3 题、28/29 章第 2 题、30 章第 3 题**：实际实现就是 `exercises_test.go` 中对应的 `TestExerciseN`。原文件已提供可运行实现，所以这些题可能直接显示完成，表示现有测试通过。
- **30 章第 2 题**：检查 `lesson.go` 中的 `runDeployment` 和 `buildTime`，测试实际构建临时产物并核对版本输出。

打开本章 `exercises_test.go`，顶部的 `TestExam` 记录三道题与函数的对应关系。检查器在同一文件内，参考答案集中在文件末尾。原来分开的讲解测试已并入该文件，函数名与测试行为保留。

## 检查范围

`go run . exam` 通过本章 `runExercises` 启动 `TestExam`。普通 `go test` 不运行这个汇总入口，仍直接执行指定的练习测试。

检查器逐题执行 `go test -json -count=1 -run '^TestExerciseN$' -timeout 30s .`，编译整个章节，运行 `exercises_test.go` 中对应的测试及其子测试。一道题失败不会阻止检查后面的题。单次命令另有 90 秒上限，包含编译时间。

`exam` 输出检查报告；需要用退出码判断通过与否时，仍可直接运行：

```sh
go test ./05-if -run '^TestExercise' -count=1 -timeout 30s
go test -race ./17-concurrency-safety -run '^TestExercise' -count=1 -timeout 30s
```

`go run` 必须先编译成功，才能进入检查器。语法错误、未使用的 import、带返回值函数直接留空造成的 `missing return`，会先显示 Go 编译错误。暂未填写的有返回值函数请保留原来的占位 `panic`，这样可以编译并显示“还未完成”。

“完成”表示通过当前测试覆盖的行为；题目要求的语法、思考题及未覆盖的情况仍需结合代码核对。
