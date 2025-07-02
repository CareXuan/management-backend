package port

import "switchboard-backend/model"

type Network struct {
	Id             int     `json:"id" xorm:"pk autoincr INT(11)"`
	Type           int     `json:"type" xorm:"INT(3) not null default 0 comment('端口类型 1：网口 2：光口')"`
	NetworkType    int     `json:"network_type" xorm:"INT(3) not null default 0 comment('网络配置类型 1：IP 2：网桥')"`
	Name           string  `json:"name" xorm:"VARCHAR(64) not null default '' comment('设备名称')"`
	Port           string  `json:"port" xorm:"VARCHAR(24) not null default '' comment('端口')"`
	Network        string  `json:"network" xorm:"VARCHAR(24) not null default '' comment('网络IP配置')"`
	BridgeId       int     `json:"bridge_id" xorm:"INT(10) not null default 0 comment('bridge_id')"`
	Bridge         *Bridge `json:"bridge" xorm:"-"`
	model.TimeBase `xorm:"extends"`
}

type ChangeNetworkReq struct {
	NetworkType int    `json:"network_type"`
	Port        string `json:"port"`
	Network     string `json:"network"`
	BridgeId    int    `json:"bridge_id"`
}
