package routes

import (
	"ginQA/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)


func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// gin设置成发布模式
		gin.SetMode(gin.ReleaseMode)
	}
	// 创建一个自定义的引擎
	r := gin.New()
	// 注册全局中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册路由组
	api := r.Group("/api")
	// api下的所有路由在这里写
	{
		api.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"data": "seccess",
			})
		})
	}
	return r
}
