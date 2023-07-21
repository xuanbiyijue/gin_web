package logic

import (
	"ginQAP/dao/mysql"
	"ginQAP/models"
	"ginQAP/pkg/jwt"
	"go.uber.org/zap"
)

func Login(p *models.UserLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	err = mysql.Login(user)
	if err != nil {
		zap.L().Error("MySQL query error", zap.Error(err))
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		zap.L().Error("JWT failed to generate the token. Procedure", zap.Error(err))
		return nil, err
	}
	user.Token = token
	return
}
