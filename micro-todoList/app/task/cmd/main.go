package main

import (
	"context"
	"fmt"

	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	"github.com/CocaineCong/micro-todoList/app/task/repository/db/dao"
	"github.com/CocaineCong/micro-todoList/app/task/repository/mq"
	"github.com/CocaineCong/micro-todoList/app/task/script"
	"github.com/CocaineCong/micro-todoList/app/task/service"
	"github.com/CocaineCong/micro-todoList/config"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
)

func main() {
	// 配置初始化
	config.Init()
	// Dao 层初始化
	dao.InitDB()
	// Rabbit MQ 初始化
	mq.InitRabbitMQ()
	// 日志初始化
	log.InitLog()

	// 启动一些脚本, 从 MQ 到 MySQL
	loadingScript()

	// etcd注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcTaskService"), // 微服务名字
		micro.Address(config.TaskServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterTaskServiceHandler(microService.Server(), service.GetTaskSrv())
	// 启动微服务
	_ = microService.Run()
}

// loadingScript 开一个线程去执行MQ到MySQL的任务
func loadingScript() {
	ctx := context.Background()
	go script.TaskCreateSync(ctx)
}
