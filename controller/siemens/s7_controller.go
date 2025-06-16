package siemens

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	siemens2 "switchboard-backend/model/siemens"
	"switchboard-backend/service/siemens"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	name := c.Query("name")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	siemens.ListSer(c, name, pageInt, pageSizeInt)
}

func Info(c *gin.Context) {
	id := c.Query("id")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	idInt, _ := strconv.Atoi(id)
	siemens.InfoSer(c, idInt, pageInt, pageSizeInt)
}

func Add(c *gin.Context) {
	var addReq siemens2.AddSiemensS7Req
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	siemens.AddSer(c, addReq)
}

func AddData(c *gin.Context) {
	var addReq siemens2.AddSiemensDataReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	siemens.AddSiemensDataSer(c, addReq)
}
