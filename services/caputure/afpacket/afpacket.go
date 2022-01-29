// +build linux

package afpacket

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"regexp"
	"runtime"
	"server-monitoring/services/caputure/utils"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type TPacket struct {
	// fd is the C file descriptor.
	fd             int
	interfaceIndex int
	mu             sync.Mutex // guards below
}

func NewTPacket(device string) (h *TPacket, err error) {
	h = &TPacket{}

	fd, err := unix.Socket(unix.AF_PACKET, int(unix.SOCK_RAW), int(htons(unix.ETH_P_ALL)))
	if err != nil {
		return nil, err
	}
	h.fd = fd

	if err = h.bindToInterface(device); err != nil {
		h.Close()
		return nil, err
	}

	runtime.SetFinalizer(h, (*TPacket).Close)
	return h, err
}

// bindToInterface binds the TPacket socket to a particular named interface.
func (h *TPacket) bindToInterface(ifaceName string) error {
	ifIndex := 0
	// An empty string / Any here means to listen to all interfaces
	if ifaceName != "" && ifaceName != "any" {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return fmt.Errorf("InterfaceByName: %v", err)
		}
		ifIndex = iface.Index
	}
	s := &unix.SockaddrLinklayer{
		Protocol: htons(uint16(unix.ETH_P_ALL)),
		Ifindex:  ifIndex,
	}
	h.interfaceIndex = ifIndex
	return unix.Bind(h.fd, s)
}

func (h *TPacket) SetBPFFilter(filter string) error {
	if filter == "" {
		return nil
	}

	// pcapBPF, err := pcap.CompileBPFFilter(layers.LinkTypeEthernet, 65535, filter)
	// if err != nil {
	// 	return err
	// }
	// bpfIns := []bpf.RawInstruction{}
	// for _, ins := range pcapBPF {
	// 	bpfIns2 := bpf.RawInstruction{
	// 		Op: ins.Code,
	// 		Jt: ins.Jt,
	// 		Jf: ins.Jf,
	// 		K:  ins.K,
	// 	}
	// 	bpfIns = append(bpfIns, bpfIns2)
	// }
	// fmt.Println(string(utils.ToJSON(bpfIns)))

	bpfIns, err := h.parseBPFFilter(filter)
	if err != nil {
		return err
	}
	// fmt.Println(string(utils.ToJSON(bpfIns)))

	return h.setBPF(bpfIns)
}

