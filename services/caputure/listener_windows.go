// +build !nocgo

package capture

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func newPacketSource(device string, bpf string) (*gopacket.PacketSource, error) {
	if device == "" || device == "any" {
		return nil, fmt.Errorf("Windows not support listen to all interface. Please use (-i) to specific an interface")
	}

	if ip := net.ParseIP(device); ip == nil {
		return nil, fmt.Errorf("Windows only support ip. Please use (-i) to specific an ip")
	}

	deviceId, err := getDeviceId(device)
	if err != nil {
		return nil, err
	}

	if handle, err := pcap.OpenLive(deviceId, 1600, true, pcap.BlockForever); err != nil {
		return nil, err
	} else if err := handle.SetBPFFilter(bpf); err != nil { // optional
		return nil, err
	} else {
		return gopacket.NewPacketSource(handle, handle.LinkType()), nil
	}

}

func getDeviceId(device string) (string, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", err
	}

	for _, dev := range devices {
		addrs := dev.Addresses
		for _, addr := range addrs {
			if addr.IP.String() == device {
				return dev.Name, nil
			}
		}
	}

	return "", fmt.Errorf("Not found device: %s", device)
}
