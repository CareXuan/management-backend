package modbus

import "switchboard-backend/model"

type Modbus struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	Ip             string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Port           int    `json:"port" xorm:"VARCHAR(8) not null default 0 comment('端口')"`
	BaudRate       int    `json:"baud_rate" xorm:"VARCHAR(8) not null default 0 comment('波特率')"`
	StartAddress   int    `json:"start_address" xorm:"VARCHAR(16) not null default 0 comment('起始地址')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}
