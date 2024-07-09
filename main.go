package main

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/conf"
	"my-gpt-server/router"
	"net/http"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	conf.NewConfig("./conf/config.yaml")
	err := http.ListenAndServe(":8233", r)
	if err != nil {
		return
	}
}
