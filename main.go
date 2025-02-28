package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"prize-draw/conf"
	"prize-draw/router"
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
