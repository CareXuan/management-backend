package service

import (
	"fmt"
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
	var userRole model.UserRole
	_, err = conf.Mysql.Where("id = ?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var role model.Role
	_, err = conf.Mysql.Where("id = ?", userRole.RoleId).Get(&role)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	user.RoleStr = role
	common.ResOk(c, "ok", user)
}

func GetUserPermissionSer(c *gin.Context, userId int) {
	var userRole model.UserRole
	_, err := conf.Mysql.Where("id = ?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var rolePermission []*model.RolePermission
	err = conf.Mysql.Where("role_id = ?", userRole.RoleId).Find(&rolePermission)
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "获取用户权限失败")
		return
	}
	var permissionIds []int
	for _, rp := range rolePermission {
		permissionIds = append(permissionIds, rp.PermissionId)
	}
	var permissions []*model.Permission
	err = conf.Mysql.In("id", permissionIds).OrderBy("sort").Find(&permissions)
	if err != nil {
		common.ResError(c, "获取权限信息失败")
		return
	}
	var fatherPermissions []*model.Permission
	var childrenPermissions = make(map[int][]*model.Permission)
	for _, p := range permissions {
		if p.ParentId == 0 {
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
