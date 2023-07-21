package controller

import (
	"errors"
	"ginQAP/dao/mysql"
	"ginQAP/logic"
	"ginQAP/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandleSignUp 注册
func HandleSignUp(c *gin.Context) {
	// 获取参数并校验
	var param models.UserSignup
	err := c.ShouldBindJSON(&param)
	if err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 业务处理
	err = logic.SignUp(&param)
	if err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}