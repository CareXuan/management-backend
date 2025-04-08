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
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&tasks)
	if err != nil {
		common.ResError(c, "获取任务列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: tasks})
}

func TaskInfo(c *gin.Context, id int) {
	var taskItem model.Task
	_, err := conf.Mysql.Where("id = ?", id).Where("delete_at = 0").Get(&taskItem)
	if err != nil {
		common.ResError(c, "获取任务详情失败")
		return
	}
	var taskGifts []*model.TaskGift
	err = conf.Mysql.Where("task_id = ?", id).Where("delete_at = 0").Find(&taskGifts)
	if err != nil {
		common.ResError(c, "获取任务关联礼物失败")
		return
	}
	var giftIds []int
	for _, i := range taskGifts {
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
	for _, i := range taskGifts {
		i.GiftItem = giftMapping[i.TaskId]
	}
	taskItem.Gifts = taskGifts
	common.ResOk(c, "ok", taskItem)
}

func TaskAdd(c *gin.Context, req model.TaskAddReq) {
	if req.Id != 0 {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("id = ?", req.Id).Where("delete_at = 0").Update(&model.Task{
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Deadline:    req.Deadline,
			StartTime:   req.StartTime,
			RepeatCnt:   req.RepeatCnt,
			Status:      1,
			Star:        req.Star,
			Year:        req.Year,
		})
		if err != nil {
			common.ResError(c, "修改任务失败")
			return
		}
		_, err = sess.Where("task_id = ?", req.Id).Update(&model.TaskGift{
			DeleteAt: int(time.Now().Unix()),
		})
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
		_, err = sess.Where("task_id = ?", req.Id).Where("status = ?", 1).Where("delete_at = 0").Update(&model.TaskDo{
			DeleteAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "删除未提交任务失败")
			return
		}
		deadlineInt := 0
		needAdd := true
		taskStartTime := int(time.Now().Unix())
		if req.Type == 2 {
			nextDayTime := common.GetNextDay(req.Deadline)
			needAdd = common.IsToday(nextDayTime)
			deadlineInt = int(nextDayTime.Unix())
		} else {
			deadlineInt = req.Deadline
		}
		if req.Type == 4 {
			taskStartTime = req.StartTime
		}
		if needAdd {
			_, err = conf.Mysql.Insert(model.TaskDo{
				TaskId:    req.Id,
				Status:    1,
				StartTime: taskStartTime,
				Deadline:  deadlineInt,
				CreateAt:  int(time.Now().Unix()),
			})
			if err != nil {
				common.ResError(c, "添加任务失败")
				return
			}
		}
		sess.Commit()
	} else {
		var taskAdd = &model.Task{
			Name:        req.Name,
			Description: req.Description,
			Type:        req.Type,
			Deadline:    req.Deadline,
			StartTime:   req.StartTime,
			RepeatCnt:   req.RepeatCnt,
			Star:        req.Star,
			Status:      1,
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
		deadlineInt := 0
		needAdd := true
		taskStartTime := int(time.Now().Unix())
		if req.Type == 2 {
			nextDayTime := common.GetNextDay(req.Deadline)
			needAdd = common.IsToday(nextDayTime)
			deadlineInt = int(nextDayTime.Unix())
		} else {
			deadlineInt = req.Deadline
		}
		if req.Type == 4 {
			taskStartTime = req.StartTime
		}
		if needAdd {
			_, err = conf.Mysql.Insert(model.TaskDo{
				TaskId:    taskAdd.Id,
				Status:    1,
				StartTime: taskStartTime,
				Deadline:  deadlineInt,
				CreateAt:  int(time.Now().Unix()),
			})
			if err != nil {
				common.ResError(c, "添加任务失败")
				return
			}
		}
	}
	common.ResOk(c, "ok", nil)
}

func TaskDelete(c *gin.Context, req model.TaskDeleteReq) {
	sess := conf.Mysql.NewSession()
	_, err := sess.In("id", req.Ids).Where("delete_at = 0").Update(model.Task{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除任务失败")
		sess.Rollback()
		return
	}
	_, err = sess.In("task_id", req.Ids).Where("delete_at = 0").Update(model.TaskGift{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除任务关联礼物关系失败")
		sess.Rollback()
		return
	}

	_, err = sess.In("task_id", req.Ids).Where("status = ?", 1).Where("delete_at = 0").Update(model.TaskDo{
		DeleteAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "删除任务记录失败")
		sess.Rollback()
		return
	}
	common.ResOk(c, "ok", nil)
}

func TaskChangeStatus(c *gin.Context, req model.TaskChangeStatusReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Where("delete_at = 0").Update(model.Task{
		Status: req.Status,
	})
	if err != nil {
		common.ResError(c, "修改任务状态失败")
		return
	}
	if req.Status == 2 {
		_, err = conf.Mysql.Where("task_id = ?", req.Id).Where("status = ?", 1).Where("delete_at = 0").Update(model.TaskDo{
			DeleteAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "禁用任务关联礼物关系失败")
			return
		}
	} else {
		_, err = conf.Mysql.Where("task_id = ?", req.Id).Where("status = ?", 1).Where("delete_at = 0").Update(model.TaskDo{
			DeleteAt: 0,
		})
		if err != nil {
			common.ResError(c, "启用任务记录失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func TaskCheckList(c *gin.Context, taskId, page, pageSize int) {
	var taskDos []*model.TaskDo
	sess := conf.Mysql.NewSession()
	if taskId != 0 {
		sess.Where("task_id =?", taskId)
	}
	count, err := sess.Where("delete_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&taskDos)
	if err != nil {
		common.ResError(c, "获取任务执行情况列表失败")
		return
	}
	var taskIds []int
	for _, i := range taskDos {
		taskIds = append(taskIds, i.TaskId)
	}
	var taskItems []*model.Task
	err = sess.In("id", taskIds).Find(&taskItems)
	if err != nil {
		common.ResError(c, "获取任务信息失败")
		return
	}
	var taskMapping = make(map[int]*model.Task)
	for _, i := range taskItems {
		taskMapping[i.Id] = i
	}
	for _, i := range taskDos {
		i.Task = taskMapping[i.TaskId]
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: taskDos})
}

func TaskCheckInfo(c *gin.Context, id int) {
	var taskDoItem model.TaskDo
	_, err := conf.Mysql.Where("id = ?", id).Get(&taskDoItem)
	if err != nil {
		common.ResError(c, "获取任务执行情况失败")
		return
	}
	var taskItem *model.Task
	_, err = conf.Mysql.Where("id = ?", taskDoItem.TaskId).Get(taskItem)
	if err != nil {
		common.ResError(c, "获取任务信息失败")
		return
	}
	taskDoItem.Task = taskItem
	common.ResOk(c, "ok", taskDoItem)
}

func TaskCheck(c *gin.Context, req model.TaskCheckReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Where("delete_at = 0").Update(&model.TaskDo{
		Status:  req.Status,
		Reason:  req.Reason,
		CheckAt: int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "审核失败")
		return
	}
	if req.Status == 3 {
		var taskDoItem model.TaskDo
		_, err := conf.Mysql.Where("id = ?", req.Id).Where("delete_at = 0").Get(&taskDoItem)
		if err != nil {
			common.ResError(c, "获取任务信息失败")
			return
		}
		err = sendTaskFinishGift(taskDoItem.TaskId)
		if err != nil {
			common.ResError(c, "发放礼物失败")
			return
		}
		err = checkTaskAchievement(taskDoItem.TaskId)
		if err != nil {
			common.ResError(c, "验证成就失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

/*=====================================app=====================================*/

func TaskDo(c *gin.Context, req model.TaskDoReq) {
	nowTime := int(time.Now().Unix())
	var taskDoItem model.TaskDo
	_, err := conf.Mysql.Where("id = ?", req.Id).Get(&taskDoItem)
	if err != nil {
		common.ResError(c, "搜索可用任务失败")
		return
	}
	if taskDoItem.Deadline != 0 && (nowTime <= taskDoItem.StartTime || nowTime >= taskDoItem.Deadline) {
		common.ResForbidden(c, "当前任务不在可完成时间内")
		return
	}
	if taskDoItem.Status == 2 || taskDoItem.Status == 3 {
		common.ResForbidden(c, "任务已提交审核或已通过，请勿重复提交")
		return
	}
	_, err = conf.Mysql.Where("id = ?", taskDoItem.Id).Update(model.TaskDo{
		Pic:    req.Pic,
		Status: 2,
	})
	if err != nil {
		common.ResError(c, "提交任务失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

/*=====================================func=====================================*/

func sendTaskFinishGift(taskId int) error {
	var taskGifts []*model.TaskGift
	err := conf.Mysql.Where("task_id = ?", taskId).Where("delete_at = 0").Find(&taskGifts)
	if err != nil {
		return err
	}
	for _, i := range taskGifts {
		var giftPackageItem model.GiftPackage
		sess := conf.Mysql.NewSession()
		has, err := sess.Where("gift_id = ?", i.GiftId).Where("delete_at = 0").Get(&giftPackageItem)
		if err != nil {
			sess.Rollback()
			return err
		}
		var giftItem model.Gift
		_, err = sess.Where("id = ?", i.GiftId).Where("delete_at = 0").Get(&giftItem)
		if err != nil {
			sess.Rollback()
			return err
		}
		if has {
			_, err = sess.Where("gift_id = ?", i.GiftId).Where("delete_at = 0").Update(model.GiftPackage{
				Count: giftPackageItem.Count + i.Count,
			})
			if err != nil {
				sess.Rollback()
				return err
			}
		} else {
			_, err = sess.Insert(model.GiftPackage{
				GiftId:     i.GiftId,
				Count:      i.Count,
				Consumable: giftItem.Consumable,
				CreateAt:   int(time.Now().Unix()),
			})
			if err != nil {
				sess.Rollback()
				return err
			}
		}
		_, err = sess.Insert(&model.GiftExtract{
			GiftId:   i.GiftId,
			Type:     1,
			Count:    i.Count,
			GetTime:  int(time.Now().Unix()),
			CreateAt: int(time.Now().Unix()),
		})
		if err != nil {
			sess.Rollback()
			return err
		}
	}

	return nil
}

func checkTaskAchievement(taskId int) error {
	var achievementTasks []*model.AchievementTask
	err := conf.Mysql.Where("task_id = ?", taskId).Where("delete_at = 0").Find(&achievementTasks)
	if err != nil {
		return err
	}
	var achievementIds []int
	for _, i := range achievementTasks {
		achievementIds = append(achievementIds, i.AchievementId)
	}
	var unFinishedAchievement []*model.Achievement
	err = conf.Mysql.In("id", achievementIds).Where("delete_at = 0").Find(&unFinishedAchievement)
	if err != nil {
		return err
	}
	var unFinishedAchievementTask []*model.AchievementTask
	err = conf.Mysql.In("achievement_id", achievementIds).Where("delete_at = 0").Find(&unFinishedAchievementTask)
	if err != nil {
		return err
	}
	var unFinishedAchievementMapping = make(map[int]map[int]int)
	var checkTaskIds []int
	for _, i := range unFinishedAchievementTask {
		checkTaskIds = append(checkTaskIds, i.TaskId)
		if _, ok := unFinishedAchievementMapping[i.AchievementId]; !ok {
			unFinishedAchievementMapping[i.AchievementId] = make(map[int]int)
		}
		unFinishedAchievementMapping[i.AchievementId][i.TaskId] = i.Count
	}
	var checkTaskDo []*model.TaskDo
	err = conf.Mysql.In("task_id", checkTaskIds).Where("status = 3").Where("delete_at = 0").Find(&checkTaskDo)
	if err != nil {
		return err
	}
	var checkTaskCount = make(map[int]int)
	for _, i := range checkTaskDo {
		if _, ok := checkTaskCount[i.TaskId]; ok {
			checkTaskCount[i.TaskId] += 1
		} else {
			checkTaskCount[i.TaskId] = 1
		}
	}
	var needFinishId []int
	for achievementId, i := range unFinishedAchievementMapping {
		needAdd := true
		for taskId, count := range i {
			if checkTaskCount[taskId] != count {
				needAdd = false
			}
		}
		if needAdd {
			needFinishId = append(needFinishId, achievementId)
		}
	}
	_, err = conf.Mysql.In("id", needFinishId).Where("delete_at = 0").Update(&model.Achievement{
		FinishAt: int(time.Now().Unix()),
	})
	for _, i := range needFinishId {
		err = SendFinishAchievement(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendFinishAchievement(achievementId int) error {
	var achievementItem model.Achievement
	_, err := conf.Mysql.Where("id = ?", achievementId).Where("delete_at = 0").Get(&achievementItem)
	if err != nil {
		return err
	}
	var pointGift model.GiftPackage
	has, err := conf.Mysql.Where("gift_id = 0").Where("delete_at = 0").Get(&pointGift)
	if err != nil {
		return err
	}
	if has {
		_, err = conf.Mysql.Where("gift_id = 0").Where("delete_at = 0").Update(&model.GiftPackage{
			Count: pointGift.Count + achievementItem.Point,
		})
	} else {
		_, err = conf.Mysql.Insert(&model.GiftPackage{
			GiftId:   0,
			Count:    achievementItem.Point,
			CreateAt: int(time.Now().Unix()),
		})
	}
	var achievementGifts []*model.AchievementGift
	err = conf.Mysql.Where("achievement_id = ?", achievementId).Find(&achievementGifts)
	if err != nil {
		return err
	}
	for _, j := range achievementGifts {
		var giftPackageItem model.GiftPackage
		sess := conf.Mysql.NewSession()
		has, err := sess.Where("gift_id = ?", j.GiftId).Where("delete_at = 0").Get(&giftPackageItem)
		if err != nil {
			sess.Rollback()
			return err
		}
		var giftItem model.Gift
		_, err = sess.Where("id = ?", j.GiftId).Where("delete_at = 0").Get(&giftItem)
		if err != nil {
			sess.Rollback()
			return err
		}
		if has {
			_, err = sess.Where("gift_id = ?", j.GiftId).Where("delete_at = 0").Update(model.GiftPackage{
				Count: giftPackageItem.Count + j.Count,
			})
			if err != nil {
				sess.Rollback()
				return err
			}
		} else {
			_, err = sess.Insert(model.GiftPackage{
				GiftId:     j.GiftId,
				Count:      j.Count,
				Consumable: giftItem.Consumable,
				CreateAt:   int(time.Now().Unix()),
			})
			if err != nil {
				sess.Rollback()
				return err
			}
		}
		_, err = sess.Insert(&model.GiftExtract{
			GiftId:   j.GiftId,
			Type:     1,
			Count:    j.Count,
			GetTime:  int(time.Now().Unix()),
			CreateAt: int(time.Now().Unix()),
		})
		if err != nil {
			sess.Rollback()
			return err
		}
	}
	return nil
}
