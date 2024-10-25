package service

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
	"time"
)

func DeviceListSer(c *gin.Context, page, pageSize int) {
	var deviceItems []*model.Device
	sess := conf.Mysql.NewSession()
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&deviceItems)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: deviceItems})
}

func AddDeviceSer(c *gin.Context, req model.DeviceAddReq) {
	var deviceExist model.Device
	_, err := conf.Mysql.Where("device_id = ?", req.DeviceId).Get(&deviceExist)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	if deviceExist.Id != 0 {
		common.ResForbidden(c, "当前设备ID已注册")
		return
	}
	_, err = conf.Mysql.Insert(model.Device{
		DeviceId: req.DeviceId,
		Ts:       time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		common.ResError(c, "写入设备信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func DeviceInfoSer(c *gin.Context, deviceId int) {
	var deviceItem model.Device
	_, err := conf.Mysql.Where("device_id = ?", deviceId).Get(&deviceItem)
	if err != nil {
		common.ResError(c, "获取设备详情失败")
		return
	}
	common.ResOk(c, "ok", deviceItem)
}

func DeviceCommonDataSer(c *gin.Context, deviceId, page, pageSize int) {
	var commonDatas []*model.DeviceCommonData
	sess := conf.Mysql.NewSession()
	sess.Where("device_id=?", deviceId)
	count, err := sess.OrderBy("id DESC").Limit(pageSize, (page-1)*pageSize).FindAndCount(&commonDatas)
	if err != nil {
		common.ResError(c, "获取通用数据失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: commonDatas})
}

func DeviceServiceDataSer(c *gin.Context, deviceId, page, pageSize int) {
	var serviceDatas []*model.DeviceServiceData
	sess := conf.Mysql.NewSession()
	sess.Where("device_id=?", deviceId)
	count, err := sess.OrderBy("id DESC").Limit(pageSize, (page-1)*pageSize).FindAndCount(&serviceDatas)
	if err != nil {
		common.ResError(c, "获取通用数据失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: serviceDatas})
}

func DeviceLocationHistorySer(c *gin.Context, deviceId int) {

}
