package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
)

func DeviceListSer(c *gin.Context, page int, pageSize int, name string, iccid string, deviceType string, status string) {
	var devices []*model.Device

	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name=?", name)
	}
	if iccid != "" {
		sess.Where("iccid=?", iccid)
	}
	if deviceType != "" {
		sess.Where("type=?", deviceType)
	}
	if status != "" {
		sess.Where("status=?", status)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&devices)
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: devices})
}

func SignalDetailListSer(c *gin.Context, page int, pageSize int, deviceId string) {
	var signals []*model.SignalData
	sess := conf.Mysql.NewSession()
	count, err := sess.Where("device_id=?", deviceId).OrderBy("id desc").Limit(pageSize, (page-1)*pageSize).FindAndCount(&signals)
	if err != nil {
		common.ResError(c, "获取设备信号列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: signals})
}
