package rpc

import (
	"go-micro.dev/v4"

	"github.com/CocaineCong/micro-todoList/app/gateway/wrappers"
	"github.com/CocaineCong/micro-todoList/idl/pb"
)

var (
	UserService pb.UserService
	TaskService pb.TaskService
)

// InitRPC 初始化 RPC 配置，创建了 2 个 service
func InitRPC() {
	// 创建一个用于连接用户微服务(user server)的 client 服务
	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		micro.WrapClient(wrappers.NewUserWrapper), // 熔断功能
	)
	// 用户服务调用实例，这个实例是由client和server一起组成的service
	userService := pb.NewUserService("rpcUserService", userMicroService.Client())
	UserService = userService

	// 创建一个用于连接任务微服务(task server)的 client 服务
	taskMicroService := micro.NewService(
		micro.Name("taskService.client"),
		micro.WrapClient(wrappers.NewTaskWrapper), // 熔断功能
	)
	// 任务服务调用实例，这个实例是由client和server一起组成的service
	taskService := pb.NewTaskService("rpcTaskService", taskMicroService.Client())
	TaskService = taskService
}
