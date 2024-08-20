package ammeter

type AmmeterWarning struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	AmmeterId      int    `json:"ammeter_id" xorm:"INT(11) not null"`
	Type           string `json:"type" xorm:"varchar(128) not null default ''"`
	AElectric      int    `json:"a_electric" xorm:"INT(11) not null default 0"`
	BElectric      int    `json:"b_electric" xorm:"INT(11) not null default 0"`
	CElectric      int    `json:"c_electric" xorm:"INT(11) not null default 0"`
	AVoltage       int    `json:"a_voltage" xorm:"INT(11) not null default 0"`
	BVoltage       int    `json:"b_voltage" xorm:"INT(11) not null default 0"`
	CVoltage       int    `json:"c_voltage" xorm:"INT(11) not null default 0"`
	Power          int    `json:"power" xorm:"INT(11) not null default 0"`
	Leakage        int    `json:"leakage" xorm:"INT(11) not null default 0"`
	Temperature    int    `json:"temperature" xorm:"INT(11) not null default 0"`
	Times          int    `json:"times" xorm:"INT(11) not null default 0"`
	Status         int    `json:"status" xorm:"INT(11) not null default 0"`
	DealTime       int    `json:"-" xorm:"INT(11) not null default 0"`
	DealTimeStr    string `json:"deal_time_str" xorm:"-"`
	DealUser       int    `json:"-" xorm:"INT(11) not null default 0"`
	DealUserName   string `json:"deal_user_name" xorm:"-"`
	CreateTime     int    `json:"-" xorm:"int(11) not null default 0"`
	WarningTimeStr string `json:"warning_time_str" xorm:"-"`
}

type AmmeterWarningUpdateReq struct {
	WarningId int `json:"warning_id"`
	Status    int `json:"status"`
	UserId    int `json:"user_id"`
}
