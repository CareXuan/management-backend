package router

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/controller"
)

func InitRouter(r *gin.Engine) {
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) { common.ResOk(c, "Hello World!", conf.Mysql) })
		user := v1.Group("user")
		{
			user.GET("/", controller.GetUserInfo)
		}
	}
}
