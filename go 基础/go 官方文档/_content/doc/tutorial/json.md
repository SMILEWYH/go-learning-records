<!--{
  "Title": "教程：处理 JSON",
  "Breadcrumb": true
}-->

JSON（JavaScript Object Notation，JavaScript 对象表示法）是一种简单的数据交换格式。它的语法类似 JavaScript 中的对象和列表，常用于与网络 API 服务通信，也广泛用于其他场景。[json.org](http://json.org) 网站提供了清晰、简洁的标准定义。

借助 [encoding/json/v2 包](/pkg/encoding/json/v2)，可以轻松地在 Go 程序中读写 JSON 数据。相较于较早的 [encoding/json](/pkg/encoding/json) 包，它提供了更简洁的 API 和更合理的默认行为。

## 编码 {#encoding}

使用 [`Marshal`](/pkg/encoding/json/v2#Marshal) 函数可以将数据编码为 JSON。

```
func Marshal(in any, opts ...Options) (out []byte, err error)
```

假设有如下 Go 数据结构 `Message`：

```
type Message struct {
	Name string
	Body string
	Time time.Time
}
```

以及一个 `Message` 实例：

```
m := Message{"Alice", "Hello", time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC)}
```

使用 `json.Marshal`，可以将 m 序列化为 JSON 编码形式：

```
b, err := json.Marshal(m)
```

如果一切正常，`err` 为 `nil`，`b` 则是包含以下 JSON 数据的 `[]byte`：

```
b == []byte(`{"Name":"Alice","Body":"Hello","Time":"2011-01-25T00:00:00Z"}`)
```

只有能表示为有效 JSON 的数据结构才可以被编码：

  - 指针会被编码为其指向的值；如果指针为 `nil`，则编码为 `null`。

  - JSON 对象只支持字符串键。因此，要对 Go 映射（map）进行编码，其键必须能编码为 JSON 字符串，例如 `map[string]T`，其中 `T` 可以是 json 包支持的任意 Go 类型。

  - 通道（channel）、复数和函数类型不能被编码。

  - 不支持包含循环引用的数据结构。

  - 完整的编码语义请参阅 [`Marshal`](/pkg/encoding/json/v2#Marshal) 的文档。

json 包只访问结构体类型的导出字段，也就是名称以大写字母开头的字段。因此，JSON 输出中只会包含结构体的导出字段。

## 解码 {#decoding}

使用 [`Unmarshal`](/pkg/encoding/json/v2#Unmarshal) 函数可以解码 JSON 数据。

```
func Unmarshal(in []byte, out any, opts ...Options) (err error)
```

首先创建用于存储解码结果的变量：

```
var m Message
```

然后调用 `json.Unmarshal`，传入包含 JSON 数据的 `[]byte` 和指向 `m` 的指针：

```
err := json.Unmarshal(b, &m)
```

如果 `b` 包含可以存入 `m` 的有效 JSON，调用后 `err` 为 `nil`，`b` 中的数据会存入结构体 `m`，效果相当于以下赋值：

```
m = Message{
	Name: "Alice",
	Body: "Hello",
	Time: time.Date(2011, 1, 25, 0, 0, 0, 0, time.UTC),
}
```

`Unmarshal` 如何确定解码后的数据应该存入哪个字段？对于 JSON 键 `"Foo"`，`Unmarshal` 会按以下优先顺序查找目标结构体的字段：

  - 带有 `json:"Foo"` 标签的导出字段。有关结构体标签的更多说明，请参阅 [Go 语言规范](/ref/spec#Struct_types)。

  - 名为 `Foo` 的导出字段。

如果 JSON 数据的结构与 Go 类型不完全匹配，会发生什么？

```
b := []byte(`{"Name":"Bob","Food":"Pickle"}`)
var m Message
err := json.Unmarshal(b, &m)
```

`Unmarshal` 只会解码目标类型中能找到的字段。在这个例子中，只会填充 `m` 的 `Name` 字段，而 `Food` 字段会被忽略。当你只想从大量 JSON 数据中提取少数字段时，这种行为非常实用。这也意味着，目标结构体中的非导出字段不会受到 `Unmarshal` 的影响。

如果事先不知道 JSON 数据的结构，又该怎么办？

## 使用 `any` 处理通用 JSON {#generic-json-with-any}

`any` 类型可以表示任意 Go 类型。`any` 是 `interface{}`（空接口）的别名，表示不包含任何方法的接口。所有 Go 类型都满足没有方法要求的空接口。

`any` 可以作为通用的容器类型：

	var a any
	a = "a string"
	a = 2011
	a = 2.777

通过类型断言，可以访问其底层的具体类型：

	r := a.(float64)
	fmt.Println("the circle's area", math.Pi*r*r)

如果不知道底层类型，可以用类型选择（type switch）进行判断：

	switch v := a.(type) {
	case int:
	    fmt.Println("twice a is", v*2)
	case float64:
	    fmt.Println("the reciprocal of a is", 1/v)
	case string:
	    h := len(v) / 2
	    fmt.Println("a swapped by halves is", v[h:]+v[:h])
	default:
	    // a isn't one of the types above
	}

json 包使用 `map[string]any` 和 `[]any` 值存储任意 JSON 对象和数组，能够将任意有效的 JSON 数据反序列化到一个普通的 `any` 值中。默认使用的具体 Go 类型如下：

  - JSON 布尔值对应 `bool`。

  - JSON 数字对应 `float64`。

  - JSON 字符串对应 `string`。

  - JSON null 对应 `nil`。

## 解码任意数据 {#decoding-arbitrary-data}

考虑存储在变量 `b` 中的以下 JSON 数据：

```
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
```

即使不知道数据结构，也可以使用 `Unmarshal` 将它解码到 `any` 值中：

```
var f any
err := json.Unmarshal(b, &f)
```

此时，`f` 中的 Go 值是一个映射，其键为字符串，值则分别存储在 `any` 中：

```
f = map[string]any{
	"Name": "Wednesday",
	"Age":  6,
	"Parents": []any{
		"Gomez",
		"Morticia",
	},
}
```

要访问这些数据，可以使用类型断言取得 `f` 底层的 `map[string]any`：

```
m := f.(map[string]any)
```

随后可以使用 range 语句遍历映射，再用类型选择，按各个值的具体类型访问它们：

```
for k, v := range m {
	switch vv := v.(type) {
	case string:
		fmt.Println(k, "is string", vv)
	case float64:
		fmt.Println(k, "is float64", vv)
	case []any:
		fmt.Println(k, "is an array:")
		for i, u := range vv {
			fmt.Println(i, u)
		}
	default:
		fmt.Println(k, "is of a type I don't know how to handle")
	}
}
```

通过这种方式，既能处理结构未知的 JSON 数据，又能保留类型安全的优势。

## 引用类型 {#reference-types}

定义一个 Go 类型，用来存储上一示例中的数据：

```
type FamilyMember struct {
	Name    string
	Age     int
	Parents []string
}

var m FamilyMember
err := json.Unmarshal(b, &m)
```

将数据反序列化到 `FamilyMember` 值中会得到预期的结果，但仔细观察可以发现一个值得注意的细节。我们通过 var 语句创建了 `FamilyMember` 结构体，再将指向该值的指针传给 `Unmarshal`；此时 `Parents` 字段还是一个值为 `nil` 的切片。为了填充 `Parents` 字段，`Unmarshal` 在内部创建了一个新切片。这就是 `Unmarshal` 处理其支持的引用类型（指针、切片和映射）时的典型行为。

再考虑反序列化到以下数据结构的情况：

```
type Foo struct {
	Bar *Bar
}
```

如果 JSON 对象中存在 `Bar` 字段，`Unmarshal` 就会创建一个新的 `Bar` 值并填充数据；否则，`Bar` 会保持为 `nil` 指针。

由此可以得到一种实用模式：如果应用程序会接收几种不同类型的消息，可以定义如下“接收用”结构体：

```
type IncomingMessage struct {
	Cmd *Command
	Msg *Message
}
```

发送方可以根据要传递的消息类型，填充顶层 JSON 对象中的 `Cmd` 字段、`Msg` 字段，或者同时填充两者。`Unmarshal` 将 JSON 解码到 `IncomingMessage` 结构体时，只会为 JSON 数据中存在的部分分配相应的数据结构。要确定需要处理哪些消息，只需检查 `Cmd` 或 `Msg` 是否非 `nil`。

## 流式序列化与反序列化 {#streaming-marshal-and-unmarshal}

[`io.Reader`](/pkg/io#Reader) 和 [`io.Writer`](/pkg/io#Writer) 接口在 Go 中十分常见，用于以流的形式访问 HTTP 连接、WebSocket、文件等资源。`MarshalWrite` 和 `UnmarshalRead` 函数可以直接向这些流写入序列化结果，或从流中读取并反序列化，无需创建存放完整消息的中间 `[]byte`。

```
func MarshalWrite(out io.Writer, in any, opts ...Options) (err error)
func UnmarshalRead(in io.Reader, out any, opts ...Options) (err error)
```

例如，可以直接将序列化结果写入标准输出：

```
err := json.MarshalWrite(os.Stdout, m)
```

## 自定义序列化与反序列化 {#custom-marshal-and-unmarshal}

有时，默认的序列化行为并不适合你的类型。

例如，假设有一个描述软件版本的类型：

```
type Version struct {
	Major, Minor, Patch int64
}

v := Version{1, 2, 3}
b, err := json.Marshal(s)
```

> 译注：上游示例最后一行的 `json.Marshal(s)` 中，`s` 未定义。实际运行时，应使用前面声明的 `v`，即调用 `json.Marshal(v)`。

根据结构体定义，它会被序列化为 `{"Major":1,"Minor":2,"Patch":3}`，符合默认行为。不过，对 JSON 数据的使用方来说，"1.2.3" 这样的版本字符串往往更合适。

通过实现 [`MarshalerTo`](/pkg/encoding/json/v2#MarshalerTo) 接口，可以自定义 JSON 表示形式。[`UnmarshalerFrom`](/pkg/encoding/json/v2#UnmarshalerFrom) 则用于将这种自定义 JSON 表示反序列化。

```
type MarshalerTo interface {
	MarshalJSONTo(*jsontext.Encoder) error
}

type UnmarshalerFrom interface {
	UnmarshalJSONFrom(*jsontext.Decoder) error
}
```

可以实现这些方法，完成与字符串表示之间的双向转换。

```
func (v Version) MarshalJSONTo(enc *jsontext.Encoder) error {
	return json.MarshalEncode(enc, fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
}
```

有了自定义序列化逻辑，`Version{1, 2, 3}` 现在会被序列化为 `"1.2.3"`。

```
func (v *Version) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != jsontext.KindString {
		// Value must be a string.
		return &json.SemanticError{JSONKind: k}
	}

	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	_, err := fmt.Sscanf(s, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return err
}
```

有了自定义反序列化逻辑，`"1.2.3"` 现在会被反序列化为 `Version{1, 2, 3}`。这个简单示例要求版本严格符合 `"major.minor.patch"` 格式；更复杂的 `UnmarshalJSONFrom` 实现可以提供更灵活的处理方式，例如允许省略次版本号和补丁版本号。

## 参考资料 {#references}

更多信息请参阅 [encoding/json/v2](/pkg/encoding/json/v2) 和 [encoding/json/jsontext](/pkg/encoding/json/jsontext) 的包文档。
