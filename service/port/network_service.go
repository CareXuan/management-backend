package port

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/port"
)

var NetworkFileMapping = map[string]string{
	"fe1":  "/etc/systemd/network/40-fe1.network",
	"fe2":  "/etc/systemd/network/41-fe2.network",
	"fe3":  "/etc/systemd/network/42-fe3.network",
	"fe4":  "/etc/systemd/network/43-fe4.network",
	"fe5":  "/etc/systemd/network/44-fe5.network",
	"fe6":  "/etc/systemd/network/45-fe6.network",
	"fe7":  "/etc/systemd/network/46-fe7.network",
	"fe8":  "/etc/systemd/network/47-fe8.network",
	"gs9":  "/etc/systemd/network/48-gs9.network",
	"gs10": "/etc/systemd/network/49-gs10.network",
}

func ListSer(c *gin.Context) {
	var networks []*port.Network
	err := conf.Mysql.Where("deleted_at = 0").Find(&networks)
	if err != nil {
		common.ResError(c, "获取网口配置失败")
		return
	}
	var bridgeIds []int
	for _, i := range networks {
		bridgeIds = append(bridgeIds, i.BridgeId)
	}
	var bridgeMapping = make(map[int]*port.Bridge)
	var bridgeItems []*port.Bridge
	err = conf.Mysql.In("id", bridgeIds).Find(&bridgeItems)
	if err != nil {
		common.ResError(c, "获取网桥信息失败")
		return
	}
	for _, i := range bridgeItems {
		bridgeMapping[i.Id] = i
	}
	for _, i := range networks {
		i.Bridge = bridgeMapping[i.BridgeId]
	}
	common.ResOk(c, "ok", networks)
}

func InfoSer(c *gin.Context, id int) {
	var networkItem port.Network
	_, err := conf.Mysql.Where("id = ?", id).Get(&networkItem)
	if err != nil {
		common.ResError(c, "获取网口配置失败")
		return
	}
	if networkItem.BridgeId != 0 {
		var bridge port.Bridge
		_, err := conf.Mysql.Where("id = ?", networkItem.BridgeId).Get(&bridge)
		if err != nil {
			common.ResError(c, "获取网桥失败")
			return
		}
		networkItem.Bridge = &bridge
	}
	common.ResOk(c, "ok", networkItem)
}

func ChangeSer(c *gin.Context, req port.ChangeNetworkReq) {
	_, err := conf.Mysql.Where("port = ?", req.Port).Update(&port.Network{
		NetworkType: req.NetworkType,
		Network:     req.Network,
		BridgeId:    req.BridgeId,
	})
	if err != nil {
		common.ResError(c, "修改端口信息失败")
		return
	}
	network := req.Network
	if req.NetworkType == 2 {
		var bridgeItem port.Bridge
		_, err := conf.Mysql.Where("id = ?", req.BridgeId).Get(&bridgeItem)
		if err != nil {
			common.ResError(c, "获取网桥信息失败")
			return
		}
		network = bridgeItem.EnglishName
	}
	err = changeNetworkFile(req.Port, network, req.NetworkType)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}
	common.ResOk(c, "ok", nil)
}

func changeNetworkFile(port, network string, networkType int) error {
	var newContent string
	if networkType == 1 {
		newContent = "[Match]\nName=" + port + "\n\n[Network]\nAddress=" + network + "\n"
	} else {
		newContent = "[Match]\nName=" + port + "\n\n[Network]\nBridge=" + network + "\n"
	}
	// 实际应从请求中获取或校验
	targetFile := NetworkFileMapping[port]
	backupFile := targetFile + ".bak"

	// 1. 检查root权限
	if os.Geteuid() != 0 {
		return fmt.Errorf("需要root权限才能修改" + targetFile)
	}

	// 2. 检查原文件是否存在（若不存在则创建备份会失败）
	_, err := os.Stat(targetFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("检查原文件失败: %v", err)
	}

	// 3. 备份原文件（仅当原文件存在时）
	if !os.IsNotExist(err) {
		originalContent, err := os.ReadFile(targetFile)
		if err != nil {
			return fmt.Errorf("读取原文件失败: %v", err)
		}
		// 写入备份文件（权限0644）
		if err := os.WriteFile(backupFile, originalContent, 0644); err != nil {
			return fmt.Errorf("备份原文件失败: %v", err)
		}
		log.Printf("原文件已备份至 %s", backupFile)
	} else {
		log.Printf("原文件不存在，将创建新文件 %s", targetFile)
	}

	// 4. 创建临时文件（显式设置权限为0644）
	tmpDir := filepath.Dir(targetFile)
	tmpFileName := filepath.Join(tmpDir, fmt.Sprintf("%s-", filepath.Base(targetFile))) // 临时文件名前缀
	tmpFile, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)     // 关键：显式设置权限0644
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer func() {
		// 清理临时文件（仅在未成功替换时清理）
		if err != nil {
			_ = os.Remove(tmpFileName)
		}
	}()

	// 5. 写入新内容到临时文件
	if _, err := tmpFile.WriteString(newContent); err != nil {
		return fmt.Errorf("写入临时文件失败: %v", err)
	}

	// 6. 强制同步临时文件到磁盘（确保数据持久化）
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("同步临时文件失败: %v", err)
	}
	if err := tmpFile.Close(); err != nil { // 关闭文件以便后续重命名
		return fmt.Errorf("关闭临时文件失败: %v", err)
	}

	// 7. 原子替换原文件（继承临时文件的0644权限）
	if err := os.Rename(tmpFileName, targetFile); err != nil {
		return fmt.Errorf("原子替换失败: %v", err)
	}
	log.Printf("文件 %s 已成功替换为临时文件（权限0644）", targetFile)

	err = restartNetwork(port)
	if err != nil {
		return fmt.Errorf("重启服务失败")
	}

	return nil
}

func restartNetwork(network string) error {
	out, err := exec.Command("chown", "113", NetworkFileMapping[network]).Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	out, err = exec.Command("chgrp", "119", NetworkFileMapping[network]).Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	//out, err = exec.Command("systemctl", "restart", "systemd-networkd").Output()
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println(out)
	return nil
}
