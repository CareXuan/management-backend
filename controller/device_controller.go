package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/model"
	"management-backend/service"
	"strconv"
)

func GetDeviceList(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.DeviceListSer(c, pageInt, pageSizeInt)
}

func AddOneDevice(c *gin.Context) {
	var device model.DeviceAddReq
	if err := c.ShouldBindJSON(&device); err != nil {
		log.Fatal(err)
		return
	}
	service.AddDeviceSer(c, device)
}

func GetOneDeviceInfo(c *gin.Context) {
	deviceId := c.Query("device_id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceInfoSer(c, deviceIdInt)
}

func GetOneDeviceCommonData(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	deviceId := c.Query("device_id")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceCommonDataSer(c, deviceIdInt, pageInt, pageSizeInt)
}

func GetOneDeviceServiceData(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	deviceId := c.Query("device_id")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceServiceDataSer(c, deviceIdInt, pageInt, pageSizeInt)
}

func GetAllDeviceLocation(c *gin.Context) {
	service.AllDeviceLocationSer(c)
}

func GetDeviceLocationHistory(c *gin.Context) {
	deviceId := c.Query("id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceLocationHistorySer(c, deviceIdInt)
}
