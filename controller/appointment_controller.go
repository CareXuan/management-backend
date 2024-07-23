package controller

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/service"
	"strconv"
)

func GetAppointmentList(c *gin.Context) {
	deviceId := c.Query("device_id")
	deviceIdInt, _ := strconv.Atoi(deviceId)
	memberName := c.Query("member_name")
	memberPhone := c.Query("member_phone")
	memberCard := c.Query("member_card")
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("page_size")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetAppointmentListSer(c, deviceIdInt, memberName, memberPhone, memberCard, pageInt, pageSizeInt)
}

func GetAppointmentDetail(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GetAppointmentDetailSer(c, idInt)
}

func AddAppointment(c *gin.Context) {

}
