/*
User数据库模型
使用validator库做参数校验
 */

package models


// User 用户信息
type User struct {
	UserID      int64      `db:"user_id"`
	Username    string     `db:"username"`
	Password    string     `db:"password"`
	Email       string     `db:"email"`
	Gender      int        `db:"gender"`
	Token       string
}

// UserLogin 用户登陆参数
type UserLogin struct {
	Username    string     `json:"username" binding:"required"`
	Password    string     `json:"password" binding:"required"`
}

// UserSignup 用户注册参数
type UserSignup struct {
	Username    string     `json:"username" binding:"required"`
	Password    string     `json:"password" binding:"required"`
	RePassword  string     `json:"re_password" binding:"required,eqfield=Password"`
	Email       string     `json:"email" binding:"required"`
	Gender      int        `json:"gender"`
}