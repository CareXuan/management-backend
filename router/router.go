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
			commonCtr.GET("/wechat_check", controller.WechatCheck)
		}

		// 设备
		device := v1.Group("device").Use(middleware.AuthCheck())
		{
			device.GET("/list", controller.GetDeviceList)
			device.GET("/info", controller.GetOneDeviceInfo)
			device.GET("/common_data", controller.GetOneDeviceCommonData)
			device.GET("/service_data", controller.GetOneDeviceServiceData)
			device.GET("/location_history", controller.GetDeviceLocationHistory)
			device.GET("/all_location", controller.GetAllDeviceLocation)
			device.GET("/statistic", controller.GetDeviceStatistic)
			device.POST("/add", controller.AddOneDevice)
		}
	}
}
