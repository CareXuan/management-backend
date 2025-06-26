package modbus

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	modbus2 "switchboard-backend/model/modbus"
	"switchboard-backend/service/modbus"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	name := c.Query("name")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	modbus.ListSer(c, name, pageInt, pageSizeInt)
}

func Info(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	modbus.InfoSer(c, idInt)
}

func Add(c *gin.Context) {
	var addReq modbus2.AddModbusReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	modbus.AddSer(c, addReq)
}
