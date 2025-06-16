package siemens

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/siemens"
	"switchboard-backend/utils"
)

func ListSer(c *gin.Context, name string, page, pageSize int) {
	var s7List []*siemens.SiemensS7
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	err := sess.Where("deleted_at = 0").Find(&s7List)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", s7List)
}

func InfoSer(c *gin.Context, id, page, pageSize int) {
	var s7Data []*siemens.SiemensS7Data
	count, err := conf.Mysql.Where("device_id = ?", id).Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&s7Data)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: s7Data})
}

func AddSer(c *gin.Context, req siemens.AddSiemensS7Req) {
	if req.Id != 0 {
		_, err := conf.Mysql.MustCols("name", "ip", "rack", "slot", "is_enable").Where("id = ?", req.Id).Update(&siemens.SiemensS7{
			Name:     req.Name,
			Ip:       req.Ip,
			Rack:     req.Rack,
			Slot:     req.Slot,
			IsEnable: req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "修改设备信息失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&siemens.SiemensS7{
			Name:     req.Name,
			Ip:       req.Ip,
			Rack:     req.Rack,
			Slot:     req.Slot,
			IsEnable: req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "添加设备信息失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func AddSiemensDataSer(c *gin.Context, req siemens.AddSiemensDataReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.MustCols("device_id", "type", "start", "db_num").Where("id = ?", req.Id).Update(&siemens.SiemensS7Data{
			DeviceId: req.DeviceId,
			Type:     req.Type,
			Start:    req.Start,
			DbNum:    req.DbNum,
		})
		if err != nil {
			common.ResError(c, "修改信息失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&siemens.SiemensS7Data{
			DeviceId: req.DeviceId,
			Type:     req.Type,
			Start:    req.Start,
			DbNum:    req.DbNum,
		})
		if err != nil {
			common.ResError(c, "添加信息失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
