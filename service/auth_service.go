package service

import (
	"data_verify/common"
	"data_verify/conf"
	"data_verify/model"
	"data_verify/utils"
	"fmt"
	"github.com/gin-gonic/gin"
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

func GetRolePermissionSer(c *gin.Context, roleId int) {
	var role model.Role
	_, err := conf.Mysql.Where("id = ?", roleId).Get(&role)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	var rolePermissions []*model.RolePermission
	err = conf.Mysql.Where("role_id = ?", roleId).Find(&rolePermissions)
	if err != nil {
		common.ResError(c, "获取角色关联权限失败")
		return
	}
	var permissionIds []int
	for _, p := range rolePermissions {
		permissionIds = append(permissionIds, p.PermissionId)
	}
	common.ResOk(c, "ok", model.RoleInfoRes{Id: role.Id, Name: role.Name, Permission: permissionIds})
}

func AddRoleSer(c *gin.Context, roleAdd model.RoleAddReq) {
	if roleAdd.Id == 0 {
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
	} else {
		var role model.Role
		_, err := conf.Mysql.Where("id = ?", roleAdd.Id).Get(&role)
		if err != nil {
			common.ResError(c, "获取角色信息失败")
			return
		}
		_, err = conf.Mysql.Where("id = ?", roleAdd.Id).Update(&model.Role{
			Name: roleAdd.Name,
		})
		if err != nil {
			common.ResError(c, "修改角色信息失败")
			return
		}
		_, err = conf.Mysql.Where("role_id = ?", roleAdd.Id).Delete(&model.RolePermission{})
		if err != nil {
			common.ResError(c, "删除原有角色关联权限失败")
			return
		}
		var newRolePermissions []model.RolePermission
		for _, pId := range roleAdd.Permission {
			if pId == 0 {
				continue
			}
			newRolePermissions = append(newRolePermissions, model.RolePermission{
				RoleId:       role.Id,
				PermissionId: pId,
			})
		}
		_, err = conf.Mysql.Insert(&newRolePermissions)
		if err != nil {
			common.ResError(c, "添加角色关联权限失败")
			return
		}
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

func AllGroupSer(c *gin.Context, name string, page, pageSize int) {
	sess := conf.Mysql.NewSession()
	var groupItems []*model.Group
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.OrderBy("id DESC").Limit(pageSize, (page-1)*pageSize).FindAndCount(&groupItems)
	if err != nil {
		common.ResError(c, "获取组别列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: groupItems})
}

func GroupInfoSer(c *gin.Context, id int) {
	var groupInfo model.Group
	_, err := conf.Mysql.Where("id = ?", id).Get(&groupInfo)
	if err != nil {
		common.ResError(c, "获取组别详情失败")
		return
	}
	groupId := groupInfo.Id
	var groupUsers []*model.GroupUser
	err = conf.Mysql.Where("group_id = ?", groupId).Find(&groupUsers)
	if err != nil {
		common.ResError(c, "获取组别用户关联关系失败")
		return
	}
	var userIds []int
	for _, i := range groupUsers {
		userIds = append(userIds, i.UserId)
	}
	var userItems []*model.User
	err = conf.Mysql.In("id", userIds).Find(&userItems)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	groupInfo.Users = userItems
	common.ResOk(c, "ok", groupInfo)
}

func AddGroupSer(c *gin.Context, req model.GroupAddReq) {
	if req.Id == 0 {
		var newGroupData = model.Group{Name: req.Name}
		_, err := conf.Mysql.Insert(&newGroupData)
		if err != nil {
			common.ResError(c, "添加组别失败")
			return
		}
		var groupUserData []*model.GroupUser
		for _, i := range req.UserIds {
			groupUserData = append(groupUserData, &model.GroupUser{
				GroupId: newGroupData.Id,
				UserId:  i,
			})
		}
		_, err = conf.Mysql.Insert(&groupUserData)
		if err != nil {
			common.ResError(c, "添加组别用户关联关系失败")
			return
		}

	} else {
		sess := conf.Mysql.NewSession()
		_, err := sess.Where("group_id = ?", req.Id).Delete(&model.GroupUser{})
		if err != nil {
			sess.Rollback()
			common.ResError(c, "删除组别用户关联关系失败")
			return
		}
		_, err = sess.Where("id = ?", req.Id).Update(&model.Group{Name: req.Name})
		if err != nil {
			sess.Rollback()
			common.ResError(c, "修改组别信息失败")
			return
		}
		var groupUserData []*model.GroupUser
		for _, i := range req.UserIds {
			groupUserData = append(groupUserData, &model.GroupUser{
				GroupId: req.Id,
				UserId:  i,
			})
		}
		_, err = conf.Mysql.Insert(&groupUserData)
		if err != nil {
			sess.Rollback()
			common.ResError(c, "写入组别用户关联数据失败")
			return
		}
		sess.Commit()
	}
	common.ResOk(c, "ok", nil)
}

func DeleteGroupSer(c *gin.Context, id int) {
	sess := conf.Mysql.NewSession()
	_, err := sess.Where("id = ?", id).Delete(&model.Group{})
	if err != nil {
		sess.Rollback()
		common.ResError(c, "删除组别失败")
		return
	}
	_, err = conf.Mysql.Where("group_id = ?", id).Delete(&model.GroupUser{})
	if err != nil {
		sess.Rollback()
		common.ResError(c, "删除组别用户关联关系失败")
		return
	}
	sess.Commit()
	common.ResOk(c, "ok", nil)
}
