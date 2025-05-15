package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/model"
	"management-backend/service"
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

func PermissionInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GetPermissionInfoSer(c, idInt)
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
		permissionReq.Id,
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

func AllOrganization(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AllOrganizationSer(c, pageInt, pageSizeInt)
}

func GetOrganizationInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GetOrganizationInfoSer(c, idInt)
}

func AddOrganization(c *gin.Context) {
	var addReq model.OrganizationAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddOrganizationSer(c, addReq)
}

func DeleteOrganization(c *gin.Context) {
	var deleteReq model.OrganizationDeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.DeleteOrganizationSer(c, deleteReq)
}
