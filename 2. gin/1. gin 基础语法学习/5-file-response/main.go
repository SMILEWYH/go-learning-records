package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// 第一种下载文件方式(不常用):
	// 前端请求这个get接口, 然后后端通过这个接口传递二进制的文件内容给前端, 然后前端收到的是 blob 类型的数据
	// 下一步在前端代码中就可以对blob类型的文件内容进行转换和消费, 但是要想保存到用户电脑里, 还需要额外的一步
	// 就是为在内存中生成的blob类型的文件内容创建一个临时下载地址, 然后通过前端 a 标签触发下载到电脑
	r.GET("", func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")             // 表示响应内容是通用二进制数据
		c.Header("Content-Disposition", "attachment; filename=文件下载.txt") // 告诉浏览器按附件下载，并指定下载文件名
		c.Header("fileName", "xxx.png")                                  // 自定义响应头参数-1
		c.Header("msg", "文件下载成功")                                        // 自定义响应头参数-2
		c.File("文件下载.txt")
	})

	// 	async downloadFile(row) {
	//    this.$http({
	//       method: 'post',
	//       url: 'file/upload',
	//       data:postData,
	//       responseType: "blob"
	//    }).then(res => {
	//       const _res = res.data
	//       let blob = new Blob([_res], {
	//             type: 'application/png'
	//           });
	//       let downloadElement = document.createElement("a");
	//       let href = window.URL.createObjectURL(blob); //创建下载的链接
	//       downloadElement.href = href;
	//       downloadElement.download = res.headers["fileName"]; //下载后文件名
	//       document.body.appendChild(downloadElement);
	//       downloadElement.click(); //点击下载
	//       document.body.removeChild(downloadElement); //下载完成移除元素
	//       window.URL.revokeObjectURL(href); //释放掉blob对象
	//     })}

	// 所以我们通常使用第二种下载文件方式:
	// 请求后端下载接口的请求, 后端不返回实际文件内容, 而是生成一个临时下载地址, 然后返回到前端
	// 前端构造a标签, herf指向后端返回的临时下载地址和 dawnload 配置文件名, 然后通过a标签触发下载到电脑

	// 所以第一种下载方式更加适合下载文件的目的是为了在前端js中消费, 而第二种下载方式更加适合下载文件的目的是为了保存到用户电脑里
	// 再补充一点如果下载需要认证信息的话, 那么下载链接需要带上 cookie 这样的认证信息

	r.Run(":8080")
}
