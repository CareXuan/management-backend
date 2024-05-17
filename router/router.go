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
			auth.GET("/permission/detail", controller.GetPermissionDetail)
			auth.POST("/permission/add", controller.AddPermission)
			auth.POST("/permission/delete", controller.RemovePermission)
			auth.GET("/roles", controller.AllRoles)
			auth.GET("/roles/info", controller.GetRoleInfo)
			auth.POST("/roles/add", controller.AddRole)
			auth.POST("/roles/delete", controller.DeleteRole)
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
