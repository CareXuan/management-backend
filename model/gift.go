package model

type Gift struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(24) not null default '' comment('礼物名称')"`
	Pic         string `json:"pic" xorm:"VARCHAR(128) not null default '' comment('配图路径')"`
	Description string `json:"description" xorm:"VARCHAR(256) not null default '' comment('礼物描述')"`
	Level       string `json:"level" xorm:"VARCHAR(2) not null default '' comment('礼物等级')"`
	Show        int    `json:"show" xorm:"INT(3) not null default 0 comment('是否展示在前端 1:展示 2:不展示')"`
	CanObtain   int    `json:"can_obtain" xorm:"INT(3) not null default 0 comment('是否可被获取 1：是 2：否')"`
	CrushCnt    int    `json:"crush_cnt" xorm:"INT(5) not null default 0 comment('粉碎所得抽卡点数')"`
	Consumable  int    `json:"consumable" xorm:"INT(3) not null default 0 comment('是否可消耗 1：可以 2：不可以')"`
	CreateAt    int    `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt    int    `json:"delete_at" xorm:"INT(10) not null default 0"`
}

type GiftGroup struct {
	Id        int              `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string           `json:"name" xorm:"VARCHAR(24) not null default '' comment('礼物组名称')"`
	Status    int              `json:"status" xorm:"INT(3) not null default 0 comment('礼物组状态 1：启用 2：禁用')"`
	StartTime int              `json:"start_time" xorm:"INT(10) not null default 0 comment('开始时间')"`
	EndTime   int              `json:"end_time" xorm:"INT(10) not null default 0 comment('结束时间')"`
	CreateAt  int              `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt  int              `json:"delete_at" xorm:"INT(10) not null default 0"`
	GroupGift []*GiftGroupGift `json:"group_gift" xorm:"-"`
}

type GiftGroupGift struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	GroupId     int    `json:"group_id" xorm:"INT(10) not null default 0 comment('礼物组ID')"`
	Level       string `json:"level" xorm:"VARCHAR(5) not null default 0 comment('礼物等级')"`
	Probability int    `json:"probability" xorm:"INT(5) not null default 0 comment('礼物在礼物组中被获取的概率')"`
	CreateAt    int    `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt    int    `json:"delete_at" xorm:"INT(10) not null default 0"`
}

type GiftPackage struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	GiftId     int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	Count      int `json:"count" xorm:"INT(10) not null default 0 comment('礼物数量')"`
	Consumable int `json:"consumable" xorm:"INT(3) not null default 0 comment('是否可消耗 1：可以 2：不可以')"`
	CreateAt   int `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt   int `json:"delete_at" xorm:"INT(10) not null default 0"`
}

type GiftExtract struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	GiftId   int `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	Count    int `json:"count" xorm:"INT(5) not null default 0 comment('获取数量')"`
	Type     int `json:"type" xorm:"INT(3) not null default 0 comment('获取类型 1：获取 2：销毁')"`
	GetTime  int `json:"get_time" xorm:"INT(10) not null default 0 comment('获取时间')"`
	CreateAt int `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt int `json:"delete_at" xorm:"INT(10) not null default 0"`
}

type GiftAddReq struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Consumable  int    `json:"consumable"`
	Show        int    `json:"show"`
	CanObtain   int    `json:"can_obtain"`
}

type GiftAddPointReq struct {
	Count int `json:"count"`
}

type GiftConfigSetReq struct {
	OneCount int `json:"one_count"`
	TenCount int `json:"ten_count"`
}

type GiftDeleteReq struct {
	Ids []int `json:"ids"`
}

type GiftGroupAdd struct {
	Id        int                    `json:"id"`
	Name      string                 `json:"name"`
	StartTime int                    `json:"start_time"`
	EndTime   int                    `json:"end_time"`
	Status    int                    `json:"status"`
	GiftIds   []giftLevelProbability `json:"gift_ids"`
}

type giftLevelProbability struct {
	Level       string `json:"level"`
	Probability int    `json:"probability"`
}

type GiftGroupDelete struct {
	Ids []int `json:"ids"`
}

type GiftChangeStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type GiftProbabilityRes struct {
	GetCount int `json:"get_count"`
	AllCount int `json:"all_count"`
}

type GiftRemainListItem struct {
	GiftId int    `json:"gift_id"`
	Name   string `json:"name"`
	Pic    string `json:"pic"`
	Exist  int    `json:"exist"`
}

type GiftRemainRes struct {
	List        []*GiftRemainListItem          `json:"list"`
	Probability map[string]*GiftProbabilityRes `json:"probability"`
}

var GIFT_LEVEL_MAPPING = map[string]int{
	"A": 700,
	"B": 70,
	"C": 27,
	"D": 7,
	"E": 0,
}
