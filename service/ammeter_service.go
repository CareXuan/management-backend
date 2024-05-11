package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
	"time"
)

func ListSer(c *gin.Context, page int, pageSize int, num string, userId int) {
	var userRole model.UserRole
	_, err := conf.Mysql.Where("user_id=?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if userRole.RoleId != conf.Conf.Ammeter.Manager && userRole.RoleId != conf.Conf.Ammeter.Supervisor {
		common.ResError(c, "当前用户无法访问此模块")
		return
	}
	var ammeters []*model.Ammeter
	sess := conf.Mysql.NewSession()
	if num != "" {
		sess.Where("num=?", num)
	}
	var ammeterIds []int
	if userRole.RoleId == conf.Conf.Ammeter.Manager {
		var ammeterManages []*model.AmmeterManage
		err := conf.Mysql.Where("user_id=?", userId).Find(&ammeterManages)
		if err != nil {
			common.ResError(c, "获取管理信息失败")
			return
		}
		for _, ammeterManage := range ammeterManages {
			ammeterIds = append(ammeterIds, ammeterManage.AmmeterId)
		}
	}
	if len(ammeterIds) > 0 {
		sess.In("id", ammeterIds)
	}
	count, err := sess.Where("parent_id != 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&ammeters)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: ammeters})
}

func TreeSer(c *gin.Context, userId int) {
	var userRole model.UserRole
	_, err := conf.Mysql.Where("user_id=?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if userRole.RoleId != conf.Conf.Ammeter.Manager && userRole.RoleId != conf.Conf.Ammeter.Supervisor {
		common.ResError(c, "当前用户无法访问此模块")
		return
	}
	var ammeters []*model.Ammeter
	var ammeterIds []int
	if userRole.RoleId == conf.Conf.Ammeter.Manager {
		var ammeterManages []*model.AmmeterManage
		err = conf.Mysql.Where("user_id=?", userId).Find(&ammeterManages)
		if err != nil {
			common.ResError(c, "获取管理信息失败")
			return
		}
		for _, ammeterManage := range ammeterManages {
			ammeterIds = append(ammeterIds, ammeterManage.AmmeterId)
		}
	}
	sess := conf.Mysql.NewSession()
	if len(ammeterIds) > 0 {
		sess.In("id", ammeterIds)
	}
	err = sess.Where("parent_id=0").Find(&ammeters)
	if err != nil {
		common.ResError(c, "获取节点失败")
		return
	}
	getChildNode(ammeters, ammeterIds, userRole.RoleId == conf.Conf.Ammeter.Supervisor)
	common.ResOk(c, "ok", ammeters)
}

func getChildNode(parentNode []*model.Ammeter, manageIds []int, isSupervisor bool) {
	var parentIds []int
	for _, parent := range parentNode {
		parentIds = append(parentIds, parent.Id)
	}
	if len(parentIds) == 0 {
		return
	}
	var childAmmeters []*model.Ammeter
	sess := conf.Mysql.NewSession()
	if len(manageIds) > 0 {
		sess.In("id", manageIds)
	}
	err := sess.In("parent_id", parentIds).Find(&childAmmeters)
	if err != nil {
		log.Fatal(err)
		return
	}
	var childMapping = make(map[int][]*model.Ammeter)
	for _, childNode := range childAmmeters {
		childMapping[childNode.ParentId] = append(childMapping[childNode.ParentId], childNode)
	}
	for _, pNode := range parentNode {
		if isSupervisor {
			pNode.IsSupervisor = utils.AMMETER_SUPERVISOR
		} else {
			pNode.IsSupervisor = utils.AMMETER_MANAGER
		}
		pNode.Children = childMapping[pNode.Id]
		getChildNode(pNode.Children, manageIds, isSupervisor)
	}
}

func AmmeterInfoSer(c *gin.Context, ammeterId int) {
	var ammeterInfo model.Ammeter
	_, err := conf.Mysql.Where("id=?", ammeterId).Get(&ammeterInfo)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	common.ResOk(c, "ok", ammeterInfo)
}

func AmmeterStatisticsSer(c *gin.Context, statisticsType int, ammeterId int, startAt string, endAt string) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", startAt+" 00:00:00")
	endTime, _ := time.Parse("2006-01-02 15:04:05", endAt+" 23:59:59")
	var datas []model.StatisticForm
	sess := conf.Mysql.NewSession()
	sess.Table("ammeter_data")
	sess.Where("ammeter_id=?", ammeterId)
	sess.Where("type=?", statisticsType)
	sess.Where("create_time > ? AND create_time < ?", startTime.Unix(), endTime.Unix())
	if statisticsType == utils.AMMETER_DATA_TYPE_CONSUMPTION {
		sess.Select("DATE_FORMAT(FROM_UNIXTIME(create_time), '%Y-%m-%d %H') AS label,SUM(value) AS count")
		sess.GroupBy("label")
	} else {
		sess.Select("DATE_FORMAT(FROM_UNIXTIME(create_time), '%Y-%m') AS label,SUM(value) AS count")
		sess.GroupBy("label")
	}
	err := sess.Find(&datas)
	if err != nil {
		common.ResError(c, "获取设备统计数据失败")
		return
	}
	var res model.AmmeterStatisticRes
	res.Data = datas
	if statisticsType == utils.AMMETER_DATA_TYPE_CONSUMPTION {
		var todayTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 00:00:00")
		var yesterdayTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(0, 0, -1).Format("2006-01-02")+" 00:00:00")
		var thisMonthTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01")+"-01 00:00:00")
		var lastMonthTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(0, -1, 0).Format("2006-01")+"-01 00:00:00")
		var thisYearTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().Format("2006")+"-01-01 00:00:00")
		var lastYearTime, _ = time.Parse("2006-01-02 15:04:05", time.Now().AddDate(-1, 0, 0).Format("2006")+"-01-01 00:00:00")
		var statisticsData []*model.AmmeterData
		lastYearTimeStamp := time.Date(time.Now().Year()-1, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
		err := conf.Mysql.Where("create_time > ?", lastYearTimeStamp).Where("ammeter_id=?", ammeterId).Where("type=?", statisticsType).Find(&statisticsData)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		for _, sd := range statisticsData {
			cTime := int64(sd.CreateTime)
			if cTime > todayTime.Unix() {
				res.TodayElectricityConsumption += sd.Value
			}
			if cTime < todayTime.Unix() && cTime > yesterdayTime.Unix() {
				res.YesterdayElectricityConsumption += sd.Value
			}
			if cTime > thisMonthTime.Unix() {
				res.MonthElectricityConsumption += sd.Value
			}
			if cTime < thisMonthTime.Unix() && cTime > lastMonthTime.Unix() {
				res.LastMonthElectricityConsumption += sd.Value
			}
			if cTime > thisYearTime.Unix() {
				res.YearElectricityConsumption += sd.Value
			}
			if cTime < thisYearTime.Unix() && cTime > lastYearTime.Unix() {
				res.LastYearElectricityConsumption += sd.Value
			}
		}
	}
	common.ResOk(c, "ok", res)
}

