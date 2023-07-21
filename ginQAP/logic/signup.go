package logic

import (
	"ginQAP/dao/mysql"
	"ginQAP/models"
	"ginQAP/pkg/snowflake"
	"go.uber.org/zap"
)

func SignUp(p *models.UserSignup) (err error) {
	// 判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		zap.L().Error("Username has existed", zap.Error(err))
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	//fmt.Println(userID)
	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
		Email:    p.Email,
		Gender:   p.Gender,
	}
	// 保存进数据库
	return mysql.InsertUser(user)
}
