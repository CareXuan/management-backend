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

// GroupAdd 添加礼物组
func GroupAdd(c *gin.Context) {
	var addReq model.GiftGroupAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupAdd(c, addReq)
}

/*=====================================app=====================================*/
