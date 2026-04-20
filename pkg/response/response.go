package response

import "github.com/gin-gonic/gin"

// APIResponse 定义统一返回格式。
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Fail 返回失败响应。
func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, APIResponse{
		Code: 1,
		Msg:  msg,
	})
}
