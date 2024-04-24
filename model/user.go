package model

type User struct {
	Id        int    `json:"-" xorm:"pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"varchar(128) not null"`
	Password  string `json:"password" xorm:"varchar(256) not null"`
	Phone     string `json:"phone" xorm:"varchar(20) not null"`
	Token     string `json:"token" xorm:"varchar(128) not null"`
	RoleStr   Role   `json:"role" xorm:"-"`
	CreatedAt int64  `json:"-" xorm:"int(20) not null default 0"`
	UpdatedAt int64  `json:"-" xorm:"int(20) not null default 0"`
	DeletedAt int64  `json:"-" xorm:"int(20) not null default 0"`
}
