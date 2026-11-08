package main

// go get -u gorm.io/gorm
// go get gorm.io/driver/mysql

// orm 框架的最大优点: 结构体可以映射到表, 表也可以映射到结构体

// 数据库的表中的每一行数据可以映射到一个结构体实例上
// 如果想要整个表的数据, 那么就是使用结构体切片来承接即可

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:b1Ol49pbgjEPLvQOO-gvgDwG2t7X7DAb@tcp(127.0.0.1:3306)/vvyh_sql"), &gorm.Config{})
	if err != nil {
		// 这个 err 比上一章的 sql.Open() 返回的 err 更加强大
		// 除了可以判断 dsn 错误, 还可以判断数据库账号密码是否正确等等
		// 因为它里面默认执行了 ping 操作
		fmt.Println(err)
		return
	}

	fmt.Println(db)
}
