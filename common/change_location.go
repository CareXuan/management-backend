package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// 高德 API 配置
	AMAP_KEY               = "b319d8c98dd28e171fd006c2b36d956a" // 替换为你的 API Key
	AMAP_REGEocode_URL     = "https://restapi.amap.com/v3/geocode/regeo"
	AMAP_COORD_CONVERT_URL = "https://restapi.amap.com/v3/assistant/coordinate/convert"
)

// 定义高德 API 返回的结构体
type ReGeoCodeResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
	} `json:"regeocode"`
}

type CoordConvertResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Locations string `json:"locations"`
}

// 逆地理编码：坐标 -> 地址
func ReGeoCode(location string, radius int, extensions string) (string, error) {
	// 构造请求 URL
	apiURL := fmt.Sprintf("%s?key=%s&location=%s&radius=%d&extensions=%s",
		AMAP_REGEocode_URL, AMAP_KEY, url.QueryEscape(location), radius, extensions)

	// 发送 HTTP GET 请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	body, _ := ioutil.ReadAll(resp.Body)
	var result ReGeoCodeResponse
	json.Unmarshal(body, &result)

	// 检查 API 返回状态
	if result.Status != "1" {
		return "", fmt.Errorf("API 错误: %s", result.Info)
	}

	return result.Regeocode.FormattedAddress, nil
}

// 坐标系转换：WGS84/GPS -> 高德 GCJ-02
func ConvertCoordinate(location string, coordsys string) (string, error) {
	apiURL := fmt.Sprintf("%s?key=%s&coordsys=%s&locations=%s",
		AMAP_COORD_CONVERT_URL, AMAP_KEY, coordsys, url.QueryEscape(location))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("坐标转换请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result CoordConvertResponse
	json.Unmarshal(body, &result)

	if result.Status != "1" {
		return "", fmt.Errorf("坐标转换失败: %s", result.Info)
	}

	return result.Locations, nil
}

func main() {
	// 示例 1：直接逆地理编码（需确保坐标是 GCJ-02）
	location := "116.397428,39.90923" // 北京天安门（GCJ-02 坐标）
	address, err := ReGeoCode(location, 1000, "all")
	if err != nil {
		fmt.Println("逆地理编码失败:", err)
		return
	}
	fmt.Println("地址:", address) // 输出：北京市东城区东长安街

	// 示例 2：WGS84 坐标先转换再逆地理编码
	wgs84Location := "116.397124,39.916527" // WGS84 坐标
	gcj02Location, err := ConvertCoordinate(wgs84Location, "gps")
	if err != nil {
		fmt.Println("坐标转换失败:", err)
		return
	}
	address, _ = ReGeoCode(gcj02Location, 1000, "all")
	fmt.Println("转换后地址:", address)
}
