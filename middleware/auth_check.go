package middleware

import "C"
import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 中获取 token
		token := c.GetHeader("Token")
		if token == "" {
			common.ResAuthorization(c, "请先登录")
			c.Abort()
			return
		}
		var user model.User
		_, err := conf.Mysql.Where("token=?", token).Get(&user)
		if err != nil {
			common.ResError(c, "获取用户信息失败")
			c.Abort()
			return
		}
		if user.Id == 0 {
			common.ResAuthorization(c, "非法请求")
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
	}
}
