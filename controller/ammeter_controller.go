package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/model"
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
	status := c.Query("status")
	statusInt, _ := strconv.Atoi(status)
	ammeterType := c.Query("type")
	ammeterTypeInt, _ := strconv.Atoi(ammeterType)
	service.ListSer(c, pageInt, pageSizeInt, num, statusInt, ammeterTypeInt, userIdInt)
}

func Tree(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.TreeSer(c, userIdInt)
}

func AddNode(c *gin.Context) {
	var req model.AmmeterNodeAdd
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.AddTreeNodeSer(c, req.NodeId, req.NodeType, req.NodeModel, req.Num, req.Card, req.Location, req.ParentId, req.Managers)
}

func Info(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.AmmeterInfoSer(c, ammeterIdInt)
}

func ChangeSwitch(c *gin.Context) {
	var req model.ChangeSwitchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeAmmeterSwitchSer(c, req.AmmeterId, req.Switch)
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
	warningType := c.Query("type")
	warningTypeInt, _ := strconv.Atoi(warningType)
	status := c.Query("status")
	statusInt, _ := strconv.Atoi(status)
	startDealTime := c.Query("start_deal_time")
	endDealTime := c.Query("end_deal_time")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	dealUser := c.Query("deal_user")
	service.WarningListSer(c, pageInt, pageSizeInt, warningTypeInt, statusInt, startDealTime, endDealTime, startTime, endTime, dealUser, ammeterIdInt)
}

func ChangeWarning(c *gin.Context) {
	var req model.AmmeterWarningUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeWarningStatusSer(c, req.WarningId, req.Status, req.UserId)
}

func Config(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.ConfigInfoSer(c, ammeterIdInt)
}

func UpdateConfig(c *gin.Context) {
	var req model.AmmeterConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.UpdateConfigSer(c, req)
}

func AddTestData(c *gin.Context) {
	var req model.AmmeterDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.AddTestData(c, req.AmmeterId, req.Type, req.Value, req.CreateTime)
}
