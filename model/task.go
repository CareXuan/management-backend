package model

type Task struct {
	Id          int         `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string      `json:"name" xorm:"VARCHAR(24) not null default '' comment('任务名称')"`
	Description string      `json:"description" xorm:"VARCHAR(256) not null default '' comment('任务描述')"`
	Type        int         `json:"type" xorm:"INT(3) not null default 0 comment('任务类型 1：常规任务 2：日常任务 3：限时任务')"`
	StartTime   int         `json:"start_time" xorm:"INT(10) not null default 0 comment('任务开始时间')"`
	Deadline    int         `json:"deadline" xorm:"INT(10) not null default 0 comment('任务截止时间 当type为1时无截止时间 type为2时是1-7表示周一到周日 type为3时是截止时间')"`
	Star        int         `json:"star" xorm:"INT(3) not null default 0 comment('任务星级')"`
	Status      int         `json:"status" xorm:"INT(3) not null default 0 comment('是否开启 1：开启 2：关闭')"`
	Year        string      `json:"year" xorm:"VARCHAR(4) not null default '' comment('年份')"`
	CreateAt    int         `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt    int         `json:"delete_at" xorm:"INT(10) not null default 0"`
	Gifts       []*TaskGift `json:"gifts" xorm:"-"`
}

type TaskGift struct {
	Id       int   `json:"id" xorm:"pk autoincr INT(11)"`
	TaskId   int   `json:"task_id" xorm:"INT(10) not null default 0 comment('任务ID')"`
	GiftId   int   `json:"gift_id" xorm:"INT(10) not null default 0 comment('礼物ID')"`
	Count    int   `json:"count" xorm:"INT(10) not null default 0 comment('完成任务奖励礼物数量')"`
	CreateAt int   `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt int   `json:"delete_at" xorm:"INT(10) not null default 0"`
	GiftItem *Gift `json:"gift_item" xorm:"-"`
}

type TaskDo struct {
	Id        int    `json:"id" xorm:"pk autoincr INT(11)"`
	TaskId    int    `json:"task_id" xorm:"INT(10) not null default 0 comment('任务ID')"`
	Pic       string `json:"pic" xorm:"VARCHAR(128) not null default '' comment('完成任务配的图片')"`
	Status    int    `json:"status" xorm:"INT(10) not null default 0 comment('任务状态 1：待提交 2：待审核 3：通过 4：拒绝')"`
	Reason    string `json:"reason" xorm:"VARCHAR(256) not null default '' comment('拒绝理由')"`
	StartTime int    `json:"start_time" xorm:"INT(10) not null default 0 comment('任务开始时间')"`
	Deadline  int    `json:"deadline" xorm:"INT(10) not null default 0 comment('任务截止时间')"`
	FinishAt  int    `json:"finish_at" xorm:"INT(10) not null default 0 comment('完成时间')"`
	CreateAt  int    `json:"create_at" xorm:"INT(10) not null default 0"`
	DeleteAt  int    `json:"delete_at" xorm:"INT(10) not null default 0"`
}

type TaskAddReq struct {
	Id          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        int            `json:"type"`
	StartTime   int            `json:"start_time"`
	Deadline    int            `json:"deadline"`
	Star        int            `json:"star"`
	Year        string         `json:"year"`
	BindGifts   []taskGiftBind `json:"bind_gifts"`
}

type TaskDeleteReq struct {
	Ids []int `json:"ids"`
}

type TaskChangeStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type taskGiftBind struct {
	GiftId int `json:"gift_id"`
	Count  int `json:"count"`
}

type TaskDoReq struct {
	Id  int    `json:"id"`
	Pic string `json:"pic"`
}

type TaskCheckReq struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Reason string `json:"reason"`
}
