package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"prize-draw/model"
	"prize-draw/service"
	"strconv"
)

// List 礼物列表
func List(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.List(c, name, pageInt, pageSizeInt)
}

// Info 礼物详情
func Info(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GiftInfo(c, idInt)
}

// Remain 礼物剩余情况
func Remain(c *gin.Context) {
	level := c.Query("level")
	exist := c.Query("exist")
	existInt, _ := strconv.Atoi(exist)
	service.GiftRemain(c, level, existInt)
}

// RaffleConfig 抽卡配置
func RaffleConfig(c *gin.Context) {
	service.RaffleConfig(c)
}

// PointLeft 抽卡点剩余情况
func PointLeft(c *gin.Context) {
	service.PointLeft(c)
}

// GroupList 礼物组列表
func GroupList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.GroupList(c, name, pageInt, pageSizeInt)
}

// GroupInfo 礼物组详情
func GroupInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.GroupInfo(c, idInt)
}

// AlbumList 相册列表
func AlbumList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AlbumList(c, name, pageInt, pageSizeInt)
}

// AlbumInfo 相册详情
func AlbumInfo(c *gin.Context) {
	id := c.Query("id")
	idInt, _ := strconv.Atoi(id)
	service.AlbumInfo(c, idInt)
}

// AlbumGift 获取未被相册绑定的礼物
func AlbumGift(c *gin.Context) {
	albumId := c.Query("album_id")
	albumIdInt, _ := strconv.Atoi(albumId)
	service.AlbumGift(c, albumIdInt)
}

// Add 添加礼物
func Add(c *gin.Context) {
	var addReq model.GiftAddReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.Add(c, addReq)
}

// AddPoint 添加抽卡点
func AddPoint(c *gin.Context) {
	var addReq model.GiftAddPointReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AddPoint(c, addReq)
}

// RaffleConfigSet 抽卡配置设置
func RaffleConfigSet(c *gin.Context) {
	var addReq model.GiftConfigSetReq
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.RaffleConfigSet(c, addReq)
}

// ResetPoint 一键重置
func ResetPoint(c *gin.Context) {
	service.ResetPoint(c)
}

// Delete 删除礼物(软)
func Delete(c *gin.Context) {
	var deleteReq model.GiftDeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.Delete(c, deleteReq)
}

// ChangeStatus 修改可用状态
func ChangeStatus(c *gin.Context) {
	var changeReq model.GiftChangeStatusReq
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		log.Fatal(err)
		return
	}
	service.ChangeReq(c, changeReq)
}

// GroupAdd 添加礼物组
func GroupAdd(c *gin.Context) {
	var addReq model.GiftGroupAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupAdd(c, addReq)
}

// GroupDelete 删除礼物组(软)
func GroupDelete(c *gin.Context) {
	var deleteReq model.GiftGroupDelete
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupDelete(c, deleteReq)
}

// GroupChangeStatus 修改礼物组状态
func GroupChangeStatus(c *gin.Context) {
	var addReq model.GiftGroupAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.GroupAdd(c, addReq)
}

// AlbumAdd 添加相册
func AlbumAdd(c *gin.Context) {
	var addReq model.GiftAlbumAdd
	if err := c.ShouldBindJSON(&addReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AlbumAdd(c, addReq)
}

// AlbumDelete 删除相册(软)
func AlbumDelete(c *gin.Context) {
	var deleteReq model.GiftAlbumDelete
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		log.Fatal(err)
		return
	}
	service.AlbumDelete(c, deleteReq)
}

/*=====================================app=====================================*/

func AppAlbumItems(c *gin.Context) {
	service.AppAlbumItems(c)
}

func AppAlbumList(c *gin.Context) {
	name := c.Query("name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AppAlbumList(c, name, pageInt, pageSizeInt)
}

func AppRealList(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	service.AppRealList(c, pageInt, pageSizeInt)
}
