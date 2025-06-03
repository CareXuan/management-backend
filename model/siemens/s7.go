package model

type SiemensS7 struct {
	Id   int    `json:"id" xorm:"pk autoincr INT(11)"`
	Ip   string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Rack int    `json:"rack" xorm:"INT(3) not null default 0 comment('rack')"`
	Slot int    `json:"slot" xorm:"INT(3) not null default 0 comment('slot')"`
}
