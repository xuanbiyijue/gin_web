package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/CocaineCong/micro-todoList/app/user/repository/db/model"
)

type UserDao struct {
	*gorm.DB
}

// NewUserDao 创建一个包含一个上下文为给定 context 数据库实例的 UserDao，用于操作用户表
func NewUserDao(ctx context.Context) *UserDao {
	if ctx == nil {
		ctx = context.Background()
	}
	return &UserDao{NewDBClient(ctx)}
}

// FindUserByUserName 通过用户名查询用户
func (dao *UserDao) FindUserByUserName(userName string) (r *model.User, err error) {
	err = dao.Model(&model.User{}).
		Where("user_name = ?", userName).Find(&r).Error
	if err != nil {
		return
	}

	return
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(in *model.User) (err error) {
	return dao.Model(&model.User{}).Create(&in).Error
}
