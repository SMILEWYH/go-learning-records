package main

import "github.com/gin-gonic/gin"

// Gin 框架是一个基于 Go语言的轻量级 Web 框架, 它基于 Go 语言的 net/http 包进行开发, 提供了高效、快速、易用的 Web 服务开发体验

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Index(c *gin.Context) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: "2026-11-01",
	})
}

func main() {
	// 1. 初始化
	gin.SetMode("release") // 去除调试日志
	r := gin.Default()
	// 2. 挂载路由
	r.GET("/index", Index)
	// 3. 绑定端口,  运行
	r.Run(":8080")
}
