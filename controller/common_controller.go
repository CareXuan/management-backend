package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-gpt-server/common"
	"my-gpt-server/conf"
	"net/http"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ResError(c, "文件上传失败")
		return
	}
	fileMaxSize := 4 << 20 //4M
	if int(file.Size) > fileMaxSize {
		common.ResError(c, "文件不允许大小于32KB")
		return
	}
	reader, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	b := make([]byte, 512)
	reader.Read(b)
	contentType := http.DetectContentType(b)
	if contentType != "image/jpeg" && contentType != "image/png" {
		common.ResError(c, "只能上传jpeg/png")
		return
	}
	fileName := "YQ" + common.GetOneNewCard(18) + ".png"
	dst := "./upload/" + fileName
	c.SaveUploadedFile(file, dst)
	common.ResOk(c, "ok", conf.Conf.Upload.Url+"/"+fileName)
}
