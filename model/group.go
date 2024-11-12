package model

type Group struct {
	Id    int     `json:"id" xorm:"pk autoincr INT(11)"`
	Name  string  `json:"name" xorm:"VARCHAR(64) not null default '' comment('组名')"`
	Users []*User `json:"users" xorm:"-"`
}

type GroupUser struct {
	Id      int   `json:"id" xorm:"pk autoincr INT(11)"`
	UserId  int   `json:"user_id" xorm:"INT(10) not null default 0 comment('用户ID')"`
	GroupId int   `json:"group_id" xorm:"INT(10) not null default 0 comment('组ID')"`
	User    User  `json:"user" xorm:"-"`
	Group   Group `json:"group" xorm:"-"`
}

type GroupData struct {
	Id      int    `json:"id" xorm:"pk autoincr INT(11)"`
	GroupId int    `json:"group_id" xorm:"INT(10) not null default 0 comment('组ID')"`
	Bmddm   string `json:"bmddm" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	Kmdm    string `json:"kmdm" xorm:"VARCHAR(4) not null default '' comment('科目代码')"`
}

type GroupAddReq struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	UserIds []int  `json:"user_ids"`
}
