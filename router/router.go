package router

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/controller"
	"management-backend/middleware"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) { common.ResOk(c, "Hello World!", nil) })

		// 用户相关
		user := v1.Group("user")
		{
			user.GET("/", controller.GetUserInfo)
			user.GET("/permission", controller.GetUserPermission)
		}

		// 用户相关
		auth := v1.Group("auth")
		{
			auth.POST("/login", controller.Login)
		}
	}
}
