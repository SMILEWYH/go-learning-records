//go:build ignore

package main

// 练习 1:编写分页查询接口
// 要求:GET /products?page=1&size=20, 两个参数缺省为 1 和 20;page>=1, size 为 1..100
// 成功返回 application/json 和 {"page":1,"size":20};非法输入返回 400
// 常规:?page=2&size=10 -> 200;无查询参数 -> 默认值
// 边界:page=abc、page=0、size=101 -> 400;POST 同路径 -> 405
// 不用任何真实产品数据库;用 httptest.NewRecorder/NewRequest 检查 handler

// 练习 2:实现局部修改昵称的接口
// 要求:PATCH /users/{id}, id 为正整数;正文 {"name":"小周"};成功返回 id/name, 状态 200
// 校验 application/json;限制 1 KiB;拒绝空名称、未知字段和第二个 JSON 值
// 边界:id=abc、id=0、空名称 -> 400;错误媒体类型 -> 415;正文超限 -> 413
// 示例不保存状态;用 NewRequest+Do 发请求, 响应 Body 必须关闭

// 练习 3:让客户端识别服务端失败
// 建议:func fetchMessage(client *http.Client, url string) (string, error)
// 要求:只接受 200;其他状态返回包含状态码的 error;任何成功拿到的响应都关闭 Body
// 常规:测试服务返回 200 + "ready" -> ready,nil;返回 404/500 -> error
// 边界:使用已取消的 context 创建请求并通过 Do 发送, errors.Is(err, context.Canceled) 为 true
// 客户端设置 Timeout;不要用外部网站或真实休眠很久的接口制造失败

func main() {
}
