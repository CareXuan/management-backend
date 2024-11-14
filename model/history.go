package model

type History struct {
	Id            int    `json:"id" xorm:"pk autoincr INT(11)"`
	Step          int    `json:"step" xorm:"INT(10) not null default 0 comment('步骤号')"`
	StepName      string `json:"step_name" xorm:"-"`
	Status        int    `json:"status" xorm:"INT(10) not null default 0 comment('状态 1:未进行 2:进行中 3:校验成功 4:校验失败')"`
	Remark        string `json:"remark" xorm:"VARCHAR(256) not null default '' comment('操作记录文案')"`
	Bmddm         string `json:"bmddm" xorm:"VARCHAR(8) not null default '' comment('报名点代码')"`
	Bmdmc         string `json:"bmdmc" xorm:"-"`
	UserId        int    `json:"user_id" xorm:"INT(10) not null default 0 comment('用户ID')"`
	UserName      string `json:"user_name" xorm:"VARCHAR(64) not null default 0 comment('用户昵称')"`
	CreateTime    int    `json:"create_time" xorm:"INT(10) not null default 0 comment('时间')"`
	CreateTimeStr string `json:"create_time_str" xorm:"-"`
	Year          string `json:"year" xorm:"VARCHAR(4) not null default 0 comment('年份')"`
}
