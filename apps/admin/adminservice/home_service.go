package adminservice

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"server-monitoring/domain/requests"

	"time"
)

var (
	HomeServices homeServicesInterface = &homeService{}
)

type homeServicesInterface interface {
	LoadRequests(int) ([]requests.Request, error)
	LoadRequestsFilter(int, string) ([]requests.Request, error)
	CpuInfo() (int, error)
	MemoryInfo() (mem.VirtualMemoryStat, error)
	HostInfo() (host.InfoStat, error)
	DiskInfo() (map[string]disk.UsageStat, error)
	NetInfo() (string, error)
}
type homeService struct {
}

func (h homeService) NetInfo() (string, error) {
	info, _ := net.IOCounters(true)
	for index, v := range info {
		fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	}

	return "", nil
}

func (h homeService) HostInfo() (host.InfoStat, error) {
	hInfo, _ := host.Info()
	fmt.Printf("host info:%v uptime:%v boottime:%v\n", hInfo, hInfo.Uptime, hInfo.BootTime)
	return *hInfo, nil
}

func (h homeService) DiskInfo() (map[string]disk.UsageStat, error) {
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get Partitions failed, err:%v\n", err)
		return nil, err
	}
	items := make(map[string]disk.UsageStat)
	for _, part := range parts {
		fmt.Printf("part:%v\n", part.String())
		diskInfo, _ := disk.Usage(part.Mountpoint)
		items[part.Device] = *diskInfo
		fmt.Printf("disk info:used:%v free:%v\n", diskInfo.UsedPercent, diskInfo.Free)
	}

	ioStat, _ := disk.IOCounters()
	for k, v := range ioStat {
		fmt.Printf("%v:%v\n", k, v)
	}
	return items, nil
}

func (h homeService) MemoryInfo() (mem.VirtualMemoryStat, error) {
	memInfo, _ := mem.VirtualMemory()
	return *memInfo, nil

}

func (h homeService) CpuInfo() (int, error) {
	percent, _ := cpu.Percent(time.Second, false)
	var usage float64
	for _, f := range percent {
		usage += f
	}
	return int(usage), nil

}

func (h homeService) LoadRequests(page int) ([]requests.Request, error) {
	var r requests.Request
	return r.Find(page)
}
func (h homeService) LoadRequestsFilter(page int, key string) ([]requests.Request, error) {
	var r requests.Request
	if len(key) == 0 {
		return r.Find(page)
	}
	return r.FinMultipleFilter(page, key)
}
