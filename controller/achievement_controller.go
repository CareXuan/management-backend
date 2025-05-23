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
	isFinish := c.Query("is_finish")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	isFinishInt, _ := strconv.Atoi(isFinish)
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AchievementList(c, name, isFinishInt, pageInt, pageSizeInt)
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

// AchievementFinish 成就手动完成
func AchievementFinish(c *gin.Context) {
	var finishReq model.AchievementFinishReq
	if err := c.ShouldBindJSON(&finishReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AchievementFinish(c, finishReq)
}

// AchievementDelete 成就删除
func AchievementDelete(c *gin.Context) {
	var deleteReq model.AchievementDeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AchievementDelete(c, deleteReq)
}

/*=====================================app=====================================*/

func AppAchievementList(c *gin.Context) {
	status := c.Query("status")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	statusInt, _ := strconv.Atoi(status)
	service.AppAchievementList(c, pageInt, pageSizeInt, statusInt)
}

func AppReceiveAchievement(c *gin.Context) {
	var receiveReq model.AppAchievementReceiveReq
	if err := c.ShouldBindJSON(&receiveReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AppAchievementReceive(c, receiveReq)
}
