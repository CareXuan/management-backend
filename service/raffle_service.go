package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"prize-draw/common"
	"prize-draw/conf"
	"prize-draw/model"
	"time"
)

type raffleResult struct {
	Msg  string `json:"msg"`
	Pic  string `json:"pic"`
	Time string `json:"time"`
}

func RaffleOne(c *gin.Context, req model.RaffleOneReq) {
	var pointCount model.GiftPackage
	has, err := conf.Mysql.Where("gift_id = ?", 0).Where("delete_at = 0").Get(&pointCount)
	if err != nil {
		common.ResError(c, "获取抽卡点信息失败")
		return
	}
	var raffleConfig model.Config
	_, err = conf.Mysql.Get(&raffleConfig)
	if err != nil {
		common.ResError(c, "获取配置失败")
		return
	}
	oneRaffleCount := raffleConfig.OnePoint
	if req.Count == 10 {
		oneRaffleCount = raffleConfig.TenPoint
	}
	if !has || pointCount.Count < oneRaffleCount {
		common.ResForbidden(c, "抽卡点不足")
		return
	}
	var resLogs []raffleResult
	var probabilityItems []*model.GiftGroupGift
	err = conf.Mysql.Where("group_id = ?", req.GiftGroupId).Where("delete_at = 0").Find(&probabilityItems)
	if err != nil {
		common.ResError(c, "获取奖池信息失败")
		return
	}
	k := 0
	for k < req.Count {
		maxProbability := 0
		var probabilityMapping = make(map[string]int)
		for _, i := range probabilityItems {
			if i.Probability == 0 {
				continue
			}
			maxProbability += i.Probability
			probabilityMapping[i.Level] = maxProbability
		}
		activeProbability := common.RandomClosedInterval(1, maxProbability)
		resLevel := "D"
		for level, probability := range probabilityMapping {
			if activeProbability <= probability && resLevel > level {
				resLevel = level
			}
		}
		var activeLevelGift []*model.Gift
		err = conf.Mysql.Where("level = ?", resLevel).Where("delete_at = 0").Find(&activeLevelGift)
		if err != nil {
			common.ResError(c, "获取礼物信息失败")
			return
		}
		if len(activeLevelGift) > 0 {
			giftNum := common.RandomClosedInterval(0, len(activeLevelGift)-1)
			var activeGift model.Gift
			for k, v := range activeLevelGift {
				if k == giftNum {
					activeGift = *v
					break
				}
			}
			has, err := conf.Mysql.Where("gift_id = ?", activeGift.Id).Where("delete_at = 0").Get(&model.GiftPackage{})
			if err != nil {
				common.ResError(c, "查询背包失败")
				return
			}
			resLog := ""
			if has {
				var giftPackageItem model.GiftPackage
				has, err := conf.Mysql.Where("gift_id = ?", 0).Where("delete_at = 0").Get(&giftPackageItem)
				if err != nil {
					common.ResError(c, "获取礼物背包信息失败")
					return
				}
				if has {
					_, err = conf.Mysql.Where("gift_id = ?", 0).Where("delete_at = 0").Update(model.GiftPackage{
						Count: giftPackageItem.Count + activeGift.CrushCnt,
					})
					if err != nil {
						common.ResError(c, "发放抽卡点失败")
						return
					}
				} else {
					_, err = conf.Mysql.Insert(model.GiftPackage{
						GiftId:     0,
						Count:      activeGift.CrushCnt,
						Consumable: 1,
						CreateAt:   int(time.Now().Unix()),
					})
					if err != nil {
						common.ResError(c, "发放抽卡点失败")
						return
					}
				}
				resLog = fmt.Sprintf("获取%s级【%s】X1，自动分解为抽卡点%d个。", activeGift.Level, activeGift.Name, activeGift.CrushCnt)
			} else {
				_, err := conf.Mysql.Insert(&model.GiftPackage{
					GiftId:     activeGift.Id,
					Count:      1,
					CreateAt:   int(time.Now().Unix()),
					Consumable: activeGift.Consumable,
				})
				if err != nil {
					common.ResError(c, "礼物加入背包失败")
					return
				}
				_, err = conf.Mysql.Where("id = ?", activeGift.Id).Update(&model.Gift{
					Count: 1,
				})
				if err != nil {
					common.ResError(c, "礼物信息修改失败")
					return
				}
				resLog = fmt.Sprintf("获取%s级【%s】X1。", activeGift.Level, activeGift.Name)
			}
			resLogs = append(resLogs, raffleResult{
				Msg:  resLog,
				Pic:  activeGift.Pic,
				Time: time.Now().Format("2006-01-02 15:04:05"),
			})
		} else {
			common.ResForbidden(c, "抽卡失败")
			return
		}
		k += 1
	}
	var nowPointCount model.GiftPackage
	has, err = conf.Mysql.Where("gift_id = ?", 0).Where("delete_at = 0").Get(&nowPointCount)
	if err != nil {
		common.ResError(c, "获取抽卡点信息失败")
		return
	}
	_, err = conf.Mysql.Where("gift_id = ?", 0).Where("delete_at = 0").Update(&model.GiftPackage{
		Count: nowPointCount.Count - oneRaffleCount,
	})
	if err != nil {
		common.ResError(c, "扣除抽卡点失败")
		return
	}
	common.ResOk(c, "ok", resLogs)
}

func AppRaffleConfig(c *gin.Context) {
	type appRaffleConfig struct {
		PointCount       int `json:"point_count"`
		AchievementCount int `json:"achievement_count"`
		RaffleOne        int `json:"raffle_one"`
		RaffleTen        int `json:"raffle_ten"`
	}
	var configItem model.Config
	_, err := conf.Mysql.Where("id = 1").Get(&configItem)
	if err != nil {
		common.ResError(c, "获取配置失败")
		return
	}
	var packageItem model.GiftPackage
	_, err = conf.Mysql.Where("gift_id = 0").Get(&packageItem)
	if err != nil {
		common.ResError(c, "获取抽卡点失败")
		return
	}
	var achievementItems []*model.Achievement
	err = conf.Mysql.Where("is_finish = 1").Where("is_receive = 0").Find(&achievementItems)
	if err != nil {
		common.ResError(c, "获取成就信息失败")
		return
	}
	common.ResOk(c, "ok", appRaffleConfig{
		PointCount:       packageItem.Count,
		AchievementCount: len(achievementItems),
		RaffleOne:        configItem.OnePoint,
		RaffleTen:        configItem.TenPoint,
	})
}
