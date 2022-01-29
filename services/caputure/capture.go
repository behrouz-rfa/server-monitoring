package capture

import (
	"context"
	"fmt"
	"server-monitoring/services/output"

	"server-monitoring/services/caputure/flags"
	"server-monitoring/services/protos/http"
	"server-monitoring/services/protos/tls"
	"sort"
	"strings"
	"sync"


	"github.com/google/gopacket/layers"
	"github.com/rs/zerolog/log"
)

type captureManager struct {
	connManager      sync.Map
	chRecvTimeoutMsg chan *tcpConnection
}

func newCaptureManager() *captureManager {
	return &captureManager{
		chRecvTimeoutMsg: make(chan *tcpConnection, 100),
	}
}

func (t *captureManager) Run(devicenam,filter string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go t.checkConnectionTimeout(ctx)

	//input := reader.NewRAWInput(flags.Options.InterfaceName, flags.Options.Port)
	log.Info().Msgf("Listening [%s] with BPF filter: %s", flags.Options.InterfaceName, getBPFFilter())
	packetSource, err := newPacketSource(devicenam, getBPFFilter())
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
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			// Get actual TCP data from this layer
			tcp, _ := tcpLayer.(*layers.TCP)

			// ban some common use port connection, ex: ftp/samba/telnet
			srcPort := int(tcp.SrcPort)
			dstPort := int(tcp.DstPort)
			if isBanPort(srcPort, dstPort) {
				// log.Trace().Msg("ban common use port stream")
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

func Start(devicenam,filter string) {
	newCaptureManager().Run(devicenam,filter)
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

	return fmt.Sprintf("dst 192.168.1.42 and port %d", 80)
}

// filter common use port
// 21: ftp, 23: telnet, 25: smtp, 139: samba
func isBanPort(srcPort int, dstPort int) bool {
	if srcPort < 200 && srcPort != 80 {
		return true
	}

	if dstPort < 200 && dstPort != 80 {
		return true
	}

	return false
}
