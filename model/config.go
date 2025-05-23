package model

type Config struct {
	Id       int    `json:"id" xorm:"pk autoincr INT(11)"`
	OnePoint int    `json:"one_point" xorm:"INT(10) not null default 0 comment('单抽消耗')"`
	TenPoint int    `json:"ten_point" xorm:"INT(10) not null default 0 comment('十抽消耗')"`
	LoginPwd string `json:"login_pwd" xorm:"VARCHAR(24) not null default '' comment('登录密码')"`
}
