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
			auth.GET("/permission", controller.AllPermission).Use(middleware.AuthCheck())
			auth.GET("/permission/detail", controller.GetPermissionDetail).Use(middleware.AuthCheck())
			auth.POST("/permission/add", controller.AddPermission).Use(middleware.AuthCheck())
			auth.POST("/permission/delete", controller.RemovePermission).Use(middleware.AuthCheck())
			auth.GET("/roles", controller.AllRoles).Use(middleware.AuthCheck())
			auth.GET("/roles/info", controller.GetRoleInfo).Use(middleware.AuthCheck())
			auth.POST("/roles/add", controller.AddRole).Use(middleware.AuthCheck())
			auth.POST("/roles/delete", controller.DeleteRole).Use(middleware.AuthCheck())
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
		}

		// 电动车设备
		device := v1.Group("device")
		{
			device.GET("/list", controller.DeviceList)
			device.GET("/signals", controller.SignalDetailList)
			device.POST("/report", controller.DeviceReport)
		}

		// 电表设备
		ammeter := v1.Group("ammeter")
		{
			ammeter.GET("/list", controller.List)
			ammeter.GET("/tree", controller.Tree)
			ammeter.GET("/tree/manager", controller.TreeManager)
			ammeter.POST("/tree/add", controller.AddNode)
			ammeter.POST("/tree/delete", controller.DeleteNode)
			ammeter.GET("/info", controller.Info)
			ammeter.POST("/switch/pwd", controller.SetSwitchPwd)
			ammeter.POST("/switch", controller.ChangeSwitch)
			ammeter.GET("/statistics", controller.Statistics)
			ammeter.GET("/warning", controller.Warning)
			ammeter.POST("/warning/deal", controller.ChangeWarning)
			ammeter.GET("/config", controller.Config)
			ammeter.POST("/config/update", controller.UpdateConfig)
			ammeter.POST("/data/add", controller.AddTestData)
		}
	}
}
