package modbus

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/modbus"
	"switchboard-backend/utils"
)

func ListSer(c *gin.Context, name string, page, pageSize int) {
	var modbusItems []*modbus.Modbus
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&modbusItems)
	if err != nil {
		common.ResError(c, "获取设备失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: modbusItems})
}

func InfoSer(c *gin.Context, id int) {
	var modbusItem modbus.Modbus
	_, err := conf.Mysql.Where("id = ?", id).Get(&modbusItem)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	common.ResOk(c, "ok", modbusItem)
}

func AddSer(c *gin.Context, req modbus.AddModbusReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(&modbus.Modbus{
			Name:         req.Name,
			ConnectType:  req.ConnectType,
			Ip:           req.Ip,
			Port:         req.Port,
			BaudRate:     req.BaudRate,
			StartAddress: req.StartAddress,
			Count:        req.Count,
			Slave:        req.Slave,
			IsEnable:     req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&modbus.Modbus{
			Name:         req.Name,
			ConnectType:  req.ConnectType,
			Ip:           req.Ip,
			Port:         req.Port,
			BaudRate:     req.BaudRate,
			StartAddress: req.StartAddress,
			Count:        req.Count,
			Slave:        req.Slave,
			IsEnable:     req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
