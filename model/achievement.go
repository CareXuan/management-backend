package model

type Achievement struct {
	Id          int                `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string             `json:"name" xorm:"VARCHAR(24) not null default '' comment('成就名称')"`
	Pic         string             `json:"pic" xorm:"VARCHAR(128) not null default '' comment('配图路径')"`
	Description string             `json:"description" xorm:"VARCHAR(256) not null default '' comment('成就描述')"`
	Point       int                `json:"point" xorm:"INT(10) not null default 0 comment('成就点')"`
	FinishAt    int                `json:"finish_at" xorm:"INT(10) not null default 0 comment('完成时间')"`
	CreateAt    int                `json:"create_at" xorm:"INT(10) not null default 0"`
	Tasks       []*AchievementTask `json:"tasks" xorm:"-"`
	Gifts       []*AchievementGift `json:"gifts" xorm:"-"`
}

type AchievementTask struct {
	Id            int   `json:"id" xorm:"pk autoincr INT(11)"`
	AchievementId int   `json:"achievement_id" xorm:"INT(10) not null default 0 comment('成就ID')"`
	TaskId        int   `json:"task_id" xorm:"INT(10) not null default 0 comment('任务ID')"`
	Count         int   `json:"count" xorm:"INT(10) not null default 0 comment('完成成就所需任务完成次数')"`
	CreateAt      int   `json:"create_at" xorm:"INT(10) not null default 0"`
	TaskItem      *Task `json:"task_item" xorm:"-"`
}

type AchievementGift struct {
	Id            int   `json:"id" xorm:"pk autoincr INT(11)"`
	AchievementId int   `json:"achievement_id" xorm:"INT(10) not null default 0 comment('成就ID')"`
	GiftId        int   `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	Count         int   `json:"count" xorm:"INT(10) not null default 0 comment('完成成就奖励礼物数量')"`
	CreateAt      int   `json:"create_at" xorm:"INT(10) not null default 0"`
	GiftItem      *Gift `json:"gift_item" xorm:"-"`
}

type AchievementAddReq struct {
	Id          int                  `json:"id"`
	Name        string               `json:"name"`
	Pic         string               `json:"pic"`
	Description string               `json:"description"`
	Point       int                  `json:"point"`
	Tasks       []achievementTaskAdd `json:"tasks"`
	Gifts       []achievementGiftAdd `json:"gifts"`
}

type achievementTaskAdd struct {
	TaskId int `json:"task_id"`
	Count  int `json:"count"`
}

type achievementGiftAdd struct {
	GiftId int `json:"gift_id"`
	Count  int `json:"count"`
}
