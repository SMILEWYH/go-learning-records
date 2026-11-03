package main

import (
	"database/sql"
	"log"

	// go get -u github.com/go-sql-driver/mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 使用账号连接 mysql 实例中的指定名称的数据库
	// 第一个参数: 数据库类型, 第二个参数: 数据源名称(dsn)
	db, err := sql.Open("mysql", "root:b1Ol49pbgjEPLvQOO-gvgDwG2t7X7DAb@tcp(127.0.0.1:3306)/vvyh_sql")
	if err != nil {
		log.Fatal(err) // dsn 格式错误, 这个 err 无法判断账号密码是否正确
	}
	//所以需要 ping 一下
	err = db.Ping()
	if err != nil {
		log.Fatal(err) // 这个 err 可以判断账号密码是否正确
	}
	defer db.Close()

	// 执行sql: 除查询外, 使用Exec函数
	_, err2 := db.Exec("CREATE TABLE users1(id INT NOT NULL , name VARCHAR(20), PRIMARY KEY(ID));")
	if err2 != nil {
		log.Fatal(err2)
	}

	// 查询sql: 使用Query函数
	res, err3 := db.Query("select * from users")
	if err3 != nil {
		log.Fatal(err3)
	}
	for res.Next() { // 默认是从第一行开始, 每次调用Next()都会往下移动一行, 直到没有数据为止
		var id int
		var name string

		err = res.Scan(&id, &name) // 这里scan的字段数量要对上表数据的列数, 否则会报错
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := res.Err(); err != nil { // 检查 res.Next() 迭代过程中是否有错误发生
		log.Fatal(err)
	}

	// 也可以使用QueryRow函数, 查一行
	var id int
	var name string

	err = db.QueryRow("select id, name from users").Scan(&id, &name) // 这里scan的字段数量要对上表数据的列数, 否则会报错
	if err != nil {
		log.Fatal(err)
	}
}