func WarningListSer(c *gin.Context, page int, pageSize int, ammeterId int) {
	var warnings []*model.AmmeterWarning
	sess := conf.Mysql.NewSession()
	sess.Where("ammeter_id=?", ammeterId)
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&warnings)
	if err != nil {
		common.ResError(c, "查询报警列表失败")
		return
	}
	var userIds []int
	for _, w := range warnings {
		userIds = append(userIds, w.DealUser)
	}
	var users []*model.User
	var userMapping = make(map[int]string)
	err = conf.Mysql.In("id", userIds).Find(&users)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	for _, user := range users {
		userMapping[user.Id] = user.Name
	}
	for _, w := range warnings {
		w.DealUserName = userMapping[w.DealUser]
		w.DealTimeStr = time.Unix(int64(w.DealTime), 0).Format("2006-01-02 15:04:05")
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: warnings})
}

func ConfigInfoSer(c *gin.Context, ammeterId int) {
	var config model.AmmeterConfig
	_, err := conf.Mysql.Where("ammeter_id=?", ammeterId).Get(&config)
	if err != nil {
		common.ResError(c, "获取设备配置失败")
		return
	}
	config.TimingOpenTimeStr = time.Unix(int64(config.TimingOpenTime), 0).Format("2006-01-02 15:04:05")
	config.TimingCloseTimeStr = time.Unix(int64(config.TimingCloseTime), 0).Format("2006-01-02 15:04:05")
	common.ResOk(c, "ok", config)
}
