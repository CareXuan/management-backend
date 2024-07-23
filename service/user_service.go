package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"my-gpt-server/model"
	"my-gpt-server/utils"
	"time"
)

func GetUserInfoSer(c *gin.Context, userId int) {
	var user model.User
	_, err := conf.Mysql.Where("id = ?", userId).Where("deleted_at = 0").Get(&user)
	if err != nil {
		common.ResError(c, "获取用户信息失败")
		return
	}
	var userRole model.UserRole
	_, err = conf.Mysql.Where("user_id = ?", userId).Where("deleted_at = 0").Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var role model.Role
	_, err = conf.Mysql.Where("id = ?", userRole.RoleId).Where("deleted_at = 0").Get(&role)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	user.RoleStr = role
	common.ResOk(c, "ok", user)
}

func GetUserListSrv(c *gin.Context, userName string, page int, pageSize int) {
	var users []*model.User
	sess := conf.Mysql.NewSession().Where("deleted_at = 0")
	if userName != "" {
		sess.Where("name=?", userName)
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
	err = conf.Mysql.In("user_id", userIds).Where("deleted_at = 0").Find(&userRoles)
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

func AddUserSer(c *gin.Context, id int, name string, password string, phone string, roleId int) {
	if id == 0 {
		newToken, err := utils.GenerateRandomToken(20)
		if err != nil {
			common.ResError(c, "生成token失败")
			return
		}
		var user = model.User{
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
		var userRole = model.UserRole{
			UserId:    int(user.Id),
			RoleId:    roleId,
			CreatedAt: int(time.Now().Unix()),
			UpdatedAt: int(time.Now().Unix()),
		}
		_, err = conf.Mysql.Insert(&userRole)
		if err != nil {
			common.ResError(c, "添加用户角色关联信息失败")
			return
		}
	} else {
		var user model.User
		_, err := conf.Mysql.Where("id = ?", id).Where("deleted_at = 0").Get(&user)
		if err != nil {
			common.ResError(c, "获取用户失败")
			return
		}
		_, err = conf.Mysql.Where("id = ?", id).Where("deleted_at = 0").Update(&model.User{
			Name:     name,
			Password: password,
			Phone:    phone,
		})
		if err != nil {
			common.ResError(c, "更新用户失败")
			return
		}
		_, err = conf.Mysql.Where("user_id = ?", id).Update(&model.UserRole{DeletedAt: int(time.Now().Unix())})
		if err != nil {
			common.ResError(c, "删除原有用户角色失败")
			return
		}
		var userRole = model.UserRole{
			UserId:    int(user.Id),
			RoleId:    roleId,
			CreatedAt: int(time.Now().Unix()),
			UpdatedAt: int(time.Now().Unix()),
		}
		_, err = conf.Mysql.Insert(&userRole)
		if err != nil {
			common.ResError(c, "添加用户角色关联信息失败")
			return
		}
	}

	common.ResOk(c, "ok", nil)
}

func DeleteUserSer(c *gin.Context, userId int) {
	_, err := conf.Mysql.Where("id = ?", userId).Where("deleted_at = 0").Update(model.User{DeletedAt: time.Now().Unix()})
	if err != nil {
		common.ResError(c, "删除用户失败")
		return
	}
	_, err = conf.Mysql.Where("user_id = ?", userId).Where("deleted_at = 0").Update(model.UserRole{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "删除用户角色关联信息失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func GetUserPermissionSer(c *gin.Context, userId int) {
	var userRole model.UserRole
	_, err := conf.Mysql.Where("user_id = ?", userId).Where("deleted_at = 0").Get(&userRole)
	if err != nil {
		common.ResError(c, "获取用户角色失败")
		return
	}
	var rolePermission []*model.RolePermission
	err = conf.Mysql.Where("role_id = ?", userRole.RoleId).Where("deleted_at = 0").Find(&rolePermission)
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
