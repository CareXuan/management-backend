package service

import (
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	"data_verify/utils"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"time"
)

func StepListSer(c *gin.Context, bmddm, bmdmc string, page, pageSize int) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	sess := conf.Mysql.NewSession()
	sess.Where("year = ?", year)
	if bmddm != "" {
		sess.Where("bmddm = ?", bmddm)
	}
	if bmdmc != "" {
		sess.Where("bmdmc = ?", bmdmc)
	}
	var checkDatas []*model.CheckData
	count, err := sess.OrderBy("bmddm").Limit(pageSize, (page-1)*pageSize).FindAndCount(&checkDatas)
	if err != nil {
		common.ResError(c, "获取列表失败")
		return
	}
	var checkDataRes = make(map[string]*model.CheckListRes)
	var bmddms []string
	for _, i := range checkDatas {
		bmddms = append(bmddms, i.Bmddm)
		checkDataRes[i.Bmddm] = &model.CheckListRes{
			Bmddm: i.Bmddm,
			Bmdmc: i.Bmdmc,
			Step1: model.STEP_STATUS_WAITING,
			Step2: model.STEP_STATUS_WAITING,
			Step3: model.STEP_STATUS_WAITING,
			Step4: model.STEP_STATUS_WAITING,
			Step5: model.STEP_STATUS_WAITING,
		}
	}
	var stepCheckDatas []*model.StepCheckData
	err = conf.Mysql.In("bmddm", bmddms).Find(&stepCheckDatas)
	if err != nil {
		common.ResError(c, "获取校验结果失败")
		return
	}

	for _, i := range stepCheckDatas {
		if i.Status == model.CHECK_STATUS_PASS {
			switch i.Step {
			case model.CHECK_STEP_WAITING:
				checkDataRes[i.Bmddm].Step1 = model.STEP_STATUS_PASS
			case model.CHECK_STEP_FIRST:
				checkDataRes[i.Bmddm].Step2 = model.STEP_STATUS_PASS
			case model.CHECK_STEP_SECOND:
				checkDataRes[i.Bmddm].Step3 = model.STEP_STATUS_PASS
			case model.CHECK_STEP_THIRD:
				checkDataRes[i.Bmddm].Step4 = model.STEP_STATUS_PASS
			case model.CHECK_STEP_FOURTH:
				checkDataRes[i.Bmddm].Step5 = model.STEP_STATUS_PASS
			}
		}
	}
	sort.Strings(bmddms)
	var res []*model.CheckListRes
	for _, i := range bmddms {
		res = append(res, checkDataRes[i])
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: res})
}

func StepInfoSer(c *gin.Context, step int) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	userId, _, err := common.GetUserIdByToken(c)
	if err != nil {
		common.ResError(c, "获取用户ID失败")
		return
	}
	var groupUser model.GroupUser
	_, err = conf.Mysql.Where("user_id = ?", userId).Get(&groupUser)
	if err != nil {
		common.ResError(c, "获取组配置失败")
		return
	}
	var stepCheckData model.StepCheckData
	sess := conf.Mysql.NewSession()
	sess.Where("year = ?", year)
	sess.Where("bmddm >= ?", groupUser.BmdStart)
	sess.Where("bmddm <= ?", groupUser.BmdEnd)
	sess.Where("step = ?", step)
	sess.Where("status != ?", model.CHECK_STATUS_PASS)
	_, err = sess.OrderBy("bmddm").Get(&stepCheckData)
	if err != nil {
		common.ResError(c, "获取数据失败")
		return
	}
	selectString := "bmddm,bmdmc"
	if step == model.CHECK_STEP_FOURTH {
		selectString = "bmddm,bmdmc,envelope_cnt,first_cnt,last_cnt,blue_cnt,red_cnt,black_cnt"
	}
	var bmdItem model.CheckData
	_, err = conf.Mysql.Select(selectString).Where("bmddm = ?", stepCheckData.Bmddm).Where("year = ?", year).Get(&bmdItem)
	if err != nil {
		common.ResError(c, "获取报名点信息失败")
		return
	}
	var bmdCountItems []*model.CheckData
	err = conf.Mysql.Where("year = ?", year).Where("bmddm >= ?", groupUser.BmdStart).Where("bmddm <= ?", groupUser.BmdEnd).Find(&bmdCountItems)
	if err != nil {
		common.ResError(c, "获取报名点数量失败")
		return
	}
	var finishBmdCountItems []*model.StepCheckData
	err = conf.Mysql.Where("year = ?", year).Where("bmddm >= ?", groupUser.BmdStart).Where("bmddm <= ?", groupUser.BmdEnd).Where("step = ?", step).Where("status = ?", model.CHECK_STATUS_PASS).Find(&finishBmdCountItems)
	if err != nil {
		common.ResError(c, "获取报名点完成数量失败")
		return
	}
	bmdItem.ProgressCnt = len(bmdCountItems)
	bmdItem.RemainCnt = len(bmdCountItems) - len(finishBmdCountItems)
	common.ResOk(c, "ok", bmdItem)
}

