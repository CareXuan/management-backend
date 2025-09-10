package dlt645

import "switchboard-backend/model"

type Dlt645 struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	ConnectType    int    `json:"connect_type" xorm:"INT(3) not null default 0 comment('连接类型 1：tcp 2：rs485 3:rs232')"`
	Ip             string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Port           string `json:"port" xorm:"VARCHAR(8) not null default '' comment('端口')"`
	BaudRate       string `json:"baud_rate" xorm:"VARCHAR(8) not null default '' comment('波特率')"`
	Address        string `json:"address" xorm:"VARCHAR(16) not null default '' comment('电表地址')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}

type AddDlt645Req struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ConnectType int    `json:"connect_type"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
	BaudRate    string `json:"baud_rate"`
	Address     string `json:"address"`
	IsEnable    int    `json:"is_enable"`
}
