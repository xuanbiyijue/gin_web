package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/micro-todoList/app/gateway/http"
	"github.com/CocaineCong/micro-todoList/app/gateway/middleware"
)

// NewRouter 基于 gin 注册 HTTP 路由
func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	// （全局）跨域中间件
	ginRouter.Use(middleware.Cors())
	// cookie-session 中间件，将在 context 中存放session，这两句有点费解
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	// RESTful API
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务，注册和登录
		v1.POST("/user/register", http.UserRegisterHandler)
		v1.POST("/user/login", http.UserLoginHandler)

		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("tasks", http.ListTaskHandler)
			authed.POST("task", http.CreateTaskHandler)
			authed.GET("task/:id", http.GetTaskHandler)       // task_id
			authed.PUT("task/:id", http.UpdateTaskHandler)    // task_id
			authed.DELETE("task/:id", http.DeleteTaskHandler) // task_id
		}
	}
	return ginRouter
}
