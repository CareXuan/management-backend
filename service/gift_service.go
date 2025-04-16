package service

import (
	"fmt"
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

func GiftRemain(c *gin.Context, level string, giftExist int) {
	var giftRemainRes = model.GiftRemainRes{}
	var levels = []string{"总量", "E", "D", "C", "B", "A"}
	var r = make(map[string]*model.GiftProbabilityRes)
	for _, i := range levels {
		r[i] = &model.GiftProbabilityRes{
			GetCount: 0,
			AllCount: 0,
		}
	}
	giftRemainRes.Probability = r
	var giftItem []*model.Gift
	sess := conf.Mysql.NewSession()
	if level != "" {
		sess.Where("level = ?", level)
	}
	err := sess.Where("delete_at = 0").Find(&giftItem)
	if err != nil {
		common.ResError(c, "获取礼物信息失败")
		return
	}
	var giftPackage []*model.GiftPackage
	err = conf.Mysql.Where("delete_at = 0").Find(&giftPackage)
	if err != nil {
		common.ResError(c, "获取背包信息失败")
		return
	}

	var giftPackageMapping = make(map[int]*model.GiftPackage)
	for _, i := range giftPackage {
		giftPackageMapping[i.GiftId] = i
	}
	for _, i := range giftItem {
		exist := 2
		if _, ok := giftPackageMapping[i.Id]; ok {
			exist = 1
			giftRemainRes.Probability[i.Level].GetCount += 1
			giftRemainRes.Probability["总量"].GetCount += 1
		}
		giftRemainRes.Probability[i.Level].AllCount += 1
		giftRemainRes.Probability["总量"].AllCount += 1
		if giftExist != 0 && giftExist != exist {
			continue
		}
		giftRemainRes.List = append(giftRemainRes.List, &model.GiftRemainListItem{
			GiftId: i.Id,
			Name:   i.Name,
			Pic:    i.Pic,
			Exist:  exist,
		})
	}
	common.ResOk(c, "ok", giftRemainRes)
}

func RaffleConfig(c *gin.Context) {
	var config model.Config
	_, err := conf.Mysql.Where("id = 1").Get(&config)
	if err != nil {
		common.ResError(c, "获取配置失败")
		return
	}
	common.ResOk(c, "ok", config)
}

func RaffleConfigSet(c *gin.Context, req model.GiftConfigSetReq) {
	_, err := conf.Mysql.Where("id > 0").Update(model.Config{
		OnePoint: req.OneCount,
		TenPoint: req.TenCount,
	})
	if err != nil {
		common.ResError(c, "修改配置失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func PointLeft(c *gin.Context) {
	var pointLeft model.GiftPackage
	_, err := conf.Mysql.Where("gift_id = 0").Where("delete_at = 0").Get(&pointLeft)
	if err != nil {
		common.ResError(c, "获取抽卡点失败")
		return
	}
	common.ResOk(c, "ok", pointLeft)
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
	groupItem.GroupGift = groupGifts
	common.ResOk(c, "ok", groupItem)
}

func AlbumList(c *gin.Context, name string, page, pageSize int) {
	var albums []*model.Album
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&albums)
	if err != nil {
		common.ResError(c, "获取相册列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: albums})
}

func AlbumInfo(c *gin.Context, id int) {
	var albumItem model.Album
	_, err := conf.Mysql.Where("id = ?", id).Where("delete_at = 0").Get(&albumItem)
	if err != nil {
		common.ResError(c, "获取相册失败")
		return
	}
	var albumGifts []*model.AlbumGift
	err = conf.Mysql.Where("album_id = ?", id).Where("delete_at = 0").Find(&albumGifts)
	if err != nil {
		common.ResError(c, "获取相册关联礼物失败")
		return
	}
	albumItem.AlbumGift = albumGifts
	common.ResOk(c, "ok", albumItem)
}

func AlbumGift(c *gin.Context, albumId int) {
	var existGift []*model.AlbumGift
	sess := conf.Mysql.NewSession()
	if albumId != 0 {
		sess.Where("album_id != ?", albumId)
	}
	err := sess.Where("delete_at = 0").Find(&existGift)
	if err != nil {
		common.ResError(c, "获取已关联相册相片失败")
		return
	}
	var existId []int
	for _, i := range existGift {
		existId = append(existId, i.GiftId)
	}
	var giftItems []*model.Gift
	err = conf.Mysql.NotIn("id", existId).Where("delete_at = 0").Find(&giftItems)
	if err != nil {
		common.ResError(c, "获取礼物列表失败")
		return
	}
	common.ResOk(c, "ok", giftItems)
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

func AddPoint(c *gin.Context, req model.GiftAddPointReq) {
	var pointItem model.GiftPackage
	_, err := conf.Mysql.Where("gift_id = 0").Where("delete_at = 0").Get(&pointItem)
	if err != nil {
		common.ResError(c, "获取抽卡点失败")
		return
	}
	_, err = conf.Mysql.Where("gift_id = 0").Where("delete_at = 0").Update(&model.GiftPackage{
		Count: pointItem.Count + req.Count,
	})
	if err != nil {
		common.ResError(c, "发放抽卡点失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func ResetPoint(c *gin.Context) {
	_, err := conf.Mysql.MustCols("gift_id", "count").Where("gift_id = 0").Update(&model.GiftPackage{
		Count: 0,
	})
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "重置抽卡点失败")
		return
	}
	_, err = conf.Mysql.Where("gift_id != 0").Delete(&model.GiftPackage{})
	if err != nil {
		common.ResError(c, "重置礼物获取情况失败")
		return
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
		_, err = sess.Where("group_id = ?", req.Id).Update(&model.GiftGroupGift{
			DeleteAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "删除已有关联关系失败")
			return
		}
		var giftGroupItems []*model.GiftGroupGift
		for _, i := range req.GiftIds {
			giftGroupItems = append(giftGroupItems, &model.GiftGroupGift{
				GroupId:     req.Id,
				Level:       i.Level,
				Probability: i.Probability,
				CreateAt:    int(time.Now().Unix()),
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
				GroupId:     groupAdd.Id,
				Level:       i.Level,
				Probability: i.Probability,
				CreateAt:    int(time.Now().Unix()),
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
	_, err = sess.In("group_id", req.Ids).Where("delete_at = 0").Update(&model.GiftGroupGift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除礼物组关联关系失败")
		sess.Rollback()
		return
	}
	sess.Commit()
	common.ResOk(c, "ok", nil)
}

func AlbumAdd(c *gin.Context, req model.GiftAlbumAdd) {
	if req.Id != 0 {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("id = ?", req.Id).Update(model.Album{
			Name: req.Name,
		})
		if err != nil {
			common.ResError(c, "修改相册失败")
			return
		}
		_, err = sess.Where("album_id = ?", req.Id).Update(&model.AlbumGift{
			DeleteAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "删除已有关联关系失败")
			return
		}
		var giftAlbumItems []*model.AlbumGift
		for _, i := range req.GiftIds {
			giftAlbumItems = append(giftAlbumItems, &model.AlbumGift{
				AlbumId:  req.Id,
				GiftId:   i,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = sess.Insert(giftAlbumItems)
		if err != nil {
			common.ResError(c, "绑定相册失败")
			return
		}
		sess.Commit()
	} else {
		var albumAdd = &model.Album{
			Name:     req.Name,
			CreateAt: int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(albumAdd)
		if err != nil {
			common.ResError(c, "添加相册失败")
			return
		}
		var giftAlbumItems []*model.AlbumGift
		for _, i := range req.GiftIds {
			giftAlbumItems = append(giftAlbumItems, &model.AlbumGift{
				AlbumId:  albumAdd.Id,
				GiftId:   i,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(giftAlbumItems)
		if err != nil {
			common.ResError(c, "绑定相册失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func AlbumDelete(c *gin.Context, req model.GiftAlbumDelete) {
	sess := conf.Mysql.NewSession()
	_, err := sess.In("id", req.Ids).Where("delete_at = 0").Update(&model.Album{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除相册失败")
		sess.Rollback()
		return
	}
	_, err = sess.In("album_id", req.Ids).Where("delete_at = 0").Update(&model.AlbumGift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除相册关联关系失败")
		sess.Rollback()
		return
	}
	sess.Commit()
	common.ResOk(c, "ok", nil)
}

/*=====================================app=====================================*/
