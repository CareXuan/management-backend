package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"switchboard-backend/conf"
	"switchboard-backend/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	conf.NewConfig("./conf/config.yaml")
	err := http.ListenAndServe(":8222", r)
	if err != nil {
		return
	}
}
