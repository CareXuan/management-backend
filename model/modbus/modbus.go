package modbus

import "switchboard-backend/model"

type Modbus struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	ConnectType    int    `json:"connect_type" xorm:"INT(3) not null default 0 comment('连接类型 1：tcp 2：rs485 3:rs232')"`
	Ip             string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Port           string `json:"port" xorm:"VARCHAR(8) not null default '' comment('端口')"`
	BaudRate       string `json:"baud_rate" xorm:"VARCHAR(8) not null default '' comment('波特率')"`
	StartAddress   string `json:"start_address" xorm:"VARCHAR(16) not null default '' comment('起始地址')"`
	Count          int    `json:"count" xorm:"INT(8) not null default 0 comment('读取位数')"`
	Slave          int    `json:"slave" xorm:"INT(8) not null default 0 comment('slave')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}

type AddModbusReq struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ConnectType  int    `json:"connect_type"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	BaudRate     string `json:"baud_rate"`
	StartAddress string `json:"start_address"`
	Count        int    `json:"count"`
	Slave        int    `json:"slave"`
	IsEnable     int    `json:"is_enable"`
}
