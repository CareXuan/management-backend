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

func LoginSrv(c *gin.Context, phone string, password string) {
	var user model.User
	_, err := conf.Mysql.Where("phone = ?", phone).Where("deleted_at = 0").Get(&user)
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
	_, err = conf.Mysql.Where("id = ?", user.Id).Where("deleted_at = 0").Get(&userRole)
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

func GetAllPermissionSer(c *gin.Context) {
	var permissions []*model.Permission
	err := conf.Mysql.Where("deleted_at = 0").OrderBy("sort").Find(&permissions)
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

func GetPermissionDetailSer(c *gin.Context, id int) {
	var permission model.Permission
	_, err := conf.Mysql.Where("id = ?", id).Get(&permission)
	if err != nil {
		common.ResError(c, "获取权限详情失败")
		return
	}
	common.ResOk(c, "ok", permission)
}

func GetAllRolesSer(c *gin.Context, page int, pageSize int) {
	var roles []*model.Role
	count, err := conf.Mysql.Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&roles)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: roles})
}

func GetRolePermissionSer(c *gin.Context, roleId int) {
	var role model.Role
	_, err := conf.Mysql.Where("id = ?", roleId).Where("deleted_at = 0").Get(&role)
	if err != nil {
		common.ResError(c, "获取角色信息失败")
		return
	}
	var rolePermissions []*model.RolePermission
	err = conf.Mysql.Where("role_id = ?", roleId).Where("deleted_at = 0").Find(&rolePermissions)
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
			Name:      roleAdd.Name,
			CreatedAt: int(time.Now().Unix()),
			UpdatedAt: int(time.Now().Unix()),
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
				CreatedAt:    int(time.Now().Unix()),
				UpdatedAt:    int(time.Now().Unix()),
			})
		}
		_, err = conf.Mysql.Insert(&rolePermissions)
		if err != nil {
			common.ResError(c, "添加角色权限关系失败")
			return
		}
	} else {
		var role model.Role
		_, err := conf.Mysql.Where("id = ?", roleAdd.Id).Where("deleted_at = 0").Get(&role)
		if err != nil {
			common.ResError(c, "获取角色信息失败")
			return
		}
		_, err = conf.Mysql.Where("id = ?", roleAdd.Id).Where("deleted_at = 0").Update(&model.Role{
			Name: roleAdd.Name,
		})
		if err != nil {
			common.ResError(c, "修改角色信息失败")
			return
		}
		_, err = conf.Mysql.Where("role_id = ?", roleAdd.Id).Where("deleted_at = 0").Update(&model.RolePermission{DeletedAt: int(time.Now().Unix())})
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
				CreatedAt:    int(time.Now().Unix()),
				UpdatedAt:    int(time.Now().Unix()),
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
	_, err := conf.Mysql.Where("id = ?", roleId).Where("deleted_at = 0").Update(&model.Role{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "删除角色失败")
		return
	}
	_, err = conf.Mysql.Where("role_id = ?", roleId).Where("deleted_at = 0").Update(&model.RolePermission{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "删除角色权限关系失败")
		return
	}
	_, err = conf.Mysql.Where("role_id = ?", roleId).Where("deleted_at = 0").Update(&model.UserRole{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "删除用户角色关系失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func AddPermissionSer(c *gin.Context, id int, parentId int, path string, icon string, sort int, label string, desc string, component string) {
	if id != 0 {
		_, err := conf.Mysql.Where("id = ?", id).Cols("path", "icon", "sort", "parent_id", "label", "desc", "component").Update(model.Permission{
			Path:      path,
			Icon:      icon,
			Sort:      sort,
			ParentId:  parentId,
			Label:     label,
			Desc:      desc,
			Component: component,
			UpdatedAt: int(time.Now().Unix()),
		})
		if err != nil {
			common.ResError(c, "修改权限失败")
			return
		}
	} else {
		var newPermission = model.Permission{
			Path:      path,
			Icon:      icon,
			Sort:      sort,
			ParentId:  parentId,
			Label:     label,
			Desc:      desc,
			Component: component,
			CreatedAt: int(time.Now().Unix()),
			UpdatedAt: int(time.Now().Unix()),
		}
		_, err := conf.Mysql.Insert(newPermission)
		if err != nil {
			common.ResError(c, "添加权限失败")
			return
		}
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
	_, err = conf.Mysql.In("id", permissionIds).Where("deleted_at = 0").Update(&model.Permission{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "节点删除失败")
		return
	}
	_, err = conf.Mysql.In("permission_id", permissionIds).Where("deleted_at = 0").Delete(&model.RolePermission{DeletedAt: int(time.Now().Unix())})
	if err != nil {
		common.ResError(c, "节点关联角色信息删除失败")
		return
	}

	common.ResOk(c, "ok", nil)
}

func SmsCodeSer(c *gin.Context, req model.SmsReq) {
	var existSms []*model.Sms
	err := conf.Mysql.Where("expired_at > ?", time.Now().Unix()).Find(&existSms)
	if err != nil {
		common.ResError(c, "获取短信信息失败")
		return
	}
	if len(existSms) > 0 {
		common.ResForbidden(c, "验证码有效时间5分钟，不可重复发送")
		return
	}
	smsCode := common.GetOneNewCard(4)
	data := "{\"code\":\"" + smsCode + "\"}"
	err = common.SendSms(req.Phone, conf.Conf.Sms.CodeTemplate, conf.Conf.Sms.CodeSign, data)
	if err != nil {
		common.ResError(c, "发送短信失败")
		return
	}
	_, err = conf.Mysql.Insert(model.Sms{
		Phone:     req.Phone,
		Code:      smsCode,
		UseAt:     0,
		ExpiredAt: int(time.Now().Unix()) + 300,
		CreateAt:  int(time.Now().Unix()),
	})
	if err != nil {
		common.ResError(c, "添加短信发送记录失败")
		return
	}
	common.ResOk(c, "ok", nil)
}
