package router

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/controller"
	siemens2 "switchboard-backend/controller/siemens"
	"switchboard-backend/middleware"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) { common.ResOk(c, "Hello World!", nil) })

		device := v1.Group("device")
		{
			device.POST("/tttt", controller.Tttt)
		}

		siemens := v1.Group("siemens")
		{
			s7 := siemens.Group("s7")
			{
				s7.GET("/list", siemens2.List)
				s7.GET("/info", siemens2.Info)
				s7.POST("/add", siemens2.Add)
				s7.POST("/add_data", siemens2.AddData)
			}
		}

		// 用户相关
		user := v1.Group("user")
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
			auth.GET("/permission", controller.AllPermission)
			auth.GET("/permission/info", controller.PermissionInfo)
			auth.POST("/permission/add", controller.AddPermission)
			auth.POST("/permission/delete", controller.RemovePermission)
			auth.GET("/roles", controller.AllRoles)
			auth.GET("/roles/info", controller.GetRoleInfo)
			auth.POST("/roles/add", controller.AddRole)
			auth.POST("/roles/delete", controller.DeleteRole)
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
		}
	}
}
