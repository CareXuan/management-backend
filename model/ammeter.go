package model

type Ammeter struct {
	Id              int        `json:"id" xorm:"pk autoincr INT(11)"`
	Type            int        `json:"type" xorm:"tinyint not null"`
	Model           string     `json:"model" xorm:"varchar(32) not null"`
	Num             string     `json:"num" xorm:"varchar(64) not null"`
	Card            string     `json:"card" xorm:"varchar(64) not null"`
	Location        string     `json:"location" xorm:"varchar(64) not null"`
	Status          int        `json:"status" xorm:"int not null"`
	Switch          int        `json:"switch" xorm:"INT(11) not null default 0"`
	ElectricCurrent int        `json:"electric_current" xorm:"INT(11) not null default 0"`
	LeakageCurrent  int        `json:"leakage_current" xorm:"INT(11) not null default 0"`
	Voltage         int        `json:"voltage" xorm:"INT(11) not null default 0"`
	Power           int        `json:"power" xorm:"INT(11) not null default 0"`
	ParentId        int        `json:"parent_id" xorm:"int(11) not null default 0"`
	Children        []*Ammeter `json:"children" xorm:"-"`
	IsSupervisor    int        `json:"is_supervisor" xorm:"-"`
	CreateTime      int        `json:"-" xorm:"int(11) not null default 0"`
}
