// +build nocgo

package capture

import (
	"server-monitoring/services/caputure/afpacket"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func newPacketSource(device string, bpf string) (*gopacket.PacketSource, error) {
	handle, err := afpacket.NewTPacket(device)
	if err != nil {
		return nil, err
	}

	if err := handle.SetBPFFilter(bpf); err != nil {
		return nil, err
	}
	return gopacket.NewPacketSource(handle, layers.LinkTypeEthernet),nil

}
