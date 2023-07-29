package main

import (
	"fmt"
	"time"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"

	"github.com/CocaineCong/micro-todoList/app/gateway/router"
	"github.com/CocaineCong/micro-todoList/app/gateway/rpc"
	"github.com/CocaineCong/micro-todoList/config"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
)

func main() {
	// 初始化项目配置
	config.Init()
	// 初始化 RPC 配置
	rpc.InitRPC()
	// 初始化日志
	log.InitLog()
	// 注册 etcd（etcd必须先运行起来），这里是一个通用的注册接口，只要给定了主机号和端口号，他会自己寻找服务
	// 这里 registry.NewRegistry 返回一个mDNS的注册表，mDNS可在局域网内通过广播实现资源的查询。也可使用consul
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)

	// 创建web微服务实例，使用gin暴露http接口并注册到etcd
	// web.RegisterTTL(time.Second*30): 如果服务在这个时间内没有重新注册，那么注册中心会认为服务已经下线，从而删除服务的信息。
	// 这个参数可以用来防止服务在异常情况下崩溃而没有及时注销，导致注册中心中存在无效的服务信息。可以不进行设置而使用默认值
	// web.RegisterInterval(time.Second*15): 是设置间隔多久再次注册服务。这个参数可以用来保持服务的在线状态，防止服务被注册中心误删。
	// 一般来说，web.RegisterInterval应该小于web.RegisterTTL，这样才能确保服务在过期之前重新注册。
	// web.Metadata(map[string]string{"protocol": "http"}): 设置服务的元数据，也就是一些额外的信息，比如协议、版本、标签等。
	// 这些元数据可以用来描述服务的特征或者用于服务发现的过滤条件。web.Metadata是一个map[string]string类型，可以自定义键值对。
	server := web.NewService(
		web.Name("httpService"),                             // 服务名称
		web.Address("127.0.0.1:4000"),                       // 服务地址
		web.Handler(router.NewRouter()),                     // HTTP 路由引擎
		web.Registry(etcdReg),                               // 服务发现
		web.RegisterTTL(time.Second*30),                     // 设置注册服务的过期时间
		web.RegisterInterval(time.Second*15),                // 设置间隔多久再次注册服务
		web.Metadata(map[string]string{"protocol": "http"}), // 设置元数据，表明使用http协议
	)
	// 初始化并运行（接收命令行参数）
	_ = server.Init()
	_ = server.Run()
}
