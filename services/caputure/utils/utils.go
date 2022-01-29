package utils

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"server-monitoring/services/caputure/flags"
	"strconv"
	"strings"


)

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func ShowAllInterfaces() {
	ifaces, _ := net.Interfaces()

	iplist := ""
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		ipAddrs := []string{}
		for _, addr := range addrs {
			var ip net.IP
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip = ipnet.IP
			} else if ipaddr, ok := addr.(*net.IPAddr); ok {
				ip = ipaddr.IP
			}
			if ip != nil && ip.To4() != nil && !ip.IsUnspecified() {
				ipstr := addr.String()
				idx := strings.Index(ipstr, "/")
				if idx >= 0 {
					ipstr = ipstr[:idx]
				}
				ipAddrs = append(ipAddrs, ipstr)
				name := iface.Name

				if !ip.IsLoopback() && strings.Contains(iface.Flags.String(), "up") {
					name = fmt.Sprintf("%s (%s)", iface.Name, iface.Flags.String())
				}
				iplist += fmt.Sprintf("%-7d %-50s %s\n", iface.Index, name, strings.Join(ipAddrs, ", "))
			}
		}
	}

	fmt.Printf("%-7s %-50s %s\n", "index", "interface name", "ip")
	fmt.Print(iplist)
}

func GetHostIp() string {
	ip := "127.0.0.1"

	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		ipnet := net.ParseIP(a.String())
		if ipnet != nil && !ipnet.IsLoopback() && !ipnet.IsUnspecified() {
			if ipnet.To4() != nil {
				ip = ipnet.String()
				break
			}
		}
	}

	return ip
}

func GetFirstInterface() (name string, ip string) {
	ifaces, _ := net.Interfaces()

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()

		ipV4 := false
		for _, addr := range addrs {
			var ip net.IP
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip = ipnet.IP
			} else if ipaddr, ok := addr.(*net.IPAddr); ok {
				ip = ipaddr.IP
			}
			if ip != nil && ip.To4() != nil && !ip.IsUnspecified() {
				ipstr := addr.String()
				idx := strings.Index(ipstr, "/")
				if idx >= 0 {
					ipstr = ipstr[:idx]
				}

				return iface.Name, ipstr
			}
		}
		if !ipV4 {
			continue
		}

	}

	return "", "0.0.0.0"
}

func GetDefaultInterface() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip = ipnet.IP
			} else if ipaddr, ok := addr.(*net.IPAddr); ok {
				ip = ipaddr.IP
			}
			if ip != nil && ip.To4() != nil && !ip.IsUnspecified() && !ip.IsLoopback() && strings.Contains(iface.Flags.String(), "up") {
				return iface.Name
			}
		}
	}

	return ""

}

func Debug(args ...interface{}) {
	if flags.Options.Verbose {
		log.Println(args...)
	}
}

func MustParseInt(s string) int {
	if len(s) == 0 {
		return 0
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func IP2int(ip string) uint32 {
	ipv4 := net.ParseIP(ip).To4()
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ipv4[12:16])
	}
	return binary.BigEndian.Uint32(ipv4)
}
