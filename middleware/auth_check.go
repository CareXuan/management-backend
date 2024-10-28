package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 中获取 token
		authHeader := c.GetHeader("Token")
		if authHeader == "" {
			common.ResAuthorization(c, "请先登录")
			return
		}

		fmt.Println(authHeader)
		// 继续处理请求
		c.Next()
	}
}
