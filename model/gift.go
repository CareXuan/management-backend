package model

type Gift struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(24) not null default '' comment('礼物名称')"`
	Pic         string `json:"pic" xorm:"VARCHAR(128) not null default '' comment('配图路径')"`
	Description string `json:"description" xorm:"VARCHAR(256) not null default '' comment('礼物描述')"`
	Level       int    `json:"level" xorm:"INT(3) not null default 0 comment('礼物等级')"`
	Show        int    `json:"show" xorm:"INT(3) not null default 0 comment('是否展示在前端 1:展示 2:不展示')"`
	CanObtain   int    `json:"can_obtain" xorm:"INT(3) not null default 0 comment('是否可被获取 1：是 2：否')"`
	Year        string `json:"year" xorm:"VARCHAR(4) not null default '' comment('年份')"`
	CreateAt    int    `json:"create_at" xorm:"INT(10) not null default 0"`
}

type GiftGroup struct {
	Id        int     `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string  `json:"name" xorm:"VARCHAR(24) not null default '' comment('礼物组名称')"`
	CreateAt  int     `json:"create_at" xorm:"INT(10) not null default 0"`
	GroupGift []*Gift `json:"group_gift" xorm:"-"`
}

type GiftGroupGift struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	GroupId  int `json:"group_id" xorm:"INT(10) not null default 0 comment('礼物组ID')"`
	GiftId   int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	CreateAt int `json:"create_at" xorm:"INT(10) not null default 0"`
}

type GiftUpgrade struct {
	Id            int `json:"id" xorm:"pk autoincr INT(11)"`
	GiftId        int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	UpgradeGiftId int `json:"upgrade_gift_id" xorm:"INT(10) not null default 0 comment('升级所需礼物ID')"`
	Count         int `json:"count" xorm:"INT(10) not null default 0 comment('所需数量')"`
	CreateAt      int `json:"create_at" xorm:"INT(10) not null default 0"`
}

type GiftPackage struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	GiftId   int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	Count    int `json:"count" xorm:"INT(10) not null default 0 comment('礼物数量')"`
	CreateAt int `json:"create_at" xorm:"INT(10) not null default 0"`
}

type GiftExtract struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	GiftId   int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	GetTime  int `json:"get_time" xorm:"INT(10) not null default 0 comment('获取时间')"`
	CreateAt int `json:"create_at" xorm:"INT(10) not null default 0"`
}

type GiftAddReq struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Show        int    `json:"show"`
	CanObtain   int    `json:"can_obtain"`
	Year        string `json:"year"`
}

type GiftGroupAdd struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	GiftIds []int  `json:"gift_ids"`
}
