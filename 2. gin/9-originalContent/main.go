package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// 这一章主要说关于请求的请求体的获取、常见的三种数据类型

func main() {
	r := gin.Default()

	r.POST("user", func(c *gin.Context) {
		byteData, _ := io.ReadAll(c.Request.Body)
		fmt.Println(string(byteData))
		// 读了之后, body就没了, 阅后即焚

		// 但我们可以把读取请求体得到的字节流又转换回去, 再赋值即可
		c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

		// 通过 c.Request.Header 来获取这次请求的请求体的数据类型、分割符
		fmt.Println(c.Request.Header)
		// Content-Type:[multipart/form-data; boundary=--------------------------052455317193517003536866]

		// c.Request.Body 请求体的不同数据类型:
		// 1. form-data: 通常用于携带 file 的表单数据
		// ----------------------------853882779395683818968400
		// Content-Disposition: form-data; name="age"

		// 1234
		// ----------------------------853882779395683818968400--

		// 2. x-www-form-urlencoded: 通常用于不携带 file 的表单数据
		// name=%E6%9E%AB%E6%9E%AB&age=1234

		// 3. json:
		// 		{
		//     "name": "枫枫",
		//     "age": 23
		// }

	})

	r.Run(":8080")
}
