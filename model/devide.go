package model

type Device struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string `json:"name" xorm:"varchar(64) not null"`
	Pic       string `json:"pic" xorm:"varchar(512) not null"`
	Cert      string `json:"cert" xorm:"varchar(512) not null"`
	CreatedAt int    `json:"created_at" xorm:"INT(11) not null default 0"`
}
