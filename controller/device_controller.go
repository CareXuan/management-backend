package controller

import (
	"github.com/gin-gonic/gin"
	"management-backend/service"
	"strconv"
)

func DeviceList(c *gin.Context) {
	deviceType := c.Query("type")
	name := c.Query("name")
	iccid := c.Query("iccid")
	status := c.Query("status")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.DeviceListSer(c, pageInt, pageSizeInt, name, iccid, deviceType, status)
}

func SignalDetailList(c *gin.Context) {
	deviceId := c.Param("device_id")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.SignalDetailListSer(c, pageInt, pageSizeInt, deviceId)
}
