package model

type Config struct {
	Id       int `json:"id" xorm:"pk autoincr INT(11)"`
	OnePoint int `json:"one_point" xorm:"INT(10) not null default 0 comment('单抽消耗')"`
	TenPoint int `json:"ten_point" xorm:"INT(10) not null default 0 comment('十抽消耗')"`
}
