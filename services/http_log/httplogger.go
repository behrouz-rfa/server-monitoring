package http_log

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	// Default interface set to 'eth0'. iface indicates '-i' parameter.
	iface          string = "eth0"
	// Default mode is live capture. pcapFile indicates '-r' parameter.
	pcapFile string = ""
	// Default mode is not any specific string for payload. string_payload indicates '-s' parameter.
	string_payload string = ""
	// Default mode is dumping all captured packets. expression indicates '<expression>' parameter.
	// Expression parameter should be given in "" characters.
	expression string = ""
	snapshotLen int32  = 65535
	promiscuous bool   = true
	err         error
	timeout     time.Duration = -1 * time.Second
	handle      *pcap.Handle
)

func findMaxValue(s []int) int {
	// This function used for properly split user parameters.
	maxValue := 0
	for _, element := range s {
		if element > maxValue {
			maxValue = element
		}
	}

	return maxValue
}

func Run(string_payload,iface,pcapFile,expression string) {

	// Handle CLI arguments: (-i, -r, -s can be given in mixed order!)
	//arguments := os.Args
	// To keep readed index in arguments, so can find 'expression' part properly.
	//var readed_arguments_index []int
	//for index, element := range arguments {
	//	if element == "-i" {
	//		iface = arguments[index + 1]
	//		readed_arguments_index = append(readed_arguments_index, index + 1)
	//	} else if element == "-r" {
	//		pcapFile = arguments[index + 1]
	//		readed_arguments_index = append(readed_arguments_index, index + 1)
	//	} else if element == "-s" {
	//		string_payload = arguments[index + 1]
	//		readed_arguments_index = append(readed_arguments_index, index + 1)
	//	} else if string([]rune(element)[0]) == "-" {
	//		fmt.Println("UNKNOWN PARAMETER")
	//		os.Exit(3)
	//	}
	//}

	//if len(arguments) - 1 > findMaxValue(readed_arguments_index) {
	//	for i := findMaxValue(readed_arguments_index) + 1; i < len(arguments); i++ {
	//		expression += arguments[i] + " "
	//	}
	//}

	/*
	   // Debug line for control user parameters.
	   fmt.Println("-i -> ", iface)
	   fmt.Println("-r -> ", pcapFile)
	   fmt.Println("-s -> ", string_payload)
	   fmt.Println("expression -> ", expression)
	*/

	//if len(readed_arguments_index) == 0 {
	//	fmt.Println("Suggested Usage: sudo go run main.go -i <interface> <options> <expression>")
	//	fmt.Println("Options:")
	//	fmt.Println("\t -r: Read packets from <file> in tcpdump format")
	//	fmt.Println("\t -s: Keep only packets that contain <string> in their payload (In quotation marks).")
	//	fmt.Println("<expression> is a BPF filter that specifies which packets will be dumped. If no filter is given, packets seen on the interface (or contained in the trace) should be dumped. Otherwise, only packets matching <expression> should be dumped.")
	//	os.Exit(3)
	//}

	if pcapFile != "" {
		// If '-r' parameter is specified, then read pcapFile.
		handle, err = pcap.OpenOffline(pcapFile)
		if err != nil { log.Fatal(err) }
		defer handle.Close()
	} else {
		// If there is no '-r' parameter, then capture packets live.
		handle, err = pcap.OpenLive(iface, snapshotLen, promiscuous, timeout)
		if err != nil {log.Fatal(err) }
		defer handle.Close()
	}

	if expression != "" {
		// If some expression like "tcp and port 80" is specified, then set the BPFFilter:
		err = handle.SetBPFFilter(expression)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("BPF filter detected. Only capturing ", expression)
	}
	packetChannel = make(chan packetMessage, 10240)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet, string_payload)

	}

}
type logMessage struct {
	lType string
	msg   string
}
type packetMessage struct {
	data []byte
	ci   gopacket.CaptureInfo
}

var packetChannel chan packetMessage

