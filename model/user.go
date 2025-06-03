package model

type User struct {
	Id          int      `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string   `json:"name" xorm:"varchar(128) not null"`
	Password    string   `json:"password" xorm:"varchar(256) not null"`
	Phone       string   `json:"phone" xorm:"varchar(20) not null"`
	Token       string   `json:"token" xorm:"varchar(128) not null"`
	OpenId      string   `json:"open_id" xorm:"varchar(128) not null default ''"`
	RoleStr     Role     `json:"role" xorm:"-"`
	Permissions []string `json:"permissions" xorm:"-"`
	TimeBase
}

type AddUserReq struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	RoleId   int    `json:"role_id"`
}
