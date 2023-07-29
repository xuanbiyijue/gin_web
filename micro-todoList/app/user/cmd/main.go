package main

import (
	"fmt"

	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	"github.com/CocaineCong/micro-todoList/app/user/repository/db/dao"
	"github.com/CocaineCong/micro-todoList/app/user/service"
	"github.com/CocaineCong/micro-todoList/config"
	"github.com/CocaineCong/micro-todoList/idl/pb"
)

func main() {
	// 初始化配置
	config.Init()
	// Dao 层初始化
	dao.InitDB()
	// etcd 注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcUserService"),             // 微服务名字
		micro.Address(config.UserServiceAddress), // 微服务地址
		micro.Registry(etcdReg),                  // etcd注册件
	)
	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册，这里的参数是需要一个 Server，另一个是实现了 Server 接口的结构体
	_ = pb.RegisterUserServiceHandler(microService.Server(), service.GetUserSrv())
	// 启动微服务
	_ = microService.Run()
}
