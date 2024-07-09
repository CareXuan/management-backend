package controller

import (
	"github.com/gin-gonic/gin"
	"my-gpt-server/service"
	"strconv"
)

func GetOneAnswer(c *gin.Context) {
	question := c.Query("question")
	service.GetOneAnswer(c, question)
}

func QuestionList(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	page := c.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)
	pageSize := c.DefaultQuery("page_size", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.QuestionList(c, search, pageInt, pageSizeInt)
}

func QuestionDetail(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.QuestionDetail(c, idInt)
}
