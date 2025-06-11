package siemens

import (
	"github.com/gin-gonic/gin"
	"strconv"
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

func Add(c *gin.Context) {

}

func Tttt(c *gin.Context) {
	siemens.Tttt(c)
}
