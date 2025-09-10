package dlt645

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/dlt645"
	"switchboard-backend/utils"
)

func ListSer(c *gin.Context, name string, page, pageSize int) {
	var dlt645Items []*dlt645.Dlt645
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&dlt645Items)
	if err != nil {
		common.ResError(c, "获取设备失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: dlt645Items})
}

func InfoSer(c *gin.Context, id int) {
	var dlt645Item dlt645.Dlt645
	_, err := conf.Mysql.Where("id = ?", id).Get(&dlt645Item)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	common.ResOk(c, "ok", dlt645Item)
}

func AddSer(c *gin.Context, req dlt645.AddDlt645Req) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(&dlt645.Dlt645{
			Name:        req.Name,
			ConnectType: req.ConnectType,
			Ip:          req.Ip,
			Port:        req.Port,
			BaudRate:    req.BaudRate,
			Address:     req.Address,
			IsEnable:    req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&dlt645.Dlt645{
			Name:        req.Name,
			ConnectType: req.ConnectType,
			Ip:          req.Ip,
			Port:        req.Port,
			BaudRate:    req.BaudRate,
			Address:     req.Address,
			IsEnable:    req.IsEnable,
		})
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
