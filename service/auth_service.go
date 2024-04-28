package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
)

func LoginSrv(c *gin.Context, phone string, password string) {
	var user model.User
	_, err := conf.Mysql.Where("phone = ?", phone).Get(&user)
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "获取用户失败")
		return
	}
	if user.Id == 0 {
		common.ResForbidden(c, "尚未注册")
		return
	}
	if user.Password != password {
		common.ResForbidden(c, "密码错误")
		return
	}
	var userRole model.UserRole
	_, err = conf.Mysql.Where("id = ?", user.Id).Get(&userRole)
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

func GetAllPermissionSer(c *gin.Context) {
	var permissions []*model.Permission
	err := conf.Mysql.OrderBy("sort").Find(&permissions)
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

func GetAllRolesSer(c *gin.Context, page int, pageSize int) {
	var roles []*model.Role
	count, err := conf.Mysql.Limit(pageSize, (page-1)*pageSize).FindAndCount(&roles)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: roles})
}

func AddRoleSer(c *gin.Context, roleAdd model.RoleAddReq) {
	var role = model.Role{
		Name: roleAdd.Name,
	}
	_, err := conf.Mysql.Insert(&role)
	if err != nil {
		common.ResError(c, "添加角色失败")
		return
	}
	var rolePermissions []model.RolePermission
	for _, pId := range roleAdd.Permission {
		if pId == 0 {
			continue
		}
		rolePermissions = append(rolePermissions, model.RolePermission{
			RoleId:       role.Id,
			PermissionId: pId,
		})
	}
	_, err = conf.Mysql.Insert(&rolePermissions)
	if err != nil {
		common.ResError(c, "添加角色权限关系失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func DeleteRoleSer(c *gin.Context, roleId int) {
	_, err := conf.Mysql.Where("id = ?", roleId).Delete(&model.Role{})
	if err != nil {
		common.ResError(c, "删除角色失败")
		return
	}
	_, err = conf.Mysql.Where("role_id = ?", roleId).Delete(&model.RolePermission{})
	if err != nil {
		common.ResError(c, "删除角色权限关系失败")
		return
	}
	_, err = conf.Mysql.Where("role_id = ?", roleId).Delete(&model.UserRole{})
	if err != nil {
		common.ResError(c, "删除用户角色关系失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func AddPermissionSer(c *gin.Context, parentId int, path string, icon string, sort int, label string, desc string, component string) {
	var newPermission = model.Permission{
		Path:      path,
		Icon:      icon,
		Sort:      sort,
		ParentId:  parentId,
		Label:     label,
		Desc:      desc,
		Component: component,
	}
	_, err := conf.Mysql.Insert(newPermission)
	if err != nil {
		common.ResError(c, "添加权限失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func RemovePermissionSer(c *gin.Context, id int) {
	var permissions []model.Permission
	err := conf.Mysql.Where("id = ?", id).Or("parent_id = ?", id).Find(&permissions)
	if err != nil {
		common.ResError(c, "搜索节点及其子节点失败")
		return
	}
	var permissionIds []int
	for _, p := range permissions {
		permissionIds = append(permissionIds, p.Id)
	}
	_, err = conf.Mysql.In("id", permissionIds).Delete(&model.Permission{})
	if err != nil {
		common.ResError(c, "节点删除失败")
		return
	}
	_, err = conf.Mysql.In("permission_id", permissionIds).Delete(&model.RolePermission{})
	if err != nil {
		common.ResError(c, "节点关联角色信息删除失败")
		return
	}

	common.ResOk(c, "ok", nil)
}
