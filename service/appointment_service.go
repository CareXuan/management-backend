package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetAppointmentListSer(c *gin.Context, deviceId int, memberName, memberPhone, memberCard string, page, pageSize int) {
	var appointments []*model.MemberUseRecord
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
	var resDeviceIds []int
	for _, appointment := range appointments {
		resDeviceIds = append(resDeviceIds, appointment.DeviceId)
	}
	var resDevices []*model.Device
	err = conf.Mysql.In("id", resDeviceIds).Find(&resDevices)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	var deviceMapping = make(map[int]*model.Device)
	for _, device := range resDevices {
		deviceMapping[device.Id] = device
	}
	for _, appointment := range appointments {
		appointment.Member = memberMapping[appointment.MemberId]
		appointment.Device = deviceMapping[appointment.DeviceId]
	}

	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: appointments})
}

func GetAppointmentDetailSer(c *gin.Context, appointmentId int) {
	var appointment model.MemberUseRecord
	_, err := conf.Mysql.Where("id=?", appointmentId).Get(&appointment)
	if err != nil {
		common.ResError(c, "获取预约详情失败")
		return
	}
	var member model.Member
	_, err = conf.Mysql.Where("id=?", appointment.MemberId).Get(&member)
	if err != nil {
		common.ResError(c, "获取预约用户失败")
		return
	}
	var device model.Device
	_, err = conf.Mysql.Where("id=?", appointment.DeviceId).Get(&device)
	if err != nil {
		common.ResError(c, "获取预约设备失败")
		return
	}
	appointment.Member = &member
	appointment.Device = &device
	common.ResOk(c, "ok", appointment)
}

func AddAppointmentSer(c *gin.Context, req model.AddAppointmentReq) {
	var deviceRecord []*model.MemberDeviceRecord
	err := conf.Mysql.Where("device_id = ?", req.DeviceId).Where("member_id = ?", req.MemberId).Find(&deviceRecord)
	if err != nil {
		common.ResError(c, "获取设备储值记录失败")
		return
	}
	timesRemain := 0
	monthlyRemain := 0
	for _, i := range deviceRecord {
		if i.Type == model.RECHARGE_TYPE_TIMES {
			timesRemain += i.Value
		}
		if i.Type == model.RECHARGE_TYPE_MONTHLY {
			monthlyRemain = i.Value
		}
	}
	var memberUseRecord []*model.MemberUseRecord
	err = conf.Mysql.Where("device_id = ?", req.DeviceId).Where("member_id = ?", req.MemberId).In("status", []int{
		model.MEMBER_USE_STATUS_WAITING,
		model.MEMBER_USE_STATUS_PASS,
	}).Find(&memberUseRecord)
	if err != nil {
		common.ResError(c, "获取设备使用记录失败")
		return
	}
	if len(memberUseRecord) >= timesRemain && int(time.Now().Unix()) > monthlyRemain {
		common.ResForbidden(c, "剩余使用次数不足")
		return
	}

	startTimeInt, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	endTimeInt, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)

	var deviceUseExist []*model.MemberUseRecord
	var memberUseExist []*model.MemberUseRecord
	err = conf.Mysql.Where("device_id = ?", req.DeviceId).Where("start_time <= ?", endTimeInt.Unix()).Where("end_time >= ?", startTimeInt.Unix()).Find(&deviceUseExist)
	if err != nil {
		common.ResError(c, "获取设备使用情况失败")
		return
	}
	err = conf.Mysql.Where("member_id = ?", req.MemberId).Where("start_time <= ?", endTimeInt.Unix()).Where("end_time >= ?", startTimeInt.Unix()).Find(&memberUseExist)
	if err != nil {
		common.ResError(c, "获取用户预约情况失败")
		return
	}
	if len(deviceUseExist) > 0 {
		common.ResForbidden(c, "当前时间内有其他人预约，请更换预约时间")
		return
	}
	if len(memberUseExist) > 0 {
		common.ResForbidden(c, "您已有其他预约，请更换预约时间")
		return
	}

	_, err = conf.Mysql.Insert(model.MemberUseRecord{
		MemberId:  req.MemberId,
		DeviceId:  req.DeviceId,
		Times:     1,
		StartTime: int(startTimeInt.Unix()),
		EndTime:   int(endTimeInt.Unix()),
		Remark:    req.Remark,
		Pic:       req.Pic,
		Status:    model.MEMBER_USE_STATUS_WAITING,
		CreatedAt: time.Now().Unix(),
	})
	if err != nil {
		common.ResError(c, "写入预约信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func VerifyAppointmentSer(c *gin.Context, req model.VerifyAppointmentReq) {
	_, err := conf.Mysql.Where("id=?", req.AppointmentId).Update(&model.MemberUseRecord{
		Status: req.Status,
	})
	if err != nil {
		common.ResError(c, "修改预约状态失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
