package port

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"switchboard-backend/model/port"
	port2 "switchboard-backend/service/port"
)

func BridgeList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	port2.BridgeListSer(c, name, pageInt, pageSizeInt)
}

func BridgeInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	port2.BridgeInfoSer(c, idInt)
}

func BridgeAdd(c *gin.Context) {
	var addReq port.AddBridgeReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	port2.BridgeAddSer(c, addReq)
}
