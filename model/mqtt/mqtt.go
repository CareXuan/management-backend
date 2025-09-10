package mqtt

import "switchboard-backend/model"

type Mqtt struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('连接名称')"`
	Ip             string `json:"ip" xorm:"VARCHAR(32) not null default '' comment('端口')"`
	UserName       string `json:"user_name" xorm:"VARCHAR(32) not null default '' comment('用户名')"`
	Password       string `json:"password" xorm:"VARCHAR(32) not null default '' comment('密码')"`
	ClientId       string `json:"client_id" xorm:"VARCHAR(64) not null default '' comment('client-id')"`
	Theme          string `json:"theme" xorm:"VARCHAR(64) not null default '' comment('主题')"`
	Topic          string `json:"topic" xorm:"VARCHAR(64) not null default '' comment('topic')"`
	model.TimeBase `xorm:"extends"`
}
