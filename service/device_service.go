package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetDeviceListSer(c *gin.Context, name string, page, pageSize int) {
	var devices []*model.Device
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&devices)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: devices})
}

func GetDeviceInfoSer(c *gin.Context, id int) {
	var device model.Device
	_, err := conf.Mysql.Where("id = ?", id).Get(&device)
	if err != nil {
		common.ResError(c, "获取设备详情失败")
		return
	}
	common.ResOk(c, "ok", device)
}

func AddDeviceSer(c *gin.Context, device model.Device) {
	if device.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", device.Id).Update(device)
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(model.Device{
			Name:      device.Name,
			Pic:       device.Pic,
			Cert:      device.Cert,
			UseTime:   device.UseTime,
			CreatedAt: time.Now().Unix(),
		})
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
