package port

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"switchboard-backend/model/port"
	port2 "switchboard-backend/service/port"
)

func List(c *gin.Context) {
	port2.ListSer(c)
}

func Info(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	port2.InfoSer(c, idInt)
}

func Change(c *gin.Context) {
	var changeReq port.ChangeNetworkReq
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		log.Fatal(err)
		return
	}
	port2.ChangeSer(c, changeReq)
}
