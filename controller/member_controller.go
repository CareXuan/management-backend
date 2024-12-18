package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-gpt-server/model"
	"my-gpt-server/service"
	"strconv"
)

func GetMemberList(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	card := c.Query("card")
	name := c.Query("name")
	phone := c.Query("phone")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetMemberList(c, name, card, phone, pageInt, pageSizeInt)
}

func GetMemberDetail(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GetMemberDetail(c, idInt)
}

func GetMemberRechargeDetail(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GetMemberRechargeDetail(c, idInt)
}

func AddMember(c *gin.Context) {
	var member model.MemberAddReq
	if err := c.ShouldBindJSON(&member); err != nil {
		log.Fatal(err)
		return
	}
	service.AddMember(c, member)
}

func MemberRecharge(c *gin.Context) {
	var memberRecharge model.MemberRechargeReq
	if err := c.ShouldBindJSON(&memberRecharge); err != nil {
		log.Fatal(err)
		return
	}
	service.Recharge(c, memberRecharge)
}

func UniappLogin(c *gin.Context) {
	var req model.UniappLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.UniappLoginSer(c, req)
}

func UniappUpdate(c *gin.Context) {
	var req model.UniappUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.UniappUpdateSer(c, req)
}

func UniappPhoneBind(c *gin.Context) {
	var req model.UniappPhoneBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Fatal(err)
		return
	}
	service.UniappPhoneBindSer(c, req)
}

func UniappInfo(c *gin.Context) {
	memberId := c.Query("member_id")
	memberIdInt, _ := strconv.Atoi(memberId)
	service.UniappInfoSer(c, memberIdInt)
}

func UniappGetMemberRechargeDetail(c *gin.Context) {
	id := c.Query("member_id")
	idInt, _ := strconv.Atoi(id)
	service.GetMemberRechargeDetail(c, idInt)
}
