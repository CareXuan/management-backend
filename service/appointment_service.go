package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"strconv"
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

func GetAppointmentChartSer(c *gin.Context, appointmentIds []string, date string) {
	var allAppointment []*model.MemberUseRecord
	location, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" 00:00:00", location)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" 23:59:59", location)
	if len(appointmentIds) <= 0 || appointmentIds[0] == "" {
		var allDevices []*model.Device
		err := conf.Mysql.Find(&allDevices)
		if err != nil {
			common.ResError(c, "获取设备信息失败")
			return
		}
		for _, i := range allDevices {
			appointmentIds = append(appointmentIds, strconv.Itoa(i.Id))
		}
	}
	err := conf.Mysql.In("device_id", appointmentIds).In("status", []int{model.MEMBER_USE_STATUS_PASS, model.MEMBER_USE_STATUS_WAITING}).Where("start_time > ?", startTime.Unix()).Where("end_time < ?", endTime.Unix()).OrderBy("start_time").Find(&allAppointment)
	if err != nil {
		common.ResError(c, "获取预约信息失败")
		return
	}
	var memberIds []int
	var deviceIds []int
	for _, i := range allAppointment {
		memberIds = append(memberIds, i.MemberId)
		deviceIds = append(deviceIds, i.DeviceId)
	}
	var memberItems []*model.Member
	err = conf.Mysql.In("id", memberIds).Find(&memberItems)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	var memberMapping = make(map[int]*model.Member)
	for _, i := range memberItems {
		memberMapping[i.Id] = i
	}
	var deviceItems []*model.Device
	var deviceMapping = make(map[int]*model.Device)
	err = conf.Mysql.In("id", deviceIds).Find(&deviceItems)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	for _, i := range deviceItems {
		deviceMapping[i.Id] = i
	}
	var res = make(map[int]*model.MemberAppointmentData)
	for _, i := range allAppointment {
		if _, ok := res[i.DeviceId]; !ok {
			res[i.DeviceId] = &model.MemberAppointmentData{
				DeviceId:   i.DeviceId,
				DeviceName: deviceMapping[i.DeviceId].Name,
				Data:       []model.AppointmentDataTimes{},
			}
		}
		res[i.DeviceId].Data = append(res[i.DeviceId].Data, model.AppointmentDataTimes{
			StartTime: i.StartTime,
			Duration:  (float64(i.EndTime) - float64(i.StartTime)) / 3600.0,
			EndTime:   i.EndTime,
			Member: model.AppointmentMember{
				MemberId:     i.MemberId,
				MemberName:   memberMapping[i.MemberId].Name,
				MemberColor:  "#000000",
				MemberRemark: memberMapping[i.MemberId].Remark,
			},
		})
	}
	var response []*model.MemberAppointmentData
	for _, i := range res {
		response = append(response, i)
	}
	common.ResOk(c, "ok", response)
}

func AddAppointmentSer(c *gin.Context, req model.AddAppointmentReq) {
	var deviceItem model.Device
	_, err := conf.Mysql.Where("id=?", req.DeviceId).Get(&deviceItem)
	if err != nil {
		common.ResError(c, "获取设备详情失败")
		return
	}
	if deviceItem.Status == model.DEVICE_STATUS_CLOSE {
		common.ResForbidden(c, "当前设备不可用，请联系商家")
		return
	}
	var deviceRecord []*model.MemberDeviceRecord
	err = conf.Mysql.Where("device_id = ?", req.DeviceId).Where("member_id = ?", req.MemberId).Find(&deviceRecord)
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

	location, _ := time.LoadLocation("Asia/Shanghai")
	startTimeInt, _ := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, location)
	endTimeInt, _ := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, location)

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
		Reason: req.Reason,
	})
	if err != nil {
		common.ResError(c, "修改预约状态失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
