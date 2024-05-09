package model

type VehicleConfig struct {
	CampusId   int    `json:"campusId"`
	CampusName string `json:"campusName"`
	CreateTime string `json:"createTime"`
	DeviceId   string `json:"deviceId"`
	DeviceSN   string `json:"deviceSN"`
	DeviceSite string `json:"deviceSite"`
	DeviceType int    `json:"deviceType"`
	Id         int    `json:"id"`
	OperGid    int    `json:"operGid"`
	SenceIn    string `json:"senceIn"`
	SenceOut   string `json:"senceOut"`
	ZoneId     int    `json:"zoneId"`
	ZoneName   string `json:"zoneName"`
}
