package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"net/http"
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
			Src:           device.Pic,
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

func UniappLoginSer(c *gin.Context, req model.UniappLoginReq) {
	// 请求微信 API 获取 openid 和 session_key
	resp, err := http.Get(fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		conf.Conf.Wechat.AppId, conf.Conf.Wechat.AppSecret, req.Code,
	))
	if err != nil {
		common.ResError(c, "请求微信失败")
		return
	}
	defer resp.Body.Close()

	// 解析微信 API 返回的 JSON 数据
	body, _ := ioutil.ReadAll(resp.Body)
	var weChatResp model.WeChatLoginResponse
	err = json.Unmarshal(body, &weChatResp)
	if err != nil {
		common.ResError(c, "解析返回数据失败")
		return
	}

	// 处理错误
	if weChatResp.ErrCode != 0 {
		common.ResError(c, "处理错误")
		return
	}

	var wechatUser model.MemberWechat
	_, err = conf.Mysql.Where("open_id = ?", weChatResp.OpenID).Get(&wechatUser)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if wechatUser.Id == 0 {
		userAdd := model.Member{
			Name:      "微信用户" + common.GetOneNewCard(12),
			Card:      common.GetOneNewCard(12),
			Phone:     "",
			Emergency: "",
			Birthday:  "",
			Gender:    0,
			Remark:    "",
			CreatedAt: time.Now().Unix(),
		}
		_, err := conf.Mysql.Insert(&userAdd)
		if err != nil {
			common.ResError(c, "添加用户失败")
			return
		}
		wechatUserAdd := model.MemberWechat{
			MemberId:  userAdd.Id,
			OpenId:    weChatResp.OpenID,
			CreatedAt: time.Now().Unix(),
		}
		_, err = conf.Mysql.Insert(&wechatUserAdd)
		if err != nil {
			common.ResError(c, "添加微信用户失败")
			return
		}
		common.ResOk(c, "ok", userAdd)
		return
	}

	var userItem model.Member
	_, err = conf.Mysql.Where("id = ?", wechatUser.MemberId).Get(&userItem)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}

	common.ResOk(c, "ok", userItem)
}

func UniappUpdateSer(c *gin.Context, req model.UniappUpdateReq) {
	_, err := conf.Mysql.Where("id = ?", req.MemberId).Update(model.Member{
		Name:      req.Name,
		Emergency: req.Emergency,
		Birthday:  req.Birthday,
		Gender:    req.Gender,
		Remark:    req.Remark,
	})
	if err != nil {
		common.ResError(c, "修改用户信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func UniappPhoneBindSer(c *gin.Context, req model.UniappPhoneBindReq) {
	var nowSms model.Sms
	has, err := conf.Mysql.Where("phone = ?", req.Phone).Where("use_at = 0").Where("expired_at > ?", time.Now().Unix()).OrderBy("create_at DESC").Get(&nowSms)
	if err != nil {
		common.ResError(c, "获取验证码信息失败")
		return
	}
	if !has {
		common.ResForbidden(c, "请先发送验证码")
		return
	}
	_, err = conf.Mysql.Where("phone = ?", req.Phone).Where("use_at = 0").Update(model.Sms{
		UseAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "修改使用时间失败")
		return
	}
	if nowSms.Code != req.Code {
		common.ResForbidden(c, "验证码错误")
		return
	}
	var existMember model.Member
	has, err = conf.Mysql.Where("phone = ?", req.Phone).Get(&existMember)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if has {
		var wechatMember model.MemberWechat
		_, err := conf.Mysql.Where("member_id = ?", req.MemberId).Get(&wechatMember)
		if err != nil {
			common.ResError(c, "获取微信绑定用户失败")
			return
		}
		sess := conf.Mysql.NewSession()
		_, err = sess.Where("open_id = ?", wechatMember.OpenId).Update(model.MemberWechat{
			MemberId: existMember.Id,
		})
		if err != nil {
			common.ResError(c, "修改用户信息失败")
			sess.Rollback()
			return
		}
		_, err = sess.Where("id = ?", req.MemberId).Delete(&model.Member{})
		if err != nil {
			common.ResError(c, "删除微信绑定用户失败")
			sess.Rollback()
			return
		}
		sess.Commit()
	} else {
		_, err := conf.Mysql.Where("id = ?", req.MemberId).Update(model.Member{
			Phone: req.Phone,
		})
		if err != nil {
			common.ResError(c, "修改用户信息失败")
			return
		}
	}
	var resUser model.Member
	_, err = conf.Mysql.Where("phone = ?", req.Phone).Get(&resUser)
	if err != nil {
		common.ResError(c, "获取用户失败")
		return
	}
	common.ResOk(c, "ok", resUser)
}

func UniappInfoSer(c *gin.Context, memberId int) {
	var memberInfo model.Member
	_, err := conf.Mysql.Where("id = ?", memberId).Get(&memberInfo)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	common.ResOk(c, "ok", memberInfo)
}
