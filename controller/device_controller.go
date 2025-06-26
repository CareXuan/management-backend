package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"switchboard-backend/model"
	"switchboard-backend/service"
)

func List(c *gin.Context) {
	service.ListSer(c)
}

func Change(c *gin.Context) {
	var changeReq model.ChangeDeviceReq
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeSer(c, changeReq)
}
