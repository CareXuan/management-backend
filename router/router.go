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
			user.POST("/delete", controller.DeleteUser)
			user.GET("/permission", controller.GetUserPermission)
		}

		// 权限
		auth := v1.Group("auth")
		{
			auth.POST("/login", controller.Login)
			auth.GET("/permission", controller.AllPermission).Use(middleware.AuthCheck())
			auth.GET("/permission/info", controller.PermissionInfo).Use(middleware.AuthCheck())
			auth.POST("/permission/add", controller.AddPermission).Use(middleware.AuthCheck())
			auth.POST("/permission/delete", controller.RemovePermission).Use(middleware.AuthCheck())
			auth.GET("/roles", controller.AllRoles).Use(middleware.AuthCheck())
			auth.GET("/roles/info", controller.GetRoleInfo).Use(middleware.AuthCheck())
			auth.POST("/roles/add", controller.AddRole).Use(middleware.AuthCheck())
			auth.POST("/roles/delete", controller.DeleteRole).Use(middleware.AuthCheck())
			auth.GET("/organization", controller.AllOrganization).Use(middleware.AuthCheck())
			auth.GET("/organization/info", controller.GetOrganizationInfo).Use(middleware.AuthCheck())
			auth.POST("/organization/add", controller.AddOrganization).Use(middleware.AuthCheck())
			auth.POST("/organization/delete", controller.DeleteOrganization).Use(middleware.AuthCheck())
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
			commonCtr.GET("/wechat_check", controller.WechatCheck)
			commonCtr.Any("/wechat/callback", controller.WechatHandler)
		}

		// 设备
		device := v1.Group("device").Use(middleware.AuthCheck())
		{
			device.GET("/list", controller.GetDeviceList)
			device.GET("/info", controller.GetOneDeviceInfo)
			device.GET("/common_data", controller.GetOneDeviceCommonData)
			device.GET("/service_data", controller.GetOneDeviceServiceData)
			device.GET("/new_service_data", controller.GetOneDeviceNewServiceData)
			device.GET("/location_history", controller.GetDeviceLocationHistory)
			device.GET("/all_location", controller.GetAllDeviceLocation)
			device.GET("/statistic", controller.GetDeviceStatistic)
			device.GET("/get_all_warning", controller.GetAllWarning)
			device.GET("/get_single_warning", controller.GetSingleWarning)
			device.GET("/special_info/log", controller.GetSpecialInfoLog)
			device.POST("/add", controller.AddOneDevice)
			device.POST("/special_info", controller.UpdateSpecialInfo)
			device.POST("/special_info/read", controller.ReadSpecialInfo)
			device.POST("/report", controller.DeviceReport)
		}
	}
}
