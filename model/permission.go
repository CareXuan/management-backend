package model

type Permission struct {
	Id        int           `json:"id" xorm:"pk autoincr INT(11)"`
	Uuid      int           `json:"uuid" xorm:"varchar(128) not null"`
	Label     string        `json:"label" xorm:"varchar(64) not null"`
	Path      string        `json:"path" xorm:"varchar(256) not null"`
	Icon      string        `json:"icon" xorm:"varchar(512) not null"`
	Component string        `json:"component" xorm:"varchar(128) not null"`
	Desc      string        `json:"desc" xorm:"varchar(128) not null"`
	Sort      int           `json:"sort" xorm:"int(11) not null default 100"`
	Children  []*Permission `json:"children" xorm:"-"`
	ParentId  int           `json:"parent_id" xorm:"int(11) not null default 0"`
	CreatedAt int           `json:"-" xorm:"int(11) not null default 0"`
	UpdatedAt int           `json:"-" xorm:"int(11) not null default 0"`
	DeletedAt int           `json:"-" xorm:"int(11) not null default 0"`
}
