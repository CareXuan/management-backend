package siemens

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"switchboard-backend/common"
)

func ListSer(c *gin.Context, name string, page, pageSize int) {

}

func Tttt(c *gin.Context) {
	newContent := "[Match]\nName=fe5\n\n[Network]\nAddress=192.168.5.43/24\n" // 实际应从请求中获取或校验
	targetFile := "/etc/systemd/network/44-fe5.network"
	backupFile := targetFile + ".bak"

	// 1. 检查root权限
	if os.Geteuid() != 0 {
		common.ResError(c, "需要root权限才能修改"+targetFile)
		return
	}

	// 2. 检查原文件是否存在（若不存在则创建备份会失败）
	_, err := os.Stat(targetFile)
	if err != nil && !os.IsNotExist(err) {
		common.ResError(c, fmt.Sprintf("检查原文件失败: %v", err))
		return
	}

	// 3. 备份原文件（仅当原文件存在时）
	if !os.IsNotExist(err) {
		originalContent, err := os.ReadFile(targetFile)
		if err != nil {
			common.ResError(c, fmt.Sprintf("读取原文件失败: %v", err))
			return
		}
		// 写入备份文件（权限0644）
		if err := os.WriteFile(backupFile, originalContent, 0644); err != nil {
			common.ResError(c, fmt.Sprintf("备份原文件失败: %v", err))
			return
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
		common.ResError(c, fmt.Sprintf("创建临时文件失败: %v", err))
		return
	}
	defer func() {
		// 清理临时文件（仅在未成功替换时清理）
		if err != nil {
			_ = os.Remove(tmpFileName)
		}
	}()

	// 5. 写入新内容到临时文件
	if _, err := tmpFile.WriteString(newContent); err != nil {
		common.ResError(c, fmt.Sprintf("写入临时文件失败: %v", err))
		return
	}

	// 6. 强制同步临时文件到磁盘（确保数据持久化）
	if err := tmpFile.Sync(); err != nil {
		common.ResError(c, fmt.Sprintf("同步临时文件失败: %v", err))
		return
	}
	if err := tmpFile.Close(); err != nil { // 关闭文件以便后续重命名
		common.ResError(c, fmt.Sprintf("关闭临时文件失败: %v", err))
		return
	}

	// 7. 原子替换原文件（继承临时文件的0644权限）
	if err := os.Rename(tmpFileName, targetFile); err != nil {
		common.ResError(c, fmt.Sprintf("原子替换失败: %v", err))
		return
	}
	log.Printf("文件 %s 已成功替换为临时文件（权限0644）", targetFile)

	err = restartNetwork("fe1")
	if err != nil {
		common.ResError(c, "重启服务失败")
		return
	}

	common.ResOk(c, "文件修改并生效成功（权限0644）", nil)
}

func restartNetwork(network string) error {
	out, err := exec.Command("chown", "113", "/etc/systemd/network/44-fe5.network").Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	out, err = exec.Command("chgrp", "119", "/etc/systemd/network/44-fe5.network").Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	out, err = exec.Command("systemctl", "restart", "systemd-networkd").Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	return nil
}
