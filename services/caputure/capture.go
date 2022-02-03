package capture

import (
	"bytes"
	"context"
	"fmt"
	"server-monitoring/domain/requests"
	"server-monitoring/services/caputure/flags"
	"server-monitoring/services/output"
	"server-monitoring/services/protos/http"
	"server-monitoring/services/protos/tls"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
)

var filter = ""

type captureManager struct {
	connManager      sync.Map
	chRecvTimeoutMsg chan *tcpConnection
	Done             bool
}

func newCaptureManager() *captureManager {
	return &captureManager{
		chRecvTimeoutMsg: make(chan *tcpConnection, 100),
	}

}

func (t *captureManager) Run(devicenam, f string) {
	filter = f
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go t.checkConnectionTimeout(ctx)

	//input := reader.NewRAWInput(flags.Options.InterfaceName, flags.Options.Port)
	log.Info().Msgf("Listening [%s] with BPF filter: %s", flags.Options.InterfaceName, filter)
	packetSource, err := newPacketSource(devicenam, filter)

	if err != nil {
		log.Err(err).Msg("")
		return
	}
	for packet := range packetSource.Packets() {

		if packet.NetworkLayer() == nil {

			continue
		}
		srcAddr := packet.NetworkLayer().NetworkFlow().Src().String()
		dstAddr := packet.NetworkLayer().NetworkFlow().Dst().String()

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		tcpLayer := packet.Layer(layers.LayerTypeTCP) // It is 'nil' in case of packet is not TCP.
		udpLayer := packet.Layer(layers.LayerTypeUDP) // It is 'nil' in case of packet is not UDP.
		icmLayer := packet.Layer(layers.LayerTypeICMPv4)

		if tcpLayer != nil {
			// Get actual TCP data from this layer
			tcp, _ := tcpLayer.(*layers.TCP)

			// ban some common use port connection, ex: ftp/samba/telnet
			srcPort := int(tcp.SrcPort)
			dstPort := int(tcp.DstPort)
			if isBanPort(srcPort, dstPort) {
				// log.Trace().Msg("ban common use port stream")
				if dstPort == 21 {
					app := packet.ApplicationLayer()
					loginType := ""
					if app != nil {
						payload := app.Payload()
						dst := packet.NetworkLayer().NetworkFlow().Dst()

						if bytes.Contains(payload, []byte("USER")) {
							loginType = string(payload)
							fmt.Printf("%v:%d-> %s\n", dst, dstPort, string(payload))
						} else if bytes.Contains(app.Payload(), []byte("PASS")) {
							loginType = string(payload)
							fmt.Printf("%v:%d-> %s\n", dst, dstPort, string(payload))
						} else if bytes.Contains(app.Payload(), []byte("AUTH TLS")) {
							fmt.Printf("%v:%d-> %s\n", dst, dstPort, string(payload))
							loginType = string(payload)
						} else {
							loginType = ""
							fmt.Printf("%v:%d-> %s\n", dst, dstPort, string(payload))
						}

						// App is present
					}

					fmt.Println()

					request := requests.Request{
						SrcAddr:       srcAddr,
						SrcPort:       srcPort,
						DstAddr:       dstAddr,
						DstPort:       dstPort,
						ContentLength: 0,
						Url:           fmt.Sprintf("%v:%d", dstAddr, dstPort),
						UserAgent:     "0",
						Method:        "FTP",
						Ts:            time.Now(),
						Body:          []byte(loginType),
						Response:      []byte(""),
					}

					//all_flags := [9]string{"FIN", "SYN", "RST", "PSH", "ACK", "URG", "ECE", "CWR", "NS"}
					//all_flags_value := [9]bool{tcp.FIN, tcp.SYN, tcp.RST, tcp.PSH, tcp.ACK, tcp.URG, tcp.ECE, tcp.CWR, tcp.NS}
					fmt.Printf(" FIN: %v - SYN: %v  - RST: %v - PSH: %v -ACK: %v\n", tcp.FIN, tcp.SYN, tcp.RST, tcp.PSH, tcp.ACK)
					if tcp.RST && tcp.ACK {
						request.StatusCode = -1
						fmt.Printf("srcAddr: %v dstAddr: %v ,FTP connection closed by reset", srcAddr, dstAddr)
					} else if tcp.FIN && tcp.ACK {
						request.StatusCode = 0
						fmt.Printf("srcAddr: %v dstAddr: %v ,FTP connection closed by user", srcAddr, dstAddr)
					} else if tcp.PSH && tcp.ACK {
						request.StatusCode = 1
						fmt.Printf("srcAddr: %v dstAddr: %v ,FTP connection establish", srcAddr, dstAddr)
					}
					request.InsertConsoleLog()
					//for index, element := range all_flags_value {
					//	if element {
					//		if all_flags[index] ==  "FIN" {
					//			//fmt.Println("FTP closed")
					//		}
					//		//fmt.Print(all_flags[index], " \n")
					//	}
					//}
				}
				continue
			}

			// receive FIN to close connection
			if tcp.FIN {
				t.handleFin(srcAddr, srcPort, dstAddr, dstPort)
				continue
			}

			// filter SYN,FIN,ACK-only packets not have data inside and Keepalive hearbeat packets with no data inside
			if len(tcp.Payload) == 0 {
				continue
			}

			// filter Keepalive hearbeat packets with 1-byte segment on Windows
			if tcp.ACK && len(tcp.Payload) == 1 {
				continue
			}

			streamId := fmt.Sprintf("%s@%d", srcAddr, tcp.Ack)
			// log.Debug().Msgf("%-30s -> %-40s seq:%-15d ack:%-15v payload:%-10d id:%s", fmt.Sprintf("%s:%s", srcAddr, tcp.SrcPort), fmt.Sprintf("%s:%s", dstAddr, tcp.DstPort), tcp.Seq, tcp.Ack, len(tcp.Payload), streamId)

			// Trying to add packet to existing message or creating new message
			//
			// For TCP message unique id is Acknowledgment number (see tcp_packet.go)
			connId := t.getConnId(srcAddr, srcPort, dstAddr, dstPort)

			newConn := NewTcpConnection(connId, streamId, srcAddr, srcPort, dstAddr, dstPort)
			v, _ := t.connManager.LoadOrStore(connId, newConn)
			conn := v.(*tcpConnection)
			conn.RegisterTimeoutEvent(t.chRecvTimeoutMsg)

			// handle same tcp connection, multiple request: HTTP/1.1
			// TODO: HTTP2
			isRequest := conn.IsRequest(srcAddr, srcPort)
			if isRequest && !conn.IsSameRequest(streamId) {
				log.Trace().Msgf("same tcp connection, multi request. %s", connId)
				// handle old request
				t.handleStream(conn, tcp, isRequest)
				t.handleFin(srcAddr, srcPort, dstAddr, dstPort)

				// handle new request
				conn = NewTcpConnection(connId, streamId, srcAddr, srcPort, dstAddr, dstPort)
				conn.RegisterTimeoutEvent(t.chRecvTimeoutMsg)
				t.connManager.Store(connId, conn)
			}
			conn.MarkActive()
			t.handleStream(conn, tcp, isRequest)
		}

		if ipLayer != nil {

		}
		if udpLayer != nil {

		}
		if icmLayer != nil {

		}
		//handlePacket(packet)  // do something with each packet
	}

	//go CopyMulty(input)
}

