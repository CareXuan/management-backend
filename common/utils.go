package common

import (
	"crypto/rand"
	"encoding/base64"
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

func GetNextDay(weekday int) time.Time {
	now := time.Now()
	today := now.Weekday()

	// 计算距离目标星期几的天数
	daysUntilTarget := (weekday - int(today) + 7) % 7
	targetDate := now.AddDate(0, 0, daysUntilTarget)

	// 设置时间为 23:59:59
	targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, targetDate.Location())

	return targetDate
}

func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
