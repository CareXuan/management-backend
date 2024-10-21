package ammeter

type AmmeterManage struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId  int `json:"ammeter_id" xorm:"INT(11) not null default 0"`
	UserId     int `json:"user_id" xorm:"INT(11) not null default 0"`
	CreateTime int `json:"-" xorm:"int(11) not null default 0"`
}

type AmmeterManageConfig struct {
	Id         int    `json:"id" xorm:"pk autoincr INT(11)"`
	UserId     int    `json:"user_id" xorm:"INT(11) not null default 0 comment('用户ID')"`
	AmmeterId  int    `json:"ammeter_id" xorm:"INT(11) not null default 0 comment('设备ID')"`
	Type       int    `json:"type" xorm:"INT(3) not null default 0 comment('配置类型 1：密码')"`
	Value      string `json:"value" xorm:"VARCHAR(256) not null default '' comment('配置值')"`
	CreateTime int    `json:"create_time" xorm:"INT(11) not null default 0"`
}

type AmmeterManageLog struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	UserId     int `json:"user_id" xorm:"INT(11) not null default 0 comment('用户ID')"`
	AmmeterId  int `json:"ammeter_id" xorm:"INT(11) not null default 0 comment('设备ID')"`
	Status     int `json:"status" xorm:"INT(11) not null default 0 comment('状态 1：开 2：停')"`
	CreateTime int `json:"create_time" xorm:"INT(11) not null default 0"`
}

type TreeManagerRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
