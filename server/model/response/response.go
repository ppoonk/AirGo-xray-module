package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERROR   = 1
	SUCCESS = 0
)

// 序列化器
type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`

	//Error  string      `json:"error"`
}

func Result(code int64, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}
func OK(message string, data interface{}, c *gin.Context) {
	Result(SUCCESS, message, data, c)
}

func Fail(message string, data interface{}, c *gin.Context) {
	Result(ERROR, message, data, c)
}
