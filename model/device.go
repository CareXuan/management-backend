package model

type Device struct {
	Id       int    `json:"id" xorm:"pk autoincr INT(11)"`
	Type     int    `json:"type" xorm:"INT(3) not null default 0 comment('端口类型 1：网口 2：光口')"`
	Name     string `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	Port     string `json:"port" xorm:"VARCHAR(24) not null default '' comment('端口')"`
	Network  string `json:"network" xorm:"VARCHAR(24) not null default '' comment('网络IP配置')"`
	TimeBase `xorm:"extends"`
}

type ChangeDeviceReq struct {
	Port    string `json:"port"`
	Network string `json:"network"`
}
