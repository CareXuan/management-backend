package controller

import (
	"data_verify/model"
	"data_verify/service"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type LoginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var loginReq LoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		log.Fatal(err)
		return
	}
	service.LoginSrv(c, loginReq.Phone, loginReq.Password)
}

func AllPermission(c *gin.Context) {
	service.GetAllPermissionSer(c)
}

func AllRoles(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GetAllRolesSer(c, pageInt, pageSizeInt)
}

func GetRoleInfo(c *gin.Context) {
	roleId := c.Query("role_id")
	roleIdInt, _ := strconv.Atoi(roleId)
	service.GetRolePermissionSer(c, roleIdInt)
}

func AddRole(c *gin.Context) {
	var roleAddReq model.RoleAddReq
	if err := c.ShouldBindJSON(&roleAddReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddRoleSer(c, roleAddReq)
}

func DeleteRole(c *gin.Context) {
	var roleDeleteReq model.RoleAddReq
	if err := c.ShouldBindJSON(&roleDeleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteRoleSer(c, roleDeleteReq.Id)
}

func AddPermission(c *gin.Context) {
	var permissionReq model.Permission
	if err := c.ShouldBindJSON(&permissionReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddPermissionSer(
		c,
		permissionReq.ParentId,
		permissionReq.Path,
		permissionReq.Icon,
		permissionReq.Sort,
		permissionReq.Label,
		permissionReq.Desc,
		permissionReq.Component,
	)
}

func RemovePermission(c *gin.Context) {
	var permissionReq model.Permission
	if err := c.ShouldBindJSON(&permissionReq); err != nil {
		log.Fatal(err)
		return
	}
	service.RemovePermissionSer(c, permissionReq.Id)
}

func AllGroups(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AllGroupSer(c, name, pageInt, pageSizeInt)
}

func GroupInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GroupInfoSer(c, idInt)
}

func AddGroup(c *gin.Context) {
	var groupAddReq model.GroupAddReq
	if err := c.ShouldBindJSON(&groupAddReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddGroupSer(c, groupAddReq)
}

func DeleteGroup(c *gin.Context) {
	var groupDeleteReq model.GroupDeleteReq
	if err := c.ShouldBindJSON(&groupDeleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteGroupSer(c, groupDeleteReq)
}
