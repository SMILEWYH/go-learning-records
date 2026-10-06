<!--{
  "Title": "教程：使用 Go 和 Gin 开发 RESTful API",
  "Breadcrumb": true
}-->

本教程介绍如何使用 Go 和 [Gin Web 框架](https://gin-gonic.com/en/docs/)（简称 Gin）编写 RESTful Web 服务 API。

如果已经初步了解 Go 及其工具，学习本教程会更顺利。如果这是你第一次接触 Go，请先阅读[教程：Go 入门](/doc/tutorial/getting-started)，快速了解基础知识。

Gin 简化了构建 Web 应用程序（包括 Web 服务）时的许多编码任务。本教程将使用 Gin 进行请求路由、读取请求信息，并将响应序列化为 JSON。

在本教程中，你将构建一个具有两个端点路径的 RESTful API 服务器。示例项目用于存储经典爵士唱片的数据。

本教程包含以下内容：

1. 设计 API 端点。
2. 为代码创建文件夹。
3. 创建数据。
4. 编写返回所有条目的处理函数。
5. 编写添加新条目的处理函数。
6. 编写返回指定条目的处理函数。

**注意：**其他教程请参阅[教程列表](/doc/tutorial/index.html)。

如果想在 Google Cloud Shell 中以交互形式完成本教程，请点击下面的按钮。

[![在 Cloud Shell 中打开](https://gstatic.com/cloudssh/images/open-btn.png)](https://ide.cloud.google.com/?cloudshell_workspace=~&walkthrough_tutorial_url=https://raw.githubusercontent.com/golang/tour/master/tutorial/web-service-gin.md)

## 准备工作 {#prerequisites}

*   **Go。**建议使用最新版本的 Go 学习本教程。安装步骤请参阅[安装 Go](/doc/install)。
*   **代码编辑工具。**任何文本编辑器都可以。
*   **命令行终端。**Go 可以在 Linux 和 Mac 的各种终端，以及 Windows 的 PowerShell 或 cmd 中正常使用。
*   **curl 工具。**Linux 和 Mac 通常已安装这个工具。Windows 10 Insider build 17063 及更高版本也包含它；更早的 Windows 版本可能需要自行安装。更多说明请参阅 [Tar 和 Curl 登陆 Windows](https://docs.microsoft.com/en-us/virtualization/community/team-blog/2017/20171219-tar-and-curl-come-to-windows)。

## 设计 API 端点 {#design_endpoints}

你将构建一个 API，让客户端访问一家销售经典黑胶唱片的商店。因此，需要提供端点，让客户端能够为用户获取和添加专辑。

开发 API 时，通常先设计端点。端点越容易理解，API 的使用者就越容易正确使用它。

本教程将创建以下端点：

/albums

*   `GET`——获取全部专辑列表，以 JSON 返回。
*   `POST`——根据请求发送的 JSON 数据添加新专辑。

/albums/:id

*   `GET`——根据 ID 获取一张专辑，以 JSON 返回专辑数据。

接下来，你将为代码创建文件夹。

## 为代码创建文件夹 {#create_folder}

首先，为即将编写的代码创建一个项目。

1. 打开命令行终端，切换到用户主目录。

    在 Linux 或 Mac 上：

    ```
    $ cd
    ```

    在 Windows 上：

    ```
    C:\> cd %HOMEPATH%
    ```

2. 在命令行中创建名为 web-service-gin 的目录，用来存放代码。

    ```
    $ mkdir web-service-gin
    $ cd web-service-gin
    ```

3. 创建模块，用于管理依赖。

    运行 `go mod init` 命令，并传入代码所属模块的路径。

    ```
    $ go mod init example/web-service-gin
    go: creating new go.mod: module example/web-service-gin
    ```

    这个命令会创建 go.mod 文件，用于记录后续添加的依赖。有关使用模块路径命名模块的更多说明，请参阅[依赖管理](/doc/modules/managing-dependencies#naming_module)。

接下来，你将设计用于处理数据的数据结构。

## 创建数据 {#create_data}

为简化教程，这里将数据保存在内存中。实际的 API 通常会与数据库交互。

注意，将数据保存在内存中意味着：每次停止服务器时，当前专辑数据都会丢失；重新启动时，将重新创建初始数据。

#### 编写代码 {#write-the-code}

1. 使用文本编辑器，在 web-service-gin 目录中创建 main.go 文件。你将在这个文件中编写 Go 代码。
2. 将以下包声明粘贴到 main.go 文件顶部。

    ```
    package main
    ```

    独立运行的程序使用 `main` 包，而库则使用其他包名。

3. 在包声明下方粘贴以下 `album` 结构体声明，用于在内存中存储专辑数据。

    ``json:"artist"`` 这样的结构体标签用于指定将结构体序列化为 JSON 时对应的字段名称。如果不加这些标签，JSON 就会使用结构体中以大写字母开头的字段名，而这种命名风格在 JSON 中不太常见。

    ```
    // album represents data about a record album.
    type album struct {
    	ID     string  `json:"id"`
    	Title  string  `json:"title"`
    	Artist string  `json:"artist"`
    	Price  float64 `json:"price"`
    }
    ```

4. 在刚添加的结构体声明下方，粘贴以下元素类型为 `album` 的切片，作为初始数据。

    ```
    // albums slice to seed record album data.
    var albums = []album{
    	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
    }
    ```

接下来，你将编写代码，实现第一个端点。

## 编写返回所有条目的处理函数 {#all_items}

当客户端请求 `GET /albums` 时，需要以 JSON 返回所有专辑。

为此，需要编写以下代码：

*   准备响应的逻辑。
*   将请求路径映射到这段逻辑的代码。

注意，这里的编写顺序与运行时的执行顺序相反：先实现被依赖的代码，再实现依赖它的代码。

#### 编写代码 {#write-the-code-1}

1. 在上一节添加的结构体相关代码下方，粘贴以下代码，用于获取专辑列表。

    `getAlbums` 函数将 `album` 结构体切片转换为 JSON，并写入响应。

    ```
    // getAlbums responds with the list of all albums as JSON.
    func getAlbums(c *gin.Context) {
    	c.IndentedJSON(http.StatusOK, albums)
    }
    ```

    在这段代码中，你完成了以下操作：

    *   编写 `getAlbums` 函数，接收一个指向 [`gin.Context`](https://pkg.go.dev/github.com/gin-gonic/gin#Context) 的参数。这个函数可以使用其他名称；Gin 和 Go 都没有规定它必须采用某种特定命名格式。

        `gin.Context` 是 Gin 的核心组成部分，用于承载请求详情、校验和序列化 JSON 等。（虽然名称相似，但它与 Go 标准库的 [`context`](/pkg/context/) 包不同。）

    *   调用 [`Context.IndentedJSON`](https://pkg.go.dev/github.com/gin-gonic/gin#Context.IndentedJSON)，将结构体数据序列化为 JSON 并写入响应。

        第一个实参是要发送给客户端的 HTTP 状态码。这里传入 `net/http` 包中的 [`StatusOK`](https://pkg.go.dev/net/http#StatusOK) 常量，表示 `200 OK`。

        也可以将 `Context.IndentedJSON` 换成 [`Context.JSON`](https://pkg.go.dev/github.com/gin-gonic/gin#Context.JSON)，发送更紧凑的 JSON。调试时，带缩进的形式更易阅读，而两者的大小差异通常不大。

2. 在 main.go 中，紧接 `albums` 切片声明的位置，粘贴以下代码，将处理函数关联到端点路径。

    这会建立一项路由关联，让 `getAlbums` 处理 `/albums` 端点路径的请求。

    ```
    func main() {
    	router := gin.Default()
    	router.GET("/albums", getAlbums)

    	router.Run("localhost:8080")
    }
    ```

    在这段代码中，你完成了以下操作：

    *   使用 [`Default`](https://pkg.go.dev/github.com/gin-gonic/gin#Default) 初始化 Gin 路由器。
    *   使用 [`GET`](https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.GET) 方法，将 HTTP `GET` 方法与 `/albums` 路径的组合关联到处理函数。

        注意，这里传入的是 `getAlbums` 函数本身，写的是它的_名称_，而不是调用函数后的_结果_。写成 `getAlbums()`（带括号）才是调用函数。

    *   使用 [`Run`](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.Run) 方法，将路由器接入 `http.Server` 并启动服务器。

3. 在 main.go 顶部，紧接包声明的位置，导入刚才所写代码需要的包。

    代码开头应如下所示：

    ```
    package main

    import (
    	"net/http"

    	"github.com/gin-gonic/gin"
    )
    ```

4. 保存 main.go。

#### 运行代码 {#run-the-code}

1. 将 Gin 模块添加为依赖。

    在命令行中使用 [`go get`](/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them)，将 github.com/gin-gonic/gin 模块添加为当前模块的依赖。参数中的点表示“获取当前目录中代码所需的依赖”。

    ```
    $ go get .
    go get: added github.com/gin-gonic/gin v1.7.2
    ```

    Go 会解析并下载这个依赖，以满足上一步添加的 `import` 声明。

2. 在 main.go 所在目录的命令行中运行代码。参数中的点表示“运行当前目录中的代码”。

    ```
    $ go run .
    ```

    代码开始运行后，HTTP 服务器就启动了，可以向它发送请求。

3. 打开一个新的命令行窗口，使用 `curl` 向正在运行的 Web 服务发送请求。

    ```
    $ curl http://localhost:8080/albums
    ```

    这个命令应该显示服务的初始数据。

    ```
    [
            {
                    "id": "1",
                    "title": "Blue Train",
                    "artist": "John Coltrane",
                    "price": 56.99
            },
            {
                    "id": "2",
                    "title": "Jeru",
                    "artist": "Gerry Mulligan",
                    "price": 17.99
            },
            {
                    "id": "3",
                    "title": "Sarah Vaughan and Clifford Brown",
                    "artist": "Sarah Vaughan",
                    "price": 39.99
            }
    ]
    ```

API 已经运行起来了！下一节将添加处理 `POST` 请求的端点，用于新增条目。

## 编写添加新条目的处理函数 {#add_item}

当客户端向 `/albums` 发送 `POST` 请求时，需要将请求体中描述的专辑加入已有专辑数据。

为此，需要编写以下代码：

*   将新专辑添加到现有列表的逻辑。
*   将 `POST` 请求路由到这段逻辑的代码。

#### 编写代码 {#write-the-code-2}

1. 添加将专辑数据加入列表的代码。

    在 `import` 语句之后的合适位置粘贴以下代码。（可以放在文件末尾；Go 不要求按特定顺序声明函数。）

    ```
    // postAlbums adds an album from JSON received in the request body.
    func postAlbums(c *gin.Context) {
    	var newAlbum album

    	// Call BindJSON to bind the received JSON to
    	// newAlbum.
    	if err := c.BindJSON(&newAlbum); err != nil {
    		return
    	}

    	// Add the new album to the slice.
    	albums = append(albums, newAlbum)
    	c.IndentedJSON(http.StatusCreated, newAlbum)
    }
    ```

    在这段代码中，你完成了以下操作：

    *   使用 [`Context.BindJSON`](https://pkg.go.dev/github.com/gin-gonic/gin#Context.BindJSON)，将请求体中的 JSON 绑定到 `newAlbum`。
    *   将根据 JSON 初始化的 `album` 结构体追加到 `albums` 切片。
    *   在响应中返回 `201` 状态码，以及表示新添加专辑的 JSON 数据。

2. 修改 `main` 函数，添加对 `router.POST` 的调用，如下所示。

    ```
    func main() {
    	router := gin.Default()
    	router.GET("/albums", getAlbums)
    	router.POST("/albums", postAlbums)

    	router.Run("localhost:8080")
    }
    ```

    在这段代码中，你完成了以下操作：

    *   将 `/albums` 路径的 `POST` 方法关联到 `postAlbums` 函数。

        Gin 可以将处理函数关联到 HTTP 方法与路径的组合。这样，即使请求发送到同一个路径，也可以根据客户端使用的方法分别路由。

#### 运行代码 {#run-the-code-1}

1. 如果上一节的服务器仍在运行，先将它停止。
2. 在 main.go 所在目录的命令行中运行代码。

    ```
    $ go run .
    ```

3. 在另一个命令行窗口中，使用 `curl` 向正在运行的 Web 服务发送请求。

    ```
    $ curl http://localhost:8080/albums \
        --include \
        --header "Content-Type: application/json" \
        --request "POST" \
        --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
    ```

    这个命令应该显示响应头，以及新添加专辑的 JSON 数据。

    ```
    HTTP/1.1 201 Created
    Content-Type: application/json; charset=utf-8
    Date: Wed, 02 Jun 2021 00:34:12 GMT
    Content-Length: 116

    {
        "id": "4",
        "title": "The Modern Sound of Betty Carter",
        "artist": "Betty Carter",
        "price": 49.99
    }
    ```

4. 和上一节一样，使用 `curl` 获取完整的专辑列表，确认新专辑已经加入。

    ```
    $ curl http://localhost:8080/albums \
        --header "Content-Type: application/json" \
        --request "GET"
    ```

    这个命令应该显示专辑列表。

    ```
    [
            {
                    "id": "1",
                    "title": "Blue Train",
                    "artist": "John Coltrane",
                    "price": 56.99
            },
            {
                    "id": "2",
                    "title": "Jeru",
                    "artist": "Gerry Mulligan",
                    "price": 17.99
            },
            {
                    "id": "3",
                    "title": "Sarah Vaughan and Clifford Brown",
                    "artist": "Sarah Vaughan",
                    "price": 39.99
            },
            {
                    "id": "4",
                    "title": "The Modern Sound of Betty Carter",
                    "artist": "Betty Carter",
                    "price": 49.99
            }
    ]
    ```

下一节将添加代码，处理获取指定条目的 `GET` 请求。

## 编写返回指定条目的处理函数 {#specific_item}

当客户端请求 `GET /albums/[id]` 时，需要返回 ID 与路径参数 `id` 匹配的专辑。

为此，需要完成以下操作：

*   添加获取指定专辑的逻辑。
*   将路径映射到这段逻辑。

#### 编写代码 {#write-the-code-3}

1. 在上一节添加的 `postAlbums` 函数下方，粘贴以下代码，用于获取指定专辑。

    `getAlbumByID` 函数提取请求路径中的 ID，然后查找匹配的专辑。

    ```
    // getAlbumByID locates the album whose ID value matches the id
    // parameter sent by the client, then returns that album as a response.
    func getAlbumByID(c *gin.Context) {
    	id := c.Param("id")

    	// Loop over the list of albums, looking for
    	// an album whose ID value matches the parameter.
    	for _, a := range albums {
    		if a.ID == id {
    			c.IndentedJSON(http.StatusOK, a)
    			return
    		}
    	}
    	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
    }
    ```

    在这段代码中，你完成了以下操作：

    *   使用 [`Context.Param`](https://pkg.go.dev/github.com/gin-gonic/gin#Context.Param)，从 URL 中获取 `id` 路径参数。将这个处理函数映射到路径时，需要在路径中包含相应的参数占位符。
    *   遍历切片中的 `album` 结构体，查找 `ID` 字段与 `id` 参数值匹配的专辑。如果找到，就将这个 `album` 结构体序列化为 JSON，并随 HTTP `200 OK` 状态码一起返回。

        如前所述，实际服务通常会通过数据库查询完成这种查找。

    *   如果没有找到专辑，则使用 [`http.StatusNotFound`](https://pkg.go.dev/net/http#StatusNotFound) 返回 HTTP `404` 错误。

2. 最后修改 `main`，新增一次 `router.GET` 调用，路径使用 `/albums/:id`，如下所示。

    ```
    func main() {
    	router := gin.Default()
    	router.GET("/albums", getAlbums)
    	router.GET("/albums/:id", getAlbumByID)
    	router.POST("/albums", postAlbums)

    	router.Run("localhost:8080")
    }
    ```

    在这段代码中，你完成了以下操作：

    *   将 `/albums/:id` 路径关联到 `getAlbumByID` 函数。在 Gin 中，路径片段前面的冒号表示这个片段是路径参数。

#### 运行代码 {#run-the-code-2}

1. 如果上一节的服务器仍在运行，先将它停止。
2. 在 main.go 所在目录的命令行中运行代码，启动服务器。

    ```
    $ go run .
    ```

3. 在另一个命令行窗口中，使用 `curl` 向正在运行的 Web 服务发送请求。

    ```
    $ curl http://localhost:8080/albums/2
    ```

    这个命令应该显示指定 ID 对应专辑的 JSON 数据。如果未找到专辑，则会返回包含错误消息的 JSON。

    ```
    {
            "id": "2",
            "title": "Jeru",
            "artist": "Gerry Mulligan",
            "price": 17.99
    }
    ```

## 总结 {#conclusion}

恭喜！你已经使用 Go 和 Gin 编写了一个简单的 RESTful Web 服务。

建议继续学习以下内容：

*   如果刚接触 Go，可以阅读[高效 Go 编程](/doc/effective_go)和[如何编写 Go 代码](/doc/code)，了解实用的最佳实践。
*   [Go 语言之旅](/tour/)循序渐进地介绍 Go 基础知识。
*   关于 Gin 的更多信息，请参阅 [Gin Web 框架包文档](https://pkg.go.dev/github.com/gin-gonic/gin)或 [Gin Web 框架使用文档](https://gin-gonic.com/en/docs/)。

## 完整代码 {#completed_code}

下面是本教程所构建应用程序的完整代码。

```
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
```
