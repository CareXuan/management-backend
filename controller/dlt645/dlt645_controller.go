package dlt645

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"switchboard-backend/model/dlt645"
	dlt6452 "switchboard-backend/service/dlt645"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	name := c.Query("name")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	dlt6452.ListSer(c, name, pageInt, pageSizeInt)
}

func Info(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	dlt6452.InfoSer(c, idInt)
}

func Add(c *gin.Context) {
	var addReq dlt645.AddDlt645Req
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	dlt6452.AddSer(c, addReq)
}
