package common

import (
	"crypto/rand"
	"data_verify/conf"
	"data_verify/model"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func GetOneNewCard(length int) string {
	// 创建一个字节切片用于存储随机数据
	byteSlice := make([]byte, length)

	// 使用加密安全的随机数填充字节切片
	rand.Read(byteSlice)

	// 使用base64编码将字节切片转换为字符串
	return base64.URLEncoding.EncodeToString(byteSlice)
}

func IgnoreError[T any](value T, _ error) T {
	return value
}

func GetCurrentYear() (string, error) {
	var yearConfig model.Config
	_, err := conf.Mysql.Where("type = ?", model.CONFIG_TYPE_YEAR).Get(&yearConfig)
	if err != nil {
		return "", err
	}
	return yearConfig.Value, nil
}

func GetUserIdByToken(c *gin.Context) (int, string, error) {
	authorization := c.GetHeader("Authorization")
	part := strings.Split(authorization, " ")
	var user model.User
	_, err := conf.Mysql.Where("token=?", part[1]).Get(&user)
	if err != nil {
		return 0, "", err
	}
	return user.Id, user.Name, nil
}

func SetHistory(step, status, userId int, bmddm, userName, remark, year string) error {
	_, err := conf.Mysql.Insert(model.History{
		Step:       step,
		Status:     status,
		UserId:     userId,
		Bmddm:      bmddm,
		UserName:   userName,
		Remark:     remark,
		Year:       year,
		CreateTime: int(time.Now().Unix()),
	})
	if err != nil {
		return err
	}
	return nil
}
