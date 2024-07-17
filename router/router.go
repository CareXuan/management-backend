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
		user := v1.Group("user")
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
			auth.GET("/permission", controller.AllPermission)
			auth.POST("/permission/add", controller.AddPermission)
			auth.POST("/permission/delete", controller.RemovePermission)
			auth.GET("/roles", controller.AllRoles)
			auth.GET("/roles/info", controller.GetRoleInfo)
			auth.POST("/roles/add", controller.AddRole)
			auth.POST("/roles/delete", controller.DeleteRole)
		}

		// 会员
		member := v1.Group("member")
		{
			member.GET("/list", controller.GetMemberList)
			member.GET("/detail", controller.GetMemberDetail)
			member.POST("/add", controller.AddMember)
			member.POST("/recharge", controller.MemberRecharge)
		}

		// 设备
		device := v1.Group("device")
		{
			device.GET("/list", controller.GetDeviceList)
			device.GET("/info", controller.GetDeviceInfo)
			device.POST("/info", controller.AddDevice)
		}

		// 预约
		appointment := v1.Group("appointment")
		{
			appointment.GET("/list", controller.GetAppointmentList)
			appointment.GET("/detail", controller.GetAppointmentDetail)
		}

		// gpt
		gpt := v1.Group("gpt")
		{
			gpt.GET("/one", controller.GetOneAnswer)
			gpt.GET("/list", controller.QuestionList)
			gpt.GET("/detail", controller.QuestionDetail)
		}
	}
}
