package model

type AmmeterManage struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId  int `json:"ammeter_id" xorm:"INT(11) not null default 0"`
	UserId     int `json:"user_id" xorm:"INT(11) not null default 0"`
	CreateTime int `json:"-" xorm:"int(11) not null default 0"`
}

type TreeManagerRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
