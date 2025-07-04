package firewall

import "switchboard-backend/model"

type Firewall struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Ip             string `json:"ip" xorm:"VARCHAR(64) not null default '' comment('IP')"`
	Type           int    `json:"type" xorm:"INT(3) not null default 0 comment('方向 1：入方向 2：出方向')"`
	AllowType      int    `json:"allow_type" xorm:"INT(3) not null default 0 comment('1：允许 2：禁止')"`
	model.TimeBase `xorm:"extends"`
}

type AddFirewallReq struct {
	Ip        string `json:"ip"`
	Type      int    `json:"type"`
	AllowType int    `json:"allow_type"`
}

type DeleteFirewallReq struct {
	Id int `json:"id"`
}
