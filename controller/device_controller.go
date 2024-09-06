package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-gpt-server/model"
	"my-gpt-server/service"
	"strconv"
)

func GetDeviceList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetDeviceListSer(c, name, pageInt, pageSizeInt)
}

func GetDeviceInfo(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	service.GetDeviceInfoSer(c, idInt)
}

func AddDevice(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		log.Fatal(err)
		return
	}
	service.AddDeviceSer(c, device)
}

// ================================== 套餐 ==================================

func GetPackageList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("pageSize")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetPackageListSer(c, name, pageInt, pageSizeInt)
}

func GetPackageInfo(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	service.GetPackageInfoSer(c, idInt)
}

func AddPackage(c *gin.Context) {
	var packageAdd model.DevicePackageAddReq
	if err := c.ShouldBindJSON(&packageAdd); err != nil {
		log.Fatal(err)
		return
	}
	service.AddPackageSer(c, packageAdd)
}

func PackageChangeStatus(c *gin.Context) {
	var packageChangeStatus model.DevicePackageChangeStatusReq
	if err := c.ShouldBindJSON(&packageChangeStatus); err != nil {
		log.Fatal(err)
		return
	}
	service.PackageChangeStatusSer(c, packageChangeStatus)
}
