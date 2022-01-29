package packet

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"reflect"
	"server-monitoring/domain/nodes"
	"server-monitoring/shared/database"
	"time"
)

var (
	devuceName        = ""
	snpplan     int32 = 65535
	prmisc_bool       = false
	startInsert       = false
	err         error
	items                     = []nodes.Node{}
	timout      time.Duration = -1 * time.Second
	handle      *pcap.Handle
)

type httpStreamFactory struct{}

func Run() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	//IPAddress := net.ParseIP("192.168.1.42")

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
	handle, err = pcap.OpenLive(devuceName, snpplan, prmisc_bool, timout)
	if err != nil {
		log.Fatalln(err)
	}
	defer handle.Close()

	var filter string = "dst host 192.168.1.42"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatalln(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//fmt.Println("somone ping me!!")
		//fmt.Println("----------------")
		//fmt.Println(packet)
		analizePacket(packet)
	}
}

func analizePacket(packet gopacket.Packet) {
	timestamp := packet.Metadata().Timestamp.Format("2006-01-02 15:04:05.999999")
	fmt.Println(timestamp)
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)

	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Print(ethernetPacket.SrcMAC, " -> ", ethernetPacket.SrcMAC)
		fmt.Print("type")
		//fmt.Printf(" type 0x%x ", uint16(ethernetPacket.EthernetType))
		//fmt.Print("len ", packet.Metadata().Length)
		//fmt.Println()
	}

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP) // It is 'nil' in case of packet is not TCP.
	udpLayer := packet.Layer(layers.LayerTypeUDP) // It is 'nil' in case of packet is not UDP.
	icmLayer := packet.Layer(layers.LayerTypeICMPv4)

	ethernetFrame := ethernetLayer.(*layers.Ethernet)

	//applicationLayer := packet.ApplicationLayer()
	//if applicationLayer != nil {
	//	fmt.Println("Application layer/Payload found.")
	//	//fmt.Printf("%s\n", applicationLayer.Payload())
	//	fmt.Printf("%v\n", string(applicationLayer.Payload()))
	//
	//	// Search for a string inside the payload
	//	if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
	//		fmt.Println("HTTP found!")
	//		fmt.Println("layer content", string(applicationLayer.LayerContents()))
	//		fmt.Println("layer payload", string(applicationLayer.Payload()))
	//		fmt.Println("layer type", string(applicationLayer.LayerType()))
	//	}
	//}

	// It is 'nil' in case of packet is not UDP.
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)

		if tcpLayer != nil {

			// If packet is TCP:
			tcp, _ := tcpLayer.(*layers.TCP)
			if len(tcp.Payload) != 0 {
				reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
				httpReq, err := http.ReadRequest(reader)
				if err == nil {
					fmt.Printf("host: %s route:%s method:%s\n", httpReq.Host, httpReq.RequestURI, httpReq.Method)
					for s := range httpReq.Header {
						fmt.Println(s)
					}
				} else {
					fmt.Println("http error")
					fmt.Println(err.Error())
				}

			}
			srcPort, dstPort := uint16(tcp.SrcPort), uint16(tcp.DstPort)
			fmt.Print(ip.SrcIP, ":", srcPort, " -> ", ip.DstIP, ":", dstPort, " ")
			fmt.Print(ip.Protocol, " ")
			items = append(items, nodes.Node{
				ID:        primitive.NewObjectID(),
				SrcPort:   fmt.Sprintf("%s", srcPort),
				DstPort:   fmt.Sprintf("%s", dstPort),
				SrcIp:     fmt.Sprintf("%s", ip.SrcIP),
				DstIp:     fmt.Sprintf("%s", ip.DstIP),
				SrcMac:    fmt.Sprintf("%s", ethernetFrame.SrcMAC),
				DstMac:    fmt.Sprintf("%s", ethernetFrame.DstMAC),
				Protocol:  "TCP",
				Timestamp: time.Now().Unix(),
			})
			//applicationLayer := packet.ApplicationLayer()
			//if applicationLayer != nil {
			//	payload1 := string(applicationLayer.Payload())
			//	b := []byte(payload1)
			//	fmt.Printf("%s\n", hex.Dump(b))
			//}
			//fmt.Println()
			//all_flags := [9]string{"FIN", "SYN", "RST", "PSH", "ACK", "URG", "ECE", "CWR", "NS"}
			//all_flags_value := [9]bool{tcp.FIN, tcp.SYN, tcp.RST, tcp.PSH, tcp.ACK, tcp.URG, tcp.ECE, tcp.CWR, tcp.NS}
			//for index, element := range all_flags_value {
			//	if element {
			//		fmt.Print(all_flags[index], " ")
			//	}
			//}
			fmt.Println()

		} else if udpLayer != nil {
			// If packet is UDP
			udp, _ := udpLayer.(*layers.UDP)

			srcPort, dstPort := uint16(udp.SrcPort), uint16(udp.DstPort)
			fmt.Print(ip.SrcIP, ":", srcPort, " -> ", ip.DstIP, ":", dstPort, " ")
			fmt.Print(ip.Protocol)

			fmt.Println("UDP data to string: ", udp.LayerPayload())

			if len(udp.Payload) != 0 {
				reader := bufio.NewReader(bytes.NewReader(udp.Payload))
				httpReq, err := http.ReadRequest(reader)
				if err == nil {
					fmt.Printf("host: %s route:%s method:%s\n", httpReq.Host, httpReq.RequestURI, httpReq.Method)
					for s := range httpReq.Header {
						fmt.Println(s)
					}
				}
			}
			items = append(items, nodes.Node{
				ID:        primitive.NewObjectID(),
				SrcPort:   fmt.Sprintf("%s", srcPort),
				DstPort:   fmt.Sprintf("%s", dstPort),
				SrcIp:     fmt.Sprintf("%s", ip.SrcIP),
				DstIp:     fmt.Sprintf("%s", ip.DstIP),
				SrcMac:    fmt.Sprintf("%s", ethernetFrame.SrcMAC),
				DstMac:    fmt.Sprintf("%s", ethernetFrame.DstMAC),
				Protocol:  "UDP",
				Timestamp: time.Now().Unix(),
			})
			fmt.Println()
		} else if uint8(ip.Protocol) == 1 {
			// If packet is ICMP (ICMP's Protocol number is 1)
			icmp_packet := icmLayer.(*layers.ICMPv4)
			if icmp_packet.TypeCode.String() == "EchoRequest" {
				if len(icmp_packet.Payload) > 0 {
					fmt.Println("Info: Echorequest receoved")
				} else {
					fmt.Println("Warning: empty echorequest received")
					ehternetFramCopy := *ethernetFrame
					ippacketCopy := *ip
					ICMpACKETcIOPY := *icmp_packet
					ehternetFramCopy.SrcMAC = ethernetFrame.DstMAC
					ehternetFramCopy.DstMAC = ethernetFrame.SrcMAC

					ippacketCopy.SrcIP = ip.DstIP
					ippacketCopy.DstIP = ip.SrcIP
					ICMpACKETcIOPY.TypeCode = layers.ICMPv4TypeEchoReply
					var buffer gopacket.SerializeBuffer
					var options gopacket.SerializeOptions
					options.ComputeChecksums = true
					gopacket.SerializeLayers(buffer, options,
						&ehternetFramCopy, &ippacketCopy, &ICMpACKETcIOPY, gopacket.Payload(ICMpACKETcIOPY.Payload),
					)
					new_message := buffer.Bytes()
					err := handle.WritePacketData(new_message)
					if err != nil {
						log.Fatal(err)
					}

				}

				items = append(items, nodes.Node{
					ID:        primitive.NewObjectID(),
					SrcPort:   "",
					DstPort:   "",
					SrcIp:     fmt.Sprintf("%s", ip.SrcIP),
					DstIp:     fmt.Sprintf("%s", ip.DstIP),
					SrcMac:    fmt.Sprintf("%s", ethernetFrame.SrcMAC),
					DstMac:    fmt.Sprintf("%s", ethernetFrame.DstMAC),
					Protocol:  "UDP",
					Timestamp: time.Now().Unix(),
				})
			}
			//fmt.Print(ip.SrcIP, " -> ", ip.DstIP, " ")
			//fmt.Print(ip.Protocol)
			//fmt.Println()
			//fmt.Println("ICMP code: ",icmp_packet.TypeCode)
			//fmt.Println("ICMP sequence number: ",icmp_packet.Seq)
			//fmt.Println("ICMP data length: ",len(icmp_packet.Payload))
			//fmt.Println("ICMP data : ",icmp_packet.Payload)
			//fmt.Println("Payload data to string: ",string(icmp_packet.Payload))

		} else {
			//fmt.Print(ip.SrcIP, " -> ", ip.DstIP, " ")
			//fmt.Print("OTHER")
			//fmt.Println()
		}

		if !startInsert && len(items) > 10 {
			startInsert = !startInsert
			if err := database.InsertLogs(items); err == nil {
				//items = []nodes.Node{}
			}
		}

	}
}
