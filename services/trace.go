package services

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"runtime"
	"server-monitoring/apps/admin/adminservice"
	"server-monitoring/domain/settings"
	capture "server-monitoring/services/caputure"
	"server-monitoring/services/output"
	"strings"
)

var (
	idace    = "Wi-Fi"
	snaplan  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	filter   = "ip and tcp and port 80"
	devFound = false
)

func Run() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	//IPAddress := net.ParseIP("192.168.1.42")
	devuceName := ""
	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows")
		if len(devices) > 0 {
			devuceName = devices[0].Addresses[0].IP.String()
		}
	} else {
		if len(devices) > 0 {
			for _, device := range devices {
				if strings.Contains(device.Name, "ens") || strings.Contains(device.Name, "eth") {
					devuceName = device.Name
					break
				}
			}
		}
	}

	//
	//handle ,err := pcap.OpenLive(devuceName,snaplan,promisc,timeout)
	//if err != nil {
	//	log.Println(err)
	//}
	//defer handle.Close()
	//if err := handle.SetBPFFilter(filter); err != nil {
	//	log.Println(err)
	//}
	//
	//sourec := gopacket.NewPacketSource(handle,handle.LinkType())
	//
	//for packet := range sourec.Packets() {
	//	fmt.Println(packet.String())
	//}
	//packet.Run()
	var setting settings.Setting
	adminservice.SettingService.Get(&setting)
	if len(setting.Interface) == 0 {
		setting.Interface = devuceName
	}
	output.Init()
	capture.Start(setting.Interface, setting.Filter)
	//httass.Run(devuceName)
	//http_log.Run("",devuceName,"","")
}
