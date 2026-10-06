<!--{
  "Title": "教程：泛型入门",
  "Breadcrumb": true
}-->

本教程介绍 Go 泛型的基础知识。借助泛型，可以声明和使用适用于一组类型的函数或类型，具体类型由调用方代码提供。

本教程将先声明两个简单的非泛型函数，再将相同的逻辑归纳为一个泛型函数。

你将依次完成以下内容：

1. 为代码创建文件夹。
2. 添加非泛型函数。
3. 添加能处理多种类型的泛型函数。
4. 在调用泛型函数时省略类型实参。
5. 声明类型约束。

**注意：**其他教程请参阅[教程列表](/doc/tutorial/index.html)。

## 准备工作 {#prerequisites}

*   **Go。**建议使用最新版本的 Go 学习本教程。安装步骤请参阅[安装 Go](/doc/install)。
*   **代码编辑工具。**任何文本编辑器都可以。
*   **命令行终端。**Go 可以在 Linux 和 Mac 的各种终端，以及 Windows 的 PowerShell 或 cmd 中正常使用。

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

2. 在命令行中创建名为 generics 的目录，用来存放代码。

    ```
    $ mkdir generics
    $ cd generics
    ```

3. 创建模块，组织代码。

    运行 `go mod init` 命令，并传入新代码的模块路径。

    ```
    $ go mod init example/generics
    go: creating new go.mod: module example/generics
    ```

    **注意：**对于生产环境的代码，应根据实际需求指定更合适的模块路径。更多说明请参阅[依赖管理](/doc/modules/managing-dependencies)。

接下来，将添加一些处理映射（map）的简单代码。

## 添加非泛型函数 {#non_generic_functions}

这一步将添加两个函数，分别将映射中的所有值相加并返回总和。

这里需要声明两个函数，是因为要处理两种不同类型的映射：一种存储 `int64` 值，另一种存储 `float64` 值。

#### 编写代码 {#write-the-code}

1. 使用文本编辑器，在 generics 目录中创建 main.go 文件。你将在这个文件中编写 Go 代码。
2. 将以下包声明粘贴到 main.go 文件顶部。

    ```
    package main
    ```

    独立运行的程序使用 `main` 包，而库则使用其他包名。

3. 在包声明下方粘贴以下两个函数声明。

    ```
    // SumInts adds together the values of m.
    func SumInts(m map[string]int64) int64 {
    	var s int64
    	for _, v := range m {
    		s += v
    	}
    	return s
    }

    // SumFloats adds together the values of m.
    func SumFloats(m map[string]float64) float64 {
    	var s float64
    	for _, v := range m {
    		s += v
    	}
    	return s
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明两个函数，分别将映射中的值相加并返回总和。
        *   `SumFloats` 接收键为 `string`、值为 `float64` 的映射。
        *   `SumInts` 接收键为 `string`、值为 `int64` 的映射。

4. 在 main.go 顶部的包声明下方，粘贴以下 `main` 函数，初始化两个映射，并将它们作为实参传给上一步声明的函数。

    ```
    func main() {
    	// Initialize a map for the integer values
    	ints := map[string]int64{
    		"first":  34,
    		"second": 12,
    	}

    	// Initialize a map for the float values
    	floats := map[string]float64{
    		"first":  35.98,
    		"second": 26.99,
    	}

    	fmt.Printf("Non-Generic Sums: %v and %v\n",
    		SumInts(ints),
    		SumFloats(floats))
    }
    ```

    在这段代码中，你完成了以下操作：

    *   初始化值类型分别为 `float64` 和 `int64` 的两个映射，每个映射各包含两个条目。
    *   调用前面声明的两个函数，分别计算各个映射中值的总和。
    *   输出结果。

5. 在 main.go 顶部，紧接包声明的位置，导入刚才所写代码需要的包。

    代码的开头应如下所示：

    ```
    package main

    import "fmt"
    ```

6. 保存 main.go。

#### 运行代码 {#run-the-code}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Non-Generic Sums: 46 and 62.97
```

