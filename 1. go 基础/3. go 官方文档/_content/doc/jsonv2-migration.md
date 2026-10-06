---
title: 迁移到 encoding/json/v2
layout: article
---

[encoding/json/v2](/pkg/encoding/json/v2) 包于 Go 1.27 引入，是原有 [encoding/json](/pkg/encoding/json) 包的一次重大修订。本指南介绍迁移到 v2 的理由，以及安全迁移的具体方法。

# 为什么要迁移？ {#why-migrate}

首先说明：迁移并非必须！`encoding/json` 包会一直保留，并受 Go 1 兼容性承诺保障。使用 `encoding/json` 的包将能够持续工作。

此外，`encoding/json` 与 `encoding/json/v2` 彼此兼容。例如，如果某个类型通过 [`encoding/json/v2.MarshalerTo`](/pkg/encoding/json/v2#MarshalerTo) 实现了自定义序列化行为，那么调用方通过 [`encoding/json.Marshal`](/pkg/encoding/json#Marshal) 序列化该类型时，仍然会使用这一自定义实现。同样，Go 1.27 随 `encoding/json/v2` 引入的新 `json` 结构体标签也受 `encoding/json` 支持。

虽然不必迁移，但迁移有以下几个充分的理由：

首先，新 API 更易于使用。例如，可以通过 [`encoding/json/v2.MarshalWrite`](/pkg/encoding/json/v2#MarshalWrite) 直接向 [`io.Writer`](/pkg/io#Writer) 写入序列化结果，无需创建 [`encoding/json.Encoder`](/pkg/encoding/json#Encoder)。[`encoding/json/v2.Marshalers`](/pkg/encoding/json/v2#Marshalers) 可以覆盖特定类型的序列化行为，即使这些类型并非由你定义。[`encoding/json/v2.MatchCaseInsenstiveNames`](/pkg/encoding/json/v2#MatchCaseInsenstiveNames) 可以控制 JSON 对象成员名与 Go 结构体字段匹配时是否区分大小写。

> 译注：上文原文将 API 名称拼写为 `MatchCaseInsenstiveNames`。正确名称是 [`MatchCaseInsensitiveNames`](/pkg/encoding/json/v2#MatchCaseInsensitiveNames)。为便于对照，原句及其链接予以保留。

这些改进固然方便，但迁移最主要的理由是：v2 的默认行为更严格，与其他 JSON 实现的互操作性也更好。[`encoding/json` 文档](/pkg/encoding/json#hdr-Migrating_to_v2) 列出了完整差异，其中包括：

* v1 会将字符串中无效的 UTF-8 字节静默替换为 Unicode 替换字符；v2 遇到无效 UTF-8 时会返回错误。
* v1 允许 JSON 对象中出现重复名称；v2 则会因此返回错误。
* v1 将 nil Go 切片或 map 序列化为 JSON null；v2 则分别序列化为空 JSON 数组或空 JSON 对象。
* v1 不会在运行时报告 Go 结构体类型的某些结构性错误，例如格式错误的字段标签；v2 则会在运行时报告影响 JSON 序列化的 Go 类型错误。

这些变化旨在让 `encoding/json/v2` 更好地融入整个 JSON 生态，减少意外行为和出错机会。不过，它们并不向后兼容：一些应用程序可能依赖 v1 的行为。因此，迁移到 v2 时必须仔细测试，确保兼容性。

# API 变化与行为变化 {#api-vs-behavior-changes}

对于简单的序列化和反序列化，Go API 在语言层面基本兼容，所以迁移 API 很容易。例如，只需将 `import "encoding/json"` 改为 `import "encoding/json/v2"`，`b, err := json.Marshal(v)` 就仍然可以编译。

从 v1 迁移到 v2 的难点在于序列化和反序列化的行为变化，而非 API 差异。请看下面的程序：

```
package main

import (
	"encoding/json"
	"fmt"
)

type Pet struct {
	Name      string
	Nicknames []string
}

func main() {
	pets := []Pet{
		{Name: "Oliver", Nicknames: []string{"Ollie", "Olliepop"}},
		{Name: "Remi"},
	}
	b, err := json.Marshal(pets)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```

[运行](/play/p/Us887UVmEwm)后输出：

```
[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":null}]
```

如果只修改导入路径，将这个程序迁移到 `encoding/json/v2`，程序仍能编译。[运行](/play/p/E-SprmrVHZF)后输出：

```
[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":[]}]
```

注意，Remi 的 `Nicknames` 字段从 `null` 变成了 `[]`。对于新程序，使用空数组通常是个不错的改进；但对于已有应用程序，使用该输出的下游系统可能依赖 `null`，因此这项变化可能导致它们无法正常工作。

# 选项 {#options}

v2 与 v1 的所有行为差异都有相应的 [`Options`](/pkg/encoding/json/v2#Options)，可以在使用 v2 API 时指定 v1 行为。要使用 v2 API 并保留全部 v1 行为，请使用 [`DefaultOptionsV1`](/pkg/encoding/json#DefaultOptionsV1)：

```
package main

import (
	jsonv1 "encoding/json"
	"encoding/json/v2"
	"fmt"
)

type Pet struct {
	Name      string
	Nicknames []string
}

func main() {
	pets := []Pet{
		{Name: "Oliver", Nicknames: []string{"Ollie", "Olliepop"}},
		{Name: "Remi"},
	}
	b, err := json.Marshal(pets, jsonv1.DefaultOptionsV1())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```

[运行](/play/p/8wg09vDeNN6)后，又会输出 `null`：

```
[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":null}]
```

[DefaultOptionsV1](/pkg/encoding/json#DefaultOptionsV1) 文档列出了用于保持 v1 兼容性的全部选项。后面的选项会覆盖前面的选项，因此可以参考这份列表，逐项启用 v2 行为。

# 迁移 {#migration}

根据应用程序类型及其对风险的容忍程度，可以采用不同的 v2 迁移方式：

* 风险容忍度高：[一次性迁移](#all-at-once)
* 风险容忍度低，或需要排查问题：[逐项迁移](#option-by-option)
* 生产环境中的服务端应用：[jsonsplit](#jsonsplit)

## 一次性迁移 {#all-at-once}

如果应用程序很简单，或能容忍较高风险（迁移出现问题的影响较小），就不一定需要复杂的迁移流程。直接更新调用位置，改用 `encoding/json/v2`，确认测试通过后提交即可。

如果遇到兼容性问题，可能从输出变化或错误信息就能看出原因，此时可以设置[相应的兼容性选项](/pkg/encoding/json#DefaultOptionsV1)。如果原因不明确，可以使用下面的方法辅助排查。

也可以先用这种方式快速找出并修复明显问题（例如单元测试能发现的问题），再用更细致的方法定位剩余问题。

## 逐项迁移 {#option-by-option}

如果应用程序复杂，或对风险的容忍度较低，就可能需要更缓慢、更谨慎地推进迁移。

如前所述，调用 `Marshal` 或 `Unmarshal` 时传入 [DefaultOptionsV1](/pkg/encoding/json#DefaultOptionsV1)，其行为就与 `encoding/json` 一致。第一步可以将所有调用迁移到 `encoding/json/v2`，并传入 `DefaultOptionsV1`。这项修改简单而安全；事实上，[`encoding/json` 就是这样实现 `Marshal` 和 `Unmarshal` 的](https://cs.opensource.google/go/go/+/refs/tags/go1.27rc2:src/encoding/json/v2_encode.go;l=184-186)！

后传入的选项会覆盖之前的选项，因此可以逐项关闭 [v1 兼容性选项](/pkg/encoding/json#DefaultOptionsV1)。例如，`json.Marshal(v, jsonv1.DefaultOptionsV1(), json.FormatNilSliceAsNull(false))` 的行为与 v1 相同，唯一的区别是将 `nil` 切片编码为空数组。

这样就可以逐步迁移，而不必一次改变全部行为。可以每次启用一个选项，以准确掌握变化，也可以将相似选项分组启用。排查问题时，还可以借助这些选项二分定位究竟是哪项行为变化导致不兼容；这与下文 `jsonsplit` 的选项检测机制类似。

## jsonsplit {#jsonsplit}

[`github.com/go-json-experiment/jsonsplit`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit) 是一个 JSON 包装库，可以在生产环境中报告 v1 与 v2 的差异，辅助完成迁移。

`jsonsplit` 提供 `Marshal` 和 `Unmarshal` 包装函数，默认行为与 `encoding/json` 相同，但可以在运行时配置为使用 v1、v2 或同时使用两者。

配置为同时使用两者（[`CallBothButReturnV1`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode)）时，`jsonsplit.Marshal` 会将输入序列化两次，分别使用 v1 和 v2。它会报告发现的差异，但仍向调用方返回 v1 的结果。这使生产服务能够在保持现有行为的同时发现差异。启用可选的 [`AutoDetectOptions`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#Codec) 后，`jsonsplit` 甚至可以自动找出导致差异的具体选项。

这些功能有一定开销。分别使用 v1 和 v2 序列化以检测差异，序列化成本大约会翻倍；`AutoDetectOptions` 为了缩小相关选项的范围，还会进行更多次序列化。为降低成本，`jsonsplit` 可以通过 [`SetMarshalCallRatio`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#Codec.SetMarshalCallRatio) 随机抽取部分调用进行检查。

下面把 `jsonsplit` 应用到最初的示例中：

```
package main

import (
	"fmt"

	"github.com/go-json-experiment/jsonsplit"
)

func init() {
	// Call both v1 and v2 so we can detect differences, but continue using
	// v1 output.
	jsonsplit.GlobalCodec.SetMarshalCallMode(jsonsplit.CallBothButReturnV1)

	// Specify that when a difference is detected, to auto-detect which
	// options are causing the difference.
	jsonsplit.GlobalCodec.AutoDetectOptions = true

	// Log every time we detect a difference between v1 and v2.
	jsonsplit.GlobalCodec.ReportDifference = func(d jsonsplit.Difference) {
		fmt.Printf("detected jsonv1-to-jsonv2 difference: %v\n", d)
	}
}

type Pet struct {
	Name      string
	Nicknames []string
}

func main() {
	pets := []Pet{
		{Name: "Oliver", Nicknames: []string{"Ollie", "Olliepop"}},
		{Name: "Remi"},
	}
	b, err := jsonsplit.Marshal(pets)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
```

[运行](/play/p/wagp5v8V1w-)后，它会报告差异，甚至指出差异由 `FormatNilSliceAsNull` 引起：

```
detected jsonv1-to-jsonv2 difference: {"Caller":"main.main+5","Func":"Marshal","GoType":"[]main.Pet","JSONValueV1":[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":null}],"JSONValueV2":[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":[]}],"Options":["jsonv2.FormatNilSliceAsNull"]}
[{"Name":"Oliver","Nicknames":["Ollie","Olliepop"]},{"Name":"Remi","Nicknames":null}]
```

可以按以下流程平稳迁移生产服务：

1. 将调用位置改为使用 `jsonsplit`，设置 [`CallBothButReturnV1`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode) 和可选的 [`AutoDetectOptions`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode)，并为 [`ReportDifference`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode) 接入所需的监控方式，例如日志或上报指标。

2. 监控生产环境中报告的差异。

3. 将差异处理落实到代码中。

针对 `jsonsplit` 报告的差异，调整选项或类型，使输出保持一致。

例如，上例可以传入 `json.FormatNilSliceAsNull(true)` 选项。有时，v2 报告的问题也可以直接修复。例如，将 JSON 结构体字段标签的 `"string"` 选项用于不适用的类型（如结构体）时，v1 会忽略它，而 v2 会报错。虽然 [`ReportErrorsWithLegacySemantics`](/pkg/encoding/json#ReportErrorsWithLegacySemantics) 可以抑制该错误，但删除 `"string"` 标签更合理：它原本就没有发挥作用。

输出不同不一定会破坏下游行为，但意味着存在这种可能。先通过调整选项保持兼容，可以完成迁移的大部分工作，无需在每处细微变化前停下来评估。切换到 v2 后，应重新检查这些位置，判断是否可以改用新行为。

4. 切换到 v2。

当生产环境不再报告新的差异时，可以设置 [`OnlyCallV2`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode) 来迁移到 v2 行为；也可以设置 [`CallBothButReturnV2`](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#CallMode)，继续检查差异。

5. 清理迁移代码。

安全部署 v2 切换后，将 `jsonsplit` 改为直接使用 `encoding/json/v2`，清理迁移阶段的代码。此时可以评估第 3 步中保留 v1 行为的情况，判断是否能安全切换到 v2 行为。

关于使用 `jsonsplit` 迁移的更多细节，请参阅 [`jsonsplit` 文档](https://pkg.go.dev/github.com/go-json-experiment/jsonsplit#hdr-Example_usage_and_migration)。