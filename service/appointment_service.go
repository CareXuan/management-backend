package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
)

func GetAppointmentListSer(c *gin.Context, deviceId int, memberName, memberPhone, memberCard string, page, pageSize int) {
	var appointments []*model.Appointment
	sess := conf.Mysql.NewSession()
	if deviceId != 0 {
		sess.Where("device_id=?", deviceId)
	}
	var memberIds []int
	if memberName != "" || memberPhone != "" || memberCard != "" {
		var members []*model.Member
		memberSess := conf.Mysql.NewSession()
		if memberName != "" {
			memberSess.Where("name like ?", "%"+memberName+"%")
		}
		if memberPhone != "" {
			memberSess.Where("phone like ?", "%"+memberPhone+"%")
		}
		if memberCard != "" {
			memberSess.Where("card like ?", "%"+memberCard+"%")
		}
		err := memberSess.Find(&members)
		if err != nil {
			common.ResError(c, "获取用户信息失败")
			return
		}
		for _, member := range members {
			memberIds = append(memberIds, member.Id)
		}
	}
	if len(memberIds) > 0 {
		sess.In("member_id", memberIds)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&appointments)
	if err != nil {
		common.ResError(c, "获取预约信息失败")
		return
	}
	var resMemberIds []int
	for _, appointment := range appointments {
		resMemberIds = append(resMemberIds, appointment.MemberId)
	}
	var resMembers []*model.Member
	err = conf.Mysql.In("id", resMemberIds).Find(&resMembers)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	var memberMapping = make(map[int]*model.Member)
	for _, member := range resMembers {
		memberMapping[member.Id] = member
	}
	for _, appointment := range appointments {
		appointment.Member = memberMapping[appointment.MemberId]
	}

	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: appointments})
}

func GetAppointmentDetailSer(c *gin.Context, appointmentId int) {
	var appointment model.Appointment
	_, err := conf.Mysql.Where("id=?", appointmentId).Get(&appointment)
	if err != nil {
		common.ResError(c, "获取预约详情失败")
		return
	}
	var member *model.Member
	_, err = conf.Mysql.Where("id=?", appointment.MemberId).Get(member)
	if err != nil {
		common.ResError(c, "获取预约用户失败")
		return
	}
	appointment.Member = member
	common.ResOk(c, "ok", appointment)
}
