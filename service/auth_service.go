package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/common"
	"management-backend/conf"
	"management-backend/model"
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
