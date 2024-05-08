package model

type Device struct {
	Id              int    `json:"id" xorm:"pk autoincr INT(11)"`
	Type            int    `json:"type" xorm:"tinyint not null"`
	Iccid           string `json:"iccid" xorm:"varchar(31) not null"`
	Name            string `json:"name" xorm:"varchar(255) not null"`
	SignalIntensity int    `json:"signal_intensity" xorm:"int not null"`
	Status          int    `json:"status" xorm:"int not null"`
	Comment         string `json:"comment" xorm:"varchar(255) not null"`
	Ts              string `json:"-" xorm:"int(11) not null default 0"`
}
