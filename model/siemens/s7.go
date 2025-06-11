package siemens

import "switchboard-backend/model"

type SiemensS7 struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	Ip             string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Rack           int    `json:"rack" xorm:"INT(3) not null default 0 comment('rack')"`
	Slot           int    `json:"slot" xorm:"INT(3) not null default 0 comment('slot')"`
	DbNum          int    `json:"db_num" xorm:"INT(3) not null default 0 comment('DB块编号')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}

type SiemensS7Data struct {
	Id             int `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId       int `json:"device_id" xorm:"INT(10) not null default 0 comment('设备ID')"`
	Type           int `json:"type" xorm:"INT(3) not null default 0 comment('数据类型 1：整型 2：浮点型 3：布尔型')"`
	Start          int `json:"start" xorm:"INT(8) not null default 0 comment('开始地址')"`
	model.TimeBase `xorm:"extends"`
}
