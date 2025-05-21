package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/model"
	"management-backend/service"
	"strconv"
)

func GetDeviceList(c *gin.Context) {
	deviceId := c.Query("device_id")
	name := c.Query("name")
	phone := c.Query("phone")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.DeviceListSer(c, deviceId, name, phone, pageInt, pageSizeInt)
}

func AddOneDevice(c *gin.Context) {
	var device model.DeviceAddReq
	if err := c.ShouldBindJSON(&device); err != nil {
		log.Fatal(err)
		return
	}
	service.AddDeviceSer(c, device)
}

func UpdateSpecialInfo(c *gin.Context) {
	var device model.UpdateSpecialInfoReq
	if err := c.ShouldBindJSON(&device); err != nil {
		log.Fatal(err)
		return
	}
	service.UpdateSpecialInfoSer(c, device)
}

func ReadSpecialInfo(c *gin.Context) {
	var device model.ReadSpecialInfoReq
	if err := c.ShouldBindJSON(&device); err != nil {
		log.Fatal(err)
		return
	}
	service.ReadSpecialInfoSer(c, device)
}

func GetSpecialInfoLog(c *gin.Context) {
	deviceId := c.Query("device_id")
	service.GetSpecialInfoLogSer(c, deviceId)
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

func GetOneDeviceNewServiceData(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	deviceId := c.Query("device_id")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceNewServiceDataSer(c, deviceIdInt, pageInt, pageSizeInt)
}

func GetAllDeviceLocation(c *gin.Context) {
	service.AllDeviceLocationSer(c)
}

func GetDeviceStatistic(c *gin.Context) {
	startTime := c.Query("start_at")
	endTime := c.Query("end_at")
	deviceId := c.Query("device_id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceStatisticSer(c, deviceIdInt, startTime, endTime)
}

func GetAllWarning(c *gin.Context) {
	service.GetDeviceAllWarningSer(c)
}

func GetSingleWarning(c *gin.Context) {
	deviceId := c.Query("device_id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.GetDeviceSingleWarningSer(c, deviceIdInt)
}

func GetDeviceLocationHistory(c *gin.Context) {
	deviceId := c.Query("id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	service.DeviceLocationHistorySer(c, deviceIdInt)
}

func DeviceReport(c *gin.Context) {
	var req model.DeviceReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.DeviceReportSer(c, req.ReportType, req.DeviceId, req.ControlType, req.Msg)
}
