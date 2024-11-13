package controller

import (
	"data_verify/model"
	"data_verify/service"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func GetConfig(c *gin.Context) {
	searchType := c.Query("type")
	searchTypeInt, _ := strconv.Atoi(searchType)
	service.GetConfigSer(c, searchTypeInt)
}

func SetConfig(c *gin.Context) {
	var setConfigReq model.SetConfigReq
	if err := c.ShouldBindJSON(&setConfigReq); err != nil {
		log.Fatal(err)
		return
	}
	service.SetConfigSer(c, setConfigReq)
}
