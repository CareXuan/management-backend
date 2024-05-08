package model

type SignalData struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId    int    `json:"device_id" xorm:"INT(11) not null"`
	SignalId    int    `json:"signal_id" xorm:"INT(11) not null"`
	StationId   int    `json:"station_id" xorm:"INT(11) not null"`
	CoilNum     int    `json:"coil_num" xorm:"INT(11) not null"`
	SameCount   int    `json:"same_count" xorm:"INT(11) not null"`
	CreatedTime string `json:"-" xorm:"datetime not null default 0"`
	UpdatedTime string `json:"-" xorm:"datetime not null default 0"`
}
