package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"prize-draw/model"
	"prize-draw/service"
	"strconv"
)

// TaskList 任务列表
func TaskList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.TaskList(c, name, pageInt, pageSizeInt)
}

// TaskInfo 任务详情
func TaskInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.TaskInfo(c, idInt)
}

// TaskAdd 任务添加
func TaskAdd(c *gin.Context) {
	var addReq model.TaskAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.TaskAdd(c, addReq)
}
