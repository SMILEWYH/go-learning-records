package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	// 加载 templates 文件夹下全部的 html 文件, 然后在 r.HTML() 中通过第二个参数进行指定消费
	r.LoadHTMLGlob("templates/*")
	// 加载 templates 文件夹下对应名称的 html 文件
	//r.LoadHTMLFiles("templates/index.html")

	r.GET("", func(c *gin.Context) {
		// 第三个参数为传递到 html 文件中使用, 使用场景为: 轻量版 SSR 服务端渲染
		c.HTML(200, "index.html", map[string]any{
			"frameworkName": "gin",
		})
	})

	r.Run(":8080")
}
