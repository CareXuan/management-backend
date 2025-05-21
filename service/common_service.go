package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"management-backend/conf"
	"management-backend/model"
	"strings"
)

func IsSupervisor(c *gin.Context) (int, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return 0, fmt.Errorf("获取token失败")
	}
	part := strings.Split(authorization, " ")
	if len(part) < 2 {
		return 0, fmt.Errorf("非法请求")
	}
	var user model.User
	_, err := conf.Mysql.Where("token=?", part[1]).Get(&user)
	if err != nil {
		return 0, fmt.Errorf("获取用户信息失败")
	}
	var userRole model.UserRole
	_, err = conf.Mysql.Where("user_id = ?", user.Id).Get(&userRole)
	if err != nil {
		return 0, fmt.Errorf("获取用户角色失败")
	}
	var roleItem model.Role
	_, err = conf.Mysql.Where("id = ?", userRole.RoleId).Get(&roleItem)
	if err != nil {
		return 0, fmt.Errorf("获取角色信息失败")
	}
	if roleItem.Name != "管理员" {
		return 0, nil
	}
	return 1, nil
}
