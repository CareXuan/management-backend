package monitor

import (
	"github.com/gin-gonic/gin"
	"switchboard-backend/service/monitor"
)

func Collect(c *gin.Context) {
	monitor.Collect(c)
}
