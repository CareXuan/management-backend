package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetMemberList(c *gin.Context, name, card, phone string, page, pageSize int) {
	var members []*model.Member
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	if card != "" {
		sess.Where("card = ?", card)
	}
	if phone != "" {
		sess.Where("phone = ?", phone)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&members)
	if err != nil {
		common.ResError(c, "查询会员列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: members})
}

func GetMemberDetail(c *gin.Context, memberId int) {
	var member model.Member
	_, err := conf.Mysql.Where("id = ?", memberId).Get(&member)
	if err != nil {
		common.ResError(c, "查询会员详情失败")
		return
	}
	common.ResOk(c, "ok", member)
}

func GetMemberRechargeDetail(c *gin.Context, memberId int) {
	var deviceRecords []model.MemberDeviceRecord
	err := conf.Mysql.Where("member_id = ?", memberId).Find(&deviceRecords)
	if err != nil {
		common.ResError(c, "获取余额失败")
		return
	}
	var timesRecordMapping = make(map[int]int)
	var monthlyRecordMapping = make(map[int]int)
	for _, deviceRecord := range deviceRecords {
		if deviceRecord.Type == model.RECHARGE_TYPE_TIMES {
			timesRecordMapping[deviceRecord.DeviceId] = deviceRecord.Value
		}
		if deviceRecord.Type == model.RECHARGE_TYPE_MONTHLY {
			monthlyRecordMapping[deviceRecord.DeviceId] = deviceRecord.Value
		}
	}
	var devices []model.Device
	err = conf.Mysql.Where("created_at > 0").Find(&devices)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	var recordRes []model.MemberRechargeRecordRes
	for _, device := range devices {
		timeRemain := 0
		monthlyRemain := ""
		if _, ok := timesRecordMapping[device.Id]; ok && timesRecordMapping[device.Id] > 0 {
			timeRemain = timesRecordMapping[device.Id]
		}
		if _, ok := monthlyRecordMapping[device.Id]; ok && monthlyRecordMapping[device.Id] > int(time.Now().Unix()) {
			monthlyRemain = time.Unix(int64(monthlyRecordMapping[device.Id]), 0).Format("2006-01-02 15:04:05")
		}
		recordRes = append(recordRes, model.MemberRechargeRecordRes{
			Name:          device.Name,
			TimesRemain:   timeRemain,
			MonthlyRemain: monthlyRemain,
		})
	}
	common.ResOk(c, "ok", recordRes)
}

func AddMember(c *gin.Context, member model.MemberAddReq) {
	var memberExist model.Member
	_, err := conf.Mysql.Where("phone = ?", member.Phone).Get(&memberExist)
	if err != nil {
		common.ResError(c, "获取会员信息失败")
		return
	}
	if memberExist.Id != 0 {
		common.ResForbidden(c, "当前手机号已创建会员")
		return
	}
	newMember := model.Member{
		Name:      member.Name,
		Card:      common.GetOneNewCard(12),
		Phone:     member.Phone,
		Emergency: member.Emergency,
		Birthday:  member.Birthday,
		Gender:    member.Gender,
		Remark:    member.UserRemark,
		CreatedAt: time.Now().Unix(),
	}
	_, err = conf.Mysql.Insert(&newMember)
	if err != nil {
		common.ResError(c, "添加新会员失败")
		return
	}
	//newMemberRecord := model.MemberRecord{
	//	CardId:    newMember.Id,
	//	PackageId: member.PackageId,
	//	Type:      member.Type,
	//	Price:     member.Price,
	//	Cost:      member.Cost,
	//	Remark:    member.RechargeRemark,
	//	Pic:       member.Pic,
	//	CreatedAt: time.Now().Unix(),
	//}
	//_, err = conf.Mysql.Insert(&newMemberRecord)
	//if err != nil {
	//	common.ResError(c, "添加会员记录失败")
	//	return
	//}
	//_, err = addRechargeDetail(newMemberRecord.Id, member.RechargeDetail)
	//if err != nil {
	//	common.ResError(c, "添加充值记录详情失败")
	//	return
	//}
	common.ResOk(c, "ok", nil)
}

func Recharge(c *gin.Context, memberRecharge model.MemberRechargeReq) {
	_, errMsg := updateDeviceRecord(memberRecharge.MemberId, memberRecharge.RechargeDetail)
	if errMsg != "" {
		common.ResForbidden(c, errMsg)
		return
	}
	newMemberRecord := model.MemberRecord{
		MemberId:  memberRecharge.MemberId,
		PackageId: memberRecharge.PackageId,
		Type:      memberRecharge.Type,
		Price:     memberRecharge.Price,
		Cost:      memberRecharge.Cost,
		Remark:    memberRecharge.Remark,
		Pic:       memberRecharge.Pic,
		CreatedAt: time.Now().Unix(),
	}
	_, err := conf.Mysql.Insert(&newMemberRecord)
	if err != nil {
		common.ResError(c, "添加充值记录失败")
		return
	}
	_, err = addRechargeDetail(newMemberRecord.Id, memberRecharge.RechargeDetail)
	if err != nil {
		common.ResError(c, "添加充值记录详情失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func addRechargeDetail(recordId int, rechargeDetail []model.MemberRecordDetailAdd) (int64, error) {
	var insertData []model.MemberRecordDetail
	for _, detail := range rechargeDetail {
		value := detail.Value
		if detail.Type == model.RECHARGE_TYPE_MONTHLY {
			value = detail.Value * 86400
		}
		insertData = append(insertData, model.MemberRecordDetail{
			RecordId:  recordId,
			DeviceId:  detail.DeviceId,
			Type:      detail.Type,
			Value:     value,
			CreatedAt: time.Now().Unix(),
		})
	}
	num, err := conf.Mysql.Insert(insertData)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func updateDeviceRecord(memberId int, rechargeDetail []model.MemberRecordDetailAdd) (int64, string) {
	var memberDeviceRecords []model.MemberDeviceRecord
	err := conf.Mysql.Where("member_id = ?", memberId).Find(&memberDeviceRecords)
	if err != nil {
		return -1, "查询用户剩余使用记录失败"
	}
	var timesRecord = make(map[int]int)
	var monthlyRecord = make(map[int]int)
	for _, record := range memberDeviceRecords {
		if record.Type == model.RECHARGE_TYPE_TIMES {
			timesRecord[record.DeviceId] = record.Value
		}
		if record.Type == model.RECHARGE_TYPE_MONTHLY {
			monthlyRecord[record.DeviceId] = record.Value
		}
	}

	sess := conf.Mysql.NewSession()
	if err := sess.Begin(); err != nil {
		return -1, "开启事务失败"
	}
	for _, detail := range rechargeDetail {
		if detail.Type == model.RECHARGE_TYPE_TIMES {
			if _, ok := timesRecord[detail.DeviceId]; ok {
				if monthlyRecord[detail.DeviceId] > 0 {
					return -1, "已有包月的设备无法添加计次"
				}
				_, err := conf.Mysql.Where("member_id = ?", memberId).Where("device_id = ?", detail.DeviceId).Where("type = ?", model.RECHARGE_TYPE_TIMES).Update(model.MemberDeviceRecord{
					Value: timesRecord[detail.DeviceId] + detail.Value,
				})
				if err != nil {
					sess.Rollback()
					return -1, "写入数据失败"
				}
			} else {
				_, err := conf.Mysql.Insert(model.MemberDeviceRecord{
					MemberId:  memberId,
					DeviceId:  detail.DeviceId,
					Type:      model.RECHARGE_TYPE_TIMES,
					Value:     detail.Value,
					CreatedAt: time.Now().Unix(),
				})
				if err != nil {
					sess.Rollback()
					return -1, "写入数据失败"
				}
			}
		}
		if detail.Type == model.RECHARGE_TYPE_MONTHLY {
			if _, ok := monthlyRecord[detail.DeviceId]; ok && monthlyRecord[detail.DeviceId] > int(time.Now().Unix()) {
				_, err := conf.Mysql.Where("member_id = ?", memberId).Where("device_id = ?", detail.DeviceId).Where("type = ?", model.RECHARGE_TYPE_MONTHLY).Update(model.MemberDeviceRecord{
					Value: monthlyRecord[detail.DeviceId] + detail.Value*86400,
				})
				if err != nil {
					sess.Rollback()
					return -1, "写入数据失败"
				}
			} else {
				_, err := conf.Mysql.Insert(model.MemberDeviceRecord{
					MemberId:  memberId,
					DeviceId:  detail.DeviceId,
					Type:      model.RECHARGE_TYPE_MONTHLY,
					Value:     int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.Now().Location()).Unix()) + detail.Value*86400,
					CreatedAt: time.Now().Unix(),
				})
				if err != nil {
					sess.Rollback()
					return -1, "写入数据失败"
				}
			}
		}
	}
	if err := sess.Commit(); err != nil {
		sess.Rollback()
		return -1, "提交事务失败"
	}
	return 0, ""
}