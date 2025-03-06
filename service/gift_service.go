package service

import (
	"github.com/gin-gonic/gin"
	"prize-draw/common"
	"prize-draw/conf"
	"prize-draw/model"
	"prize-draw/utils"
	"time"
)

func List(c *gin.Context, name string, page, pageSize int) {
	var gifts []*model.Gift
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&gifts)
	if err != nil {
		common.ResError(c, "获取礼物列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: gifts})
}

func GiftInfo(c *gin.Context, id int) {
	var giftItem model.Gift
	_, err := conf.Mysql.Where("id = ?", id).Where("delete_at = 0").Get(&giftItem)
	if err != nil {
		common.ResError(c, "获取礼物详情失败")
		return
	}
	common.ResOk(c, "ok", giftItem)
}

func GroupList(c *gin.Context, name string, page, pageSize int) {
	var giftGroups []*model.GiftGroup
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&giftGroups)
	if err != nil {
		common.ResError(c, "获取礼物列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: giftGroups})
}

func GroupInfo(c *gin.Context, id int) {
	var groupItem model.GiftGroup
	_, err := conf.Mysql.Where("id = ?", id).Where("delete_at = 0").Get(&groupItem)
	if err != nil {
		common.ResError(c, "获取礼物组失败")
		return
	}
	var groupGifts []*model.GiftGroupGift
	err = conf.Mysql.Where("group_id = ?", id).Where("delete_at = 0").Find(&groupGifts)
	if err != nil {
		common.ResError(c, "获取礼物组关联礼物失败")
		return
	}
	var giftIds []int
	for _, i := range groupGifts {
		giftIds = append(giftIds, i.GiftId)
	}
	var gifts []*model.Gift
	err = conf.Mysql.In("id", giftIds).Where("delete_at = 0").Find(&gifts)
	if err != nil {
		common.ResError(c, "获取礼物失败")
		return
	}
	groupItem.GroupGift = gifts
	common.ResOk(c, "ok", groupItem)
}

func Add(c *gin.Context, req model.GiftAddReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(model.Gift{
			Name:        req.Name,
			Pic:         req.Pic,
			Description: req.Description,
			Level:       req.Level,
			Consumable:  req.Consumable,
			Show:        req.Show,
			CanObtain:   req.CanObtain,
			CrushCnt:    model.GIFT_LEVEL_MAPPING[req.Level],
			CreateAt:    int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "修改礼物失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&model.Gift{
			Name:        req.Name,
			Pic:         req.Pic,
			Description: req.Description,
			Level:       req.Level,
			Consumable:  req.Consumable,
			Show:        req.Show,
			CanObtain:   req.CanObtain,
			CrushCnt:    model.GIFT_LEVEL_MAPPING[req.Level],
			CreateAt:    int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "添加礼物失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func Delete(c *gin.Context, req model.GiftDeleteReq) {
	_, err := conf.Mysql.In("id", req.Ids).Update(&model.Gift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除礼物失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func ChangeReq(c *gin.Context, req model.GiftChangeStatusReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Where("delete_at = 0").Update(&model.Gift{
		CanObtain: req.Status,
	})
	if err != nil {
		common.ResError(c, "修改状态失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func GroupAdd(c *gin.Context, req model.GiftGroupAdd) {
	if req.Id != 0 {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("id = ?", req.Id).Update(model.GiftGroup{
			Name:      req.Name,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Status:    req.Status,
		})
		if err != nil {
			common.ResError(c, "修改礼物组失败")
			return
		}
		_, err = sess.Where("group_id = ?", req.Id).Delete(&model.GiftGroupGift{})
		if err != nil {
			common.ResError(c, "删除已有关联关系失败")
			return
		}
		var giftGroupItems []*model.GiftGroupGift
		for _, i := range req.GiftIds {
			giftGroupItems = append(giftGroupItems, &model.GiftGroupGift{
				GroupId:  req.Id,
				GiftId:   i,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = sess.Insert(giftGroupItems)
		if err != nil {
			common.ResError(c, "绑定礼物组失败")
			return
		}
		sess.Commit()
	} else {
		var groupAdd = &model.GiftGroup{
			Name:      req.Name,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Status:    req.Status,
			CreateAt:  int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(groupAdd)
		if err != nil {
			common.ResError(c, "添加礼物组失败")
			return
		}
		var giftGroupItems []*model.GiftGroupGift
		for _, i := range req.GiftIds {
			giftGroupItems = append(giftGroupItems, &model.GiftGroupGift{
				GroupId:  groupAdd.Id,
				GiftId:   i,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(giftGroupItems)
		if err != nil {
			common.ResError(c, "绑定礼物组失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func GroupDelete(c *gin.Context, req model.GiftGroupDelete) {
	sess := conf.Mysql.NewSession()
	_, err := sess.In("id", req.Ids).Where("delete_at = 0").Update(&model.GiftGroup{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除礼物组失败")
		sess.Rollback()
		return
	}
	_, err = sess.In("group_id = ?", req.Ids).Where("delete_at = 0").Update(&model.GiftGroupGift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除礼物组关联关系失败")
		sess.Rollback()
		return
	}
	common.ResOk(c, "ok", nil)
}

/*=====================================app=====================================*/
