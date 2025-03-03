package model

type Config struct {
	Id              int `json:"id" xorm:"pk autoincr INT(11)"`
	GroupId         int `json:"group_id" xorm:"INT(10) not null default 0 comment('礼物组ID')"`
	GiftProbability int `json:"gift_probability" xorm:"INT(3) not null default 0 comment('抽中卡牌概率')"`
}
