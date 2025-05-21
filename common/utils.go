package common

import (
	"crypto/rand"
	"encoding/base64"
	"math"
	"strconv"
	"strings"
)

const (
	a  = 6378245.0              // 长半轴
	ee = 0.00669342162296594323 // 扁率
)

func GetOneNewCard(length int) string {
	// 创建一个字节切片用于存储随机数据
	byteSlice := make([]byte, length)

	// 使用加密安全的随机数填充字节切片
	rand.Read(byteSlice)

	// 使用base64编码将字节切片转换为字符串
	return base64.URLEncoding.EncodeToString(byteSlice)
}

// WGS84ToGCJ02 WGS84 转 GCJ02（火星坐标）
func WGS84ToGCJ02(wgLat, wgLon float64) (cgLat, cgLon float64) {
	if OutOfChina(wgLat, wgLon) {
		return wgLat, wgLon
	}
	dLat := transformLat(wgLon-105.0, wgLat-35.0)
	dLon := transformLon(wgLon-105.0, wgLat-35.0)
	radLat := wgLat / 180.0 * math.Pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * math.Pi)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * math.Pi)
	cgLat = wgLat + dLat
	cgLon = wgLon + dLon
	return cgLat, cgLon
}

// 判断是否在中国境内
func OutOfChina(lat, lon float64) bool {
	return lon < 73.66 || lon > 135.05 || lat < 3.86 || lat > 53.55
}

// 转换纬度偏移
func transformLat(x, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*math.Pi) + 40.0*math.Sin(y/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*math.Pi) + 320*math.Sin(y*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

// 转换经度偏移
func transformLon(x, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*math.Pi) + 20.0*math.Sin(2.0*x*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*math.Pi) + 40.0*math.Sin(x/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*math.Pi) + 300.0*math.Sin(x/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

// ConvertBDSLatitude 将北斗纬度格式 dd.mmmmm 转换为十进制度
// 输入示例：39.12345（北纬39度12.345分）或 "-39.12345"（南纬39度12.345分）
func ConvertBDSLatitude(latStr string) float64 {
	// 处理符号
	sign := 1.0
	if strings.HasPrefix(latStr, "-") {
		sign = -1.0
		latStr = latStr[1:]
	} else if strings.HasPrefix(latStr, "+") {
		latStr = latStr[1:]
	}

	// 分割度分部分
	parts := strings.Split(latStr, ".")
	if len(parts) != 2 {
		return 0
	}

	// 解析度部分
	degrees, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	minutesStr := parts[1][:2] + "." + parts[1][2:] // 转换为 XX.XXX 格式
	minutes, err := strconv.ParseFloat(minutesStr, 64)
	if err != nil {
		return 0
	}

	// 计算十进制度
	total := float64(degrees) + minutes/60
	return sign * total
}
