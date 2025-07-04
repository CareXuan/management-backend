package opcua

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/opcua"
	"switchboard-backend/utils"
)

func ListSer(c *gin.Context, name string, page, pageSize int) {
	var opcuaList []*opcua.Opcua
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	err := sess.Where("deleted_at = 0").Find(&opcuaList)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", opcuaList)
}

func InfoSer(c *gin.Context, id, page, pageSize int) {
	var opcuaData []*opcua.OpcuaData
	count, err := conf.Mysql.Where("device_id = ?", id).Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&opcuaData)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: opcuaData})
}

func AddSer(c *gin.Context, req opcua.AddOpcuaReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.MustCols("name", "ip", "port", "is_enable").Where("id = ?", req.Id).Update(&opcua.Opcua{
			Name:     req.Name,
			Ip:       req.Ip,
			Port:     req.Port,
			IsEnable: req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "修改设备信息失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&opcua.Opcua{
			Name:     req.Name,
			Ip:       req.Ip,
			Port:     req.Port,
			IsEnable: req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "添加设备信息失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
