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
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: count})
}

func GetMemberDetail(c *gin.Context, memberId int) {
	var member model.Member
	_, err := conf.Mysql.Where("id = ?", memberId).Get(&member)
	if err != nil {
		common.ResError(c, "查询会员详情失败")
		return
	}
	common.ResOk(c, "ok", memberId)
}

func AddMember(c *gin.Context, member model.MemberAddReq) {
	newMember := model.Member{
		Name:      member.Name,
		Card:      common.GetOneNewCard(24),
		Phone:     member.Phone,
		Emergency: member.Emergency,
		Birthday:  member.Birthday,
		Gender:    member.Gender,
		Remark:    member.UserRemark,
		CreatedAt: time.Now().Unix(),
	}
	_, err := conf.Mysql.Insert(&newMember)
	if err != nil {
		common.ResError(c, "添加新会员失败")
		return
	}
	newMemberRecord := model.MemberRecord{
		CardId:    newMember.Id,
		PackageId: member.PackageId,
		Type:      member.Type,
		Price:     member.Price,
		Cost:      member.Cost,
		Remark:    member.RechargeRemark,
		Pic:       member.Pic,
		CreatedAt: time.Now().Unix(),
	}
	_, err = conf.Mysql.Insert(&newMemberRecord)
	if err != nil {
		common.ResError(c, "添加会员记录失败")
		return
	}
	_, err = addRechargeDetail(newMemberRecord.Id, member.RechargeDetail)
	if err != nil {
		common.ResError(c, "添加充值记录详情失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func Recharge(c *gin.Context, memberRecharge model.MemberRechargeReq) {
	newMemberRecord := model.MemberRecord{
		CardId:    memberRecharge.CardId,
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
