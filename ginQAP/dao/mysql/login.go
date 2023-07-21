package mysql

import (
	"database/sql"
	"ginQAP/models"
	"ginQAP/utils"
	"go.uber.org/zap"
)


// Login 登录
func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录的密码
	sqlStr := `select user_id, username, password, email, gender from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		zap.L().Error("MySQL query failed", zap.Error(err))
		return err
	}
	// 判断密码是否正确
	password := utils.EncryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

