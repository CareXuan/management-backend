package service

import (
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
)

func GetUserInfoSer(c *gin.Context, userId int) {
	var user model.User
	_, err := conf.Mysql.Where("id = ?", userId).Get(&user)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	common.ResOk(c, "ok", user)
}

func GetUserPermissionSer(c *gin.Context, userId int) {
	var userPermission []*model.UserPermission
	err := conf.Mysql.Where("id = ?", userId).OrderBy("sort").Find(&userPermission)
	if err != nil {
		common.ResError(c, "获取用户权限失败")
		return
	}
	var permissionIds []int
	for _, up := range userPermission {
		permissionIds = append(permissionIds, up.PermissionId)
	}
	var permissions []*model.Permission
	err = conf.Mysql.In("id", permissionIds).Find(&permissions)
	if err != nil {
		common.ResError(c, "获取权限信息失败")
		return
	}
	var fatherPermissions []*model.Permission
	var childrenPermissions = make(map[int][]*model.Permission)
	for _, p := range permissions {
		if p.ParentId != 0 {
			fatherPermissions = append(fatherPermissions, p)
		} else {
			childrenPermissions[p.ParentId] = append(childrenPermissions[p.ParentId], p)
		}
	}

	for _, fp := range fatherPermissions {
		fp.Children = childrenPermissions[fp.Id]
	}
	common.ResOk(c, "ok", fatherPermissions)
}
