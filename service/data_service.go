package service

import (
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	"data_verify/utils"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
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
	sess.Where("status IN (?,?)", model.CHECK_STATUS_WAITING, model.CHECK_STATUS_ERROR)
	_, err = sess.OrderBy("bmddm").Get(&stepCheckData)
	if err != nil {
		common.ResError(c, "获取数据失败")
		return
	}
	var bmdItem model.CheckData
	_, err = conf.Mysql.Select("bmddm,bmdmc").Where("bmddm = ?", stepCheckData.Bmddm).Get(&bmdItem)
	if err != nil {
		common.ResError(c, "获取报名点信息失败")
		return
	}
	common.ResOk(c, "ok", bmdItem)
}

func CheckSer(c *gin.Context, req model.CheckReq) {

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
	if stepCheckData.Step != model.CHECK_STEP_WAITING {
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
	err = common.SetHistory(req.Step, 0, userId, userName, userName+"确认了当前阶段", year)
	if err != nil {
		common.ResError(c, "记录操作失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
