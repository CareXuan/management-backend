package firewall

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model/firewall"
	"switchboard-backend/utils"
)

func ListSer(c *gin.Context, ip string, page, pageSize int) {
	var List []*firewall.Firewall
	sess := conf.Mysql.NewSession()
	if ip != "" {
		sess.Where("ip like ?", "%"+ip+"%")
	}
	count, err := sess.Where("deleted_at = 0").Limit(pageSize, (page-1)*pageSize).FindAndCount(&List)
	if err != nil {
		common.ResError(c, "获取防火墙配置列表失败")
		return
	}
	common.ResOk(c, "ok", utils.CommonListRes{Count: count, Data: List})
}

func AddSer(c *gin.Context, req firewall.AddFirewallReq) {
	_, err := conf.Mysql.Insert(&firewall.Firewall{
		Ip:        req.Ip,
		Type:      req.Type,
		AllowType: req.AllowType,
	})
	if err != nil {
		common.ResError(c, "添加防火墙配置失败")
		return
	}
	err = updateFirewallConfig(req.Ip, 1, req.Type, req.AllowType)
	if err != nil {
		common.ResError(c, "执行系统命令失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func DeleteSer(c *gin.Context, req firewall.DeleteFirewallReq) {
	var firewallItem firewall.Firewall
	_, err := conf.Mysql.Where("id = ?", req.Id).Get(&firewallItem)
	if err != nil {
		common.ResError(c, "获取防火墙配置失败")
		return
	}
	_, err = conf.Mysql.Where("id = ?", req.Id).Delete(&firewall.Firewall{})
	if err != nil {
		common.ResError(c, "删除防火墙配置失败")
		return
	}
	err = updateFirewallConfig(firewallItem.Ip, 2, firewallItem.Type, firewallItem.AllowType)
	if err != nil {
		common.ResError(c, "执行系统命令失败")
		return
	}
	common.ResOk(c, "ok", nil)
}

func updateFirewallConfig(ip string, commandType, directionType, allowType int) error {
	command := "-A"
	direction := "INPUT"
	allow := "ACCEPT"
	if commandType == 2 {
		command = "-D"
	}
	if directionType == 2 {
		direction = "OUTPUT"
	}
	if allowType == 2 {
		allow = "DROP"
	}

	out, err := exec.Command("iptables", command, direction, "-s", ip, "-j", allow).Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(out)
	return nil
}
