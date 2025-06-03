package model

type TimeBase struct {
	CreatedAt int `json:"-" xorm:"int(11) not null default 0 created"`
	UpdatedAt int `json:"-" xorm:"int(11) not null default 0 updated"`
	DeletedAt int `json:"-" xorm:"int(11) not null default 0 index"`
}
