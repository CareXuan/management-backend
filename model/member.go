package model

const RECHARGE_TYPE_TIMES = 1
const RECHARGE_TYPE_MONTHLY = 2

const USE_STATUS_RESERVATION = 1
const USE_STATUS_FINISH = 2

const MEMBER_GENDER_MALE = 1
const MEMBER_GENDER_FEMALE = 2

const MEMBER_RECORD_TYPE_RECHARGE = 1
const MEMBER_RECORD_TYPE_REFUND = 2

const MEMBER_USE_STATUS_WAITING = 1
const MEMBER_USE_STATUS_PASS = 2
const MEMBER_USE_STATUS_REFUSE = 3

type Member struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"VARCHAR(64) not null default '' comment('姓名')"`
	Card      string `json:"card" xorm:"VARCHAR(128) not null default '' comment('卡号')"`
	Phone     string `json:"phone" xorm:"VARCHAR(16) not null default '' comment('手机号')"`
	Emergency string `json:"emergency" xorm:"VARCHAR(16) not null default '' comment('紧急联系人')"`
	Birthday  string `json:"birthday" xorm:"VARCHAR(16) not null default '' comment('生日')"`
	Gender    int    `json:"gender" xorm:"INT(3) not null default 0 comment('性别 1：男 2：女')"`
	Remark    string `json:"remark" xorm:"VARCHAR(256) not null default '' comment('备注')"`
	CreatedAt int64  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberRecord struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	MemberId  int    `json:"member_id" xorm:"INT(11) not null"`
	Type      int    `json:"type" xorm:"INT(11) not null default 0 comment('记录类型 1：充值 2：退款')"`
	Price     int    `json:"price" xorm:"INT(11) not null default 0 comment('原价')"`
	Cost      int    `json:"cost" xorm:"INT(11) not null default 0 comment('实际支付价格')"`
	PackageId int    `json:"package_id" xorm:"INT(11) not null default 0 comment('套餐ID 如未使用套餐为0')"`
	Remark    string `json:"remark" xorm:"VARCHAR(256) not null default '' comment('备注')"`
	Pic       string `json:"pic" xorm:"VARCHAR(512) not null default '' comment('图片')"`
	CreatedAt int64  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberRecordDetail struct {
	Id        int   `json:"id" xorm:"pk autoincr INT(11)"`
	RecordId  int   `json:"record_id" xorm:"INT(11) not null comment('充值记录ID')"`
	DeviceId  int   `json:"device_id" xorm:"INT(10) not null comment('设备ID')"`
	Type      int   `json:"type" xorm:"INT(3) not null default 0 comment('类型 1：计次 2：计时')"`
	Value     int   `json:"value" xorm:"INT(11) not null default 0 comment('type=1 充值次数 type=2 结束时间')"`
	CreatedAt int64 `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberUseRecord struct {
	Id        int     `json:"id" xorm:"pk autoincr INT(11)"`
	MemberId  int     `json:"member_id" xorm:"INT(11) not null"`
	DeviceId  int     `json:"device_id" xorm:"INT(11) not null"`
	Times     int     `json:"times" xorm:"INT(11) not null"`
	StartTime int     `json:"start_time" xorm:"INT(11) not null"`
	EndTime   int     `json:"end_time" xorm:"INT(11) not null"`
	Remark    string  `json:"remark" xorm:"VARCHAR(256) not null"`
	Pic       string  `json:"pic" xorm:"VARCHAR(512) not null"`
	Status    int     `json:"status" xorm:"INT(5) not null"`
	CreatedAt int64   `json:"created_at" xorm:"INT(11) not null default 0"`
	Member    *Member `json:"member" xorm:"-"`
	Device    *Device `json:"device" xorm:"-"`
}

type MemberDeviceRecord struct {
	Id        int   `json:"id" xorm:"pk autoincr INT(11)"`
	MemberId  int   `json:"member_id" xorm:"INT(11) not null default 0 comment('会员ID')"`
	DeviceId  int   `json:"device_id" xorm:"INT(11) not null default 0 comment('设备ID')"`
	Type      int   `json:"type" xorm:"INT(3) not null default 0 comment('类型 1：计次 2：计时')"`
	Value     int   `json:"value" xorm:"INT(11) not null default 0 comment('type=1 充值次数 type=2 结束时间')"`
	CreatedAt int64 `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberAddReq struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Emergency  string `json:"emergency"`
	Birthday   string `json:"birthday"`
	Gender     int    `json:"gender"`
	UserRemark string `json:"user_remark"`
	//PackageId      int                     `json:"package_id"`
	//Type           int                     `json:"type"`
	//Price          int                     `json:"price"`
	//Cost           int                     `json:"cost"`
	//RechargeRemark string                  `json:"recharge_remark"`
	//Pic            string                  `json:"pic"`
	//RechargeDetail []MemberRecordDetailAdd `json:"recharge_detail"`
}

type MemberRechargeReq struct {
	MemberId       int                     `json:"member_id"`
	PackageId      int                     `json:"package_id"`
	Type           int                     `json:"type"`
	Price          int                     `json:"price"`
	Cost           int                     `json:"cost"`
	Remark         string                  `json:"remark"`
	Pic            string                  `json:"pic"`
	RechargeDetail []MemberRecordDetailAdd `json:"recharge_detail"`
}

type MemberRecordDetailAdd struct {
	DeviceId int `json:"device_id"`
	Type     int `json:"type"`
	Value    int `json:"value"`
}

type MemberRechargeRecordRes struct {
	Name          string `json:"name"`
	TimesRemain   int    `json:"times_remain"`
	MonthlyRemain string `json:"monthly_remain"`
}
