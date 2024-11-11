package middleware

import (
	"data_verify/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LimitUploadSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		if err := c.Request.ParseMultipartForm(maxSize); err != nil {
			common.ResError(c, "文件过大")
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Next()
	}
}
