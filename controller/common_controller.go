package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"prize-draw/common"
	"prize-draw/conf"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.ResError(c, "文件上传失败")
		return
	}

	fileMaxSize := 128 << 20 // 128MB
	if int(file.Size) > fileMaxSize {
		common.ResError(c, "文件不允许大于8MB")
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		fmt.Println(err)
		common.ResError(c, "打开文件失败")
		return
	}
	defer fileReader.Close()

	b := make([]byte, 512)
	fileReader.Read(b)
	contentType := http.DetectContentType(b)

	if contentType != "image/jpeg" && contentType != "image/png" {
		common.ResError(c, "只能上传jpeg/png")
		return
	}

	fileReader.Seek(0, 0) // 重置reader指针

	// 生成文件名
	fileName := "YQ" + common.GetOneNewCard(18) + ".png"
	dst := "./upload/" + fileName

	scale := 0.2 // 这里控制缩放比例，比如 0.2 = 缩成原图的20%大小，你可以设置成传参

	if file.Size > 2<<20 { // 超过2MB，进行缩放压缩
		img, format, err := image.Decode(fileReader)
		if err != nil {
			common.ResError(c, "图片解码失败")
			return
		}

		// 计算新的宽高
		newWidth := uint(float64(img.Bounds().Dx()) * scale)
		newHeight := uint(float64(img.Bounds().Dy()) * scale)

		// 使用 resize 库进行缩放
		resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

		out, err := os.Create(dst)
		if err != nil {
			common.ResError(c, "创建文件失败")
			return
		}
		defer out.Close()

		if format == "jpeg" {
			jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 85})
		} else if format == "png" {
			png.Encode(out, resizedImg)
		}
	} else {
		// 小于2MB，直接保存
		c.SaveUploadedFile(file, dst)
	}

	common.ResOk(c, "ok", conf.Conf.Upload.Url+"/"+fileName)
}

func WechatCheck(c *gin.Context) {
	common.ResOk(c, "ok", nil)
}

func ResizeUpload(c *gin.Context) {
	common.ResizeUploadImages()
}
