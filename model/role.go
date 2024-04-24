package model

type Role struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"varchar(64) not null"`
	CreatedAt int    `json:"-" xorm:"int(10) not null default 0"`
	UpdatedAt int    `json:"-" xorm:"int(10) not null default 0"`
	DeletedAt int    `json:"-" xorm:"int(10) not null default 0"`
}
