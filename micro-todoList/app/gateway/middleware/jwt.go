package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/micro-todoList/pkg/ctl"
	"github.com/CocaineCong/micro-todoList/pkg/utils"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code uint32
		code = 200
		// 拿到 token
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
			c.JSON(500, gin.H{
				"code": code,
				"msg":  "鉴权失败",
			})
		}
		// 解析 token
		claims, err := utils.ParseToken(token)
		if err != nil {
			code = 401
			c.JSON(500, gin.H{
				"code": code,
				"msg":  "鉴权失败",
			})
			c.Abort()
		}
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.Id}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}
