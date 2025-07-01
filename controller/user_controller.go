package controller

import (
	"env-backend/model"
	"env-backend/service"
	"github.com/gin-gonic/gin"
	"log"
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

func GetUserList(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	userName := c.Query("user_name")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetUserListSrv(c, userName, pageInt, pageSizeInt)
}

func AddUser(c *gin.Context) {
	var req model.AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.AddUserSer(c, req)
}

func DeleteUser(c *gin.Context) {
	var req model.DeleteUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteUserSer(c, req)
}
