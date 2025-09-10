package data_sync

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/data_sync"
)

func ConfigInfoSer(c *gin.Context) {
	var configInfo data_sync.DataSyncConfig
	_, err := conf.Mysql.Get(&configInfo)
	if err != nil {
		common.ResError(c, "获取数据库配置失败")
		return
	}
	common.ResOk(c, "ok", configInfo)
}

func UpdateConfigSer(c *gin.Context, req data_sync.DataSyncAddReq) {
	var configInfo data_sync.DataSync
	has, err := conf.Mysql.Get(&configInfo)
	if err != nil {
		common.ResError(c, "获取数据库配置失败")
		return
	}

	if !has {
		_, err := conf.Mysql.Insert(&data_sync.DataSyncConfig{
			SqlAddress:   req.SqlAddress,
			SqlDatabase:  req.SqlDatabase,
			SqlUser:      req.SqlUser,
			SqlPassword:  req.SqlPassword,
			SqlTableName: req.SqlTableName,
			SqlRate:      req.SqlRate,
		})
		if err != nil {
			common.ResError(c, "修改配置失败")
			return
		}
	} else {
		_, err := conf.Mysql.Where("id = 1").Update(&data_sync.DataSyncConfig{
			SqlAddress:   req.SqlAddress,
			SqlDatabase:  req.SqlDatabase,
			SqlUser:      req.SqlUser,
			SqlPassword:  req.SqlPassword,
			SqlTableName: req.SqlTableName,
			SqlRate:      req.SqlRate,
		})
		if err != nil {
			common.ResError(c, "修改配置失败")
			return
		}
	}

	common.ResOk(c, "ok", nil)
}