func printPacketInfo(packet gopacket.Packet, payload_filter string) {

	 result :=""
	 println(result)
	// This is the decider for printing the packet.
	printThePacket := false

	// Firstly check if '-s' parameter has entered or not.
	if payload_filter != "" {
		app := packet.ApplicationLayer()
		if app != nil {
			packet_payload := app.Payload()
			if strings.Contains(string(packet_payload), payload_filter) {
				// If there is '-s' parameter and payload has that, then print.
				printThePacket = true
			}
		}
	} else {
		// If there is no '-s' parameter, then print all packets.
		printThePacket = true
	}

	if printThePacket {
		timestamp := packet.Metadata().Timestamp.Format("2006-01-02 15:04:05.999999")
		fmt.Print(timestamp, " ")

		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)

		if ethernetLayer != nil {

			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

			fmt.Print(ethernetPacket.SrcMAC, " -> ", ethernetPacket.DstMAC)
			fmt.Printf(" type 0x%x ", uint16(ethernetPacket.EthernetType))
			fmt.Print("len ", packet.Metadata().Length)
			fmt.Println()


			data := ethernetPacket.LayerContents()
			fmt.Println("datacontent")
			fmt.Println(string(data[0:5]))
			if bytes.Equal(data[0:5],[]byte("HTTP/")) {
				headers := strings.Split(string(data),"\n")
				statusLine  := strings.TrimSpace(headers[0])
				var contentLength string
				var contentType string
				for i := 1; i < len(headers); i++ {
					if len(headers[i]) > 16 && strings.EqualFold(headers[i][0:16], "Content-Length: ") {
						contentLength = strings.TrimSpace(headers[i][16:])
					} else if len(headers[i]) > 14 && strings.EqualFold(headers[i][0:14], "Content-Type: ") {
						contentType = strings.TrimSpace(headers[i][14:])
					} else if len(headers[i]) < 2 && strings.TrimSpace(headers[i]) == "" {
						break
					}
				}
				result =fmt.Sprintf("%s %s %s",statusLine,contentLength,contentType)

			}else if   bytes.Equal(data[0:4], []byte("GET ")) ||
				bytes.Equal(data[0:5], []byte("POST ")) ||
				bytes.Equal(data[0:5], []byte("HEAD ")) ||
				bytes.Equal(data[0:4], []byte("PUT ")) ||
				bytes.Equal(data[0:7], []byte("DELETE ")) ||
				bytes.Equal(data[0:8], []byte("OPTIONS ")) ||
				bytes.Equal(data[0:6], []byte("TRACE ")) ||
				bytes.Equal(data[0:6], []byte("PATCH ")) ||
				bytes.Equal(data[0:8], []byte("CONNECT ")) {
				headers := strings.Split(string(data), "\n")
				requestLine := strings.TrimSpace(headers[0])
				var hostData string
				var userAgent string

				for i := 1; i < len(headers); i++ {
					//log.Printf("%d %s", i, headers[i])
					if len(headers[i]) > 6 && strings.EqualFold(headers[i][0:6], "Host: ") {
						hostData = strings.TrimSpace(headers[i][7:])
					} else if len(headers[i]) > 12 && strings.EqualFold(headers[i][0:12], "User-Agent: ") {
						userAgent = strings.TrimSpace(headers[i][12:])
					} else if len(headers[i]) < 2 && strings.TrimSpace(headers[i]) == "" {
						break
					}
				}

				/*ch := bytes.IndexByte(data, byte('\r'))
				if ch == -1 {
					ch = bytes.IndexByte(data, byte('\n'))
					if ch == -1 {
						continue
					}
				}*/
				result = fmt.Sprintf("%s %s %s", requestLine, hostData, userAgent)

			}else {
				//contrin
			}

		}

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		tcpLayer := packet.Layer(layers.LayerTypeTCP) // It is 'nil' in case of packet is not TCP.
		udpLayer := packet.Layer(layers.LayerTypeUDP) // It is 'nil' in case of packet is not UDP.

		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)

			if tcpLayer != nil {
				// If packet is TCP:
				tcp, _ := tcpLayer.(*layers.TCP)

				srcPort, dstPort := uint16(tcp.SrcPort), uint16(tcp.DstPort)
				fmt.Print(ip.SrcIP, ":", srcPort, " -> ", ip.DstIP, ":", dstPort, " ")
				fmt.Print(ip.Protocol, " ")

				all_flags := [9]string{"FIN", "SYN", "RST", "PSH", "ACK", "URG", "ECE", "CWR", "NS"}
				all_flags_value := [9]bool{tcp.FIN, tcp.SYN, tcp.RST, tcp.PSH, tcp.ACK, tcp.URG, tcp.ECE, tcp.CWR, tcp.NS}
				for index, element := range all_flags_value {
					if element {
						fmt.Print(all_flags[index], " ")
					}
				}
				fmt.Println()

			} else if udpLayer != nil {
				// If packet is UDP
				udp, _ := udpLayer.(*layers.UDP)

				srcPort, dstPort := uint16(udp.SrcPort), uint16(udp.DstPort)
				fmt.Print(ip.SrcIP, ":", srcPort, " -> ", ip.DstIP, ":", dstPort, " ")
				fmt.Print(ip.Protocol)
				fmt.Println()
			} else if uint8(ip.Protocol) == 1 {
				// If packet is ICMP (ICMP's Protocol number is 1)
				fmt.Print(ip.SrcIP, " -> ", ip.DstIP, " ")
				fmt.Print(ip.Protocol)
				fmt.Println()
			}else {
				fmt.Print(ip.SrcIP, " -> ", ip.DstIP, " ")
				fmt.Print("OTHER")
				fmt.Println()
			}

		}

		app := packet.ApplicationLayer()
		if app != nil {
			packet_payload := app.Payload()
			// Hex dump of byte[] type payload
			fmt.Println(hex.Dump(packet_payload))
		} else {
			fmt.Println()
		}
	}
}