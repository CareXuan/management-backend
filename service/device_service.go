package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
	"strconv"
	"time"
)

func DeviceListSer(c *gin.Context, deviceId, name, phone string, page, pageSize int) {
	userId := c.Param("user_id")
	var userRoleItem model.UserRole
	_, err := conf.Mysql.Where("user_id = ?", userId).Get(&userRoleItem)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var deviceItems []*model.Device
	sess := conf.Mysql.NewSession()
	needFilter := false
	var province []string
	var city []string
	var zone []string
	// 1 超级管理员 8 商户
	if userRoleItem.RoleId != 1 && userRoleItem.RoleId != 8 {
		needFilter = true
		var organizationUser model.OrganizationUser
		_, err := conf.Mysql.Where("user_id = ?", userId).Get(&organizationUser)
		if err != nil {
			common.ResError(c, "获取用户关联组织失败")
			return
		}
		var organizationItem []*model.Organization
		err = conf.Mysql.Where("id = ?", organizationUser.OrganizationId).Find(&organizationItem)
		if err != nil {
			common.ResError(c, "获取组织信息失败")
			return
		}
		for _, i := range organizationItem {
			province = append(province, i.Province)
			city = append(city, i.City)
			zone = append(zone, i.Zone)
		}
	} else if userRoleItem.RoleId == 8 {
		var userItem model.User
		_, err := conf.Mysql.Where("id = ?", userId).Get(&userItem)
		if err != nil {
			common.ResError(c, "获取用户信息失败")
			return
		}
		sess.Where("phone = ?", userItem.Phone)
	}
	if needFilter {
		sess.In("province", province)
		sess.In("city", city)
		sess.In("zone", zone)
	}
	if deviceId != "" {
		sess.Where("device_id = ?", deviceId)
	}
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	if phone != "" {
		sess.Where("phone = ?", phone)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&deviceItems)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	isSupervisor, err := IsSupervisor(c)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	for _, i := range deviceItems {
		i.IsSupervisor = isSupervisor
		if isSupervisor != 1 {
			i.Pwd1 = 0
			i.Pwd1Base = 0
			i.Pwd2 = 0
			i.Pwd2Base = 0
			i.Pwd3 = 0
			i.Pwd3Base = 0
			i.Pwd4 = 0
			i.Pwd4Base = 0
			i.Ip1 = ""
			i.Port1 = 0
			i.Ip2 = ""
			i.Port2 = 0
			i.PwdChangeDate = 0
		}
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: deviceItems})
}