func CheckSer(c *gin.Context, req model.CheckReq) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	userId, userName, err := common.GetUserIdByToken(c)
	if err != nil {
		common.ResError(c, "获取用户ID失败")
		return
	}
	var groupUser model.GroupUser
	_, err = conf.Mysql.Where("user_id = ?", userId).Get(&groupUser)
	if err != nil {
		common.ResError(c, "获取组配置失败")
		return
	}
	bmdStart, _ := strconv.Atoi(groupUser.BmdStart)
	bmdEnd, _ := strconv.Atoi(groupUser.BmdEnd)
	bmdNow, _ := strconv.Atoi(req.Bmddm)
	if bmdNow < bmdStart || bmdNow > bmdEnd {
		common.ResForbidden(c, "您无权处理此报名点")
		return
	}
	var existCheckItem model.StepCheckData
	_, err = conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Get(&existCheckItem)
	if err != nil {
		common.ResError(c, "获取当前阶段配置失败")
		return
	}
	if existCheckItem.Id == 0 {
		common.ResForbidden(c, "当前报名点尚未到达此阶段")
		return
	}
	if existCheckItem.Status == model.CHECK_STATUS_PASS {
		common.ResForbidden(c, "当前报名点此阶段已校验通过")
		return
	}
	var checkDataItem model.CheckData
	_, err = conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("year = ?", year).Get(&checkDataItem)
	if err != nil {
		common.ResError(c, "获取校验数据失败")
		return
	}
	switch req.Step {
	case model.CHECK_STEP_FIRST:
		allCnt := req.Data[0]
		if checkDataItem.EnvelopeCnt != allCnt {
			common.SetHistory(req.Step, model.CHECK_STATUS_ERROR, userId, req.Bmddm, userName, userName+"检查此阶段失败", year)
			_, err := conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Update(model.StepCheckData{
				Status: model.CHECK_STATUS_ERROR,
			})
			if err != nil {
				common.ResError(c, "修改检查状态失败")
				return
			}
			common.ResForbidden(c, "校验失败")
			return
		}
	case model.CHECK_STEP_SECOND:
		startCnt := req.Data[0]
		endCnt := req.Data[1]
		if checkDataItem.FirstCnt != startCnt || checkDataItem.LastCnt != endCnt {
			common.SetHistory(req.Step, model.CHECK_STATUS_ERROR, userId, req.Bmddm, userName, userName+"检查此阶段失败", year)
			_, err := conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Update(model.StepCheckData{
				Status: model.CHECK_STATUS_ERROR,
			})
			if err != nil {
				common.ResError(c, "修改检查状态失败")
				return
			}
			common.ResForbidden(c, "校验失败")
			return
		}
	case model.CHECK_STEP_THIRD:
		blueCnt := req.Data[0]
		redCnt := req.Data[1]
		blackCnt := req.Data[2]
		if checkDataItem.BlueCnt != blueCnt || checkDataItem.RedCnt != redCnt || checkDataItem.BlackCnt != blackCnt {
			common.SetHistory(req.Step, model.CHECK_STATUS_ERROR, userId, req.Bmddm, userName, userName+"检查此阶段失败", year)
			errStatus := model.CHECK_STATUS_ERROR
			if checkDataItem.BlueCnt != blueCnt && checkDataItem.RedCnt == redCnt && checkDataItem.BlackCnt == blackCnt {
				errStatus = model.CHECK_STATUS_BLUE_ERROR
			}
			if checkDataItem.BlueCnt == blueCnt && checkDataItem.RedCnt != redCnt && checkDataItem.BlackCnt == blackCnt {
				errStatus = model.CHECK_STATUS_RED_ERROR
			}
			if checkDataItem.BlueCnt == blueCnt && checkDataItem.RedCnt == redCnt && checkDataItem.BlackCnt != blackCnt {
				errStatus = model.CHECK_STATUS_BLACK_ERROR
			}
			if checkDataItem.BlueCnt != blueCnt && checkDataItem.RedCnt != redCnt && checkDataItem.BlackCnt == blackCnt {
				errStatus = model.CHECK_STATUS_RED_BLUE_ERROR
			}
			if checkDataItem.BlueCnt == blueCnt && checkDataItem.RedCnt != redCnt && checkDataItem.BlackCnt != blackCnt {
				errStatus = model.CHECK_STATUS_BLACK_RED_ERROR
			}
			if checkDataItem.BlueCnt != blueCnt && checkDataItem.RedCnt == redCnt && checkDataItem.BlackCnt != blackCnt {
				errStatus = model.CHECK_STATUS_BLACK_BLUE_ERROR
			}
			if checkDataItem.BlueCnt != blueCnt && checkDataItem.RedCnt != redCnt && checkDataItem.BlackCnt != blackCnt {
				errStatus = model.CHECK_STATUS_ALL_ERROR
			}
			_, err := conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Update(model.StepCheckData{
				Status: errStatus,
			})
			if err != nil {
				common.ResError(c, "修改检查状态失败")
				return
			}
			common.ResForbidden(c, "校验失败")
		}
	}
	common.SetHistory(req.Step, model.CHECK_STATUS_PASS, userId, req.Bmddm, userName, userName+"检查此阶段通过", year)
	_, err = conf.Mysql.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Update(model.StepCheckData{
		Status: model.CHECK_STATUS_PASS,
	})
	if err != nil {
		common.ResError(c, "修改检查状态失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func NextSer(c *gin.Context, req model.NextStepReq) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	userId, userName, err := common.GetUserIdByToken(c)
	if err != nil {
		common.ResError(c, "获取用户ID失败")
		return
	}
	var groupUser model.GroupUser
	_, err = conf.Mysql.Where("user_id = ?", userId).Get(&groupUser)
	if err != nil {
		common.ResError(c, "获取组配置失败")
		return
	}
	bmdStart, _ := strconv.Atoi(groupUser.BmdStart)
	bmdEnd, _ := strconv.Atoi(groupUser.BmdEnd)
	bmdNow, _ := strconv.Atoi(req.Bmddm)
	if bmdNow < bmdStart || bmdNow > bmdEnd {
		common.ResForbidden(c, "您无权处理此报名点")
		return
	}
	var stepCheckData model.StepCheckData
	_, err = conf.Mysql.Where("year = ?", year).Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Get(&stepCheckData)
	if err != nil {
		common.ResError(c, "获取当前阶段信息失败")
		return
	}
	if stepCheckData.Id == 0 {
		common.ResForbidden(c, "您尚未进入此阶段")
		return
	}
	if stepCheckData.Step != model.CHECK_STEP_WAITING && stepCheckData.Step != model.CHECK_STEP_FOURTH {
		if stepCheckData.Status != model.CHECK_STATUS_PASS {
			common.ResForbidden(c, "您的当前流程尚未校验通过")
			return
		}
	}
	sess := conf.Mysql.NewSession()
	_, err = sess.Where("bmddm = ?", req.Bmddm).Where("step = ?", req.Step).Where("year = ?", year).Update(&model.StepCheckData{
		Status: model.CHECK_STATUS_PASS,
	})
	if err != nil {
		sess.Rollback()
		common.ResError(c, "通过失败")
		return
	}
	_, err = sess.Insert(model.StepCheckData{
		Bmddm:  req.Bmddm,
		Step:   req.Step + 1,
		Status: model.CHECK_STATUS_WAITING,
		Year:   year,
	})
	if err != nil {
		sess.Rollback()
		common.ResError(c, "写入新流程失败")
		return
	}
	_, err = sess.Where("bmddm = ?", req.Bmddm).Update(model.CheckData{
		Step: req.Step + 1,
	})
	if err != nil {
		sess.Rollback()
		common.ResError(c, "写入新流程失败")
		return
	}
	sess.Commit()
	err = common.SetHistory(req.Step, 0, userId, req.Bmddm, userName, userName+"确认了当前阶段", year)
	if err != nil {
		common.ResError(c, "记录操作失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func HistoryListSer(c *gin.Context, page, pageSize int) {
	year, err := common.GetCurrentYear()
	if err != nil {
		common.ResError(c, "获取年份信息失败")
		return
	}
	var historyItems []*model.History
	sess := conf.Mysql.NewSession()
	sess.Where("year = ?", year)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).OrderBy("create_time DESC").FindAndCount(&historyItems)
	if err != nil {
		common.ResError(c, "获取操作历史失败")
		return
	}
	var bmddms []string
	for _, i := range historyItems {
		bmddms = append(bmddms, i.Bmddm)
	}
	var bmdItems []*model.CheckData
	err = conf.Mysql.In("bmddm", bmddms).Where("year = ?", year).Find(&bmdItems)
	if err != nil {
		common.ResError(c, "获取报名点信息失败")
		return
	}
	var bmdMapping = make(map[string]string)
	for _, i := range bmdItems {
		bmdMapping[i.Bmddm] = i.Bmdmc
	}
	for _, i := range historyItems {
		i.Bmdmc = bmdMapping[i.Bmddm]
		i.CreateTimeStr = time.Unix(int64(i.CreateTime), 0).Format("2006-01-02 15:04:05")
		i.StepName = model.CHECK_STEP_NAME_MAPPING[i.Step]
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: historyItems})
}
