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

func AllDeviceLocationSer(c *gin.Context) {
	var deviceItems []*model.Device
	err := conf.Mysql.Find(&deviceItems)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	var locationRes []model.DeviceLocationRes
	for _, i := range deviceItems {
		if i.Latitude == "" && i.Longitude == "" {
			continue
		}
		locationRes = append(locationRes, model.DeviceLocationRes{
			DeviceId: i.DeviceId,
			Iccid:    i.Iccid,
			Lat:      i.Latitude,
			Lng:      i.Longitude,
		})
	}
	common.ResOk(c, "ok", locationRes)
}

func DeviceStatisticSer(c *gin.Context, deviceId, statisticType int, startTime, endTime string) {
	var serviceDatas []*model.DeviceServiceData
	sess := conf.Mysql.NewSession()
	sess.Where("device_id = ?", deviceId)
	sess.Where("ts >= ? and ts <= ?", startTime, endTime)
	err := sess.Find(&serviceDatas)
	if err != nil {
		common.ResError(c, "获取数据失败")
		return
	}
	var statisticRes model.DeviceStatisticRes
	var datah []int
	var datal []int
	for _, i := range serviceDatas {
		statisticRes.Columns = append(statisticRes.Columns, i.Ts)
		switch statisticType {
		case 1:
			datah = append(datah, i.HighVoltageH)
			datal = append(datal, i.HighVoltageL)
		case 2:
			datah = append(datah, i.HighCurrentH)
			datal = append(datal, i.HighCurrentL)
		case 3:
			datah = append(datah, i.SwitchCurrent)
		}
	}
	statisticRes.Datas = append(statisticRes.Datas, datah)
	statisticRes.Datas = append(statisticRes.Datas, datal)
	common.ResOk(c, "ok", statisticRes)
}

func DeviceLocationHistorySer(c *gin.Context, deviceId int) {

}
