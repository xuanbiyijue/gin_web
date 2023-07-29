package rpc

import (
	"context"
	"errors"

	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/e"
)

// UserLogin 用户登陆，通过调用 service 的 client api
func UserLogin(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	// 调用 service 的 client api 完成通信
	resp, err = UserService.UserLogin(ctx, req)
	if err != nil {
		return
	}
	if resp.Code != e.SUCCESS {
		err = errors.New(e.GetMsg(int(resp.Code)))
		return
	}
	return
}

// UserRegister 用户注册，通过调用 service 的 client api
func UserRegister(ctx context.Context, req *pb.UserRequest) (resp *pb.UserDetailResponse, err error) {
	resp, err = UserService.UserRegister(ctx, req)
	if err != nil {
		return
	}
	return
}
