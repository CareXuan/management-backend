package port

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/port"
	"switchboard-backend/utils"
)

func BridgeListSer(c *gin.Context, name string, page, pageSize int) {
	var bridgeItems []*port.Bridge
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&bridgeItems)
	if err != nil {
		common.ResError(c, "获取网桥列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: bridgeItems})
}

func BridgeInfoSer(c *gin.Context, id int) {
	var bridgeItem port.Bridge
	_, err := conf.Mysql.Where("id = ?", id).Get(&bridgeItem)
	if err != nil {
		common.ResError(c, "获取网桥信息失败")
		return
	}
	common.ResOk(c, "ok", bridgeItem)
}

func BridgeAddSer(c *gin.Context, req port.AddBridgeReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(&port.Bridge{
			Name:        req.Name,
			Ip:          req.Ip,
			EnglishName: req.EnglishName,
		})
		if err != nil {
			common.ResError(c, "修改网桥失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&port.Bridge{
			Name:        req.Name,
			Ip:          req.Ip,
			EnglishName: req.EnglishName,
		})
		if err != nil {
			common.ResError(c, "新增网桥失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
