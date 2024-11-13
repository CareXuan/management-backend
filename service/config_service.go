package service

import (
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetConfigSer(c *gin.Context, searchType int) {
	var configItem model.Config
	_, err := conf.Mysql.Where("type = ?", searchType).Get(&configItem)
	if err != nil {
		common.ResError(c, "获取配置失败")
		return
	}
	common.ResOk(c, "ok", configItem)
}

func SetConfigSer(c *gin.Context, req model.SetConfigReq) {
	var configItem model.Config
	_, err := conf.Mysql.Where("type = ?", req.Type).Get(&configItem)
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "获取配置失败")
		return
	}
	if configItem.Id != 0 {
		_, err := conf.Mysql.Where("type = ?", req.Type).Update(&model.Config{
			Value: req.Value,
		})
		if err != nil {
			common.ResError(c, "写入配置失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&model.Config{
			Type:  req.Type,
			Value: req.Value,
		})
		if err != nil {
			common.ResError(c, "写入配置失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
