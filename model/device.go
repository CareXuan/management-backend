package model

type Device struct {
	Id              int    `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId        int    `json:"device_id" xorm:"INT(10) not null default 0"`
	Iccid           string `json:"iccid" xorm:"VARCHAR(64) not null default ''"`
	Latitude        string `json:"latitude" xorm:"VARCHAR(12) not null default '0'"`
	Longitude       string `json:"longitude" xorm:"VARCHAR(12) not null default '0'"`
	SignalIntensity int    `json:"signal_intensity" xorm:"INT(10) not null default 0"`
	Name            string `json:"name" xorm:"VARCHAR(30) not null default '' comment('商家名称')"`
	Province        string `json:"province" xorm:"VARCHAR(20) not null default '' comment('省份')"`
	City            string `json:"city" xorm:"VARCHAR(20) not null default '' comment('市')"`
	Zone            string `json:"zone" xorm:"VARCHAR(20) not null default '' comment('区县')"`
	Address         string `json:"address" xorm:"VARCHAR(100) not null default '' comment('详细地址')"`
	Manager         string `json:"manager" xorm:"VARCHAR(20) not null default '' comment('管理人员')"`
	Phone           string `json:"phone" xorm:"VARCHAR(15) not null default '' comment('联系电话')"`
	Remark          string `json:"remark" xorm:"VARCHAR(200) not null default '' comment('备注')"`
	Status          int    `json:"status" xorm:"INT(10) not null default 0"`
	Pwd1            int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd1Base        int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd2            int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd2Base        int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd3            int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd3Base        int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd4            int    `json:"-" xorm:"INT(10) not null default 0"`
	Pwd4Base        int    `json:"-" xorm:"INT(10) not null default 0"`
	Ts              string `json:"ts" xorm:"timestamp not null"`
}

type DeviceCommonData struct {
	Id       int    `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId int    `json:"device_id" xorm:"INT(10) not null default 0"`
	IoStatus int    `json:"io_status" xorm:"INT(10) not null default 0"`
	Tem1h    int    `json:"tem1h" xorm:"INT(10) not null default 0"`
	Tem1l    int    `json:"tem1l" xorm:"INT(10) not null default 0"`
	Tem2h    int    `json:"tem2h" xorm:"INT(10) not null default 0"`
	Tem2l    int    `json:"tem2l" xorm:"INT(10) not null default 0"`
	Ts       string `json:"ts" xorm:"timestamp not null"`
}

type DeviceLocationHistory struct {
	Id            int    `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId      int    `json:"device_id" xorm:"INT(10) not null default 0"`
	LocationValid int    `json:"location_valid" xorm:"INT(10) not null default 0"`
	Latitude      string `json:"latitude" xorm:"VARCHAR(12) not null default '0'"`
	Longitude     string `json:"longitude" xorm:"VARCHAR(12) not null default '0'"`
	Ts            string `json:"ts" xorm:"timestamp not null"`
}

type DeviceServiceData struct {
	Id            int    `json:"id" xorm:"pk autoincr INT(11)"`
	DeviceId      int    `json:"device_id" xorm:"INT(10) not null default 0"`
	HighVoltageH  int    `json:"high_voltage_h" xorm:"INT(10) not null default 0"`
	HighVoltageL  int    `json:"high_voltage_l" xorm:"INT(10) not null default 0"`
	HighCurrentH  int    `json:"high_current_h" xorm:"INT(10) not null default 0"`
	HighCurrentL  int    `json:"high_current_l" xorm:"INT(10) not null default 0"`
	SwitchCurrent int    `json:"switch_current" xorm:"INT(0) not null default 0"`
	CurrentBak1   int    `json:"current_bak_1" xorm:"current_bak_1 INT(10) not null default 0"`
	CurrentBak2   int    `json:"current_bak_2" xorm:"current_bak_2 INT(10) not null default 0"`
	DataTime      string `json:"data_time" xorm:"VARCHAR(12) not null default ''"`
	Ts            string `json:"ts" xorm:"timestamp not null"`
}

type DeviceAddReq struct {
	Id       int    `json:"id"`
	DeviceId int    `json:"device_id"`
	Name     string `json:"name"`
	Province string `json:"province"`
	City     string `json:"city"`
	Zone     string `json:"zone"`
	Address  string `json:"address"`
	Manager  string `json:"manager"`
	Phone    string `json:"phone"`
	Remark   string `json:"remark"`
}

type DeviceLocationRes struct {
	DeviceId int    `json:"device_id"`
	Iccid    string `json:"iccid"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
}

type DeviceStatisticFourRes struct {
	Hour string `json:"hour"`
	Sum  int    `json:"sum"`
}

type DeviceStatisticRes struct {
	Columns []string        `json:"columns"`
	Datas   [][]interface{} `json:"datas"`
}
