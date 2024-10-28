package common

import (
	"github.com/gin-gonic/gin"
	"management-backend/utils"
)

func ResOk(c *gin.Context, msg string, data interface{}) {
	commonResponse(c, utils.RESPONSE_OK, msg, data)
}

func ResAuthorization(c *gin.Context, msg string) {
	commonResponse(c, utils.RESPONSE_AUTHORIZATION, msg, nil)
}

func ResForbidden(c *gin.Context, msg string) {
	commonResponse(c, utils.RESPONSE_FORBIDDEN, msg, nil)
}

func ResNotFound(c *gin.Context, msg string) {
	commonResponse(c, utils.RESPONSE_NOT_FOUND, msg, nil)
}

func ResError(c *gin.Context, msg string) {
	commonResponse(c, utils.RESPONSE_ERROR, msg, nil)
}

func commonResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
