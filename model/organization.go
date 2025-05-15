package model

type Organization struct {
	Id       int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name     string `json:"name" xorm:"VARCHAR(64) NOT NULL DEFAULT '' COMMENT('组织名称')"`
	Province string `json:"province" xorm:"VARCHAR(20) not null default '' comment('省份')"`
	City     string `json:"city" xorm:"VARCHAR(20) not null default '' comment('市')"`
	Zone     string `json:"zone" xorm:"VARCHAR(20) not null default '' comment('区县')"`
	UserIds  []int  `json:"user_ids" xorm:"-"`
	Ts       string `json:"ts" xorm:"timestamp not null"`
}

type OrganizationUser struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	OrganizationId int    `json:"organization_id" xorm:"INT(11) not null default 0 comment('组织ID')"`
	UserId         int    `json:"user_id" xorm:"INT(11) not null default 0 comment('用户ID')"`
	Ts             string `json:"ts" xorm:"timestamp not null"`
}

type OrganizationAddReq struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Province string `json:"province"`
	City     string `json:"city"`
	Zone     string `json:"zone"`
	UserIds  []int  `json:"user_ids"`
}

type OrganizationDeleteReq struct {
	Id int `json:"id"`
}
