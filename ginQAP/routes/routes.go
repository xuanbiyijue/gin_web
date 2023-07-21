package routes

import (
	"ginQAP/controller"
	"ginQAP/logger"
	"ginQAP/middlewares"
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
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.Cors())
	// 注册路由组
	api := r.Group("/api")
	// api下的所有路由在这里写
	{
		// 路由-登录
		api.POST("/login", controller.HandleLogin)
		// 路由-注册
		api.POST("/signup", controller.HandleSignUp)
		// 应用JWT中间件
		api.Use(middlewares.JWTAuthMiddleware())
		{
			// 获得所有帖子
			api.GET("/", controller.HandleGetIssues)
			// 获得某个帖子
			api.GET("/:issueID", controller.HandleGetIssueByID)
			// 发布帖子
			api.POST("/post_issue", controller.HandlePostIssue)
			// 发表答案
			api.POST("/:issueID/post_answer", controller.HandlePostAnswer)
			api.GET("/test", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"data": "seccess",
				})
			})
		}
	}
	return r
}
