package ammeter

type Ammeter struct {
	Id              int        `json:"id" xorm:"pk autoincr INT(11)"`
	Type            int        `json:"type" xorm:"tinyint not null comment('设备类型 1：电表')"`
	Model           string     `json:"model" xorm:"varchar(32) not null comment('设备型号')"`
	Num             string     `json:"num" xorm:"varchar(64) not null comment('设备编号')"`
	Card            string     `json:"card" xorm:"varchar(64) not null comment('设备卡号')"`
	Location        string     `json:"location" xorm:"varchar(64) not null comment('设备位置')"`
	Status          int        `json:"status" xorm:"int not null comment('设备状态 1：在线 2：离线')"`
	Switch          int        `json:"switch" xorm:"INT(11) not null default 0 comment('设备开关 1：开 2：关')"`
	ElectricCurrent int        `json:"electric_current" xorm:"INT(11) not null default 0 comment('实时电流')"`
	LeakageCurrent  int        `json:"leakage_current" xorm:"INT(11) not null default 0 comment('实时漏电电流')"`
	Voltage         int        `json:"voltage" xorm:"INT(11) not null default 0 comment('电压')"`
	Power           int        `json:"power" xorm:"INT(11) not null default 0 comment('功率')"`
	ParentId        int        `json:"parent_id" xorm:"int(11) not null default 0 comment('父级节点ID')"`
	Children        []*Ammeter `json:"children" xorm:"-"`
	IsSupervisor    int        `json:"is_supervisor" xorm:"-"`
	CreateTime      int        `json:"-" xorm:"int(11) not null default 0"`
}

type AmmeterNodeAdd struct {
	NodeId    int    `json:"node_id"`
	NodeType  int    `json:"node_type"`
	NodeModel string `json:"node_model"`
	Num       string `json:"num"`
	Card      string `json:"card"`
	Location  string `json:"location"`
	ParentId  int    `json:"parent_id"`
	Managers  []int  `json:"managers"`
}

type ChangeSwitchReq struct {
	AmmeterId int `json:"ammeter_id"`
	Switch    int `json:"switch"`
}
