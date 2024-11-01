package router

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/controller"
	"my-gpt-server/middleware"
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
			auth.GET("/permission/detail", controller.PermissionDetail).Use(middleware.AuthCheck())
			auth.POST("/permission/add", controller.AddPermission).Use(middleware.AuthCheck())
			auth.POST("/permission/delete", controller.RemovePermission).Use(middleware.AuthCheck())
			auth.GET("/roles", controller.AllRoles).Use(middleware.AuthCheck())
			auth.GET("/roles/info", controller.GetRoleInfo).Use(middleware.AuthCheck())
			auth.POST("/roles/add", controller.AddRole).Use(middleware.AuthCheck())
			auth.POST("/roles/delete", controller.DeleteRole).Use(middleware.AuthCheck())
		}

		// 会员
		member := v1.Group("member").Use(middleware.AuthCheck())
		{
			member.GET("/list", controller.GetMemberList)
			member.GET("/detail", controller.GetMemberDetail)
			member.GET("/recharge/detail", controller.GetMemberRechargeDetail)
			member.POST("/add", controller.AddMember)
			member.POST("/recharge", controller.MemberRecharge)
		}

		// 设备
		device := v1.Group("device").Use(middleware.AuthCheck())
		{
			device.GET("/list", controller.GetDeviceList)
			device.GET("/info", controller.GetDeviceInfo)
			device.GET("/package/list", controller.GetPackageList)
			device.GET("/package/info", controller.GetPackageInfo)
			device.POST("/add", controller.AddDevice)
			device.POST("/status", controller.ChangeStatus)
			device.POST("/package/add", controller.AddPackage)
			device.POST("/package/status", controller.PackageChangeStatus)
		}

		// 预约
		appointment := v1.Group("appointment").Use(middleware.AuthCheck())
		{
			appointment.GET("/list", controller.GetAppointmentList)
			appointment.GET("/detail", controller.GetAppointmentDetail)
			appointment.POST("/add", controller.AddAppointment)
			appointment.POST("/verify", controller.VerifyAppointment)
		}

		// gpt
		gpt := v1.Group("gpt").Use(middleware.AuthCheck())
		{
			gpt.GET("/one", controller.GetOneAnswer)
			gpt.GET("/list", controller.QuestionList)
			gpt.GET("/detail", controller.QuestionDetail)
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
		}
	}
}
