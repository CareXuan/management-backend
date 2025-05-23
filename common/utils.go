package common

import (
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
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

func CompareDaysSinceTimestamp(timestamp int64, number int) bool {
	t := time.Unix(timestamp, 0)
	now := time.Now()

	// 去除时间部分，只保留日期
	startDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	days := int(endDate.Sub(startDate).Hours() / 24)
	if days < 0 {
		days = -days
	}

	return days >= number
}

func RandomClosedInterval(min, max int) int {
	// 参数校验与修正（支持倒序调用）
	if min > max {
		min, max = max, min
	}

	// 核心算法（确保包含两端）
	rangeSize := max - min + 1 // 计算范围长度（关键！）
	return rand.Intn(rangeSize) + min
}

func ResizeUploadImages() {
	dirPath := "./upload" // 要遍历的文件夹
	scale := 0.2          // 缩放比例，比如 0.2 = 缩小到20%

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("访问文件出错:", err)
			return err
		}
		if d.IsDir() {
			return nil // 如果是文件夹就跳过
		}

		info, err := d.Info()
		if err != nil {
			fmt.Println("读取文件信息失败:", err)
			return err
		}

		if info.Size() <= 2<<20 { // 2MB = 2*2^20
			return nil // 小于2MB的文件跳过
		}

		// 打开大文件
		fmt.Println("压缩大文件:", path)
		err = compressImage(path, scale)
		if err != nil {
			fmt.Println("压缩失败:", err)
		}
		return nil
	})

	if err != nil {
		fmt.Println("遍历文件夹出错:", err)
	}
}

func compressImage(filePath string, scale float64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("图片解码失败: %w", err)
	}

	newWidth := uint(float64(img.Bounds().Dx()) * scale)
	newHeight := uint(float64(img.Bounds().Dy()) * scale)

	if newWidth == 0 || newHeight == 0 {
		return fmt.Errorf("新尺寸太小，不进行压缩")
	}

	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	out, err := os.Create(filePath) // 直接覆盖原文件
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	if format == "jpeg" {
		err = jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 85})
	} else if format == "png" {
		err = png.Encode(out, resizedImg)
	} else {
		return fmt.Errorf("不支持的图片格式: %s", format)
	}

	return err
}
