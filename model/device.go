package model

type Device struct {
	Id           int            `json:"id" xorm:"pk autoincr INT(11)"`
	Name         string         `json:"name" xorm:"varchar(64) not null"`
	Pic          string         `json:"pic" xorm:"varchar(512) not null"`
	Cert         string         `json:"cert" xorm:"varchar(512) not null"`
	UseTime      int            `json:"use_time" xorm:"INT(11) not null"`
	CreatedAt    int64          `json:"created_at" xorm:"INT(11) not null default 0"`
	Appointments []*Appointment `json:"appointments" xorm:"-"`
}
