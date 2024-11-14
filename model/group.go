package model

type Group struct {
	Id        int      `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string   `json:"name" xorm:"VARCHAR(64) not null default '' comment('组名')"`
	Users     []int    `json:"users" xorm:"-"`
	UserNames []string `json:"user_names" xorm:"-"`
	BmdStart  string   `json:"bmd_start" xorm:"-"`
	BmdEnd    string   `json:"bmd_end" xorm:"-"`
}

type GroupUser struct {
	Id       int    `json:"id" xorm:"pk autoincr INT(11)"`
	UserId   int    `json:"user_id" xorm:"INT(10) not null default 0 comment('用户ID')"`
	GroupId  int    `json:"group_id" xorm:"INT(10) not null default 0 comment('组ID')"`
	BmdStart string `json:"bmd_start" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	BmdEnd   string `json:"bmd_end" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	Year     string `json:"year" xorm:"VARCHAR(4) not null default 0 comment('年份')"`
	User     User   `json:"user" xorm:"-"`
	Group    Group  `json:"group" xorm:"-"`
}

type GroupAddReq struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UserIds  []int  `json:"user_ids"`
	BmdStart string `json:"bmd_start"`
	BmdEnd   string `json:"bmd_end"`
}

type GroupDeleteReq struct {
	Id int `json:"id"`
}
