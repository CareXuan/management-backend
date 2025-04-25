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
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&achievements)
	if err != nil {
		common.ResError(c, "获取成就列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: achievements})
}

func AchievementInfo(c *gin.Context, id int) {
	var achievementItem model.Achievement
	_, err := conf.Mysql.Where("id = ?", id).Where("delete_at = 0").Get(&achievementItem)
	if err != nil {
		common.ResError(c, "获取成就详情失败")
		return
	}
	var achievementTasks []*model.AchievementTask
	err = conf.Mysql.Where("achievement_id = ?", id).Where("delete_at = 0").Find(&achievementTasks)
	if err != nil {
		common.ResError(c, "获取成就关联任务失败")
		return
	}
	var taskIds []int
	for _, i := range achievementTasks {
		taskIds = append(taskIds, i.TaskId)
	}
	var taskItems []*model.Task
	err = conf.Mysql.In("id", taskIds).Where("delete_at = 0").Find(&taskItems)
	if err != nil {
		common.ResError(c, "获取任务信息失败")
		return
	}
	var taskMapping = make(map[int]*model.Task)
	for _, i := range taskItems {
		taskMapping[i.Id] = i
	}
	var achievementGifts []*model.AchievementGift
	err = conf.Mysql.Where("achievement_id = ?", id).Where("delete_at = 0").Find(&achievementGifts)
	if err != nil {
		common.ResError(c, "获取成就关联礼物失败")
		return
	}
	var giftIds []int
	for _, i := range achievementGifts {
		giftIds = append(giftIds, i.GiftId)
	}
	var giftItems []*model.Gift
	err = conf.Mysql.In("id", giftIds).Where("delete_at = 0").Find(&giftItems)
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
		_, err := sess.Where("id = ?", req.Id).Where("delete_at = 0").Update(&model.Achievement{
			Name:        req.Name,
			Pic:         req.Pic,
			Description: req.Description,
			Point:       req.Point,
		})
		if err != nil {
			common.ResError(c, "修改成就失败")
			return
		}
		_, err = sess.Where("achievement_id = ?", req.Id).Where("delete_at = 0").Update(&model.AchievementTask{
			DeleteAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "删除成就关联任务失败")
			return
		}
		_, err = sess.Where("achievement_id = ?", req.Id).Where("delete_at = 0").Update(&model.AchievementGift{
			DeleteAt: int(time.Now().Unix()),
		})
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

func AchievementFinish(c *gin.Context, req model.AchievementFinishReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Update(&model.Achievement{
		IsFinish: 1,
		FinishAt: req.FinishTime,
	})
	if err != nil {
		common.ResError(c, "修改成就失败")
		return
	}
	//err = SendFinishAchievement(req.Id)
	//if err != nil {
	//	common.ResError(c, "发放成就奖励失败")
	//	return
	//}
	common.ResOk(c, "ok", nil)
}

func AchievementDelete(c *gin.Context, req model.AchievementDeleteReq) {
	sess := conf.Mysql.NewSession()
	_, err := sess.In("id", req.Ids).Where("delete_at = 0").Update(model.Achievement{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除成就失败")
		sess.Rollback()
		return
	}
	_, err = sess.In("achievement_id", req.Ids).Where("delete_at = 0").Update(model.AchievementGift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除成就关联礼物失败")
		sess.Rollback()
		return
	}
	_, err = sess.In("achievement_id", req.Ids).Where("delete_at = 0").Update(model.AchievementTask{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除成就关联任务失败")
		sess.Rollback()
		return
	}
	common.ResOk(c, "ok", nil)
}

/*=====================================app=====================================*/

func AppAchievementList(c *gin.Context, statusInt int) {
	var achievementList []*model.Achievement
	sess := conf.Mysql.NewSession()
	if statusInt == 1 {
		sess.Where("is_finish = 1")
	}
	if statusInt == 2 {
		sess.Where("is_finish = 0")
	}
	err := sess.Where("delete_at = 0").Find(&achievementList)
	if err != nil {
		common.ResError(c, "获取成就列表失败")
		return
	}
	var unFinishAchievementId []int
	for _, i := range achievementList {
		if i.IsFinish == 0 {
			unFinishAchievementId = append(unFinishAchievementId, i.Id)
		}
	}
	var achievementTask []*model.AchievementTask
	err = conf.Mysql.In("achievement_id", unFinishAchievementId).Find(&achievementTask)
	if err != nil {
		common.ResError(c, "获取关联任务失败")
		return
	}
	var taskIds []int
	for _, i := range achievementTask {
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
	var taskProgress []*model.TaskDo
	err = conf.Mysql.In("task_id", taskIds).Find(&taskProgress)
	if err != nil {
		common.ResError(c, "搜索任务进度失败")
		return
	}
	var taskProgressMapping = make(map[int]int)
	for _, i := range taskProgress {
		if _, ok := taskProgressMapping[i.TaskId]; ok && i.Status == 3 {
			taskProgressMapping[i.TaskId] += 1
		} else {
			taskProgressMapping[i.TaskId] = 1
		}
	}
	var achievementTaskMapping = make(map[int][]*model.AchievementTask)
	for _, i := range achievementTask {
		i.Progress = taskProgressMapping[i.TaskId]
		i.TaskItem = taskMapping[i.TaskId]
		if _, ok := achievementTaskMapping[i.AchievementId]; ok {
			achievementTaskMapping[i.AchievementId] = append(achievementTaskMapping[i.AchievementId], i)
		} else {
			achievementTaskMapping[i.AchievementId] = []*model.AchievementTask{i}
		}
	}

	for _, i := range achievementList {
		i.Tasks = achievementTaskMapping[i.Id]
	}
	common.ResOk(c, "ok", achievementList)
}

func AppAchievementReceive(c *gin.Context, req model.AppAchievementReceiveReq) {
	achievementId := req.AchievementId
	var achievementItem model.Achievement
	_, err := conf.Mysql.Where("id = ?", achievementId).Get(&achievementItem)
	if err != nil {
		common.ResError(c, "获取成就信息失败")
		return
	}
	if achievementItem.IsFinish == 0 {
		common.ResForbidden(c, "成就尚未完成")
		return
	}
	if achievementItem.IsReceive == 1 {
		common.ResForbidden(c, "奖励已经领取")
		return
	}
	_, err = conf.Mysql.Where("id = ?", achievementId).Update(&model.Achievement{
		IsReceive: 1,
	})
	err = SendFinishAchievement(achievementId)
	if err != nil {
		common.ResError(c, "发放成就奖励失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
