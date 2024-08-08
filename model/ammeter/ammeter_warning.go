package ammeter

type AmmeterWarning struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId      int    `json:"ammeterId" xorm:"INT(11) not null comment('设备ID')"`
	Type           int    `json:"type" xorm:"INT(11) not null default 0 comment('报警类型 1：在线 2：离线')"`
	Status         int    `json:"status" xorm:"INT(11) not null default 0 comment('处理状态 1：已处理 2：未处理')"`
	DealTime       int    `json:"-" xorm:"INT(11) not null default 0 comment('处理时间')"`
	DealTimeStr    string `json:"deal_time_str" xorm:"-"`
	DealUser       int    `json:"-" xorm:"INT(11) not null default 0 comment('处理用户ID')"`
	DealUserName   string `json:"deal_user_name" xorm:"-"`
	CreateTime     int    `json:"-" xorm:"int(11) not null default 0"`
	WarningTimeStr string `json:"warning_time_str" xorm:"-"`
}

type AmmeterWarningUpdateReq struct {
	WarningId int `json:"warning_id"`
	Status    int `json:"status"`
	UserId    int `json:"user_id"`
}
