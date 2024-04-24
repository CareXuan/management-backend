package router

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
)

func InitRouter(r *gin.Engine) {
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) { common.ResOk(c, "Hello World!", conf.Mysql) })
	}
}
