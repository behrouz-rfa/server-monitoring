// +build !nocgo

package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func newPacketSource(device string, bpf string) (*gopacket.PacketSource, error) {
	if handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever); err != nil {
		return nil, err
	} else if err := handle.SetBPFFilter(bpf); err != nil { // optional
		return nil, err
	} else {
		return gopacket.NewPacketSource(handle, handle.LinkType()), nil
	}
}
