package service

import (
	"github.com/gin-gonic/gin"
	"prize-draw/common"
	"prize-draw/conf"
	"prize-draw/model"
	"prize-draw/utils"
	"time"
)

func AchievementList(c *gin.Context, name string, page, pageSize int) {
	var achievements []*model.Achievement
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&achievements)
	if err != nil {
		common.ResError(c, "获取成就列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: achievements})
}

func AchievementInfo(c *gin.Context, id int) {
	var achievementItem model.Achievement
	_, err := conf.Mysql.Where("id = ?", id).Get(&achievementItem)
	if err != nil {
		common.ResError(c, "获取成就详情失败")
		return
	}
	var achievementTasks []*model.AchievementTask
	err = conf.Mysql.Where("achievement_id = ?", id).Find(&achievementTasks)
	if err != nil {
		common.ResError(c, "获取成就关联任务失败")
		return
	}
	var taskIds []int
	for _, i := range achievementTasks {
		taskIds = append(taskIds, i.TaskId)
	}
	var taskItems []*model.Task
	err = conf.Mysql.In("id", taskIds).Find(&taskItems)
	if err != nil {
		common.ResError(c, "获取任务信息失败")
		return
	}
	var taskMapping = make(map[int]*model.Task)
	for _, i := range taskItems {
		taskMapping[i.Id] = i
	}
	var achievementGifts []*model.AchievementGift
	err = conf.Mysql.Where("achievement_id = ?", id).Find(&achievementGifts)
	if err != nil {
		common.ResError(c, "获取成就关联礼物失败")
		return
	}
	var giftIds []int
	for _, i := range achievementGifts {
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
	for _, i := range achievementTasks {
		i.TaskItem = taskMapping[i.TaskId]
	}
	for _, i := range achievementGifts {
		i.GiftItem = giftMapping[i.GiftId]
	}
	achievementItem.Tasks = achievementTasks
	achievementItem.Gifts = achievementGifts
	common.ResOk(c, "ok", achievementItem)
}

func AchievementAdd(c *gin.Context, req model.AchievementAddReq) {
	if req.Id != 0 {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("id = ?", req.Id).Update(&model.Achievement{
			Name:        req.Name,
			Pic:         req.Pic,
			Description: req.Description,
			Point:       req.Point,
		})
		if err != nil {
			common.ResError(c, "修改成就失败")
			return
		}
		_, err = sess.Where("achievement_id = ?", req.Id).Delete(&model.AchievementTask{})
		if err != nil {
			common.ResError(c, "删除成就关联任务失败")
			return
		}
		_, err = sess.Where("achievement_id = ?", req.Id).Delete(&model.AchievementGift{})
		if err != nil {
			common.ResError(c, "删除成就关联礼物失败")
			return
		}
		var achievementTask []*model.AchievementTask
		for _, i := range req.Tasks {
			achievementTask = append(achievementTask, &model.AchievementTask{
				AchievementId: req.Id,
				TaskId:        i.TaskId,
				Count:         i.Count,
				CreateAt:      int(time.Now().Unix()),
			})
		}
		_, err = sess.Insert(&achievementTask)
		if err != nil {
			common.ResError(c, "添加成就关联任务失败")
			return
		}
		var achievementGift []*model.AchievementGift
		for _, i := range req.Gifts {
			achievementGift = append(achievementGift, &model.AchievementGift{
				AchievementId: req.Id,
				GiftId:        i.GiftId,
				Count:         i.Count,
				CreateAt:      int(time.Now().Unix()),
			})
		}
		_, err = sess.Insert(&achievementGift)
		if err != nil {
			common.ResError(c, "添加成就关联礼物失败")
			return
		}
		sess.Commit()
	} else {
		var achievementAdd = &model.Achievement{
			Name:        req.Name,
			Pic:         req.Pic,
			Description: req.Description,
			Point:       req.Point,
			CreateAt:    int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(achievementAdd)
		if err != nil {
			common.ResError(c, "添加成就失败")
			return
		}
		var achievementTask []*model.AchievementTask
		for _, i := range req.Tasks {
			achievementTask = append(achievementTask, &model.AchievementTask{
				AchievementId: achievementAdd.Id,
				TaskId:        i.TaskId,
				Count:         i.Count,
				CreateAt:      int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(&achievementTask)
		if err != nil {
			common.ResError(c, "添加成就关联任务失败")
			return
		}
		var achievementGift []*model.AchievementGift
		for _, i := range req.Gifts {
			achievementGift = append(achievementGift, &model.AchievementGift{
				AchievementId: achievementAdd.Id,
				GiftId:        i.GiftId,
				Count:         i.Count,
				CreateAt:      int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(&achievementGift)
		if err != nil {
			common.ResError(c, "添加成就关联礼物失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}