func (t *captureManager) handleStream(conn *tcpConnection, packet *layers.TCP, isRequest bool) {
	if conn.State == stateClosed {
		return
	}

	// handle https clienthello
	if isRequest {
		if clientHello, ok := tls.Parse(packet.Payload); ok {
			output.Print(clientHello, conn.srcAddr, conn.srcPort, conn.dstAddr, conn.dstPort)
			conn.Close()
			return
		}

	}

	// handle http
	if isRequest {
		log.Trace().Msgf("handle http reqeust: %s", conn.ID)
		http.ParseRequest(conn.ReqID, packet.Payload, flags.Options.Body)
	} else {
		log.Trace().Msgf("handle http response: %s", conn.ID)
		result, completed := http.ParseReponse(conn.ReqID, packet.Payload, flags.Options.Body)
		if completed && result != nil {
			output.Print(result, conn.srcAddr, conn.srcPort, conn.dstAddr, conn.dstPort)
			http.Close(conn.ReqID)
			conn.Close()
			return
		}
	}

}

// sometime may not receive response FIN
func (t *captureManager) handleFin(srcAddr string, srcPort int, dstAddr string, dstPort int) {
	connId := t.getConnId(srcAddr, srcPort, dstAddr, dstPort)
	log.Trace().Msgf("handleFin. connId: %s", connId)
	v, loaded := t.connManager.Load(connId)
	if !loaded {
		return
	}

	conn := v.(*tcpConnection)
	if conn.State == stateClosed {
		t.connManager.Delete(conn.ID)
		return
	}

	// send whatever data we got so far as complete. This
	// is needed for the HTTP/1.0 without Content-Length situation.
	result := http.GetResult(conn.ReqID)
	if result != nil {
		output.Print(result, conn.srcAddr, conn.srcPort, conn.dstAddr, conn.dstPort)
	}
	http.Close(conn.ReqID)
	conn.Close()
	t.connManager.Delete(conn.ID)
}

func (t *captureManager) checkConnectionTimeout(ctx context.Context) {
	for {
		select {
		case conn := <-t.chRecvTimeoutMsg:
			log.Trace().Msgf("Recieve timeout: %s:%d -> %s:%d", conn.srcAddr, conn.srcPort, conn.dstAddr, conn.dstPort)
			t.handleFin(conn.srcAddr, conn.srcPort, conn.dstAddr, conn.dstPort)
			conn.Close()
			t.connManager.Delete(conn.ID)
		case <-ctx.Done():
			return
		}
	}
}

func (t *captureManager) getConnId(srcAddr string, srcPort int, dstAddr string, dstPort int) string {
	keys := []string{srcAddr, fmt.Sprintf("%d", srcPort), dstAddr, fmt.Sprintf("%d", dstPort)}
	sortKeys := sort.StringSlice(keys)
	sortKeys.Sort()
	return strings.Join(sortKeys, "_")
}

func Start(devicenam, filter string) {
	newCaptureManager().Run(devicenam, filter)
}

func getBPFFilter() string {
	//if flags.Options.Port > 0 && flags.Options.Ip != "" {
	//	return fmt.Sprintf("tcp and port %d and host %s", flags.Options.Port, flags.Options.Ip)
	//} else if flags.Options.Ip != "" {
	//	return fmt.Sprintf("tcp and host %s", flags.Options.Ip)
	//} else if flags.Options.Port > 0 {
	//	return fmt.Sprintf("tcp and port %d", flags.Options.Port)
	//} else {
	//	return "tcp"
	//}

	return filter
}

// filter common use port
// 21: ftp, 23: telnet, 25: smtp, 139: samba
func isBanPort(srcPort int, dstPort int) bool {
	//fmt.Printf("srcport: %d  destPor: %d\n", srcPort, dstPort)
	if srcPort < 200 && srcPort != 80 {
		return true
	}

	if dstPort < 200 && dstPort != 80 {
		return true
	}

	return false
}
