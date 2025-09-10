package data_sync

type DataSyncConfig struct {
	Id           int    `json:"id" xorm:"pk autoincr INT(11)"`
	SqlAddress   string `json:"sql_address" xorm:"VARCHAR(24) not null default '' comment('数据库地址')"`
	SqlDatabase  string `json:"sql_database" xorm:"VARCHAR(24) not null default '' comment('默认数据库')"`
	SqlUser      string `json:"sql_user" xorm:"VARCHAR(24) not null default '' comment('数据库用户')"`
	SqlPassword  string `json:"sql_password" xorm:"VARCHAR(24) not null default '' comment('数据库密码')"`
	SqlTableName string `json:"sql_table_name" xorm:"VARCHAR(24) not null default '' comment('数据库表名')"`
	SqlRate      int    `json:"sql_rate" xorm:"INT(10) not null default 0 comment('推送频率')"`
}

type DataSync struct {
	Id         int    `json:"id" xorm:"pk autoincr INT(11)"`
	DataType   int    `json:"data_type" xorm:"INT(3) not null default 0 comment('数据类型 1：s7 2：modbus 3：opcua 4：dlt645')"`
	Value      string `json:"value" xorm:"VARCHAR(128) not null default '' comment('数据值')"`
	PointId    string `json:"point_id" xorm:"VARCHAR(32) not null default '' comment('点位ID')"`
	PointName  string `json:"point_name" xorm:"VARCHAR(32) not null default '' comment('点位名称')"`
	Status     int    `json:"status" xorm:"INT(3) not null default 0 comment('是否推送完成 0：未完成 1：已完成')"`
	CreateTime int    `json:"create_time" xorm:"INT(10) not null default 0"`
	SyncTime   int    `json:"sync_time" xorm:"INT(10) not null default 0"`
}

type DataSyncAddReq struct {
	SqlAddress   string `json:"sql_address"`
	SqlDatabase  string `json:"sql_database"`
	SqlUser      string `json:"sql_user"`
	SqlPassword  string `json:"sql_password"`
	SqlTableName string `json:"sql_table_name"`
	SqlRate      int    `json:"sql_rate"`
}
