package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"github.com/CocaineCong/micro-todoList/app/gateway/rpc"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	"github.com/CocaineCong/micro-todoList/pkg/ctl"
)

// ListTaskHandler 获得任务清单
func ListTaskHandler(ctx *gin.Context) {
	// 请求参数绑定
	var taskReq pb.TaskRequest
	if err := ctx.Bind(&taskReq); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数失败"))
		return
	}

	// 获得用户信息
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	// 请求参数绑定用户ID
	taskReq.Uid = uint64(user.Id)
	// rpc 通信，调用服务端的函数
	taskResp, err := rpc.TaskList(ctx, &taskReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "taskResp RPC 调用失败"))
		return
	}
	// 返回响应
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, taskResp))
}

// CreateTaskHandler 创建任务
func CreateTaskHandler(ctx *gin.Context) {
	// 请求参数绑定
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数失败"))
		return
	}
	// 获得用户信息
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	// 请求参数绑定用户ID
	req.Uid = uint64(user.Id)
	// rpc 通信
	taskRes, err := rpc.TaskCreate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskList RPC 调度失败"))
		return
	}
	// 返回响应
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, taskRes))
}

// GetTaskHandler 获得单个任务详情
func GetTaskHandler(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数失败"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.Id = cast.ToUint64(ctx.Param("id"))
	req.Uid = uint64(user.Id)
	taskRes, err := rpc.TaskGet(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskList RPC 调度失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, taskRes))
}

// UpdateTaskHandler 修改任务
func UpdateTaskHandler(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数失败"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.Id = cast.ToUint64(ctx.Param("id"))
	req.Uid = uint64(user.Id)
	taskRes, err := rpc.TaskUpdate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskUpdate RPC 调度失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, taskRes))
}

// DeleteTaskHandler 删除任务
func DeleteTaskHandler(ctx *gin.Context) {
	var req pb.TaskRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.RespError(ctx, err, "绑定参数失败"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.Id = cast.ToUint64(ctx.Param("id"))
	req.Uid = uint64(user.Id)
	taskRes, err := rpc.TaskDelete(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ctl.RespError(ctx, err, "TaskDelete RPC 调度失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, taskRes))
}
