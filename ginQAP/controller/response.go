/*
封装要返回的数据，其中包括响应码、响应码提示信息、数据
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": 10000, // 程序中的错误码
	"msg": xx,     // 提示信息
	"data": {},    // 数据
}
*/

type ResponseData struct {
	Code ResponseCode     `json:"code"`
	Msg  interface{}      `json:"msg"`            // 对code的解释信息
	Data interface{}      `json:"data,omitempty"` // 要返回的数据 omitempty意思是没有数据时可以不要
}

func ResponseError(c *gin.Context, code ResponseCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResponseCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
