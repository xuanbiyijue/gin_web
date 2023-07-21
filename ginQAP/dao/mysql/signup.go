package mysql

import (
	"ginQAP/models"
	"ginQAP/utils"
	"go.uber.org/zap"
)

func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}


func InsertUser(user *models.User) (err error) {
	// 先对密码进行加密
	password := utils.EncryptPassword(user.Password)
	// 插入
	sqlStr := `insert into user (user_id, username, password, email, gender) values(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, password, user.Email, user.Gender)
	if err != nil {
		zap.L().Error("Insert user failded", zap.Error(err))
	}
	return
}