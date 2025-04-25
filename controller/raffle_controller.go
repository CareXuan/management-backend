package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"prize-draw/model"
	"prize-draw/service"
)

func RaffleOne(c *gin.Context) {
	var raffleReq model.RaffleOneReq
	if err := c.ShouldBindJSON(&raffleReq); err != nil {
		log.Fatal(err)
		return
	}
	service.RaffleOne(c, raffleReq)
}

func AppRaffleConfig(c *gin.Context) {
	service.AppRaffleConfig(c)
}
