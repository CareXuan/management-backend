package port

import "switchboard-backend/model"

type Bridge struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('网桥名称')"`
	EnglishName    string `json:"english_name" xorm:"VARCHAR(64) not null default '' comment('英文名称')"`
	Ip             string `json:"ip" xorm:"VARCHAR(16) not null default '' comment('IP')"`
	FileName       string `json:"file_name" xorm:"VARCHAR(64) not null default '' comment('文件路径')"`
	model.TimeBase `xorm:"extends"`
}

type AddBridgeReq struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Ip          string `json:"ip"`
}
