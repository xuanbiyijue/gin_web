package main

import (
	"context"
	"fmt"
	"ginQAP/config"
	"ginQAP/dao/mysql"
	"ginQAP/dao/redis"
	"ginQAP/logger"
	"ginQAP/pkg/snowflake"
	"ginQAP/routes"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// configFilePath 配置文件的路径
const configFilePath = "./config.yaml"

func main() {
	// 初始化配置
	err := config.Init(configFilePath)
	if err != nil {
		panic(fmt.Errorf("init config failed, err: %v", err))
		return
	}

	// 初始化日志
	err = logger.Init(config.Config.LogConfig, config.Config.Mode)
	if err != nil {
		panic(fmt.Errorf("init logger failed, err: %v", err))
		return
	}
	// 退出前需要调用，用来刷新缓冲
	defer zap.L().Sync()

	// 初始化Mysql
	err = mysql.Init(config.Config.MySQLConfig)
	if err != nil {
		zap.L().Error("init mysql failed!", zap.Error(err))
		return
	}
	defer mysql.Close()

	// 初始化Redis
	err = redis.Init(config.Config.RedisConfig)
	if err != nil {
		zap.L().Error("init redis failed!", zap.Error(err))
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(config.Config.StartTime, config.Config.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err: %v\n", err)
		return
	}

	// 注册路由组
	r := routes.Setup(config.Config.Mode)

	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Config.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
