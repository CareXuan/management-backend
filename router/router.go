package router

import (
	"github.com/gin-gonic/gin"
	"prize-draw/common"
	"prize-draw/controller"
	"prize-draw/middleware"
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

		// 礼物
		gift := v1.Group("gift")
		{
			gift.GET("/list", controller.List)
			gift.GET("/info", controller.Info)
			gift.GET("/remain", controller.Remain)
			gift.GET("/raffle/config", controller.RaffleConfig)
			gift.GET("/point/left", controller.PointLeft)
			gift.GET("/group/list", controller.GroupList)
			gift.GET("/group/info", controller.GroupInfo)
			gift.GET("/album/list", controller.AlbumList)
			gift.GET("/album/info", controller.AlbumInfo)
			gift.GET("/album/gift", controller.AlbumGift)
			gift.POST("/add", controller.Add)
			gift.POST("/add_point", controller.AddPoint)
			gift.POST("/reset_point", controller.ResetPoint)
			gift.POST("/raffle/config/set", controller.RaffleConfigSet)
			gift.POST("/delete", controller.Delete)
			gift.POST("/change_status", controller.ChangeStatus)
			gift.POST("/group/add", controller.GroupAdd)
			gift.POST("/group/delete", controller.GroupDelete)
			gift.POST("/group/change_status", controller.GroupChangeStatus)
			gift.POST("/album/add", controller.AlbumAdd)
			gift.POST("/album/delete", controller.AlbumDelete)
		}

		// 任务
		task := v1.Group("task")
		{
			task.GET("/list", controller.TaskList)
			task.GET("/check/list", controller.TaskCheckList)
			task.GET("/info", controller.TaskInfo)
			task.GET("/check/info", controller.TaskCheckInfo)
			task.POST("/add", controller.TaskAdd)
			task.POST("/delete", controller.TaskDelete)
			task.POST("/change_status", controller.TaskChangeStatus)
			task.POST("/check", controller.TaskCheck)
		}

		// 成就
		achievement := v1.Group("achievement")
		{
			achievement.GET("/list", controller.AchievementList)
			achievement.GET("/info", controller.AchievementInfo)
			achievement.POST("/add", controller.AchievementAdd)
			achievement.POST("/finish", controller.AchievementFinish)
			achievement.POST("/delete", controller.AchievementDelete)
		}

		// 小程序端
		app := v1.Group("app")
		{
			task := app.Group("task")
			{
				task.GET("/info")
				task.POST("/do", controller.TaskDo)
			}
			raffle := app.Group("raffle")
			{
				raffle.POST("/one", controller.RaffleOne)
			}
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
			commonCtr.GET("/wechat_check", controller.WechatCheck)
		}
	}
}