借助泛型，这里只需要编写一个函数。接下来，你将添加一个泛型函数，同时处理值为整数或浮点数的映射。

## 添加能处理多种类型的泛型函数 {#add_generic_function}

本节将添加一个泛型函数，让它既能接收值为整数的映射，也能接收值为浮点数的映射，用一个函数实现刚才两个函数的功能。

为了支持这两种类型的值，函数需要一种方式来声明它支持哪些类型。另一方面，调用方代码也需要能够指定传入的是整数映射还是浮点数映射。

为此，你将编写一个除了普通函数形参之外，还声明了_类型形参_的函数。类型形参使函数成为泛型函数，能够处理不同类型的实参。调用时，需要提供_类型实参_和普通函数实参。

每个类型形参都有一个_类型约束_，可以把它理解为描述类型形参的“元类型”。类型约束规定了调用方可以为相应类型形参提供哪些类型实参。

虽然类型形参的约束通常表示一组类型，但在具体调用的编译过程中，类型形参对应的是单个类型，即调用方提供的类型实参。如果该类型不满足类型形参的约束，代码就无法编译。

请注意，类型约束允许的类型必须都支持泛型代码对该类型形参执行的操作。例如，若某个类型形参的约束包含数值类型，而函数却试图对它执行 `string` 操作（如索引），代码就无法编译。

接下来编写的代码将使用一个允许整数或浮点数类型的约束。

#### 编写代码 {#write-the-code-1}

1. 在前面添加的两个函数下方，粘贴以下泛型函数。

    ```
    // SumIntsOrFloats sums the values of map m. It supports both int64 and float64
    // as types for map values.
    func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
        var s V
        for _, v := range m {
            s += v
        }
        return s
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明 `SumIntsOrFloats` 函数，在方括号内声明两个类型形参 `K` 和 `V`，并声明一个使用这些类型形参的普通形参 `m`，其类型为 `map[K]V`。函数返回 `V` 类型的值。
    *   为类型形参 `K` 指定 `comparable` 类型约束。`comparable` 是 Go 的预声明约束，适用于这里这样的场景，允许值可以作为 `==` 和 `!=` 比较运算符操作数的类型。Go 要求映射的键可比较，因此必须将 `K` 约束为 `comparable`，才能用它作为映射的键类型。这也保证了调用方为映射键提供的类型是合法的。
    *   为类型形参 `V` 指定由 `int64` 和 `float64` 两种类型构成的联合约束。`|` 表示这两种类型的并集，也就是说，这个约束允许其中任意一种类型。调用方使用这两种类型中的任意一种作为类型实参，都能通过编译。
    *   将形参 `m` 的类型指定为 `map[K]V`，其中 `K` 和 `V` 就是前面声明的类型形参。因为 `K` 可比较，所以 `map[K]V` 是合法的映射类型。如果没有将 `K` 约束为可比较类型，编译器就会拒绝 `map[K]V`。

2. 在 main.go 的 `main` 函数中，已有代码下方粘贴以下代码。

    ```
    fmt.Printf("Generic Sums: %v and %v\n",
    	SumIntsOrFloats[string, int64](ints),
    	SumIntsOrFloats[string, float64](floats))
    ```

    在这段代码中，你完成了以下操作：

    *   分别传入创建好的两个映射，调用刚才声明的泛型函数。
    *   指定类型实参，即方括号中的类型名称，明确被调用函数中的类型形参应替换为哪些类型。

        下一节将介绍，在许多情况下可以省略调用中的类型实参，因为 Go 能根据代码推断它们。
    *   输出函数返回的总和。

#### 运行代码 {#run-the-code-1}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
```

为了执行代码，编译器会针对每次调用，使用该调用指定的具体类型替换类型形参。

调用刚才编写的泛型函数时，你指定了类型实参，告诉编译器应该用哪些类型替换函数的类型形参。下一节将介绍，很多情况下可以省略这些类型实参，让编译器进行推断。

## 调用泛型函数时省略类型实参 {#remove_type_arguments}

本节将添加一版略作修改的泛型函数调用，用一个小改动简化调用方代码：移除本例中不必显式指定的类型实参。

