package controller

import (
	"github.com/gin-gonic/gin"
	"management-backend/service"
	"strconv"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("page_size")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	num := c.Query("num")
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.ListSer(c, pageInt, pageSizeInt, num, userIdInt)
}

func Tree(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.TreeSer(c, userIdInt)
}

func Info(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.AmmeterInfoSer(c, ammeterIdInt)
}

func Statistics(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	statisticsType := c.Query("statistics_type")
	statisticsTypeInt, _ := strconv.Atoi(statisticsType)
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	service.AmmeterStatisticsSer(c, statisticsTypeInt, ammeterIdInt, startAt, endAt)
}

func Warning(c *gin.Context) {
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("page_size")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.WarningListSer(c, pageInt, pageSizeInt, ammeterIdInt)
}

func Config(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.ConfigInfoSer(c, ammeterIdInt)
}