func (h *TPacket) parseBPFFilter(filter string) ([]bpf.RawInstruction, error) {
	bpfIns := []bpf.RawInstruction{}
	rePort := regexp.MustCompile(`port\s+?(\d+)`)
	reHost := regexp.MustCompile(`host\s+?(\d+\.\d+\.\d+\.\d+)`)
	if filter == "tcp" {
		data := []byte(`[
				{"Op":40,"Jt":0,"Jf":0,"K":12},
				{"Op":21,"Jt":0,"Jf":5,"K":34525},
				{"Op":48,"Jt":0,"Jf":0,"K":20},
				{"Op":21,"Jt":6,"Jf":0,"K":6},
				{"Op":21,"Jt":0,"Jf":6,"K":44},
				{"Op":48,"Jt":0,"Jf":0,"K":54},
				{"Op":21,"Jt":3,"Jf":4,"K":6},
				{"Op":21,"Jt":0,"Jf":3,"K":2048},
				{"Op":48,"Jt":0,"Jf":0,"K":23},
				{"Op":21,"Jt":0,"Jf":1,"K":6},
				{"Op":6,"Jt":0,"Jf":0,"K":65535},
				{"Op":6,"Jt":0,"Jf":0,"K":0}
			]`)
		if err := utils.FromJSON(data, &bpfIns); err != nil {
			fmt.Println(err)
		}
	} else if strings.Contains(filter, "port ") && strings.Contains(filter, "host ") {
		match := rePort.FindStringSubmatch(filter)
		port := match[1]
		match = reHost.FindStringSubmatch(filter)
		host := match[1]
		hostInt := utils.IP2int(host)
		data := []byte(fmt.Sprintf(`[
				{"Op":40,"Jt":0,"Jf":0,"K":12},
				{"Op":21,"Jt":15,"Jf":0,"K":34525},
				{"Op":21,"Jt":0,"Jf":14,"K":2048}
				,{"Op":48,"Jt":0,"Jf":0,"K":23},
				{"Op":21,"Jt":0,"Jf":12,"K":6},
				{"Op":40,"Jt":0,"Jf":0,"K":20},
				{"Op":69,"Jt":10,"Jf":0,"K":8191},
				{"Op":177,"Jt":0,"Jf":0,"K":14},
				{"Op":72,"Jt":0,"Jf":0,"K":14},
				{"Op":21,"Jt":2,"Jf":0,"K":%s},
				{"Op":72,"Jt":0,"Jf":0,"K":16},
				{"Op":21,"Jt":0,"Jf":5,"K":%s},
				{"Op":32,"Jt":0,"Jf":0,"K":26},
				{"Op":21,"Jt":2,"Jf":0,"K":%d},
				{"Op":32,"Jt":0,"Jf":0,"K":30},
				{"Op":21,"Jt":0,"Jf":1,"K":%d},
				{"Op":6,"Jt":0,"Jf":0,"K":65535},
				{"Op":6,"Jt":0,"Jf":0,"K":0}
			]`, port, port, hostInt, hostInt))
		if err := utils.FromJSON(data, &bpfIns); err != nil {
			fmt.Println(err)
		}
	} else if strings.Contains(filter, "port ") {
		match := rePort.FindStringSubmatch(filter)
		port := match[1]
		data := []byte(fmt.Sprintf(`[
				{"Op":40,"Jt":0,"Jf":0,"K":12},
				{"Op":21,"Jt":0,"Jf":6,"K":34525},
				{"Op":48,"Jt":0,"Jf":0,"K":20},
				{"Op":21,"Jt":0,"Jf":15,"K":6},
				{"Op":40,"Jt":0,"Jf":0,"K":54},
				{"Op":21,"Jt":12,"Jf":0,"K":%s},
				{"Op":40,"Jt":0,"Jf":0,"K":56},
				{"Op":21,"Jt":10,"Jf":11,"K":%s},
				{"Op":21,"Jt":0,"Jf":10,"K":2048},
				{"Op":48,"Jt":0,"Jf":0,"K":23},
				{"Op":21,"Jt":0,"Jf":8,"K":6},
				{"Op":40,"Jt":0,"Jf":0,"K":20},
				{"Op":69,"Jt":6,"Jf":0,"K":8191},
				{"Op":177,"Jt":0,"Jf":0,"K":14},
				{"Op":72,"Jt":0,"Jf":0,"K":14},
				{"Op":21,"Jt":2,"Jf":0,"K":%s},
				{"Op":72,"Jt":0,"Jf":0,"K":16},
				{"Op":21,"Jt":0,"Jf":1,"K":%s},
				{"Op":6,"Jt":0,"Jf":0,"K":65535},
				{"Op":6,"Jt":0,"Jf":0,"K":0}
			]`, port, port, port, port))
		if err := utils.FromJSON(data, &bpfIns); err != nil {
			fmt.Println(err)
		}
	} else if strings.Contains(filter, "host ") {
		match := reHost.FindStringSubmatch(filter)
		host := match[1]
		hostInt := utils.IP2int(host)
		data := []byte(fmt.Sprintf(`[
			{"Op":40,"Jt":0,"Jf":0,"K":12},
			{"Op":21,"Jt":8,"Jf":0,"K":34525},
			{"Op":21,"Jt":0,"Jf":7,"K":2048},
			{"Op":48,"Jt":0,"Jf":0,"K":23},
			{"Op":21,"Jt":0,"Jf":5,"K":6},
			{"Op":32,"Jt":0,"Jf":0,"K":26},
			{"Op":21,"Jt":2,"Jf":0,"K":%d},
			{"Op":32,"Jt":0,"Jf":0,"K":30},
			{"Op":21,"Jt":0,"Jf":1,"K":%d},
			{"Op":6,"Jt":0,"Jf":0,"K":65535},
			{"Op":6,"Jt":0,"Jf":0,"K":0}
			]`, hostInt, hostInt))
		if err := utils.FromJSON(data, &bpfIns); err != nil {
			fmt.Println(err)
		}
	} else {
		return nil, fmt.Errorf("bpf filter only support tcp/port/host: %s", filter)
	}

	return bpfIns, nil
}

// SetBPF attaches a BPF filter to the underlying socket
func (h *TPacket) setBPF(filter []bpf.RawInstruction) error {
	var p unix.SockFprog
	if len(filter) > int(^uint16(0)) {
		return errors.New("filter too large")
	}
	p.Len = uint16(len(filter))
	p.Filter = (*unix.SockFilter)(unsafe.Pointer(&filter[0]))

	return setsockopt(h.fd, unix.SOL_SOCKET, unix.SO_ATTACH_FILTER, unsafe.Pointer(&p), unix.SizeofSockFprog)
}

// ReadPacketData reads the next packet, copies it into a new buffer, and returns
// that buffer.  Since the buffer is allocated by ReadPacketData, it is safe for long-term
// use.  This implements gopacket.PacketDataSource.
func (h *TPacket) ReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	buf := make([]byte, 65536)
	n, _, err := unix.Recvfrom(h.fd, buf, 0)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	if n <= 0 {
		err = errors.New("no data to read")
		return
	}

	ci.Timestamp = time.Now()
	ci.CaptureLength = n
	ci.Length = n
	ci.InterfaceIndex = h.interfaceIndex

	data = make([]byte, n)
	copy(data, buf)
	return
}

// Close cleans up the TPacket.  It should not be used after the Close call.
func (h *TPacket) Close() {
	if h.fd == -1 {
		return // already closed.
	}

	unix.Close(h.fd)
	h.fd = -1
	runtime.SetFinalizer(h, nil)
}
