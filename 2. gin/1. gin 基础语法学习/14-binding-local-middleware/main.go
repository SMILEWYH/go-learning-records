package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 局部中间件

func Home(c *gin.Context) {
	fmt.Println("Home")
	c.String(200, "Home")
}
func M1(c *gin.Context) {
	fmt.Println("M1 请求部分")
	c.Next()

	// c.Abort() // 只输出 "M1 请求部分" 和 "M1 响应部分", 如果只想要前者, 不想要后者则补一个 return 即可
	// 因为中断流程后不会再走 Hone 函数, 导致接口没有返回任何数据, 所以我们可以在这里补一个 c.String(500, "fail")
	// c.String(500, "fail")
	// return // 这样就不会输出 "M1 响应部分", 注意如果 M1 前还有其他中间件, 那么其他中间件的 c.Next() 后的代码还是会正常运行

	fmt.Println("M1 响应部分")
}
func M2(c *gin.Context) {
	fmt.Println("M2 请求部分")
	c.Next()
	fmt.Println("M2 响应部分")
}

// 还可以在中间件中使用 c.Abort() 进行拦截，使用后接口流程会从当前这个中间件原路返回

func main() {
	r := gin.Default()
	r.GET("", M1, M2, Home) // 洋葱模型
	r.Run(":8080")

	// 输出顺序:
	// 	M1 请求部分
	// M2 请求部分
	// Home
	// M2 响应部分
	// M1 响应部分
}
