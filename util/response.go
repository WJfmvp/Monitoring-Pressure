package util

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, ResponseData{
		Code: ErrCodeSuccess,
		Msg:  GetMessage(ErrCodeSuccess),
		Data: data,
	})
}

func ResponseError(c *gin.Context, code int) {
	c.JSON(200, ResponseData{
		Code: code,
		Msg:  GetMessage(code),
	})
}

func ResponseErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(200, ResponseData{
		Code: code,
		Msg:  msg,
	})
}
