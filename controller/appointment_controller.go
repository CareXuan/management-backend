package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-gpt-server/model"
	"my-gpt-server/service"
	"strconv"
	"strings"
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

func GetAppointmentChart(c *gin.Context) {
	ids := c.Query("ids")
	date := c.Query("date")
	idsArr := strings.Split(ids, ",")
	service.GetAppointmentChartSer(c, idsArr, date)
}

func AddAppointment(c *gin.Context) {
	var addAppointmentReq model.AddAppointmentReq
	if err := c.ShouldBindJSON(&addAppointmentReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddAppointmentSer(c, addAppointmentReq)
}

func VerifyAppointment(c *gin.Context) {
	var verifyAppointmentReq model.VerifyAppointmentReq
	if err := c.ShouldBindJSON(&verifyAppointmentReq); err != nil {
		log.Fatal(err)
		return
	}
	service.VerifyAppointmentSer(c, verifyAppointmentReq)
}

func UniappAppointment(c *gin.Context) {
	memberId := c.Query("member_id")
	memberIdInt, _ := strconv.Atoi(memberId)
	searchType := c.Query("type")
	page := c.Query("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.Query("page_size")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.UniappAppointmentSer(c, memberIdInt, pageInt, pageSizeInt, searchType)
}

func UniappAppointmentDevice(c *gin.Context) {
	memberId := c.Query("member_id")
	memberIdInt, _ := strconv.Atoi(memberId)
	service.UniappAppointmentDeviceSer(c, memberIdInt)
}

func UniappAppointmentAdd(c *gin.Context) {
	var addAppointmentReq model.AddAppointmentReq
	if err := c.ShouldBindJSON(&addAppointmentReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddAppointmentSer(c, addAppointmentReq)
}

func UniappAppointmentCancel(c *gin.Context) {
	var uniappAppointmentCancel model.UniappAppointmentCancelReq
	if err := c.ShouldBindJSON(&uniappAppointmentCancel); err != nil {
		log.Fatal(err)
		return
	}
	service.UniappAppointmentCancelSer(c, uniappAppointmentCancel)
}
