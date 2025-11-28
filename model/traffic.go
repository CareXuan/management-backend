package model

type Traffic struct {
	Id          int    `json:"id" xorm:"pk autoincr INT(11)"`
	Name        string `json:"name" xorm:"VARCHAR(64) not null default ''"`
	BytesSent   uint64 `json:"bytes_sent" xorm:"BIGINT(40) not null default 0"`
	BytesRecv   uint64 `json:"bytes_recv" xorm:"BIGINT(40) not null default 0"`
	PacketsSent uint64 `json:"packets_sent" xorm:"BIGINT(40) not null default 0"`
	PacketsRecv uint64 `json:"packets_recv" xorm:"BIGINT(40) not null default 0"`
	Errin       uint64 `json:"errin" xorm:"BIGINT(40) not null default 0"`
	Errout      uint64 `json:"errout" xorm:"BIGINT(40) not null default 0"`
	Dropin      uint64 `json:"dropin" xorm:"BIGINT(40) not null default 0"`
	Dropout     uint64 `json:"dropout" xorm:"BIGINT(40) not null default 0"`
	Fifoin      uint64 `json:"fifoin" xorm:"BIGINT(40) not null default 0"`
	Fifoout     uint64 `json:"fifoout" xorm:"BIGINT(40) not null default 0"`
	CreateTime  int    `json:"create_time" xorm:"INT(10) not null default 0"`
}
