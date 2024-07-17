package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-gpt-server/model"
	"my-gpt-server/service"
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
	var user model.AddUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Fatal(err)
		return
	}
	service.AddUserSer(c, user.Id, user.Name, user.Password, user.Phone, user.RoleId)
}

func DeleteUser(c *gin.Context) {
	var user model.AddUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteUserSer(c, user.Id)
}
