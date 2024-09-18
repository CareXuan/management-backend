package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetMemberList(c *gin.Context, name, card string, page, pageSize int) {
	var members []*model.Member
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	if card != "" {
		sess.Where("card = ?", card)
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
	var records []model.MemberRecord
	err := conf.Mysql.Where("member_id = ?", memberId).Find(&records)
	if err != nil {
		common.ResError(c, "查询会员信息失败")
		return
	}
	if len(records) <= 0 {
		common.ResForbidden(c, "当前用户未查询到充值记录")
		return
	}
	var recordIds []int
	for _, record := range records {
		recordIds = append(recordIds, record.Id)
	}
	var timeRecordMapping = make(map[int]int)
	var monthlyRecordMapping = make(map[int]string)
	var recordDetail []model.MemberRecordDetail
	err = conf.Mysql.In("record_id", recordIds).Find(&recordDetail)
	if err != nil {
		common.ResError(c, "查询充值详情失败")
		return
	}
	for _, detail := range recordDetail {
		if detail.Type == model.RECHARGE_TYPE_TIMES {
			if _, ok := timeRecordMapping[detail.DeviceId]; !ok {
				timeRecordMapping[detail.DeviceId] = detail.Value
			} else {
				timeRecordMapping[detail.DeviceId] += detail.Value
			}
		}
		if detail.Type == model.RECHARGE_TYPE_MONTHLY {
			endTime := time.Unix(int64(detail.Value), 0).Format("2006-01-02 15:04:05")
			monthlyRecordMapping[detail.DeviceId] = endTime
		}
	}
	var memberUseRecords []model.MemberUseRecord
	err = conf.Mysql.Where("member_id = ?", memberId).Find(&memberUseRecords)
	if err != nil {
		common.ResError(c, "获取用户使用记录失败")
		return
	}
	var memberUses = make(map[int]int)
	for _, useRecord := range memberUseRecords {
		if _, ok := memberUses[useRecord.DeviceId]; !ok {
			memberUses[useRecord.DeviceId] = useRecord.Times
		} else {
			memberUses[useRecord.DeviceId] += useRecord.Times
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
		if _, ok := timeRecordMapping[device.Id]; ok {
			if _, useOk := memberUses[device.Id]; useOk {
				if timeRecordMapping[device.Id] > memberUses[device.Id] {
					timeRemain = timeRecordMapping[device.Id] - memberUses[device.Id]
				}
			} else {
				timeRemain = timeRecordMapping[device.Id]
			}
		}
		if _, ok := monthlyRecordMapping[device.Id]; ok {
			endTime, _ := time.Parse("2006-01-02 15:04:05", monthlyRecordMapping[device.Id])
			if endTime.After(time.Now()) {
				monthlyRemain = monthlyRecordMapping[device.Id]
			}
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
		insertData = append(insertData, model.MemberRecordDetail{
			RecordId:  recordId,
			DeviceId:  detail.DeviceId,
			Type:      detail.Type,
			Value:     detail.Value,
			CreatedAt: time.Now().Unix(),
		})
	}
	num, err := conf.Mysql.Insert(insertData)
	if err != nil {
		return -1, err
	}
	return num, nil
}
