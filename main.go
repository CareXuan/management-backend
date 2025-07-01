package main

import (
	"env-backend/conf"
	"env-backend/router"
	"github.com/gin-gonic/gin"
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
