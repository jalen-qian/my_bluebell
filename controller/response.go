package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
由于每个请求的响应，都会返回固定的结构，导致代码冗余比较严重，这里对响应参数做一个封装
*/

// 定义响应结构体
type ResponseData struct {
	Code  ResCode     `json:"code"`  //响应状态码
	BCode int         `json:"bcode"` //同一状态下，用于区分不同的业务逻辑（一般不用）
	Msg   interface{} `json:"msg"`   //消息提示
	Data  interface{} `json:"data"`  //数据
}

// 返回错误信息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code:  code,
		BCode: 100,
		Msg:   code.Msg(),
		Data:  nil,
	})
}

// 返回错误信息，同时指定自己的提示信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code:  code,
		BCode: 100,
		Msg:   msg,
		Data:  nil,
	})
}

// 返回成功信息，这里使用可变参数来实现bcode的默认值
func ResponseSuccess(c *gin.Context, data interface{}, bcode ...int) {
	b := 100
	if len(bcode) > 0 {
		b = bcode[0]
	}
	c.JSON(http.StatusOK, &ResponseData{
		Code:  CodeSuccess,
		BCode: b,
		Msg:   CodeSuccess.Msg(),
		Data:  data,
	})
}

func ResponseSuccessWithMsg(c *gin.Context, data interface{}, msg interface{}, bcode ...int) {
	b := 100
	if len(bcode) > 0 {
		b = bcode[0]
	}
	c.JSON(http.StatusOK, &ResponseData{
		Code:  CodeSuccess,
		BCode: b,
		Msg:   msg,
		Data:  data,
	})
}
