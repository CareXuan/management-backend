package router

import (
	"data_verify/common"
	"data_verify/controller"
	"data_verify/middleware"
	"github.com/gin-gonic/gin"
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

		// 数据
		data := v1.Group("data").Use(middleware.AuthCheck())
		{
			data.POST("/sbk/upload", controller.UploadAndReadDBF).Use(middleware.LimitUploadSize(int64(100 << 20)))
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.GET("/")
			//commonCtr.POST("/upload", controller.Upload)
			//commonCtr.GET("/wechat_check", controller.WechatCheck)
		}
	}
}
