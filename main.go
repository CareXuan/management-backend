package main

import (
	"data_verify/conf"
	"data_verify/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	conf.NewConfig("./conf/config.yaml")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
