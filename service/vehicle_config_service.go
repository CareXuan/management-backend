package service

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"management-backend/conf"
	"management-backend/model/vehicle"
	"management-backend/utils"
)

func GetVehicleConfigFromRedis(deviceId int) (*vehicle.VehicleConfig, error) {
	var config vehicle.VehicleConfig
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
