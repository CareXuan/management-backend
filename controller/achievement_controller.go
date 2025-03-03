package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"prize-draw/model"
	"prize-draw/service"
	"strconv"
)

// AchievementList 成就列表
func AchievementList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AchievementList(c, name, pageInt, pageSizeInt)
}

// AchievementInfo 成就详情
func AchievementInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.AchievementInfo(c, idInt)
}

// AchievementAdd 成就添加
func AchievementAdd(c *gin.Context) {
	var addReq model.AchievementAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AchievementAdd(c, addReq)
}
