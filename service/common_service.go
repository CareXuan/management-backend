package service

import (
	"encoding/json"
	"env-backend/common"
	"env-backend/conf"
	"env-backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func IsSupervisor(c *gin.Context) (int, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return 0, fmt.Errorf("获取token失败")
	}
	part := strings.Split(authorization, " ")
	if len(part) < 2 {
		return 0, fmt.Errorf("非法请求")
	}
	var user model.User
	_, err := conf.Mysql.Where("token=?", part[1]).Get(&user)
	if err != nil {
		return 0, fmt.Errorf("获取用户信息失败")
	}
	var userRole model.UserRole
	_, err = conf.Mysql.Where("user_id = ?", user.Id).Get(&userRole)
	if err != nil {
		return 0, fmt.Errorf("获取用户角色失败")
	}
	var roleItem model.Role
	_, err = conf.Mysql.Where("id = ?", userRole.RoleId).Get(&roleItem)
	if err != nil {
		return 0, fmt.Errorf("获取角色信息失败")
	}
	if roleItem.Name != "管理员" {
		return 0, nil
	}
	return 1, nil
}

func WechatCodeSer(c *gin.Context, code string) {
	type WechatResult struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Openid       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	appId := "wxb334388d276b5763"
	secret := "e115cbf4f38b4f0686628b565e45a498"
	grantType := "authorization_code"
	res, err := common.DoGet("https://api.weixin.qq.com/sns/oauth2/access_token", map[string]string{
		"appid":      appId,
		"secret":     secret,
		"code":       code,
		"grant_type": grantType,
	})
	if err != nil {
		common.ResError(c, "发送请求失败")
		return
	}
	var wechatResult WechatResult
	_ = json.Unmarshal(res, &wechatResult)
	common.ResOk(c, "ok", wechatResult)
}

func WechatBindSer(c *gin.Context, req model.WechatBindReq) {
	var smsExist model.Sms
	has, err := conf.Mysql.Where("phone = ?", req.Phone).Where("used_at = 0").Where("expired_at > ?", time.Now().Unix()).Get(&smsExist)
	if err != nil {
		common.ResError(c, "检查验证码状态失败")
		return
	}
	if !has {
		common.ResForbidden(c, "请先发送验证码")
		return
	}
	if smsExist.SmsCode != req.SmsCode {
		common.ResForbidden(c, "验证码错误")
		return
	}
	_, err = conf.Mysql.Where("phone = ?", req.Phone).Where("used_at = 0").Where("expired_at > ?", time.Now().Unix()).Update(&model.Sms{
		UsedAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "使用验证码失败")
		return
	}
	_, err = conf.Mysql.Where("phone = ?", req.Phone).Update(&model.User{
		OpenId: req.Openid,
	})
	if err != nil {
		common.ResError(c, "修改用户信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func SmsCodeOne(c *gin.Context, req model.SmsOneReq) {
	var userExist model.User
	has, err := conf.Mysql.Where("phone = ?", req.Phone).Get(&userExist)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if !has {
		common.ResForbidden(c, "当前手机号并未关联商户")
		return
	}
	var smsExist model.Sms
	has, err = conf.Mysql.Where("phone = ?", req.Phone).Where("used_at = 0").Where("expired_at > ?", time.Now().Unix()).Get(&smsExist)
	if err != nil {
		common.ResError(c, "检查验证码状态失败")
		return
	}
	if has {
		common.ResForbidden(c, "验证码已发送，有效期五分钟")
		return
	}
	_, err = conf.Mysql.Insert(model.Sms{
		Phone:     req.Phone,
		SmsCode:   strconv.Itoa(rand.Intn(9000) + 1000),
		ExpiredAt: int(time.Now().Unix()) + 300,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "发送验证码失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
