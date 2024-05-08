package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"management-backend/conf"
	"management-backend/utils"
)

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
	SenceId    string `json:"senceId"`
	SenceOut   string `json:"senceOut"`
	ZoneId     int    `json:"zoneId"`
	ZoneName   string `json:"zoneName"`
}

func GetVehicleConfigFromRedis(deviceId int) (*VehicleConfig, error) {
	var config VehicleConfig
	vehicleConfig, err := redis.String(conf.Redis.Do("GET", fmt.Sprintf(utils.REDIS_KEY_VEHICLE_CONFIG, deviceId)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(vehicleConfig), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