如果 Go 编译器能够推断出要使用的类型，就可以省略调用中的类型实参。编译器会根据普通函数实参的类型推断类型实参。

注意，并非所有情况都能省略。例如，调用一个没有普通函数参数的泛型函数时，就需要在调用中提供类型实参。

#### 编写代码 {#write-the-code-2}

*   在 main.go 的 `main` 函数中，已有代码下方粘贴以下代码。

    ```
    fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
    	SumIntsOrFloats(ints),
    	SumIntsOrFloats(floats))
    ```

    在这段代码中，你完成了以下操作：

    *   调用泛型函数，并省略类型实参。

#### 运行代码 {#run-the-code-2}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
Generic Sums, type parameters inferred: 46 and 62.97
```

接下来，你将把整数与浮点数的联合约束提取为一个可复用的类型约束，进一步简化函数，并使其他代码也能使用它。

## 声明类型约束 {#declare_type_constraint}

在最后这一节中，你将把前面定义的约束移到一个独立的接口中，以便在多个位置复用。以这种方式声明约束有助于简化代码，尤其是在约束比较复杂时。

_类型约束_通过接口声明，允许满足该接口的类型。例如，如果声明一个包含三个方法的约束接口，并将其用于泛型函数的类型形参，那么调用该函数时提供的类型实参就必须具有这三个方法。

正如本节即将展示的，约束接口也可以包含具体的类型。

#### 编写代码 {#write-the-code-3}

1. 在 import 语句之后、`main` 函数之前，粘贴以下代码，声明类型约束。

    ```
    type Number interface {
        int64 | float64
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明 `Number` 接口类型，将其用作类型约束。
    *   在接口中声明 `int64` 与 `float64` 的联合约束。

        实际上，你将联合约束从函数声明移到了一个新的类型约束中。以后需要将类型形参限制为 `int64` 或 `float64` 时，就可以直接使用 `Number` 约束，而不必重复写出 `int64 | float64`。

2. 在已有函数下方粘贴以下泛型函数 `SumNumbers`。

    ```
    // SumNumbers sums the values of map m. It supports both integers
    // and floats as map values.
    func SumNumbers[K comparable, V Number](m map[K]V) V {
        var s V
        for _, v := range m {
            s += v
        }
        return s
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明一个逻辑与前面的泛型函数相同的函数，但使用新的接口类型作为类型约束，替代直接写出的联合约束。与之前一样，函数的参数和返回值都使用类型形参。

3. 在 main.go 的 `main` 函数中，已有代码下方粘贴以下代码。

    ```
    fmt.Printf("Generic Sums with Constraint: %v and %v\n",
    	SumNumbers(ints),
    	SumNumbers(floats))
    ```

    在这段代码中，你完成了以下操作：

    *   分别传入两个映射调用 `SumNumbers`，并输出各个映射中值的总和。

        与上一节相同，调用泛型函数时省略类型实参，即方括号中的类型名称。Go 编译器可以根据其他实参推断类型实参。

#### 运行代码 {#run-the-code-3}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Non-Generic Sums: 46 and 62.97
Generic Sums: 46 and 62.97
Generic Sums, type parameters inferred: 46 and 62.97
Generic Sums with Constraint: 46 and 62.97
```

## 总结 {#conclusion}

做得不错！你已经初步了解了 Go 泛型。

建议继续学习以下内容：

*   [Go 语言之旅](/tour/)循序渐进地介绍 Go 基础知识。
*   [高效 Go 编程](/doc/effective_go)和[如何编写 Go 代码](/doc/code)介绍了实用的 Go 最佳实践。

## 完整代码 {#completed_code}

你可以在 [Go 在线运行环境](/play/p/apNmfVwogK0)中运行这个程序，只需点击 **Run（运行）**按钮。

```
package main

import "fmt"

type Number interface {
	int64 | float64
}

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first": 34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first": 35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))

	fmt.Printf("Generic Sums, type parameters inferred: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}

// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

// SumIntsOrFloats sums the values of map m. It supports both floats and integers
// as map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

