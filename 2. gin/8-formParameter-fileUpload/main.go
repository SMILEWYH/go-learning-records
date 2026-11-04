package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 表单参数, 请求体的数据结构为 form-data、 x-www-form-urlencoded
	r.POST("user", func(c *gin.Context) {
		name := c.PostForm("name")
		age, ok := c.GetPostForm("age")
		if !ok {
			// 请求没有携带 age 参数
		}
		fmt.Println(name, age)
	})

	// 单个文件上传
	r.POST("file-single", func(c *gin.Context) {
		fileHeader, err := c.FormFile("vvyhs headerAvatar")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(fileHeader.Filename) // 文件名
		fmt.Println(fileHeader.Size)     // 文件大小，单位是字节

		file, _ := fileHeader.Open()
		byteData, _ := io.ReadAll(file)

		err = os.WriteFile("xxx.jpg", byteData, 0666)
		fmt.Println(err)

		// 或者直接使用 SaveUploadedFile 方法消费 fileHeader
		err = c.SaveUploadedFile(fileHeader, "uploads/xxx/yyy/"+fileHeader.Filename)
		fmt.Println(err)
	})

	// 多文件上传
	r.POST("file-many", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, headers := range form.File {
			for _, header := range headers {
				c.SaveUploadedFile(header, "uploads/"+header.Filename)
			}
		}
	})

	r.Run(":8080")
}
