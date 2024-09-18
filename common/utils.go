package common

import (
	"crypto/rand"
	"math/big"
)

func GetOneNewCard(length int) string {
	// 创建一个大整数，设置其为length位的随机数
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(10)).Exp(big.NewInt(10), big.NewInt(int64(length)), nil))
	if err != nil {
		return ""
	}

	// 将大整数转换为字符串
	return randomInt.String()
}
