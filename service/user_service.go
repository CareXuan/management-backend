package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model/rbac"
	"management-backend/utils"
)

func GetUserInfoSer(c *gin.Context, userId int) {
	var user rbac.User
	_, err := conf.Mysql.Where("id = ?", userId).Get(&user)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	var userRole rbac.UserRole
	_, err = conf.Mysql.Where("id = ?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var role rbac.Role
	_, err = conf.Mysql.Where("id = ?", userRole.RoleId).Get(&role)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	user.RoleStr = role
	common.ResOk(c, "ok", user)
}

func GetUserListSrv(c *gin.Context, userName string, page int, pageSize int) {
	var users []*rbac.User
	sess := conf.Mysql.NewSession()
	if userName != "" {
		sess.Where("name=?", userName)
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&users)
	if err != nil {
		common.ResError(c, "获取用户列表失败")
		return
	}
	if len(users) <= 0 {
		common.ResOk(c, "ok", utils.CommonListRes{Count: 0, Data: []*rbac.User{}})
		return
	}
	var userIds []int
	for _, user := range users {
		userIds = append(userIds, user.Id)
	}
	var userRoles []rbac.UserRole
	err = conf.Mysql.In("user_id", userIds).Find(&userRoles)
	if err != nil {
		common.ResError(c, "获取用户角色关联关系失败")
		return
	}
	var rolesIds []int
	var userRoleIdMapping = make(map[int]int)
	for _, userRole := range userRoles {
		rolesIds = append(rolesIds, userRole.RoleId)
		userRoleIdMapping[userRole.UserId] = userRole.RoleId
	}
	var roles []rbac.Role
	err = conf.Mysql.In("id", rolesIds).Find(&roles)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	var roleMapping = make(map[int]rbac.Role)
	for _, role := range roles {
		roleMapping[role.Id] = role
	}
	for _, u := range users {
		u.RoleStr = roleMapping[userRoleIdMapping[u.Id]]
	}

	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: users})
}

func AddUserSer(c *gin.Context, id int, name string, password string, phone string, roleId int) {
	if id == 0 {
		newToken, err := utils.GenerateRandomToken(20)
		if err != nil {
			common.ResError(c, "生成token失败")
			return
		}
		var user = rbac.User{
			Name:     name,
			Password: password,
			Phone:    phone,
			Token:    newToken,
		}
		_, err = conf.Mysql.Insert(&user)
		if err != nil {
			common.ResError(c, "添加用户失败")
			return
		}
		var userRole = rbac.UserRole{
			UserId: int(user.Id),
			RoleId: roleId,
		}
		_, err = conf.Mysql.Insert(&userRole)
		if err != nil {
			common.ResError(c, "添加用户角色关联信息失败")
			return
		}
	} else {
		var user rbac.User
		_, err := conf.Mysql.Where("id = ?", id).Get(&user)
		if err != nil {
			common.ResError(c, "获取用户失败")
			return
		}
		_, err = conf.Mysql.Where("id = ?", id).Update(&rbac.User{
			Name:     name,
			Password: password,
			Phone:    phone,
		})
		if err != nil {
			common.ResError(c, "更新用户失败")
			return
		}
		_, err = conf.Mysql.Where("user_id = ?", id).Delete(&rbac.UserRole{})
		if err != nil {
			common.ResError(c, "删除原有用户角色失败")
			return
		}
		var userRole = rbac.UserRole{
			UserId: int(user.Id),
			RoleId: roleId,
		}
		_, err = conf.Mysql.Insert(&userRole)
		if err != nil {
			common.ResError(c, "添加用户角色关联信息失败")
			return
		}
	}

	common.ResOk(c, "ok", nil)
}

func GetUserPermissionSer(c *gin.Context, userId int) {
	var userRole rbac.UserRole
	_, err := conf.Mysql.Where("id = ?", userId).Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var rolePermission []*rbac.RolePermission
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
	var permissions []*rbac.Permission
	err = conf.Mysql.In("id", permissionIds).OrderBy("sort").Find(&permissions)
	if err != nil {
		common.ResError(c, "获取权限信息失败")
		return
	}
	var fatherPermissions []*rbac.Permission
	var childrenPermissions = make(map[int][]*rbac.Permission)
	for _, p := range permissions {
		p.Title = p.Label
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
