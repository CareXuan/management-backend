package ammeter

type AmmeterManage struct {
	Id         int `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId  int `json:"ammeter_id" xorm:"INT(11) not null default 0 comment('设备ID')"`
	UserId     int `json:"user_id" xorm:"INT(11) not null default 0 comment('管理员用户ID')"`
	CreateTime int `json:"-" xorm:"int(11) not null default 0"`
}

type TreeManagerRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
