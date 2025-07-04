package firewall

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	firewall2 "switchboard-backend/model/firewall"
	"switchboard-backend/service/firewall"
)

func List(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	ip := c.Query("ip")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	firewall.ListSer(c, ip, pageInt, pageSizeInt)
}

func Add(c *gin.Context) {
	var addReq firewall2.AddFirewallReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	firewall.AddSer(c, addReq)
}

func Delete(c *gin.Context) {
	var deleteReq firewall2.DeleteFirewallReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	firewall.DeleteSer(c, deleteReq)
}
