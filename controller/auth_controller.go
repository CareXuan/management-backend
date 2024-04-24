package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"management-backend/service"
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
	fmt.Println(loginReq)
	service.LoginSrv(c, loginReq.Phone, loginReq.Password)
}
