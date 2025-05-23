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
	status := c.Query("status")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	statusInt, _ := strconv.Atoi(status)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.TaskList(c, name, statusInt, pageInt, pageSizeInt)
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

// TaskDelete 任务删除
func TaskDelete(c *gin.Context) {
	var deleteReq model.TaskDeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.TaskDelete(c, deleteReq)
}

// TaskChangeStatus 任务修改状态
func TaskChangeStatus(c *gin.Context) {
	var changeReq model.TaskChangeStatusReq
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		log.Fatal(err)
		return
	}
	service.TaskChangeStatus(c, changeReq)
}

// TaskCheckList 任务执行记录列表
func TaskCheckList(c *gin.Context) {
	taskId := c.Query("task_id")
	status := c.Query("status")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	taskIdInt, _ := strconv.Atoi(taskId)
	statusInt, _ := strconv.Atoi(status)
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.TaskCheckList(c, taskIdInt, statusInt, pageInt, pageSizeInt)
}

// TaskCheckInfo 任务执行记录详情
func TaskCheckInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.TaskCheckInfo(c, idInt)
}

// TaskCheck 任务审核
func TaskCheck(c *gin.Context) {
	var checkReq model.TaskCheckReq
	if err := c.ShouldBindJSON(&checkReq); err != nil {
		log.Fatal(err)
		return
	}
	service.TaskCheck(c, checkReq)
}

/*=====================================app=====================================*/

// TaskDo 任务提交
func TaskDo(c *gin.Context) {
	var doReq model.TaskDoReq
	if err := c.ShouldBindJSON(&doReq); err != nil {
		log.Fatal(err)
		return
	}
	service.TaskDo(c, doReq)
}

func AppTaskList(c *gin.Context) {
	name := c.Query("name")
	service.AppTaskListSer(c, name)
}
