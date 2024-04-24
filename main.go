package main

import (
	"github.com/gin-gonic/gin"
	"management-backend/conf"
	"management-backend/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)
	conf.NewConfig("./conf/config.yaml")
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
