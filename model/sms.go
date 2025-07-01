package model

type Sms struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Phone     string `json:"phone" xorm:"VARCHAR(16) not null default '' comment('手机号')"`
	SmsCode   string `json:"sms_code" xorm:"VARCHAR(8) not null default '' comment('验证码')"`
	ExpiredAt int    `json:"expired_at" xorm:"INT(11) not null default 0 comment('过期时间')"`
	UsedAt    int    `json:"used_at" xorm:"INT(11) not null default 0 comment('使用时间')"`
	CreatedAt int    `json:"-" xorm:"INT(11) not null default 0"`
	UpdatedAt int    `json:"-" xorm:"INT(11) not null default 0"`
	DeletedAt int    `json:"-" xorm:"INT(11) not null default 0"`
}

type SmsOneReq struct {
	Phone string `json:"phone"`
}

type WechatBindReq struct {
	Phone   string `json:"phone"`
	Openid  string `json:"openid"`
	SmsCode string `json:"sms_code"`
}
