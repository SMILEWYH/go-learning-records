package res

import "github.com/gin-gonic/gin"

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type Code int

const (
	RoleErrCode    Code = 1001
	NetworkErrCode Code = 1002
)

var codeMap = map[Code]string{
	RoleErrCode:    "权限错误",
	NetworkErrCode: "网络错误",
}

func init() {
	// 可能是一个
}

func response(c *gin.Context, r Response) {
	c.JSON(200, r)
}

func Ok(c *gin.Context, data any, msg string) {
	response(c, Response{
		Code: 0,
		Data: data,
		Msg:  msg,
	})
}

func OkWithData(c *gin.Context, data any) {
	Ok(c, data, "成功")
}

func OkWithMsg(c *gin.Context, msg string) {
	Ok(c, map[string]any{}, msg)
}

func Fail(c *gin.Context, code int, data any, msg string) {
	response(c, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func FailWithMsg(c *gin.Context, msg string) {
	response(c, Response{
		Code: 7,
		Data: nil,
		Msg:  msg,
	})
}

func FailWithCode(c *gin.Context, code Code) {
	// 去找code对应的msg
	msg, ok := codeMap[code]
	if !ok {
		msg = "未知错误"
	}
	response(c, Response{
		Code: int(code),
		Data: nil,
		Msg:  msg,
	})
}
