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
		user := v1.Group("user").Use(middleware.AuthCheck())
		{
			user.GET("/info", controller.GetUserInfo)
			user.GET("/list", controller.GetUserList)
			user.POST("/add", controller.AddUser)
			user.GET("/permission", controller.GetUserPermission)
		}

		// 权限
		auth := v1.Group("auth")
		{
			auth.POST("/login", controller.Login)
			authCheck := auth.Group("/").Use(middleware.AuthCheck())
			{
				authCheck.GET("/permission", controller.AllPermission)
				authCheck.POST("/permission/add", controller.AddPermission)
				authCheck.POST("/permission/delete", controller.RemovePermission)
				authCheck.GET("/roles", controller.AllRoles)
				authCheck.GET("/roles/info", controller.GetRoleInfo)
				authCheck.POST("/roles/add", controller.AddRole)
				authCheck.POST("/roles/delete", controller.DeleteRole)
			}
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
			commonCtr.GET("/wechat_check", controller.WechatCheck)
		}
	}
}
