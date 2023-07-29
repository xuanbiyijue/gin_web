package service

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"

	"github.com/CocaineCong/micro-todoList/app/user/repository/db/dao"
	"github.com/CocaineCong/micro-todoList/app/user/repository/db/model"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/e"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

// UserSrv 实现用户微服务的 Server 接口
type UserSrv struct{}

// GetUserSrv 只一次地创建一个 UserSrv 并返回其指针
func GetUserSrv() *UserSrv {
	// 保证只执行一次
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// UserLogin 用户登录
func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	// 获得用户信息
	user, err := dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		resp.Code = e.ERROR
		return
	}
	// 校验密码
	if !user.CheckPassword(req.Password) {
		resp.Code = e.InvalidParams
		return
	}
	// 构造响应，填入用户的详细信息
	resp.UserDetail = BuildUser(user)
	return
}

// UserRegister 用户注册
func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	// 检查两次密码输入是否正确
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		return
	}
	resp.Code = e.SUCCESS
	// 检查用户是否存在
	_, err = dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.ERROR
			return
		}
	}
	// 创建用户模型
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.ERROR
		return
	}
	// 插入到数据库
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.ERROR
		return
	}
	// 构造响应
	resp.UserDetail = BuildUser(user)
	return
}

// BuildUser 构造proto文件定义的用户模型
func BuildUser(item *model.User) *pb.UserModel {
	userModel := pb.UserModel{
		Id:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return &userModel
}
