package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	// 第一个参数是路径参数, 第二个才是实际路径
	// 也就是说前端请求的路径是 http://localhost:8080/static/abc.txt, 而实际的文件路径是 static/abc.txt
	r.Static("static", "static")
	// 第一个参数是路径参数, 第二个才是实际单个文件路径
	// 也就是说前端请求的路径是 http://localhost:8080/abcd, 而实际的文件路径是 static/abc.txt
	r.StaticFile("abcd", "static/abc.txt")

	// 注意: 静态文件的路径, 不能再被路由使用了

	r.Run(":8080")
}
