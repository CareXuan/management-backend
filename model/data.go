package model

type SbkData struct {
	Id      int    `json:"id" xorm:"pk autoincr INT(11)"`
	Idxx    string `json:"idxx" xorm:"VARCHAR(128) not null default '' comment('ID')"`
	Bxyj    string `json:"bxyj" xorm:"VARCHAR(2) not null default '' comment('BXYJ')"`
	Bmddm   string `json:"bmddm" xorm:"VARCHAR(4) not null default '' comment('报名点代码')"`
	Bmdmc   string `json:"bmdmc" xorm:"VARCHAR(64) not null default '' comment('报名点名称')"`
	Hbbmdmc string `json:"hbbmdmc" xorm:"VARCHAR(64) not null default '' comment('合并报名点名称')"`
	Xm      string `json:"xm" xorm:"VARCHAR(64) not null default '' comment('姓名')"`
	Ksbh    string `json:"ksbh" xorm:"VARCHAR(16) not null default '' comment('考生编号')"`
	Kmdm    string `json:"kmdm" xorm:"VARCHAR(4) not null default '' comment('科目代码')"`
	Kmmc    string `json:"kmmc" xorm:"VARCHAR(64) not null default '' comment('科目名称')"`
	Kmdy    string `json:"kmdy" xorm:"VARCHAR(2) not null default '' comment('科目单元')"`
	Dy2     int64  `json:"dy2" xorm:"INT(2) not null default 0 comment('单元2')"`
	Dy3     int64  `json:"dy3" xorm:"INT(2) not null default 0 comment('单元3')"`
	Dy4     int64  `json:"dy4" xorm:"INT(2) not null default 0 comment('单元4')"`
	Kssj    string `json:"kssj" xorm:"VARCHAR(32) not null default '' comment('考试时间')"`
	Xlsh    int64  `json:"xlsh" xorm:"INT(8) not null default 0 comment('小流水号')"`
	Dlsh    int64  `json:"dlsh" xorm:"INT(8) not null default 0 comment('大流水号')"`
	Djtxm1  string `json:"djtxm1" xorm:"djtxm1 VARCHAR(16) not null default '' comment('djtxm1')"`
	Djtxm2  string `json:"djtxm2" xorm:"djtxm2 VARCHAR(16) not null default '' comment('djtxm2')"`
	Djtxm3  string `json:"djtxm3" xorm:"djtxm3 VARCHAR(16) not null default '' comment('djtxm3')"`
	Djtxm4  string `json:"djtxm4" xorm:"djtxm4 VARCHAR(16) not null default '' comment('djtxm4')"`
	Year    int    `json:"year" xorm:"int(4) not null default 0 comment('年份')"`
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
}
