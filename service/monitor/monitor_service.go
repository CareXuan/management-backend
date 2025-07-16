package monitor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"switchboard-backend/common"
	"switchboard-backend/conf"
	"switchboard-backend/model"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type Metric struct {
	Timestamp string         `json:"timestamp"`
	CPUUsage  string         `json:"cpu_usage"`
	MemUsage  string         `json:"mem_usage"`
	DiskUsage string         `json:"disk_usage"`
	NetSentMB string         `json:"net_sent_mb"`
	NetRecvMB string         `json:"net_recv_mb"`
	Traffic   []*TrafficItem `json:"traffic"`
}

type TrafficItem struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	Errin       uint64 `json:"errin"`
	Errout      uint64 `json:"errout"`
	Dropin      uint64 `json:"dropin"`
	Dropout     uint64 `json:"dropout"`
	Fifoin      uint64 `json:"fifoin"`
	Fifoout     uint64 `json:"fifoout"`
}

func Collect(c *gin.Context) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	netStats, err := net.IOCounters(false)
	if err != nil || len(netStats) == 0 {
		common.ResError(c, err.Error())
		return
	}

	metric := &Metric{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		CPUUsage:  fmt.Sprintf("%.2f", cpuPercent[0]),
		MemUsage:  fmt.Sprintf("%.2f", vmStat.UsedPercent),
		DiskUsage: fmt.Sprintf("%.2f", diskStat.UsedPercent),
		NetSentMB: fmt.Sprintf("%.2f", float64(netStats[0].BytesSent)/1024/1024),
		NetRecvMB: fmt.Sprintf("%.2f", float64(netStats[0].BytesRecv)/1024/1024),
	}

	var trafficNameItems []*model.Traffic
	err = conf.Mysql.GroupBy("name").Find(&trafficNameItems)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	var trafficStatisticItem []*model.Traffic
	err = conf.Mysql.OrderBy("id DESC").Limit(len(trafficNameItems) * 2).Find(&trafficStatisticItem)
	if err != nil {
		common.ResError(c, err.Error())
		return
	}

	var trafficMapping = make(map[string]*TrafficItem)
	for _, i := range trafficStatisticItem {
		if _, ok := trafficMapping[i.Name]; ok {
			trafficMapping[i.Name].BytesSent -= i.BytesSent
			trafficMapping[i.Name].BytesRecv -= i.BytesRecv
			trafficMapping[i.Name].PacketsSent -= i.PacketsSent
			trafficMapping[i.Name].PacketsRecv -= i.PacketsRecv
			trafficMapping[i.Name].Errin -= i.Errin
			trafficMapping[i.Name].Errout -= i.Errout
			trafficMapping[i.Name].Dropin -= i.Dropin
			trafficMapping[i.Name].Dropout -= i.Dropout
			trafficMapping[i.Name].Fifoin -= i.Fifoin
			trafficMapping[i.Name].Fifoout -= i.Fifoout
		} else {
			trafficMapping[i.Name] = &TrafficItem{
				Name:        i.Name,
				BytesSent:   i.BytesSent,
				BytesRecv:   i.BytesRecv,
				PacketsSent: i.PacketsSent,
				PacketsRecv: i.PacketsRecv,
				Errin:       i.Errin,
				Errout:      i.Errout,
				Dropin:      i.Dropin,
				Dropout:     i.Dropout,
				Fifoin:      i.Fifoin,
				Fifoout:     i.Fifoout,
			}
		}
	}

	for _, i := range trafficMapping {
		metric.Traffic = append(metric.Traffic, i)
	}

	common.ResOk(c, "ok", metric)
}
