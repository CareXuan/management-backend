package data_sync

import (
	"github.com/gin-gonic/gin"
	"log"
	data_sync2 "switchboard-backend/model/data_sync"
	"switchboard-backend/service/data_sync"
)

func ConfigInfo(c *gin.Context) {
	data_sync.ConfigInfoSer(c)
}

func ConfigUpdate(c *gin.Context) {
	var addReq data_sync2.DataSyncAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	data_sync.UpdateConfigSer(c, addReq)
}
