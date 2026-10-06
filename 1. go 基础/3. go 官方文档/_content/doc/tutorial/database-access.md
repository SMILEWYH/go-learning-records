<!--{
  "Title": "教程：访问关系型数据库",
  "Breadcrumb": true
}-->

本教程介绍如何使用 Go 及其标准库中的 `database/sql` 包访问关系型数据库。

如果已经初步了解 Go 及其工具，学习本教程会更顺利。如果这是你第一次接触 Go，请先阅读[教程：Go 入门](/doc/tutorial/getting-started)，快速了解基础知识。

接下来使用的 [`database/sql`](https://pkg.go.dev/database/sql) 包提供了连接数据库、执行事务、取消正在进行的操作等功能所需的类型和函数。关于这个包的详细用法，请参阅[访问数据库](/doc/database/index)。

在本教程中，你将创建一个数据库，然后编写访问它的代码。示例项目用于存储经典爵士唱片的数据。

你将依次完成以下内容：

1. 为代码创建文件夹。
2. 设置数据库。
3. 导入数据库驱动。
4. 获取数据库句柄并建立连接。
5. 查询多行数据。
6. 查询单行数据。
7. 添加数据。

**注意：**其他教程请参阅[教程列表](/doc/tutorial/index.html)。

## 准备工作 {#prerequisites}

*   **安装 [MySQL](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/) 关系型数据库管理系统（DBMS）。**
*   **安装 Go。**安装步骤请参阅[安装 Go](/doc/install)。
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

2. 在命令行中创建名为 data-access 的目录，用来存放代码。

    ```
    $ mkdir data-access
    $ cd data-access
    ```

3. 创建一个模块，管理本教程中即将添加的依赖。

    运行 `go mod init` 命令，并传入新代码的模块路径。

    ```
    $ go mod init example/data-access
    go: creating new go.mod: module example/data-access
    ```

    这个命令会创建 go.mod 文件，用于记录后续添加的依赖。更多说明请参阅[依赖管理](/doc/modules/managing-dependencies)。

    **注意：**实际开发时，应根据自己的需求指定更合适的模块路径。更多说明请参阅[依赖管理](/doc/modules/managing-dependencies#naming_module)。

接下来，你将创建数据库。

## 设置数据库 {#set_up_database}

这一步将创建后续要使用的数据库。你将使用数据库管理系统自带的命令行工具，创建数据库和表，并添加数据。

这个数据库将存储经典爵士黑胶唱片的信息。

这里使用的是 [MySQL 命令行工具](https://dev.mysql.com/doc/refman/8.0/en/mysql.html)，大多数数据库管理系统也提供了功能类似的命令行工具。

1. 打开一个新的命令行终端。
2. 在命令行中登录数据库管理系统。以下以 MySQL 为例：

    ```
    $ mysql -u root -p
    Enter password:

    mysql>
    ```

3. 在 `mysql` 命令提示符下创建数据库。

    ```
    mysql> create database recordings;
    ```

4. 切换到刚创建的数据库，以便添加表。

    ```
    mysql> use recordings;
    Database changed
    ```

5. 在文本编辑器中，在 data-access 文件夹下创建 create-tables.sql 文件，用于保存创建表的 SQL 脚本。
6. 将以下 SQL 代码粘贴到文件中，然后保存。

    ```
    DROP TABLE IF EXISTS album;
    CREATE TABLE album (
      id         INT AUTO_INCREMENT NOT NULL,
      title      VARCHAR(128) NOT NULL,
      artist     VARCHAR(255) NOT NULL,
      price      DECIMAL(5,2) NOT NULL,
      PRIMARY KEY (`id`)
    );

    INSERT INTO album
      (title, artist, price)
    VALUES
      ('Blue Train', 'John Coltrane', 56.99),
      ('Giant Steps', 'John Coltrane', 63.99),
      ('Jeru', 'Gerry Mulligan', 17.99),
      ('Sarah Vaughan', 'Sarah Vaughan', 34.98);
    ```

    在这段 SQL 代码中，你完成了以下操作：

    *   删除（drop）名为 `album` 的表。先执行这条命令，便于以后需要重新创建表时再次运行整个脚本。

    *   创建 `album` 表，包含 `title`、`artist`、`price` 和 `id` 四列。每行的 `id` 值由数据库管理系统自动生成。

    *   添加四行数据。

7. 在 `mysql` 命令提示符下运行刚创建的脚本。

    使用以下形式的 `source` 命令：

    ```
    mysql> source /path/to/create-tables.sql
    ```

8. 在数据库管理系统的命令提示符下，使用 `SELECT` 语句确认表及其中的数据已成功创建。

    ```
    mysql> select * from album;
    +----+---------------+----------------+-------+
    | id | title         | artist         | price |
    +----+---------------+----------------+-------+
    |  1 | Blue Train    | John Coltrane  | 56.99 |
    |  2 | Giant Steps   | John Coltrane  | 63.99 |
    |  3 | Jeru          | Gerry Mulligan | 17.99 |
    |  4 | Sarah Vaughan | Sarah Vaughan  | 34.98 |
    +----+---------------+----------------+-------+
    4 rows in set (0.00 sec)
    ```

接下来，将编写 Go 代码连接数据库，以便执行查询。

## 查找并导入数据库驱动 {#import_driver}

现在已经有了包含数据的数据库，可以开始编写 Go 代码了。

查找并导入合适的数据库驱动。它会将通过 `database/sql` 包中的函数发出的请求，转换为数据库能够理解的请求。

1. 在浏览器中访问 [SQLDrivers](/wiki/SQLDrivers) Wiki 页面，查找可用的驱动。

    根据页面上的列表选择驱动。本教程访问 MySQL 时使用 [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql/)。

2. 记下驱动的包路径，这里是 `github.com/go-sql-driver/mysql`。

3. 在文本编辑器中创建用于编写 Go 代码的文件，将它保存为前面创建的 data-access 目录中的 main.go。

4. 将以下代码粘贴到 main.go 中，导入驱动包。

    ```
    package main

    import "github.com/go-sql-driver/mysql"
    ```

    在这段代码中，你完成了以下操作：

    *   将代码放入 `main` 包，使其能够作为独立程序执行。

    *   导入 MySQL 驱动 `github.com/go-sql-driver/mysql`。

导入驱动后，就可以开始编写访问数据库的代码了。

## 获取数据库句柄并建立连接 {#get_handle}

现在编写一些 Go 代码，通过数据库句柄访问数据库。

你将使用指向 `sql.DB` 结构体的指针，它表示对某个数据库的访问入口。

#### 编写代码 {#write-the-code}

1. 在 main.go 中刚添加的 `import` 代码下方，粘贴以下 Go 代码，创建数据库句柄。

    ```
    var db *sql.DB

    func main() {
    	// Capture connection properties.
    	cfg := mysql.NewConfig()
    	cfg.User = os.Getenv("DBUSER")
    	cfg.Passwd = os.Getenv("DBPASS")
    	cfg.Net = "tcp"
    	cfg.Addr = "127.0.0.1:3306"
    	cfg.DBName = "recordings"

    	// Get a database handle.
    	var err error
    	db, err = sql.Open("mysql", cfg.FormatDSN())
    	if err != nil {
    		log.Fatal(err)
    	}

    	pingErr := db.Ping()
    	if pingErr != nil {
    		log.Fatal(pingErr)
    	}
    	fmt.Println("Connected!")
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明类型为 [`*sql.DB`](https://pkg.go.dev/database/sql#DB) 的变量 `db`，作为数据库句柄。

        将 `db` 定义为全局变量是为了简化示例。在生产代码中，应避免使用全局变量，例如将它传给需要它的函数，或封装到结构体中。

    *   使用 MySQL 驱动的 [`Config`](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config) 类型以及该类型的 [`FormatDSN`](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config.FormatDSN) 方法，收集连接属性并格式化为 DSN（数据源名称）连接字符串。

        相比直接拼写连接字符串，使用 `Config` 结构体能让代码更易读。

    *   调用 [`sql.Open`](https://pkg.go.dev/database/sql#Open)，传入 `FormatDSN` 的返回值，初始化 `db` 变量。

    *   检查 `sql.Open` 是否返回错误。例如，数据库连接参数格式不正确时，它可能失败。

        为了简化代码，这里调用 `log.Fatal` 将错误输出到控制台并终止程序。生产代码应采用更妥善的错误处理方式。

    *   调用 [`DB.Ping`](https://pkg.go.dev/database/sql#DB.Ping)，确认能够连接数据库。运行时，`sql.Open` 是否立即建立连接取决于驱动。这里通过 `Ping` 确认 `database/sql` 包在需要时能够连接数据库。

    *   检查 `Ping` 是否返回错误，以判断连接是否失败。

    *   如果 `Ping` 连接成功，输出一条消息。

2. 在 main.go 顶部，紧接包声明的位置，导入刚才所写代码需要的包。

    文件开头现在应如下所示：

    ```
    package main

    import (
    	"database/sql"
    	"fmt"
    	"log"
    	"os"

    	"github.com/go-sql-driver/mysql"
    )
    ```

3. 保存 main.go。

#### 运行代码 {#run-the-code}

1. 将 MySQL 驱动模块添加为依赖。

    使用 [`go get`](/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them)，将 github.com/go-sql-driver/mysql 模块添加为当前模块的依赖。参数中的点表示“获取当前目录中代码所需的依赖”。

    ```
    $ go get .
    go: added filippo.io/edwards25519 v1.1.0
    go: added github.com/go-sql-driver/mysql v1.8.1
    ```

    因为上一步在 `import` 声明中导入了这个包，Go 会下载相应的依赖。有关依赖管理的更多说明，请参阅[添加依赖](/doc/modules/managing-dependencies#adding_dependency)。

2. 在命令行中设置 `DBUSER` 和 `DBPASS` 环境变量，供 Go 程序读取。

    在 Linux 或 Mac 上：

    ```
    $ export DBUSER=username
    $ export DBPASS=password
    ```

    在 Windows 上：

    ```
    C:\Users\you\data-access> set DBUSER=username
    C:\Users\you\data-access> set DBPASS=password
    ```

3. 在 main.go 所在目录的命令行中运行代码。输入 `go run` 并加上点作为参数，表示“运行当前目录中的包”。

    ```
    $ go run .
    Connected!
    ```

连接成功！接下来，你将查询一些数据。

## 查询多行数据 {#multiple_rows}

本节将使用 Go 执行一个返回多行数据的 SQL 查询。

对于可能返回多行的 SQL 语句，使用 `database/sql` 包提供的 `Query` 方法，再遍历返回的各行。（稍后的[查询单行数据](#single_row)一节会介绍如何查询一行。）

#### 编写代码 {#write-the-code-1}

1. 在 main.go 中，紧挨 `func main` 上方的位置，粘贴以下 `Album` 结构体定义，用于存储查询返回的行数据。

    ```
    type Album struct {
    	ID     int64
    	Title  string
    	Artist string
    	Price  float32
    }
    ```

2. 在 `func main` 下方粘贴以下 `albumsByArtist` 函数，用来查询数据库。

    ```
    // albumsByArtist queries for albums that have the specified artist name.
    func albumsByArtist(name string) ([]Album, error) {
    	// An albums slice to hold data from returned rows.
    	var albums []Album

    	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
    	if err != nil {
    		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    	}
    	defer rows.Close()
    	// Loop through rows, using Scan to assign column data to struct fields.
    	for rows.Next() {
    		var alb Album
    		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
    			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    		}
    		albums = append(albums, alb)
    	}
    	if err := rows.Err(); err != nil {
    		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    	}
    	return albums, nil
    }
    ```

    在这段代码中，你完成了以下操作：

    *   声明一个元素类型为前面定义的 `Album` 的切片 `albums`，存储返回的行数据。结构体字段的名称和类型与数据库列对应。

    *   使用 [`DB.Query`](https://pkg.go.dev/database/sql#DB.Query) 执行 `SELECT` 语句，查询指定艺术家的专辑。

        `Query` 的第一个参数是 SQL 语句，后面可以传入零个或多个任意类型的参数，用于提供 SQL 语句中各个参数的值。将 SQL 语句与参数值分开传入，而不是使用 `fmt.Sprintf` 等方式拼接，能够让 `database/sql` 包将值与 SQL 文本分别传递，避免由参数拼接引入的 SQL 注入风险。

    *   使用 defer 延迟关闭 `rows`，使函数退出时释放它持有的资源。

    *   遍历返回的各行，使用 [`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) 将每行的列值赋给 `Album` 结构体字段。

        `Scan` 接收一组指向 Go 值的指针，列值将写入这些指针指向的位置。这里使用 `&` 运算符取得 `alb` 变量各字段的指针，再传给 `Scan`，由它通过指针更新结构体字段。

    *   在循环内，检查将列值读入结构体字段时是否发生错误。

    *   在循环内，将新的 `alb` 追加到 `albums` 切片中。

    *   循环结束后，使用 `rows.Err` 检查整个查询过程中是否发生错误。注意，如果遍历结果时查询失败，必须在这里检查错误，才能发现结果并不完整。

3. 修改 `main` 函数，调用 `albumsByArtist`。

    在 `func main` 的末尾添加以下代码。

    ```
    albums, err := albumsByArtist("John Coltrane")
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("Albums found: %v\n", albums)
    ```

    在新增代码中，你完成了以下操作：

    *   调用刚添加的 `albumsByArtist` 函数，并将返回结果赋给新的 `albums` 变量。

    *   输出结果。

#### 运行代码 {#run-the-code-1}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
```

接下来，你将查询单行数据。

## 查询单行数据 {#single_row}

本节将使用 Go 查询数据库中的单行数据。

对于已知最多返回一行的 SQL 语句，可以使用 `QueryRow`，比使用 `Query` 再遍历结果更简单。

#### 编写代码 {#write-the-code-2}

1. 在 `albumsByArtist` 下方粘贴以下 `albumByID` 函数。

    ```
    // albumByID queries for the album with the specified ID.
    func albumByID(id int64) (Album, error) {
    	// An album to hold data from the returned row.
    	var alb Album

    	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
    	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
    		if err == sql.ErrNoRows {
    			return alb, fmt.Errorf("albumsById %d: no such album", id)
    		}
    		return alb, fmt.Errorf("albumsById %d: %v", id, err)
    	}
    	return alb, nil
    }
    ```

    在这段代码中，你完成了以下操作：

    *   使用 [`DB.QueryRow`](https://pkg.go.dev/database/sql#DB.QueryRow) 执行 `SELECT` 语句，查询指定 ID 的专辑。

        它返回一个 `sql.Row`。为了简化调用方代码，`QueryRow` 不直接返回错误，而是将查询错误（例如 `sql.ErrNoRows`）留到后续调用 `Row.Scan` 时返回。

    *   使用 [`Row.Scan`](https://pkg.go.dev/database/sql#Row.Scan) 将列值复制到结构体字段中。

    *   检查 `Scan` 是否返回错误。

        特殊错误 `sql.ErrNoRows` 表示查询没有返回任何行。通常可以将它替换为更具体的错误说明，例如本例中的“no such album（不存在该专辑）”。

2. 修改 `main`，调用 `albumByID`。

    在 `func main` 的末尾添加以下代码。

    ```
    // Hard-code ID 2 here to test the query.
    alb, err := albumByID(2)
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("Album found: %v\n", alb)
    ```

    在新增代码中，你完成了以下操作：

    *   调用刚添加的 `albumByID` 函数。

    *   输出返回的专辑信息，包括专辑 ID。

#### 运行代码 {#run-the-code-2}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
Album found: {2 Giant Steps John Coltrane 63.99}
```

