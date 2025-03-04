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
			Status:      1,
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
		_, err = sess.Where("task_id = ?", req.Id).Where("status = ?", 1).Delete(&model.TaskDo{})
		if err != nil {
			common.ResError(c, "删除未提交任务失败")
			return
		}
		deadlineInt := 0
		needAdd := true
		if req.Type == 2 {
			nextDayTime := common.GetNextDay(req.Deadline)
			needAdd = common.IsToday(nextDayTime)
			deadlineInt = int(nextDayTime.Unix())
		} else {
			deadlineInt = req.Deadline
		}
		if needAdd {
			_, err = conf.Mysql.Insert(model.TaskDo{
				TaskId:    req.Id,
				Status:    1,
				StartTime: int(time.Now().Unix()),
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
		deadlineInt := 0
		needAdd := true
		if req.Type == 2 {
			nextDayTime := common.GetNextDay(req.Deadline)
			needAdd = common.IsToday(nextDayTime)
			deadlineInt = int(nextDayTime.Unix())
		} else {
			deadlineInt = req.Deadline
		}
		if needAdd {
			_, err = conf.Mysql.Insert(model.TaskDo{
				TaskId:    taskAdd.Id,
				Status:    1,
				StartTime: int(time.Now().Unix()),
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

func TaskCheck(c *gin.Context, req model.TaskCheckReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Update(&model.TaskDo{
		Status: req.Status,
		Reason: req.Reason,
	})
	if err != nil {
		common.ResError(c, "审核失败")
		return
	}
	if req.Status == 3 {
		var taskDoItem model.TaskDo
		_, err := conf.Mysql.Where("id = ?", req.Id).Get(&taskDoItem)
		if err != nil {
			common.ResError(c, "获取任务信息失败")
			return
		}
		fmt.Println(taskDoItem.TaskId)
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
	err := conf.Mysql.Where("task_id = ?", taskId).Find(&taskGifts)
	if err != nil {
		return err
	}
	for _, i := range taskGifts {
		var giftPackageItem model.GiftPackage
		has, err := conf.Mysql.Where("gift_id = ?", i.GiftId).Get(&giftPackageItem)
		if err != nil {
			return err
		}
		if has {
			_, err = conf.Mysql.Where("gift_id = ?", i.GiftId).Update(model.GiftPackage{
				Count: giftPackageItem.Count + i.Count,
			})
		} else {
			_, err = conf.Mysql.Insert(model.GiftPackage{
				GiftId:   i.GiftId,
				Count:    i.Count,
				CreateAt: int(time.Now().Unix()),
			})
		}
	}

	return nil
}

func checkTaskAchievement(taskId int) error {
	var achievementTasks []*model.AchievementTask
	err := conf.Mysql.Where("task_id = ?", taskId).Find(&achievementTasks)
	if err != nil {
		return err
	}
	var achievementIds []int
	for _, i := range achievementTasks {
		achievementIds = append(achievementIds, i.AchievementId)
	}
	var unFinishedAchievement []*model.Achievement
	err = conf.Mysql.In("id", achievementIds).Find(&unFinishedAchievement)
	if err != nil {
		return err
	}
	var unFinishedAchievementTask []*model.AchievementTask
	err = conf.Mysql.In("achievement_id", achievementIds).Find(&unFinishedAchievementTask)
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
	err = conf.Mysql.In("task_id", checkTaskIds).Where("status = 3").Find(&checkTaskDo)
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
	_, err = conf.Mysql.In("id", needFinishId).Update(&model.Achievement{
		FinishAt: int(time.Now().Unix()),
	})
	for _, i := range needFinishId {
		var achievementItem model.Achievement
		_, err := conf.Mysql.Where("id = ?", i).Get(&achievementItem)
		if err != nil {
			return err
		}
		var pointGift model.GiftPackage
		has, err := conf.Mysql.Where("gift_id = 0").Get(&pointGift)
		if err != nil {
			return err
		}
		if has {
			_, err = conf.Mysql.Where("gift_id = 0").Update(&model.GiftPackage{
				Count: pointGift.Count + achievementItem.Point,
			})
		} else {
			_, err = conf.Mysql.Insert(&model.GiftPackage{
				GiftId:   0,
				Count:    achievementItem.Point,
				CreateAt: int(time.Now().Unix()),
			})
		}
	}
	return nil
}
