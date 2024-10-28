package middleware

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model/rbac"
	"strings"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 中获取 token
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			common.ResAuthorized(c, "请先登录")
			c.Abort()
			return
		}
		part := strings.Split(authorization, " ")
		if len(part) < 2 {
			common.ResAuthorized(c, "非法请求")
			c.Abort()
			return
		}
		var user rbac.User
		_, err := conf.Mysql.Where("token=?", part[1]).Get(&user)
		if err != nil {
			common.ResError(c, "获取用户信息失败")
			c.Abort()
			return
		}
		if user.Id == 0 {
			common.ResAuthorized(c, "非法请求")
			c.Abort()
			return
		}
		// 继续处理请求
		c.Next()
	}
}