func AddDeviceSer(c *gin.Context, req model.DeviceAddReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(model.Device{
			DeviceId: req.DeviceId,
			Name:     req.Name,
			Province: req.Province,
			City:     req.City,
			Zone:     req.Zone,
			Address:  req.Address,
			Manager:  req.Manager,
			Phone:    req.Phone,
			Remark:   req.Remark,
			UpdateTs: time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			common.ResError(c, "写入设备信息失败")
			return
		}
	} else {
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
			Name:     req.Name,
			Province: req.Province,
			City:     req.City,
			Zone:     req.Zone,
			Address:  req.Address,
			Manager:  req.Manager,
			Phone:    req.Phone,
			Remark:   req.Remark,
			Ts:       time.Now().Format("2006-01-02 15:04:05"),
			UpdateTs: time.Now().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			common.ResError(c, "写入设备信息失败")
			return
		}
	}
	var userExist model.User
	has, err := conf.Mysql.Where("phone = ?", req.Phone).Get(&userExist)
	if err != nil {
		common.ResError(c, "查询用户失败")
		return
	}
	if !has {
		newToken, err := utils.GenerateRandomToken(20)
		if err != nil {
			common.ResError(c, "生成token失败")
			return
		}
		var newUser = model.User{
			Name:     req.Manager,
			Password: "123456",
			Phone:    req.Phone,
			Token:    newToken,
		}
		_, err = conf.Mysql.Insert(&newUser)
		if err != nil {
			common.ResError(c, "生成用户失败")
			return
		}
		_, err = conf.Mysql.Insert(&model.UserRole{
			UserId: newUser.Id,
			RoleId: 8,
		})
		if err != nil {
			common.ResError(c, "分配用户角色失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func UpdateSpecialInfoSer(c *gin.Context, req model.UpdateSpecialInfoReq) {
	_, err := conf.Mysql.Where("device_id = ?", req.DeviceId).Update(model.Device{
		HeartBeat:     req.HeartBeat,
		HeartBeatMin:  req.HeartBeatMin,
		Pwd1:          req.Pwd1,
		Pwd2:          req.Pwd2,
		Pwd3:          req.Pwd3,
		Pwd4:          req.Pwd4,
		Pwd1Base:      req.Pwd1Base,
		Pwd2Base:      req.Pwd2Base,
		Pwd3Base:      req.Pwd3Base,
		Pwd4Base:      req.Pwd4Base,
		Ip1:           req.Ip1,
		Port1:         req.Port1,
		Ip2:           req.Ip2,
		Port2:         req.Port2,
		PwdChangeDate: req.PwdChangeDate,
	})
	if err != nil {
		common.ResError(c, "修改数据失败")
		return
	}
	var msgData = make(map[string]interface{})
	msgData["beat_time"] = req.HeartBeat
	msgData["beat_min_time"] = req.HeartBeatMin
	msgData["statistic_pwd"] = req.Pwd1
	msgData["statistic_base"] = req.Pwd1Base
	msgData["pwd1"] = req.Pwd2
	msgData["pwd1_base"] = req.Pwd2Base
	msgData["pwd2"] = req.Pwd3
	msgData["pwd2_base"] = req.Pwd3Base
	msgData["pwd3"] = req.Pwd4
	msgData["pwd3_base"] = req.Pwd4Base
	msgData["ip1"] = req.Ip1
	msgData["port1"] = req.Port1
	msgData["ip2"] = req.Ip2
	msgData["port2"] = req.Port2
	msgData["pwd_change_date"] = req.PwdChangeDate
	msg, err := json.Marshal(msgData)
	if err != nil {
		common.ResError(c, "数据转换失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host1, conf.Conf.Tcp.Port1, 1, req.DeviceId, 1, string(msg))
	if err != nil {
		common.ResError(c, "发送控制命令1失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host2, conf.Conf.Tcp.Port2, 1, req.DeviceId, 1, string(msg))
	if err != nil {
		common.ResError(c, "发送控制命令2失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host3, conf.Conf.Tcp.Port3, 1, req.DeviceId, 1, string(msg))
	if err != nil {
		common.ResError(c, "发送控制命令2失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func ReadSpecialInfoSer(c *gin.Context, req model.ReadSpecialInfoReq) {
	err := common.CommonSendDeviceReport(conf.Conf.Tcp.Host1, conf.Conf.Tcp.Port1, 2, req.DeviceId, 1, "")
	if err != nil {
		common.ResError(c, "发送控制命令1失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host2, conf.Conf.Tcp.Port2, 2, req.DeviceId, 1, "")
	if err != nil {
		common.ResError(c, "发送控制命令2失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host3, conf.Conf.Tcp.Port3, 2, req.DeviceId, 1, "")
	if err != nil {
		common.ResError(c, "发送控制命令2失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func GetSpecialInfoLogSer(c *gin.Context, deviceId string) {
	var specialInfoLog []*model.DeviceChangeLog
	err := conf.Mysql.Where("device_id = ?", deviceId).Where("type = 2").OrderBy("id DESC").Find(&specialInfoLog)
	if err != nil {
		common.ResError(c, "获取日志失败")
		return
	}
	common.ResOk(c, "ok", specialInfoLog)
}

func DeviceInfoSer(c *gin.Context, deviceId int) {
	var deviceItem model.Device
	_, err := conf.Mysql.Where("device_id = ?", deviceId).Get(&deviceItem)
	if err != nil {
		common.ResError(c, "获取设备详情失败")
		return
	}
	caLat, caLng := common.WGS84ToGCJ02(common.ConvertBDSLatitude(deviceItem.Latitude), common.ConvertBDSLatitude(deviceItem.Longitude))
	deviceItem.RealLocation, err = common.ReGeoCode(strconv.FormatFloat(caLng, 'f', -1, 64)+","+strconv.FormatFloat(caLat, 'f', -1, 64), 1000, "all")
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

func DeviceNewServiceDataSer(c *gin.Context, deviceId, page, pageSize int) {
	var serviceDatas []*model.DeviceNewServiceData
	sess := conf.Mysql.NewSession()
	sess.Where("device_id=?", deviceId)
	sess.Where("current_h != 255")
	sess.Where("voltage_h != 255")
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
		caLat, caLng := common.WGS84ToGCJ02(common.ConvertBDSLatitude(i.Latitude), common.ConvertBDSLatitude(i.Longitude))
		locationRes = append(locationRes, model.DeviceLocationRes{
			DeviceId: i.DeviceId,
			Name:     i.Name,
			Iccid:    i.Iccid,
			Lat:      strconv.FormatFloat(caLat, 'f', -1, 64),
			Lng:      strconv.FormatFloat(caLng, 'f', -1, 64),
			//Lat: i.Latitude,
			//Lng: i.Longitude,
		})
	}
	common.ResOk(c, "ok", locationRes)
}

func DeviceStatisticSer(c *gin.Context, deviceId int, startTime, endTime string) {
	var serviceData []*model.DeviceNewServiceData
	var statisticRes model.DeviceStatisticRes
	sess := conf.Mysql.NewSession()
	var highVoltage []interface{}
	var highCurrent []interface{}
	var warningCnt []interface{}
	var power []interface{}
	var status1 []interface{}
	var status2 []interface{}
	var voltage []interface{}
	var current []interface{}
	var temperature []interface{}
	var signalIntensity []interface{}
	sess.Where("device_id = ?", deviceId)
	sess.Where("ts >= ? and ts <= ?", startTime, endTime)
	err := sess.Find(&serviceData)
	if err != nil {
		common.ResError(c, "获取数据失败")
		return
	}
	for _, i := range serviceData {
		statisticRes.Columns = append(statisticRes.Columns, i.Ts)
		highVoltage = append(highVoltage, fmt.Sprintf("%.2f", float64((i.HighVoltageH<<8+i.HighVoltageL)*15)/4095.0))
		highCurrent = append(highCurrent, fmt.Sprintf("%.2f", float64((i.HighCurrentH<<8+i.HighCurrentL)*20)/4095.0))
		warningCnt = append(warningCnt, i.WarningCount)
		power = append(power, fmt.Sprintf("%.2f", float64((i.PowerH<<8)+i.Power)/100.0))
		status1Item := 0
		if i.IoStatus&64 > 0 {
			status1Item = 1
		}
		status1 = append(status1, status1Item)
		status2Item := 0
		if i.IoStatus&32 > 0 {
			status2Item = 1
		}
		status2 = append(status2, status2Item)
		voltage = append(voltage, fmt.Sprintf("%.2f", float64((i.VoltageH<<8)+i.Voltage)/100.0))
		current = append(current, fmt.Sprintf("%.2f", float64((i.CurrentH<<8)+i.Current)/100.0))
		tem := (i.Tem1h << 8) + i.Tem1l
		var realTem float64
		if tem < 0x8000 {
			realTem = float64(tem) * 0.0625
		} else {
			tem = -((^tem & 0xffff) + 1)
			realTem = float64(tem) * 0.0625
		}
		if realTem >= -55 {
			temperature = append(temperature, realTem)
		}
		signalIntensity = append(signalIntensity, i.SignalIntensity)
	}
	statisticRes.Datas = append(statisticRes.Datas, highVoltage)
	statisticRes.Datas = append(statisticRes.Datas, highCurrent)
	statisticRes.Datas = append(statisticRes.Datas, warningCnt)
	statisticRes.Datas = append(statisticRes.Datas, power)
	statisticRes.Datas = append(statisticRes.Datas, status1)
	statisticRes.Datas = append(statisticRes.Datas, status2)
	statisticRes.Datas = append(statisticRes.Datas, voltage)
	statisticRes.Datas = append(statisticRes.Datas, current)
	statisticRes.Datas = append(statisticRes.Datas, temperature)
	statisticRes.Datas = append(statisticRes.Datas, signalIntensity)
	common.ResOk(c, "ok", statisticRes)
}

func GetDeviceAllWarningSer(c *gin.Context) {
	isSupervisor, err := IsSupervisor(c)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	if isSupervisor == 0 {
		common.ResOk(c, "ok", nil)
		return
	}
	var warnings []*model.DeviceChangeLog
	err = conf.Mysql.Where("has_all_warning = ?", 0).Find(&warnings)
	if err != nil {
		common.ResError(c, "获取报警记录失败")
		return
	}
	_, err = conf.Mysql.Where("has_all_warning = ?", 0).Update(model.DeviceChangeLog{
		HasAllWarning: 1,
	})
	if err != nil {
		common.ResError(c, "修改报警记录状态失败")
		return
	}
	common.ResOk(c, "ok", warnings)
}

func GetDeviceSingleWarningSer(c *gin.Context, deviceId int) {
	isSupervisor, err := IsSupervisor(c)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	if isSupervisor == 0 {
		common.ResOk(c, "ok", nil)
		return
	}
	var warnings []*model.DeviceChangeLog
	err = conf.Mysql.Where("device_id = ?", deviceId).Where("has_single_warning = ?", 0).Find(&warnings)
	if err != nil {
		common.ResError(c, "获取报警记录失败")
		return
	}
	_, err = conf.Mysql.Where("device_id = ?", deviceId).Where("has_single_warning = ?", 0).Update(model.DeviceChangeLog{
		HasSingleWarning: 1,
	})
	if err != nil {
		common.ResError(c, "修改报警记录状态失败")
		return
	}
	common.ResOk(c, "ok", warnings)
}

func DeviceLocationHistorySer(c *gin.Context, deviceId int) {

}

func DeviceReportSer(c *gin.Context, reportType int, deviceId int, controlType int, msg string) {
	err := common.CommonSendDeviceReport(conf.Conf.Tcp.Host1, conf.Conf.Tcp.Port1, reportType, deviceId, controlType, msg)
	if err != nil {
		common.ResError(c, "发送控制命令1失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host2, conf.Conf.Tcp.Port2, reportType, deviceId, controlType, msg)
	if err != nil {
		common.ResError(c, "发送控制命令2失败")
		return
	}
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host3, conf.Conf.Tcp.Port3, reportType, deviceId, controlType, msg)
	if err != nil {
		common.ResError(c, "发送控制命令3失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
