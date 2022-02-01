package services

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"reflect"
	capture "server-monitoring/services/caputure"
	"server-monitoring/services/output"
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
	for _, device := range devices {

		for _, address := range device.Addresses {
			ip := fmt.Sprintf("%s", address.IP)
			if reflect.DeepEqual(ip, "192.168.1.42") {
				devuceName = device.Name
				break
			}

		}
	}

	fmt.Println(devuceName)
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
	output.Init()
	capture.Start("ens33", filter)
	//httass.Run(devuceName)
	//http_log.Run("",devuceName,"","")
}
