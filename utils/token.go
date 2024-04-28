package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomToken(length int) (string, error) {
	// 计算生成的字节数
	numBytes := length * 3 / 4 // base64 编码每 4 个字符表示 3 个字节

	// 生成随机的字节序列
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 使用 base64 进行编码
	token := base64.URLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的 token 字符串
	return token[:length], nil
}
