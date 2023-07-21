package controller

import (
	"errors"
	"fmt"
	"ginQAP/dao/mysql"
	"ginQAP/logic"
	"ginQAP/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


// HandleLogin 登录
func HandleLogin(c *gin.Context)  {
	// 暂存参数
	param := new(models.UserLogin)
	// 获取请求参数并校验
	err := c.ShouldBindJSON(param)
	if err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 请求参数有误，返回响应
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 登录逻辑处理，登陆成功获得token
	user, err := logic.Login(param)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", param.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	fmt.Println(user.Token)
	// 登录成功，返回用户信息
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // id值大于1<<53-1  int64类型的最大值是1<<63-1
		"username":  user.Username,
		"email":     user.Email,
		"gender":    user.Gender,
		"token":     user.Token,
	})
}