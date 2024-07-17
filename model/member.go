package model

const MEMBER_TYPE_TIMES = 1
const MEMBER_TYPE_MONTHLY = 2

const USE_STATUS_RESERVATION = 1
const USE_STATUS_FINISH = 2

type Member struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"VARCHAR(64) not null"`
	Card      string `json:"card" xorm:"VARCHAR(128) not null"`
	Phone     string `json:"phone" xorm:"VARCHAR(16) not null"`
	CreatedAt int64  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberRecord struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	CardId    int    `json:"card_id" xorm:"INT(11) not null"`
	Type      int    `json:"type" xorm:"INT(11) not null"`
	Price     int    `json:"price" xorm:"INT(11) not null"`
	Cost      int    `json:"cost" xorm:"INT(11) not null"`
	Remark    string `json:"remark" xorm:"VARCHAR(256) not null"`
	Pic       string `json:"pic" xorm:"VARCHAR(512) not null"`
	CreatedAt int64  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberUseRecord struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	CardId    int    `json:"card_id" xorm:"INT(11) not null"`
	UseType   int    `json:"use_type" xorm:"INT(5) not null"`
	Times     int    `json:"times" xorm:"INT(11) not null"`
	StartTime int    `json:"start_time" xorm:"INT(11) not null"`
	EndTime   int    `json:"end_time" xorm:"INT(11) not null"`
	Remark    string `json:"remark" xorm:"VARCHAR(256) not null"`
	Pic       string `json:"pic" xorm:"VARCHAR(512) not null"`
	Status    int    `json:"status" xorm:"INT(5) not null"`
	CreatedAt int64  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type MemberAddReq struct {
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Price  int    `json:"price"`
	Cost   int    `json:"cost"`
	Remark string `json:"remark"`
	Pic    string `json:"pic"`
}

type MemberRechargeReq struct {
	CardId int    `json:"card_id"`
	Type   int    `json:"type"`
	Price  int    `json:"price"`
	Cost   int    `json:"cost"`
	Remark string `json:"remark"`
	Pic    string `json:"pic"`
}
