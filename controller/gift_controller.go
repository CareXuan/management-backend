package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"prize-draw/model"
	"prize-draw/service"
	"strconv"
)

// List 礼物列表
func List(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.List(c, name, pageInt, pageSizeInt)
}

// Info 礼物详情
func Info(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GiftInfo(c, idInt)
}

// GroupList 礼物组列表
func GroupList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GroupList(c, name, pageInt, pageSizeInt)
}

// GroupInfo 礼物组详情
func GroupInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GroupInfo(c, idInt)
}

// Add 添加礼物
func Add(c *gin.Context) {
	var addReq model.GiftAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.Add(c, addReq)
}

// Delete 删除礼物(软)
func Delete(c *gin.Context) {
	var deleteReq model.GiftDeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.Delete(c, deleteReq)
}

// ChangeStatus 修改可用状态
func ChangeStatus(c *gin.Context) {
	var changeReq model.GiftChangeStatusReq
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeReq(c, changeReq)
}

// GroupAdd 添加礼物组
func GroupAdd(c *gin.Context) {
	var addReq model.GiftGroupAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupAdd(c, addReq)
}

// GroupDelete 删除礼物组(软)
func GroupDelete(c *gin.Context) {
	var deleteReq model.GiftGroupDelete
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupDelete(c, deleteReq)
}

// GroupChangeStatus 修改礼物组状态
func GroupChangeStatus(c *gin.Context) {
	var addReq model.GiftGroupAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupAdd(c, addReq)
}

/*=====================================app=====================================*/
