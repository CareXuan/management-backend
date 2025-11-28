package modbus

import "switchboard-backend/model"

type Modbus struct {
	Id             int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name           string `json:"name" xorm:"VARCHAR(64) not null default '' comment('点位名称')"`
	ConnectType    int    `json:"connect_type" xorm:"INT(3) not null default 0 comment('连接类型 1：tcp 2：rs485 3:rs232')"`
	Ip             string `json:"ip" xorm:"VARCHAR(24) not null default '' comment('设备IP')"`
	Port           string `json:"port" xorm:"VARCHAR(8) not null default '' comment('端口')"`
	BaudRate       string `json:"baud_rate" xorm:"VARCHAR(8) not null default '' comment('波特率')"`
	RegisterType   int    `json:"register_type" xorm:"INT(3) not null default 0 comment('功能码 1：03 2：04 3：485服务')"`
	StartAddress   int    `json:"start_address" xorm:"INT(3) not null default 0' comment('起始地址')"`
	Count          int    `json:"count" xorm:"INT(8) not null default 0 comment('读取位数')"`
	Slave          int    `json:"slave" xorm:"INT(8) not null default 0 comment('slave')"`
	RangeMin       int    `json:"range_min" xorm:"INT(10) not null default 0 comment('最小量程')"`
	RangeMax       int    `json:"range_max" xorm:"INT(10) not null default 0 comment('最大量程')"`
	Position       int    `json:"position" xorm:"INT(10) not null default 0 comment('获取位置')"`
	Proportion     int    `json:"proportion" xorm:"INT(10) not null default 0 comment('换算比例')"`
	PointDw        string `json:"point_dw" xorm:"VARCHAR(64) not null default '' comment('厂区')"`
	PointQy        string `json:"point_qy" xorm:"VARCHAR(64) not null default '' comment('车间')"`
	PointSb        string `json:"point_sb" xorm:"VARCHAR(64) not null default '' comment('设备')"`
	QsMingcheng    string `json:"qs_mingcheng" xorm:"VARCHAR(64) not null default '' comment('点位地址')"`
	Leixing        string `json:"leixing" xorm:"VARCHAR(10) not null default '' comment('类型 模拟量/开关量')"`
	Miaoshu        string `json:"miaoshu" xorm:"VARCHAR(64) not null default '' comment('点位中文注释')"`
	DataType       string `json:"data_type" xorm:"VARCHAR(64) not null default '' comment('数据类型')"`
	IsEnable       int    `json:"is_enable" xorm:"INT(3) not null default 0 comment('是否可用 1：可用 2：不可用')"`
	model.TimeBase `xorm:"extends"`
}

type AddModbusReq struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ConnectType  int    `json:"connect_type"`
	Ip           string `json:"ip"`
	Port         string `json:"port"`
	BaudRate     string `json:"baud_rate"`
	RegisterType int    `json:"register_type"`
	StartAddress int    `json:"start_address"`
	Count        int    `json:"count"`
	RangeMin     int    `json:"range_min"`
	RangeMax     int    `json:"range_max"`
	Position     int    `json:"position"`
	Proportion   int    `json:"proportion"`
	PointDw      string `json:"point_dw"`
	PointQy      string `json:"point_qy"`
	PointSb      string `json:"point_sb"`
	QsMingcheng  string `json:"qs_mingcheng"`
	Leixing      string `json:"leixing"`
	Miaoshu      string `json:"miaoshu"`
	DataType     string `json:"data_type"`
	Slave        int    `json:"slave"`
	IsEnable     int    `json:"is_enable"`
}
