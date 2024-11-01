package model

const PACKAGE_STATUS_OPEN = 1
const PACKAGE_STATUS_CLOSE = 0

const DEVICE_STATUS_OPEN = 1
const DEVICE_STATUS_CLOSE = 2

type Device struct {
	Id           int                `json:"id" xorm:"pk autoincr INT(11)"`
	Name         string             `json:"name" xorm:"varchar(64) not null default '' comment('名称')"`
	Pic          string             `json:"pic" xorm:"varchar(512) not null default '' comment('图片')"`
	Cert         string             `json:"cert" xorm:"varchar(512) not null default '' comment('证书')"`
	Status       int                `json:"status" xorm:"INT(3) not null default 0 comment('可用状态 1:可用 2:不可用')"`
	UseTime      int                `json:"use_time" xorm:"INT(11) not null default 0 comment('每次使用时间')"`
	CreatedAt    int64              `json:"created_at" xorm:"INT(11) not null default 0"`
	Appointments []*MemberUseRecord `json:"appointments" xorm:"-"`
}

type DevicePackage struct {
	Id        int                    `json:"id" xorm:"pk autoincr INT(11)"`
	Name      string                 `json:"name" xorm:"VARCHAR(64) not null default '' comment('套餐名称')"`
	Cost      int                    `json:"cost" xorm:"INT(10) not null default 0 comment('价格')"`
	Status    int                    `json:"status" xorm:"INT(3) not null default 0 comment('状态 1：开启 0：关闭')"`
	Details   []*DevicePackageDetail `json:"details" xorm:"-"`
	CreatedAt int64                  `json:"created_at" xorm:"INT(11) not null default 0"`
}

type DevicePackageDetail struct {
	Id        int     `json:"id" xorm:"pk autoincr INT(11)"`
	PackageId int     `json:"package_id" xorm:"INT(10) not null default 0 comment('套餐ID')"`
	DeviceId  int     `json:"device_id" xorm:"INT(11) not null default 0 comment('设备ID')"`
	Device    *Device `json:"device" xorm:"-"`
	Type      int     `json:"type" xorm:"INT(3) not null default 0 comment('类型 1：计次 2：计时')"`
	Value     int     `json:"value" xorm:"INT(11) not null default 0 comment('type=1 充值次数 type=2 结束时间')"`
	CreatedAt int64   `json:"created_at" xorm:"INT(11) not null default 0"`
}

type DeviceChangeStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type DevicePackageAddReq struct {
	Id      int                         `json:"id"`
	Name    string                      `json:"name"`
	Cost    int                         `json:"cost"`
	Details []devicePackageDetailAddReq `json:"details"`
}

type DevicePackageChangeStatusReq struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

type devicePackageDetailAddReq struct {
	DeviceId int `json:"device_id"`
	Type     int `json:"type"`
	Value    int `json:"value"`
}
