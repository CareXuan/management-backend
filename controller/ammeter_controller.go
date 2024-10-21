package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/model/ammeter"
	"management-backend/service"
	"strconv"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("page_size")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	card := c.Query("card")
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	status := c.Query("status")
	statusInt, _ := strconv.Atoi(status)
	ammeterType := c.Query("type")
	ammeterTypeInt, _ := strconv.Atoi(ammeterType)
	service.ListSer(c, pageInt, pageSizeInt, card, statusInt, ammeterTypeInt, userIdInt)
}

func Tree(c *gin.Context) {
	userId := c.Query("user_id")
	userIdInt, _ := strconv.Atoi(userId)
	service.TreeSer(c, userIdInt)
}

func TreeManager(c *gin.Context) {
	service.TreeManagerSer(c)
}

func AddNode(c *gin.Context) {
	var req ammeter.AmmeterNodeAdd
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.AddTreeNodeSer(c, req.NodeId, req.NodeType, req.NodeModel, req.Num, req.Card, req.Location, req.ParentId, req.Managers)
}

func DeleteNode(c *gin.Context) {
	var req ammeter.AmmeterNodeAdd
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteTreeNodeSer(c, req.NodeId)
}

func Info(c *gin.Context) {
	ammeterId := c.Query("ammeter_id")
	ammeterIdInt, _ := strconv.Atoi(ammeterId)
	service.AmmeterInfoSer(c, ammeterIdInt)
}

func ChangeSwitch(c *gin.Context) {
	var req ammeter.ChangeSwitchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeAmmeterSwitchSer(c, req.AmmeterId, req.UserId, req.Switch, req.Pwd)
}

func SetSwitchPwd(c *gin.Context) {
	var req ammeter.SetSwitchPwdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.SetSwitchPwdSer(c, req.AmmeterId, req.UserId, req.OldPwd, req.Pwd)
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
	var req ammeter.AmmeterWarningUpdateReq
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
	var req ammeter.ConfigUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.UpdateConfigSer(c, req)
}

func AddTestData(c *gin.Context) {
	var req ammeter.AmmeterDataReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.AddTestData(c, req.AmmeterId, req.Type, req.Value, req.CreateTime)
}
