package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
	"management-backend/utils"
	"time"
)

func GetUserInfoSer(c *gin.Context, userId int) {
	var user model.User
	_, err := conf.Mysql.Where("id = ?", userId).Get(&user)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	var userRole model.UserRole
	_, err = conf.Mysql.Where("user_id = ?", userId).Get(&userRole)
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
	var organizationUsers []*model.OrganizationUser
	err = conf.Mysql.Where("user_id = ?", userId).Find(&organizationUsers)
	if err != nil {
		common.ResError(c, "获取用户关联组织信息失败")
		return
	}
	for _, i := range organizationUsers {
		user.OrganizationIds = append(user.OrganizationIds, i.OrganizationId)
	}
	common.ResOk(c, "ok", user)
}

func GetUserListSrv(c *gin.Context, userName string, page int, pageSize int) {
	userId := c.Param("user_id")
	var users []*model.User
	sess := conf.Mysql.NewSession()
	if userName != "" {
		sess.Where("name=?", userName)
	}
	if userId != "1" {
		sess.Where("id != 1")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&users)
	if err != nil {
		common.ResError(c, "获取用户列表失败")
		return
	}
	if len(users) <= 0 {
		common.ResOk(c, "ok", utils.CommonListRes{Count: 0, Data: []*model.User{}})
		return
	}
	var userIds []int
	for _, user := range users {
		userIds = append(userIds, user.Id)
	}
	var userRoles []model.UserRole
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
	var roles []model.Role
	err = conf.Mysql.In("id", rolesIds).Find(&roles)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	var roleMapping = make(map[int]model.Role)
	for _, role := range roles {
		roleMapping[role.Id] = role
	}
	for _, u := range users {
		u.RoleStr = roleMapping[userRoleIdMapping[u.Id]]
	}

	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: users})
}

func AddUserSer(c *gin.Context, req model.AddUserReq) {
	var userId int
	if req.Id == 0 {
		newToken, err := utils.GenerateRandomToken(20)
		if err != nil {
			common.ResError(c, "生成token失败")
			return
		}
		var user = model.User{
			Name:     req.Name,
			Password: req.Password,
			Phone:    req.Phone,
			Token:    newToken,
		}
		_, err = conf.Mysql.Insert(&user)
		if err != nil {
			common.ResError(c, "添加用户失败")
			return
		}
		userId = user.Id
	} else {
		var user model.User
		_, err := conf.Mysql.Where("id = ?", req.Id).Get(&user)
		if err != nil {
			common.ResError(c, "获取用户失败")
			return
		}
		_, err = conf.Mysql.Where("id = ?", req.Id).Update(&model.User{
			Name:     req.Name,
			Password: req.Password,
			Phone:    req.Phone,
		})
		if err != nil {
			common.ResError(c, "更新用户失败")
			return
		}
		userId = req.Id
	}

	_, err := conf.Mysql.Where("user_id = ?", userId).Delete(&model.UserRole{})
	if err != nil {
		common.ResError(c, "删除原有用户角色失败")
		return
	}
	var userRole = model.UserRole{
		UserId: userId,
		RoleId: req.RoleId,
	}
	_, err = conf.Mysql.Insert(&userRole)
	if err != nil {
		common.ResError(c, "添加用户角色关联信息失败")
		return
	}

	_, err = conf.Mysql.Where("user_id = ?", userId).Delete(&model.OrganizationUser{})
	if err != nil {
		common.ResError(c, "删除原有用户组织关联关系失败")
		return
	}
	var insertOrganization []*model.OrganizationUser
	for _, i := range req.OrganizationIds {
		insertOrganization = append(insertOrganization, &model.OrganizationUser{
			UserId:         userId,
			OrganizationId: i,
			Ts:             time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	_, err = conf.Mysql.Insert(&insertOrganization)
	if err != nil {
		common.ResError(c, "添加用户组织关联信息失败")
		return
	}

	common.ResOk(c, "ok", nil)
}

func DeleteUserSer(c *gin.Context, req model.DeleteUserReq) {
	_, err := conf.Mysql.Where("id = ?", req.Id).Delete(&model.User{})
	if err != nil {
		common.ResError(c, "删除用户失败")
		return
	}
	_, err = conf.Mysql.Where("user_id = ?", req.Id).Delete(&model.UserRole{})
	if err != nil {
		common.ResError(c, "删除用户关联角色失败")
		return
	}
	_, err = conf.Mysql.Where("user_id = ?", req.Id).Delete(&model.OrganizationUser{})
	if err != nil {
		common.ResError(c, "组织用户信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func GetUserPermissionSer(c *gin.Context, userId int) {
	var userRole model.UserRole
	_, err := conf.Mysql.Where("user_id = ?", userId).Get(&userRole)
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
