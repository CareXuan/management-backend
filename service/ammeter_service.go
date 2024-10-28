package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model/ammeter"
	"management-backend/model/rbac"
	"management-backend/utils"
	"strconv"
	"time"
)

func ListSer(c *gin.Context, page int, pageSize int, card string, status int, ammeterType int, userId int) {
	var userRole rbac.UserRole
	_, err := conf.Mysql.Where("user_id=?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if userRole.RoleId != conf.Conf.Ammeter.Manager && userRole.RoleId != conf.Conf.Ammeter.Supervisor && userRole.RoleId != conf.Conf.Admin {
		common.ResError(c, "当前用户无法访问此模块")
		return
	}
	var ammeters []*ammeter.Ammeter
	sess := conf.Mysql.NewSession()
	if card != "" {
		sess.Where("card=?", card)
	}
	if status != 0 {
		sess.Where("status=?", status)
	}
	if ammeterType != 0 {
		sess.Where("type=?", ammeterType)
	}
	var ammeterIds []int
	if userRole.RoleId == conf.Conf.Ammeter.Manager {
		var ammeterManages []*ammeter.AmmeterManage
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
	var userRole rbac.UserRole
	_, err := conf.Mysql.Where("user_id=?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	if userRole.RoleId != conf.Conf.Ammeter.Manager && userRole.RoleId != conf.Conf.Ammeter.Supervisor && userRole.RoleId != conf.Conf.Admin {
		common.ResError(c, "当前用户无法访问此模块")
		return
	}
	var ammeters []*ammeter.Ammeter
	var ammeterIds []int
	if userRole.RoleId == conf.Conf.Ammeter.Manager {
		var ammeterManages []*ammeter.AmmeterManage
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

func TreeManagerSer(c *gin.Context) {
	var userRoles []*rbac.UserRole
	var userIds []int
	err := conf.Mysql.Where("role_id=?", conf.Conf.Ammeter.Manager).Find(&userRoles)
	if err != nil {
		common.ResError(c, "获取管理员ID失败")
		return
	}
	for _, userRole := range userRoles {
		userIds = append(userIds, userRole.UserId)
	}
	var users []*rbac.User
	err = conf.Mysql.In("id", userIds).Find(&users)
	if err != nil {
		common.ResError(c, "获取用户失败")
	}
	var res []ammeter.TreeManagerRes
	for _, user := range users {
		res = append(res, ammeter.TreeManagerRes{
			Id:   user.Id,
			Name: user.Name,
		})
	}
	common.ResOk(c, "ok", res)
}

func getChildNode(parentNode []*ammeter.Ammeter, manageIds []int, isSupervisor bool) {
	var parentIds []int
	for _, parent := range parentNode {
		parentIds = append(parentIds, parent.Id)
	}
	if len(parentIds) == 0 {
		return
	}
	var childAmmeters []*ammeter.Ammeter
	sess := conf.Mysql.NewSession()
	if len(manageIds) > 0 {
		sess.In("id", manageIds)
	}
	err := sess.In("parent_id", parentIds).Find(&childAmmeters)
	if err != nil {
		log.Fatal(err)
		return
	}
	var childMapping = make(map[int][]*ammeter.Ammeter)
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

func AddTreeNodeSer(c *gin.Context, nodeId int, nodeType int, nodeModel string, num string, card string, location string, parentId int, managers []int) {
	if nodeId == 0 {
		var insertAmmeter = ammeter.Ammeter{
			Type:       nodeType,
			Model:      nodeModel,
			Num:        num,
			Card:       card,
			Location:   location,
			ParentId:   parentId,
			Status:     utils.AMMETER_STATUS_OFFLINE,
			Switch:     utils.AMMETER_STATUS_SWITCH_CLOSE,
			CreateTime: int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(&insertAmmeter)
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
		_, err = conf.Mysql.Insert(&ammeter.AmmeterConfig{
			AmmeterId: insertAmmeter.Id,
		})
		if err != nil {
			common.ResError(c, "添加设备配置失败")
			return
		}
		var addAmmeterManager []*ammeter.AmmeterManage
		var addAmmeterManagerPwd []*ammeter.AmmeterManageConfig
		for _, manager := range managers {
			addAmmeterManager = append(addAmmeterManager, &ammeter.AmmeterManage{
				AmmeterId:  insertAmmeter.Id,
				UserId:     manager,
				CreateTime: int(time.Now().Unix()),
			})
			addAmmeterManagerPwd = append(addAmmeterManagerPwd, &ammeter.AmmeterManageConfig{
				AmmeterId:  insertAmmeter.Id,
				UserId:     manager,
				Type:       utils.AMMETER_MANAGE_CONFIG_TYPE_PWD,
				Value:      utils.AMMETER_MANAGE_CONFIG_DEFAULT_PWD,
				CreateTime: int(time.Now().Unix()),
			})
			if parentId != 0 {
				addAmmeterManager = append(addAmmeterManager, &ammeter.AmmeterManage{
					AmmeterId:  parentId,
					UserId:     manager,
					CreateTime: int(time.Now().Unix()),
				})
				addAmmeterManagerPwd = append(addAmmeterManagerPwd, &ammeter.AmmeterManageConfig{
					AmmeterId:  parentId,
					UserId:     manager,
					Type:       utils.AMMETER_MANAGE_CONFIG_TYPE_PWD,
					Value:      utils.AMMETER_MANAGE_CONFIG_DEFAULT_PWD,
					CreateTime: int(time.Now().Unix()),
				})
			}
		}
		_, err = conf.Mysql.Insert(&addAmmeterManager)
		if err != nil {
			common.ResError(c, "添加管理员失败")
			return
		}
		_, err = conf.Mysql.Insert(&addAmmeterManagerPwd)
		if err != nil {
			common.ResError(c, "添加管理员密码失败")
			return
		}
	} else {
		_, err := conf.Mysql.Where("id=?", nodeId).Update(ammeter.Ammeter{
			Type:       nodeType,
			Model:      nodeModel,
			Num:        num,
			Card:       card,
			Location:   location,
			ParentId:   parentId,
			Status:     utils.AMMETER_STATUS_OFFLINE,
			Switch:     utils.AMMETER_STATUS_SWITCH_CLOSE,
			CreateTime: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
		_, err = conf.Mysql.Where("ammeter_id=?", nodeId).Delete(&ammeter.AmmeterManage{})
		if err != nil {
			common.ResError(c, "删除管理员失败")
			return
		}
		var addAmmeterManager []*ammeter.AmmeterManage
		var addAmmeterManagerPwd []*ammeter.AmmeterManageConfig
		for _, manager := range managers {
			addAmmeterManager = append(addAmmeterManager, &ammeter.AmmeterManage{
				AmmeterId:  nodeId,
				UserId:     manager,
				CreateTime: int(time.Now().Unix()),
			})
			addAmmeterManagerPwd = append(addAmmeterManagerPwd, &ammeter.AmmeterManageConfig{
				AmmeterId:  nodeId,
				UserId:     manager,
				Type:       utils.AMMETER_MANAGE_CONFIG_TYPE_PWD,
				Value:      utils.AMMETER_MANAGE_CONFIG_DEFAULT_PWD,
				CreateTime: int(time.Now().Unix()),
			})
			if parentId != 0 {
				addAmmeterManager = append(addAmmeterManager, &ammeter.AmmeterManage{
					AmmeterId:  parentId,
					UserId:     manager,
					CreateTime: int(time.Now().Unix()),
				})
				addAmmeterManagerPwd = append(addAmmeterManagerPwd, &ammeter.AmmeterManageConfig{
					AmmeterId:  parentId,
					UserId:     manager,
					Type:       utils.AMMETER_MANAGE_CONFIG_TYPE_PWD,
					Value:      utils.AMMETER_MANAGE_CONFIG_DEFAULT_PWD,
					CreateTime: int(time.Now().Unix()),
				})
			}
		}
		_, err = conf.Mysql.Insert(&addAmmeterManager)
		if err != nil {
			common.ResError(c, "添加管理员失败")
			return
		}
		_, err = conf.Mysql.Insert(&addAmmeterManagerPwd)
		if err != nil {
			common.ResError(c, "添加管理员密码失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

func DeleteTreeNodeSer(c *gin.Context, nodeId int) {
	var nodes []*ammeter.Ammeter
	var nodeIds []int
	err := conf.Mysql.Where("id=?", nodeId).Or("parent_id=?", nodeId).Find(&nodes)
	if err != nil {
		common.ResError(c, "获取节点失败")
		return
	}
	for _, node := range nodes {
		nodeIds = append(nodeIds, node.Id)
	}
	_, err = conf.Mysql.In("id", nodeIds).Delete(&ammeter.Ammeter{})
	if err != nil {
		common.ResError(c, "删除节点失败")
		return
	}
	_, err = conf.Mysql.In("ammeter_id", nodeIds).Delete(&ammeter.AmmeterManage{})
	if err != nil {
		common.ResError(c, "删除节点管理失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func AmmeterInfoSer(c *gin.Context, ammeterId int) {
	var ammeterInfo ammeter.Ammeter
	_, err := conf.Mysql.Where("id=?", ammeterId).Get(&ammeterInfo)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	var boardInfo ammeter.Board
	_, err = conf.Mysql.Where("board_id = ?", ammeterInfo.Card).Get(&boardInfo)
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "获取主板信息失败")
		return
	}
	ammeterInfo.Iccid = boardInfo.Iccid
	num := ""
	if len(ammeterInfo.Num) == 1 {
		num = "0" + ammeterInfo.Num
	} else {
		num = ammeterInfo.Num
	}
	ammeterInfo.Num = num
	common.ResOk(c, "ok", ammeterInfo)
}

func ChangeAmmeterSwitchSer(c *gin.Context, ammeterId int, userId int, ammeterSwitch int, pwd string) {
	var ammeterItem ammeter.Ammeter
	_, err := conf.Mysql.Where("id=?", ammeterId).Get(&ammeterItem)
	if err != nil {
		common.ResError(c, "获取电表信息失败")
		return
	}
	var ammeterManagePwd ammeter.AmmeterManageConfig
	_, err = conf.Mysql.Where("ammeter_id=?", ammeterId).Where("user_id=?", userId).Where("type=?", utils.AMMETER_MANAGE_CONFIG_TYPE_PWD).Get(&ammeterManagePwd)
	if err != nil {
		common.ResError(c, "获取电表配置信息失败")
		return
	}
	if ammeterManagePwd.Value != pwd {
		common.ResForbidden(c, "密码错误，请重试！")
		return
	}
	_, err = conf.Mysql.Where("id=?", ammeterId).Update(&ammeter.Ammeter{
		Switch: ammeterSwitch,
	})
	if err != nil {
		common.ResError(c, "修改开关状态失败")
		return
	}
	_, err = conf.Mysql.Insert(ammeter.AmmeterManageLog{
		AmmeterId:  ammeterId,
		UserId:     userId,
		Status:     ammeterSwitch,
		CreateTime: int(time.Now().Unix()),
	})
	deviceId, _ := strconv.Atoi(ammeterItem.Card)
	msg := "{\"num\":" + ammeterItem.Num + "}"
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host, conf.Conf.Tcp.Port, 1, deviceId, ammeterSwitch, msg)
	if err != nil {
		common.ResError(c, "发送控制命令失败")
		return
	}

	common.ResOk(c, "ok", nil)
}

func SetSwitchPwdSer(c *gin.Context, ammeterId int, userId int, oldPwd string, pwd string) {
	var ammeterPwd ammeter.AmmeterManageConfig
	_, err := conf.Mysql.Where("ammeter_id=?", ammeterId).Where("user_id=?", userId).Where("type=?", utils.AMMETER_MANAGE_CONFIG_TYPE_PWD).Get(&ammeterPwd)
	if err != nil {
		common.ResError(c, "获取设备配置信息失败")
		return
	}
	if ammeterPwd.Value == "" {
		_, err := conf.Mysql.Insert(ammeter.AmmeterManageConfig{
			UserId:     userId,
			AmmeterId:  ammeterId,
			Type:       utils.AMMETER_MANAGE_CONFIG_TYPE_PWD,
			Value:      utils.AMMETER_MANAGE_CONFIG_DEFAULT_PWD,
			CreateTime: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "设置默认密码失败")
			return
		}
		_, err = conf.Mysql.Where("ammeter_id=?", ammeterId).Where("user_id=?", userId).Where("type=?", utils.AMMETER_MANAGE_CONFIG_TYPE_PWD).Get(&ammeterPwd)
		if err != nil {
			common.ResError(c, "获取设备配置信息失败")
			return
		}
	}
	if ammeterPwd.Value != "" && ammeterPwd.Value != oldPwd {
		common.ResForbidden(c, "密码输入不正确，无法修改密码")
		return
	}
	_, err = conf.Mysql.Where("ammeter_id=?", ammeterId).Where("user_id=?", userId).Where("type=?", utils.AMMETER_MANAGE_CONFIG_TYPE_PWD).Update(&ammeter.AmmeterManageConfig{
		Value: pwd,
	})
	if err != nil {
		common.ResError(c, "修改设备控制密码失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func AmmeterStatisticsSer(c *gin.Context, statisticsType int, ammeterId int, startAt string, endAt string) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	var res ammeter.AmmeterStatisticRes
	if statisticsType == utils.AMMETER_DATA_TYPE_CONSUMPTION {
		var todayTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02")+" 00:00:00", location)
		var yesterdayTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, 0, -1).Format("2006-01-02")+" 00:00:00", location)
		var thisMonthTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01")+"-01 00:00:00", location)
		var lastMonthTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(0, -1, 0).Format("2006-01")+"-01 00:00:00", location)
		var thisYearTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006")+"-01-01 00:00:00", location)
		var lastYearTime, _ = time.ParseInLocation("2006-01-02 15:04:05", time.Now().AddDate(-1, 0, 0).Format("2006")+"-01-01 00:00:00", location)
		var todayStatistic ammeter.AmmeterData
		_, err := conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", todayTime.Unix()).Get(&todayStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var yesterdayStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", yesterdayTime.Unix()).Where("create_time < ?", todayTime.Unix()).Get(&yesterdayStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var monthStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", thisMonthTime.Unix()).Get(&monthStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var lastMonthStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", lastMonthTime.Unix()).Where("create_time < ?", thisMonthTime.Unix()).Get(&lastMonthStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var yearStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", thisYearTime.Unix()).Get(&yearStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var lastYearStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) - min(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Where("create_time > ?", lastYearTime.Unix()).Where("create_time < ?", thisYearTime.Unix()).Get(&lastYearStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		var mainStatistic ammeter.AmmeterData
		_, err = conf.Mysql.Select("max(value) AS value").Where("ammeter_id = ?", ammeterId).Where("type = 6").Get(&mainStatistic)
		if err != nil {
			common.ResError(c, "获取统计数据失败")
			return
		}
		res.TodayElectricityConsumption = todayStatistic.Value
		res.YesterdayElectricityConsumption = yesterdayStatistic.Value
		res.MonthElectricityConsumption = monthStatistic.Value
		res.LastMonthElectricityConsumption = lastMonthStatistic.Value
		res.YearElectricityConsumption = yearStatistic.Value
		res.LastYearElectricityConsumption = lastYearStatistic.Value
		res.MainElectricityConsumption = mainStatistic.Value
	} else {
		startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startAt, location)
		endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endAt, location)
		var datas []*ammeter.AmmeterData
		sess := conf.Mysql.NewSession()
		sess.Where("ammeter_id=?", ammeterId)
		sess.Where("type=?", statisticsType)
		sess.Where("create_time > ? AND create_time < ?", startTime.Unix(), endTime.Unix())
		sess.OrderBy("create_time")
		err := sess.Find(&datas)
		if err != nil {
			common.ResError(c, "获取设备统计数据失败")
			return
		}
		var statisticData []ammeter.StatisticForm
		for _, d := range datas {
			statisticValue := 0.0
			if d.Type == utils.AMMETER_DATA_TYPE_LEAKAGE {
				statisticValue = float64(d.Value) / 10
			}
			if d.Type == utils.AMMETER_DATA_TYPE_CURRENT {
				statisticValue = float64(d.Value) / 1000
			}
			if d.Type == utils.AMMETER_DATA_TYPE_VOLTAGE {
				statisticValue = float64(d.Value) / 1000
			}
			if d.Type == utils.AMMETER_DATA_TYPE_TEMPERATURE {
				statisticValue = float64(d.Value) / 10
			}
			if d.Type == utils.AMMETER_DATA_TYPE_POWER {
				statisticValue = float64(d.Value) / 10000
			}
			if d.Type == utils.AMMETER_DATA_TYPE_QUANTITY {
				statisticValue = float64(d.Value) / 100
			}
			if d.Type == utils.AMMETER_DATA_TYPE_SIGNAL_INTENSITY {
				statisticValue = float64(d.Value)
			}
			if d.Type == utils.AMMETER_DATA_TYPE_SWITCH {
				statisticValue = float64(d.Value)
			}
			statisticData = append(statisticData, ammeter.StatisticForm{
				Label: time.Unix(int64(d.CreateTime), 0).Format("2006-01-02 15:04:05"),
				Count: statisticValue,
			})
		}
		res.Data = statisticData
	}
	common.ResOk(c, "ok", res)
}

func WarningListSer(c *gin.Context, page int, pageSize int, warningType int, status int, startDealTime string, endDealTime string, startTime string, endTime string, dealUser string, ammeterId int) {
	var warnings []*ammeter.AmmeterWarning
	sess := conf.Mysql.NewSession()
	sess.Where("ammeter_id=?", ammeterId)
	if warningType != 0 {
		sess.Where("type=?", warningType)
	}
	if status != 0 {
		sess.Where("status=?", status)
	}
	if startDealTime != "" && endDealTime != "" {
		dealTimeStart, _ := time.Parse("2006-01-02 15:04:05", startDealTime+" 00:00:00")
		dealTimeEnd, _ := time.Parse("2006-01-02 15:04:05", endDealTime+" 23:59:59")
		fmt.Println(startDealTime, dealTimeStart)
		sess.Where("deal_time > ? AND deal_time < ?", dealTimeStart.Unix(), dealTimeEnd.Unix())
	}
	if startTime != "" && endTime != "" {
		timeStart, _ := time.Parse("2006-01-02 15:04:05", startTime+" 00:00:00")
		timeEnd, _ := time.Parse("2006-01-02 15:04:05", endTime+" 23:59:59")
		sess.Where("create_time > ? AND create_time < ?", timeStart.Unix(), timeEnd.Unix())
	}
	if dealUser != "" {
		var users []*rbac.User
		var userIds []int
		err := conf.Mysql.Where("name like ?", "%"+dealUser+"%").Find(&users)
		if err != nil {
			common.ResError(c, "获取用户列表失败")
			return
		}
		for _, user := range users {
			userIds = append(userIds, user.Id)
		}
		sess.In("deal_user", userIds)
	}
	count, err := sess.OrderBy("id DESC").Limit(pageSize, (page-1)*pageSize).FindAndCount(&warnings)
	if err != nil {
		common.ResError(c, "查询报警列表失败")
		return
	}
	var userIds []int
	var ammeterIds []int
	for _, w := range warnings {
		userIds = append(userIds, w.DealUser)
		ammeterIds = append(ammeterIds, w.AmmeterId)
	}
	var ammeterItems []*ammeter.Ammeter
	var ammeterMapping = make(map[int]*ammeter.Ammeter)
	err = conf.Mysql.In("id", ammeterIds).Find(&ammeterItems)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	for _, item := range ammeterItems {
		ammeterMapping[item.Id] = item
	}
	var users []*rbac.User
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
		w.AmmeterCard = ammeterMapping[w.AmmeterId].Card
		w.AmmeterLocation = ammeterMapping[w.AmmeterId].Location
		if w.DealTime == 0 {
			w.DealTimeStr = "-"
		} else {
			w.DealTimeStr = time.Unix(int64(w.DealTime), 0).Format("2006-01-02 15:04:05")
		}
		w.WarningTimeStr = time.Unix(int64(w.CreateTime), 0).Format("2006-01-02 15:04:05")
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: warnings})
}

func ChangeWarningStatusSer(c *gin.Context, warningId int, status int, userId int) {
	_, err := conf.Mysql.Where("id=?", warningId).Update(&ammeter.AmmeterWarning{
		Status:   status,
		DealTime: int(time.Now().Unix()),
		DealUser: userId,
	})
	if err != nil {
		common.ResError(c, "处理警告失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func ConfigInfoSer(c *gin.Context, ammeterId int) {
	var config ammeter.AmmeterConfig
	_, err := conf.Mysql.Where("ammeter_id=?", ammeterId).Get(&config)
	if err != nil {
		common.ResError(c, "获取设备配置失败")
		return
	}
	common.ResOk(c, "ok", config)
}

func UpdateConfigSer(c *gin.Context, configUpdate ammeter.ConfigUpdateReq) {
	var deviceItem ammeter.Ammeter
	_, err := conf.Mysql.Where("id=?", configUpdate.AmmeterId).Get(&deviceItem)
	if err != nil {
		common.ResError(c, "获取设备信息失败")
		return
	}
	sqlStr := fmt.Sprintf("UPDATE ammeter_config SET %s=%d WHERE ammeter_id=%d", configUpdate.UpdateType, configUpdate.UpdateValue, configUpdate.AmmeterId)
	_, err = conf.Mysql.Query(sqlStr)
	if err != nil {
		common.ResError(c, "修改配置失败")
		return
	}
	updateValueHexString := common.DecimalToHex(configUpdate.UpdateValue)
	reportType := 4
	deviceIdInt, _ := strconv.Atoi(deviceItem.Card)
	deviceNumInt, _ := strconv.Atoi(deviceItem.Num)
	msgByte, _ := json.Marshal(ammeter.ConfigUpdateMsgData{Num: deviceNumInt, ParamType: utils.UPDATE_AMMETER_CONFIG_PARAMS[configUpdate.UpdateType], ParamValue: updateValueHexString})
	err = common.CommonSendDeviceReport(conf.Conf.Tcp.Host, conf.Conf.Tcp.Port, reportType, deviceIdInt, 0, string(msgByte))
	common.ResOk(c, "ok", nil)
}

func AddTestData(c *gin.Context, ammeterId int, dataType int, value int, createTime string) {
	t, _ := time.Parse("2006-01-02 15:04:05", createTime)
	_, _ = conf.Mysql.Insert(&ammeter.AmmeterData{
		AmmeterId:  ammeterId,
		Type:       dataType,
		Value:      value,
		CreateTime: int(t.Unix()),
	})
	common.ResOk(c, "ok", nil)
}