package router

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/controller"
	"switchboard-backend/controller/data_sync"
	dlt6452 "switchboard-backend/controller/dlt645"
	firewall2 "switchboard-backend/controller/firewall"
	modbus2 "switchboard-backend/controller/modbus"
	opcua2 "switchboard-backend/controller/opcua"
	port2 "switchboard-backend/controller/port"
	siemens2 "switchboard-backend/controller/siemens"
	"switchboard-backend/middleware"
	monitor2 "switchboard-backend/service/monitor"
)

func InitRouter(r *gin.Engine) {
	r.Use(middleware.Cors())
	v1 := r.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) { common.ResOk(c, "Hello World!", nil) })

		port := v1.Group("port")
		{
			port.GET("/network/list", port2.List)
			port.GET("/network/info", port2.Info)
			port.GET("/bridge/list", port2.BridgeList)
			port.GET("/bridge/info", port2.BridgeInfo)
			port.POST("/network/change", port2.Change)
			port.POST("/bridge/add", port2.BridgeAdd)
		}

		siemens := v1.Group("siemens")
		{
			s7 := siemens.Group("s7")
			{
				s7.GET("/list", siemens2.List)
				s7.GET("/info", siemens2.Info)
				s7.POST("/add", siemens2.Add)
				s7.POST("/add_data", siemens2.AddData)
			}
		}

		modbus := v1.Group("modbus")
		{
			modbus.GET("list", modbus2.List)
			modbus.GET("info", modbus2.Info)
			modbus.POST("add", modbus2.Add)
		}

		opcua := v1.Group("opcua")
		{
			opcua.GET("/list", opcua2.List)
			opcua.GET("/info", opcua2.Info)
			opcua.POST("/add", opcua2.Add)
			opcua.POST("/add_data", opcua2.AddData)
		}

		dlt645 := v1.Group("dlt645")
		{
			dlt645.GET("list", dlt6452.List)
			dlt645.GET("info", dlt6452.Info)
			dlt645.POST("add", dlt6452.Add)
		}

		firewall := v1.Group("firewall")
		{
			firewall.GET("/list", firewall2.List)
			firewall.POST("/add", firewall2.Add)
			firewall.POST("/delete", firewall2.Delete)
		}

		monitor := v1.Group("monitor")
		{
			monitor.GET("collect", monitor2.Collect)
		}

		dataSync := v1.Group("data_sync")
		{
			dataSync.GET("/info", data_sync.ConfigInfo)
			dataSync.POST("/add", data_sync.ConfigUpdate)
		}

		mqtt := v1.Group("/mqtt")
		{
			mqtt.GET("/list")
			mqtt.POST("/add")
		}

		// 用户相关
		user := v1.Group("user")
		{
			user.GET("/info", controller.GetUserInfo)
			user.GET("/list", controller.GetUserList)
			user.POST("/add", controller.AddUser)
			user.GET("/permission", controller.GetUserPermission)
		}

		// 权限
		auth := v1.Group("auth")
		{
			auth.POST("/login", controller.Login)
			auth.GET("/permission", controller.AllPermission)
			auth.GET("/permission/info", controller.PermissionInfo)
			auth.POST("/permission/add", controller.AddPermission)
			auth.POST("/permission/delete", controller.RemovePermission)
			auth.GET("/roles", controller.AllRoles)
			auth.GET("/roles/info", controller.GetRoleInfo)
			auth.POST("/roles/add", controller.AddRole)
			auth.POST("/roles/delete", controller.DeleteRole)
		}

		commonCtr := v1.Group("common")
		{
			commonCtr.POST("/upload", controller.Upload)
		}
	}
}
