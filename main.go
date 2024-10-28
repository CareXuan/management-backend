package main

import (
	"github.com/gin-gonic/gin"
	"management-backend/conf"
	"management-backend/router"
	"net/http"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	conf.NewConfig("./conf/config.yaml")
	err := http.ListenAndServe(":8586", r)
	if err != nil {
		return
	}
}
