package port

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/port"
	"switchboard-backend/utils"
)

func BridgeListSer(c *gin.Context, name string, page, pageSize int) {
	var bridgeItems []*port.Bridge
	sess := conf.Mysql.NewSession()
	if name != "" {
		sess.Where("name like ?", "%"+name+"%")
	}
	count, err := sess.Limit(pageSize, (page-1)*pageSize).FindAndCount(&bridgeItems)
	if err != nil {
		common.ResError(c, "获取网桥列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: bridgeItems})
}

func BridgeInfoSer(c *gin.Context, id int) {
	var bridgeItem port.Bridge
	_, err := conf.Mysql.Where("id = ?", id).Get(&bridgeItem)
	if err != nil {
		common.ResError(c, "获取网桥信息失败")
		return
	}
	common.ResOk(c, "ok", bridgeItem)
}

func BridgeAddSer(c *gin.Context, req port.AddBridgeReq) {
	if req.Id != 0 {
		_, err := conf.Mysql.Where("id = ?", req.Id).Update(&port.Bridge{
			Name:        req.Name,
			Ip:          req.Ip,
			EnglishName: req.EnglishName,
			FileName:    "/etc/systemd/network/" + req.EnglishName + ".network",
		})
		if err != nil {
			common.ResError(c, "修改网桥失败")
			return
		}
	} else {
		_, err := conf.Mysql.Insert(&port.Bridge{
			Name:        req.Name,
			Ip:          req.Ip,
			EnglishName: req.EnglishName,
			FileName:    "/etc/systemd/network/" + req.EnglishName + ".network",
		})
		if err != nil {
			common.ResError(c, "新增网桥失败")
			return
		}
	}
	count, err := conf.Mysql.Where("deleted_at = 0").FindAndCount(&[]*port.Bridge{})
	if err != nil {
		common.ResError(c, "获取网桥数量失败")
		return
	}
	err = changeBridgeFile(req.EnglishName, req.Ip, fmt.Sprintf("6%03d", count))
	if err != nil {
		common.ResError(c, "创建网桥文件失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func changeBridgeFile(name, ip, pvid string) error {
	newContentDev := "[NetDev]\nName=" + name + "\nKind=bridge\n\n[Bridge]\nVLANFiltering=yes\nDefaultPVID=" + pvid + "\n" // 实际应从请求中获取或校验
	newContentWork := "[Match]\nName=" + name + "\n\n[Network]\nAddress=" + ip + "/24\n"                                   // 实际应从请求中获取或校验
	targetFileDev := "/etc/systemd/network/" + name + ".netdev"
	targetFileWork := "/etc/systemd/network/" + name + ".network"

	// netdev文件
	// 1. 检查root权限
	if os.Geteuid() != 0 {
		return fmt.Errorf("需要root权限才能修改" + targetFileDev)
	}

	// 2. 检查原文件是否存在（若不存在则创建备份会失败）
	_, err := os.Stat(targetFileDev)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("检查原文件失败: %v", err)
	}

	if !os.IsNotExist(err) {
		originalContent, err := os.ReadFile(targetFileDev)
		if err != nil {
			return fmt.Errorf("读取原文件失败: %v", err)
		}
		// 写入备份文件（权限0644）
		if err := os.WriteFile(targetFileDev+".bak", originalContent, 0644); err != nil {
			return fmt.Errorf("备份原文件失败: %v", err)
		}
		log.Printf("原文件已备份至 %s", targetFileDev+".bak")
	} else {
		log.Printf("原文件不存在，将创建新文件 %s", targetFileDev)
	}

	// 4. 创建临时文件（显式设置权限为0644）
	tmpDir := filepath.Dir(targetFileDev)
	tmpFileName := filepath.Join(tmpDir, fmt.Sprintf("%s-", filepath.Base(targetFileDev))) // 临时文件名前缀
	tmpFile, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)        // 关键：显式设置权限0644
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
	if _, err := tmpFile.WriteString(newContentDev); err != nil {
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
	if err := os.Rename(tmpFileName, targetFileDev); err != nil {
		return fmt.Errorf("原子替换失败: %v", err)
	}
	log.Printf("文件 %s 已成功替换为临时文件（权限0644）", targetFileDev)

	// network文件
	// 2. 检查原文件是否存在（若不存在则创建备份会失败）
	_, err = os.Stat(targetFileWork)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("检查原文件失败: %v", err)
	}

	if !os.IsNotExist(err) {
		originalContent, err := os.ReadFile(targetFileWork)
		if err != nil {
			return fmt.Errorf("读取原文件失败: %v", err)
		}
		// 写入备份文件（权限0644）
		if err := os.WriteFile(targetFileWork+".bak", originalContent, 0644); err != nil {
			return fmt.Errorf("备份原文件失败: %v", err)
		}
		log.Printf("原文件已备份至 %s", targetFileWork+".bak")
	} else {
		log.Printf("原文件不存在，将创建新文件 %s", targetFileWork)
	}

	// 4. 创建临时文件（显式设置权限为0644）
	tmpDir = filepath.Dir(targetFileWork)
	tmpFileName = filepath.Join(tmpDir, fmt.Sprintf("%s-", filepath.Base(targetFileWork))) // 临时文件名前缀
	tmpFile, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)         // 关键：显式设置权限0644
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
	if _, err := tmpFile.WriteString(newContentWork); err != nil {
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
	if err := os.Rename(tmpFileName, targetFileWork); err != nil {
		return fmt.Errorf("原子替换失败: %v", err)
	}
	log.Printf("文件 %s 已成功替换为临时文件（权限0644）", targetFileWork)

	return nil
}
