package common

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

func DecimalToHex(decimal int) string {
	hex := fmt.Sprintf("%x", decimal)
	if len(hex) < 4 {
		hex = strings.Repeat("0", 4-len(hex)) + hex
	}
	return hex
}

func GetOneNewCard(length int) string {
	// 创建一个字节切片用于存储随机数据
	byteSlice := make([]byte, length)

	// 使用加密安全的随机数填充字节切片
	rand.Read(byteSlice)

	// 使用base64编码将字节切片转换为字符串
	return base64.URLEncoding.EncodeToString(byteSlice)
}
