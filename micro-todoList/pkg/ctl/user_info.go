package ctl

import (
	"context"
	"errors"
)

type key int

var userKey key

type UserInfo struct {
	Id uint `json:"id"`
}

// GetUserInfo 从 context 里获得用户信息
func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := FromContext(ctx)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}

// NewContext 将用户信息存放 context
func NewContext(ctx context.Context, u *UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext 从 context 里获得用户信息
func FromContext(ctx context.Context) (*UserInfo, bool) {
	u, ok := ctx.Value(userKey).(*UserInfo)
	return u, ok
}

func InitUserInfo(ctx context.Context) {
	// TOOD 放缓存，之后的用户信息，走缓存
}
