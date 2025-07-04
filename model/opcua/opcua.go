package opcua

import "switchboard-backend/model"

type Opcua struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	Ip             string `json:"ip" xorm:"VARCHAR(32) not null default '' comment('IP')"`
	Port           string `json:"port" xorm:"VARCHAR(16) not null default '' comment('端口')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}

type OpcuaData struct {
	Id             int `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId       int `json:"device_id" xorm:"INT(10) not null default 0 comment('设备ID')"`
	ParentId       int `json:"parent_id" xorm:"INT(10) not null default 0 comment('父级节点ID')"`
	Position       int `json:"position" xorm:"INT(10) not null default 0 comment('位数')"`
	model.TimeBase `xorm:"extends"`
}

type AddOpcuaReq struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	IsEnable int    `json:"is_enable"`
}
