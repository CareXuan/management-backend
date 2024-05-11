package model

type AmmeterWarning struct {
	Id           int    `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId    int    `json:"ammeterId" xorm:"INT(11) not null"`
	Type         int    `json:"type" xorm:"INT(11) not null default 0"`
	Status       int    `json:"status" xorm:"INT(11) not null default 0"`
	DealTime     int    `json:"-" xorm:"INT(11) not null default 0"`
	DealTimeStr  string `json:"deal_time_str" xorm:"-"`
	DealUser     int    `json:"-" xorm:"INT(11) not null default 0"`
	DealUserName string `json:"deal_user_name" xorm:"-"`
	CreateTime   int    `json:"-" xorm:"int(11) not null default 0"`
}
