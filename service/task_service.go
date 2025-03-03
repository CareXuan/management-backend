package service

import (
	"github.com/gin-gonic/gin"
	"prize-draw/common"
	"prize-draw/conf"
	"prize-draw/model"
	"prize-draw/utils"
	"time"
)

func TaskList(c *gin.Context, name string, page, pageSize int) {
	var tasks []*model.Task
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&tasks)
	if err != nil {
		common.ResError(c, "获取任务列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: tasks})
}

func TaskInfo(c *gin.Context, id int) {
	var taskItem model.Task
	_, err := conf.Mysql.Where("id = ?", id).Get(&taskItem)
	if err != nil {
		common.ResError(c, "获取任务详情失败")
		return
	}
	var taskGifts []*model.TaskGift
	err = conf.Mysql.Where("task_id = ?", id).Find(&taskGifts)
	if err != nil {
		common.ResError(c, "获取任务关联礼物失败")
		return
	}
	var giftIds []int
	for _, i := range taskGifts {
		giftIds = append(giftIds, i.GiftId)
	}
	var giftItems []*model.Gift
	err = conf.Mysql.In("id", giftIds).Find(&giftItems)
	if err != nil {
		common.ResError(c, "获取礼物信息失败")
		return
	}
	var giftMapping = make(map[int]*model.Gift)
	for _, i := range giftItems {
		giftMapping[i.Id] = i
	}
	for _, i := range taskGifts {
		i.GiftItem = giftMapping[i.TaskId]
	}
	taskItem.Gifts = taskGifts
	common.ResOk(c, "ok", taskItem)
}

func TaskAdd(c *gin.Context, req model.TaskAddReq) {
	if req.Id != 0 {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("id = ?", req.Id).Update(&model.Task{
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Deadline:    req.Deadline,
			Star:        req.Star,
			Year:        req.Year,
		})
		if err != nil {
			common.ResError(c, "修改任务失败")
			return
		}
		_, err = sess.Where("task_id = ?", req.Id).Delete(&model.TaskGift{})
		if err != nil {
			common.ResError(c, "删除任务关联礼物失败")
			return
		}
		var taskGiftAdd []*model.TaskGift
		for _, i := range req.BindGifts {
			taskGiftAdd = append(taskGiftAdd, &model.TaskGift{
				TaskId:   req.Id,
				GiftId:   i.GiftId,
				Count:    i.Count,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = sess.Insert(&taskGiftAdd)
		if err != nil {
			common.ResError(c, "添加任务关联礼物失败")
			return
		}
		sess.Commit()
	} else {
		var taskAdd = &model.Task{
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Deadline:    req.Deadline,
			Star:        req.Star,
			Year:        req.Year,
			CreateAt:    int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(taskAdd)
		if err != nil {
			common.ResError(c, "添加任务失败")
			return
		}
		var taskGiftAdd []*model.TaskGift
		for _, i := range req.BindGifts {
			taskGiftAdd = append(taskGiftAdd, &model.TaskGift{
				TaskId:   taskAdd.Id,
				GiftId:   i.GiftId,
				Count:    i.Count,
				CreateAt: int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(&taskGiftAdd)
		if err != nil {
			common.ResError(c, "添加任务关联礼物失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
