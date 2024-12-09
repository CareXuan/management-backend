package model

type Sms struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Phone     string `json:"phone" xorm:"varchar(16) not null default '' comment('手机号')"`
	Code      string `json:"code" xorm:"varchar(8) not null default '' comment('验证码')"`
	UseAt     int    `json:"use_at" xorm:"int(10) not null default 0 comment('使用时间')"`
	ExpiredAt int    `json:"expired_at" xorm:"int(10) not null default 0 comment('过期时间')"`
	BizId     string `json:"biz_id" xorm:"varchar(128) not null default '' comment('biz_id')"`
	RequestId string `json:"request_id" xorm:"varchar(128) not null default '' comment('request_id')"`
	CreateAt  int    `json:"create_at" xorm:"int(10) not null default 0 comment('创建时间')"`
}

type SmsReq struct {
	Phone string `json:"phone"`
}
