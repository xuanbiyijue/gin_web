package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/micro-todoList/app/gateway/rpc"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/ctl"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
	"github.com/CocaineCong/micro-todoList/pkg/utils"
	"github.com/CocaineCong/micro-todoList/types"
)

// UserRegisterHandler 用户注册
func UserRegisterHandler(ctx *gin.Context) {
	// 请求参数绑定
	var req pb.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserRegister Bind 绑定参数失败"))
		return
	}
	// rpc 通信
	userResp, err := rpc.UserRegister(ctx, &req)
	if err != nil {
		log.LogrusObj.Errorf("UserRegister:%v", err)
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserRegister RPC 调用失败"))
		return
	}
	// 返回响应
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, userResp))
}

// UserLoginHandler 用户登录
func UserLoginHandler(ctx *gin.Context) {
	// 请求参数绑定
	var req pb.UserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "UserLogin Bind 绑定参数失败"))
		return
	}
	// rpc 通信
	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "UserLogin RPC 调用失败"))
		return
	}
	// 生成 token，token 里包含用户ID
	token, err := utils.GenerateToken(uint(userResp.UserDetail.Id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "GenerateToken 失败"))
		return
	}
	// 返回响应
	res := &types.TokenData{
		User:  userResp,
		Token: token,
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, res))
}