接下来，你将向数据库添加一张专辑。

## 添加数据 {#add_data}

本节将使用 Go 执行 SQL `INSERT` 语句，在数据库中添加一行数据。

前面已经介绍了如何使用 `Query` 和 `QueryRow` 执行返回数据的 SQL 语句。执行_不返回数据_的 SQL 语句时，应使用 `Exec`。

#### 编写代码 {#write-the-code-3}

1. 在 `albumByID` 下方粘贴以下 `addAlbum` 函数，用于向数据库插入新专辑，然后保存 main.go。

    ```
    // addAlbum adds the specified album to the database,
    // returning the album ID of the new entry
    func addAlbum(alb Album) (int64, error) {
    	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
    	if err != nil {
    		return 0, fmt.Errorf("addAlbum: %v", err)
    	}
    	id, err := result.LastInsertId()
    	if err != nil {
    		return 0, fmt.Errorf("addAlbum: %v", err)
    	}
    	return id, nil
    }
    ```

    在这段代码中，你完成了以下操作：

    *   使用 [`DB.Exec`](https://pkg.go.dev/database/sql#DB.Exec) 执行 `INSERT` 语句。

        和 `Query` 一样，`Exec` 接收 SQL 语句，以及紧随其后的 SQL 参数值。

    *   检查执行 `INSERT` 时是否发生错误。

    *   使用 [`Result.LastInsertId`](https://pkg.go.dev/database/sql#Result.LastInsertId) 获取新插入行的 ID。

    *   检查获取 ID 时是否发生错误。

2. 修改 `main`，调用新的 `addAlbum` 函数。

    在 `func main` 的末尾添加以下代码。

    ```
    albID, err := addAlbum(Album{
    	Title:  "The Modern Sound of Betty Carter",
    	Artist: "Betty Carter",
    	Price:  49.99,
    })
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("ID of added album: %v\n", albID)
    ```

    在新增代码中，你完成了以下操作：

    *   传入一张新专辑调用 `addAlbum`，并将新专辑的 ID 赋给 `albID` 变量。

#### 运行代码 {#run-the-code-3}

在 main.go 所在目录的命令行中运行代码。

```
$ go run .
Connected!
Albums found: [{1 Blue Train John Coltrane 56.99} {2 Giant Steps John Coltrane 63.99}]
Album found: {2 Giant Steps John Coltrane 63.99}
ID of added album: 5
```

## 总结 {#conclusion}

恭喜！你已经使用 Go 完成了对关系型数据库的一些基本操作。

建议继续学习以下内容：

*   阅读数据访问指南，深入了解本教程中简要涉及的各个主题。

*   如果刚接触 Go，可以阅读[高效 Go 编程](/doc/effective_go)和[如何编写 Go 代码](/doc/code)，了解实用的最佳实践。

*   [Go 语言之旅](/tour/)循序渐进地介绍 Go 基础知识。

## 完整代码 {#completed_code}

下面是本教程所构建应用程序的完整代码。

```
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "recordings"

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
```
