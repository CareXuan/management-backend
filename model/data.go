package model

const CHECK_STEP_WAITING = 1
const CHECK_STEP_FIRST = 2
const CHECK_STEP_SECOND = 3
const CHECK_STEP_THIRD = 4
const CHECK_STEP_FOURTH = 5

const CHECK_STATUS_WAITING = 1
const CHECK_STATUS_PASS = 2
const CHECK_STATUS_ERROR = 3
const CHECK_STATUS_BLUE_ERROR = 4
const CHECK_STATUS_RED_ERROR = 5
const CHECK_STATUS_BLACK_ERROR = 6
const CHECK_STATUS_BLACK_RED_ERROR = 7
const CHECK_STATUS_BLACK_BLUE_ERROR = 8
const CHECK_STATUS_RED_BLUE_ERROR = 9
const CHECK_STATUS_ALL_ERROR = 10

const STEP_STATUS_WAITING = 1
const STEP_STATUS_PASS = 2

var CHECK_STEP_NAME_MAPPING = map[int]string{
	CHECK_STEP_WAITING: "封装流程开始",
	CHECK_STEP_FIRST:   "总数清点",
	CHECK_STEP_SECOND:  "单元清点",
	CHECK_STEP_THIRD:   "首末号核对",
	CHECK_STEP_FOURTH:  "封装确认",
}

type CheckData struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	Bmddm       string `json:"bmddm" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	Bmdmc       string `json:"bmdmc" xorm:"VARCHAR(64) not null default '' comment('报名点名称')"`
	EnvelopeCnt int    `json:"envelope_cnt" xorm:"INT(8) not null default 0 comment('信封数')"`
	FirstCnt    int    `json:"first_cnt" xorm:"INT(8) not null default 0 comment('信封首号')"`
	LastCnt     int    `json:"last_cnt" xorm:"INT(8) not null default 0 comment('信封末号')"`
	BlueCnt     int    `json:"blue_cnt" xorm:"INT(8) not null default 0 comment('蓝封数量')"`
	RedCnt      int    `json:"red_cnt" xorm:"INT(8) not null default 0 comment('红封数量')"`
	BlackCnt    int    `json:"black_cnt" xorm:"INT(8) not null default 0 comment('黑封数量')"`
	Step        int    `json:"step" xorm:"INT(4) not null default 0 comment('步骤号')"`
	Year        string `json:"year" xorm:"VARCHAR(4) not null default 0 comment('年份')"`
	ProgressCnt int    `json:"progress_cnt" xorm:"-"`
	RemainCnt   int    `json:"remain_cnt" xorm:"-"`
}

type StepCheckData struct {
	Id     int    `json:"id" xorm:"pk autoincr INT(11)"`
	Bmddm  string `json:"bmddm" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	Step   int    `json:"step" xorm:"INT(4) not null default 0 comment('步骤号')"`
	Status int    `json:"status" xorm:"INT(4) not null default 0 comment('检查状态 1:未检查 2:检查通过 3:检查不通过')"`
	Year   string `json:"year" xorm:"VARCHAR(4) not null default 0 comment('年份')"`
}

type NextStepReq struct {
	Bmddm string `json:"bmddm"`
	Step  int    `json:"step"`
}

type CheckReq struct {
	Bmddm string `json:"bmddm"`
	Step  int    `json:"step"`
	Data  []int  `json:"data"`
}

type CheckListRes struct {
	Bmddm string `json:"bmddm"`
	Bmdmc string `json:"bmdmc"`
	Step1 int    `json:"step1"`
	Step2 int    `json:"step2"`
	Step3 int    `json:"step3"`
	Step4 int    `json:"step4"`
	Step5 int    `json:"step5"`
}
