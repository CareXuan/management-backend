package controller

import (
	"github.com/gin-gonic/gin"
	"management-backend/service"
	"strconv"
)

func GetUserInfo(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.GetUserInfoSer(c, userIdInt)
}

func GetUserPermission(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.GetUserPermissionSer(c, userIdInt)
}
