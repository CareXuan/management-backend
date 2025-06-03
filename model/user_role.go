package model

type UserRole struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	UserId   int `json:"user_id" xorm:"INT(11) not null"`
	RoleId   int `json:"role_id" xorm:"INT(11) not null"`
	TimeBase `xorm:"extends"`
}
