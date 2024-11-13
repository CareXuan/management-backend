package model

const CONFIG_TYPE_YEAR = 1

type Config struct {
	Id    int    `json:"id" xorm:"pk autoincr INT(11)"`
	Type  int    `json:"type" xorm:"INT(3) not null default 0 comment('配置类型 1:当前年份')"`
	Value string `json:"value" xorm:"VARCHAR(64) not null default '' comment('配置值')"`
}

type SetConfigReq struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}
