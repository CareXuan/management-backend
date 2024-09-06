package service

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetDeviceListSer(c *gin.Context, name string, page, pageSize int) {
	var devices []*model.Device
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&devices)
	if err != nil {
		common.ResError(c, "获取设备列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: devices})
}

func GetDeviceInfoSer(c *gin.Context, id int) {
	var device model.Device
	_, err := conf.Mysql.Where("id = ?", id).Get(&device)
	if err != nil {
		common.ResError(c, "获取设备详情失败")
		return
	}
	common.ResOk(c, "ok", device)
}

func AddDeviceSer(c *gin.Context, device model.Device) {
	if device.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", device.Id).Update(device)
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(model.Device{
			Name:      device.Name,
			Pic:       device.Pic,
			Cert:      device.Cert,
			UseTime:   device.UseTime,
			CreatedAt: time.Now().Unix(),
		})
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
	}
	common.ResOk(c, "ok", nil)
}

// ================================== 套餐 ==================================

func GetPackageListSer(c *gin.Context, name string, page, pageSize int) {
	var packages []*model.DevicePackage
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&packages)
	if err != nil {
		common.ResError(c, "获取套餐列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: packages})
}

func GetPackageInfoSer(c *gin.Context, id int) {
	var packageInfo model.DevicePackage
	_, err := conf.Mysql.Where("id = ?", id).Get(&packageInfo)
	if err != nil {
		common.ResError(c, "获取套餐详情失败")
		return
	}
	var details []model.DevicePackageDetail
	err = conf.Mysql.Where("package_id = ?", packageInfo.Id).Find(&details)
	if err != nil {
		common.ResError(c, "获取套餐详情失败")
		return
	}
	packageInfo.Details = details
	common.ResOk(c, "ok", packageInfo)
}

func AddPackageSer(c *gin.Context, packageAdd model.DevicePackageAddReq) {
	packageAddId := 0
	if packageAdd.Id != 0 {
		packageAddId = packageAdd.Id
		_, err := conf.Mysql.Where("id = ?", packageAddId).Update(&model.DevicePackage{
			Name:   packageAdd.Name,
			Cost:   packageAdd.Cost,
			Status: model.PACKAGE_STATUS_OPEN,
		})
		if err != nil {
			common.ResError(c, "修改设备失败")
			return
		}
		_, err = conf.Mysql.Where("package_id = ?", packageAddId).Delete(&model.DevicePackageDetail{})
		if err != nil {
			common.ResError(c, "删除现有套餐失败")
			return
		}
	} else {
		var packageAddInfo = model.DevicePackage{
			Name:      packageAdd.Name,
			Cost:      packageAdd.Cost,
			Status:    model.PACKAGE_STATUS_OPEN,
			CreatedAt: time.Now().Unix(),
		}
		_, err := conf.Mysql.Insert(&packageAddInfo)
		if err != nil {
			common.ResError(c, "添加设备失败")
			return
		}
		packageAddId = packageAddInfo.Id
	}
	var packageDetails []model.DevicePackageDetail
	for _, detail := range packageAdd.Details {
		packageDetails = append(packageDetails, model.DevicePackageDetail{
			PackageId: packageAddId,
			DeviceId:  detail.DeviceId,
			Type:      detail.Type,
			Value:     detail.Value,
		})
	}
	_, err := conf.Mysql.Insert(&packageDetails)
	if err != nil {
		common.ResError(c, "添加套餐详情失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func PackageChangeStatusSer(c *gin.Context, packageChangeStatus model.DevicePackageChangeStatusReq) {
	_, err := conf.Mysql.Where("id = ?", packageChangeStatus.Id).Update(model.DevicePackage{
		Status: packageChangeStatus.Status,
	})
	if err != nil {
		common.ResError(c, "修改套餐状态失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
