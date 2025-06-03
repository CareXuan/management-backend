package model

type TimeBase struct {
	CreatedAt int `json:"created_at" xorm:"int(11) not null default 0 created"`
	UpdatedAt int `json:"updated_at" xorm:"int(11) not null default 0 updated"`
	DeletedAt int `json:"deleted_at" xorm:"int(11) not null default 0 index"`
}
