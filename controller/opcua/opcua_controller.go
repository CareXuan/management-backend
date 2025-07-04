package opcua

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	opcua2 "switchboard-backend/model/opcua"
	"switchboard-backend/service/opcua"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	name := c.Query("name")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	opcua.ListSer(c, name, pageInt, pageSizeInt)
}

func Info(c *gin.Context) {
	id := c.Query("id")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	idInt, _ := strconv.Atoi(id)
	opcua.InfoSer(c, idInt, pageInt, pageSizeInt)
}

func Add(c *gin.Context) {
	var addReq opcua2.AddOpcuaReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	opcua.AddSer(c, addReq)
}
