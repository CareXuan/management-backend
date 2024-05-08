package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"management-backend/common"
	"management-backend/utils"
	"strconv"
)

func GetVehicleConfig(configUrl string) {
	var reqParams = make(map[string]string)
	reqParams["current"] = "0"
	reqParams["page_size"] = "9999"

	res, err := common.DoGet(configUrl, reqParams)
	if err != nil {
		log.Fatal(err)
		return
	}
	configData := res.Body.(map[string]any)["list"]
	for _, v := range configData.([]interface{}) {
		data := v.(map[string]interface{})
		deviceId, _ := strconv.Atoi(data["deviceId"].(string))
		deviceConfigJson, _ := json.Marshal(data)
		_, err := Redis.Do("SET", fmt.Sprintf(utils.REDIS_KEY_VEHICLE_CONFIG, deviceId), string(deviceConfigJson))
		if err != nil {
			log.Fatal(err)
		}
	}
}
