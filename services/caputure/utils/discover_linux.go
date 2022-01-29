package utils

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func DiscoverServices() map[int]Service {
	cmd := "netstat -tnpl | grep -E 'mysqld|redis-server|memcached|mongos|nutcracker' | awk '{print $4,$7}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Print("DiscoverServices error.")
		log.Fatal(err)
	}

	services := make(map[int]Service)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		reg := regexp.MustCompile(`\s+`)
		cols := reg.Split(strings.TrimSpace(line), -1)

		// ignore ipv6
		if len(cols) < 2 || strings.HasPrefix(cols[0], ":::") {
			continue
		}

		arr := strings.Split(cols[0], ":")
		port, _ := strconv.Atoi(arr[1])
		arr = strings.Split(cols[1], "/")
		pid, _ := strconv.Atoi(arr[0])
		exec := arr[1]

		switch exec {
		case "redis-server":
			services[port] = Service{Port: port, Type: Service_Type_Redis, Pid: pid}
		case "memcached":
			services[port] = Service{Port: port, Type: Service_Type_Memcache, Pid: pid}
		case "mongod":
			services[port] = Service{Port: port, Type: Service_Type_Mongodb, Pid: pid}
		case "mysqld":
			services[port] = Service{Port: port, Type: Service_Type_Mysql, Pid: pid}
		case "nutcracker":
			services[port] = Service{Port: port, Type: Service_Type_Twemproxy, Pid: pid}
		}

	}

	return services
}
